package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubmissionService struct {
	uploadDir string
	noonPath  string
	db        *pgxpool.Pool
}

func NewSubmissionService(db *pgxpool.Pool) *SubmissionService {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	os.MkdirAll(uploadDir, 0755)

	return &SubmissionService{
		uploadDir: uploadDir,
		noonPath:  "../../src/noon",
		db:        db,
	}
}

type SubmissionResponse struct {
	JobID      string `json:"job_id"`
	Status     string `json:"status"`
	OCRText    string `json:"ocr_text,omitempty"`
	QuestionID string `json:"question_id,omitempty"`
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
	slog.Info("starting OCR processing", "job_id", jobID, "file", filePath)

	// Try multiple OCR methods in order of preference
	var ocrText string
	var method string

	// Method 1: Tesseract (primary - local, free)
	ocrText = s.runTesseract(filePath)
	if ocrText != "" && len(ocrText) > 10 {
		method = "tesseract"
		slog.Info("OCR succeeded with tesseract", "job_id", jobID, "length", len(ocrText))
	} else {
		// Method 2: Windows OCR (Windows only)
		ocrText = s.runWindowsOCR(filePath)
		if ocrText != "" && len(ocrText) > 10 {
			method = "windows-ocr"
			slog.Info("OCR succeeded with Windows OCR", "job_id", jobID, "length", len(ocrText))
		} else {
			// Method 3: Google Cloud Vision API (requires credentials)
			ocrText = s.runGoogleVisionOCR(filePath)
			if ocrText != "" && len(ocrText) > 10 {
				method = "google-vision"
				slog.Info("OCR succeeded with Google Vision", "job_id", jobID, "length", len(ocrText))
			} else {
				// Method 4: OCR.space API (free tier)
				ocrText = s.runOCRSpaceOCR(filePath)
				if ocrText != "" && len(ocrText) > 10 {
					method = "ocr-space"
					slog.Info("OCR succeeded with OCR.space", "job_id", jobID, "length", len(ocrText))
				}
			}
		}
	}

	// Save OCR result
	s.saveJobResult(jobID, ocrText, method, "completed")

	// If database is available, save the question
	if s.db != nil && ocrText != "" {
		s.saveQuestionToDatabase(jobID, ocrText, subject, filePath)
	}
}

func (s *SubmissionService) runTesseract(imagePath string) string {
	// Preprocess image for better OCR
	preprocessedPath := s.preprocessImage(imagePath)
	if preprocessedPath != imagePath {
		defer os.Remove(preprocessedPath)
		imagePath = preprocessedPath
	}

	cmd := exec.Command("tesseract", imagePath, "stdout", "-l", "ara+fra+eng", "--psm", "6")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("tesseract failed", "error", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (s *SubmissionService) preprocessImage(imagePath string) string {
	// Use ImageMagick to preprocess for better OCR
	preprocessedPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + "_processed.png"

	// Try to improve contrast and binarize
	cmd := exec.Command("convert", imagePath, "-contrast", "-normalize", "-threshold", "50%", preprocessedPath)
	if err := cmd.Run(); err != nil {
		// If ImageMagick not available, return original
		return imagePath
	}

	return preprocessedPath
}

func (s *SubmissionService) runWindowsOCR(imagePath string) string {
	// Use PowerShell Windows.Media.Ocr
	script := fmt.Sprintf(`
Add-Type -AssemblyName System.Runtime.WindowsRuntime
$asyncOp = [Windows.Media.Ocr.OcrEngine]::TryCreateFromUserProfileLanguages()
if ($asyncOp) {
    $ocrEngine = $asyncOp.GetResults()
    $file = [Windows.Storage.StorageFile]::GetFileFromPathAsync('%s').GetAwaiter().GetResult()
    $bitmap = [Windows.Graphics.Imaging.BitmapDecoder]::CreateAsync($file.OpenAsync([Windows.Storage.FileAccessMode]::Read)).GetAwaiter().GetResult()
    $softwareBitmap = $bitmap.GetPixelDataAsync().GetAwaiter().GetResult()
    $result = $ocrEngine.RecognizeAsync($softwareBitmap).GetAwaiter().GetResult()
    Write-Output $result.Text
}
`, strings.ReplaceAll(imagePath, "\\", "\\\\"))

	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("Windows OCR failed", "error", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (s *SubmissionService) runGoogleVisionOCR(imagePath string) string {
	apiKey := os.Getenv("GOOGLE_VISION_API_KEY")
	if apiKey == "" {
		slog.Debug("Google Vision API key not set")
		return ""
	}

	// Read and encode image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString(imageData)

	// Build request
	type VisionRequest struct {
		Requests []struct {
			Image struct {
				Content string `json:"content"`
			} `json:"image"`
			Features []struct {
				Type       string `json:"type"`
				MaxResults int    `json:"maxResults"`
			} `json:"features"`
		} `json:"requests"`
	}

	var visionReq VisionRequest
	visionReq.Requests[0].Image.Content = encoded
	visionReq.Requests[0].Features = []struct {
		Type       string `json:"type"`
		MaxResults int    `json:"maxResults"`
	}{{Type: "TEXT_DETECTION", MaxResults: 10}}

	reqBody, _ := json.Marshal(visionReq)
	resp, err := http.Post(
		fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%s", apiKey),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	type VisionResponse struct {
		Responses []struct {
			TextAnnotations []struct {
				Description string `json:"description"`
			} `json:"textAnnotations"`
		} `json:"responses"`
	}

	var visionResp VisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&visionResp); err != nil {
		return ""
	}

	if len(visionResp.Responses) > 0 && len(visionResp.Responses[0].TextAnnotations) > 0 {
		return visionResp.Responses[0].TextAnnotations[0].Description
	}

	return ""
}

func (s *SubmissionService) runOCRSpaceOCR(imagePath string) string {
	apiKey := os.Getenv("OCRSPACE_API_KEY")
	if apiKey == "" {
		apiKey = "helloworld" // Free tier key
	}

	// Read image
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return ""
	}

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("language", "ara+fra+eng")
	writer.WriteField("isOverlayRequired", "false")
	writer.WriteField("filetype", "PNG")
	writer.WriteField("detectOrientation", "true")
	writer.WriteField("scale", "true")
	writer.WriteField("OCREngine", "2")

	part, _ := writer.CreateFormFile("file", filepath.Base(imagePath))
	part.Write(imageData)
	writer.Close()

	resp, err := http.Post(
		fmt.Sprintf("https://api.ocr.space/parse/image?apikey=%s", apiKey),
		writer.FormDataContentType(),
		&buf,
	)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	type OCRSpaceResponse struct {
		ParsedResults []struct {
			ParsedText string `json:"ParsedText"`
		} `json:"ParsedResults"`
	}

	var ocrResp OCRSpaceResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		return ""
	}

	if len(ocrResp.ParsedResults) > 0 {
		return ocrResp.ParsedResults[0].ParsedText
	}

	return ""
}

func (s *SubmissionService) runEasyOCR(imagePath string) string {
	// Check if Python EasyOCR is available
	cmd := exec.Command("python3", "-c", "import easyocr")
	if err := cmd.Run(); err != nil {
		slog.Debug("EasyOCR not available")
		return ""
	}

	// Run EasyOCR
	script := fmt.Sprintf(`
import easyocr
reader = easyocr.Reader(['fr', 'ar', 'en'])
result = reader.readtext('%s', detail=0)
print('\\n'.join(result))
`, imagePath)

	cmd = exec.Command("python3", "-c", script)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (s *SubmissionService) processPDF(jobID, filePath, subject string) {
	slog.Info("starting PDF processing", "job_id", jobID, "file", filePath)

	// Extract text from PDF
	text := s.extractPDFText(filePath)
	if text == "" {
		// Try OCR on PDF images
		text = s.extractPDFImagesOCR(filePath)
	}

	// Extract images if any
	images := s.extractPDFImages(filePath)
	for _, img := range images {
		imgText := s.runTesseract(img)
		text += "\n" + imgText
	}

	// Determine method
	method := "pdftotext"
	if text == "" {
		method = "ocr"
	}

	// Save result
	s.saveJobResult(jobID, text, method, "completed")

	// If database is available, save the question
	if s.db != nil && text != "" {
		s.saveQuestionToDatabase(jobID, text, subject, filePath)
	}
}

func (s *SubmissionService) extractPDFText(pdfPath string) string {
	cmd := exec.Command("pdftotext", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("pdftotext failed", "error", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (s *SubmissionService) extractPDFImages(pdfPath string) []string {
	// Create temp directory for images
	tmpDir := filepath.Join(os.TempDir(), "pdfimages-"+uuid.New().String())
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	// Use pdfimages to extract images
	cmd := exec.Command("pdfimages", "-list", pdfPath)
	_, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	// Extract images
	cmd = exec.Command("pdfimages", "-png", pdfPath, filepath.Join(tmpDir, "img"))
	if err := cmd.Run(); err != nil {
		return []string{}
	}

	// Get extracted images
	var images []string
	files, _ := os.ReadDir(tmpDir)
	for _, f := range files {
		if !f.IsDir() {
			images = append(images, filepath.Join(tmpDir, f.Name()))
		}
	}

	return images
}

func (s *SubmissionService) extractPDFImagesOCR(pdfPath string) string {
	images := s.extractPDFImages(pdfPath)
	var allText string

	for _, img := range images {
		text := s.runTesseract(img)
		allText += text + "\n"
	}

	return strings.TrimSpace(allText)
}

func (s *SubmissionService) processURL(jobID, url, subject string) {
	slog.Info("starting URL processing", "job_id", jobID, "url", url)

	// Scrape content from URL
	text := s.scrapeURL(url)

	// Determine method
	method := "curl"
	if text == "" {
		method = "failed"
	}

	// Clean and extract relevant content
	cleanedText := s.cleanWebContent(text)

	// Save result
	s.saveJobResult(jobID, cleanedText, method, "completed")

	// If database is available, save the question
	if s.db != nil && cleanedText != "" {
		s.saveQuestionToDatabase(jobID, cleanedText, subject, url)
	}
}

func (s *SubmissionService) scrapeURL(url string) string {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("URL fetch failed", "error", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func (s *SubmissionService) cleanWebContent(html string) string {
	// Simple HTML tag removal
	text := strings.ReplaceAll(html, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<p>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")

	// Remove script and style content
	text = strings.ReplaceAll(text, "<script[^>]*>[^<]*</script>", "")
	text = strings.ReplaceAll(text, "<style[^>]*>[^<]*</style>", "")

	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]+>`)
	text = re.ReplaceAllString(text, "")

	// Clean up whitespace
	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}

func (s *SubmissionService) saveJobResult(jobID, text, method, status string) {
	// Save to file system
	resultPath := filepath.Join(s.uploadDir, jobID+".json")

	type Result struct {
		JobID   string `json:"job_id"`
		Text    string `json:"text"`
		Method  string `json:"method"`
		Status  string `json:"status"`
		Updated string `json:"updated_at"`
	}

	result := Result{
		JobID:   jobID,
		Text:    text,
		Method:  method,
		Status:  status,
		Updated: time.Now().Format(time.RFC3339),
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	os.WriteFile(resultPath, data, 0644)
}

func (s *SubmissionService) saveQuestionToDatabase(jobID, text, sourceType, sourcePath string) {
	ctx := context.Background()

	// Insert into questions table
	_, err := s.db.Exec(ctx, `
		INSERT INTO questions (
			question_text, 
			source_type, 
			original_filename,
			verification_status,
			created_at
		) VALUES ($1, $2, $3, 'pending', NOW())`,
		text, sourceType, sourcePath)

	if err != nil {
		slog.Error("failed to save question to database", "error", err, "job_id", jobID)
		return
	}

	slog.Info("question saved to database", "job_id", jobID, "text_length", len(text))
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
