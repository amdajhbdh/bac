package analyzer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bac-unified/agent/internal/nlm"
	"github.com/bac-unified/agent/internal/ocr"
	"github.com/bac-unified/agent/internal/online"
)

var httpClient = &http.Client{
	Timeout: 120 * time.Second,
}

type AnalysisResult struct {
	Service  string  `json:"service"`
	Content  string  `json:"content"`
	Summary  string  `json:"summary,omitempty"`
	Success  bool    `json:"success"`
	Duration float64 `json:"duration_ms"`
	Error    string  `json:"error,omitempty"`
}

type ResourceAnalysis struct {
	ResourcePath string           `json:"resource_path"`
	ResourceType string           `json:"resource_type"`
	Results      []AnalysisResult `json:"results"`
	OCRText      string           `json:"ocr_text,omitempty"`
	Timestamp    time.Time        `json:"timestamp"`
}

type Analyzer struct {
	services []string
}

func New(services []string) *Analyzer {
	if len(services) == 0 {
		services = []string{"ollama", "nlm", "chatgpt", "claude", "grok", "deepseek", "local"}
	}
	return &Analyzer{services: services}
}

func (a *Analyzer) Analyze(ctx context.Context, resourcePath string) ResourceAnalysis {
	start := time.Now()
	slog.Info("starting resource analysis", "path", resourcePath)

	resourceType := detectResourceType(resourcePath)
	slog.Info("detected resource type", "type", resourceType)

	// Step 1: OCR/PDF extraction
	var ocrText string
	if resourceType == "image" {
		result := ocr.ProcessImage(ctx, resourcePath)
		if result.Success {
			ocrText = result.Text
			slog.Info("OCR successful", "source", result.Source, "length", len(ocrText))
		}
	} else if resourceType == "pdf" {
		results := ocr.ProcessPDF(ctx, resourcePath)
		var texts []string
		for _, r := range results {
			if r.Success {
				texts = append(texts, r.Text)
			}
		}
		ocrText = strings.Join(texts, "\n\n---PAGE---\n\n")
		slog.Info("PDF extraction complete", "pages", len(results), "length", len(ocrText))
	}

	// Step 2: Analyze with multiple services in parallel
	results := a.analyzeParallel(ctx, ocrText, resourcePath)

	// Add OCR result if available
	if ocrText != "" {
		results = append(results, AnalysisResult{
			Service:  "ocr",
			Content:  ocrText,
			Success:  true,
			Duration: float64(time.Since(start).Milliseconds()),
		})
	}

	slog.Info("analysis complete", "services", len(results), "total_time", time.Since(start).Seconds())

	return ResourceAnalysis{
		ResourcePath: resourcePath,
		ResourceType: resourceType,
		Results:      results,
		OCRText:      ocrText,
		Timestamp:    time.Now(),
	}
}

func (a *Analyzer) analyzeParallel(ctx context.Context, text, resourcePath string) []AnalysisResult {
	results := make(chan AnalysisResult, len(a.services))
	var wg sync.WaitGroup

	for _, service := range a.services {
		wg.Add(1)
		go func(svc string) {
			defer wg.Done()
			result := a.analyzeWithService(ctx, svc, text, resourcePath)
			results <- result
		}(service)
	}

	// Close channel when done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var analysisResults []AnalysisResult
	for result := range results {
		if result.Success {
			slog.Info("service analysis success", "service", result.Service)
		} else {
			slog.Warn("service analysis failed", "service", result.Service, "error", result.Error)
		}
		analysisResults = append(analysisResults, result)
	}

	return analysisResults
}

func (a *Analyzer) analyzeWithService(ctx context.Context, service, text, resourcePath string) AnalysisResult {
	start := time.Now()

	var result AnalysisResult
	result.Service = service

	switch service {
	case "ollama":
		result = a.analyzeWithOllama(ctx, text)
	case "nlm":
		result = a.analyzeWithNLM(ctx, text, resourcePath)
	case "chatgpt":
		result = a.analyzeWithChatGPT(ctx, text)
	case "claude":
		result = a.analyzeWithClaude(ctx, text)
	case "grok":
		result = a.analyzeWithGrok(ctx, text)
	case "deepseek":
		result = a.analyzeWithDeepSeek(ctx, text)
	case "local":
		result = a.analyzeLocal(text)
	default:
		result.Success = false
		result.Error = "unknown service: " + service
	}

	result.Duration = float64(time.Since(start).Milliseconds())
	return result
}

func (a *Analyzer) analyzeWithOllama(ctx context.Context, text string) AnalysisResult {
	if text == "" {
		return AnalysisResult{Service: "ollama", Success: false, Error: "empty text"}
	}

	prompt := fmt.Sprintf(`Tu es un assistant d'analyse de documents BAC.
Analyse le contenu suivant et fournit:
1. Un résumé en français (2-3 phrases)
2. Les concepts clés
3. Les matières concernées (math, pc, svt, philosophie, français)

Contenu:
%s

Réponds en français.`, text[:min(len(text), 3000)])

	type GenerateRequest struct {
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		Stream      bool    `json:"stream"`
		Temperature float64 `json:"temperature"`
	}

	reqBody, _ := json.Marshal(GenerateRequest{
		Model:       "llama3.2:3b",
		Prompt:      prompt,
		Stream:      false,
		Temperature: 0.3,
	})

	resp, err := httpClient.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return AnalysisResult{Service: "ollama", Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return AnalysisResult{Service: "ollama", Success: false, Error: err.Error()}
	}

	response, ok := respData["response"].(string)
	if !ok {
		return AnalysisResult{Service: "ollama", Success: false, Error: "invalid response"}
	}

	// Extract summary
	summary := extractSummary(response)

	return AnalysisResult{
		Service: "ollama",
		Content: response,
		Summary: summary,
		Success: true,
	}
}

func (a *Analyzer) analyzeWithNLM(ctx context.Context, text, resourcePath string) AnalysisResult {
	if resourcePath == "" {
		return AnalysisResult{Service: "nlm", Success: false, Error: "no resource path"}
	}

	// Create notebook for analysis
	notebookTitle := "BAC Analysis - " + filepath.Base(resourcePath)
	notebookID, err := nlm.CreateNotebook(ctx, notebookTitle)
	if err != nil {
		slog.Warn("failed to create notebook", "error", err)
		// Try default notebook
		notebookID = "16b01950-5766-4353-8bed-c7f67966cb6b"
	}

	// Add source
	if strings.HasSuffix(resourcePath, ".pdf") || strings.HasSuffix(resourcePath, ".jpg") || strings.HasSuffix(resourcePath, ".png") {
		// For files, we need to upload to Drive first
		// For now, just query with the OCR text
		result := nlm.Query(ctx, notebookID, text[:min(len(text), 2000)])
		if result.Success {
			return AnalysisResult{
				Service: "nlm",
				Content: result.Results,
				Summary: extractSummary(result.Results),
				Success: true,
			}
		}
	}

	return AnalysisResult{Service: "nlm", Success: false, Error: "NLM analysis not available for this resource"}
}

func (a *Analyzer) analyzeWithChatGPT(ctx context.Context, text string) AnalysisResult {
	if text == "" {
		return AnalysisResult{Service: "chatgpt", Success: false, Error: "empty text"}
	}

	result := online.AutoSolve(ctx, text)
	return AnalysisResult{
		Service:  "chatgpt",
		Content:  result.Results,
		Summary:  extractSummary(result.Results),
		Success:  result.Success,
		Duration: 0,
		Error:    result.Error,
	}
}

func (a *Analyzer) analyzeWithClaude(ctx context.Context, text string) AnalysisResult {
	if text == "" {
		return AnalysisResult{Service: "claude", Success: false, Error: "empty text"}
	}

	result := online.AutoSolve(ctx, text)
	return AnalysisResult{
		Service:  "claude",
		Content:  result.Results,
		Summary:  extractSummary(result.Results),
		Success:  result.Success,
		Duration: 0,
		Error:    result.Error,
	}
}

func (a *Analyzer) analyzeWithGrok(ctx context.Context, text string) AnalysisResult {
	if text == "" {
		return AnalysisResult{Service: "grok", Success: false, Error: "empty text"}
	}

	result := online.AutoSolve(ctx, text)
	return AnalysisResult{
		Service:  "grok",
		Content:  result.Results,
		Summary:  extractSummary(result.Results),
		Success:  result.Success,
		Duration: 0,
		Error:    result.Error,
	}
}

func (a *Analyzer) analyzeWithDeepSeek(ctx context.Context, text string) AnalysisResult {
	if text == "" {
		return AnalysisResult{Service: "deepseek", Success: false, Error: "empty text"}
	}

	result := online.AutoSolve(ctx, text)
	return AnalysisResult{
		Service:  "deepseek",
		Content:  result.Results,
		Summary:  extractSummary(result.Results),
		Success:  result.Success,
		Duration: 0,
		Error:    result.Error,
	}
}

func (a *Analyzer) analyzeLocal(text string) AnalysisResult {
	if text == "" {
		return AnalysisResult{Service: "local", Success: false, Error: "empty text"}
	}

	// Simple local analysis without AI
	lines := strings.Split(text, "\n")
	wordCount := len(strings.Fields(text))

	summary := fmt.Sprintf("Analyse locale: %d lignes, %d mots, %d caractères", len(lines), wordCount, len(text))

	return AnalysisResult{
		Service: "local",
		Content: text[:min(len(text), 1000)],
		Summary: summary,
		Success: true,
	}
}

func detectResourceType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff":
		return "image"
	case ".pdf":
		return "pdf"
	case ".txt", ".md", ".doc", ".docx":
		return "document"
	case ".mp3", ".wav", ".mp4", ".avi":
		return "media"
	default:
		return "unknown"
	}
}

func extractSummary(text string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 20 && len(line) < 200 {
			return line
		}
	}
	if len(text) > 100 {
		return text[:100] + "..."
	}
	return text
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
