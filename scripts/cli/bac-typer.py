#!/usr/bin/env python3
"""
BAC Study RAG CLI — Built with Typer.
Usage: python3 bac-typer.py [command] [args]

Flags:
    --json, -j    Output raw JSON instead of formatted text
    --quiet, -q   Minimal output (errors only)
"""

import sys
import json
import subprocess
from typing import Optional
import typer
from rich.console import Console
from rich.table import Table

app = typer.Typer(add_completion=False)
console = Console()

WORKER = "https://bac-api.amdajhbdh.workers.dev"


def curl(
    method: str, path: str, data: Optional[dict] = None, timeout: int = 60
) -> dict:
    """Call worker API via curl."""
    cmd = ["curl", "-s", "-X", method, f"{WORKER}{path}"]
    if data:
        cmd.extend(["-H", "Content-Type: application/json", "-d", json.dumps(data)])
    try:
        result = subprocess.run(
            cmd, capture_output=True, text=True, timeout=timeout + 10
        )
        return json.loads(result.stdout)
    except (json.JSONDecodeError, subprocess.TimeoutExpired) as e:
        return {"error": str(e)}


def output_json(data: dict, pretty: bool = False):
    """Output JSON."""
    if pretty:
        print(json.dumps(data, indent=2, ensure_ascii=False))
    else:
        print(json.dumps(data, ensure_ascii=False))


@app.command()
def query(
    question: str = typer.Argument(..., help="Question to ask"),
    subjects: Optional[str] = typer.Option(
        None, "--subjects", "-s", help="Filter by subjects"
    ),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Query the knowledge vault."""
    data = {"query": question}
    if subjects:
        data["subjects"] = subjects.split(",")

    r = curl("POST", "/rag/query", data)

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    console.print(f"\n[bold cyan]Question:[/bold cyan] {question}")
    console.print(f"[dim]Subjects:[/dim] {', '.join(r.get('subjects', []))}")
    console.print(f"\n{r.get('answer', 'No answer')}\n")

    sources = r.get("sources", [])
    if sources:
        console.print("[bold]Sources:[/bold]")
        for i, s in enumerate(sources, 1):
            console.print(f"  [{i}] {s['source']} (score: {s['score']:.2f})")


@app.command()
def solve(
    question: str = typer.Argument(..., help="Question to solve"),
    subject: Optional[str] = typer.Option(None, "--subject", "-s", help="Subject"),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Solve a math/physics/chemistry question."""
    data = {"question": question}
    if subject:
        data["subject"] = subject

    console.print("[yellow]Solving...[/yellow]")
    r = curl("POST", "/rag/solve", data, timeout=120)

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    console.print(
        f"\n[bold green]Solution:[/bold green]\n{r.get('solution', 'No solution')}"
    )


@app.command()
def search(
    query: str = typer.Argument(..., help="Search query"),
    subjects: Optional[str] = typer.Option(
        None, "--subjects", "-s", help="Filter by subjects"
    ),
    top_k: int = typer.Option(5, "--top", "-k", help="Number of results"),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Search the vault for relevant chunks."""
    data = {"query": query, "topK": top_k}
    if subjects:
        data["subjects"] = subjects.split(",")

    r = curl("POST", "/rag/search", data)

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    chunks = r.get("chunks", [])
    if not chunks:
        console.print("[yellow]No results found.[/yellow]")
        return

    console.print(f"\n[bold]Search Results for:[/bold] {query}\n")
    for i, c in enumerate(chunks, 1):
        text = c["metadata"]["text"][:200]
        console.print(
            f"[{i}] [cyan]{c['metadata']['subject']}[/cyan] | score: {c['score']:.3f}"
        )
        console.print(f"    {text}...\n")


@app.command()
def practice(
    subject: Optional[str] = typer.Option(
        None, "--subject", "-s", help="Subject filter"
    ),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Get practice questions."""
    url = f"/rag/practice?subject={subject}" if subject else "/rag/practice"
    r = curl("GET", url)

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    q = r.get("question")
    if not q:
        console.print("[yellow]No questions available.[/yellow]")
        return

    console.print(f"\n[bold]Question #{q['id']}[/bold]")
    console.print(f"[cyan]{q['text']}[/cyan]")
    console.print(f"\n[dim]Topic:[/dim] {q.get('topic_tags', 'N/A')}")


@app.command()
def grade(
    question_id: int = typer.Argument(..., help="Question ID"),
    answer: str = typer.Argument(..., help="Your answer"),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Grade your answer to a practice question."""
    data = {"question_id": question_id, "user_answer": answer}
    r = curl("POST", "/rag/grade", data)

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    result = r.get("result", {})
    console.print(f"\n[bold]Score:[/bold] {result.get('score', 0)}/5")
    console.print(f"[bold]Feedback:[/bold] {result.get('feedback', 'N/A')}")


@app.command()
def track(
    user_id: int = typer.Option(1, "--user", "-u", help="User ID"),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Show study progress and streaks."""
    r = curl("GET", f"/rag/track?user_id={user_id}")

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    stats = r.get("stats", {})

    table = Table(title="Study Progress")
    table.add_column("Metric", style="cyan")
    table.add_column("Value", style="green")

    table.add_row("Streak", str(stats.get("streak", 0)))
    table.add_row("Total Reviews", str(stats.get("total_reviews", 0)))
    table.add_row("Accuracy", f"{stats.get('accuracy', 0):.1f}%")
    table.add_row("Due Today", str(stats.get("due_today", 0)))

    console.print(table)


@app.command()
def status(json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON")):
    """Check system status."""
    r = curl("GET", "/rag/status")
    s = curl("GET", "/rag/subjects")

    if json_output:
        output_json({"status": r, "subjects": s})
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    table = Table(title="System Status")
    table.add_column("Metric", style="cyan")
    table.add_column("Value", style="green")

    table.add_row("Total Vectors", str(r.get("total_vectors", 0)))
    table.add_row("Dimension", str(r.get("dimension", 0)))

    for subj, count in s.items():
        table.add_row(subj.replace("subject_", ""), str(count))

    console.print(table)


@app.command()
def rate(
    question_id: int = typer.Argument(..., help="Question ID"),
    rating: str = typer.Argument(..., help="Rating: again, hard, good, easy"),
    json_output: bool = typer.Option(False, "--json", "-j", help="Output JSON"),
):
    """Rate a question for spaced repetition."""
    rating_map = {"again": 0, "hard": 3, "good": 4, "easy": 5}
    rating_val = rating_map.get(rating.lower())

    if rating_val is None:
        console.print(f"[red]Invalid rating. Use: again, hard, good, easy[/red]")
        return

    data = {"question_id": question_id, "rating": rating_val}
    r = curl("POST", "/rag/rate", data)

    if json_output:
        output_json(r)
        return

    if "error" in r:
        console.print(f"[red]Error:[/red] {r.get('error')}")
        return

    console.print(f"[green]Rated: {rating}[/green]")


if __name__ == "__main__":
    app()
