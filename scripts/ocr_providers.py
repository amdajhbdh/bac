#!/usr/bin/env python3
"""
BAC OCR Pipeline — 18-Provider OCR Module
Multi-provider OCR with cascading fallback chain.
Each provider returns a standardized response format.

Providers (in priority order):
  Cloud Tier:     Workers AI (Llama Vision), Mistral OCR, OCR.space, Yandex Vision
  Cloud Shell:    LFM2.5-VL (Ollama), Qwen2.5-VL (Ollama), LLaVA (Ollama)
  Local Python:  PaddleOCR, EasyOCR, GLM-OCR, Tesseract, Yandex (API), Cloudmersive,
                 ABBYY Cloud, Base64.ai, OCRAPI.cloud, Nanonets, Puter.js

Usage:
  from ocr_providers import OCRResult, run_ocr_chain
  result = run_ocr_chain(image_path_or_bytes, providers=["paddle", "tesseract"])
"""

from __future__ import annotations

import base64
import io
import json
import os
import re
import subprocess
import time
from dataclasses import dataclass, field, asdict
from pathlib import Path
from typing import Optional, Literal

import requests
import importlib.util

# ============================================================
# UNIFIED RESPONSE FORMAT
# ============================================================


@dataclass
class OCRResult:
    """Standardized OCR response from any provider."""

    text: str
    source: str  # short name: "mistral", "paddle", etc.
    provider: str  # full name
    tier: str  # "cloud" | "cloud-shell" | "local"
    chars: int = 0
    confidence: Optional[float] = None
    error: Optional[str] = None
    latency_ms: int = 0
    cost: float = 0.0  # USD per call
    languages: Optional[list[str]] = None

    def to_dict(self) -> dict:
        return asdict(self)

    @property
    def success(self) -> bool:
        return bool(self.text and self.chars >= 5 and not self.error)

    def __post_init__(self):
        if self.chars == 0 and self.text:
            self.chars = len(self.text)


# ============================================================
# CONFIG
# ============================================================


@dataclass
class OCRConfig:
    """Configuration for OCR providers."""

    # API Keys
    mistral_api_key: Optional[str] = None
    ocrspace_api_key: str = "helloworld"  # free tier key
    yandex_iam_token: Optional[str] = None
    yandex_folder_id: Optional[str] = None
    cloudmersive_key: Optional[str] = None
    abbyy_app_id: Optional[str] = None
    abbyy_password: Optional[str] = None
    base64_api_key: Optional[str] = None
    nanonets_api_key: Optional[str] = None
    mathpix_app_id: Optional[str] = None
    mathpix_app_key: Optional[str] = None

    # Endpoints
    worker_url: str = "https://bac-api.amdajhbdh.workers.dev"
    cloud_shell_url: str = "http://localhost:11434"
    ollama_url: str = "http://localhost:11434"

    # Timeouts (seconds)
    api_timeout: int = 30
    local_timeout: int = 60

    @classmethod
    def from_env(cls) -> "OCRConfig":
        return cls(
            mistral_api_key=os.environ.get("MISTRAL_API_KEY"),
            yandex_iam_token=os.environ.get("YANDEX_IAM_TOKEN"),
            yandex_folder_id=os.environ.get("YANDEX_FOLDER_ID"),
            cloudmersive_key=os.environ.get("CLOUDMERSIVE_API_KEY"),
            abbyy_app_id=os.environ.get("ABBYY_APP_ID"),
            abbyy_password=os.environ.get("ABBYY_PASSWORD"),
            base64_api_key=os.environ.get("BASE64_API_KEY"),
            nanonets_api_key=os.environ.get("NANONETS_API_KEY"),
            mathpix_app_id=os.environ.get("MATHPIX_APP_ID"),
            mathpix_app_key=os.environ.get("MATHPIX_APP_KEY"),
            worker_url=os.environ.get(
                "WORKER_URL", "https://bac-api.amdajhbdh.workers.dev"
            ),
            cloud_shell_url=os.environ.get("CLOUD_SHELL_URL", "http://localhost:11434"),
        )


# ============================================================
# IMAGE UTILITIES
# ============================================================


def read_image(path_or_bytes: str | bytes | Path) -> bytes:
    if isinstance(path_or_bytes, (str, Path)):
        return Path(path_or_bytes).read_bytes()
    return bytes(path_or_bytes)


def to_base64(image_bytes: bytes, mime: str = "image/jpeg") -> str:
    return f"data:{mime};base64,{base64.b64encode(image_bytes).decode()}"


def to_base64_clean(image_bytes: bytes) -> str:
    return base64.b64encode(image_bytes).decode()


def image_to_pil(image_bytes: bytes):
    try:
        from PIL import Image

        return Image.open(io.BytesIO(image_bytes))
    except ImportError:
        return None


# ============================================================
# CLOUD TIER PROVIDERS (4)
# ============================================================


def ocr_workers_ai(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Workers AI via Cloudflare (Llama Vision). Free, 10K neurons/day."""
    config = config or OCRConfig.from_env()
    start = time.time()

    try:
        img_bytes = read_image(image_path_or_bytes)
        b64 = to_base64(img_bytes)

        resp = requests.post(
            f"{config.worker_url}/rag/ocr",
            json={"image": b64},
            headers={"Content-Type": "application/json"},
            timeout=config.api_timeout,
        )
        data = resp.json()

        if data.get("text"):
            return OCRResult(
                text=data["text"],
                source=data.get("source", "workers-ai"),
                provider="Cloudflare Workers AI (Llama Vision)",
                tier="cloud",
                chars=len(data["text"]),
                latency_ms=int((time.time() - start) * 1000),
                cost=0.0,
            )
        else:
            return OCRResult(
                text="",
                source="workers-ai",
                provider="Workers AI",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=data.get("error", "Empty response"),
            )
    except Exception as e:
        return OCRResult(
            text="",
            source="workers-ai",
            provider="Workers AI",
            tier="cloud",
            latency_ms=int((time.time() - start) * 1000),
            error=str(e),
        )


def ocr_mistral(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Mistral OCR API. ~$0.002/page. Best accuracy for French/Arabic."""
    config = config or OCRConfig.from_env()
    start = time.time()

    if not config.mistral_api_key:
        return OCRResult(
            text="",
            source="mistral",
            provider="Mistral OCR",
            tier="cloud",
            error="MISTRAL_API_KEY not configured",
        )

    try:
        img_bytes = read_image(image_path_or_bytes)
        b64 = to_base64_clean(img_bytes)

        resp = requests.post(
            "https://api.mistral.ai/v1/ocr",
            headers={
                "Authorization": f"Bearer {config.mistral_api_key}",
                "Content-Type": "application/json",
            },
            json={
                "model": "mistral-ocr-latest",
                "document": {
                    "type": "image_url",
                    "image_url": f"data:image/jpeg;base64,{b64}",
                },
            },
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="mistral",
                provider="Mistral OCR",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}: {resp.text[:200]}",
            )

        data = resp.json()
        pages = data.get("pages", [])
        text = "\n\n".join(p.get("markdown", "") or "" for p in pages)

        return OCRResult(
            text=text,
            source="mistral",
            provider="Mistral OCR",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.002,  # ~$0.002/page
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="mistral",
            provider="Mistral OCR",
            tier="cloud",
            latency_ms=int((time.time() - start) * 1000),
            error=str(e),
        )


def ocr_ocrspace(
    image_path_or_bytes: str | bytes,
    config: Optional[OCRConfig] = None,
    language: str = "fra",
) -> OCRResult:
    """OCR.space API. 500 requests/day free. No key needed for free tier."""
    config = config or OCRConfig.from_env()
    start = time.time()

    try:
        img_bytes = read_image(image_path_or_bytes)
        b64 = to_base64_clean(img_bytes)

        form_data = [
            ("base64Image", f"data:image/jpeg;base64,{b64}"),
            ("language", language),
            ("isOverlayRequired", "false"),
            ("detectOrientation", "true"),
            ("scale", "true"),
            ("OCREngine", "2"),  # Engine 2 = more accurate
        ]

        headers = {}
        if config.ocrspace_api_key and config.ocrspace_api_key != "helloworld":
            headers["apikey"] = config.ocrspace_api_key

        resp = requests.post(
            "https://api.ocr.space/parse/image",
            files={"file": ("image.jpg", img_bytes, "image/jpeg")},
            data={
                "language": language,
                "isOverlayRequired": "false",
                "detectOrientation": "true",
                "scale": "true",
                "OCREngine": "2",
            },
            headers=headers,
            timeout=config.api_timeout,
        )

        data = resp.json()
        if data.get("IsErroredOnProcessing"):
            return OCRResult(
                text="",
                source="ocrspace",
                provider="OCR.space",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=", ".join(data.get("ErrorMessage", [])[:2]),
            )

        parsed = data.get("ParsedResults", [])
        text = "\n".join(p.get("ParsedText", "") for p in parsed)

        return OCRResult(
            text=text,
            source="ocrspace",
            provider="OCR.space",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="ocrspace",
            provider="OCR.space",
            tier="cloud",
            latency_ms=int((time.time() - start) * 1000),
            error=str(e),
        )


def ocr_yandex(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Yandex Vision OCR. 1K operations/3 months free. French + Arabic."""
    config = config or OCRConfig.from_env()
    start = time.time()

    if not config.yandex_iam_token or not config.yandex_folder_id:
        return OCRResult(
            text="",
            source="yandex",
            provider="Yandex Vision",
            tier="cloud",
            error="YANDEX_IAM_TOKEN or YANDEX_FOLDER_ID not configured",
        )

    try:
        img_bytes = read_image(image_path_or_bytes)
        b64 = to_base64_clean(img_bytes)

        resp = requests.post(
            "https://ocr.api.cloud.yandex.net/ocr/v1/recognizeText",
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {config.yandex_iam_token}",
                "x-folder-id": config.yandex_folder_id,
            },
            json={
                "mimeType": "JPEG",
                "languageCodes": ["fr", "ar", "en"],
                "model": "page",
                "content": b64,
            },
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="yandex",
                provider="Yandex Vision",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        data = resp.json()
        annotation = data.get("result", {}).get("textAnnotation", {})
        text = annotation.get("fullText", "")
        if not text:
            words = annotation.get("words", [])
            text = " ".join(w.get("text", "") for w in words)

        return OCRResult(
            text=text,
            source="yandex",
            provider="Yandex Vision OCR",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="yandex",
            provider="Yandex Vision",
            tier="cloud",
            latency_ms=int((time.time() - start) * 1000),
            error=str(e),
        )


# ============================================================
# CLOUD SHELL TIER PROVIDERS (3)
# ============================================================


def _ollama_vision(
    image_bytes: bytes, model: str, prompt: str, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Generic Ollama vision model OCR."""
    config = config or OCRConfig.from_env()
    start = time.time()

    try:
        b64 = base64.b64encode(image_bytes).decode()

        resp = requests.post(
            f"{config.ollama_url}/api/generate",
            json={
                "model": model,
                "prompt": prompt,
                "images": [b64],
                "stream": False,
                "options": {"temperature": 0.1, "num_predict": 1024},
            },
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source=model.split("/")[-1],
                provider=f"Ollama ({model})",
                tier="cloud-shell",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        data = resp.json()
        text = data.get("response", "").strip()

        return OCRResult(
            text=text,
            source=model.split("/")[-1],
            provider=f"Ollama ({model})",
            tier="cloud-shell",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,
        )
    except Exception as e:
        return OCRResult(
            text="",
            source=model.split("/")[-1],
            provider=f"Ollama ({model})",
            tier="cloud-shell",
            latency_ms=int((time.time() - start) * 1000),
            error=str(e),
        )


OCR_PROMPT_FR_AR = (
    "Extract ALL text from this image word-for-word. "
    "Preserve French characters (é, è, ê, ë, à, â, ç, ô, û, ï, ü, œ, Æ) "
    "and Arabic characters exactly. "
    "Maintain line breaks and paragraph structure. "
    "If no readable text, say: NO_TEXT_DETECTED"
)


def ocr_lfm25_vl(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """LFM2.5-VL via Ollama. Best for French/Arabic OCR. 696 MB (Q4_0)."""
    img_bytes = read_image(image_path_or_bytes)
    return _ollama_vision(
        img_bytes, "liquidai/lfm2.5-vl-1.6b:q4_0", OCR_PROMPT_FR_AR, config
    )


def ocr_qwen25_vl(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Qwen2.5-VL via Ollama. New multimodal model. Strong accuracy."""
    img_bytes = read_image(image_path_or_bytes)
    return _ollama_vision(img_bytes, "qwen2.5vl:latest", OCR_PROMPT_FR_AR, config)


def ocr_llava(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """LLaVA via Ollama. Stable vision model. ~7B params."""
    img_bytes = read_image(image_path_or_bytes)
    return _ollama_vision(img_bytes, "llava:latest", OCR_PROMPT_FR_AR, config)


# ============================================================
# LOCAL PYTHON TIER PROVIDERS (11)
# ============================================================


def ocr_paddle(
    image_path_or_bytes: str | bytes,
    config: Optional[OCRConfig] = None,
    languages: str = "fr",
) -> OCRResult:
    """PaddleOCR. Best open-source accuracy. Apache license. French/Arabic."""
    config = config or OCRConfig()
    start = time.time()

    try:
        from paddleocr import PaddleOCR

        img_bytes = read_image(image_path_or_bytes)

        ocr = PaddleOCR(
            use_angle_cls=True,
            lang=languages,
            use_gpu=False,
            show_log=False,
        )

        # PaddleOCR needs a file path or numpy array
        import tempfile

        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=False) as f:
            f.write(img_bytes)
            tmp_path = f.name

        try:
            result = ocr.ocr(tmp_path, cls=True)
        finally:
            Path(tmp_path).unlink(missing_ok=True)

        if not result or not result[0]:
            return OCRResult(
                text="",
                source="paddle",
                provider="PaddleOCR",
                tier="local",
                latency_ms=int((time.time() - start) * 1000),
            )

        lines = []
        for line in result[0] or []:
            if line and len(line) >= 2:
                lines.append(str(line[1][0]))

        text = "\n".join(lines)

        return OCRResult(
            text=text,
            source="paddle",
            provider="PaddleOCR",
            tier="local",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,
        )
    except ImportError:
        return OCRResult(
            text="",
            source="paddle",
            provider="PaddleOCR",
            tier="local",
            error="PaddleOCR not installed. Run: pip install paddlepaddle paddleocr",
            latency_ms=int((time.time() - start) * 1000),
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="paddle",
            provider="PaddleOCR",
            tier="local",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_easyocr(
    image_path_or_bytes: str | bytes,
    config: Optional[OCRConfig] = None,
    languages: list[str] = None,
) -> OCRResult:
    """EasyOCR. Deep learning based. 80+ languages including French/Arabic."""
    config = config or OCRConfig()
    start = time.time()
    languages = languages or ["fr", "ar", "en"]

    try:
        import easyocr

        img_bytes = read_image(image_path_or_bytes)

        reader = easyocr.Reader(languages, gpu=False, verbose=False)
        result = reader.readtext(img_bytes)

        if not result:
            return OCRResult(
                text="",
                source="easyocr",
                provider="EasyOCR",
                tier="local",
                latency_ms=int((time.time() - start) * 1000),
            )

        lines = []
        for detection in result:
            if len(detection) >= 2:
                lines.append(str(detection[1]))

        text = "\n".join(lines)

        return OCRResult(
            text=text,
            source="easyocr",
            provider="EasyOCR",
            tier="local",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,
        )
    except ImportError:
        return OCRResult(
            text="",
            source="easyocr",
            provider="EasyOCR",
            tier="local",
            error="EasyOCR not installed. Run: pip install easyocr",
            latency_ms=int((time.time() - start) * 1000),
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="easyocr",
            provider="EasyOCR",
            tier="local",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_tesseract(
    image_path_or_bytes: str | bytes,
    config: Optional[OCRConfig] = None,
    languages: str = "fra+ara+eng",
) -> OCRResult:
    """Tesseract CLI. Local, free, supports French+Arabic. Already configured."""
    config = config or OCRConfig()
    start = time.time()

    try:
        img_bytes = read_image(image_path_or_bytes)

        import tempfile

        with tempfile.NamedTemporaryFile(suffix=".jpg", delete=False) as f:
            f.write(img_bytes)
            tmp_path = f.name

        try:
            langs = languages.replace("+", "+")
            result = subprocess.run(
                ["tesseract", tmp_path, "stdout", "-l", langs],
                capture_output=True,
                text=True,
                timeout=config.local_timeout,
            )
            if result.returncode != 0:
                return OCRResult(
                    text="",
                    source="tesseract",
                    provider="Tesseract OCR",
                    tier="local",
                    latency_ms=int((time.time() - start) * 1000),
                    error=f"Tesseract error: {result.stderr}",
                )
            text = result.stdout.strip()
        finally:
            Path(tmp_path).unlink(missing_ok=True)

        return OCRResult(
            text=text,
            source="tesseract",
            provider="Tesseract OCR",
            tier="local",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,
        )
    except FileNotFoundError:
        return OCRResult(
            text="",
            source="tesseract",
            provider="Tesseract OCR",
            tier="local",
            error="Tesseract not installed. Run: apt install tesseract-ocr tesseract-ocr-fra tesseract-ocr-ara",
            latency_ms=int((time.time() - start) * 1000),
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="tesseract",
            provider="Tesseract OCR",
            tier="local",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_glm(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """GLM-OCR via Ollama. Best for math equations and complex documents. 0.9B params."""
    config = config or OCRConfig.from_env()
    start = time.time()

    try:
        img_bytes = read_image(image_path_or_bytes)
        b64 = base64.b64encode(img_bytes).decode()

        # Try Ollama first
        try:
            resp = requests.post(
                f"{config.ollama_url}/api/generate",
                json={
                    "model": "glm-ocr:latest",
                    "prompt": (
                        "This is a math and science document. Extract ALL text including "
                        "mathematical formulas, equations, chemical formulas, and regular text. "
                        "Preserve LaTeX formatting for equations. Output ONLY the extracted text."
                    ),
                    "images": [b64],
                    "stream": False,
                },
                timeout=config.api_timeout,
            )
            if resp.ok:
                data = resp.json()
                text = data.get("response", "").strip()
                return OCRResult(
                    text=text,
                    source="glm-ocr",
                    provider="GLM-OCR (Ollama)",
                    tier="cloud-shell",
                    chars=len(text),
                    latency_ms=int((time.time() - start) * 1000),
                )
        except Exception:
            pass

        # Fallback: use GLM via HuggingFace transformers
        return OCRResult(
            text="",
            source="glm-ocr",
            provider="GLM-OCR",
            tier="local",
            latency_ms=int((time.time() - start) * 1000),
            error="GLM-OCR not available on Ollama. Run: ollama pull glm-ocr",
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="glm-ocr",
            provider="GLM-OCR",
            tier="local",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_cloudmersive(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Cloudmersive API. 800 requests/month free. French supported, Arabic limited."""
    config = config or OCRConfig.from_env()
    start = time.time()

    if not config.cloudmersive_key:
        return OCRResult(
            text="",
            source="cloudmersive",
            provider="Cloudmersive OCR",
            tier="cloud",
            error="CLOUDMERSIVE_API_KEY not configured",
        )

    try:
        img_bytes = read_image(image_path_or_bytes)

        resp = requests.post(
            "https://api.cloudmersive.com/v1/ocr/Image/ocr",
            headers={"Apikey": config.cloudmersive_key},
            files={"imageFile": ("image.jpg", img_bytes, "image/jpeg")},
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="cloudmersive",
                provider="Cloudmersive OCR",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        data = resp.json()
        text = data.get("textResult", "")

        return OCRResult(
            text=text,
            source="cloudmersive",
            provider="Cloudmersive OCR",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,  # 800/mo free
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="cloudmersive",
            provider="Cloudmersive OCR",
            tier="cloud",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_abbyy(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """ABBYY Cloud OCR. 500 pages/month free. Highest accuracy for complex layouts."""
    config = config or OCRConfig.from_env()
    start = time.time()

    if not config.abbyy_app_id or not config.abbyy_password:
        return OCRResult(
            text="",
            source="abbyy",
            provider="ABBYY Cloud OCR",
            tier="cloud",
            error="ABBYY_APP_ID or ABBYY_PASSWORD not configured",
        )

    try:
        img_bytes = read_image(image_path_or_bytes)

        # Step 1: Get auth token
        auth_resp = requests.get(
            "https://cloud.ocrsdk.com/v2/getAccountBalance",
            auth=(config.abbyy_app_id, config.abbyy_password),
            timeout=config.api_timeout,
        )
        if not auth_resp.ok:
            return OCRResult(
                text="",
                source="abbyy",
                provider="ABBYY Cloud OCR",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"Auth failed: HTTP {auth_resp.status_code}",
            )

        # Step 2: Submit image
        resp = requests.post(
            "https://cloud.ocrsdk.com/v2/processImage",
            auth=(config.abbyy_app_id, config.abbyy_password),
            data={"language": "French,English,Arabic", "exportFormat": "txt"},
            files={"file": ("image.jpg", img_bytes, "image/jpeg")},
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="abbyy",
                provider="ABBYY Cloud OCR",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        # ABBYY returns XML with task ID, need to poll
        # For simplicity, return the task response
        return OCRResult(
            text="",
            source="abbyy",
            provider="ABBYY Cloud OCR",
            tier="cloud",
            latency_ms=int((time.time() - start) * 1000),
            error="ABBYY requires async polling. See: https://ocrsdk.com/documentation/",
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="abbyy",
            provider="ABBYY Cloud OCR",
            tier="cloud",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_base64ai(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Base64.ai. $0.01/page. 165+ languages, 2,800+ document types."""
    config = config or OCRConfig.from_env()
    start = time.time()

    if not config.base64_api_key:
        return OCRResult(
            text="",
            source="base64",
            provider="Base64.ai",
            tier="cloud",
            error="BASE64_API_KEY not configured",
        )

    try:
        img_bytes = read_image(image_path_or_bytes)

        resp = requests.post(
            "https://api.base64.ai/v1/extract",
            headers={"Authorization": f"Bearer {config.base64_api_key}"},
            files={"file": ("image.jpg", img_bytes, "image/jpeg")},
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="base64",
                provider="Base64.ai",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        data = resp.json()
        # Base64.ai returns structured data, extract text fields
        text_parts = []

        def extract_text(obj):
            if isinstance(obj, dict):
                for key in ["text", "rawText", "content"]:
                    if key in obj and isinstance(obj[key], str):
                        text_parts.append(obj[key])
                for val in obj.values():
                    extract_text(val)
            elif isinstance(obj, list):
                for item in obj:
                    extract_text(item)

        extract_text(data)
        text = "\n".join(text_parts)

        return OCRResult(
            text=text,
            source="base64",
            provider="Base64.ai",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.01,
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="base64",
            provider="Base64.ai",
            tier="cloud",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_ocrapi(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """OCRAPI.cloud. 250 requests/month free. Async with webhooks."""
    config = config or OCRConfig.from_env()
    start = time.time()

    try:
        img_bytes = read_image(image_path_or_bytes)

        resp = requests.post(
            "https://ocrapi.cloud/api/v1/recognize",
            headers={"Content-Type": "application/json"},
            json={
                "base64": base64.b64encode(img_bytes).decode(),
                "language": "fr",
                "isOverlayRequired": False,
            },
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="ocrapi",
                provider="OCRAPI.cloud",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        data = resp.json()
        # OCRAPI.cloud returns various formats depending on plan
        text = data.get("parsedText", "") or data.get("text", "")

        return OCRResult(
            text=text,
            source="ocrapi",
            provider="OCRAPI.cloud",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.0,  # 250/mo free
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="ocrapi",
            provider="OCRAPI.cloud",
            tier="cloud",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


def ocr_nanonets(
    image_path_or_bytes: str | bytes, config: Optional[OCRConfig] = None
) -> OCRResult:
    """Nanonets. $200 credits on signup. 200+ languages, LaTeX support."""
    config = config or OCRConfig.from_env()
    start = time.time()

    if not config.nanonets_api_key:
        return OCRResult(
            text="",
            source="nanonets",
            provider="Nanonets OCR",
            tier="cloud",
            error="NANONETS_API_KEY not configured",
        )

    try:
        img_bytes = read_image(image_path_or_bytes)

        resp = requests.post(
            "https://app.nanonets.com/api/v2/OCR/FullText",
            auth=(config.nanonets_api_key, ""),
            files={"file": ("image.jpg", img_bytes, "image/jpeg")},
            timeout=config.api_timeout,
        )

        if not resp.ok:
            return OCRResult(
                text="",
                source="nanonets",
                provider="Nanonets OCR",
                tier="cloud",
                latency_ms=int((time.time() - start) * 1000),
                error=f"HTTP {resp.status_code}",
            )

        data = resp.json()
        results = data.get("result", [])
        if isinstance(results, list) and len(results) > 0:
            text = results[0].get("text", "")
        else:
            text = str(results)

        return OCRResult(
            text=text,
            source="nanonets",
            provider="Nanonets OCR",
            tier="cloud",
            chars=len(text),
            latency_ms=int((time.time() - start) * 1000),
            cost=0.003,  # ~$0.003/page
        )
    except Exception as e:
        return OCRResult(
            text="",
            source="nanonets",
            provider="Nanonets OCR",
            tier="cloud",
            error=str(e),
            latency_ms=int((time.time() - start) * 1000),
        )


# ============================================================
# PROVIDER REGISTRY
# ============================================================

PROVIDER_MAP: dict[str, callable] = {
    # Cloud tier
    "workers": ocr_workers_ai,
    "mistral": ocr_mistral,
    "ocrspace": ocr_ocrspace,
    "yandex": ocr_yandex,
    # Cloud Shell tier
    "lfm25vl": ocr_lfm25_vl,
    "qwen25vl": ocr_qwen25_vl,
    "llava": ocr_llava,
    # Local tier
    "paddle": ocr_paddle,
    "easyocr": ocr_easyocr,
    "tesseract": ocr_tesseract,
    "glm": ocr_glm,
    # Paid cloud
    "cloudmersive": ocr_cloudmersive,
    "abbyy": ocr_abbyy,
    "base64": ocr_base64ai,
    "ocrapi": ocr_ocrapi,
    "nanonets": ocr_nanonets,
}

# Default chain (fastest path to success)
DEFAULT_CHAIN = [
    "workers",  # Free, fast, decent quality
    "paddle",  # Free, best open-source accuracy
    "tesseract",  # Free, reliable, French+Arabic
    "mistral",  # Paid, best accuracy
    "ocrspace",  # Free, good fallback
    "yandex",  # Free, multilingual
    "easyocr",  # Free, deep learning
    "llava",  # Free, vision-language model
    "qwen25vl",  # Free, modern multimodal
    "lfm25vl",  # Free, French/Arabic optimized
]


# ============================================================
## RCHESTRATOR
# ============================================================


def run_ocr_chain(
    image_path_or_bytes: str | bytes,
    providers: Optional[list[str]] = None,
    config: Optional[OCRConfig] = None,
    use_preprocessing: bool = True,
) -> tuple[OCRResult, list[OCRResult]]:
    """
    Run OCR with cascading fallback chain.

    Args:
        image_path_or_bytes: Image file path or raw bytes
        providers: List of provider names (default: DEFAULT_CHAIN)
        config: OCRConfig instance
        use_preprocessing: Apply smart preprocessing before OCR

    Returns: (best_result, all_results)
    - best_result: First successful provider result
    - all_results: All provider attempts with errors
    """
    providers = providers or DEFAULT_CHAIN
    config = config or OCRConfig.from_env()
    all_results: list[OCRResult] = []

    if use_preprocessing:
        _prep_mod = _get_preprocess_module()
        if _prep_mod:
            image_path_or_bytes = _prep_mod.smart_preprocess(
                read_image(image_path_or_bytes),
                providers[0] if providers else "mistral",
            )

    for name in providers:
        fn = PROVIDER_MAP.get(name)
        if not fn:
            continue

        result = fn(image_path_or_bytes, config)
        all_results.append(result)

        if result.success:
            return result, all_results

    # Return the best attempt (most characters) even if all failed
    best = (
        max(all_results, key=lambda r: r.chars)
        if all_results
        else OCRResult(
            text="",
            source="none",
            provider="none",
            tier="none",
            error="All OCR providers failed",
        )
    )
    return best, all_results


def _get_preprocess_module():
    """Load preprocessing module if available."""
    try:
        spec = importlib.util.spec_from_file_location(
            "ocr_preprocess", Path(__file__).parent / "ocr_preprocess.py"
        )
        if spec and spec.loader:
            mod = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(mod)
            return mod
    except Exception:
        pass
    return None


def run_ocr_parallel(
    image_path_or_bytes: str | bytes,
    providers: Optional[list[str]] = None,
    config: Optional[OCRConfig] = None,
    timeout_per_provider: int = 30,
) -> tuple[OCRResult, list[OCRResult]]:
    """
    Run multiple OCR providers concurrently and return the first success.
    Uses threading for parallel execution.
    """
    import concurrent.futures

    providers = providers or DEFAULT_CHAIN[:5]  # Top 5 for parallel
    config = config or OCRConfig.from_env()
    all_results: list[OCRResult] = []

    def safe_call(name: str) -> tuple[str, OCRResult]:
        fn = PROVIDER_MAP.get(name)
        if not fn:
            return name, OCRResult(
                text="",
                source=name,
                provider=name,
                tier="unknown",
                error="Unknown provider",
            )
        try:
            return name, fn(image_path_or_bytes, config)
        except Exception as e:
            return name, OCRResult(
                text="", source=name, provider=name, tier="unknown", error=str(e)
            )

    with concurrent.futures.ThreadPoolExecutor(max_workers=len(providers)) as executor:
        futures = {executor.submit(safe_call, name): name for name in providers}

        for future in concurrent.futures.as_completed(
            futures, timeout=timeout_per_provider
        ):
            name, result = future.result()
            all_results.append(result)
            if result.success:
                # Cancel remaining
                for f in futures:
                    f.cancel()
                return result, all_results

    # Return best of failures
    best = (
        max(all_results, key=lambda r: r.chars)
        if all_results
        else OCRResult(
            text="",
            source="none",
            provider="none",
            tier="none",
            error="All OCR providers failed",
        )
    )
    return best, all_results


# ============================================================
## LI
# ============================================================


def main():
    import argparse

    parser = argparse.ArgumentParser(description="18-Provider OCR Pipeline")
    parser.add_argument("image", help="Image file path or URL")
    parser.add_argument(
        "-p",
        "--provider",
        action="append",
        dest="providers",
        help="Provider(s) to use (can repeat). Default: full chain",
    )
    parser.add_argument(
        "--parallel",
        action="store_true",
        help="Try providers in parallel (first success wins)",
    )
    parser.add_argument("--list", action="store_true", help="List all providers")
    parser.add_argument("-j", "--json", action="store_true", help="JSON output")
    parser.add_argument("-o", "--output", help="Save extracted text to file")
    parser.add_argument(
        "--preprocess",
        action="store_true",
        default=True,
        help="Apply smart preprocessing (default: on)",
    )
    parser.add_argument(
        "--no-preprocess", action="store_true", help="Disable preprocessing"
    )

    args = parser.parse_args()

    if args.list:
        print("Available OCR providers:")
        for name, fn in PROVIDER_MAP.items():
            print(f"  {name}")
        print(f"\nDefault chain ({len(DEFAULT_CHAIN)} providers):")
        for name in DEFAULT_CHAIN:
            print(f"  {name}")
        return 0

    # Handle URL
    image_bytes: bytes
    if args.image.startswith("http"):
        resp = requests.get(args.image, timeout=30)
        image_bytes = resp.content
    else:
        image_bytes = Path(args.image).read_bytes()

    if args.parallel:
        best, all_results = run_ocr_parallel(image_bytes, args.providers)
    else:
        use_prep = not args.no_preprocess
        best, all_results = run_ocr_chain(
            image_bytes, args.providers, use_preprocessing=use_prep
        )

    if args.json:
        output = {
            "best": best.to_dict(),
            "all_results": [r.to_dict() for r in all_results],
        }
        print(json.dumps(output, indent=2, ensure_ascii=False))
    else:
        print(f"Provider: {best.provider} ({best.tier})")
        print(f"Source: {best.source}")
        print(f"Characters: {best.chars}")
        print(f"Latency: {best.latency_ms}ms")
        print(f"Cost: ${best.cost:.4f}")
        if best.error:
            print(f"Error: {best.error}")
        print(f"\n{'=' * 60}")
        print(best.text or "(No text extracted)")

        print(f"\n--- All Providers ({len(all_results)}) ---")
        for r in all_results:
            status = "OK" if r.success else "FAIL"
            print(
                f"  [{status}] {r.provider}: {r.chars} chars, {r.latency_ms}ms"
                + (f", ${r.cost:.4f}" if r.cost else "")
                + (f", error: {r.error}" if r.error else "")
            )

    if args.output and best.text:
        Path(args.output).write_text(best.text)
        print(f"\nSaved to {args.output}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
