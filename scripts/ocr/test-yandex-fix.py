#!/usr/bin/env python3
"""
Test script to verify and fix Yandex Vision OCR implementation
"""

import os
import sys
import time
import base64
import requests
from pathlib import Path

# Add the scripts directory to Python path
sys.path.insert(0, str(Path(__file__).parent / "scripts"))

from ocr_providers import OCRConfig


def get_new_iam_token():
    """Get a new IAM token from Yandex Cloud"""
    # This would normally require a service account key file
    # For now, we'll just validate the existing token
    config = OCRConfig.from_env()
    return config.yandex_iam_token


def test_yandex_direct_api(image_path, config):
    """Direct test of Yandex Vision API"""
    print("Testing Yandex Vision API directly...")

    # Read image
    with open(image_path, "rb") as f:
        img_bytes = f.read()

    # Convert to base64
    b64 = base64.b64encode(img_bytes).decode()

    # Try the API call
    url = "https://ocr.api.cloud.yandex.net/ocr/v1/recognizeText"
    headers = {
        "Authorization": f"Bearer {config.yandex_iam_token}",
        "x-folder-id": config.yandex_folder_id or "default",
        "Content-Type": "application/json",
    }

    payload = {
        "mimeType": "JPEG",
        "languageCodes": ["fr", "ar", "en"],
        "model": "page",
        "content": b64,
    }

    print(f"Making request to {url}")
    print(f"Folder ID: {config.yandex_folder_id}")
    print(
        f"Token preview: {config.yandex_iam_token[:10]}...{config.yandex_iam_token[-10:]}"
    )

    try:
        response = requests.post(url, headers=headers, json=payload, timeout=30)
        print(f"Response status: {response.status_code}")
        print(f"Response headers: {dict(response.headers)}")

        if response.status_code == 200:
            data = response.json()
            print("Success! Response data:")
            print(data)
            return data
        else:
            print(f"Error response: {response.text}")
            return {"error": f"HTTP {response.status_code}: {response.text}"}
    except Exception as e:
        print(f"Exception: {e}")
        return {"error": str(e)}


def main():
    """Test Yandex OCR with fixes"""
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
        f"Yandex IAM Token length: {len(config.yandex_iam_token) if config.yandex_iam_token else 0}"
    )

    # Try to get a new token if the current one seems invalid
    if not config.yandex_iam_token or len(config.yandex_iam_token) < 50:
        print("⚠️  Token seems invalid, trying to refresh...")
        new_token = get_new_iam_token()
        if new_token and new_token != config.yandex_iam_token:
            print("New token obtained")
            config.yandex_iam_token = new_token
        else:
            print("Could not obtain new token")

    # Test direct API call
    result = test_yandex_direct_api(test_image, config)

    if "error" in result:
        print(f"\n❌ Yandex OCR failed: {result['error']}")
        return 1
    else:
        print("\n✅ Yandex OCR succeeded!")
        return 0


if __name__ == "__main__":
    exit(main())
