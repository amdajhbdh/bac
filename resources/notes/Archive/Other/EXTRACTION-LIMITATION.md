# PDF Extraction Limitation Report
**Date:** 2026-03-11 18:40:38

## ❌ Automated Extraction Failed

### Issue: Filename Encoding
**Problem:** PDF filename contains special character (é) causing tool failures
- Filename: `Moha Série solutions aqueuses 2024-1.pdf`
- Character: é (UTF-8: 0xC3 0xA9)
- All CLI tools fail to open file

### Tools Attempted
✅ Available: tesseract, pdftotext, ocrmypdf
❌ All failed: Filename encoding issue

### Workaround Options

**Option 1: Rename File (Recommended)**
```bash
cd 03-Resources
mv "Moha Série solutions aqueuses 2024-1.pdf" "Moha-Serie-solutions-aqueuses-2024-1.pdf"
# Then extract normally
```

**Option 2: Manual Extraction**
- Open PDF in viewer
- Copy text exercise by exercise
- Create markdown files manually
- Time: 2-3 hours

**Option 3: Google Lens (Fallback)**
- Upload PDF to Google Drive
- Use Google Docs OCR (Drive API)
- Requires API setup
- Time: 1 hour setup + extraction

## 📋 Recommendation

**Immediate:** Rename file to ASCII-safe name, then re-run extraction
**Alternative:** Manual extraction (most reliable for 6-page PDF)

**Status:** Awaiting user decision on approach

**Coordinator:** Technical limitation identified ✓
