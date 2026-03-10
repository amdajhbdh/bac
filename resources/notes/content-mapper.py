#!/usr/bin/env python3
"""
Content Mapping Agent - Task B
Maps textbook chapters to vault topics, finds gaps and overlaps
"""
import sqlite3
from pathlib import Path
from typing import Dict, List
import re

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
DB = VAULT / ".vault-index.db"
EXTRACTED = VAULT / "05-Extracted/Resources"

class ContentMapper:
    def __init__(self):
        self.conn = sqlite3.connect(DB)
    
    def get_vault_topics(self) -> Dict[str, List[str]]:
        """Get existing vault topics by subject"""
        c = self.conn.cursor()
        c.execute("SELECT subject, title, path FROM notes WHERE type='concept'")
        
        topics = {}
        for subject, title, path in c.fetchall():
            if subject not in topics:
                topics[subject] = []
            topics[subject].append({'title': title, 'path': path})
        
        return topics
    
    def get_textbook_chapters(self) -> Dict[str, List[str]]:
        """Get all textbook chapters"""
        chapters = {}
        
        for subject_dir in EXTRACTED.iterdir():
            if subject_dir.is_dir():
                subject = subject_dir.name
                chapters[subject] = []
                
                for note in subject_dir.glob("*.md"):
                    chapters[subject].append({
                        'title': note.stem,
                        'path': str(note.relative_to(VAULT))
                    })
        
        return chapters
    
    def find_matches(self, vault_topics: Dict, textbook_chapters: Dict) -> Dict:
        """Find overlaps between vault and textbooks"""
        matches = {}
        
        for subject in textbook_chapters:
            if subject not in vault_topics:
                continue
            
            matches[subject] = []
            
            for chapter in textbook_chapters[subject]:
                ch_words = set(re.findall(r'\w+', chapter['title'].lower()))
                
                for topic in vault_topics[subject]:
                    topic_words = set(re.findall(r'\w+', topic['title'].lower()))
                    
                    # Find common words
                    common = ch_words & topic_words
                    if len(common) >= 2:  # At least 2 words match
                        matches[subject].append({
                            'chapter': chapter['title'],
                            'chapter_path': chapter['path'],
                            'topic': topic['title'],
                            'topic_path': topic['path'],
                            'match_score': len(common)
                        })
        
        return matches
    
    def find_gaps(self, vault_topics: Dict, textbook_chapters: Dict) -> Dict:
        """Find chapters in textbooks not covered in vault"""
        gaps = {}
        
        for subject, chapters in textbook_chapters.items():
            vault_titles = set()
            if subject in vault_topics:
                vault_titles = {t['title'].lower() for t in vault_topics[subject]}
            
            gaps[subject] = []
            for chapter in chapters:
                ch_words = set(re.findall(r'\w+', chapter['title'].lower()))
                
                # Check if any vault topic covers this
                covered = False
                for vault_title in vault_titles:
                    vault_words = set(re.findall(r'\w+', vault_title))
                    if len(ch_words & vault_words) >= 2:
                        covered = True
                        break
                
                if not covered:
                    gaps[subject].append(chapter)
        
        return gaps
    
    def create_cross_references(self, matches: Dict):
        """Add cross-reference links to notes"""
        for subject, match_list in matches.items():
            for match in match_list:
                # Add link to chapter note
                chapter_path = VAULT / match['chapter_path']
                if chapter_path.exists():
                    content = chapter_path.read_text()
                    if match['topic'] not in content:
                        # Add related note link
                        content += f"\n\n## Related Vault Notes\n\n- [[{match['topic']}]]\n"
                        chapter_path.write_text(content)
    
    def generate_report(self) -> str:
        """Generate mapping report"""
        vault_topics = self.get_vault_topics()
        textbook_chapters = self.get_textbook_chapters()
        matches = self.find_matches(vault_topics, textbook_chapters)
        gaps = self.find_gaps(vault_topics, textbook_chapters)
        
        report = f"""# Content Mapping Report

**Generated:** {Path(__file__).stat().st_mtime}

---

## Overview

**Vault Topics:** {sum(len(v) for v in vault_topics.values())}  
**Textbook Chapters:** {sum(len(v) for v in textbook_chapters.values())}  
**Matches Found:** {sum(len(v) for v in matches.values())}  
**Gaps Identified:** {sum(len(v) for v in gaps.values())}

---

## Matches (Reinforcement Opportunities)

"""
        
        for subject, match_list in matches.items():
            if match_list:
                report += f"\n### {subject}\n\n"
                for m in match_list[:10]:  # Top 10
                    report += f"- **{m['chapter'][:50]}** ↔ [[{m['topic']}]]\n"
        
        report += "\n---\n\n## Gaps (New Content to Study)\n\n"
        
        for subject, gap_list in gaps.items():
            if gap_list:
                report += f"\n### {subject} ({len(gap_list)} gaps)\n\n"
                for g in gap_list[:15]:  # Top 15
                    report += f"- {g['title'][:60]}\n"
        
        report += "\n---\n\n*Generated by Content Mapping Agent*\n"
        
        return report

if __name__ == "__main__":
    mapper = ContentMapper()
    
    print("🔍 Mapping content...")
    
    vault_topics = mapper.get_vault_topics()
    textbook_chapters = mapper.get_textbook_chapters()
    
    print(f"  Vault topics: {sum(len(v) for v in vault_topics.values())}")
    print(f"  Textbook chapters: {sum(len(v) for v in textbook_chapters.values())}")
    
    matches = mapper.find_matches(vault_topics, textbook_chapters)
    gaps = mapper.find_gaps(vault_topics, textbook_chapters)
    
    print(f"  Matches: {sum(len(v) for v in matches.values())}")
    print(f"  Gaps: {sum(len(v) for v in gaps.values())}")
    
    print("\n📝 Creating cross-references...")
    mapper.create_cross_references(matches)
    
    print("📄 Generating report...")
    report = mapper.generate_report()
    
    report_path = VAULT / "00-Meta/CONTENT-MAPPING.md"
    report_path.write_text(report)
    
    print(f"\n✅ Complete! Report: {report_path}")
