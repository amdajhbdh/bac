#!/usr/bin/env python3
"""
Optimized parallel upload to Upstash Vector with dynamic batch sizing
Uses Cloudflare Workers for high-throughput uploads
"""
import os
import sys
import json
import time
import yaml
import threading
import urllib.request
from pathlib import Path
# Configuration
WORKER_URL = "https://bac-api.amdajhbdh.workers.dev"
VAULT = Path("/home/med/Documents/bac/resources/notes")
RAG_DIR = Path.home() / ".config/aichat/rags"
BATCH_SIZE = 5  # Chunks per request
MAX_WORKERS = 10  # Parallel upload workers
def upload_chunk(chunk_data, subject, path, worker_id):
    """Upload single chunk with worker identification"""
    tries = 3
    delay = 2  # Initial delay in seconds
    
    for attempt in range(tries):
        try:
            payload = {
                "text": chunk_data["text"],
                "metadata": {
                    "subject": subject,
                    "path": path,
                    "chunk_index": chunk_data["chunk_index"],
                    "worker_id": worker_id,
                    "attempt": attempt
                }
            }
            
            req = urllib.request.Request(
                f"{WORKER_URL}/rag/add",
                data=json.dumps(payload).encode(),
                headers={"Content-Type": "application/json"},
                method="POST"
            )
            
            with urllib.request.urlopen(req, timeout=30) as r:
                if r.getcode() == 200:
                    print(f"Worker {worker_id}: Uploaded chunk {chunk_data['chunk_index']}")
                    return True
            delay *= 2
            time.sleep(delay)
            
        except Exception as e:
            print(f"Worker {worker_id}: Upload failed for chunk {chunk_data['chunk_index']} - {str(e)}")
            delay *= 2
            time.sleep(delay)
    
    return False
def process_rag_file(subject):
    """Process and upload chunks from a single RAG file"""
    rag_path = RAG_DIR / f"{subject}.yaml"
    if not rag_path.exists():
        print(f"⚠ {subject}.yaml not found, skipping")
        return 0, 0
    print(f"\nProcessing {subject}...")
    with open(rag_path) as f:
        rag = yaml.safe_load(f)
    total_uploaded = 0
    chunks_processed = 0
    # Collect all chunks first
    all_chunks = []
    for fid, fdata in rag["files"].items():
        path = fdata["path"]
        if path.startswith("./"):
            path = path[2:]
        
        for doc in fdata.get("documents", []):
            if content := doc.get("page_content", "").strip():
                all_chunks.append({
                    "text": content,
                    "chunk_index": len(all_chunks)  # Recalculate indices
                })
                chunks_processed += 1
    print(f"Found {len(all_chunks)} chunks to upload for {subject}")
    
    # Dynamic batch sizing
    batch_size = BATCH_SIZE
    if chunks_processed > 1000:
        batch_size = max(2, BATCH_SIZE // 2)
    elif chunks_processed > 100:
        batch_size = max(3, BATCH_SIZE)
    # Upload in parallel batches
    worker_id = 1
    vector_batches = [all_chunks[i:i + batch_size] for i in range(0, len(all_chunks), batch_size)]
    
    def worker_upload(worker_id, vector_batch):
        """Upload a batch of vectors with worker identification"""
        batch_payload = {
            "subject": subject,
            "path": vector_batch[0]["path"] if vector_batch else "",
            "chunks": vector_batch,
            "worker_id": worker_id
        }
        req = urllib.request.Request(
            f"{WORKER_URL}/rag/add",
            data=json.dumps(batch_payload).encode(),
            headers={"Content-Type": "application/json"},
            method="POST"
        )
        tries = 3
        delay = 2
        success = False
        
        for attempt in range(tries):
            try:
                with urllib.request.urlopen(req, timeout=45) as r:
                    if r.getcode() == 200:
                        success = True
                        return r.json()["chunksAdded"]
            except Exception as e:
                if attempt == tries - 1:
                    print(f"Worker {worker_id}: Batch upload attempt {attempt+1} failed - {str(e)}")
                    time.sleep(delay)
                    delay *= 2
            except KeyboardInterrupt:
                print("Upload interrupted by user")
                return 0
        
        print(f"Worker {worker_id}: Batch upload failed after {tries} attempts")
        return 0
    # Start parallel upload workers
    def worker_func(worker_id, vector_batch):
        uploaded = worker_upload(worker_id, vector_batch)
        nonlocal total_uploaded
        total_uploaded += uploaded
    threads = []
    for i, batch in enumerate(vector_batches):
        thread = threading.Thread(target=worker_func, args=(worker_id + i, batch))
        threads.append(thread)
        thread.start()
    for thread in threads:
        thread.join()
    print(f"Completed upload of {subject} with {total_uploaded}/{chunks_processed} chunks")
    return total_uploaded, chunks_processed
def main():
    """Main execution function"""
    subjects = sys.argv[1:] if len(sys.argv) > 1 else ["math", "physics"]
    print(f"Uploading subjects: {subjects}")
    total_uploaded = 0
    total_chunks = 0
    for subj in subjects:
        uploaded, processed = process_rag_file(subj)
        total_uploaded += uploaded
        total_chunks += processed
    print(f"\nTotal uploaded: {total_uploaded} out of {total_chunks} chunks")
if __name__ == "__main__":
    main()
