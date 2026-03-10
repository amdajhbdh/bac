# 📥 BAC PDF DOWNLOAD GUIDE

**Status:** Manual download guide (no authentication needed)  
**Last Updated:** March 10, 2026

---

## 🚀 QUICK START

### Option 1: Automatic Download (Python Script)
```bash
cd /home/med/Documents/bac/resources/notes
python3 bac-pdf-downloader.py
```

This will download all publicly available PDFs to `BAC-Exams-PDFs/` folder.

### Option 2: Manual Download (Browser)
Follow the links below and save to `BAC-Exams-PDFs/` folder.

---

## 📚 DIRECT DOWNLOAD LINKS

### MATHÉMATIQUES (Most Important - Highest Coefficient)

#### 2025 (Current Year - PRIORITY)
- **Bac C (M & TMGM):** https://maurimath.net/documents/sbac/BacC2025SN.pdf
- **Bac D (SN & TSGE):** https://maurimath.net/documents/sbac/BacD2025SN.pdf

#### 2024 (Recent - Important)
- **Bac C:** https://maurimath.net/documents/sbac/SujetBacC2024SN.pdf
- **Bac D:** https://maurimath.net/documents/sbac/SujetBacD2024SN.pdf

#### 2022 (Experimental)
- **Bac C Exp SN:** https://maurimath.net/documents/sbac/BacCExpSN2022.pdf
- **Bac C Exp SC:** https://maurimath.net/documents/sbac/BacCExpSC2022.pdf
- **Bac D Exp SC:** https://maurimath.net/documents/sbac/BacDExpSC2022.pdf

#### 2021
- **Bac C:** https://maurimath.net/documents/sbac/BacC2021SN.pdf

#### 2020
- **Bac C & TMGM:** https://maurimath.net/documents/sbac/BacC2020sn.pdf

#### Collections (2000-2022 with Corrections)
- **2000-2012 (C & TMGM) + Corrections:** https://www.sigmaths.net/bac2/mauritanie/baccmaths.pdf (195 pages)
- **2022 M-TMGM + Correction:** https://www.sigmaths.net/bac2/mauritanie/M_TMGM2022_sujet_correction.pdf
- **2020 & 2021 Mirror:** https://www.sigmaths.net/bac2/mauritanie/BacC2020sn.pdf

---

### PHYSIQUE-CHIMIE (Second Priority)

#### Collections (2005-2020)
- **BAC D (Sciences) 2005-2020:** https://prepabac.app/sc/dwnld/bac_d_pc.pdf (118 pages)
- **BAC C 2010-2020:** https://fr.scribd.com/document/826837707/bac-c-pc (120 pages)
- **BAC D 2010-2020:** https://fr.scribd.com/document/616665835/bac-d-pc

---

### SVT / SCIENCES NATURELLES (Third Priority)

#### Collections (2010-2024)
- **BAC C (Sciences Naturelles) 2010-2018:** https://prepabac.app/sc/dwnld/bac_c_sn.pdf
- **SN 7C Sujets Corrigés 2010-2024:** https://www.scribd.com/document/889600924/Bac-Sn-7c-Sujets-Corriges-2010-2024

---

## 📋 DOWNLOAD CHECKLIST

### Priority 1: Download These First (Mathématiques)
- [ ] 2025 Bac C
- [ ] 2025 Bac D
- [ ] 2024 Bac C
- [ ] 2024 Bac D
- [ ] 2000-2012 Collection (sigmaths)

### Priority 2: Download Next (Physique-Chimie)
- [ ] BAC D 2005-2020 (prepabac)
- [ ] BAC C 2010-2020 (Scribd)

### Priority 3: Download Last (SVT)
- [ ] BAC C 2010-2018 (prepabac)
- [ ] SN 7C 2010-2024 (Scribd)

---

## 🔧 AFTER DOWNLOADING

### Step 1: Organize Files
```
BAC-Exams-PDFs/
├── Mathématiques/
│   ├── 2025/
│   │   ├── C.pdf
│   │   └── D.pdf
│   ├── 2024/
│   │   ├── C.pdf
│   │   └── D.pdf
│   └── Collections/
│       └── 2000-2012.pdf
├── Physique-Chimie/
│   └── Collections/
│       └── BAC-D-2005-2020.pdf
└── SVT/
    └── Collections/
        └── BAC-C-2010-2018.pdf
```

### Step 2: Extract & Analyze
```bash
python3 pdf-extractor-analyzer.py
```

This will:
- Extract text from all PDFs
- Identify topics and exercises
- Create analysis report
- Save to `analysis-results.json`

### Step 3: Study
Use the extracted content with your study guides:
- `COMPREHENSIVE-BAC-STUDY-GUIDE.md`
- `DETAILED-TOPIC-BREAKDOWN.md`
- `DAILY-STUDY-SCHEDULE.md`

---

## ⚠️ IMPORTANT NOTES

### Scribd Documents
- Free account needed for full download
- Alternative: Use the "download" button if available
- Some documents may require subscription

### Direct Links
- All maurimath.net and prepabac.app links are direct (no authentication)
- No login required
- Safe to download

### File Sizes
- Individual papers: 1-5 MB
- Collections: 50-200 MB
- Total: ~500 MB for all

---

## 🆘 TROUBLESHOOTING

### Link Not Working?
1. Check if URL is correct (copy-paste carefully)
2. Try alternative mirror (sigmaths.net for maths)
3. Wait a few minutes and retry
4. Check your internet connection

### PDF Won't Open?
1. Make sure it downloaded completely
2. Try opening with different PDF reader
3. Check file size (should be > 100 KB)

### Python Script Fails?
1. Install requests: `pip install requests`
2. Install PyPDF2: `pip install PyPDF2`
3. Check internet connection
4. Run with: `python3 bac-pdf-downloader.py`

---

## 📊 WHAT YOU'LL GET

After downloading and analyzing:

✅ **Mathématiques:**
- 2025 & 2024 recent exams
- 2000-2012 historical collection
- Full corrections (sigmaths)
- ~50+ years of problems

✅ **Physique-Chimie:**
- 2005-2020 complete collection
- 118 pages of exams
- Both BAC C and D series

✅ **SVT:**
- 2010-2024 complete collection
- Corrections included
- All sessions (normale + complémentaire)

---

## 🎯 NEXT STEPS

1. **Download** all PDFs using the links above
2. **Organize** into BAC-Exams-PDFs/ folder
3. **Extract** using `python3 pdf-extractor-analyzer.py`
4. **Analyze** the results
5. **Study** using your 12-week plan

---

## 💡 TIPS

- Start with **2025 & 2024 Maths** (most relevant to current program)
- Download **collections** for historical patterns
- Use **corrections** (sigmaths) to understand solutions
- Extract **French vocabulary** from actual exams
- Practice **timed conditions** (4 hours per exam)

---

**Ready to download? Start with the Priority 1 links above! 🚀**

**Good luck! 🍀✨**
