#!/usr/bin/env python3
import sqlite3
import sys
from pathlib import Path

VAULT = Path(__file__).parent
DB = VAULT / ".vault-index.db"

def search(query, limit=10):
    conn = sqlite3.connect(DB)
    c = conn.cursor()
    c.execute('''SELECT path, title FROM fts_notes 
                 WHERE fts_notes MATCH ? ORDER BY rank LIMIT ?''', (query, limit))
    results = c.fetchall()
    conn.close()
    
    print(f"\n🔍 Search results for '{query}':\n")
    for i, (path, title) in enumerate(results, 1):
        print(f"{i}. {title}")
        print(f"   📁 {path}\n")

def stats():
    conn = sqlite3.connect(DB)
    c = conn.cursor()
    
    c.execute("SELECT COUNT(*) FROM notes")
    notes = c.fetchone()[0]
    c.execute("SELECT COUNT(*) FROM links")
    links = c.fetchone()[0]
    c.execute("SELECT COUNT(*) FROM tags")
    tags = c.fetchone()[0]
    c.execute("SELECT COUNT(*) FROM formulas")
    formulas = c.fetchone()[0]
    
    c.execute("SELECT subject, COUNT(*) FROM notes WHERE subject IS NOT NULL GROUP BY subject")
    by_subject = c.fetchall()
    
    c.execute("SELECT tag, count FROM tags ORDER BY count DESC LIMIT 10")
    top_tags = c.fetchall()
    
    conn.close()
    
    print("\n📊 Vault Statistics:\n")
    print(f"  Total Notes: {notes}")
    print(f"  Total Links: {links}")
    print(f"  Total Tags: {tags}")
    print(f"  Total Formulas: {formulas}")
    
    print("\n📚 Notes by Subject:")
    for subj, count in by_subject:
        print(f"  {subj}: {count}")
    
    print("\n🏷️  Top Tags:")
    for tag, count in top_tags:
        print(f"  #{tag}: {count}")

def orphans():
    conn = sqlite3.connect(DB)
    c = conn.cursor()
    c.execute('''SELECT path, title FROM notes 
                 WHERE path NOT IN (SELECT DISTINCT source_path FROM links)
                 AND path NOT IN (SELECT DISTINCT target_path FROM links)
                 AND type != 'moc' ''')
    results = c.fetchall()
    conn.close()
    
    print(f"\n🔗 Orphaned Notes ({len(results)}):\n")
    for path, title in results:
        print(f"  • {title} ({path})")

def formulas_search(query):
    conn = sqlite3.connect(DB)
    c = conn.cursor()
    c.execute('''SELECT DISTINCT path, formula FROM formulas 
                 WHERE formula LIKE ? LIMIT 20''', (f'%{query}%',))
    results = c.fetchall()
    conn.close()
    
    print(f"\n📐 Formulas containing '{query}':\n")
    for path, formula in results:
        print(f"  {formula}")
        print(f"  📁 {path}\n")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: ./vault-cli.py <command> [args]")
        print("\nCommands:")
        print("  search <query>    - Search notes")
        print("  stats             - Show vault statistics")
        print("  orphans           - Find unlinked notes")
        print("  formulas <query>  - Search formulas")
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "search" and len(sys.argv) > 2:
        search(' '.join(sys.argv[2:]))
    elif cmd == "stats":
        stats()
    elif cmd == "orphans":
        orphans()
    elif cmd == "formulas" and len(sys.argv) > 2:
        formulas_search(' '.join(sys.argv[2:]))
    else:
        print(f"Unknown command: {cmd}")
