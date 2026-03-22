#!/usr/bin/env python3
"""
Simple RAG Evaluation Script for BAC Study System.
Uses the Worker API directly - no heavy dependencies needed.
"""

import json
import subprocess
import sys
import time

WORKER = "https://bac-api.amdajhbdh.workers.dev"

SAMPLE_QUESTIONS = [
    # Biology
    ("explique la mitose", "biology"),
    ("qu'est-ce que l'ADN", "biology"),
    ("comment fonctionne la photosynthèse", "biology"),
    ("décris le cycle de Krebs", "biology"),
    ("c'est quoi une enzyme", "biology"),
    # Chemistry
    ("définis un atome", "chemistry"),
    ("qu'est-ce qu'une molécule", "chemistry"),
    ("explique la réaction de combustion", "chemistry"),
    ("c'est quoi un ácido-base", "chemistry"),
    ("définis le pH", "chemistry"),
    # Math
    ("résous x² - 5x + 6 = 0", "math"),
    ("calcule la dérivée de x³", "math"),
    ("qu'est-ce qu'une fonction", "math"),
    ("définis une intégrale", "math"),
    ("résous 2x + 3 = 7", "math"),
    # Physics
    ("définis la force", "physics"),
    ("c'est quoi l'énergie cinétique", "physics"),
    ("explique la gravité", "physics"),
    ("qu'est-ce que la vitesse", "physics"),
    ("définis le travail en physique", "physics"),
]


def call_api(endpoint, payload):
    """Call worker API via curl."""
    result = subprocess.run(
        [
            "curl",
            "-s",
            "-X",
            "POST",
            f"{WORKER}/{endpoint}",
            "-H",
            "Content-Type: application/json",
            "-d",
            json.dumps(payload),
        ],
        capture_output=True,
        text=True,
        timeout=30,
    )
    return json.loads(result.stdout)


def evaluate_question(question, expected_subject):
    """Evaluate a single question."""
    # Call RAG query
    response = call_api("rag/query", {"query": question})

    answer = response.get("answer", "")
    sources = response.get("sources", [])
    model = response.get("model", "unknown")
    detected_subjects = response.get("subjects", [])

    # Check if answer is non-empty
    has_answer = len(answer) > 20

    # Check if sources include the expected subject
    source_subjects = [s.get("subject", "") for s in sources]
    correct_subject = (
        expected_subject in source_subjects or expected_subject in detected_subjects
    )

    # Check if answer is relevant (contains keywords from question)
    relevant = (
        any(word.lower() in answer.lower() for word in question.split()[:3])
        if has_answer
        else False
    )

    return {
        "question": question[:30] + "...",
        "expected": expected_subject,
        "detected": detected_subjects,
        "has_answer": has_answer,
        "correct_subject": correct_subject,
        "relevant": relevant,
        "model": model,
        "num_sources": len(sources),
    }


def main():
    print("=" * 60)
    print("  BAC RAG EVALUATION")
    print("=" * 60)

    sample = SAMPLE_QUESTIONS[:20]  # Run 20 questions

    results = []
    subjects = {"biology": [], "chemistry": [], "math": [], "physics": []}

    for i, (question, expected) in enumerate(sample):
        print(f"[{i + 1}/20] {question[:40]}...", end=" ")

        try:
            result = evaluate_question(question, expected)
            results.append(result)
            subjects[expected].append(result)

            status = "✅" if result["has_answer"] and result["relevant"] else "❌"
            print(status)

        except Exception as e:
            print(f"❌ Error: {e}")
            results.append(
                {"question": question, "expected": expected, "error": str(e)}
            )

        time.sleep(0.5)  # Rate limit

    # Calculate metrics
    total = len(results)
    has_answer = sum(1 for r in results if r.get("has_answer", False))
    relevant = sum(1 for r in results if r.get("relevant", False))
    correct_subject = sum(1 for r in results if r.get("correct_subject", False))

    print("\n" + "=" * 60)
    print("  RESULTS")
    print("=" * 60)

    print(f"\nOverall ({total} questions):")
    print(f"  - Has Answer:   {has_answer}/{total} ({100 * has_answer // total}%)")
    print(f"  - Relevant:     {relevant}/{total} ({100 * relevant // total}%)")
    print(
        f"  - Correct Subj: {correct_subject}/{total} ({100 * correct_subject // total}%)"
    )

    print("\nPer Subject:")
    for subj, subj_results in subjects.items():
        count = len(subj_results)
        if count > 0:
            has_ans = sum(1 for r in subj_results if r.get("has_answer", False))
            rel = sum(1 for r in subj_results if r.get("relevant", False))
            print(f"  {subj:12}: {has_ans}/{count} answers, {rel}/{count} relevant")

    # Model usage
    models = {}
    for r in results:
        m = r.get("model", "unknown")
        models[m] = models.get(m, 0) + 1

    print("\nModel Usage:")
    for m, c in models.items():
        print(f"  - {m}: {c}")

    print("\n" + "=" * 60)

    # Simple score (0-100)
    score = has_answer * 0.4 + relevant * 0.4 + correct_subject * 0.2
    print(f"  Score: {score:.0f}/100")
    print("=" * 60)

    return 0 if score > 60 else 1


if __name__ == "__main__":
    sys.exit(main())
