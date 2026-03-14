#!/usr/bin/env python3
"""
BAC PDF Extractor & Analyzer
Extracts text from PDFs and analyzes exam content
"""

import os
import json
from pathlib import Path
from collections import Counter
import re

try:
    import PyPDF2
    PDF_AVAILABLE = True
except ImportError:
    PDF_AVAILABLE = False
    print("⚠️  PyPDF2 not installed. Install with: pip install PyPDF2")

class PDFAnalyzer:
    def __init__(self, pdf_dir):
        self.pdf_dir = Path(pdf_dir)
        self.results = {
            'files_processed': 0,
            'total_pages': 0,
            'topics': Counter(),
            'vocabulary': Counter(),
            'exercises': [],
            'errors': []
        }
    
    def extract_text_from_pdf(self, pdf_path):
        """Extract text from PDF file"""
        if not PDF_AVAILABLE:
            return None
        
        try:
            text = ""
            with open(pdf_path, 'rb') as f:
                reader = PyPDF2.PdfReader(f)
                for page in reader.pages:
                    text += page.extract_text() + "\n"
            return text
        except Exception as e:
            self.results['errors'].append(f"{pdf_path.name}: {str(e)}")
            return None
    
    def analyze_text(self, text, filename):
        """Analyze extracted text"""
        if not text:
            return
        
        # Count pages (rough estimate)
        pages = len(text.split('\n')) // 30
        self.results['total_pages'] += pages
        
        # Extract exercises
        exercises = re.findall(r'Exercice\s+(\d+)', text, re.IGNORECASE)
        if exercises:
            self.results['exercises'].append({
                'file': filename,
                'count': len(set(exercises)),
                'numbers': sorted(set(exercises))
            })
        
        # Extract topics
        topics = {
            'acide-base': r'acide|base|pH|pKa|titration',
            'organique': r'alcool|aldéhyde|cétone|ester|amine',
            'cinétique': r'cinétique|vitesse.*réaction|catalyseur',
            'redox': r'oxydation|réduction|redox',
            'électromagnétisme': r'champ.*magnétique|induction|bobine',
            'circuits': r'RLC|résonance|impédance',
            'mécanique': r'énergie|ressort|oscillation',
        }
        
        for topic, pattern in topics.items():
            count = len(re.findall(pattern, text, re.IGNORECASE))
            if count > 0:
                self.results['topics'][topic] += count
        
        # Extract French vocabulary
        words = re.findall(r'\b[a-zàâäçéèêëïîôùûü]{4,}\b', text.lower())
        for word in words[:100]:  # Sample
            self.results['vocabulary'][word] += 1
    
    def process_all_pdfs(self):
        """Process all PDFs in directory"""
        print("📄 PDF Analyzer")
        print("=" * 60)
        
        pdf_files = list(self.pdf_dir.rglob('*.pdf'))
        
        if not pdf_files:
            print(f"❌ No PDF files found in {self.pdf_dir}")
            return
        
        print(f"Found {len(pdf_files)} PDF files\n")
        
        for pdf_path in pdf_files:
            print(f"Processing: {pdf_path.relative_to(self.pdf_dir)}")
            text = self.extract_text_from_pdf(pdf_path)
            if text:
                self.analyze_text(text, pdf_path.name)
                self.results['files_processed'] += 1
                print(f"  ✅ Extracted and analyzed")
            else:
                print(f"  ❌ Failed to extract")
        
        self._print_summary()
    
    def _print_summary(self):
        """Print analysis summary"""
        print("\n" + "=" * 60)
        print("📊 ANALYSIS SUMMARY")
        print("=" * 60)
        print(f"Files processed: {self.results['files_processed']}")
        print(f"Total pages: {self.results['total_pages']}")
        
        print(f"\n🎯 TOP TOPICS:")
        for topic, count in self.results['topics'].most_common(10):
            print(f"  {topic}: {count} occurrences")
        
        print(f"\n📝 EXERCISES FOUND:")
        for ex in self.results['exercises']:
            print(f"  {ex['file']}: {ex['count']} exercises")
        
        if self.results['errors']:
            print(f"\n⚠️  ERRORS ({len(self.results['errors'])}):")
            for error in self.results['errors'][:5]:
                print(f"  - {error}")
        
        # Save results
        output_file = self.pdf_dir / 'analysis-results.json'
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump({
                'files_processed': self.results['files_processed'],
                'total_pages': self.results['total_pages'],
                'topics': dict(self.results['topics'].most_common(20)),
                'exercises': self.results['exercises'],
                'errors': self.results['errors']
            }, f, ensure_ascii=False, indent=2)
        
        print(f"\n💾 Results saved: {output_file}")

def main():
    pdf_dir = Path(__file__).parent / "BAC-Exams-PDFs"
    
    if not pdf_dir.exists():
        print(f"❌ Directory not found: {pdf_dir}")
        print(f"Run bac-pdf-downloader.py first to download PDFs")
        return
    
    analyzer = PDFAnalyzer(pdf_dir)
    analyzer.process_all_pdfs()

if __name__ == '__main__':
    main()
