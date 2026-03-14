#!/usr/bin/env python3
"""
Spaced Repetition System - SM-2 Algorithm
Agent-driven review scheduling
"""
import sqlite3
from pathlib import Path
from datetime import datetime, timedelta
from typing import Dict, List

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
SR_DB = VAULT / ".spaced-repetition.db"

class SpacedRepetitionAgent:
    def __init__(self):
        self.init_db()
    
    def init_db(self):
        conn = sqlite3.connect(SR_DB)
        c = conn.cursor()
        
        c.execute('''CREATE TABLE IF NOT EXISTS cards (
            id INTEGER PRIMARY KEY,
            note_path TEXT UNIQUE,
            easiness REAL DEFAULT 2.5,
            interval INTEGER DEFAULT 1,
            repetitions INTEGER DEFAULT 0,
            due_date TEXT,
            last_review TEXT
        )''')
        
        c.execute('''CREATE TABLE IF NOT EXISTS reviews (
            id INTEGER PRIMARY KEY,
            card_id INTEGER,
            quality INTEGER,
            timestamp TEXT,
            FOREIGN KEY(card_id) REFERENCES cards(id)
        )''')
        
        conn.commit()
        conn.close()
    
    def add_card(self, note_path: str):
        """Add note to spaced repetition"""
        conn = sqlite3.connect(SR_DB)
        c = conn.cursor()
        
        due = datetime.now().isoformat()
        c.execute('''INSERT OR IGNORE INTO cards (note_path, due_date, last_review)
                     VALUES (?, ?, ?)''', (note_path, due, due))
        
        conn.commit()
        conn.close()
    
    def review(self, note_path: str, quality: int) -> Dict:
        """
        Review card with SM-2 algorithm
        quality: 0-5 (0=complete blackout, 5=perfect)
        """
        conn = sqlite3.connect(SR_DB)
        c = conn.cursor()
        
        c.execute('SELECT * FROM cards WHERE note_path=?', (note_path,))
        card = c.fetchone()
        
        if not card:
            conn.close()
            return {"error": "Card not found"}
        
        card_id, path, ef, interval, reps, due, last = card
        
        # SM-2 algorithm
        if quality >= 3:
            if reps == 0:
                interval = 1
            elif reps == 1:
                interval = 6
            else:
                interval = int(interval * ef)
            reps += 1
        else:
            reps = 0
            interval = 1
        
        # Update easiness factor
        ef = ef + (0.1 - (5 - quality) * (0.08 + (5 - quality) * 0.02))
        ef = max(1.3, ef)
        
        # Calculate next review
        next_review = datetime.now() + timedelta(days=interval)
        
        # Update card
        c.execute('''UPDATE cards 
                     SET easiness=?, interval=?, repetitions=?, 
                         due_date=?, last_review=?
                     WHERE id=?''',
                  (ef, interval, reps, next_review.isoformat(), 
                   datetime.now().isoformat(), card_id))
        
        # Log review
        c.execute('''INSERT INTO reviews (card_id, quality, timestamp)
                     VALUES (?, ?, ?)''',
                  (card_id, quality, datetime.now().isoformat()))
        
        conn.commit()
        conn.close()
        
        return {
            "note": path,
            "quality": quality,
            "next_review": next_review.isoformat(),
            "interval_days": interval,
            "easiness": round(ef, 2)
        }
    
    def get_due(self, limit: int = 20) -> List[Dict]:
        """Get cards due for review"""
        conn = sqlite3.connect(SR_DB)
        c = conn.cursor()
        
        now = datetime.now().isoformat()
        c.execute('''SELECT note_path, due_date, interval, repetitions
                     FROM cards WHERE due_date <= ? 
                     ORDER BY due_date LIMIT ?''', (now, limit))
        
        due = [{"path": p, "due": d, "interval": i, "reps": r} 
               for p, d, i, r in c.fetchall()]
        
        conn.close()
        return due
    
    def stats(self) -> Dict:
        """Get SR statistics"""
        conn = sqlite3.connect(SR_DB)
        c = conn.cursor()
        
        c.execute('SELECT COUNT(*) FROM cards')
        total = c.fetchone()[0]
        
        now = datetime.now().isoformat()
        c.execute('SELECT COUNT(*) FROM cards WHERE due_date <= ?', (now,))
        due = c.fetchone()[0]
        
        c.execute('SELECT AVG(easiness), AVG(interval) FROM cards')
        avg_ef, avg_int = c.fetchone()
        
        conn.close()
        
        return {
            "total_cards": total,
            "due_today": due,
            "avg_easiness": round(avg_ef or 2.5, 2),
            "avg_interval": round(avg_int or 1, 1)
        }

if __name__ == "__main__":
    import sys
    import json
    
    sr = SpacedRepetitionAgent()
    
    if len(sys.argv) < 2:
        print("Usage: spaced-repetition.py <command> [args]")
        print("\nCommands:")
        print("  add <path>           - Add note to SR")
        print("  review <path> <0-5>  - Review card")
        print("  due [limit]          - Get due cards")
        print("  stats                - Show statistics")
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "add" and len(sys.argv) > 2:
        sr.add_card(sys.argv[2])
        print(f"Added: {sys.argv[2]}")
    
    elif cmd == "review" and len(sys.argv) > 3:
        result = sr.review(sys.argv[2], int(sys.argv[3]))
        print(json.dumps(result, indent=2))
    
    elif cmd == "due":
        limit = int(sys.argv[2]) if len(sys.argv) > 2 else 20
        due = sr.get_due(limit)
        print(json.dumps(due, indent=2))
    
    elif cmd == "stats":
        stats = sr.stats()
        print(json.dumps(stats, indent=2))
