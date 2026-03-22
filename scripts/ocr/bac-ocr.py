#!/usr/bin/env python3
"""
OCR Pipeline for BAC Study System.
Supports Tesseract (local) and Workers AI (edge).
"""

import os
import subprocess
import sys
import tempfile
from pathlib import Path


def ocr_tesseract(image_path, languages="fra", verbose=False):
    """OCR using Tesseract (local)."""
    languages = languages.split("+") if "+" in languages else [languages]

    # Build language string
    lang_str = "+".join(languages)

    try:
        result = subprocess.run(
            ["tesseract", image_path, "stdout", "-l", lang_str],
            capture_output=True,
            text=True,
            timeout=60,
        )

        if result.returncode != 0:
            return {"error": f"Tesseract error: {result.stderr}"}

        text = result.stdout.strip()

        if verbose:
            print(
                f"[Tesseract] Extracted {len(text)} chars from {Path(image_path).name}"
            )

        return {"text": text, "source": "tesseract", "languages": languages}

    except subprocess.TimeoutExpired:
        return {"error": "OCR timeout"}
    except FileNotFoundError:
        return {"error": "Tesseract not installed"}


def ocr_workers_ai(image_path, api_token=None, account_id=None, verbose=False):
    """OCR using Workers AI (edge)."""
    return {"error": "Use the Worker /rag/ocr endpoint for cloud OCR"}


def ocr_paddle(image_path, languages="fr", verbose=False):
    """OCR using PaddleOCR (higher accuracy for French)."""
    try:
        from paddleocr import PaddleOCR

        ocr = PaddleOCR(
            use_angle_cls=True, lang=languages, use_gpu=False, show_log=False
        )
        result = ocr.ocr(image_path, cls=True)

        if not result or not result[0]:
            return {"text": "", "source": "paddleocr", "languages": [languages]}

        # Extract all text lines
        text_lines = []
        for line in result[0]:
            if line and len(line) >= 2:
                text_lines.append(str(line[1][0]))

        text = "\n".join(text_lines)

        if verbose:
            print(
                f"[PaddleOCR] Extracted {len(text)} chars from {Path(image_path).name}"
            )

        return {"text": text, "source": "paddleocr", "languages": [languages]}

    except ImportError:
        return {
            "error": "PaddleOCR not installed. Run: pip install paddlepaddle paddleocr"
        }
    except Exception as e:
        return {"error": f"PaddleOCR failed: {str(e)}"}


def ocr_file(file_path, method="tesseract", languages="fra"):
    """Main OCR function."""
    if not os.path.exists(file_path):
        return {"error": f"File not found: {file_path}"}

    if method == "tesseract":
        return ocr_tesseract(file_path, languages)
    elif method == "paddle":
        return ocr_paddle(file_path, languages)
    elif method == "workers":
        return ocr_workers_ai(file_path)
    elif method == "workers":
        return ocr_workers_ai(file_path)
    else:
        return {"error": f"Unknown method: {method}"}


def batch_ocr(directory, extension=".png", method="tesseract", languages="fra"):
    """Process all images in a directory."""
    results = []

    for path in Path(directory).glob(f"*{extension}"):
        print(f"Processing: {path.name}")
        result = ocr_file(str(path), method, languages)
        results.append({"file": path.name, "result": result})

    return results


def main():
    import argparse

    parser = argparse.ArgumentParser(description="OCR for BAC Study System")
    parser.add_argument("file", nargs="?", help="Image file to OCR")
    parser.add_argument(
        "-m",
        "--method",
        choices=["tesseract", "paddle", "workers"],
        default="tesseract",
    )
    parser.add_argument(
        "-l", "--languages", default="fra", help="Languages: fra, ara, fra+ara"
    )
    parser.add_argument("-d", "--directory", help="Process all images in directory")
    parser.add_argument("-o", "--output", help="Output file")
    parser.add_argument("-v", "--verbose", action="store_true")

    args = parser.parse_args()

    if args.directory:
        results = batch_ocr(
            args.directory, method=args.method, languages=args.languages
        )
        print(f"\nProcessed {len(results)} files")

        if args.output:
            import json

            with open(args.output, "w") as f:
                json.dump(results, f, indent=2)
            print(f"Results saved to {args.output}")

    elif args.file:
        result = ocr_file(args.file, args.method, args.languages)

        if "error" in result:
            print(f"Error: {result['error']}")
            return 1

        print(result["text"])

        if args.output:
            with open(args.output, "w") as f:
                f.write(result["text"])
            print(f"Saved to {args.output}")

    else:
        parser.print_help()

    return 0


if __name__ == "__main__":
    sys.exit(main())
