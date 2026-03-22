#!/usr/bin/env python3
"""
Build partial RAGs from BAC Notes vault using Jina API embeddings.
Builds 4 subject RAGs: math, physics, chemistry, biology.
"""

import os, json, time, yaml, hashlib, subprocess, base64
from pathlib import Path

JINA_API = "https://api.jina.ai/v1/embeddings"
JINA_KEY = "jina_7ae28f7c2b08498a966f7cb7657d8898jXBBZ07B3TmTm_YNMqwIAGxy0gm9"
EMBED_MODEL = "jina-embeddings-v2-base-en"
VAULT = Path("/home/med/Documents/bac/resources/notes")
OUT_DIR = Path.home() / ".config/aichat/rags"
CHUNK_SIZE = 2000
CHUNK_OVERLAP = 200
BATCH_SIZE = 25
SLEEP_BETWEEN_BATCHES = 1

# =========== SUBJECTS ===========
SUBJECTS = {
    "math": {
        "paths": [
            "01-Concepts/Math/**/*.md",
            "02-Practice/Math/**/*.md",
            "02-Practice/FromTextbooks/Math/**/*.md",
            "05-Extracted/Math/**/*.md",
            "05-Extracted/Resources/Math/**/*.md",
            "03-Resources/Math/**/*.md",
        ]
    },
    "physics": {
        "paths": [
            "01-Concepts/Physics/**/*.md",
            "02-Practice/Physics/**/*.md",
            "02-Practice/FromTextbooks/Physics/**/*.md",
            "05-Extracted/Physics/**/*.md",
            "05-Extracted/Resources/Physics/**/*.md",
            "03-Resources/Physics/**/*.md",
        ]
    },
    "chemistry": {
        "paths": [
            "01-Concepts/Chemistry/**/*.md",
            "02-Practice/Chemistry/**/*.md",
            "02-Practice/FromTextbooks/Chemistry/**/*.md",
            "05-Extracted/Chemistry/**/*.md",
            "05-Extracted/Resources/Chemistry/**/*.md",
            "03-Resources/Chemistry/**/*.md",
        ]
    },
    "biology": {
        "paths": [
            "01-Concepts/Biology/**/*.md",
            "02-Practice/Biology/**/*.md",
            "05-Extracted/Biology/**/*.md",
            "01-Concepts/Natural-Sciences/**/*.md",
            "02-Practice/Natural-Sciences/**/*.md",
            "04-Exams/**/*.md",
        ]
    },
}


# =========== HELPERS ===========
def estimate_tokens(chars: int) -> int:
    return max(1, chars // 4)


def chunk_text(text: str, size: int = CHUNK_SIZE, overlap: int = CHUNK_OVERLAP):
    if not text.strip():
        return []
    start = 0
    chunks = []
    while start < len(text):
        end = start + size
        chunk = text[start:end]
        chunks.append(chunk)
        start = end - overlap
    return chunks


def file_hash(path: Path) -> str:
    m = hashlib.sha256()
    m.update(path.read_bytes())
    return m.hexdigest()[:16]


def encode_vector(vec: list[float]) -> str:
    """Encode float32 vector as base64 (aichat format)."""
    import numpy as np

    arr = np.array(vec, dtype=np.float32)
    return base64.b64encode(arr.tobytes()).decode("ascii")


def embed_batch(texts: list[str], retries: int = 5) -> list[list[float]]:
    payload = json.dumps({"model": EMBED_MODEL, "input": texts})
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
                    "Authorization: Bearer " + JINA_KEY,
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


def process_file(filepath: Path, file_id: int):
    """Process a single file: chunk it."""
    try:
        text = filepath.read_text(encoding="utf-8", errors="ignore")
        chunks = chunk_text(text)
        return {
            "file_id": file_id,
            "path": str(filepath.relative_to(VAULT)),
            "full_path": str(filepath),
            "hash": file_hash(filepath),
            "chunks": chunks,
        }
    except Exception as e:
        return None


def collect_files(patterns: list[str]) -> list[Path]:
    """Collect all matching files."""
    files = []
    for pattern in patterns:
        for p in VAULT.glob(pattern):
            if p.is_file():
                files.append(p)
    return sorted(set(files))


# =========== MAIN BUILD ===========
def build_rag(subject: str, config: dict):
    print(f"\n{'=' * 50}")
    print(f"Building RAG: {subject}")
    print(f"{'=' * 50}")

    # Collect files
    all_files = []
    for pattern in config["paths"]:
        files = collect_files([pattern])
        print(f"  Pattern {pattern}: {len(files)} files")
        all_files.extend(files)
    all_files = sorted(set(all_files))
    print(f"  Total files: {len(all_files)}")

    if not all_files:
        print(f"  No files found for {subject}!")
        return

    # Process files (chunking)
    print(f"  Chunking {len(all_files)} files...")
    file_data = []
    next_file_id = 1
    for f in all_files:
        result = process_file(f, next_file_id)
        if result and result["chunks"]:
            file_data.append(result)
            next_file_id += 1

    total_chunks = sum(len(d["chunks"]) for d in file_data)
    total_chars = sum(sum(len(c) for c in d["chunks"]) for d in file_data)
    total_tokens = estimate_tokens(total_chars)
    print(f"  Files with content: {len(file_data)}")
    print(f"  Total chunks: {total_chunks}")
    print(f"  Total chars: {total_chars:,}")
    print(f"  Est. tokens: {total_tokens:,}")
    est_rate = 4.6  # chunks/s via Jina
    print(f"  Est. time: ~{total_chunks / est_rate / 60:.0f}min")

    # Flatten all chunks with metadata
    all_chunks = []
    for fd in file_data:
        for ci, chunk in enumerate(fd["chunks"]):
            all_chunks.append(
                {
                    "file_id": fd["file_id"],
                    "path": fd["path"],
                    "chunk_index": ci,
                    "text": chunk,
                    "tokens": estimate_tokens(len(chunk)),
                }
            )

    # Embed in batches via Ollama
    print(f"  Embedding {len(all_chunks)} chunks...")
    vectors = {}
    batch = []
    processed = 0
    start_time = time.time()

    def flush_batch():
        nonlocal batch, processed
        if not batch:
            return
        texts = [item["text"] for item in batch]
        embeddings = embed_batch(texts)
        for item, emb in zip(batch, embeddings):
            chunk_id = f"{item['file_id']}-{item['chunk_index']}"
            vectors[chunk_id] = emb
        elapsed = time.time() - start_time
        rate = processed / elapsed if elapsed > 0 else 0
        remaining = len(all_chunks) - processed
        eta = remaining / rate if rate > 0 else 0
        print(
            f"    {processed}/{len(all_chunks)} ({100 * processed / len(all_chunks):.0f}%) - {rate:.0f}/s - ETA: {eta:.0f}s"
        )
        batch = []
        time.sleep(SLEEP_BETWEEN_BATCHES)

    for chunk in all_chunks:
        if len(batch) >= BATCH_SIZE:
            flush_batch()
        batch.append(chunk)
        processed += 1

    if batch:
        flush_batch()

    print(
        f"  Embedded {len(vectors)}/{len(all_chunks)} chunks in {time.time() - start_time:.0f}s"
    )

    # Build aichat RAG YAML
    print(f"  Building RAG YAML...")
    rag_files = []
    for fd in file_data:
        doc_chunks = []
        for ci, chunk_text in enumerate(fd["chunks"]):
            chunk_id = f"{fd['file_id']}-{ci}"
            if chunk_id in vectors:
                doc_chunks.append(
                    {
                        "page_content": chunk_text,
                        "metadata": {"source": fd["path"], "chunk": ci},
                    }
                )

        rag_files.append(
            {"hash": fd["hash"], "path": "./" + fd["path"], "documents": doc_chunks}
        )

    # Build vectors dict - encode as base64 strings
    rag_vectors = {}
    for chunk_id, vec in vectors.items():
        rag_vectors[chunk_id] = encode_vector(vec)

    rag = {
        "embedding_model": "jina:jina-embeddings-v2-base-en",
        "chunk_size": CHUNK_SIZE,
        "chunk_overlap": CHUNK_OVERLAP,
        "reranker_model": None,
        "top_k": 5,
        "batch_size": 25,
        "next_file_id": next_file_id,
        "document_paths": [
            str(VAULT / p.replace("**/*", "").rstrip("/")) for p in config["paths"]
        ],
        "files": {
            i + 1: f for i, f in enumerate(rag_files)
        },  # 1-indexed to match next_file_id
        "vectors": rag_vectors,
    }

    out_path = OUT_DIR / f"{subject}.yaml"
    with open(out_path, "w") as f:
        yaml.dump(rag, f, allow_unicode=True, default_flow_style=False, sort_keys=False)

    size_kb = os.path.getsize(out_path) // 1024
    print(f"  ✅ Saved: {out_path} ({size_kb} KB)")
    print(f"  Files: {len(rag_files)}, Vectors: {len(rag_vectors)}")


if __name__ == "__main__":
    import sys

    if len(sys.argv) > 1:
        subject = sys.argv[1]
        if subject in SUBJECTS:
            build_rag(subject, SUBJECTS[subject])
        else:
            print(f"Unknown subject: {subject}")
            print(f"Available: {list(SUBJECTS.keys())}")
    else:
        # Build all in sequence
        for subject, config in SUBJECTS.items():
            build_rag(subject, config)
            time.sleep(5)  # brief pause between subjects
