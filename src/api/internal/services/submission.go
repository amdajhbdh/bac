package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type SubmissionService struct {
	uploadDir string
	noonPath  string
}

type SubmissionResponse struct {
	JobID      string `json:"job_id"`
	Status     string `json:"status"`
	OCRText    string `json:"ocr_text,omitempty"`
	QuestionID string `json:"question_id,omitempty"`
}

func NewSubmissionService() *SubmissionService {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	os.MkdirAll(uploadDir, 0755)

	return &SubmissionService{
		uploadDir: uploadDir,
		noonPath:  "../../src/noon",
	}
}

func (s *SubmissionService) SubmitImage(fileData []byte, filename string, subject string) (*SubmissionResponse, error) {
	jobID := uuid.New().String()

	// Save uploaded file
	filePath := filepath.Join(s.uploadDir, jobID+"_"+filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, err
	}

	// Start async OCR processing
	go s.processImageOCR(jobID, filePath, subject)

	return &SubmissionResponse{
		JobID:  jobID,
		Status: "processing",
	}, nil
}

func (s *SubmissionService) SubmitPDF(fileData []byte, filename string, subject string) (*SubmissionResponse, error) {
	jobID := uuid.New().String()

	// Save uploaded file
	filePath := filepath.Join(s.uploadDir, jobID+"_"+filename)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return nil, err
	}

	// Start async PDF processing
	go s.processPDF(jobID, filePath, subject)

	return &SubmissionResponse{
		JobID:  jobID,
		Status: "processing",
	}, nil
}

func (s *SubmissionService) SubmitURL(url string, subject string) (*SubmissionResponse, error) {
	jobID := uuid.New().String()

	// Start async URL scraping
	go s.processURL(jobID, url, subject)

	return &SubmissionResponse{
		JobID:  jobID,
		Status: "processing",
	}, nil
}

func (s *SubmissionService) processImageOCR(jobID, filePath, subject string) {
	// Try multiple OCR methods
	var ocrText string

	// Method 1: Tesseract (primary)
	ocrText = s.runTesseract(filePath)

	if ocrText == "" {
		// Method 2: EasyOCR fallback
		ocrText = s.runEasyOCR(filePath)
	}

	// Save OCR result
	s.saveJobResult(jobID, ocrText, "completed")
}

func (s *SubmissionService) runTesseract(imagePath string) string {
	cmd := exec.Command("tesseract", imagePath, "stdout", "-l", "ara+fra+eng")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (s *SubmissionService) runEasyOCR(imagePath string) string {
	// Would call Python EasyOCR service
	return ""
}

func (s *SubmissionService) processPDF(jobID, filePath, subject string) {
	// Extract text from PDF
	text := s.extractPDFText(filePath)

	// Extract images if any
	_ = s.extractPDFImages(filePath)

	// Save result
	s.saveJobResult(jobID, text, "completed")
}

func (s *SubmissionService) extractPDFText(pdfPath string) string {
	cmd := exec.Command("pdftotext", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (s *SubmissionService) extractPDFImages(pdfPath string) []string {
	// Use pdfimages to extract images
	return []string{}
}

func (s *SubmissionService) processURL(jobID, url, subject string) {
	// Scrape content from URL
	text := s.scrapeURL(url)

	// Save result
	s.saveJobResult(jobID, text, "completed")
}

func (s *SubmissionService) scrapeURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (s *SubmissionService) saveJobResult(jobID, text, status string) {
	// Save to database/cache
	resultPath := filepath.Join(s.uploadDir, jobID+".txt")
	os.WriteFile(resultPath, []byte(text), 0644)
}

func (s *SubmissionService) GetStatus(jobID string) (*SubmissionResponse, error) {
	resultPath := filepath.Join(s.uploadDir, jobID+".txt")
	data, err := os.ReadFile(resultPath)

	status := "processing"
	var ocrText string

	if err == nil {
		status = "completed"
		ocrText = string(data)
	}

	return &SubmissionResponse{
		JobID:   jobID,
		Status:  status,
		OCRText: ocrText,
	}, nil
}

func (s *SubmissionService) AnalyzeWithAI(text, subject string) (map[string]interface{}, error) {
	// Send to AI for analysis - prompt would be used in production
	_ = fmt.Sprintf(`Analyze this question and extract:
1. Subject
2. Chapter
3. Topic tags
4. Question type
5. Difficulty estimate (1-5)
6. Key concepts

Text: %s

Format as JSON.`, text)

	// Would call LLM service
	return map[string]interface{}{
		"subject":       subject,
		"question_type": "calculation",
		"difficulty":    3,
		"concepts":      []string{},
	}, nil
}

func (s *SubmissionService) generateNoonAnimationPrompt(solution string) string {
	// Generate noon code from solution
	return ""
}
