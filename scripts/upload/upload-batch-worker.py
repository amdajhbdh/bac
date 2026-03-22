#!/usr/bin/env python3
"""
Upload BAC RAG content to Cloudflare Worker efficiently.
Handles large batches and includes proper error handling.
"""

import json, time, yaml, sys
from pathlib import Path

WORKER_URL = "https://bac-api.amdajhbdh.workers.dev"
RAG_DIR = Path.home() / ".config/aichat/rags"


def upload_batch(chunks_data, subject: str, path: str) -> dict:
    """Upload a batch of chunks to Worker endpoint."""
    data = json.dumps(
        {"chunks": chunks_data, "subject": subject, "path": path}
    ).encode()

    try:
        import urllib.request

        req = urllib.request.Request(
            f"{WORKER_URL}/rag/add",
            data=data,
            headers={"Content-Type": "application/json"},
            method="POST",
        )
        with urllib.request.urlopen(req, timeout=120) as r:
            return json.loads(r.read())
    except Exception as e:
        return {"error": str(e)}


def process_rag_file(filepath: Path, subject: str):
    with open(filepath) as f:
        rag = yaml.safe_load(f)

    total_uploaded = 0
    chunks_processed = 0

    for fid, fdata in rag["files"].items():
        path = fdata["path"]
        if path.startswith("./"):
            path = path[2:]

        # Collect all chunks from this file
        file_chunks = []
        for doc in fdata.get("documents", []):
            if doc["page_content"].strip():
                file_chunks.append(doc["page_content"])
                chunks_processed += 1

        # Upload in batches of 20 chunks
        batch_size = 20
        for i in range(0, len(file_chunks), batch_size):
            batch = file_chunks[i : i + batch_size]

            # Prepare batch data
            chunks_data = []
            for chunk in batch:
                chunks_data.append({"text": chunk, "chunk_index": i + len(chunks_data)})

            print(
                f"  Uploading batch {i // batch_size + 1} ({len(batch)} chunks) from {path}"
            )

            result = upload_batch(chunks_data, subject, path)

            if "error" in result:
                print(f"    ⚠ Error: {result['error']}")
            else:
                uploaded = result.get("chunksAdded", 0)
                total_uploaded += uploaded
                print(f"    ✅ Uploaded {uploaded} chunks")

        # Brief pause between files
        time.sleep(1)

    return total_uploaded, chunks_processed


def main():
    subjects_arg = sys.argv[1:] if len(sys.argv) > 1 else ["math"]

    for subj in subjects_arg:
        rag_path = RAG_DIR / f"{subj}.yaml"
        if not rag_path.exists():
            print(f"⚠ {subj}.yaml not found, skipping")
            continue

        print(f"\nProcessing {subj}...")
        uploaded, processed = process_rag_file(rag_path, subj)
        print(f"  Done: {uploaded} chunks uploaded out of {processed} total")
        time.sleep(2)

    # Check final status
    try:
        import urllib.request

        req = urllib.request.Request(f"{WORKER_URL}/rag/status")
        with urllib.request.urlopen(req, timeout=30) as r:
            status = json.loads(r.read())
        print(f"\nTotal vectors in Upstash: {status.get('total_vectors', '?')}")
    except Exception as e:
        print(f"\nCould not check status: {e}")


if __name__ == "__main__":
    main()
