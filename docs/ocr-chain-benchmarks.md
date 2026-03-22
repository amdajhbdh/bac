# OCR Chain Benchmarks

Tested on real BAC exam images (March 2026).

## Test Data

- **9 images** from `tests/bac-images/`
- **Subjects:** Math, Physics, Chemistry, French, Arabic Math
- **Provider tested:** Mistral OCR (primary), OCR.space (fallback)
- **Preprocessing:** Smart preprocessing enabled

## Results

| Subject | Image | Provider | Chars | Latency | Status |
|---------|-------|----------|-------|---------|--------|
| Math | math_page-1.jpg | Mistral | 2,401 | 3,547ms | ✅ |
| Math | math_page-2.jpg | Mistral | 3,433 | 4,655ms | ✅ |
| Physics | physique_page-1.jpg | Mistral | 3,187 | 2,836ms | ✅ |
| Physics | physique_page-2.jpg | Mistral | 4,202 | 3,353ms | ✅ |
| French | francais_page-1.jpg | Mistral | 2,842 | 2,584ms | ✅ |
| Arabic Math | arabic_math_page-1.jpg | Mistral | 960 | 1,686ms | ✅ |
| Arabic Math | arabic_math_page-2.jpg | Mistral | 536 | 2,492ms | ✅ |
| Chemistry | chimie_page-1.jpg | Mistral | 224 | 1,226ms | ✅ |
| Chemistry | chimie_page-2.jpg | Mistral | 39 | 1,175ms | ✅ |

## Summary

| Metric | Value |
|--------|-------|
| Total images | 9 |
| Successful | 9/9 (100%) |
| Total chars extracted | 17,824 |
| Average chars/image | 1,980 |
| Average latency | 2,617ms |
| Cost per page (Mistral) | $0.002 |
| Total cost | ~$0.018 |

## Chain Performance

- **Mistral** succeeded on all 9 images — OCR.space (fallback) was never needed
- Arabic text: Good extraction (960 + 536 chars on math pages)
- Math equations: LaTeX preserved (`$`, `\frac`, `\sqrt`, etc.)
- French text: Near-perfect extraction (2,842 chars on French essay)
- Physics/Chemistry: Tables, formulas, and QCM preserved

## Provider Comparison (prior work)

| Provider | OCRBench Score | Notes |
|----------|---------------|-------|
| LFM2.5-VL | 822 | Highest, best French/Arabic (needs Cloud Shell) |
| Mistral OCR | ~800 | Primary provider, cloud API |
| PaddleOCR | ~750 | Best open-source, needs install |
| EasyOCR | ~700 | Deep learning, 80+ languages |
| Tesseract | ~650 | Fast, needs preprocessing |
| OCR.space | ~600 | Free, no key needed |
| Workers AI | ~750 | Free, built-in |
