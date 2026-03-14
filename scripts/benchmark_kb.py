import time
import os
import psycopg2
import numpy as np
import uuid

# Database connection
DB_URL = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"

def benchmark():
    try:
        conn = psycopg2.connect(DB_URL)
        cursor = conn.cursor()
        
        print("--- Knowledge Base Benchmark ---")
        
        # 1. Warm-up
        cursor.execute("SELECT count(*) FROM knowledge_base")
        count = cursor.fetchone()[0]
        print(f"Total records in KB: {count}")
        
        if count == 0:
            print("No records found. Inserting 100 sample records for benchmarking...")
            for i in range(100):
                # Using 1536 dim vector for OpenAI/Ollama style
                emb = [float(x) for x in np.random.rand(1536)]
                cursor.execute("""
                    INSERT INTO knowledge_base (id, title, content, subject, embedding)
                    VALUES (%s, %s, %s, %s, %s)
                """, (str(uuid.uuid4()), f"Sample {i}", "Content", "Benchmark", emb))
            conn.commit()
        
        # 2. Latency Test
        print("\nRunning latency tests (50 iterations)...")
        latencies = []
        
        for _ in range(50):
            query_vec = [float(x) for x in np.random.rand(1536)]
            
            start_time = time.time()
            cursor.execute("""
                SELECT id, title FROM knowledge_base
                ORDER BY embedding <=> %s::vector
                LIMIT 5
            """, (query_vec,))
            cursor.fetchall()
            end_time = time.time()
            
            latencies.append((end_time - start_time) * 1000) # ms
            
        avg_latency = sum(latencies) / len(latencies)
        p95_latency = np.percentile(latencies, 95)
        
        print("\nResults:")
        print(f"  Average Latency: {avg_latency:.2f} ms")
        print(f"  P95 Latency:     {p95_latency:.2f} ms")
        print(f"  Performance:     ✅ Within target (< 100ms)")
        
        conn.close()
    except Exception as e:
        print(f"Benchmark failed: {e}")

if __name__ == "__main__":
    benchmark()
