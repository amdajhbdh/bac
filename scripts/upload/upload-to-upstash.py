#!/usr/bin/env python3
"""
Upload vault RAG content to Upstash Vector.
Re-embeds all vault content with jina-embeddings-v3 (1024-dim) and upserts to Upstash.
"""

import os, json, time, yaml, subprocess, sys
from pathlib import Path

JINA_API = "https://api.jina.ai/v1/embeddings"
JINA_KEY = "jina_7ae28f7c2b08498a966f7cb7657d8898jXBBZ07B3TmTm_YNMqwIAGxy0gm9"
JINA_MODEL = "jina-embeddings-v3"
UPSTASH_URL = "https://proud-sawfly-58811-us1-vector.upstash.io"
UPSTASH_TOKEN = "ABYFMHByb3VkLXNhd2ZseS01ODgxMS11czFhZG1pbll6RTRaRE00Wm1ZdFltUXdOQzAwWkRRd0xUa3lPRFl0TTJabE16TTJNRGRqTlRVMg=="
VAULT = Path("/home/med/Documents/bac/resources/notes")
RAG_DIR = Path.home() / ".config/aichat/rags"
BATCH_SIZE = 10
SLEEP = 5.0


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


def process_rag_file(filepath, subject):
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

    return all_chunks


def main():
    subjects_arg = (
        sys.argv[1:]
        if len(sys.argv) > 1
        else ["biology", "chemistry", "math", "physics"]
    )
    subjects = {s: s for s in subjects_arg}
    print(f"Uploading subjects: {list(subjects.keys())}")

    all_chunks = []
    for subj in subjects:
        rag_path = RAG_DIR / f"{subj}.yaml"
        if not rag_path.exists():
            print(f"⚠ {subj}.yaml not found, skipping")
            continue
        print(f"\nProcessing {subj}...")
        chunks = process_rag_file(rag_path, subj)
        print(f"  Found {len(chunks)} chunks")
        all_chunks.extend(chunks)

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

    # Upsert to Upstash in batches of 100
    print(f"\nUploading to Upstash...")
    upstash_batch_size = 100
    for i in range(0, len(vectors_to_upsert), upstash_batch_size):
        batch = vectors_to_upsert[i : i + upstash_batch_size]
        print(
            f"  Uploading {i + 1}-{min(i + upstash_batch_size, len(vectors_to_upsert))}..."
        )
        result = upsert_to_upstash(batch)
        print(f"  Result: {result}")
        time.sleep(0.5)

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
