#!/usr/bin/env python3
"""
BAC RAG CLI - Query the Worker API
Usage:
    python3 bac-rag.py query "question" [subject]
    python3 bac-rag.py solve "question" [subject]
    python3 bac-rag.py grade "question_id" "answer"
    python3 bac-rag.py practice [subject]
    python3 bac-rag.py status
"""

import argparse
import json
import subprocess
import sys

API = "https://bac-api.amdajhbdh.workers.dev"


def curl(method, path, data=None, timeout=120):
    cmd = ["curl", "-s", "-X", method, f"{API}{path}"]
    if data:
        cmd.extend(["-H", "Content-Type: application/json", "-d", json.dumps(data)])
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout + 10)
    try:
        return json.loads(result.stdout)
    except:
        return {"error": result.stdout[:200]}


def cmd_query(args):
    payload = {"query": args.query}
    if args.subject:
        payload["subjects"] = [args.subject]
    r = curl("POST", "/rag/query", payload)
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(r.get("answer", "No answer"))
    if r.get("sources"):
        print(f"\nSources: {len(r['sources'])} chunks ({r.get('model', '?')})")


def cmd_solve(args):
    payload = {"question": args.question}
    if args.subject:
        payload["subject"] = args.subject
    r = curl("POST", "/rag/solve", payload)
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(r.get("solution", "No solution"))


def cmd_grade(args):
    if args.preview:
        r = curl(
            "POST", "/rag/grade", {"question_id": args.question_id, "preview": True}
        )
        print(f"Solution: {r.get('solution', 'N/A')}")
        return
    r = curl(
        "POST",
        "/rag/grade",
        {"question_id": args.question_id, "user_answer": args.answer},
    )
    print(f"Score: {r.get('score', '?')}/5 {'✅' if r.get('correct') else '❌'}")
    print(f"Feedback: {r.get('feedback', '')}")


def cmd_practice(args):
    subject = args.subject
    params = f"?subject={subject}" if subject else ""
    r = curl("GET", f"/rag/practice{params}")
    q = r.get("question")
    if not q:
        print("No questions available")
        return
    print(f"Question #{q['id']}: {q['text']}")
    print(f'\nUsage: bac-rag grade {q["id"]} "your answer"')


def cmd_status(args):
    r = curl("GET", "/rag/status")
    print(f"Vectors: {r.get('total_vectors', '?')}")
    t = curl("GET", "/rag/track")
    s = t.get("stats", {})
    print(f"Questions: {s.get('attempted', 0)}/{s.get('total_questions', 0)}")
    print(f"Accuracy: {s.get('accuracy', 0)}%")


def main():
    p = argparse.ArgumentParser(description="BAC RAG CLI")
    sub = p.add_subparsers(dest="cmd", required=True)

    q = sub.add_parser("query", help="Query the RAG")
    q.add_argument("query", help="Question")
    q.add_argument(
        "subject", nargs="?", choices=["biology", "chemistry", "math", "physics"]
    )

    s = sub.add_parser("solve", help="Solve a question")
    s.add_argument("question", help="Question to solve")
    s.add_argument(
        "subject", nargs="?", choices=["biology", "chemistry", "math", "physics"]
    )

    g = sub.add_parser("grade", help="Grade an answer")
    g.add_argument("question_id", type=int)
    g.add_argument("answer", nargs="?")
    g.add_argument("--preview", action="store_true")

    p_cmd = sub.add_parser("practice", help="Get practice question")
    p_cmd.add_argument(
        "subject", nargs="?", choices=["biology", "chemistry", "math", "physics"]
    )

    sub.add_parser("status", help="System status")

    args = p.parse_args(sys.argv[1:] if len(sys.argv) > 1 else ["--help"])

    if args.cmd == "query":
        cmd_query(args)
    elif args.cmd == "solve":
        cmd_solve(args)
    elif args.cmd == "grade":
        cmd_grade(args)
    elif args.cmd == "practice":
        cmd_practice(args)
    elif args.cmd == "status":
        cmd_status(args)


if __name__ == "__main__":
    main()
