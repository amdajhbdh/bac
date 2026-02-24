#!/usr/bin/env python3
"""
BAC UNIFIED SUPER SOLVER
Uses ministral-3:14b-cloud with vision capabilities
No API keys needed - runs locally with Ollama
"""

import sys
import json
import base64
import requests
from pathlib import Path

OLLAMA_ENDPOINT = "http://localhost:11434"
MODEL = "ministral-3:14b-cloud"


def check_ollama():
    """Check if Ollama is running"""
    try:
        response = requests.get(f"{OLLAMA_ENDPOINT}/api/tags", timeout=5)
        if response.status_code == 200:
            models = response.json().get("models", [])
            model_names = [m.get("name", "") for m in models]
            print(f"✓ Ollama connected. Available models: {', '.join(model_names)}")
            return MODEL in model_names
    except Exception as e:
        print(f"✗ Ollama not available: {e}")
        return False
    return False


def solve_text(problem: str, subject: str = "math") -> dict:
    """Solve a text problem using local Ollama"""

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

    payload = {
        "model": MODEL,
        "prompt": prompt,
        "stream": False,
        "options": {"temperature": 0.3, "num_predict": 4096},
    }

    try:
        response = requests.post(
            f"{OLLAMA_ENDPOINT}/api/generate", json=payload, timeout=120
        )

        if response.status_code == 200:
            result = response.json()
            return {
                "success": True,
                "solution": result.get("response", ""),
                "solver": "ministral-3:14b-cloud",
                "type": "text",
            }
        else:
            return {
                "success": False,
                "error": f"Ollama error: {response.status_code}",
                "solver": MODEL,
            }
    except Exception as e:
        return {"success": False, "error": str(e), "solver": MODEL}


def solve_image(image_path: str, question: str = "") -> dict:
    """Solve a problem from image using vision capabilities"""

    # Read and encode image
    try:
        with open(image_path, "rb") as f:
            image_data = base64.b64encode(f.read()).decode("utf-8")
    except Exception as e:
        return {"success": False, "error": f"Cannot read image: {e}"}

    # Check if model supports vision
    # For now, use text prompt with image description
    prompt = f"""Tu es un professeur de mathématiques pour le BAC C.
Regarde cette image et résous le problème.

Question: {question or "Résous ce problème"}

Donne:
1. Le problème identifié
2. La solution étape par étape
3. La réponse finale
"""

    # Note: ministral-3 has vision, but we need to use correct API
    # Using chat API for vision models
    payload = {
        "model": MODEL,
        "messages": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": prompt},
                    {
                        "type": "image_url",
                        "image_url": {"url": f"data:image/jpeg;base64,{image_data}"},
                    },
                ],
            }
        ],
        "stream": False,
    }

    try:
        response = requests.post(
            f"{OLLAMA_ENDPOINT}/api/chat", json=payload, timeout=120
        )

        if response.status_code == 200:
            result = response.json()
            return {
                "success": True,
                "solution": result.get("message", {}).get("content", ""),
                "solver": MODEL,
                "type": "vision",
            }
        else:
            return {
                "success": False,
                "error": f"Vision API error: {response.status_code}",
                "solver": MODEL,
            }
    except Exception as e:
        return {"success": False, "error": str(e), "solver": MODEL}


def solve_pdf(pdf_path: str) -> dict:
    """Extract and solve from PDF - simplified"""
    # For PDFs, we'd use pdftotext first
    import subprocess

    try:
        result = subprocess.run(
            ["pdftotext", "-layout", pdf_path, "-"],
            capture_output=True,
            text=True,
            timeout=30,
        )

        if result.returncode == 0:
            text = result.stdout
            # Take first 2000 chars for the problem
            problem_preview = text[:2000]
            return solve_text(problem_preview, "math")
        else:
            return {"success": False, "error": "Cannot extract PDF text"}
    except FileNotFoundError:
        return {"success": False, "error": "pdftotext not installed"}
    except Exception as e:
        return {"success": False, "error": str(e)}


def solve(problem: str, input_type: str = "text", subject: str = "math") -> dict:
    """Main solve function"""

    if input_type == "image":
        return solve_image(problem)
    elif input_type == "pdf":
        return solve_pdf(problem)
    else:
        return solve_text(problem, subject)


def main():
    """CLI entry point"""
    import argparse

    parser = argparse.ArgumentParser(description="BAC Unified Super Solver")
    parser.add_argument("problem", help="Problem text or file path")
    parser.add_argument(
        "--type", choices=["text", "image", "pdf"], default="text", help="Input type"
    )
    parser.add_argument(
        "--subject", default="math", help="Subject: math, pc, svt, chimie"
    )
    parser.add_argument("--check", action="store_true", help="Check Ollama status only")

    args = parser.parse_args()

    if args.check:
        if check_ollama():
            print(f"✓ Model {MODEL} is available!")
            sys.exit(0)
        else:
            print(f"✗ Model {MODEL} not found")
            sys.exit(1)

    print(f"🎯 Solving with {MODEL}...")

    result = solve(args.problem, args.type, args.subject)

    if result.get("success"):
        print("\n" + "=" * 50)
        print("SOLUTION:")
        print("=" * 50)
        print(result.get("solution", ""))
        print("=" * 50)
        print(f"Solver: {result.get('solver')}")
    else:
        print(f"✗ Error: {result.get('error', 'Unknown error')}")
        sys.exit(1)


if __name__ == "__main__":
    main()
