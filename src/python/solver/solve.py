#!/usr/bin/env python3
"""
BAC UNIFIED SOLVER - Unified interface
Tries: Local Ollama → Fallback solvers
"""

import sys
import json
import subprocess


def run_ollama(prompt: str, model: str = "llama3.2:3b") -> dict:
    """Run local Ollama"""
    import requests

    payload = {
        "model": model,
        "prompt": prompt,
        "stream": False,
        "options": {"temperature": 0.3, "num_predict": 2048},
    }

    try:
        response = requests.post(
            "http://localhost:11434/api/generate", json=payload, timeout=120
        )
        if response.status_code == 200:
            return {
                "success": True,
                "response": response.json().get("response", ""),
                "solver": model,
            }
        else:
            return {"success": False, "error": f"Ollama error: {response.status_code}"}
    except Exception as e:
        return {"success": False, "error": str(e)}


def solve(problem: str, subject: str = "math") -> dict:
    """Main solve function with fallback chain"""

    # Build prompt
    prompt = f"""Tu es un professeur de mathématiques et sciences pour le BAC C Mauritanie.
Résous ce problème {subject} étape par étape en français.

Problème: {problem}

Requirements:
1. Montre TOUS les calculs
2. Explique chaque étape en français  
3. Donne la réponse finale clairement
4. Si c'est un problème de géométrie, décris la construction
5. Si c'est un problème de physique, identifie les lois utilisées

Format:
**Solution:**
[ta solution]

**Étapes:**
1. [Étape 1]
2. [Étape 2]
3. [Étape 3]
"""

    # Try local model first
    result = run_ollama(prompt, "llama3.2:3b")
    if result.get("success"):
        return result

    # Try other local models
    for model in ["llama3.2:1b", "qwen2.5:1.5b"]:
        result = run_ollama(prompt, model)
        if result.get("success"):
            return result

    return {"success": False, "error": "No local model available"}


def main():
    if len(sys.argv) < 2:
        print("Usage: python solver.py <problem>")
        sys.exit(1)

    problem = " ".join(sys.argv[1:])
    subject = "math"

    # Detect subject from problem
    problem_lower = problem.lower()
    if any(
        w in problem_lower
        for w in ["force", "courant", "tension", "physique", "mécanique"]
    ):
        subject = "pc"
    elif any(w in problem_lower for w in ["chimie", "réaction", "molécules", "atome"]):
        subject = "chimie"
    elif any(w in problem_lower for w in ["cellule", "adn", "svt", "biologie"]):
        subject = "svt"

    print(f"🔍 Subject detected: {subject}")
    print(f"🎯 Solving...")

    result = solve(problem, subject)

    if result.get("success"):
        print("\n" + "=" * 50)
        print(result.get("response", ""))
        print("=" * 50)
        print(f"✨ Solver: {result.get('solver', 'unknown')}")
    else:
        print(f"❌ Error: {result.get('error', 'Unknown')}")
        sys.exit(1)


if __name__ == "__main__":
    main()
