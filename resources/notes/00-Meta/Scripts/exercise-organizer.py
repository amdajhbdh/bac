#!/usr/bin/env python3
"""
Exercise Organization Agent - Task C
Categorize exercises, link to concepts, add to SR
"""
import sqlite3
from pathlib import Path
import re

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
EXERCISES = VAULT / "02-Practice/FromTextbooks"
DB = VAULT / ".vault-index.db"
SR_DB = VAULT / ".spaced-repetition.db"

class ExerciseOrganizer:
    def __init__(self):
        self.conn = sqlite3.connect(DB)
        self.sr_conn = sqlite3.connect(SR_DB)
    
    def assess_difficulty(self, content: str) -> str:
        """Assess exercise difficulty"""
        score = 0
        
        if len(content) > 1000: score += 1
        if content.count('$$') > 3: score += 1
        if any(w in content.lower() for w in ['prove', 'demonstrate', 'derive']): score += 2
        if re.search(r'\d+\)', content):
            score += min(len(re.findall(r'\d+\)', content)) // 3, 2)
        
        return 'easy' if score <= 2 else 'medium' if score <= 4 else 'hard'
    
    def find_related_concepts(self, content: str) -> list:
        """Find related concept notes"""
        terms = re.findall(r'\b[a-z]{4,}\b', content.lower())
        
        c = self.conn.cursor()
        related = []
        
        for term in set(terms[:20]):
            c.execute("SELECT title, path FROM notes WHERE type='concept' AND title LIKE ? LIMIT 2",
                     (f'%{term}%',))
            related.extend(c.fetchall())
        
        return list(set(related))[:5]
    
    def add_to_sr(self, path: str):
        """Add to spaced repetition"""
        c = self.sr_conn.cursor()
        from datetime import datetime
        c.execute('''INSERT OR IGNORE INTO cards (note_path, due_date, last_review)
                     VALUES (?, ?, ?)''',
                  (path, datetime.now().isoformat(), datetime.now().isoformat()))
        self.sr_conn.commit()
    
    def organize_exercise(self, ex_path: Path):
        """Organize single exercise"""
        content = ex_path.read_text()
        
        difficulty = self.assess_difficulty(content)
        related = self.find_related_concepts(content)
        
        # Update content
        content = re.sub(r'\*\*Difficulty:\*\* \[.*?\]',
                        f'**Difficulty:** {difficulty.title()}', content)
        
        if related and '## Related Concepts' not in content:
            content += "\n## Related Concepts\n\n"
            for title, _ in related:
                content += f"- [[{title}]]\n"
        
        ex_path.write_text(content)
        
        # Add to SR
        self.add_to_sr(str(ex_path.relative_to(VAULT)))
        
        return difficulty
    
    def organize_all(self):
        """Organize all exercises"""
        exercises = list(EXERCISES.rglob("*.md"))[:100]
        
        print(f"\n📚 Organizing {len(exercises)} exercises...\n")
        
        stats = {'easy': 0, 'medium': 0, 'hard': 0}
        
        for i, ex in enumerate(exercises, 1):
            if i % 20 == 0:
                print(f"  Progress: {i}/{len(exercises)}")
            
            try:
                diff = self.organize_exercise(ex)
                stats[diff] += 1
            except:
                pass
        
        print(f"\n{'='*60}")
        print(f"✅ Complete!")
        print(f"{'='*60}")
        print(f"Easy: {stats['easy']}")
        print(f"Medium: {stats['medium']}")
        print(f"Hard: {stats['hard']}")
        print(f"Added to SR: {sum(stats.values())}")

if __name__ == "__main__":
    organizer = ExerciseOrganizer()
    organizer.organize_all()
