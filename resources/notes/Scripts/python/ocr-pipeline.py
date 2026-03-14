#!/usr/bin/env python3
"""
OCR Pipeline - Extract text from all PDFs in the vault
Uses Tesseract OCR with Arabic + French + English support
"""

import os
import subprocess
import sys
from pathlib import Path

VAULT_DIR = Path("/home/med/Documents/bac/resources/notes")
OUTPUT_DIR = VAULT_DIR / "03-Resources" / "OCR"
OUTPUT_DIR.mkdir(parents=True, exist_ok=True)


def find_pdfs():
    """Find all PDFs in the vault"""
    patterns = [
        VAULT_DIR / "03-Resources" / "*.pdf",
        VAULT_DIR / "03-Resources" / "**" / "*.pdf",
        VAULT_DIR / "07-Assets" / "PDFs" / "*.pdf",
        VAULT_DIR / "mauritanian-bac" / "**" / "*.pdf",
    ]

    pdfs = set()
    for pattern in patterns:
        for pdf in VAULT_DIR.glob(str(pattern.relative_to(VAULT_DIR))):
            if pdf.is_file():
                pdfs.add(pdf)
    return sorted(pdfs)


def ocr_pdf(pdf_path: Path) -> Path:
    """OCR a single PDF and return the output markdown path"""
    output_base = OUTPUT_DIR / pdf_path.stem
    print(f"  Processing: {pdf_path.name}")

    # Convert PDF to PNG images
    subprocess.run(
        ["pdftoppm", "-r", "200", "-png", str(pdf_path), str(output_base)],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )

    # OCR each image
    text_parts = []
    for img in sorted(output_base.glob("*.png")):
        result = subprocess.run(
            ["tesseract", "-l", "ara+fra+eng", str(img), "stdout"],
            capture_output=True,
            text=True,
        )

        if result.stdout.strip():
            text_parts.append(result.stdout)

        img.unlink()  # Cleanup image

    # Write markdown with front matter
    output_md = output_base.with_suffix(".md")
    relative_source = pdf_path.relative_to(VAULT_DIR)

    frontmatter = f"""---
title: {pdf_path.stem}
date: {Path(__file__).stat().st_mtime}
tags: [ocr, extracted]
source: {relative_source}
---

"""

    with open(output_md, "w", encoding="utf-8") as f:
        f.write(frontmatter)
        f.write("\n\n".join(text_parts))

    print(f"    -> {output_md.name}")
    return output_md


def main():
    print("OCR Pipeline for BAC Study Vault")
    print("=" * 50)

    pdfs = find_pdfs()
    print(f"Found {len(pdfs)} PDFs to process\n")

    for i, pdf in enumerate(pdfs, 1):
        print(f"[{i}/{len(pdfs)}]")
        try:
            ocr_pdf(pdf)
        except Exception as e:
            print(f"    ERROR: {e}")

    print("\n" + "=" * 50)
    print(f"OCR complete! Output: {OUTPUT_DIR}")


if __name__ == "__main__":
    main()
