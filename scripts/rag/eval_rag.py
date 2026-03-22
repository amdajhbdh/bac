#!/usr/bin/env python3
"""
RAG Evaluation Script for BAC Study System.
Measures: coverage, retrieval quality, answer accuracy, topic gaps.

Usage:
    python3 eval_rag.py                    # full eval
    python3 eval_rag.py --coverage        # vault coverage only
    python3 eval_rag.py --retrieval       # retrieval quality only
    python3 eval_rag.py --accuracy         # answer accuracy only
    python3 eval_rag.py --subject biology  # eval specific subject
    python3 eval_rag.py --verbose         # show details
"""

import argparse
import json
import subprocess
import sys
import time
import re
import yaml
from collections import defaultdict
from pathlib import Path

WORKER = "https://bac-api.amdajhbdh.workers.dev"
VAULT_DIR = Path.home() / ".config/aichat/rags"

JINA_KEY = ""
OPENROUTER_KEY = ""

try:
    import os as _os

    cfg_path = Path.home() / ".config/aichat/config.yaml"
    if cfg_path.exists():
        cfg = cfg_path.read_text()
        m = re.search(r"api_key:\s*(jina_[^\s\n]+)", cfg)
        if m:
            JINA_KEY = m.group(1).strip()
        or_key = re.search(r"sk-or-v1-[a-f0-9]+", cfg)
        if or_key:
            OPENROUTER_KEY = or_key.group(0)
    if not JINA_KEY:
        JINA_KEY = _os.environ.get("JINA_API_KEY", "")
    if not OPENROUTER_KEY:
        OPENROUTER_KEY = _os.environ.get("OPENROUTER_API_KEY", "")
except Exception:
    import os as _os

    JINA_KEY = _os.environ.get("JINA_API_KEY", "")
    OPENROUTER_KEY = _os.environ.get("OPENROUTER_API_KEY", "")


def embed_jina(texts, model="jina-embeddings-v2-base-en"):
    """Embed texts using Jina API. Returns 768-dim vectors."""
    if not JINA_KEY:
        return None
    try:
        resp = subprocess.run(
            [
                "curl",
                "-s",
                "-X",
                "POST",
                "https://api.jina.ai/v1/embeddings",
                "-H",
                "Content-Type: application/json",
                "-H",
                f"Authorization: Bearer {JINA_KEY}",
                "-d",
                json.dumps({"model": model, "input": texts}),
            ],
            capture_output=True,
            text=True,
            timeout=60,
        )
        data = json.loads(resp.stdout)
        embeddings = data.get("data", [])
        return [e["embedding"] for e in sorted(embeddings, key=lambda x: x["index"])]
    except Exception as e:
        print(f"   Jina embed failed: {e}")
        return None


def rerank_jina(query, documents, top_n=5):
    """Rerank documents using Jina reranker."""
    if not JINA_KEY or not documents:
        return documents[:top_n]
    try:
        resp = subprocess.run(
            [
                "curl",
                "-s",
                "-X",
                "POST",
                "https://api.jina.ai/v1/rerank",
                "-H",
                "Content-Type: application/json",
                "-H",
                f"Authorization: Bearer {JINA_KEY}",
                "-d",
                json.dumps(
                    {
                        "model": "jina-reranker-v1-base-en",
                        "query": query,
                        "documents": documents,
                        "top_n": top_n,
                    }
                ),
            ],
            capture_output=True,
            text=True,
            timeout=60,
        )
        data = json.loads(resp.stdout)
        indices = [r["index"] for r in data.get("results", [])]
        return [documents[i] for i in indices if i < len(documents)]
    except:
        return documents[:top_n]


def curl(method, path, data=None, timeout=120):
    cmd = ["curl", "-s", "-X", method, f"{WORKER}{path}"]
    if data:
        cmd.extend(["-H", "Content-Type: application/json", "-d", json.dumps(data)])
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout + 10)
    try:
        return json.loads(result.stdout)
    except:
        return {}


def load_vault_manifest():
    """Parse YAML RAG files (files: {id: {path, documents: [{page_content}]}}) to get counts."""
    subjects = {}
    total_files = 0
    total_chunks = 0

    for yaml_path in sorted(VAULT_DIR.glob("*.yaml")):
        subject = yaml_path.stem
        try:
            data = yaml.safe_load(yaml_path.read_text())
        except Exception:
            continue

        files_dict = data.get("files", {})
        file_count = len(files_dict)
        chunk_count = sum(len(fi.get("documents", [])) for fi in files_dict.values())

        if file_count > 0:
            subjects[subject] = {
                "files": file_count,
                "chunks": chunk_count,
                "yaml_path": str(yaml_path),
            }
            total_files += file_count
            total_chunks += chunk_count

    return subjects, total_files, total_chunks, 0


def get_indexed_counts():
    """Get indexed vectors and questions from Worker."""
    status = curl("GET", "/rag/status")
    track = curl("GET", "/rag/track")
    result = {
        "vectors": status.get("total_vectors", 0),
        "questions": track.get("stats", {}).get("total_questions", 0),
        "attempted": track.get("stats", {}).get("attempted", 0),
        "accuracy": track.get("stats", {}).get("accuracy", 0),
        "due_now": track.get("stats", {}).get("due_now", 0),
    }
    return result


def eval_coverage(manifest, indexed, verbose=False):
    """Evaluate vault coverage — how much is indexed vs. total."""
    print(f"\n{'=' * 60}")
    print("  COVERAGE ANALYSIS")
    print(f"{'=' * 60}")

    results = []
    for subj, info in manifest.items():
        expected = info["chunks"]
        indexed_subj = indexed.get(f"subject_{subj}", "?")
        if isinstance(indexed_subj, int):
            pct = min(100, round(indexed_subj / max(expected, 1) * 100, 1))
            status = "✅" if pct >= 90 else "🟡" if pct >= 50 else "🔴"
        else:
            pct = indexed_subj
            status = "❓"
        results.append((subj, expected, indexed_subj, pct, status))

    results.sort(key=lambda x: x[3] if isinstance(x[3], float) else 0)
    total_expected = sum(r[1] for r in results)
    total_indexed = sum(r[2] for r in results if isinstance(r[2], int))

    print(f"\n{'Subject':<15} {'Expected':>10} {'% of Total':>12}  Status")
    print("-" * 50)
    for subj, expected, idx, pct, status in results:
        vault_pct = round(expected / max(total_expected, 1) * 100, 1)
        print(f"{subj:<15} {expected:>10} {vault_pct:>11.1f}%  {status}")
    print("-" * 50)
    total_pct = round(indexed["vectors"] / max(total_expected, 1) * 100, 1)
    print(
        f"{'INDEXED':<15} {indexed['vectors']:>10} / {total_expected:<10} {total_pct:>7.1f}%"
    )
    remaining = total_expected - indexed["vectors"]
    print(f"  → {remaining} chunks remaining to upload")

    if verbose:
        for subj, info in manifest.items():
            print(
                f"\n  {subj}: {info['files']} files, {info['chunks']} chunks ({info.get('no_embed', 0)} excluded)"
            )


def eval_retrieval(subjects=None, verbose=False, sample=20):
    """Test query quality via /rag/query (OpenRouter fallback)."""
    print(f"\n{'=' * 60}")
    print("  RETRIEVAL QUALITY")
    print(f"{'=' * 60}")

    params = "limit=100"
    if subjects:
        params += f"&subject={subjects[0]}"

    qs = curl("GET", f"/rag/questions?{params}")
    questions = qs.get("questions", [])

    if not questions:
        print("No questions found.")
        return {}

    import random

    random.seed(42)
    sample_qs = random.sample(questions, min(sample, len(questions)))

    results = []
    for q in sample_qs:
        qid = q["id"]
        text = q["question_text"]
        topic = q.get("topic_tags", "") or "unknown"
        solution = q.get("solution_text", "") or ""

        r = curl("POST", "/rag/query", {"query": text}, timeout=120)
        answer = r.get("answer", "")
        sources = r.get("sources", [])
        fallback = r.get("fallback", False)

        if not answer:
            results.append(
                {"id": qid, "has_answer": False, "sources": 0, "topic": topic}
            )
            continue

        sol_words = set(w.lower() for w in re.findall(r"\w{3,}", solution.lower()))
        ans_words = set(w.lower() for w in re.findall(r"\w{3,}", answer.lower()))
        overlap = (
            len(sol_words & ans_words) / max(len(sol_words), 1) if sol_words else 1.0
        )

        results.append(
            {
                "id": qid,
                "has_answer": True,
                "sources": len(sources),
                "overlap": round(overlap * 100, 1),
                "topic": topic[:30],
                "fallback": fallback,
            }
        )

    answered = [r for r in results if r["has_answer"]]
    avg_overlap = (
        round(sum(r["overlap"] for r in answered) / len(answered), 1) if answered else 0
    )
    rag_pct = (
        round(
            sum(1 for r in answered if not r.get("fallback")) / len(answered) * 100, 1
        )
        if answered
        else 0
    )

    print(f"\n📈 Avg answer overlap:  {avg_overlap}%")
    print(f"📦 RAG mode:           {rag_pct}% (rest: OpenRouter fallback)")
    print(f"📝 Questions tested:   {len(results)}")
    print(f"✅ Answered:           {len(answered)}/{len(results)}")

    if verbose:
        print(f"\n{'ID':>5}  {'Overlap':>8}  {'Srcs':>5}  Mode     Topic")
        print("-" * 60)
        for r in sorted(results, key=lambda x: x.get("overlap", 0)):
            mode = "RAG" if not r.get("fallback") else "LLM"
            srcs = r.get("sources", "-")
            print(
                f"{r['id']:>5}  {str(r.get('overlap', '-')):>8}  {srcs:>5}  {mode:<7} {r['topic']}"
            )

    return {"avg_overlap": avg_overlap, "rag_pct": rag_pct}


def eval_accuracy(subjects=None, verbose=False, sample=15):
    """Test answer accuracy using OpenRouter-powered /rag/solve endpoint."""
    print(f"\n{'=' * 60}")
    print("  ANSWER ACCURACY")
    print(f"{'=' * 60}")
    print("  Uses OpenRouter (qwen/qwen3-coder) via /rag/solve")

    params = "limit=100"
    if subjects:
        params += f"&subject={subjects[0]}"

    qs = curl("GET", f"/rag/questions?{params}")
    questions = qs.get("questions", [])

    if not questions:
        print("No questions found.")
        return {}

    import random

    random.seed(42)
    sample_qs = [
        q for q in questions if q.get("solution_text") and len(q["solution_text"]) > 5
    ]
    sample_qs = random.sample(sample_qs, min(sample, len(sample_qs)))

    results = []
    for q in sample_qs:
        qid = q["id"]
        text = q["question_text"]
        solution = q["solution_text"]
        topic = q.get("topic_tags", "") or "unknown"

        r = curl("POST", "/rag/solve", {"question": text}, timeout=120)
        answer = r.get("solution", "")

        if not answer:
            results.append(
                {"id": qid, "score": 0, "topic": topic, "status": "no_answer"}
            )
            continue

        score = grade_vs_solution(answer, solution)
        results.append(
            {
                "id": qid,
                "score": score,
                "topic": topic,
                "answer": answer[:80],
                "solution": solution[:80],
            }
        )

    avg_score = (
        round(sum(r["score"] for r in results) / len(results), 1) if results else 0
    )
    good_pct = (
        round(sum(1 for r in results if r["score"] >= 4) / len(results) * 100, 1)
        if results
        else 0
    )
    medium_pct = (
        round(sum(1 for r in results if r["score"] >= 3) / len(results) * 100, 1)
        if results
        else 0
    )

    print(f"\n📈 Average score:    {avg_score}/5")
    print(f"✅ Good (>=4):       {good_pct}%")
    print(f"📊 Acceptable (>=3): {medium_pct}%")
    print(f"📝 Questions tested: {len(results)}")

    if verbose:
        print(f"\n{'ID':>5}  {'Score':>6}  Topic")
        print("-" * 55)
        for r in sorted(results, key=lambda x: x["score"]):
            print(f"{r['id']:>5}  {r['score']:>5.1f}/5  {r['topic'][:35]}")

    topic_scores = defaultdict(list)
    for r in results:
        topic_scores[r["topic"]].append(r["score"])

    weak = [
        (t, round(sum(s) / len(s), 1))
        for t, s in topic_scores.items()
        if sum(s) / len(s) < 2.5
    ]
    if weak:
        print(f"\n🔴 WEAK TOPICS (avg < 2.5):")
        for t, sc in sorted(weak, key=lambda x: x[1]):
            print(f"   {sc:>5.1f}/5  {t[:50]}")

    return {"avg_score": avg_score, "good_pct": good_pct, "medium_pct": medium_pct}


def grade_vs_solution(answer, solution):
    """Grade answer against solution (numeric + overlap)."""

    def digits(s):
        return set(re.findall(r"-?\d+(?:[.,]\d+)?", str(s)))

    a_digits = digits(answer)
    s_digits = digits(solution)
    if a_digits and s_digits:
        overlap = len(a_digits & s_digits) / max(len(s_digits), 1)
        if overlap > 0.7 and a_digits == s_digits:
            return 5.0
        if overlap > 0.5:
            return 4.0
        if overlap > 0.3:
            return 3.0
    wa = set(re.findall(r"\w{3,}", answer.lower()))
    ws = set(re.findall(r"\w{3,}", solution.lower()))
    if not wa or not ws:
        return 0.0
    overlap = len(wa & ws) / len(ws)
    if overlap > 0.7:
        return 4.5
    elif overlap > 0.5:
        return 3.5
    elif overlap > 0.3:
        return 2.5
    elif overlap > 0.1:
        return 1.5
    return 0.5


def eval_topic_gaps(manifest, indexed, verbose=False):
    """Identify topics with low coverage."""
    print(f"\n{'=' * 60}")
    print("  TOPIC GAPS")
    print(f"{'=' * 60}")

    total = sum(info["chunks"] for info in manifest.values())
    remaining = total - indexed["vectors"]
    indexed_pct = round(indexed["vectors"] / max(total, 1) * 100, 1)

    gaps = []
    for subj, info in manifest.items():
        vault_pct = round(info["chunks"] / max(total, 1) * 100, 1)
        gaps.append((subj, info["files"], info["chunks"], vault_pct))

    gaps.sort(key=lambda x: x[2], reverse=True)

    print(f"\n📊 Total vault chunks: {total}")
    print(f"📦 Indexed:            {indexed['vectors']} ({indexed_pct}%)")
    print(f"❌ Remaining:          {remaining}")
    print(f"\n{'Subject':<15} {'Files':>7} {'Chunks':>8} {'Vault%':>8}")
    print("-" * 45)
    for subj, files, chunks, vault_pct in gaps:
        status = "🔴" if vault_pct > 30 else "🟡" if vault_pct > 10 else "🟢"
        print(f"{subj:<15} {files:>7} {chunks:>8} {vault_pct:>7.1f}%  {status}")

    print(f"\n💡 To upload remaining ~{remaining} chunks:")
    print(f"   python3 upload-via-worker.py physics math biology chemistry")


def eval_vector_diversity(questions_sample=30, verbose=False):
    """Check if queries return subject-appropriate answers."""
    print(f"\n{'=' * 60}")
    print("  SUBJECT ROUTING")
    print(f"{'=' * 60}")
    print("  Tests if the system correctly identifies the subject of each query")

    test_queries = [
        ("biology", "mitose cellule"),
        ("chemistry", "réaction chimique"),
        ("math", "équation polynôme"),
        ("physics", "force mouvement"),
        ("biology", "ADN génétique"),
        ("chemistry", "acide base pH"),
        ("math", "fonction dérivée"),
        ("physics", "énergie cinétique"),
    ]

    results = []
    for subj, query_text in test_queries:
        r = curl("POST", "/rag/query", {"query": query_text}, timeout=120)
        returned_subjs = r.get("subjects", [])
        correct = subj in returned_subjs
        results.append(
            {
                "query": query_text,
                "expected": subj,
                "returned": returned_subjs,
                "correct": correct,
            }
        )

    accuracy = (
        round(sum(1 for r in results if r["correct"]) / len(results) * 100, 1)
        if results
        else 0
    )
    print(
        f"\n📊 Subject accuracy: {accuracy}% ({sum(1 for r in results if r['correct'])}/{len(results)})"
    )

    if verbose:
        for r in results:
            status = "✅" if r["correct"] else "❌"
            print(f"  {status} [{r['expected']}] {r['query']:<25} → {r['returned']}")


def main():
    parser = argparse.ArgumentParser(description="BAC RAG Evaluation")
    parser.add_argument(
        "--coverage", action="store_true", help="Coverage analysis only"
    )
    parser.add_argument(
        "--retrieval", action="store_true", help="Retrieval quality only"
    )
    parser.add_argument("--accuracy", action="store_true", help="Answer accuracy only")
    parser.add_argument("--gaps", action="store_true", help="Topic gaps only")
    parser.add_argument(
        "--diversity", action="store_true", help="Vector diversity only"
    )
    parser.add_argument(
        "--subject", help="Filter by subject (math/biology/chemistry/physics)"
    )
    parser.add_argument("--verbose", "-v", action="store_true", help="Show details")
    parser.add_argument("--sample", type=int, default=20, help="Sample size for tests")
    args = parser.parse_args()

    print(f"\n{'=' * 60}")
    print("  BAC RAG EVALUATION")
    print(f"{'=' * 60}")

    run_all = not any(
        [args.coverage, args.retrieval, args.accuracy, args.gaps, args.diversity]
    )

    start = time.time()

    # Load data
    print("\n📂 Loading vault manifest...")
    manifest, total_files, total_chunks, no_embed = load_vault_manifest()
    print(
        f"   {len(manifest)} subjects, {total_files} files, {total_chunks} chunks ({no_embed} excluded)"
    )

    print("\n📡 Fetching Worker status...")
    indexed = get_indexed_counts()
    print(f"   Vectors: {indexed['vectors']}, Questions: {indexed['questions']}")

    subjects = [args.subject] if args.subject else None

    if run_all or args.coverage:
        eval_coverage(manifest, indexed, verbose=args.verbose)

    if run_all or args.gaps:
        eval_topic_gaps(manifest, indexed, verbose=args.verbose)

    if run_all or args.retrieval:
        eval_retrieval(subjects=subjects, verbose=args.verbose, sample=args.sample)

    if run_all or args.accuracy:
        eval_accuracy(subjects=subjects, verbose=args.verbose, sample=args.sample)

    if run_all or args.diversity:
        eval_vector_diversity(verbose=args.verbose)

    elapsed = time.time() - start
    print(f"\n{'=' * 60}")
    print(f"  Done in {elapsed:.0f}s")
    print(f"{'=' * 60}")


if __name__ == "__main__":
    main()
