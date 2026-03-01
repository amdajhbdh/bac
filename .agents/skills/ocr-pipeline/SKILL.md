# Skill: OCR Pipeline (Tesseract + Surya)

## Purpose

OCR pipeline using Tesseract and Surya for document extraction in BAC Unified.

## When to use

- Extracting text from uploaded images
- Processing PDF documents
- Converting handwritten notes to text

## Priority Order

1. **Surya** - OCR (primary, better for math)
2. **Tesseract** - OCR (fallback)

## NEVER Use

- ❌ Google Vision API
- ❌ AWS Textract
- ❌ Azure Computer Vision

## Surya (Preferred)

Python-based OCR with layout analysis.

```python
from surya.ocr import run_ocr
from surya.model.detection.segformer import load_model as load_det_model
from surya.model.recognition.model import load_model as load_rec_model
from surya.model.detection.config import get_det_config
from surya.model.recognition.config import get_rec_config

# Load models
det_model = load_det_config()
rec_model = load_rec_config()

# Run OCR
predictions = run_ocr(image, [lang], det_model, rec_model)

for pred in predictions:
    print(pred.text)
```

## Tesseract (Fallback)

```go
import "github.com/tesseract-org/tesseract"

func OCR(imagePath string) (string, error) {
    client, _ := tesseract.NewClient("")
    defer client.Close()
    
    text, err := client.Text(imagePath)
    return text, err
}
```

## Installation

```bash
# Tesseract
# macOS
brew install tesseract tesseract-lang

# Linux
sudo apt install tesseract tesseract-ocr-fra

# Surya
pip install surya-ocr
```

## Language Support

| Language | Code | Use |
|----------|------|-----|
| French | fra | Primary |
| Arabic | ara | Secondary |
| English | eng | Fallback |

## Pipeline Flow

```
Image/PDF → Preprocess → Surya OCR → Text Output
                         ↓
                  Tesseract (fallback)
```

## Preprocessing

```python
import cv2

def preprocess(image_path):
    img = cv2.imread(image_path)
    
    # Grayscale
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    # Denoise
    denoised = cv2.fastNlMeansDenoising(gray)
    
    # Threshold
    _, thresh = cv2.threshold(denoised, 0, 255, 
        cv2.THRESH_BINARY + cv2.THRESH_OTSU)
    
    return thresh
```

## Environment

| Variable | Purpose |
|----------|---------|
| `TESSDATA_PREFIX` | Tesseract data path |
| `OCR_LANG` | Default language (fra) |

## Go Integration

```go
// From src/agent/internal/ocr/ocr.go
func ProcessImage(ctx context.Context, imagePath string) (string, error) {
    // Try Surya first
    text, err := runSuryaOCR(imagePath)
    if err == nil && text != "" {
        return text, nil
    }
    
    // Fallback to Tesseract
    return runTesseractOCR(imagePath)
}
```

## Anti-Patterns

- ❌ Using paid OCR APIs
- ❌ Not handling image preprocessing
- ❌ Not supporting multiple languages
- ❌ Ignoring OCR confidence scores
