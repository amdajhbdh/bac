#!/usr/bin/env python3
"""
OCR Chain Test Suite
===================
Tests each provider independently and the full fallback chain.

Usage:
    python test-ocr-chain.py                    # Test all providers
    python test-ocr-chain.py --provider mistral  # Test single provider
    python test-ocr-chain.py --url IMAGE_URL    # Test with URL
    python test-ocr-chain.py --parallel          # Test parallel mode
    python test-ocr-chain.py --chain worker,paddle,tesseract  # Custom chain
    python test-ocr-chain.py --json             # JSON output
"""

import argparse
import base64
import json
import sys
import time
from pathlib import Path
from typing import Optional

sys.path.insert(0, str(Path(__file__).parent))

from ocr_providers import (
    OCRResult,
    OCRConfig,
    run_ocr_chain,
    run_ocr_parallel,
    PROVIDER_MAP,
    DEFAULT_CHAIN,
    ocr_workers_ai,
    ocr_mistral,
    ocr_ocrspace,
    ocr_yandex,
    ocr_paddle,
    ocr_easyocr,
    ocr_tesseract,
    ocr_llava,
    ocr_qwen25_vl,
    ocr_lfm25_vl,
    ocr_glm,
    ocr_cloudmersive,
    ocr_abbyy,
    ocr_base64ai,
    ocr_ocrapi,
    ocr_nanonets,
)

TEST_IMAGES = {
    "french_text": "https://httpbin.org/image/jpeg",
    "sample": "https://httpbin.org/image/png",
}


def download_image(url: str) -> bytes:
    import requests

    resp = requests.get(url, timeout=30)
    resp.raise_for_status()
    return resp.content


def test_provider(name: str, fn, image_bytes: bytes, config: OCRConfig) -> dict:
    """Test a single provider and return results."""
    start = time.time()
    try:
        result = fn(image_bytes, config)
        return {
            "name": name,
            "success": result.success,
            "chars": result.chars,
            "latency_ms": result.latency_ms or int((time.time() - start) * 1000),
            "source": result.source,
            "provider": result.provider,
            "tier": result.tier,
            "error": result.error,
            "text_preview": result.text[:100] if result.text else "",
        }
    except Exception as e:
        return {
            "name": name,
            "success": False,
            "chars": 0,
            "latency_ms": int((time.time() - start) * 1000),
            "error": str(e),
            "exception": True,
        }


def test_all_providers(image_bytes: bytes, config: OCRConfig) -> list[dict]:
    """Test all 18 providers and return results."""
    results = []

    # Cloud tier
    providers = [
        ("workers_ai", ocr_workers_ai),
        ("mistral", ocr_mistral),
        ("ocrspace", ocr_ocrspace),
        ("yandex", ocr_yandex),
    ]

    # Cloud Shell tier
    providers += [
        ("lfm25vl", ocr_lfm25_vl),
        ("qwen25vl", ocr_qwen25_vl),
        ("llava", ocr_llava),
    ]

    # Local Python tier
    providers += [
        ("paddle", ocr_paddle),
        ("easyocr", ocr_easyocr),
        ("tesseract", ocr_tesseract),
        ("glm", ocr_glm),
    ]

    # Paid cloud tier
    providers += [
        ("cloudmersive", ocr_cloudmersive),
        ("abbyy", ocr_abbyy),
        ("base64ai", ocr_base64ai),
        ("ocrapi", ocr_ocrapi),
        ("nanonets", ocr_nanonets),
    ]

    print(f"\nTesting {len(providers)} providers...")
    print("-" * 60)

    for name, fn in providers:
        print(f"  Testing {name}...", end=" ", flush=True)
        r = test_provider(name, fn, image_bytes, config)
        results.append(r)
        status = "OK" if r["success"] else "FAIL"
        print(
            f"[{status}] {r['chars']} chars, {r['latency_ms']}ms"
            + (f" — {r['error'][:50]}" if r.get("error") else "")
        )

    return results


def test_chain(
    image_bytes: bytes, config: OCRConfig, providers: Optional[list[str]] = None
) -> dict:
    """Test the full fallback chain."""
    print(f"\nTesting OCR chain (providers: {providers or DEFAULT_CHAIN})...")
    print("-" * 60)

    start = time.time()
    best, all_results = run_ocr_chain(image_bytes, providers, config)
    total_ms = int((time.time() - start) * 1000)

    print(f"\nChain Results:")
    print(f"  Best: {best.provider} ({best.source})")
    print(f"  Chars: {best.chars}")
    print(f"  Latency: {total_ms}ms")
    print(f"  Total providers tried: {len(all_results)}")

    return {
        "best": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_ms": total_ms,
    }


def test_parallel(image_bytes: bytes, config: OCRConfig) -> dict:
    """Test parallel OCR mode."""
    print(f"\nTesting parallel OCR...")
    print("-" * 60)

    start = time.time()
    best, all_results = run_ocr_parallel(image_bytes)
    total_ms = int((time.time() - start) * 1000)

    print(f"\nParallel Results:")
    print(f"  Best: {best.provider} ({best.source})")
    print(f"  Chars: {best.chars}")
    print(f"  Latency: {total_ms}ms")
    print(f"  Total providers tried: {len(all_results)}")
    print(f"  Providers attempted: {', '.join(r.provider for r in all_results)}")

    return {
        "best": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_ms": total_ms,
    }


def print_summary(results: list[dict]):
    """Print a summary table of results."""
    print(f"\n{'=' * 70}")
    print(f"SUMMARY — {len(results)} providers tested")
    print(f"{'=' * 70}")

    # Sort by success then chars
    results.sort(key=lambda r: (r["success"], r["chars"]), reverse=True)

    print(f"\n{'Provider':<20} {'Tier':<15} {'Status':<8} {'Chars':<8} {'Latency':<10}")
    print("-" * 70)

    success_count = 0
    total_chars = 0

    for r in results:
        status = "OK" if r["success"] else "FAIL"
        if r["success"]:
            success_count += 1
            total_chars += r["chars"]

        print(
            f"{r['name']:<20} {r.get('tier', '?'):<15} {status:<8} "
            f"{r['chars']:<8} {r['latency_ms']:<10}ms"
            + (f" ⚠ {r.get('error', '')[:30]}" if r.get("error") else "")
        )

    print("-" * 70)
    print(
        f"Success rate: {success_count}/{len(results)} ({100 * success_count / len(results):.0f}%)"
    )
    print(f"Total chars extracted: {total_chars}")
    avg_latency = sum(r["latency_ms"] for r in results) / len(results) if results else 0
    print(f"Average latency: {avg_latency:.0f}ms")


def main():
    parser = argparse.ArgumentParser(description="OCR Chain Test Suite")
    parser.add_argument("--image", help="Local image file path")
    parser.add_argument("--url", help="Image URL")
    parser.add_argument("--provider", help="Test single provider only")
    parser.add_argument("--chain", help="Comma-separated chain (e.g. worker,paddle)")
    parser.add_argument("--parallel", action="store_true", help="Test parallel mode")
    parser.add_argument("--all", action="store_true", help="Test all modes")
    parser.add_argument("--json", action="store_true", help="JSON output")
    parser.add_argument("--output", help="Save results to JSON file")

    args = parser.parse_args()

    # Load image
    if args.image:
        image_bytes = Path(args.image).read_bytes()
        source = f"file://{args.image}"
    elif args.url:
        print(f"Downloading image from {args.url}...")
        image_bytes = download_image(args.url)
        source = args.url
    else:
        print(f"Using default test image: {TEST_IMAGES['french_text']}")
        image_bytes = download_image(TEST_IMAGES["french_text"])
        source = TEST_IMAGES["french_text"]

    print(
        f"Image size: {len(image_bytes)} bytes ({len(base64.b64encode(image_bytes))} base64)"
    )

    config = OCRConfig.from_env()
    all_results = {}
    chain_results = []

    # Single provider test
    if args.provider:
        fn = PROVIDER_MAP.get(args.provider)
        if not fn:
            print(f"Unknown provider: {args.provider}")
            print(f"Available: {', '.join(PROVIDER_MAP.keys())}")
            return 1
        r = test_provider(args.provider, fn, image_bytes, config)
        all_results[args.provider] = r
        chain_results = [r]

    # Test all providers
    if args.all or (not args.provider and not args.chain and not args.parallel):
        results = test_all_providers(image_bytes, config)
        all_results["all"] = results
        chain_results = results

    # Chain test
    if args.chain or (args.all and not args.provider):
        providers = [p.strip() for p in args.chain.split(",")] if args.chain else None
        r = test_chain(image_bytes, config, providers)
        all_results["chain"] = r
        if "all" not in all_results:
            chain_results = r["all_results"]

    # Parallel test
    if args.parallel or (args.all and not args.provider):
        r = test_parallel(image_bytes, config)
        all_results["parallel"] = r

    # Print summary
    if chain_results and not args.json:
        print_summary(
            chain_results if isinstance(chain_results, list) else [chain_results[0]]
        )

    # Output
    output = {
        "source": source,
        "image_size": len(image_bytes),
        "results": all_results,
        "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
    }

    if args.json:
        print(json.dumps(output, indent=2, ensure_ascii=False))
    elif not args.all:
        print(
            f"\nBest result: {all_results.get('all', [{}])[0] if isinstance(all_results.get('all'), list) else all_results}"
        )

    if args.output:
        Path(args.output).write_text(json.dumps(output, indent=2, ensure_ascii=False))
        print(f"\nResults saved to {args.output}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
