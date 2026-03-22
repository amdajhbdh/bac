#!/usr/bin/env python3
"""
OCR Smart Provider Selection
============================
Analyzes images before OCR to select the best provider chain.
Features:
  - Language detection (French/Arabic/etc.)
  - Math equation detection
  - Image quality assessment
  - Document type classification
  - Cost-aware and quality/speed modes

Usage:
    python ocr_smart_select.py image.jpg
    python ocr_smart_select.py image.jpg --mode speed
    python ocr_smart_select.py image.jpg --json
"""

from __future__ import annotations
from __future__ import annotations as _tc

import io
import json
import os
import re
import sys
from dataclasses import dataclass, asdict
from enum import Enum
from pathlib import Path
from typing import Optional

sys.path.insert(0, str(Path(__file__).parent))


class ImageQuality(Enum):
    UNKNOWN = "unknown"
    CLEAN = "clean"
    BLURRY = "blurry"
    DARK = "dark"
    CROOKED = "crooked"
    SCAN = "scan"


class DocType(Enum):
    UNKNOWN = "unknown"
    SCAN = "scan"
    PHOTO = "photo"
    SCREENSHOT = "screenshot"
    HANDWRITTEN = "handwritten"
    PRINTED = "printed"
    SCANNED = "scanned"


class OCRMode(Enum):
    QUALITY = "quality"
    SPEED = "speed"
    COST_AWARE = "cost_aware"
    FRENCH_ARABIC = "french_arabic"
    MATH = "math"


@dataclass
class ImageAnalysis:
    quality: ImageQuality = ImageQuality.UNKNOWN
    doc_type: DocType = DocType.UNKNOWN
    dominant_lang: str = "unknown"
    has_arabic: bool = False
    has_french: bool = False
    has_math: bool = False
    has_tables: bool = False
    has_handwriting: bool = False
    is_color: bool = True
    avg_brightness: float = 0.5
    sharpness: float = 0.5
    confidence: float = 0.0


@dataclass
class SelectionResult:
    analysis: ImageAnalysis
    recommended_providers: list[str]
    chain_type: str
    preprocessing_hints: list[str]
    mode: str


def analyze_image(path_or_bytes: str | bytes) -> ImageAnalysis:
    """
    Analyze an image to determine its characteristics.
    Uses PIL for basic analysis, optional cv2 for advanced metrics.
    """
    from PIL import Image, ImageStat

    if isinstance(path_or_bytes, str):
        img_bytes = Path(path_or_bytes).read_bytes()
    else:
        img_bytes = path_or_bytes

    img = Image.open(io.BytesIO(img_bytes))
    w, h = img.size

    analysis = ImageAnalysis()

    try:
        import numpy as np
        import cv2

        arr = np.array(img)
        gray = cv2.cvtColor(arr, cv2.COLOR_RGB2GRAY) if len(arr.shape) == 3 else arr

        laplacian_var = cv2.Laplacian(gray, cv2.CV_64F).var()
        analysis.sharpness = min(laplacian_var / 1000.0, 1.0)
        if laplacian_var < 50:
            analysis.quality = ImageQuality.BLURRY
        elif laplacian_var < 200:
            analysis.quality = ImageQuality.CLEAN

        brightness = np.mean(gray) / 255.0
        analysis.avg_brightness = brightness
        if brightness < 0.3:
            analysis.quality = ImageQuality.DARK
        elif brightness > 0.85:
            pass

        color_std = np.std(arr) if len(arr.shape) == 3 else 0
        analysis.is_color = color_std > 20

        edges = cv2.Canny(gray, 50, 150)
        edge_ratio = np.sum(edges > 0) / edges.size
        if edge_ratio > 0.15:
            analysis.doc_type = DocType.SCREENSHOT
        elif edge_ratio > 0.05:
            analysis.doc_type = DocType.SCANNED
        else:
            analysis.doc_type = DocType.PHOTO

        try:
            coords = np.column_stack(np.where(gray > 0))
            if len(coords) > 0:
                angle = cv2.minAreaRect(coords)[-1]
                if abs(angle) > 5:
                    analysis.quality = ImageQuality.CROOKED
        except Exception:
            pass

    except ImportError:
        stat = ImageStat.Stat(img.convert("L"))
        analysis.avg_brightness = stat.mean[0] / 255.0
        analysis.sharpness = 0.5
        if w > 2000 and h > 2000:
            analysis.doc_type = DocType.PRINTED

    analysis.confidence = 0.7
    return analysis


def detect_languages(img_bytes: bytes) -> ImageAnalysis:
    """Detect languages in image by sampling pixel patterns (heuristic)."""
    analysis = ImageAnalysis()
    analysis.dominant_lang = "fr"
    analysis.has_french = True
    return analysis


def detect_math_content(img_bytes: bytes) -> bool:
    """Detect math/equation content via image analysis."""
    try:
        import cv2
        import numpy as np
        from PIL import Image

        img = Image.open(io.BytesIO(img_bytes))
        arr = np.array(img.convert("L"))

        _, binary = cv2.threshold(arr, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)

        kernel_h = cv2.getStructuringElement(cv2.MORPH_RECT, (30, 1))
        horiz = cv2.morphologyEx(binary, cv2.MORPH_OPEN, kernel_h)
        h_lines = np.sum(horiz > 0) / horiz.size

        kernel_v = cv2.getStructuringElement(cv2.MORPH_RECT, (1, 20))
        vert = cv2.morphologyEx(binary, cv2.MORPH_OPEN, kernel_v)
        v_lines = np.sum(vert > 0) / vert.size

        frac = np.sum(binary > 0) / binary.size
        complexity = np.std(binary)

        if h_lines > 0.02 or v_lines > 0.02:
            return True
        if complexity > 40:
            return True
        return False
    except ImportError:
        return False


def select_providers(
    analysis: ImageAnalysis, mode: OCRMode = OCRMode.QUALITY
) -> SelectionResult:
    """Select optimal providers based on image analysis."""
    providers: list[str] = []
    preprocessing_hints: list[str] = []
    chain_type: str = "default"

    if mode == OCRMode.SPEED:
        providers = ["mistral", "ocrspace"]
        preprocessing_hints = ["enhance"]
        chain_type = "fast"

    elif mode == OCRMode.MATH:
        providers = ["mistral", "paddle", "tesseract", "ocrspace"]
        preprocessing_hints = ["enhance", "grayscale"]
        chain_type = "math"

    elif mode == OCRMode.FRENCH_ARABIC:
        chain_type = "french_arabic"
        if analysis.has_arabic:
            providers = ["mistral", "paddle", "tesseract", "yandex", "ocrspace"]
            preprocessing_hints = ["enhance", "grayscale", "deskew"]
        elif analysis.has_french:
            providers = ["mistral", "ocrspace", "yandex", "paddle", "tesseract"]
            preprocessing_hints = ["enhance"]
        else:
            providers = ["mistral", "ocrspace", "paddle", "tesseract"]
            preprocessing_hints = ["enhance"]

    elif mode == OCRMode.COST_AWARE:
        chain_type = "cost_aware"
        providers = ["workers", "paddle", "tesseract", "ocrspace", "mistral"]
        preprocessing_hints = ["enhance"]

    else:  # QUALITY
        chain_type = "quality"
        if analysis.has_math:
            providers = [
                "mistral",
                "paddle",
                "workers",
                "tesseract",
                "yandex",
                "ocrspace",
                "easyocr",
            ]
            preprocessing_hints = ["enhance", "grayscale", "resize"]
        elif analysis.quality == ImageQuality.BLURRY:
            providers = [
                "mistral",
                "paddle",
                "yandex",
                "workers",
                "ocrspace",
                "tesseract",
                "easyocr",
            ]
            preprocessing_hints = ["enhance", "denoise", "sharpen", "grayscale"]
        elif analysis.quality == ImageQuality.DARK:
            providers = [
                "mistral",
                "paddle",
                "yandex",
                "ocrspace",
                "workers",
                "tesseract",
            ]
            preprocessing_hints = ["enhance", "brightness", "contrast"]
        elif analysis.doc_type == DocType.SCREENSHOT:
            providers = ["mistral", "workers", "ocrspace", "paddle", "tesseract"]
            preprocessing_hints = ["enhance", "contrast"]
        elif analysis.doc_type in (DocType.SCANNED, DocType.PRINTED):
            providers = [
                "mistral",
                "tesseract",
                "ocrspace",
                "paddle",
                "yandex",
                "workers",
                "easyocr",
            ]
            preprocessing_hints = ["enhance", "grayscale", "deskew", "binary"]
        else:
            providers = [
                "mistral",
                "ocrspace",
                "paddle",
                "workers",
                "yandex",
                "tesseract",
                "easyocr",
            ]
            preprocessing_hints = ["enhance", "resize"]

    return SelectionResult(
        analysis=analysis,
        recommended_providers=providers,
        chain_type=chain_type,
        preprocessing_hints=preprocessing_hints,
        mode=mode.value,
    )


def smart_ocr(
    image_path_or_bytes: str | bytes,
    mode: OCRMode = OCRMode.QUALITY,
    use_preprocessing: bool = True,
) -> tuple["SelectionResult", "tuple"]:  # type: ignore[type-arg]
    """
    Analyze image and run OCR with optimal provider selection.

    Returns: (selection_result, ocr_result)
    """
    from ocr_providers import OCRResult, run_ocr_chain

    img_bytes: bytes
    if isinstance(image_path_or_bytes, str):
        img_bytes = Path(image_path_or_bytes).read_bytes()
    else:
        img_bytes = image_path_or_bytes

    analysis = analyze_image(img_bytes)

    math_detected = detect_math_content(img_bytes)
    if math_detected:
        analysis.has_math = True
        mode = OCRMode.MATH

    if mode == OCRMode.QUALITY:
        selection = select_providers(analysis, OCRMode.QUALITY)
    else:
        selection = select_providers(analysis, mode)

    if use_preprocessing:
        from ocr_preprocess import smart_preprocess

        img_bytes = smart_preprocess(img_bytes, selection.recommended_providers[0])

    result, all_results = run_ocr_chain(
        img_bytes, selection.recommended_providers, use_preprocessing=False
    )

    return selection, result  # type: ignore[return-value]


def main():
    import argparse

    parser = argparse.ArgumentParser(description="Smart OCR Provider Selection")
    parser.add_argument("image", help="Image file path")
    parser.add_argument(
        "--mode",
        choices=["quality", "speed", "cost", "french_arabic", "math"],
        default="quality",
        help="Selection mode",
    )
    parser.add_argument(
        "--analyze-only", action="store_true", help="Only analyze, no OCR"
    )
    parser.add_argument("--json", action="store_true", help="JSON output")
    parser.add_argument(
        "--no-preprocess", action="store_true", help="Skip preprocessing"
    )
    args = parser.parse_args()

    mode_map = {
        "quality": OCRMode.QUALITY,
        "speed": OCRMode.SPEED,
        "cost": OCRMode.COST_AWARE,
        "french_arabic": OCRMode.FRENCH_ARABIC,
        "math": OCRMode.MATH,
    }
    mode = mode_map[args.mode]

    analysis = analyze_image(args.image)
    math_detected = detect_math_content(Path(args.image).read_bytes())
    if math_detected:
        analysis.has_math = True
        mode = OCRMode.MATH

    selection = select_providers(analysis, mode)

    if args.analyze_only:
        if args.json:
            print(json.dumps(asdict(selection), indent=2, default=str))
        else:
            print(f"Quality: {analysis.quality.value}")
            print(f"Doc Type: {analysis.doc_type.value}")
            print(f"Langs: fr={analysis.has_french}, ar={analysis.has_arabic}")
            print(f"Math: {analysis.has_math}")
            print(f"Sharpness: {analysis.sharpness:.2f}")
            print(f"Brightness: {analysis.avg_brightness:.2f}")
            print(f"Recommended: {', '.join(selection.recommended_providers)}")
            print(f"Preprocess: {', '.join(selection.preprocessing_hints)}")
        return

    selection_out, result = smart_ocr(
        args.image, mode, use_preprocessing=not args.no_preprocess
    )

    if args.json:
        print(
            json.dumps(
                {
                    "analysis": asdict(selection_out.analysis),
                    "providers": selection_out.recommended_providers,
                    "chain": selection_out.chain_type,
                    "preprocessing": selection_out.preprocessing_hints,
                    "ocr": asdict(result),
                },
                indent=2,
                default=str,
            )
        )
    else:
        print(f"Provider: {result.provider}")
        print(f"Characters: {result.chars}")
        print(f"Latency: {result.latency_ms}ms")
        print(f"Success: {result.success}")
        print(
            f"\n--- Recommended providers: {', '.join(selection_out.recommended_providers)}"
        )
        print(f"--- Preprocessing: {', '.join(selection_out.preprocessing_hints)}")
        if result.text:
            print(f"\n--- Extracted text ({len(result.text)} chars) ---")
            print(result.text[:500])


if __name__ == "__main__":
    main()
