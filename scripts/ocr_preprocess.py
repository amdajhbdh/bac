#!/usr/bin/env python3
"""
OCR Image Preprocessing Pipeline
==============================
Prepares images for optimal OCR results across all providers.
Handles: rotation, deskewing, contrast, noise removal, cropping, format conversion.

Usage:
    python ocr-preprocess.py input.jpg output.jpg --enhance
    python ocr-preprocess.py input.jpg --auto-deskew --denoise
"""

from __future__ import annotations

import base64
import io
import sys
from pathlib import Path
from typing import Optional

sys.path.insert(0, str(Path(__file__).parent))


def load_image(path_or_bytes) -> Optional["Image.Image"]:
    """Load image from path or bytes."""
    try:
        from PIL import Image

        if isinstance(path_or_bytes, (str, Path)):
            return Image.open(path_or_bytes)
        return Image.open(io.BytesIO(path_or_bytes))
    except ImportError:
        return None


def save_image(img, path_or_bytes, format: str = "JPEG", quality: int = 95):
    """Save image to path or bytes."""
    buf = io.BytesIO()
    img.save(buf, format=format, quality=quality)
    if isinstance(path_or_bytes, (str, Path)):
        Path(path_or_bytes).write_bytes(buf.getvalue())
    return buf.getvalue()


def auto_contrast(img: "Image.Image", cutoff: int = 2) -> "Image.Image":
    """Auto-adjust contrast to improve OCR readability."""
    from PIL import Image, ImageOps

    return ImageOps.autocontrast(img, cutoff=cutoff)


def to_grayscale(img: "Image.Image") -> "Image.Image":
    """Convert to grayscale (better for most OCR engines)."""
    return img.convert("L").convert("RGB")


def denoise(img: "Image.Image", strength: int = 1) -> "Image.Image":
    """Remove noise while preserving text edges."""
    try:
        import cv2
        import numpy as np
        import PIL.Image

        # Convert to numpy
        arr = np.array(img)
        if len(arr.shape) == 3:
            gray = cv2.cvtColor(arr, cv2.COLOR_RGB2GRAY)
        else:
            gray = arr

        # Denoise
        denoised = cv2.fastNlMeansDenoising(gray, None, strength * 10, 7, 21)
        return PIL.Image.fromarray(denoised)
    except ImportError:
        return img


def deskew(img: "Image.Image") -> "Image.Image":
    """Detect and correct image skew."""
    try:
        import cv2
        import numpy as np
        import PIL.Image

        arr = np.array(img.convert("L"))
        coords = np.column_stack(np.where(arr > 0))
        if len(coords) == 0:
            return img

        angle = cv2.minAreaRect(coords)[-1]
        if angle < -45:
            angle = -(90 + angle)
        else:
            angle = -angle

        if abs(angle) < 0.5:
            return img  # Already straight

        # Rotate
        h, w = arr.shape
        center = (w / 2, h / 2)
        M = cv2.getRotationMatrix2D(center, angle, 1.0)
        rotated = cv2.warpAffine(
            np.array(img),
            M,
            (w, h),
            flags=cv2.INTER_CUBIC,
            borderMode=cv2.BORDER_REPLICATE,
        )
        return PIL.Image.fromarray(rotated)
    except ImportError:
        return img


def resize_for_ocr(img: "Image.Image", max_size: int = 2048) -> "Image.Image":
    """Resize image to optimal size for OCR (not too large, not too small)."""
    from PIL import Image

    w, h = img.size
    if max(w, h) <= max_size:
        return img

    scale = max_size / max(w, h)
    new_w, new_h = int(w * scale), int(h * scale)
    return img.resize((new_w, new_h), Image.LANCZOS)


def to_binary(img: "Image.Image", threshold: int = 128) -> "Image.Image":
    """Convert to pure black and white (good for low-quality scans)."""
    from PIL import Image

    gray = img.convert("L")
    return gray.point(lambda x: 0 if x < threshold else 255, "1").convert("RGB")


def preprocess_pipeline(
    image_bytes: bytes,
    enhance: bool = True,
    do_grayscale: bool = False,
    do_denoise: bool = False,
    do_deskew: bool = False,
    resize: bool = True,
    binary: bool = False,
) -> bytes:
    """
    Apply preprocessing pipeline to image bytes.

    Returns processed image as JPEG bytes.
    """
    img = load_image(image_bytes)
    if img is None:
        return image_bytes

    if enhance:
        img = auto_contrast(img)

    if do_deskew:
        img = deskew(img)

    if do_denoise:
        img = denoise(img)

    if do_grayscale:
        img = to_grayscale(img)

    if resize:
        img = resize_for_ocr(img)

    if binary:
        img = to_binary(img)

    return save_image(img, io.BytesIO(), format="JPEG", quality=90)


def smart_preprocess(image_bytes: bytes, provider: str) -> bytes:
    """
    Smart preprocessing based on provider characteristics.

    Provider-specific optimizations:
      - tesseract: binary + deskew (best)
      - paddle: grayscale + contrast (already handles this)
      - easyocr: minimal processing (deep learning handles variation)
      - mistral: no preprocessing needed (vision model)
      - workers (llama-vision): no preprocessing needed
      - arabic text: grayscale + contrast enhancement
    """
    if provider in (
        "mistral",
        "llava",
        "lfm25vl",
        "qwen25vl",
        "workers",
        "llama-vision",
        "uform",
    ):
        # Vision-language models handle various conditions well
        return preprocess_pipeline(image_bytes, enhance=True, resize=True)
    elif provider in ("tesseract", "ocrspace"):
        # Traditional OCR benefits from binarization
        return preprocess_pipeline(
            image_bytes, enhance=True, do_deskew=True, do_grayscale=True, binary=True
        )
    elif provider in ("paddle", "easyocr"):
        # Deep learning OCR likes grayscale with contrast
        return preprocess_pipeline(
            image_bytes, enhance=True, do_grayscale=True, resize=True
        )
    else:
        # Default: light preprocessing
        return preprocess_pipeline(image_bytes, enhance=True, resize=True)


# ============================================================
# CLI
# ============================================================


def main():
    import argparse

    parser = argparse.ArgumentParser(description="OCR Image Preprocessing")
    parser.add_argument("input", help="Input image file")
    parser.add_argument("output", nargs="?", help="Output image file")
    parser.add_argument("--enhance", action="store_true", help="Auto-contrast")
    parser.add_argument("--grayscale", action="store_true", help="Convert to grayscale")
    parser.add_argument("--denoise", action="store_true", help="Remove noise")
    parser.add_argument("--deskew", action="store_true", help="Correct skew")
    parser.add_argument("--resize", action="store_true", help="Resize to optimal size")
    parser.add_argument("--binary", action="store_true", help="Black and white")
    parser.add_argument(
        "--smart", help="Smart preprocess for provider (tesseract, paddle, etc.)"
    )
    parser.add_argument("--preview", action="store_true", help="Show preview")
    args = parser.parse_args()

    img_bytes = Path(args.input).read_bytes()

    if args.smart:
        result = smart_preprocess(img_bytes, args.smart)
        print(f"Smart preprocess for {args.smart}: {len(result)} bytes")
    else:
        result = preprocess_pipeline(
            img_bytes,
            enhance=args.enhance,
            do_grayscale=args.grayscale,
            do_denoise=args.denoise,
            do_deskew=args.deskew,
            resize=args.resize or not args.output,
            binary=args.binary,
        )

    if args.output:
        Path(args.output).write_bytes(result)
        print(f"Saved to {args.output} ({len(result)} bytes)")
    else:
        import tempfile
        import os

        tmp = tempfile.NamedTemporaryFile(suffix=".jpg", delete=False)
        tmp.write(result)
        print(f"Preview: {tmp.name}")
        if args.preview:
            try:
                from PIL import Image

                Image.open(tmp.name).show()
            except Exception:
                pass
        os.unlink(tmp.name)


if __name__ == "__main__":
    main()
