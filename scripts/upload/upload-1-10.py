#!/usr/bin/env python3
"""
Upload 1/10 of BAC RAG content to Upstash Vector to avoid rate limiting.
Only processes one subject at a time, uploads 1/10 of chunks.
"""

import os, json, time, yaml, subprocess, sys
from pathlib import Path

JINA_API = "https://api.jina.ai/v1/embeddings"
JINA_KEY = "jina_7ae28f7c2b08498a966f7cb7657d8898jXBBZ07B3TmTm_YNMqwIAGxy0gm9"
JINA_MODEL = "jina-embeddings-v3"
UPSTASH_URL = "https://proud-sawfly-58811-us1-vector.upstash.io"
UPSTASH_TOKEN = "ABYFMHByb3VkLXNhd2ZseS01ODgxMS11czFhZG1pbll6RTRaRE00Wm1ZdFltUXdOQzAwWkRRd0xUa3lPRFl0s"
VAULT = Path("/home/med/Documents/bac/resources/notes")
RAG_DIR = Path.home() / ".config/aichat/rags"
BATCH_SIZE = 25
SLEEP = 1.0


def embed_batch(texts, retries=5):
    payload = json.dumps({"model": JINA_MODEL, "input": texts})
    for attempt in range(retries):
        try:
            result = subprocess.run(
                [
                    "curl",
                    "-s",
                    "-X",
                    "POST",
                    JINA_API,
                    "-H",
                    f"Authorization: Bearer {JINA_KEY}",
                    "-H",
                    "Content-Type: application/json",
                    "-d",
                    payload,
                    "--max-time",
                    "30",
                ],
                capture_output=True,
                text=True,
                timeout=35,
            )
            if result.returncode != 0:
                raise Exception(f"curl failed: {result.stderr}")
            data = json.loads(result.stdout)
            if "data" not in data:
                raise Exception(f"No data: {data}")
            embeddings = sorted(data["data"], key=lambda x: x["index"])
            return [e["embedding"] for e in embeddings]
        except Exception as e:
            wait = (2**attempt) * 2
            print(f"    ⚠ Attempt {attempt + 1} failed: {e}, retrying in {wait}s...")
            time.sleep(wait)
    raise Exception("All retries failed")


def upsert_to_upstash(vectors, retries=5):
    payload = json.dumps({"vectors": vectors})
    for attempt in range(retries):
        try:
            result = subprocess.run(
                [
                    "curl",
                    "-s",
                    "-X",
                    "POST",
                    f"{UPSTASH_URL}/upsert",
                    "-H",
                    f"Authorization: Bearer {UPSTASH_TOKEN}",
                    "-H",
                    "Content-Type: application/json",
                    "-d",
                    payload,
                    "--max-time",
                    "60",
                ],
                capture_output=True,
                text=True,
                timeout=65,
            )
            if result.returncode != 0:
                raise Exception(f"Upsert failed: {result.stderr}")
            data = json.loads(result.stdout)
            if data.get("status") == "success":
                return data
            if "keys" in data:
                return data
            raise Exception(f"Unexpected response: {data}")
        except Exception as e:
            wait = (2**attempt) * 3
            print(
                f"    ⚠ Upsert attempt {attempt + 1} failed: {e}, retrying in {wait}s..."
            )
            time.sleep(wait)
    raise Exception("All upsert retries failed")


def process_rag_file(filepath, subject, fraction=0.1):
    with open(filepath) as f:
        rag = yaml.safe_load(f)

    all_chunks = []
    for fid, fdata in rag["files"].items():
        path = fdata["path"]
        if path.startswith("./"):
            path = path[2:]

        for ci, doc in enumerate(fdata.get("documents", [])):
            chunk_text_content = doc["page_content"]
            if not chunk_text_content.strip():
                continue
            metadata = doc.get("metadata", {})
            all_chunks.append(
                {
                    "file_id": str(fid),
                    "path": path,
                    "chunk_index": ci,
                    "text": chunk_text_content,
                    "subject": subject,
                    "source": metadata.get("source", path),
                }
            )

    # Only process 1/10 of chunks to avoid rate limiting
    total_chunks = len(all_chunks)
    sample_size = max(1, int(total_chunks * fraction))
    sampled_chunks = all_chunks[:sample_size]

    print(f"  Subject: {subject}")
    print(f"  Total chunks: {total_chunks}")
    print(f"  Processing: {sample_size} chunks (1/10)")

    return sampled_chunks


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 upload-1-10.py <subject>")
        print("Available subjects: biology, chemistry, math, physics")
        sys.exit(1)

    subject = sys.argv[1]
    rag_path = RAG_DIR / f"{subject}.yaml"

    if not rag_path.exists():
        print(f"⚠ {subject}.yaml not found, skipping")
        sys.exit(1)

    print(f"\nProcessing {subject}...")
    all_chunks = process_rag_file(rag_path, subject)

    if not all_chunks:
        print("No chunks to process")
        sys.exit(0)

    print(f"\nTotal chunks to embed: {len(all_chunks)}")

    vectors_to_upsert = []
    start = time.time()
    batch = []

    def flush():
        nonlocal batch, vectors_to_upsert
        if not batch:
            return
        texts = [item["text"] for item in batch]
        embeddings = embed_batch(texts)

        for item, emb in zip(batch, embeddings):
            vectors_to_upsert.append(
                {
                    "id": f"{item['subject']}-{item['file_id']}-{item['chunk_index']}",
                    "vector": emb,
                    "metadata": {
                        "subject": item["subject"],
                        "path": item["path"],
                        "chunk_index": item["chunk_index"],
                        "text": item["text"],
                    },
                }
            )

        elapsed = time.time() - start
        rate = len(vectors_to_upsert) / elapsed if elapsed > 0 else 0
        remaining = len(all_chunks) - len(vectors_to_upsert)
        eta = remaining / rate if rate > 0 else 0
        print(
            f"  {len(vectors_to_upsert)}/{len(all_chunks)} ({100 * len(vectors_to_upsert) / len(all_chunks):.0f}%) - {rate:.0f}/s - ETA: {eta:.0f}s"
        )
        batch = []
        time.sleep(SLEEP)

    for i, chunk in enumerate(all_chunks):
        batch.append(chunk)
        if len(batch) >= BATCH_SIZE:
            flush()

    if batch:
        flush()

    print(f"\nEmbedded {len(vectors_to_upsert)} chunks in {time.time() - start:.0f}s")

    # Upsert to Upstash in small batches to avoid rate limiting
    print(f"\nUploading to Upstash...")
    upstash_batch_size = 10  # Small batches to avoid rate limiting
    for i in range(0, len(vectors_to_upsert), upstash_batch_size):
        batch = vectors_to_upsert[i : i + upstash_batch_size]
        print(
            f"  Uploading {i + 1}-{min(i + upstash_batch_size, len(vectors_to_upsert))}..."
        )
        result = upsert_to_upstash(batch)
        print(f"  Result: {result}")
        time.sleep(1.0)  # Extra sleep between batches

    # Verify
    print(f"\nVerifying...")
    result = subprocess.run(
        [
            "curl",
            "-s",
            "-X",
            "POST",
            f"{UPSTASH_URL}/info",
            "-H",
            f"Authorization: Bearer {UPSTASH_TOKEN}",
            "-H",
            "Content-Type: application/json",
        ],
        capture_output=True,
        text=True,
        timeout=15,
    )
    info = json.loads(result.stdout)
    print(f"Upstash info: {json.dumps(info.get('result', info), indent=2)}")


if __name__ == "__main__":
    main()
