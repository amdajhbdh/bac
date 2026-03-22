#!/usr/bin/env python3
"""
Upload vault content to RAG via Worker batch endpoint.
Workers AI (free) for embeddings — no rate limits on compute.
Note: Upstash Vector free tier allows 10K writes/day.
"""

import json, time, yaml, sys, subprocess
from pathlib import Path

WORKER_URL = "https://bac-api.amdajhbdh.workers.dev"
RAG_DIR = Path.home() / ".config/aichat/rags"
BATCH_SIZE = 10
MAX_RETRIES = 5
INITIAL_SLEEP = 5


def curl_post(url, payload, timeout=60):
    try:
        result = subprocess.run(
            [
                "curl",
                "-s",
                "-X",
                "POST",
                url,
                "-H",
                "Content-Type: application/json",
                "-d",
                json.dumps(payload),
                "--max-time",
                str(timeout),
            ],
            capture_output=True,
            text=True,
            timeout=timeout + 15,
        )
        if result.returncode == 28:
            return {"error": "curl timeout", "stdout": ""}
        if result.returncode != 0:
            return {
                "error": f"curl rc={result.returncode}",
                "stdout": result.stdout[:200],
            }
        try:
            return json.loads(result.stdout)
        except json.JSONDecodeError:
            return {"error": f"bad json: {result.stdout[:200]}"}
    except subprocess.TimeoutExpired:
        return {"error": "python timeout"}
    except Exception as e:
        return {"error": str(e)}


def chunk_text(text, size=2000, overlap=200):
    if not text.strip():
        return []
    chunks = []
    for start in range(0, len(text), size - overlap):
        c = text[start : start + size]
        if c.strip():
            chunks.append(c)
    return chunks


def process_rag_file(filepath: Path, subject: str):
    with open(filepath) as f:
        rag = yaml.safe_load(f)
    print(f"  {len(rag['files'])} files, building chunks...", flush=True)

    all_chunks_data = []
    for fid, fdata in rag["files"].items():
        path = fdata.get("path", fid)
        if path.startswith("./"):
            path = path[2:]
        for doc in fdata.get("documents", []):
            chunks = chunk_text(doc["page_content"])
            for i, c in enumerate(chunks):
                all_chunks_data.append(
                    {"text": c, "subject": subject, "path": f"{path}::{i}"}
                )

    print(f"  {len(all_chunks_data)} chunks", flush=True)

    total = 0
    total_errors = 0
    sleep_time = INITIAL_SLEEP

    for start in range(0, len(all_chunks_data), BATCH_SIZE):
        batch = all_chunks_data[start : start + BATCH_SIZE]
        batch_chunks = [
            {
                "text": d["text"],
                "subject": subject,
                "path": d["path"].rsplit("::", 1)[0],
            }
            for d in batch
        ]

        attempt = 0
        result = {"error": "never tried"}
        while attempt < MAX_RETRIES:
            result = curl_post(f"{WORKER_URL}/rag/add-batch", {"chunks": batch_chunks})
            if "error" not in result:
                break
            attempt += 1
            err = str(result.get("error", ""))[:60]
            print(f"  Attempt {attempt} failed: {err}", flush=True)
            if "Exceeded" in err:
                print(
                    f"  Hit Upstash daily limit. Sleeping {sleep_time}s before retry.",
                    flush=True,
                )
                time.sleep(sleep_time)
                sleep_time = min(sleep_time * 2, 120)
            elif "timeout" in err.lower():
                time.sleep(3 * attempt)
            else:
                time.sleep(2 * attempt)

        if "error" in result:
            print(
                f"  Batch error after {MAX_RETRIES} retries: {str(result.get('error', ''))[:80]}",
                flush=True,
            )
            total_errors += len(batch_chunks)
        else:
            added = int(result.get("chunksAdded", 0))
            total += added
            if total % 100 == 0 or added > 0:
                print(f"  +{added:3d} ({total:4d} total)", flush=True)

        time.sleep(0.05)

    return total, total_errors


def main():
    subjects = (
        sys.argv[1:]
        if len(sys.argv) > 1
        else ["biology", "chemistry", "math", "physics"]
    )
    grand_total = 0
    grand_errors = 0

    for subj in subjects:
        rag_path = RAG_DIR / f"{subj}.yaml"
        if not rag_path.exists():
            print(f"[{subj}] YAML not found, skipping", flush=True)
            continue

        print(f"\n=== {subj.upper()} ===", flush=True)
        total, errors = process_rag_file(rag_path, subj)
        grand_total += total
        grand_errors += errors
        print(f"[{subj}] Done: {total} added, {errors} errors", flush=True)
        time.sleep(1)

    status = curl_post(f"{WORKER_URL}/rag/status", {})
    print(f"\nTotal vectors in Upstash: {status.get('total_vectors', '?')}", flush=True)
    print(f"Grand total: {grand_total} added, {grand_errors} errors", flush=True)


if __name__ == "__main__":
    main()
