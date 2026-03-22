#!/usr/bin/env python3
"""
Test script to verify the complete OCR setup
"""

import requests
import json
import os
from pathlib import Path


def test_local_server():
    """Test the local OCR server"""
    try:
        response = requests.get("http://localhost:8000/health")
        if response.status_code == 200:
            print("✅ Local OCR server is running")
            data = response.json()
            print(f"   Available providers: {len(data['providers'])}")
            return True
        else:
            print("❌ Local OCR server is not responding")
            return False
    except Exception as e:
        print(f"❌ Local OCR server error: {e}")
        return False


def test_worker_endpoint():
    """Test the Cloudflare Worker endpoint"""
    try:
        # Test with a simple health check
        response = requests.get("https://bac-ocr-worker.amdajhbdh.workers.dev/")
        if response.status_code == 200:
            print("✅ Cloudflare Worker is accessible")
            return True
        else:
            print("❌ Cloudflare Worker is not responding")
            return False
    except Exception as e:
        print(f"❌ Cloudflare Worker error: {e}")
        return False


def test_ocr_with_sample():
    """Test OCR with a sample image"""
    try:
        # Find a test image
        test_images = list(Path("tests/bac-images").glob("*.jpg"))
        if not test_images:
            print("❌ No test images found")
            return False

        # Test with first image
        with open(test_images[0], "rb") as f:
            response = requests.post(
                "http://localhost:8000/ocr",
                files={"image": f},
                data={"providers": "workers,paddle,tesseract"},
            )

        if response.status_code == 200:
            data = response.json()
            result = data["result"]
            print(f"✅ OCR test successful")
            print(f"   Provider: {result['provider']}")
            print(f"   Characters extracted: {result['chars']}")
            print(f"   Process time: {data['total_ms']}ms")
            return True
        else:
            print("❌ OCR test failed")
            return False
    except Exception as e:
        print(f"❌ OCR test error: {e}")
        return False


def main():
    """Run all tests"""
    print("Testing complete OCR setup...\n")

    tests = [
        ("Local server", test_local_server),
        ("Cloudflare Worker", test_worker_endpoint),
        ("OCR functionality", test_ocr_with_sample),
    ]

    passed = 0
    total = len(tests)

    for name, test_func in tests:
        print(f"Testing {name}...")
        if test_func():
            passed += 1
        print()

    print(f"Results: {passed}/{total} tests passed")

    if passed == total:
        print("\n🎉 All tests passed! Your OCR system is fully operational.")
    else:
        print(f"\n⚠️  {total - passed} test(s) failed. Please check the setup.")


if __name__ == "__main__":
    main()
