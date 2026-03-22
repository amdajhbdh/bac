#!/usr/bin/env python3
"""
Test script to isolate and verify Yandex Vision OCR functionality
"""

import os
import sys
from pathlib import Path

# Add the scripts directory to Python path
sys.path.insert(0, str(Path(__file__).parent / "scripts"))

from ocr_providers import ocr_yandex, OCRConfig


def main():
    """Test Yandex OCR with a sample image"""
    # Find a test image
    test_dir = Path("tests/bac-images")
    if not test_dir.exists():
        print("❌ Test images directory not found")
        return 1

    test_images = list(test_dir.glob("*.jpg"))
    if not test_images:
        print("❌ No test images found")
        return 1

    # Test image
    test_image = test_images[0]
    print(f"Testing Yandex OCR with image: {test_image.name}")

    # Load config
    config = OCRConfig.from_env()
    print(f"Yandex Folder ID: {config.yandex_folder_id}")
    print(
        f"Yandex IAM Token: {'*' * 10}{config.yandex_iam_token[-10:] if config.yandex_iam_token else 'None'}"
    )

    # Run Yandex OCR
    result = ocr_yandex(test_image, config)

    print(f"\nResult:")
    print(f"  Success: {result.success}")
    print(f"  Text length: {result.chars} characters")
    print(f"  Error: {result.error or 'None'}")
    print(f"  Latency: {result.latency_ms}ms")

    if result.text:
        print(f"\nExtracted text preview:")
        preview = result.text[:200] + "..." if len(result.text) > 200 else result.text
        print(f"  {preview}")

    return 0 if result.success else 1


if __name__ == "__main__":
    exit(main())
