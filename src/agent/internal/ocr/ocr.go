package ocr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 120 * time.Second,
}

type OCRResult struct {
	Text       string
	Source     string
	Confidence float64
	Success    bool
	Error      string
}

func ProcessImage(ctx context.Context, imagePath string) OCRResult {
	slog.Info("processing image", "path", imagePath)

	// Try multiple OCR methods and return best result
	results := make(chan OCRResult, 4)

	// 1. Try Tesseract
	go func() {
		results <- processTesseract(ctx, imagePath)
	}()

	// 2. Try Ollama vision
	go func() {
		results <- processOllamaVision(ctx, imagePath)
	}()

	// 3. Try Python OCR as fallback
	go func() {
		results <- processPythonOCR(ctx, imagePath)
	}()

	// 4. Try Google Lens / Cloud Vision
	go func() {
		results <- processGoogleLens(ctx, imagePath)
	}()

	// Collect best result
	var best OCRResult
	best.Confidence = 0

	for i := 0; i < 3; i++ {
		result := <-results
		slog.Debug("OCR attempt", "source", result.Source, "success", result.Success, "confidence", result.Confidence)
		if result.Success && result.Confidence > best.Confidence {
			best = result
		}
	}

	if best.Success {
		slog.Info("OCR successful", "source", best.Source, "confidence", best.Confidence)
	} else {
		slog.Error("all OCR methods failed")
		best = OCRResult{Success: false, Error: "all OCR methods failed"}
	}

	return best
}

func processTesseract(ctx context.Context, imagePath string) OCRResult {
	cmd := exec.CommandContext(ctx, "tesseract", imagePath, "stdout", "-l", "fra+eng")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Debug("Tesseract failed", "error", err)
		return OCRResult{Source: "tesseract", Success: false}
	}

	text := strings.TrimSpace(stdout.String())
	if len(text) < 5 {
		return OCRResult{Source: "tesseract", Success: false}
	}

	slog.Info("Tesseract succeeded")
	return OCRResult{
		Text:       text,
		Source:     "tesseract",
		Confidence: 0.7,
		Success:    true,
	}
}

func processOllamaVision(ctx context.Context, imagePath string) OCRResult {
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		return OCRResult{Source: "ollama-vision", Success: false}
	}
	b64 := base64.StdEncoding.EncodeToString(imgData)

	models := []string{"llava", "llama3.2-vision", "qwen2-vl"}

	for _, model := range models {
		reqBody, _ := json.Marshal(map[string]interface{}{
			"model":  model,
			"prompt": "Extract all text from this image in French.",
			"images": []string{b64},
		})

		resp, err := httpClient.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			continue
		}

		response, ok := result["response"].(string)
		if !ok || len(response) < 10 {
			continue
		}

		slog.Info("Ollama vision succeeded", "model", model)
		return OCRResult{
			Text:       response,
			Source:     "ollama:" + model,
			Confidence: 0.85,
			Success:    true,
		}
	}

	return OCRResult{Source: "ollama-vision", Success: false}
}

func processPythonOCR(ctx context.Context, imagePath string) OCRResult {
	// Try with Python and pytesseract
	cmd := exec.CommandContext(ctx, "python3", "-c", fmt.Sprintf(`
import pytesseract
from PIL import Image
img = Image.open('%s')
text = pytesseract.image_to_string(img, lang='fra+eng')
print(text)
`, imagePath))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return OCRResult{Source: "python-ocr", Success: false}
	}

	text := strings.TrimSpace(stdout.String())
	if len(text) < 5 {
		return OCRResult{Source: "python-ocr", Success: false}
	}

	return OCRResult{
		Text:       text,
		Source:     "python-ocr",
		Confidence: 0.75,
		Success:    true,
	}
}

func ProcessPDF(ctx context.Context, pdfPath string) []OCRResult {
	slog.Info("processing PDF", "path", pdfPath)

	// Try Python PDF processing first
	cmd := exec.CommandContext(ctx, "python3", "-c", fmt.Sprintf(`
import fitz
doc = fitz.open('%s')
for page in doc:
    print(page.get_text())
    print('---PAGE---')
`, pdfPath))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		slog.Warn("Python PDF failed, trying pdftotext", "error", err)
		return processPDFCLI(ctx, pdfPath)
	}

	text := strings.TrimSpace(stdout.String())
	if len(text) < 10 {
		return []OCRResult{{Source: "pdf", Success: false}}
	}

	// Split by pages
	pages := strings.Split(text, "---PAGE---")
	var results []OCRResult
	for _, page := range pages {
		page = strings.TrimSpace(page)
		if len(page) > 10 {
			results = append(results, OCRResult{
				Text:       page,
				Source:     "pdf-python",
				Confidence: 0.8,
				Success:    true,
			})
		}
	}

	return results
}

func processPDFCLI(ctx context.Context, pdfPath string) []OCRResult {
	cmd := exec.CommandContext(ctx, "pdftotext", pdfPath, "-")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return []OCRResult{{Source: "pdf", Success: false}}
	}

	text := strings.TrimSpace(stdout.String())
	if len(text) < 10 {
		return []OCRResult{}
	}

	return []OCRResult{{
		Text:       text,
		Source:     "pdf-cli",
		Confidence: 0.8,
		Success:    true,
	}}
}

func DownloadAndProcess(ctx context.Context, url string) OCRResult {
	slog.Info("downloading from URL", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return OCRResult{Source: "url", Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(os.TempDir(), "input_image")
	out, err := os.Create(tmpFile)
	if err != nil {
		return OCRResult{Source: "url", Success: false, Error: err.Error()}
	}
	defer out.Close()

	io.Copy(out, resp.Body)

	// Detect file type and process
	if strings.HasSuffix(url, ".pdf") || resp.Header.Get("Content-Type") == "application/pdf" {
		results := ProcessPDF(ctx, tmpFile)
		if len(results) > 0 && results[0].Success {
			return results[0]
		}
	}

	return ProcessImage(ctx, tmpFile)
}

func processGoogleLens(ctx context.Context, imagePath string) OCRResult {
	// Try Google Lens / Cloud Vision OCR
	lensClient := NewGoogleLensFromEnv()

	data, err := os.ReadFile(imagePath)
	if err != nil {
		return OCRResult{Source: "google-lens", Success: false, Error: err.Error()}
	}

	result, err := lensClient.OCR(ctx, data)
	if err != nil || result.Text == "" {
		return OCRResult{Source: "google-lens", Success: false}
	}

	return OCRResult{
		Text:       result.Text,
		Source:     "google-lens",
		Confidence: result.Confidence,
		Success:    true,
	}
}
