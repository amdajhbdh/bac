#!/usr/bin/env python3
"""
Local OCR FastAPI Server
========================
Exposes OCR providers as an HTTP API for local fallback chain.
Also serves as the bridge between Telegram bot and Cloud Shell Ollama.

Usage:
    python local-ocr-server.py --port 8081
    uvicorn local-ocr-server:app --port 8081 --host 0.0.0.0

Endpoints:
    POST /ocr          — Full OCR chain (parallel or sequential)
    POST /ocr/chain    — Sequential chain with all providers tried
    POST /ocr/parallel — Parallel try, return first success
    GET  /providers     — List all available providers
    GET  /health        — Health check
"""

import base64
import hashlib
import io
import os
import time
from pathlib import Path
from typing import Optional

try:
    from fastapi import FastAPI, File, Form, HTTPException, UploadFile
    from fastapi.middleware.cors import CORSMiddleware
    from fastapi.responses import JSONResponse
    import uvicorn

    FASTAPI_AVAILABLE = True
except ImportError:
    FASTAPI_AVAILABLE = False
    print("WARNING: FastAPI not installed. Run: pip install fastapi uvicorn")

import sys

sys.path.insert(0, str(Path(__file__).parent))

try:
    from ocr_providers import (
        OCRResult,
        OCRConfig,
        run_ocr_chain,
        run_ocr_parallel,
        PROVIDER_MAP,
        DEFAULT_CHAIN,
    )
except ImportError:
    from .ocr_providers import (
        OCRResult,
        OCRConfig,
        run_ocr_chain,
        run_ocr_parallel,
        PROVIDER_MAP,
        DEFAULT_CHAIN,
    )


# ============================================================
# CACHE
# ============================================================


class OCRCache:
    """Simple in-memory cache for OCR results keyed by image hash."""

    def __init__(self, max_size: int = 100, ttl_seconds: int = 3600):
        self.cache: dict[str, tuple[OCRResult, float]] = {}
        self.max_size = max_size
        self.ttl = ttl_seconds

    def get(self, image_bytes: bytes) -> Optional[OCRResult]:
        key = hashlib.sha256(image_bytes).hexdigest()[:32]
        if key in self.cache:
            result, timestamp = self.cache[key]
            if time.time() - timestamp < self.ttl:
                result.source += " (cached)"
                return result
            del self.cache[key]
        return None

    def set(self, image_bytes: bytes, result: OCRResult):
        key = hashlib.sha256(image_bytes).hexdigest()[:32]
        if len(self.cache) >= self.max_size:
            oldest = min(self.cache.items(), key=lambda x: x[1][1])
            del self.cache[oldest[0]]
        self.cache[key] = (result, time.time())


cache = OCRCache()


# ============================================================
# APP
# ============================================================

app = FastAPI(
    title="BAC OCR Server",
    version="1.0.0",
    description="18-Provider OCR with fallback chain",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
def health():
    """Health check endpoint."""
    return {
        "status": "ok",
        "providers": list(PROVIDER_MAP.keys()),
        "cache_size": len(cache.cache),
    }


@app.get("/providers")
def list_providers():
    """List all available OCR providers."""
    providers = []
    for name, fn in PROVIDER_MAP.items():
        try:
            # Try calling with empty bytes to check if it errors
            result = fn(b"test", OCRConfig.from_env())
            providers.append(
                {
                    "name": name,
                    "tier": result.tier,
                    "success": result.success or result.error is not None,
                    "error": result.error if not result.success else None,
                }
            )
        except Exception as e:
            providers.append(
                {
                    "name": name,
                    "tier": "unknown",
                    "success": False,
                    "error": str(e),
                }
            )

    return {
        "providers": providers,
        "default_chain": DEFAULT_CHAIN,
        "total": len(providers),
    }


@app.post("/ocr")
def ocr_full(
    image: UploadFile = File(...),
    chain: Optional[str] = Form(None),
    parallel: bool = Form(False),
    use_cache: bool = Form(True),
):
    """
    Full OCR with configurable chain.

    Args:
        image: Image file (jpg, png, etc.)
        chain: Comma-separated provider names (default: all in order)
        parallel: Try providers in parallel (first success wins)
        use_cache: Use image hash cache
    """
    if not FASTAPI_AVAILABLE:
        raise HTTPException(503, "FastAPI not installed")

    image_bytes = image.file.read()
    start_time = time.time()

    # Check cache
    if use_cache:
        cached = cache.get(image_bytes)
        if cached:
            cached.latency_ms = int((time.time() - start_time) * 1000)
            return {
                "result": cached.to_dict(),
                "cached": True,
                "total_ms": cached.latency_ms,
            }

    # Parse chain
    providers = None
    if chain:
        providers = [p.strip() for p in chain.split(",")]

    # Run chain
    if parallel:
        best, all_results = run_ocr_parallel(image_bytes, providers)
    else:
        best, all_results = run_ocr_chain(image_bytes, providers)

    total_ms = int((time.time() - start_time) * 1000)

    # Cache successful result
    if best.success and use_cache:
        cache.set(image_bytes, best)

    return {
        "result": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_providers_tried": len(all_results),
        "total_ms": total_ms,
        "cached": False,
    }


@app.post("/ocr/chain")
def ocr_chain(
    image: UploadFile = File(...),
    providers: Optional[str] = Form(None),
):
    """Sequential chain (all providers tried until success)."""
    if not FASTAPI_AVAILABLE:
        raise HTTPException(503, "FastAPI not installed")

    image_bytes = image.file.read()
    start = time.time()

    provider_list = None
    if providers:
        provider_list = [p.strip() for p in providers.split(",")]

    best, all_results = run_ocr_chain(image_bytes, provider_list)
    total_ms = int((time.time() - start) * 1000)

    return {
        "result": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_ms": total_ms,
    }


@app.post("/ocr/parallel")
def ocr_parallel(image: UploadFile = File(...)):
    """Parallel try (first success wins)."""
    if not FASTAPI_AVAILABLE:
        raise HTTPException(503, "FastAPI not installed")

    image_bytes = image.file.read()
    start = time.time()

    best, all_results = run_ocr_parallel(image_bytes)
    total_ms = int((time.time() - start) * 1000)

    return {
        "result": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_ms": total_ms,
    }


@app.post("/ocr/base64")
def ocr_base64(image_b64: str = Form(...), providers: Optional[str] = Form(None)):
    """OCR from base64 string."""
    if not FASTAPI_AVAILABLE:
        raise HTTPException(503, "FastAPI not installed")

    try:
        image_bytes = base64.b64decode(image_b64)
    except Exception as e:
        raise HTTPException(400, f"Invalid base64: {e}")

    start = time.time()
    provider_list = None
    if providers:
        provider_list = [p.strip() for p in providers.split(",")]

    best, all_results = run_ocr_chain(image_bytes, provider_list)
    total_ms = int((time.time() - start) * 1000)

    return {
        "result": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_ms": total_ms,
    }


@app.post("/ocr/url")
def ocr_from_url(url: str = Form(...), providers: Optional[str] = Form(None)):
    """OCR from URL."""
    if not FASTAPI_AVAILABLE:
        raise HTTPException(503, "FastAPI not installed")

    import requests as req

    try:
        resp = req.get(url, timeout=30)
        resp.raise_for_status()
        image_bytes = resp.content
    except Exception as e:
        raise HTTPException(400, f"Failed to fetch URL: {e}")

    start = time.time()
    provider_list = None
    if providers:
        provider_list = [p.strip() for p in providers.split(",")]

    best, all_results = run_ocr_chain(image_bytes, provider_list)
    total_ms = int((time.time() - start) * 1000)

    return {
        "result": best.to_dict(),
        "all_results": [r.to_dict() for r in all_results],
        "total_ms": total_ms,
    }


@app.get("/cache/stats")
def cache_stats():
    """Get cache statistics."""
    return {
        "size": len(cache.cache),
        "max_size": cache.max_size,
        "ttl_seconds": cache.ttl,
        "keys": list(cache.cache.keys()),
    }


@app.post("/cache/clear")
def cache_clear():
    """Clear the OCR cache."""
    cache.cache.clear()
    return {"status": "cleared", "size": 0}


def main():
    import argparse

    parser = argparse.ArgumentParser(description="BAC Local OCR Server")
    parser.add_argument("--host", default="0.0.0.0", help="Host to bind to")
    parser.add_argument("--port", type=int, default=8081, help="Port to bind to")
    parser.add_argument("--reload", action="store_true", help="Enable auto-reload")
    parser.add_argument("--workers", type=int, default=4, help="Number of workers")
    args = parser.parse_args()

    if not FASTAPI_AVAILABLE:
        print("ERROR: FastAPI is required. Install with: pip install fastapi uvicorn")
        sys.exit(1)

    print(f"Starting BAC OCR Server on {args.host}:{args.port}")
    print(f"Available providers: {', '.join(PROVIDER_MAP.keys())}")
    print(f"Endpoints:")
    print(f"  POST /ocr          — Full OCR chain")
    print(f"  POST /ocr/chain    — Sequential chain")
    print(f"  POST /ocr/parallel — Parallel try")
    print(f"  POST /ocr/base64   — From base64 string")
    print(f"  POST /ocr/url      — From URL")
    print(f"  GET  /providers    — List providers")
    print(f"  GET  /health       — Health check")
    print(f"  GET  /cache/stats  — Cache stats")
    print(f"  POST /cache/clear  — Clear cache")

    uvicorn.run(
        "local-ocr-server:app",
        host=args.host,
        port=args.port,
        reload=args.reload,
        workers=1 if args.reload else args.workers,
    )


if __name__ == "__main__":
    main()
