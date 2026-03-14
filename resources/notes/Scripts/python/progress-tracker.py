#!/usr/bin/env python3
"""
Progress Tracking System
Real-time study analytics
"""
import sqlite3
from pathlib import Path
from datetime import datetime, timedelta
from typing import Dict

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
PROGRESS_DB = VAULT / ".progress.db"

class ProgressTracker:
    def __init__(self):
        self.init_db()
    
    def init_db(self):
        conn = sqlite3.connect(PROGRESS_DB)
        c = conn.cursor()
        
        c.execute('''CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY,
            date TEXT,
            subject TEXT,
            duration_min INTEGER,
            notes_read INTEGER,
            problems_solved INTEGER
        )''')
        
        c.execute('''CREATE TABLE IF NOT EXISTS note_progress (
            note_path TEXT PRIMARY KEY,
            status TEXT,
            last_viewed TEXT,
            view_count INTEGER DEFAULT 0
        )''')
        
        c.execute('''CREATE TABLE IF NOT EXISTS problem_progress (
            problem_id TEXT PRIMARY KEY,
            note_path TEXT,
            solved BOOLEAN,
            attempts INTEGER DEFAULT 0,
            last_attempt TEXT
        )''')
        
        conn.commit()
        conn.close()
    
    def log_session(self, subject: str, duration: int, notes: int, problems: int):
        """Log study session"""
        conn = sqlite3.connect(PROGRESS_DB)
        c = conn.cursor()
        
        c.execute('''INSERT INTO sessions (date, subject, duration_min, notes_read, problems_solved)
                     VALUES (?, ?, ?, ?, ?)''',
                  (datetime.now().isoformat(), subject, duration, notes, problems))
        
        conn.commit()
        conn.close()
    
    def mark_note_read(self, path: str):
        """Mark note as read"""
        conn = sqlite3.connect(PROGRESS_DB)
        c = conn.cursor()
        
        c.execute('''INSERT INTO note_progress (note_path, status, last_viewed, view_count)
                     VALUES (?, 'read', ?, 1)
                     ON CONFLICT(note_path) DO UPDATE SET
                     last_viewed=?, view_count=view_count+1''',
                  (path, datetime.now().isoformat(), datetime.now().isoformat()))
        
        conn.commit()
        conn.close()
    
    def mark_problem_solved(self, problem_id: str, note_path: str):
        """Mark problem as solved"""
        conn = sqlite3.connect(PROGRESS_DB)
        c = conn.cursor()
        
        c.execute('''INSERT INTO problem_progress (problem_id, note_path, solved, attempts, last_attempt)
                     VALUES (?, ?, 1, 1, ?)
                     ON CONFLICT(problem_id) DO UPDATE SET
                     solved=1, attempts=attempts+1, last_attempt=?''',
                  (problem_id, note_path, datetime.now().isoformat(), datetime.now().isoformat()))
        
        conn.commit()
        conn.close()
    
    def get_dashboard(self) -> Dict:
        """Generate progress dashboard"""
        conn = sqlite3.connect(PROGRESS_DB)
        c = conn.cursor()
        
        # Today's stats
        today = datetime.now().date().isoformat()
        c.execute('''SELECT SUM(duration_min), SUM(notes_read), SUM(problems_solved)
                     FROM sessions WHERE date LIKE ?''', (f'{today}%',))
        today_stats = c.fetchone()
        
        # Week stats
        week_ago = (datetime.now() - timedelta(days=7)).isoformat()
        c.execute('''SELECT SUM(duration_min), SUM(notes_read), SUM(problems_solved)
                     FROM sessions WHERE date >= ?''', (week_ago,))
        week_stats = c.fetchone()
        
        # By subject
        c.execute('''SELECT subject, SUM(duration_min), SUM(problems_solved)
                     FROM sessions GROUP BY subject''')
        by_subject = {s: {"time": t, "problems": p} for s, t, p in c.fetchall()}
        
        # Streak
        c.execute('''SELECT DISTINCT date(date) FROM sessions 
                     ORDER BY date DESC LIMIT 30''')
        dates = [r[0] for r in c.fetchall()]
        streak = 0
        for i, d in enumerate(dates):
            expected = (datetime.now().date() - timedelta(days=i)).isoformat()
            if d == expected:
                streak += 1
            else:
                break
        
        conn.close()
        
        return {
            "today": {
                "time_min": today_stats[0] or 0,
                "notes": today_stats[1] or 0,
                "problems": today_stats[2] or 0
            },
            "week": {
                "time_min": week_stats[0] or 0,
                "notes": week_stats[1] or 0,
                "problems": week_stats[2] or 0
            },
            "by_subject": by_subject,
            "streak_days": streak
        }
    
    def generate_report(self) -> str:
        """Generate markdown report"""
        dash = self.get_dashboard()
        
        report = f"""# Progress Dashboard

**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M')}

## Today
- ⏱️ Time: {dash['today']['time_min']} minutes
- 📖 Notes read: {dash['today']['notes']}
- ✅ Problems solved: {dash['today']['problems']}

## This Week
- ⏱️ Time: {dash['week']['time_min']} minutes
- 📖 Notes read: {dash['week']['notes']}
- ✅ Problems solved: {dash['week']['problems']}

## Streak
🔥 **{dash['streak_days']} days**

## By Subject
"""
        for subject, stats in dash['by_subject'].items():
            report += f"\n### {subject}\n"
            report += f"- Time: {stats['time']} min\n"
            report += f"- Problems: {stats['problems']}\n"
        
        return report

if __name__ == "__main__":
    import sys
    
    tracker = ProgressTracker()
    
    if len(sys.argv) < 2:
        print("Usage: progress-tracker.py <command> [args]")
        print("\nCommands:")
        print("  session <subject> <duration> <notes> <problems>")
        print("  read <path>")
        print("  solve <problem_id> <note_path>")
        print("  dashboard")
        print("  report")
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "session" and len(sys.argv) > 5:
        tracker.log_session(sys.argv[2], int(sys.argv[3]), int(sys.argv[4]), int(sys.argv[5]))
        print("Session logged")
    
    elif cmd == "read" and len(sys.argv) > 2:
        tracker.mark_note_read(sys.argv[2])
        print(f"Marked as read: {sys.argv[2]}")
    
    elif cmd == "solve" and len(sys.argv) > 3:
        tracker.mark_problem_solved(sys.argv[2], sys.argv[3])
        print(f"Problem solved: {sys.argv[2]}")
    
    elif cmd == "dashboard":
        import json
        print(json.dumps(tracker.get_dashboard(), indent=2))
    
    elif cmd == "report":
        print(tracker.generate_report())
