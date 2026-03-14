#!/usr/bin/env python3
"""
OCR Processing Pipeline
Extract text from PDFs and images
"""
import subprocess
from pathlib import Path
from typing import Dict
import json

VAULT = Path("/home/med/Documents/bac/resources/notes/Study-Vault")
ASSETS = VAULT / "07-Assets/PDFs"
EXTRACTED = VAULT / "05-Extracted"

class OCRProcessor:
    def process_pdf(self, pdf_path: Path) -> Dict:
        """Extract text from PDF using pdftotext"""
        output = pdf_path.stem + ".txt"
        output_path = Path("/tmp") / output
        
        # Try pdftotext first (faster)
        result = subprocess.run(
            ['pdftotext', str(pdf_path), str(output_path)],
            capture_output=True
        )
        
        if result.returncode == 0 and output_path.exists():
            text = output_path.read_text(errors='ignore')
            output_path.unlink()
            return {"text": text, "method": "pdftotext", "success": True}
        
        # Fallback to OCR with tesseract
        return self.ocr_pdf(pdf_path)
    
    def ocr_pdf(self, pdf_path: Path) -> Dict:
        """OCR PDF using tesseract"""
        # Convert PDF to images then OCR
        # Simplified - would use pdf2image + tesseract
        return {
            "text": f"[OCR placeholder for {pdf_path.name}]",
            "method": "tesseract",
            "success": False,
            "note": "Full OCR requires pdf2image + tesseract"
        }
    
    def process_image(self, img_path: Path) -> Dict:
        """Extract text from image"""
        output = img_path.stem
        
        result = subprocess.run(
            ['tesseract', str(img_path), output],
            capture_output=True,
            cwd='/tmp'
        )
        
        output_file = Path(f"/tmp/{output}.txt")
        if output_file.exists():
            text = output_file.read_text(errors='ignore')
            output_file.unlink()
            return {"text": text, "success": True}
        
        return {"text": "", "success": False}
    
    def create_note(self, pdf_path: Path, text: str, subject: str):
        """Create markdown note from extracted text"""
        note_name = pdf_path.stem + ".md"
        subject_dir = EXTRACTED / subject
        subject_dir.mkdir(exist_ok=True)
        
        note_path = subject_dir / note_name
        
        content = f"""# {pdf_path.stem}

**Source:** {pdf_path.name}  
**Extracted:** {Path(__file__).name}  
**Subject:** {subject}

---

## Content

{text[:5000]}

{'...(truncated)' if len(text) > 5000 else ''}

---

*Extracted from PDF*
"""
        
        note_path.write_text(content)
        return note_path
    
    def batch_process(self, subject: str = None, limit: int = 10) -> Dict:
        """Process multiple PDFs"""
        pdfs = list(ASSETS.glob("*.pdf"))[:limit]
        
        results = []
        for pdf in pdfs:
            result = self.process_pdf(pdf)
            if result['success'] and subject:
                note_path = self.create_note(pdf, result['text'], subject)
                results.append({
                    "pdf": pdf.name,
                    "note": str(note_path.relative_to(VAULT)),
                    "chars": len(result['text'])
                })
        
        return {"processed": len(results), "results": results}

if __name__ == "__main__":
    import sys
    
    processor = OCRProcessor()
    
    if len(sys.argv) < 2:
        print("Usage: ocr-processor.py <command> [args]")
        print("\nCommands:")
        print("  pdf <path>              - Process single PDF")
        print("  batch <subject> [limit] - Process multiple PDFs")
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "pdf" and len(sys.argv) > 2:
        result = processor.process_pdf(Path(sys.argv[2]))
        print(json.dumps(result, indent=2))
    
    elif cmd == "batch" and len(sys.argv) > 2:
        subject = sys.argv[2]
        limit = int(sys.argv[3]) if len(sys.argv) > 3 else 10
        result = processor.batch_process(subject, limit)
        print(json.dumps(result, indent=2))
