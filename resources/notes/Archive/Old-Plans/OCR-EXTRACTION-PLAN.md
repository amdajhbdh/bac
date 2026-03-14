# OCR Extraction Plan
**Created:** 2026-03-11 18:37:06
**Status:** Ready for execution

---

## 🚨 IMMEDIATE (Priority 1)

### Chemistry - Moha Série
**File:** `03-Resources/Moha Série solutions aqueuses 2024-1.pdf`
**Size:** 468 KB
**Pages:** ~10-15 (estimated)
**Target:** `05-Extracted/Chemistry/Moha-Serie-Solutions-Aqueuses-2024.md`

**Extraction Method:**
```bash
# Option 1: OCR with tesseract (if available)
tesseract input.pdf output -l fra

# Option 2: pdftotext (if text-based)
pdftotext "Moha Série..." output.txt

# Option 3: Manual extraction
# Open PDF, copy exercises to markdown
```

**Expected Content:**
- 15-25 exercises on aqueous solutions
- Topics: pH, buffers, titrations, equilibria
- Difficulty: Level 4/5 (Advanced)

**Post-Extraction:**
1. Create exercise files in `02-Practice/Chemistry/Solutions/`
2. Link to Buffer-Solutions.md, Titration-Advanced.md
3. Update Chemistry MOC
4. Mark TODO as complete

---

## 📅 SHORT-TERM (Priority 2)

### BAC Exams - Mathematics (14 files)

**Location:** `mauritanian-bac/mathematiques/`
**Target:** `04-Exams/BAC-Recent/`

| File | Year | Type | Priority |
|------|------|------|----------|
| BacD2025SN.pdf | 2025 | Normal | ⭐⭐⭐⭐⭐ |
| BacC2025SN.pdf | 2025 | Normal | ⭐⭐⭐⭐⭐ |
| BacD2024SN.pdf | 2024 | Normal | ⭐⭐⭐⭐⭐ |
| BacC2024SN.pdf | 2024 | Normal | ⭐⭐⭐⭐⭐ |
| BacCExpSN2022.pdf | 2022 | Comp | ⭐⭐⭐⭐ |
| BacCExpSC2022.pdf | 2022 | Comp | ⭐⭐⭐⭐ |
| BacDExpSC2022.pdf | 2022 | Comp | ⭐⭐⭐⭐ |
| BacC2021SN.pdf | 2021 | Normal | ⭐⭐⭐⭐ |
| BacC2020sn.pdf | 2020 | Normal | ⭐⭐⭐ |
| BacD2020SN.pdf | 2020 | Normal | ⭐⭐⭐ |
| BacC2018comp.pdf | 2018 | Comp | ⭐⭐ |
| BacLM2019.pdf | 2019 | LM | ⭐⭐ |
| baccmaths_2000_2012.pdf | Archive | Multi | ⭐⭐ |
| annale_serie_D_2005_2014_corr.pdf | Archive | Corr | ⭐⭐⭐ |

**Extraction Strategy:**
- Start with 2024-2025 (most relevant)
- Extract to markdown format
- Organize by subject (Math/Physics/Chemistry/Biology)
- Create index file

### BAC Exams - Physics & Biology (2 files)

| File | Subject | Priority |
|------|---------|----------|
| bac_d_pc_2005_2020.pdf | Physics-Chemistry | ⭐⭐⭐⭐ |
| bac_c_sn_2010_2018.pdf | Biology | ⭐⭐⭐⭐ |

---

## 📚 LONG-TERM (Priority 3)

### Textbooks (5 files)

| File | Subject | Size | Priority |
|------|---------|------|----------|
| 7c_bac_pc_كامل.pdf | Physics-Chem | Large | ⭐⭐⭐ |
| Les_Cahiers_De_Math*.pdf (2) | Math | Medium | ⭐⭐ |
| Livre-PHYSIQUE_TERMINALE_S*.pdf (2) | Physics | Large | ⭐⭐ |

**Note:** Large textbooks - extract key chapters only

---

## 🧹 CLEANUP (Priority 4)

### Duplicate PDFs to Delete (11)

**07-Assets/PDFs:**
- Joussour SN 7D7C 2024-2.pdf
- Joussour-PC-7M-SN vf (1)-2.pdf
- JoussourMaths7C-2.pdf
- Raja ressort 2024-2.pdf
- Système nerveux-2.pdf
- série projectile -1.pdf
- série projectile -2.pdf

**Action:**
```bash
cd 07-Assets/PDFs
rm *-2.pdf *-1.pdf
```

---

## 📊 Execution Timeline

**Week 1 (Immediate):**
- [ ] Extract Moha Série Chemistry
- [ ] Create 15-25 exercise files
- [ ] Update Chemistry section
- [ ] Complete Chemistry upgrade

**Week 2-3 (Short-term):**
- [ ] Extract BAC 2024-2025 (4 exams)
- [ ] Extract BAC 2020-2022 (6 exams)
- [ ] Extract Physics & Biology BAC (2 exams)
- [ ] Create exam index

**Week 4+ (Long-term):**
- [ ] Extract key textbook chapters
- [ ] Delete duplicate PDFs
- [ ] Final organization

---

## 🛠️ Tools Required

**Available:**
- Manual extraction (copy-paste)
- Text editor

**Recommended (if available):**
- tesseract-ocr (OCR engine)
- pdftotext (text extraction)
- pdf2image (image conversion)

**Fallback:**
- Manual transcription
- Screenshot + type

---

## ✅ Success Criteria

**Immediate:**
- Moha Série extracted → Chemistry upgrade complete

**Short-term:**
- 16 BAC exams available → Exam prep ready

**Long-term:**
- All resources organized → Vault 100% complete

**Coordinator:** Extraction plan ready for execution ✓
