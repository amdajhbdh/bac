#!/bin/bash

# ============================================================================
# BAC EXAM COLLECTION - EXECUTABLE GUIDE
# Multi-Simultaneous Search & Collection Strategy
# ============================================================================

echo "🔍 BAC EXAM COLLECTION SYSTEM - STARTING..."
echo "=================================================="
echo ""

# Create directory structure
echo "📁 Creating folder structure..."
mkdir -p ~/BAC-Exams-Complete/{2002..2025}/{Normal,Complementary}
mkdir -p ~/BAC-Exams-Complete/INDEX
mkdir -p ~/BAC-Exams-Complete/METADATA

echo "✓ Folders created at: ~/BAC-Exams-Complete/"
echo ""

# ============================================================================
# DIRECT WEBSITE LINKS - PRIORITY 1
# ============================================================================

echo "🌐 DIRECT WEBSITE LINKS (Copy & Paste in Browser):"
echo "=================================================="
echo ""
echo "1. MAURITANIAN MATH SITE (HIGHEST PRIORITY)"
echo "   URL: https://maurimath.net/Sujetbac.php"
echo "   Content: BAC exam subjects (2002-2024)"
echo "   Action: Download all PDFs"
echo ""

echo "2. MAURITANIAN STUDENT PORTAL"
echo "   URL: https://rimbac.com"
echo "   Content: Exam papers and solutions"
echo "   Action: Search 'exam' or 'examen'"
echo ""

echo "3. SENEGAL EXAM ARCHIVE"
echo "   URL: https://examens.sn/mauritanie/"
echo "   Content: West African exams"
echo "   Action: Browse Mauritania section"
echo ""

echo "4. AFRICAN EDUCATION PLATFORM"
echo "   URL: https://edukely.com/mauritania/"
echo "   Content: Educational resources"
echo "   Action: Search for BAC papers"
echo ""

echo "5. ARCHIVE.ORG WAYBACK MACHINE"
echo "   URL: https://archive.org/web/"
echo "   Search: maurimath.net OR rimbac.com"
echo "   Action: Find historical versions"
echo ""

# ============================================================================
# TERMINAL COMMANDS FOR SEARCHING
# ============================================================================

echo "💻 TERMINAL COMMANDS (Copy & Paste):"
echo "=================================================="
echo ""

echo "COMMAND 1: Search maurimath.net for PDFs"
echo "---"
echo "curl -s 'https://maurimath.net/Sujetbac.php' | grep -oP 'href=\"\K[^\"]*\.pdf' | head -20"
echo ""

echo "COMMAND 2: Search rimbac.com for exam links"
echo "---"
echo "curl -s 'https://rimbac.com' | grep -i 'exam\|sujet\|bac' | head -20"
echo ""

echo "COMMAND 3: Download all PDFs from maurimath (if available)"
echo "---"
echo "wget -r -A.pdf 'https://maurimath.net/Sujetbac.php' -P ~/BAC-Exams-Complete/"
echo ""

echo "COMMAND 4: Search for PDF files with 'BAC' in name"
echo "---"
echo "find ~ -name '*BAC*.pdf' -o -name '*bac*.pdf' 2>/dev/null"
echo ""

echo "COMMAND 5: List all downloaded files"
echo "---"
echo "ls -lh ~/BAC-Exams-Complete/"
echo ""

# ============================================================================
# GOOGLE SEARCH COMMANDS
# ============================================================================

echo "🔎 GOOGLE SEARCH QUERIES (Copy & Paste in Google):"
echo "=================================================="
echo ""

echo "QUERY 1: Find all BAC PDFs"
echo "---"
echo "\"BAC Mauritanie\" filetype:pdf 2002..2025"
echo ""

echo "QUERY 2: Search maurimath.net"
echo "---"
echo "site:maurimath.net \"Sujet BAC\" filetype:pdf"
echo ""

echo "QUERY 3: Search rimbac.com"
echo "---"
echo "site:rimbac.com \"exam\" OR \"examen\" 2024"
echo ""

echo "QUERY 4: Search by year"
echo "---"
echo "\"BAC Mauritanie\" \"2002\" OR \"2003\" OR \"2004\" filetype:pdf"
echo ""

echo "QUERY 5: Search by subject"
echo "---"
echo "\"Sujets BAC\" \"7C\" OR \"7D\" filetype:pdf"
echo ""

echo "QUERY 6: Search Scribd"
echo "---"
echo "site:scribd.com \"BAC Mauritanie\""
echo ""

echo "QUERY 7: Search Archive.org"
echo "---"
echo "site:archive.org \"maurimath\" OR \"rimbac\""
echo ""

# ============================================================================
# ORGANIZATION STRUCTURE
# ============================================================================

echo ""
echo "📁 FOLDER ORGANIZATION STRUCTURE:"
echo "=================================================="
echo ""

cat << 'STRUCTURE'
~/BAC-Exams-Complete/
├── 2002/
│   ├── Normal/
│   │   ├── Mathematics.pdf
│   │   ├── Physics.pdf
│   │   ├── Chemistry.pdf
│   │   └── Biology.pdf
│   └── Complementary/
│       ├── Mathematics.pdf
│       ├── Physics.pdf
│       ├── Chemistry.pdf
│       └── Biology.pdf
├── 2003/
│   ├── Normal/
│   └── Complementary/
├── 2004/
│   └── ...
├── ...
├── 2024/
│   ├── Normal/
│   └── Complementary/
├── INDEX/
│   └── INDEX.md (metadata file)
└── METADATA/
    └── collection_log.txt
STRUCTURE

echo ""

# ============================================================================
# FILE NAMING CONVENTION
# ============================================================================

echo "📝 FILE NAMING CONVENTION:"
echo "=================================================="
echo ""
echo "Format: BAC_[YEAR]_[SESSION]_[SUBJECT]_[LANGUAGE].pdf"
echo ""
echo "Examples:"
echo "  BAC_2024_Normal_Mathematics_French.pdf"
echo "  BAC_2024_Complementary_Physics_French.pdf"
echo "  BAC_2023_Normal_Chemistry_Arabic.pdf"
echo "  BAC_2022_Normal_Biology_French.pdf"
echo ""

# ============================================================================
# COLLECTION CHECKLIST
# ============================================================================

echo "✅ COLLECTION CHECKLIST:"
echo "=================================================="
echo ""

cat << 'CHECKLIST'
STEP 1: IDENTIFY SOURCES (5 min)
  [ ] Visit maurimath.net/Sujetbac.php
  [ ] Visit rimbac.com
  [ ] Visit examens.sn
  [ ] Visit edukely.com
  [ ] Google search results

STEP 2: COLLECT LINKS (10 min)
  [ ] Create list of all PDF links
  [ ] Verify links are working
  [ ] Note year and subject
  [ ] Check file size

STEP 3: DOWNLOAD FILES (20 min)
  [ ] Download all PDFs
  [ ] Organize by year/subject
  [ ] Verify file integrity
  [ ] Check for duplicates

STEP 4: ORGANIZE (15 min)
  [ ] Create folder structure
  [ ] Rename files consistently
  [ ] Add metadata
  [ ] Create index

STEP 5: VERIFY (10 min)
  [ ] Check all years covered (2002-2025)
  [ ] Verify both sessions (Normal & Complementary)
  [ ] Confirm all subjects included
  [ ] Test searchability
CHECKLIST

echo ""

# ============================================================================
# METADATA TEMPLATE
# ============================================================================

echo "📋 METADATA TEMPLATE (Save as INDEX.md):"
echo "=================================================="
echo ""

cat << 'METADATA'
# BAC Exam Collection Index

## Collection Date: [DATE]
## Total Papers: [NUMBER]
## Total Size: [SIZE]

### By Year:
- 2002: [X] papers
- 2003: [X] papers
- ...
- 2024: [X] papers

### By Subject:
- Mathematics: [X] papers
- Physics: [X] papers
- Chemistry: [X] papers
- Biology: [X] papers
- Other: [X] papers

### By Session:
- Normal: [X] papers
- Complementary: [X] papers

### By Language:
- French: [X] papers
- Arabic: [X] papers

### Files:
| Year | Session | Subject | Language | File | Status |
|------|---------|---------|----------|------|--------|
| 2024 | Normal | Math | FR | BAC_2024_Normal_Mathematics_French.pdf | ✓ |
| 2024 | Normal | Physics | FR | BAC_2024_Normal_Physics_French.pdf | ✓ |
| ... | ... | ... | ... | ... | ... |
METADATA

echo ""

# ============================================================================
# EXPECTED RESULTS
# ============================================================================

echo "📊 EXPECTED RESULTS:"
echo "=================================================="
echo ""
echo "Total Papers to Find:"
echo "  • Years: 24 (2002-2025)"
echo "  • Sessions: 2 per year (Normal + Complementary)"
echo "  • Subjects: 4-5 per session"
echo "  • Total: 46+ complete exam sets"
echo "  • Total Files: 200+ individual papers"
echo ""
echo "By Subject:"
echo "  • Mathematics: 48 papers"
echo "  • Physics: 48 papers"
echo "  • Chemistry: 48 papers"
echo "  • Biology: 36 papers"
echo "  • Other: 24 papers"
echo "  • TOTAL: 204 papers"
echo ""

# ============================================================================
# QUICK START COMMANDS
# ============================================================================

echo "🚀 QUICK START (Copy & Paste These):"
echo "=================================================="
echo ""

echo "1. Create directories:"
echo "   mkdir -p ~/BAC-Exams-Complete/{2002..2025}/{Normal,Complementary}"
echo ""

echo "2. Search maurimath.net:"
echo "   curl -s 'https://maurimath.net/Sujetbac.php' | grep -i 'pdf\\|download'"
echo ""

echo "3. List created folders:"
echo "   ls -la ~/BAC-Exams-Complete/"
echo ""

echo "4. Create index file:"
echo "   touch ~/BAC-Exams-Complete/INDEX/INDEX.md"
echo ""

# ============================================================================
# TROUBLESHOOTING
# ============================================================================

echo ""
echo "🔧 TROUBLESHOOTING:"
echo "=================================================="
echo ""

cat << 'TROUBLESHOOTING'
If links don't work:
  → Try Archive.org Wayback Machine
  → Search Google for cached versions
  → Check alternative sites (rimbac.com, examens.sn)
  → Look for YouTube videos showing papers

If PDFs are in Arabic:
  → Use Google Translate
  → Keep both versions
  → Note language in filename

If files are corrupted:
  → Try re-downloading
  → Check alternative sources
  → Verify file integrity with: md5sum filename.pdf

If download is slow:
  → Use wget with multiple connections
  → Download during off-peak hours
  → Try different mirror sites
TROUBLESHOOTING

echo ""

# ============================================================================
# FINAL SUMMARY
# ============================================================================

echo "📌 FINAL SUMMARY:"
echo "=================================================="
echo ""
echo "✓ Folder structure created at: ~/BAC-Exams-Complete/"
echo "✓ Direct website links provided above"
echo "✓ Terminal commands ready to copy & paste"
echo "✓ Google search queries ready to use"
echo "✓ Organization structure defined"
echo "✓ Metadata template provided"
echo ""
echo "⏱️  Expected Collection Time: 60 minutes"
echo "💾 Expected Storage: 500 MB - 1 GB"
echo ""
echo "🎯 NEXT STEPS:"
echo "1. Copy a direct website link and open in browser"
echo "2. Copy a terminal command and run it"
echo "3. Copy a Google search query and search"
echo "4. Download PDFs to ~/BAC-Exams-Complete/"
echo "5. Organize by year and subject"
echo "6. Create INDEX.md with metadata"
echo ""
echo "✅ READY TO COLLECT! 🚀"
echo ""

