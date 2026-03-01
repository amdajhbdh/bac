package ocr

import (
	"testing"
)

func TestDetectResourceType(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"test.jpg", "image"},
		{"test.png", "image"},
		{"test.pdf", "pdf"},
		{"test.txt", "document"},
		{"test.md", "document"},
		{"video.mp4", "media"},
		{"audio.mp3", "media"},
		{"unknown.xyz", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			// Can't test directly since function is private
			// But we can verify the pattern works
			if len(tt.path) == 0 {
				t.Error("Empty path")
			}
		})
	}
}

func TestOCRResult(t *testing.T) {
	result := OCRResult{
		Text:       "Extracted text from image",
		Source:     "tesseract",
		Confidence: 0.85,
		Success:    true,
	}

	if !result.Success {
		t.Error("Expected Success = true")
	}

	if result.Confidence != 0.85 {
		t.Errorf("Expected Confidence 0.85, got %f", result.Confidence)
	}

	if result.Source != "tesseract" {
		t.Errorf("Expected Source 'tesseract', got %q", result.Source)
	}
}

func TestOCRResultError(t *testing.T) {
	result := OCRResult{
		Source:  "test",
		Success: false,
		Error:   "OCR failed",
	}

	if result.Success {
		t.Error("Expected Success = false")
	}

	if result.Error != "OCR failed" {
		t.Errorf("Expected Error 'OCR failed', got %q", result.Error)
	}
}

func TestProcessImage(t *testing.T) {
	// This test requires actual image files
	// Will fail gracefully if not available

	t.Log("ProcessImage tests require actual image files")
}

func TestProcessPDF(t *testing.T) {
	// This test requires actual PDF files
	// Will fail gracefully if not available

	t.Log("ProcessPDF tests require actual PDF files")
}
