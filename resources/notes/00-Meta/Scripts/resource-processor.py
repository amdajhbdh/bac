#!/usr/bin/env python3
"""
Resource Processing Agent
Extract exercises, create chapter notes, link to vault
"""
import subprocess
import re
from pathlib import Path
from typing import Dict, List
import json

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
RESOURCES = VAULT / "03-Resources"
EXTRACTED = VAULT / "05-Extracted/Resources"
EXERCISES = VAULT / "02-Practice/FromTextbooks"

class ResourceProcessor:
    def __init__(self):
        EXTRACTED.mkdir(parents=True, exist_ok=True)
        EXERCISES.mkdir(parents=True, exist_ok=True)
    
    def extract_pdf_text(self, pdf_path: Path) -> str:
        """Extract text from PDF"""
        result = subprocess.run(
            ['pdftotext', '-layout', str(pdf_path), '-'],
            capture_output=True, text=True, errors='ignore'
        )
        return result.stdout if result.returncode == 0 else ""
    
    def identify_chapters(self, text: str) -> List[Dict]:
        """Identify chapter structure"""
        chapters = []
        
        # Common chapter patterns
        patterns = [
            r'(?:Chapter|Chapitre|الفصل)\s*(\d+)[:\s]*(.+)',
            r'(?:Leçon|Lesson|درس)\s*(\d+)[:\s]*(.+)',
            r'^(\d+)\.\s*([A-Z][^.]+)$'
        ]
        
        lines = text.split('\n')
        for i, line in enumerate(lines):
            for pattern in patterns:
                match = re.search(pattern, line, re.IGNORECASE | re.MULTILINE)
                if match:
                    chapters.append({
                        'number': match.group(1),
                        'title': match.group(2).strip(),
                        'line': i
                    })
                    break
        
        return chapters
    
    def extract_exercises(self, text: str) -> List[Dict]:
        """Extract exercise problems"""
        exercises = []
        
        # Exercise patterns
        patterns = [
            r'(?:Exercise|Exercice|تمرين)\s*(\d+)[:\s]*(.+?)(?=(?:Exercise|Exercice|تمرين)|\Z)',
            r'(?:Problem|Problème)\s*(\d+)[:\s]*(.+?)(?=(?:Problem|Problème)|\Z)',
            r'^(\d+)\)\s*(.+?)(?=^\d+\)|\Z)'
        ]
        
        for pattern in patterns:
            matches = re.finditer(pattern, text, re.MULTILINE | re.DOTALL)
            for match in matches:
                exercises.append({
                    'number': match.group(1),
                    'content': match.group(2).strip()[:500]  # First 500 chars
                })
        
        return exercises
    
    def process_resource(self, pdf_path: Path) -> Dict:
        """Process single resource"""
        print(f"\n📖 Processing: {pdf_path.name}")
        
        # Extract text
        text = self.extract_pdf_text(pdf_path)
        if not text:
            return {"error": "Failed to extract text"}
        
        # Identify structure
        chapters = self.identify_chapters(text)
        exercises = self.extract_exercises(text)
        
        # Determine subject
        name = pdf_path.stem.lower()
        if any(x in name for x in ['math', 'رياضيات', 'matrice']):
            subject = 'Math'
        elif any(x in name for x in ['phys', 'فيزياء']):
            subject = 'Physics'
        elif any(x in name for x in ['chi', 'كيمياء']):
            subject = 'Chemistry'
        else:
            subject = 'General'
        
        # Create chapter notes
        chapter_notes = []
        for ch in chapters:
            note_path = self.create_chapter_note(pdf_path, ch, subject)
            if note_path:
                chapter_notes.append(str(note_path.relative_to(VAULT)))
        
        # Create exercise notes
        exercise_notes = []
        for ex in exercises[:50]:  # Limit to first 50
            note_path = self.create_exercise_note(pdf_path, ex, subject)
            if note_path:
                exercise_notes.append(str(note_path.relative_to(VAULT)))
        
        return {
            'pdf': pdf_path.name,
            'subject': subject,
            'chapters': len(chapters),
            'exercises': len(exercises),
            'chapter_notes': chapter_notes,
            'exercise_notes': exercise_notes
        }
    
    def create_chapter_note(self, pdf: Path, chapter: Dict, subject: str) -> Path:
        """Create note for chapter"""
        safe_title = re.sub(r'[^\w\s-]', '', chapter['title'])[:50]
        note_name = f"{pdf.stem}-Ch{chapter['number']}-{safe_title}.md"
        note_path = EXTRACTED / subject / note_name
        note_path.parent.mkdir(parents=True, exist_ok=True)
        
        content = f"""# {chapter['title']}

**Source:** [[{pdf.name}]] - Chapter {chapter['number']}  
**Subject:** {subject}  
**Type:** Textbook Chapter

---

## Overview

This chapter covers: {chapter['title']}

## Key Concepts

[To be filled by studying the chapter]

## Related Notes

- See [[{subject} MOC]] for related topics

## Exercises

See practice problems extracted from this chapter.

---

*Extracted from textbook*
"""
        
        note_path.write_text(content)
        return note_path
    
    def create_exercise_note(self, pdf: Path, exercise: Dict, subject: str) -> Path:
        """Create note for exercise"""
        note_name = f"{pdf.stem}-Ex{exercise['number']}.md"
        note_path = EXERCISES / subject / note_name
        note_path.parent.mkdir(parents=True, exist_ok=True)
        
        content = f"""# Exercise {exercise['number']}

**Source:** [[{pdf.name}]]  
**Subject:** {subject}  
**Difficulty:** [To be assessed]

---

## Problem

{exercise['content']}

## Solution

[Work through this problem]

## Notes

- Related concepts: 
- Key formulas used:

---

*From textbook exercises*
"""
        
        note_path.write_text(content)
        return note_path
    
    def create_resource_index(self, results: List[Dict]):
        """Create master resource index"""
        index_path = RESOURCES / "INDEX.md"
        
        content = """# Resource Index

**Total Resources:** {total}  
**Last Updated:** {date}

---

## By Subject

""".format(total=len(results), date=Path(__file__).stat().st_mtime)
        
        by_subject = {}
        for r in results:
            subj = r.get('subject', 'General')
            if subj not in by_subject:
                by_subject[subj] = []
            by_subject[subj].append(r)
        
        for subject, resources in sorted(by_subject.items()):
            content += f"\n### {subject}\n\n"
            for r in resources:
                content += f"**[[{r['pdf']}]]**\n"
                content += f"- Chapters: {r.get('chapters', 0)}\n"
                content += f"- Exercises: {r.get('exercises', 0)}\n"
                if r.get('chapter_notes'):
                    content += f"- Notes: {len(r['chapter_notes'])} created\n"
                content += "\n"
        
        content += """
---

## Quick Access

### All Exercises
See [[02-Practice/FromTextbooks/]]

### All Chapter Notes
See [[05-Extracted/Resources/]]

---

*Auto-generated resource index*
"""
        
        index_path.write_text(content)
        print(f"\n✓ Created index: {index_path}")
    
    def process_all(self):
        """Process all resources"""
        pdfs = list(RESOURCES.glob("*.pdf"))
        print(f"\n📚 Found {len(pdfs)} resources to process\n")
        
        results = []
        for pdf in pdfs:
            try:
                result = self.process_resource(pdf)
                results.append(result)
                print(f"  ✓ {result.get('chapters', 0)} chapters, {result.get('exercises', 0)} exercises")
            except Exception as e:
                print(f"  ✗ Error: {e}")
                results.append({'pdf': pdf.name, 'error': str(e)})
        
        # Create index
        self.create_resource_index(results)
        
        # Summary
        total_chapters = sum(r.get('chapters', 0) for r in results)
        total_exercises = sum(r.get('exercises', 0) for r in results)
        
        print(f"\n{'='*60}")
        print(f"✅ Processing Complete!")
        print(f"{'='*60}")
        print(f"Resources processed: {len(pdfs)}")
        print(f"Chapters identified: {total_chapters}")
        print(f"Exercises extracted: {total_exercises}")
        print(f"Notes created: Check 05-Extracted/Resources/ and 02-Practice/FromTextbooks/")
        print(f"Index: 03-Resources/INDEX.md")
        
        return results

if __name__ == "__main__":
    processor = ResourceProcessor()
    processor.process_all()
