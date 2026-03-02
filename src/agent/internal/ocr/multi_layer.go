package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type MultiLayerResult struct {
	Text         string
	Confidence   float64
	Source       string
	AllResults   []OCRResult
	QualityFlags []string
	NLMAnalysis  *NLMAnalysis
	Complete     bool
}

type NLMAnalysis struct {
	Subject    string
	Concepts   []string
	Accuracy   float64
	Missing    []string
	Questions  []string
	NotebookID string
}

func MultiLayerExtract(ctx context.Context, filePath string) (*MultiLayerResult, error) {
	slog.Info("starting multi-layer extraction", "file", filePath)

	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".pdf" {
		return extractPDFMultiLayer(ctx, filePath)
	}
	return extractImageMultiLayer(ctx, filePath)
}

func extractImageMultiLayer(ctx context.Context, imagePath string) (*MultiLayerResult, error) {
	slog.Info("running multi-layer image extraction", "path", imagePath)

	results := make(chan OCRResult, 4)
	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		results <- processTesseract(ctx, imagePath)
	}()

	go func() {
		defer wg.Done()
		results <- processOllamaVision(ctx, imagePath)
	}()

	go func() {
		defer wg.Done()
		results <- processPythonOCR(ctx, imagePath)
	}()

	go func() {
		defer wg.Done()
		results <- processGoogleLens(ctx, imagePath)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []OCRResult
	for result := range results {
		if result.Success {
			slog.Debug("OCR result", "source", result.Source, "confidence", result.Confidence, "text_len", len(result.Text))
			allResults = append(allResults, result)
		} else {
			slog.Debug("OCR failed", "source", result.Source)
		}
	}

	if len(allResults) == 0 {
		return &MultiLayerResult{
			Complete: false,
			Text:     "",
		}, fmt.Errorf("all OCR methods failed")
	}

	bestWithNLM, err := verifyWithNLM(ctx, allResults, imagePath)
	if err != nil {
		slog.Warn("NLM verification failed, using fallback", "error", err)
		return fallbackToBest(allResults), nil
	}

	return bestWithNLM, nil
}

func extractPDFMultiLayer(ctx context.Context, pdfPath string) (*MultiLayerResult, error) {
	slog.Info("running multi-layer PDF extraction", "path", pdfPath)

	results := make(chan OCRResult, 3)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		results <- processTesseractPDF(ctx, pdfPath)
	}()

	go func() {
		defer wg.Done()
		results <- processOllamaVisionPDF(ctx, pdfPath)
	}()

	go func() {
		defer wg.Done()
		results <- processPythonPDF(ctx, pdfPath)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []OCRResult
	for result := range results {
		if result.Success {
			slog.Debug("PDF OCR result", "source", result.Source, "confidence", result.Confidence)
			allResults = append(allResults, result)
		}
	}

	if len(allResults) == 0 {
		return &MultiLayerResult{
			Complete: false,
			Text:     "",
		}, fmt.Errorf("all PDF OCR methods failed")
	}

	bestWithNLM, err := verifyWithNLM(ctx, allResults, pdfPath)
	if err != nil {
		slog.Warn("NLM verification failed, using fallback", "error", err)
		return fallbackToBest(allResults), nil
	}

	return bestWithNLM, nil
}

func processTesseractPDF(ctx context.Context, pdfPath string) OCRResult {
	tmpDir := filepath.Join(os.TempDir(), "bac-ocr")
	os.MkdirAll(tmpDir, 0755)

	cmd := exec.CommandContext(ctx, "convert", "-density", "300", pdfPath+"[0]", filepath.Join(tmpDir, "page.png"))
	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Debug("ImageMagick convert failed", "error", err)
		return OCRResult{Source: "tesseract-pdf", Success: false}
	}

	imgPath := filepath.Join(tmpDir, "page.png")
	defer os.RemoveAll(tmpDir)

	return processTesseract(ctx, imgPath)
}

func processOllamaVisionPDF(ctx context.Context, pdfPath string) OCRResult {
	tmpDir := filepath.Join(os.TempDir(), "bac-ocr")
	os.MkdirAll(tmpDir, 0755)

	cmd := exec.CommandContext(ctx, "convert", "-density", "300", pdfPath+"[0]", filepath.Join(tmpDir, "page.png"))
	_, err := cmd.CombinedOutput()
	if err != nil {
		return OCRResult{Source: "ollama-pdf", Success: false}
	}

	imgPath := filepath.Join(tmpDir, "page.png")
	defer os.RemoveAll(tmpDir)

	return processOllamaVision(ctx, imgPath)
}

func processPythonPDF(ctx context.Context, pdfPath string) OCRResult {
	cmd := exec.CommandContext(ctx, "python3", "-c", fmt.Sprintf(`
import fitz
doc = fitz.open('%s')
text = ""
for page in doc:
    text += page.get_text() + "\\n---PAGE---\\n"
print(text)
`, pdfPath))

	output, err := cmd.Output()
	if err != nil {
		return OCRResult{Source: "python-pdf", Success: false}
	}

	text := strings.TrimSpace(string(output))
	if len(text) < 10 {
		return OCRResult{Source: "python-pdf", Success: false}
	}

	return OCRResult{
		Text:       text,
		Source:     "python-pdf",
		Confidence: 0.8,
		Success:    true,
	}
}

func verifyWithNLM(ctx context.Context, results []OCRResult, filePath string) (*MultiLayerResult, error) {
	slog.Info("verifying OCR results with NLM", "count", len(results))

	if len(results) == 0 {
		return nil, fmt.Errorf("no results to verify")
	}

	notebookID, err := createVerificationNotebook(ctx)
	if err != nil {
		slog.Warn("failed to create verification notebook", "error", err)
		return nil, err
	}

	prompt := buildVerificationPrompt(results)

	verificationResult := queryNLMForVerification(ctx, notebookID, prompt)
	if !verificationResult.Success {
		return nil, fmt.Errorf("NLM verification failed: %s", verificationResult.Error)
	}

	selected, confidence, qualityFlags := parseNLMVerification(verificationResult.RawOutput, results)

	slog.Info("NLM verification complete",
		"selected_source", selected.Source,
		"confidence", confidence,
		"quality_flags", qualityFlags)

	subject, concepts := extractSubjectAndConcepts(verificationResult.RawOutput)

	return &MultiLayerResult{
		Text:         selected.Text,
		Confidence:   confidence,
		Source:       "nlm-verified:" + selected.Source,
		AllResults:   results,
		QualityFlags: qualityFlags,
		NLMAnalysis: &NLMAnalysis{
			Subject:    subject,
			Concepts:   concepts,
			Accuracy:   confidence,
			NotebookID: notebookID,
		},
		Complete: true,
	}, nil
}

func createVerificationNotebook(ctx context.Context) (string, error) {
	slog.Info("creating verification notebook")

	cmd := exec.CommandContext(ctx, "nlm", "notebook", "create", "OCR-Verification-"+fmt.Sprint(time.Now().Unix()))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("create notebook: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(string(output)), ":")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[1]), nil
	}

	return strings.TrimSpace(string(output)), nil
}

func buildVerificationPrompt(results []OCRResult) string {
	var sb strings.Builder
	sb.WriteString("Analyze the following OCR extraction results and determine which is the most accurate and complete.\n\n")

	for i, r := range results {
		sb.WriteString(fmt.Sprintf("--- RESULT %d (Source: %s, Confidence: %.2f) ---\n", i+1, r.Source, r.Confidence))
		sb.WriteString(r.Text)
		sb.WriteString("\n\n")
	}

	sb.WriteString(`Now provide your analysis in the following JSON format:
{
  "best_index": <index of best result 1-N>,
  "reasoning": "<brief explanation of why this is best>",
  "quality_assessment": {
    "completeness": "<complete/incomplete/partial>",
    "accuracy": "<high/medium/low>",
    "readability": "<good/medium/poor>"
  },
  "suggested_improvements": ["<any issues found>"],
  "detected_subject": "<math/pc/svt/philosophy/other>",
  "key_concepts": ["<list of concepts detected>"]
}`)

	return sb.String()
}

func queryNLMForVerification(ctx context.Context, notebookID, prompt string) VerificationResult {
	slog.Debug("querying NLM for verification")

	cmd := exec.CommandContext(ctx, "nlm", "notebook", "query", notebookID, prompt)
	output, err := cmd.Output()
	if err != nil {
		return VerificationResult{Success: false, Error: err.Error()}
	}

	var vr VerificationResult
	if err := json.Unmarshal(output, &vr); err != nil {
		return VerificationResult{
			Success:   true,
			RawOutput: string(output),
		}
	}

	vr.Success = true
	return vr
}

type VerificationResult struct {
	Success           bool     `json:"success"`
	Error             string   `json:"error,omitempty"`
	RawOutput         string   `json:"raw_output,omitempty"`
	BestIndex         int      `json:"best_index"`
	Reasoning         string   `json:"reasoning"`
	QualityAssessment Quality  `json:"quality_assessment"`
	Subject           string   `json:"detected_subject"`
	Concepts          []string `json:"key_concepts"`
}

type Quality struct {
	Completeness string `json:"completeness"`
	Accuracy     string `json:"accuracy"`
	Readability  string `json:"readability"`
}

func parseNLMVerification(nlmOutput string, results []OCRResult) (OCRResult, float64, []string) {
	var vr VerificationResult
	if err := json.Unmarshal([]byte(nlmOutput), &vr); err != nil {
		slog.Warn("failed to parse NLM verification JSON", "error", err)
		return fallbackToBest(results).AllResults[0], 0.5, []string{"unverified"}
	}

	if vr.BestIndex < 1 || vr.BestIndex > len(results) {
		slog.Warn("invalid best index from NLM", "index", vr.BestIndex)
		return fallbackToBest(results).AllResults[0], 0.5, []string{"unverified"}
	}

	selected := results[vr.BestIndex-1]

	confidence := calculateNLMConfidence(vr)

	var flags []string
	if vr.QualityAssessment.Completeness != "" {
		flags = append(flags, vr.QualityAssessment.Completeness)
	}
	if vr.QualityAssessment.Accuracy != "" {
		flags = append(flags, vr.QualityAssessment.Accuracy)
	}
	if vr.QualityAssessment.Readability != "" {
		flags = append(flags, vr.QualityAssessment.Readability)
	}

	return selected, confidence, flags
}

func calculateNLMConfidence(vr VerificationResult) float64 {
	conf := 0.5

	switch vr.QualityAssessment.Accuracy {
	case "high":
		conf += 0.3
	case "medium":
		conf += 0.15
	case "low":
		conf += 0.0
	}

	switch vr.QualityAssessment.Completeness {
	case "complete":
		conf += 0.2
	case "partial":
		conf += 0.1
	case "incomplete":
		conf += 0.0
	}

	if conf > 1.0 {
		conf = 1.0
	}

	return conf
}

func extractSubjectAndConcepts(nlmOutput string) (string, []string) {
	var vr VerificationResult
	if err := json.Unmarshal([]byte(nlmOutput), &vr); err != nil {
		return "unknown", []string{}
	}

	subject := vr.Subject
	if subject == "" {
		subject = detectSubjectFromText(nlmOutput)
	}

	concepts := vr.Concepts
	if len(concepts) == 0 {
		concepts = extractConceptsFromText(nlmOutput)
	}

	return subject, concepts
}

func detectSubjectFromText(text string) string {
	lower := strings.ToLower(text)

	subjects := map[string][]string{
		"math":        {"equation", "function", "derivative", "integral", "matrix", "vector", "calcul", "algebra", "limite"},
		"pc":          {"physique", "chimie", "force", "energie", "mouvement", "thermodynamique"},
		"svt":         {"biologie", "cellule", "ADN", "ecosysteme", "photosynthese", "genetique"},
		"philosophie": {"ethique", "morale", "justice", "verite", "conscience", "existence"},
	}

	for subject, keywords := range subjects {
		count := 0
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				count++
			}
		}
		if count >= 2 {
			return subject
		}
	}

	return "unknown"
}

func extractConceptsFromText(text string) []string {
	keywords := []string{
		"equation", "function", "derivative", "integral", "matrix", "vector",
		"physique", "chimie", "biologie", "cellule", "ADN",
		"force", "energie", "mouvement", "thermodynamique",
		"ecosysteme", "photosynthese", "genetique",
	}

	var found []string
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			found = append(found, kw)
		}
	}

	if len(found) == 0 {
		found = []string{"general"}
	}

	return found
}

func fallbackToBest(results []OCRResult) *MultiLayerResult {
	var best OCRResult
	best.Confidence = 0

	for _, r := range results {
		if r.Success && r.Confidence > best.Confidence {
			best = r
		}
	}

	flags := []string{"fallback"}
	if best.Confidence > 0.8 {
		flags = []string{"high-confidence"}
	} else if best.Confidence > 0.5 {
		flags = []string{"medium-confidence"}
	} else {
		flags = []string{"low-confidence"}
	}

	return &MultiLayerResult{
		Text:         best.Text,
		Confidence:   best.Confidence,
		Source:       "fallback:" + best.Source,
		AllResults:   results,
		QualityFlags: flags,
		Complete:     best.Success,
	}
}
