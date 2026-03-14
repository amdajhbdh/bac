#!/usr/bin/env python3
"""AI Content Enhancement System - Task 12"""
import sqlite3
import json
import re
from pathlib import Path
from datetime import datetime

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
DB = VAULT / ".vault-index.db"

class ContentEnhancer:
    def __init__(self):
        self.conn = sqlite3.connect(DB)
        self.conn.row_factory = sqlite3.Row
        
    def enhance_note(self, note_path):
        """Add tables, connections, and metadata to a note"""
        path = VAULT / note_path
        if not path.exists():
            return {"error": "Note not found"}
            
        content = path.read_text()
        metadata = self._extract_metadata(content)
        related = self._find_related(note_path, content)
        enhanced = self._add_enhancements(content, metadata, related)
        path.write_text(enhanced)
        
        return {
            "note": note_path,
            "related_count": len(related),
            "metadata": metadata,
            "enhanced": True
        }
    
    def _extract_metadata(self, content):
        """Extract formulas, concepts, difficulty"""
        formulas = re.findall(r'\$\$([^$]+)\$\$', content)
        concepts = re.findall(r'#(\w+)', content)
        
        difficulty = "medium"
        if len(formulas) > 5 or "advanced" in content.lower():
            difficulty = "hard"
        elif len(formulas) < 2:
            difficulty = "easy"
            
        return {
            "formulas": len(formulas),
            "concepts": list(set(concepts)),
            "difficulty": difficulty,
            "word_count": len(content.split())
        }
    
    def _find_related(self, note_path, content):
        """Find related notes using FTS5 search"""
        terms = re.findall(r'\b[A-Z][a-z]+\b', content)[:10]
        related = []
        
        for term in terms:
            try:
                cursor = self.conn.execute("""
                    SELECT path, snippet(fts_notes, 1, '<b>', '</b>', '...', 32) as snippet
                    FROM fts_notes 
                    WHERE content MATCH ? AND path != ?
                    LIMIT 3
                """, (term, note_path))
                
                for row in cursor:
                    related.append({
                        "path": row["path"],
                        "snippet": row["snippet"],
                        "term": term
                    })
            except:
                pass
        
        return related[:10]
    
    def _add_enhancements(self, content, metadata, related):
        """Add all enhancements with callout styling"""
        if "> [!summary]" in content:
            return content
            
        # Add frontmatter
        if not content.startswith("---"):
            frontmatter = "---\n"
            frontmatter += f"difficulty: {metadata['difficulty']}\n"
            frontmatter += f"concepts: {json.dumps(metadata['concepts'])}\n"
            frontmatter += f"enhanced: {datetime.now().isoformat()}\n"
            frontmatter += f"project: BAC Preparation 2026\n"
            frontmatter += f"status: not-started\n"
            frontmatter += f"progress: 0\n"
            frontmatter += "---\n\n"
            content = frontmatter + content
        
        # Add summary callout
        summary = f"\n> [!summary] 📊 Note Summary\n"
        summary += f"> \n"
        summary += f"> | Property | Value |\n"
        summary += f"> |----------|-------|\n"
        summary += f"> | **Difficulty** | `{metadata['difficulty']}` #difficulty/{metadata['difficulty']} |\n"
        summary += f"> | **Formulas** | {metadata['formulas']} |\n"
        summary += f"> | **Concepts** | {len(metadata['concepts'])} |\n"
        summary += f"> | **Related Notes** | {len(related)} |\n"
        summary += f"> | **Word Count** | {metadata['word_count']} |\n"
        summary += f"> | **Last Enhanced** | {datetime.now().strftime('%Y-%m-%d')} |\n"
        
        # Insert after frontmatter
        parts = content.split("---", 2)
        if len(parts) >= 3:
            content = "---" + parts[1] + "---" + summary + "\n" + parts[2]
        else:
            content = summary + "\n" + content
        
        # Add related notes callout
        if related:
            related_section = f"\n> [!related] 🔗 Related Notes\n"
            related_section += f"> \n"
            for r in related:
                name = Path(r["path"]).stem
                related_section += f"> - [[{r['path']}|{name}]]\n"
            content += "\n" + related_section
        
        return content
    
    def enhance_all(self, folder="01-Concepts"):
        """Enhance all notes in a folder"""
        results = []
        folder_path = VAULT / folder
        
        for note in folder_path.rglob("*.md"):
            rel_path = note.relative_to(VAULT)
            result = self.enhance_note(str(rel_path))
            results.append(result)
            print(f"✓ {rel_path}")
        
        return results
    
    def generate_progress_table(self):
        """Generate progress tracking table with callouts"""
        cursor = self.conn.execute("""
            SELECT 
                SUBSTR(path, 1, INSTR(path, '/') - 1) as folder,
                COUNT(*) as total
            FROM notes
            GROUP BY folder
        """)
        
        table = "# 📈 Vault Progress\n\n"
        table += "> [!progress] Progress Overview\n"
        table += "> \n"
        table += "> | Folder | Notes | Status |\n"
        table += "> |--------|-------|--------|\n"
        
        for row in cursor:
            folder = row["folder"]
            total = row["total"]
            table += f"> | {folder} | {total} | ✅ |\n"
        
        return table

def main():
    import sys
    
    enhancer = ContentEnhancer()
    
    if len(sys.argv) < 2:
        print("Usage: content-enhancer.py <command> [args]")
        print("Commands:")
        print("  enhance <note-path>  - Enhance single note")
        print("  enhance-all [folder] - Enhance all notes in folder")
        print("  progress             - Generate progress table")
        return
    
    cmd = sys.argv[1]
    
    if cmd == "enhance" and len(sys.argv) > 2:
        result = enhancer.enhance_note(sys.argv[2])
        print(json.dumps(result, indent=2))
    
    elif cmd == "enhance-all":
        folder = sys.argv[2] if len(sys.argv) > 2 else "01-Concepts"
        results = enhancer.enhance_all(folder)
        print(f"\n✅ Enhanced {len(results)} notes")
    
    elif cmd == "progress":
        table = enhancer.generate_progress_table()
        output = VAULT / "00-Meta" / "PROGRESS-TABLE.md"
        output.write_text(table)
        print(f"✓ Progress table: {output}")

if __name__ == "__main__":
    main()
