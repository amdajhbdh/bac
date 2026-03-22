#!/usr/bin/env python3
"""
BAC Study RAG CLI — query vault, practice, solve, track progress.
Usage: python3 bac-cli.py [command] [args]
"""

import sys, json, subprocess, argparse
from pathlib import Path
from typing import cast, Any

WORKER = "https://bac-api.amdajhbdh.workers.dev"


def curl(method, path, data=None, timeout=60):
    cmd = ["curl", "-s", "-X", method, f"{WORKER}{path}"]
    if data:
        cmd.extend(["-H", "Content-Type: application/json", "-d", json.dumps(data)])
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout + 10)
    try:
        return cast(dict[str, Any], json.loads(result.stdout))
    except json.JSONDecodeError:
        return cast(
            dict[str, Any], {"error": result.stdout[:200] or result.stderr[:200]}
        )


def cmd_query(args):
    data = {"query": args.query}
    if args.subjects:
        data["subjects"] = args.subjects.split(",")
    if args.model:
        data["model"] = args.model
    r = curl("POST", "/rag/query", data)
    if isinstance(r, str):
        print(f"Error: {r}")
        return
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"\n{'=' * 60}")
    print(f"Question: {args.query}")
    print(f"Subjects: {', '.join(r.get('subjects', []))}")
    print(f"{'=' * 60}")
    print(f"\n{r.get('answer', 'No answer')}\n")
    sources = r.get("sources")
    if sources:
        print("Sources:")
        for i, s in enumerate(sources, 1):
            print(f"  [{i}] {s['source']} (score: {s['score']:.2f})")


def cmd_search(args):
    data = {"query": args.query}
    if args.subjects:
        data["subjects"] = args.subjects.split(",")
    data["topK"] = args.top_k
    r = curl("POST", "/rag/search", data)
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"\n{'=' * 60}")
    print(f"Query: {args.query}")
    print(f"{'=' * 60}")
    if not r.get("chunks"):
        print("No results found.")
        return
    for i, c in enumerate(r["chunks"], 1):
        text = c["metadata"]["text"][:300]
        print(f"\n[{i}] {c['metadata']['subject']} | score: {c['score']:.3f}")
        print(f"    {text}...")


def cmd_solve(args):
    data = {"question": args.question}
    if args.subject:
        data["subject"] = args.subject
    print(f"Solving...")
    r = curl("POST", "/rag/solve", data, timeout=120)
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"\n{'=' * 60}")
    print(f"Question: {args.question}")
    print(f"{'=' * 60}")
    print(f"\n{r.get('solution', 'No solution')}\n")


def cmd_practice(args):
    params = ""
    if args.subject:
        params = f"?subject={args.subject}"
    r = curl("GET", f"/rag/practice{params}")
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    q = r.get("question")
    if not q:
        print(f"No questions available. {r.get('message', '')}")
        return
    print(f"\n{'=' * 60}")
    print(f"Question #{q['id']} [{'⭐' * q['difficulty']}]")
    print(f"{'=' * 60}")
    print(f"\n{q['text']}\n")
    print(f"Hint: {r.get('solution_preview', '')}")

    if not args.answer:
        print(f"\nUsage: bac-cli.py grade {q['id']} 'your answer'")
        print(f"       bac-cli.py grade {q['id']} --preview")
        return

    grade_r = curl(
        "POST", "/rag/grade", {"question_id": q["id"], "user_answer": args.answer}
    )
    print(f"\n{'=' * 60}")
    print(f"Grade: {grade_r.get('score', '?')}/5")
    print(f"FSRS Rating: {grade_r.get('rating', 'N/A')}")
    print(f"{'=' * 60}")
    print(f"\n{grade_r.get('feedback', '')}")
    print(f"\nSolution: {grade_r.get('solution', '')}")
    print(f"\nNext review: {grade_r.get('next_review', 'N/A')[:10]}")
    print(f"State: {grade_r.get('state', '')}")


def cmd_grade(args):
    chosen_rating = getattr(args, "rate", None) or getattr(args, "rating", None)
    if chosen_rating:
        r = curl(
            "POST",
            "/rag/rate",
            {"question_id": args.question_id, "rating": chosen_rating},
        )
        if "error" in r:
            print(f"Error: {r['error']}")
            return
        print(f"\n{'=' * 40}")
        print(f"  Rating: {chosen_rating}")
        print(f"  {r.get('message', '')}")
        print(f"  State: {r.get('state', '')}")
        print(f"{'=' * 40}")
        return
        print(f"\n{'=' * 40}")
        print(f"  Rating: {args.rate}")
        print(f"  {r.get('message', '')}")
        print(f"  State: {r.get('state', '')}")
        print(f"{'=' * 40}")
        return
    if args.preview:
        r = curl(
            "POST", "/rag/grade", {"question_id": args.question_id, "preview": True}
        )
        print(f"\nQuestion: {r.get('question', '')}")
        print(f"Solution: {r.get('solution', 'N/A')}")
        return
    if not args.answer:
        print("Error: provide answer as positional arg or --answer flag")
        return
    r = curl(
        "POST",
        "/rag/grade",
        {"question_id": args.question_id, "user_answer": args.answer},
    )
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"\n{'=' * 60}")
    print(f"Grade: {r.get('score', '?')}/5  {'✅' if r.get('correct') else '❌'}")
    print(f"FSRS Rating: {r.get('rating', 'N/A')}")
    print(f"{'=' * 60}")
    print(f"\n{r.get('feedback', '')}")
    print(f"\nSolution: {r.get('solution', '')}")
    print(f"\nNext review: {r.get('next_review', 'N/A')[:10]}")
    print(f"State: {r.get('state', '')}")


def cmd_track(args):
    r = curl("GET", "/rag/track")
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    s = r.get("stats", {})
    print(f"\n{'=' * 40}")
    print(f"  Progress Overview")
    print(f"{'=' * 40}")
    print(
        f"  Questions:    {s.get('attempted', 0)}/{s.get('total_questions', 0)} attempted"
    )
    print(f"  Accuracy:     {s.get('accuracy', 0)}%")
    print(f"  Due now:      {s.get('due_now', 0)}")
    print(f"  Total reviews:{s.get('total_reviews', 0)}")
    print(f"  Lapses:       {s.get('lapses', 0)}")
    print(f"{'=' * 40}")
    if r.get("due_questions"):
        print(f"\nDue questions:")
        for q in r["due_questions"]:
            print(f"  [{q['id']}] {q['question_text'][:60]}")


def cmd_status(args):
    r = curl("GET", "/rag/status")
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"Vectors in Upstash: {r.get('total_vectors', 0)}")
    print(f"Dimension: {r.get('dimension', 0)}")

    r2 = curl("GET", "/rag/track")
    s = r2.get("stats", {})
    print(f"D1 Questions: {s.get('total_questions', 0)}")
    print(f"Attempted: {s.get('attempted', 0)}")
    print(f"Accuracy: {s.get('accuracy', 0)}%")
    print(f"Due now: {s.get('due_now', 0)}")
    print(f"Lapses: {s.get('lapses', 0)}")


def cmd_list(args):
    params = []
    if args.subject:
        params.append(f"subject={args.subject}")
    if args.difficulty:
        params.append(f"difficulty={args.difficulty}")
    if args.limit:
        params.append(f"limit={args.limit}")
    qs = curl("GET", "/rag/questions?" + "&".join(params))
    if "error" in qs:
        print(f"Error: {qs['error']}")
        return
    print(f"\n{'=' * 50}")
    print(f"  Total: {qs.get('total', 0)} questions")
    print(f"{'=' * 50}")
    for q in qs.get("questions", []):
        diff = "⭐" * q.get("difficulty", 1)
        tags = (q.get("topic_tags") or "").replace(",", " · ")
        print(f"\n[{q['id']}] {diff} {q['question_text'][:65]}")
        if tags:
            print(f"     {tags}")


def cmd_loop(args):
    subject = args.subject
    limit = args.limit
    count = 0
    correct = 0
    import time

    session_start = time.time()

    try:
        while count < limit:
            params = f"?subject={subject}" if subject else ""
            r = curl("GET", f"/rag/practice{params}")
            if "error" in r:
                print(f"Error: {r['error']}")
                return
            q = r.get("question")
            if not q:
                print("No more questions available!")
                break

            count += 1
            print(f"\n{'=' * 60}")
            print(f"Question #{q['id']} [{'⭐' * q['difficulty']}] ({count}/{limit})")
            print(f"{'=' * 60}")
            print(f"\n{q['text']}\n")

            print(f"  [p] Preview solution  [s] Skip")
            ans = input("Ta réponse: ").strip()

            if ans == "p":
                pv = curl(
                    "POST", "/rag/grade", {"question_id": q["id"], "preview": True}
                )
                print(f"\n  → Solution: {pv.get('solution', 'N/A')}")
                print(f"  (not graded)")
                continue

            if ans == "s":
                print("  Skipped.")
                continue

            if not ans:
                continue

            gr = curl(
                "POST", "/rag/grade", {"question_id": q["id"], "user_answer": ans}
            )
            ok = gr.get("correct", False)
            if ok:
                correct += 1
            print(f"\n{'=' * 60}")
            print(
                f"Grade: {gr.get('score', '?')}/5  {'✅' if ok else '❌'}  [{gr.get('rating', 'N/A')}]"
            )
            print(f"Solution: {gr.get('solution', '')}")
            print(
                f"Next review: {(gr.get('next_review') or '')[:10]}  State: {gr.get('state', '')}"
            )
            print(f"\n  [r] Override rating   [Enter] Accept & continue")
            rerate = input("Choice: ").strip().lower()
            if rerate == "r":
                print(f"  Rating: Again=1  Hard=2  Good=3  Easy=4")
                rating_map = {"1": "Again", "2": "Hard", "3": "Good", "4": "Easy"}
                choice = input("  Choose (1-4): ").strip()
                if choice in rating_map:
                    rrt = curl(
                        "POST",
                        "/rag/rate",
                        {
                            "question_id": q["id"],
                            "rating": rating_map[choice],
                            "correct": ok,
                        },
                    )
                    print(f"  Updated → {rrt.get('message', '')}")

    except (EOFError, KeyboardInterrupt):
        print("\n\nSession ended.")

    elapsed = int(time.time() - session_start) if "time" in dir() else 0
    if count > 0:
        curl(
            "POST",
            "/rag/session",
            {"questions_reviewed": count, "correct": correct, "time_spent": elapsed},
        )
    print(f"\n{'=' * 40}")
    print(f"  Session: {count} questions, {correct} correct")
    print(f"  Accuracy: {round(correct / count * 100) if count else 0}%")
    if elapsed > 0:
        print(f"  Time: {elapsed // 60}m {elapsed % 60}s")
    print(f"{'=' * 40}")


def cmd_quiz(args):
    subject = args.subject
    limit = args.limit
    difficulty = args.difficulty

    params = []
    if subject:
        params.append(f"subject={subject}")
    if difficulty:
        params.append(f"difficulty={difficulty}")
    params.append(f"limit=100")
    qs = curl("GET", "/rag/questions?" + "&".join(params))
    if "error" in qs:
        print(f"Error: {qs['error']}")
        return
    questions = qs.get("questions", [])
    if len(questions) < 4:
        print(f"Not enough questions ({len(questions)}) for a quiz.")
        return

    import random

    random.shuffle(questions)

    quiz_qs = []
    for q in questions:
        if len(quiz_qs) >= limit:
            break
        correct = q.get("solution_text", "").strip()
        if not correct or len(correct) < 3:
            continue
        tag = q.get("topic_tags", "") or ""
        candidates = [
            r
            for r in questions
            if r["id"] != q["id"]
            and r.get("topic_tags", "") == tag
            and r.get("solution_text")
        ]
        if len(candidates) < 3:
            candidates = [
                r for r in questions if r["id"] != q["id"] and r.get("solution_text")
            ]
        if len(candidates) < 3:
            continue
        others = random.sample(candidates, 3)
        options = [(correct, True)] + [
            (r.get("solution_text", "").strip(), False) for r in others
        ]
        random.shuffle(options)
        quiz_qs.append(
            {
                "id": q["id"],
                "text": q["question_text"],
                "options": options,
            }
        )

    if len(quiz_qs) < limit:
        print(f"Warning: only {len(quiz_qs)} valid quiz questions.\n")

    count = 0
    correct = 0
    try:
        for qi in quiz_qs:
            count += 1
            print(f"\n{'=' * 50}")
            print(f"Q{count}/{len(quiz_qs)} [{qi['id']}]")
            print(f"{'=' * 50}")
            print(f"{qi['text']}\n")
            for i, (opt, _) in enumerate(qi["options"]):
                print(f"  {i + 1}. {opt[:80]}")
            print()

            try:
                ans = input("Choice (1-4): ").strip()
                if ans not in ("1", "2", "3", "4"):
                    print("  Skipped.")
                    continue
                chosen_is_correct = qi["options"][int(ans) - 1][1]
                if chosen_is_correct:
                    correct += 1
                print(
                    f"\n  {'✅' if chosen_is_correct else '❌'} {qi['options'][int(ans) - 1][0][:60]}"
                )
                if not chosen_is_correct:
                    correct_opt = next(o[0] for o in qi["options"] if o[1])
                    print(f"  → {correct_opt[:60]}")
            except (EOFError, KeyboardInterrupt):
                break

    except (EOFError, KeyboardInterrupt):
        print("\n\nSession ended.")

    print(f"\n{'=' * 40}")
    print(
        f"  Score: {correct}/{count} ({round(correct / count * 100) if count else 0}%)"
    )
    print(f"{'=' * 40}")


def cmd_stats(args):
    r = curl("GET", "/rag/session")
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"\n{'=' * 40}")
    print(f"  🔥 Streak: {r.get('streak', 0)} jours")
    if r.get("today"):
        t = r["today"]
        print(
            f"  Aujourd'hui: {t.get('reviews', 0)} questions, {t.get('correct', 0)} correctes"
        )
    else:
        print(f"  Aujourd'hui: pas encore practiced")
    print(f"  Total reviews: {r.get('total_reviews', 0)}")
    print(f"  Accuracy: {r.get('overall_accuracy', 0)}%")
    print(f"{'=' * 40}")
    if r.get("recent"):
        print(f"\nDerniers jours:")
        for d in r["recent"]:
            date = d.get("date", "")[5:]
            rev = d.get("reviews", 0)
            corr = d.get("correct", 0)
            bar = "█" * min(rev, 20)
            print(f"  {date}: {rev:3d} {bar} ({corr} ok)")


def cmd_due(args):
    r = curl("GET", "/rag/track")
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    due = r.get("due_questions", [])
    if not due:
        print("No due questions!")
        return
    for q in due:
        print(f"  [{q['id']}] {q['question_text'][:60]} (due: {q['due_date'][:16]})")


def cmd_retry(args):
    r = curl("GET", "/rag/track")
    due = r.get("due_questions", [])
    if not due:
        print("No due questions!")
        return
    import random

    random.shuffle(due)
    limit = args.limit if hasattr(args, "limit") else len(due)
    for qi in due[:limit]:
        pv = curl("POST", "/rag/grade", {"question_id": qi["id"], "preview": True})
        print(f"\n[{qi['id']}] {qi['question_text'][:65]}")
        print(f"  Sol: {pv.get('solution', 'N/A')[:80]}")
        ans = input("  Ta réponse: ").strip()
        if not ans:
            continue
        gr = curl("POST", "/rag/grade", {"question_id": qi["id"], "user_answer": ans})
        print(
            f"  Grade: {gr.get('score', '?')}/5 {'✅' if gr.get('correct') else '❌'} [{gr.get('rating', 'N/A')}]"
        )
        print(f"  Next: {(gr.get('next_review') or '')[:10]} {gr.get('state', '')}")


def cmd_add_question(args):
    data = {
        "question": args.question,
        "subject": args.subject or "general",
        "difficulty": args.difficulty or 2,
    }
    if args.solution:
        data["solution"] = args.solution
    elif args.solve:
        data["solve"] = True
    print("Adding question...")
    r = curl("POST", "/rag/add-question", data, timeout=120)
    if "error" in r:
        print(f"Error: {r['error']}")
        return
    print(f"Added question #{r.get('id')}: {r.get('question', '')[:60]}")
    if r.get("solution"):
        print(f"Solution: {r.get('solution', '')[:200]}...")


def main():
    parser = argparse.ArgumentParser(description="BAC Study RAG CLI")
    sub = parser.add_subparsers(dest="cmd")

    p = sub.add_parser("query", help="Query the knowledge vault")
    p.add_argument("query", help="Your question")
    p.add_argument("-s", "--subjects", help="Subjects (e.g., biology,chemistry)")
    p.add_argument("-m", "--model", help="Model override")
    p.set_defaults(func=cmd_query)

    p = sub.add_parser("search", help="Search for chunks")
    p.add_argument("query", help="Search query")
    p.add_argument("-s", "--subjects", help="Subjects")
    p.add_argument("-k", "--top-k", type=int, default=5, help="Number of results")
    p.set_defaults(func=cmd_search)

    p = sub.add_parser("solve", help="Solve a question")
    p.add_argument("question", help="Question to solve")
    p.add_argument("-s", "--subject", help="Subject")
    p.set_defaults(func=cmd_solve)

    p = sub.add_parser("practice", help="Get practice question")
    p.add_argument("-s", "--subject", help="Subject filter")
    p.add_argument("-a", "--answer", help="Answer immediately")
    p.set_defaults(func=cmd_practice)

    p = sub.add_parser("loop", help="Interactive practice loop")
    p.add_argument("-s", "--subject", help="Subject filter")
    p.add_argument("-n", "--limit", type=int, default=10, help="Number of questions")
    p.set_defaults(func=cmd_loop)

    p = sub.add_parser("quiz", help="Multiple choice quiz")
    p.add_argument("-s", "--subject", help="Subject filter")
    p.add_argument("-d", "--difficulty", type=int, help="Difficulty filter")
    p.add_argument("-n", "--limit", type=int, default=5, help="Number of questions")
    p.set_defaults(func=cmd_quiz)

    p = sub.add_parser("grade", help="Grade an answer or rate directly")
    p.add_argument("question_id", type=int, help="Question ID")
    p.add_argument("answer", nargs="?", help="Your answer")
    p.add_argument("--preview", action="store_true", help="Show solution only")
    p.add_argument(
        "--rate",
        choices=["Again", "Hard", "Good", "Easy"],
        help="Self-rate directly (skip grading)",
    )
    p.set_defaults(func=cmd_grade)

    p = sub.add_parser("rate", help="Self-rate a question (FSRS)")
    p.add_argument("question_id", type=int, help="Question ID")
    p.add_argument(
        "rating", choices=["Again", "Hard", "Good", "Easy"], help="Your self-assessment"
    )
    p.set_defaults(func=cmd_grade)

    p = sub.add_parser("track", help="Track progress")
    p.set_defaults(func=cmd_track)

    p = sub.add_parser("list", help="List questions")
    p.add_argument("-s", "--subject", help="Filter by subject")
    p.add_argument("-d", "--difficulty", type=int, help="Filter by difficulty")
    p.add_argument("-n", "--limit", type=int, help="Max results")
    p.set_defaults(func=cmd_list)

    p = sub.add_parser("due", help="List due questions")
    p.set_defaults(func=cmd_due)

    p = sub.add_parser("retry", help="Practice due questions interactively")
    p.add_argument("-n", "--limit", type=int, default=999, help="Max questions")
    p.set_defaults(func=cmd_retry)

    p = sub.add_parser("stats", help="Study streak and stats")
    p.set_defaults(func=cmd_stats)

    p = sub.add_parser("status", help="Check system status")
    p.set_defaults(func=cmd_status)

    p = sub.add_parser("add", help="Add a question")
    p.add_argument("question", help="Question text")
    p.add_argument("-s", "--subject", help="Subject")
    p.add_argument("-d", "--difficulty", type=int, help="Difficulty 1-5")
    p.add_argument("--solution", help="Solution text")
    p.add_argument("--solve", action="store_true", help="Auto-solve with AI")
    p.set_defaults(func=cmd_add_question)

    args = parser.parse_args(sys.argv[1:] if len(sys.argv) > 1 else ["--help"])
    if args.cmd is None:
        parser.print_help()
        return
    args.func(args)


if __name__ == "__main__":
    main()
