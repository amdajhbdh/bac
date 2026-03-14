#!/usr/bin/env python3
import os
import shutil
from pathlib import Path

vault_root = Path("/home/med/Documents/bac/resources/notes/Study-Vault")

# Define new structure
new_structure = {
    "00-Meta": ["INDEX.md", "START-HERE.md", "QUICK-START.md", "SYSTEM-STATUS.md"],
    "01-Concepts": ["Math", "Physics", "Chemistry", "Natural-Sciences"],
    "02-Practice": ["Math", "Physics", "Chemistry", "Natural-Sciences"],
    "03-Resources": ["PDFs"],
    "04-Exams": ["BAC-2002-2012", "School-Exams"],
    "05-Extracted": ["Math", "Physics", "Chemistry", "Biology", "OCR", "Resources"],
    "06-Daily": ["Progress-Dashboard.md"],
    "07-Assets": ["PDFs", "Images", "Diagrams"],
    "08-Templates": ["Concept-Note-Template.md", "Practice-Template.md", "Study-Session-Template.md"],
}

# Clean up duplicates and consolidate
cleanup_tasks = [
    ("03-Resources/*.pdf", "07-Assets/PDFs/"),
    ("Daily-Notes", "06-Daily/"),
    ("Assets/PDFs", "07-Assets/PDFs/"),
]

# Move PDFs from 03-Resources to 07-Assets/PDFs
src_pdfs = vault_root / "03-Resources"
dst_pdfs = vault_root / "07-Assets/PDFs"

if src_pdfs.exists():
    for pdf in src_pdfs.glob("*.pdf"):
        dst = dst_pdfs / pdf.name
        if not dst.exists():
            shutil.move(str(pdf), str(dst))
            print(f"Moved: {pdf.name}")

# Remove empty directories
for item in vault_root.rglob("*"):
    if item.is_dir() and not any(item.iterdir()):
        try:
            item.rmdir()
            print(f"Removed empty dir: {item.relative_to(vault_root)}")
        except:
            pass

print("✓ Vault restructuring complete")
