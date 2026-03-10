#!/usr/bin/env python3
"""
Unified Study CLI
All systems in one interface
"""
import sys
import json
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent))

from spaced_repetition import SpacedRepetitionAgent
from progress_tracker import ProgressTracker
from ocr_processor import OCRProcessor

class StudyCLI:
    def __init__(self):
        self.sr = SpacedRepetitionAgent()
        self.progress = ProgressTracker()
        self.ocr = OCRProcessor()
    
    def review_session(self):
        """Interactive review session"""
        due = self.sr.get_due(20)
        
        if not due:
            print("🎉 No cards due! Great job!")
            return
        
        print(f"\n📚 {len(due)} cards due for review\n")
        
        for i, card in enumerate(due, 1):
            print(f"\n[{i}/{len(due)}] {card['path']}")
            print("Rate your recall (0-5):")
            print("  0 = Complete blackout")
            print("  3 = Correct with difficulty")
            print("  5 = Perfect recall")
            
            try:
                quality = int(input("Quality: "))
                result = self.sr.review(card['path'], quality)
                print(f"✓ Next review: {result['next_review'][:10]} ({result['interval_days']} days)")
            except (ValueError, KeyboardInterrupt):
                print("\nSession ended")
                break
        
        # Log session
        self.progress.log_session("Review", len(due) * 2, len(due), 0)
    
    def study_session(self, subject: str, duration: int):
        """Guided study session"""
        print(f"\n📖 Starting {duration}-minute {subject} session\n")
        
        # Get study plan from agent
        import subprocess
        result = subprocess.run(
            ['python3', '../agent-system.py', 'study_coach', 'create_plan',
             json.dumps({"duration": duration/60, "focus": [subject]})],
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            plan = json.loads(result.stdout)
            print(f"Plan: {len(plan.get('sessions', []))} topics\n")
            
            for session in plan.get('sessions', []):
                print(f"• {session['note']}")
                self.progress.mark_note_read(session['path'])
        
        # Log session
        notes = len(plan.get('sessions', []))
        self.progress.log_session(subject, duration, notes, 0)
        
        print(f"\n✓ Session complete! Logged {duration} minutes")
    
    def dashboard(self):
        """Show progress dashboard"""
        print(self.progress.generate_report())
        
        # SR stats
        sr_stats = self.sr.stats()
        print(f"\n## Spaced Repetition")
        print(f"- Cards: {sr_stats['total_cards']}")
        print(f"- Due today: {sr_stats['due_today']}")
        print(f"- Avg interval: {sr_stats['avg_interval']} days")

if __name__ == "__main__":
    cli = StudyCLI()
    
    if len(sys.argv) < 2:
        print("Study CLI - Unified Interface")
        print("\nCommands:")
        print("  review              - Start review session")
        print("  study <subject> <min> - Start study session")
        print("  dashboard           - Show progress")
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "review":
        cli.review_session()
    elif cmd == "study" and len(sys.argv) > 3:
        cli.study_session(sys.argv[2], int(sys.argv[3]))
    elif cmd == "dashboard":
        cli.dashboard()
