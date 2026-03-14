package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Config holds the pipeline configuration
type Config struct {
	OllamaURL   string
	OllamaModel string
	Workers     int
	UseCloud    bool
	CloudModel  string
	Verbose     bool
}

// OCRResult holds the result of OCR processing
type OCRResult struct {
	Page       int    `json:"page"`
	RawText    string `json:"raw_text"`
	Fixed      string `json:"fixed"`
	Structured string `json:"structured"`
	Error      string `json:"error,omitempty"`
}

// OllamaRequest represents a request to Ollama
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents a response from Ollama
type OllamaResponse struct {
	Response string `json:"response"`
}

var cfg = Config{
	OllamaURL:   getEnv("OLLAMA_URL", "http://localhost:11434"),
	OllamaModel: getEnv("OLLAMA_MODEL", "llama3.2:3b"),
	Workers:     8,
	UseCloud:    getEnv("USE_CLOUD_MODEL", "false") == "true",
	CloudModel:  getEnv("CLOUD_MODEL", "minimax-m2.5:cloud"),
	Verbose:     false,
}

var logger *slog.Logger
var startTime time.Time

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Progress printing functions (only output when verbose mode is enabled)
func vPrint(args ...interface{}) {
	if cfg.Verbose {
		fmt.Print(args...)
	}
}

func vPrintln(args ...interface{}) {
	if cfg.Verbose {
		fmt.Println(args...)
	}
}

func vPrintf(format string, args ...interface{}) {
	if cfg.Verbose {
		fmt.Printf(format, args...)
	}
}

func printProgressBar(current, total int) {
	if !cfg.Verbose || total <= 0 {
		return
	}
	percentage := float64(current) / float64(total) * 100
	width := 30
	filled := int(float64(width) * float64(current) / float64(total))
	bar := strings.Repeat("=", filled) + strings.Repeat("-", width-filled)
	fmt.Printf("\r[%s] %d%% (%d/%d)", bar, int(percentage), current, total)
	if current == total {
		fmt.Println()
	}
}

func getElapsedTime() time.Duration {
	return time.Since(startTime)
}

func main() {
	// Parse verbose flag
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-v" || os.Args[i] == "--verbose" {
			cfg.Verbose = true
			// Remove the flag from args for proper command parsing
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			break
		}
	}

	if len(os.Args) < 2 {
		logger.Info("AI-OCR Pipeline - Smart document extraction")
		logger.Info("Usage: ai-ocr [options] <command>")
		logger.Info("Options:")
		logger.Info("  -v, --verbose    Show real-time progress")
		logger.Info("Commands:")
		logger.Info("  extract <pdf>     - Extract single PDF")
		logger.Info("  batch <dir>       - Batch process directory")
		os.Exit(1)
	}

	startTime = time.Now()
	cmd := os.Args[1]
	ctx := context.Background()

	vPrintln("=== AI-OCR Pipeline Started ===")
	vPrintf("Verbose mode: enabled\n")
	vPrintf("Ollama URL: %s\n", cfg.OllamaURL)
	vPrintf("Ollama Model: %s\n", cfg.OllamaModel)
	if cfg.UseCloud {
		vPrintf("Cloud Model: %s\n", cfg.CloudModel)
	}
	vPrintln()

	switch cmd {
	case "extract":
		if len(os.Args) < 3 {
			logger.Error("extract command requires a PDF path")
			logger.Info("Usage: ai-ocr extract <pdf>")
			os.Exit(1)
		}
		if err := extractPDF(ctx, os.Args[2]); err != nil {
			logger.Error("failed to extract PDF", "error", err)
			os.Exit(1)
		}
	case "batch":
		if len(os.Args) < 3 {
			logger.Error("batch command requires a directory path")
			logger.Info("Usage: ai-ocr batch <directory>")
			os.Exit(1)
		}
		if err := batchProcess(ctx, os.Args[2]); err != nil {
			logger.Error("failed to batch process", "error", err)
			os.Exit(1)
		}
	default:
		logger.Error("unknown command", "command", cmd)
		os.Exit(1)
	}
}

func extractPDF(ctx context.Context, pdfPath string) error {
	vPrintf("[%s] Processing PDF...\n", getElapsedTime().Round(time.Second))
	vPrintf("  Path: %s\n", pdfPath)

	// Validate input file
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF file not found: %w", err)
	}

	// Step 1: Raw OCR
	vPrintf("[%s] Converting PDF to images...\n", getElapsedTime().Round(time.Second))
	rawText, totalPages, err := runOCR(ctx, pdfPath)
	if err != nil {
		return fmt.Errorf("running OCR: %w", err)
	}
	if rawText == "" {
		return fmt.Errorf("OCR produced no text")
	}

	vPrintf("[%s] OCR completed (%d pages, %d chars)\n", getElapsedTime().Round(time.Second), totalPages, len(rawText))

	// Step 2: AI Fix - correct errors
	vPrintf("[%s] Running AI fix (step 1/2)...\n", getElapsedTime().Round(time.Second))
	fixed, err := fixOCR(ctx, rawText)
	if err != nil {
		vPrintf("[%s] AI fix failed, using raw text: %v\n", getElapsedTime().Round(time.Second), err)
		fixed = rawText
	} else {
		vPrintf("[%s] AI fix completed (%d chars)\n", getElapsedTime().Round(time.Second), len(fixed))
	}

	// Step 3: AI Structure - convert to markdown
	vPrintf("[%s] Running AI structuring (step 2/2)...\n", getElapsedTime().Round(time.Second))
	structured, err := structureText(ctx, fixed)
	if err != nil {
		vPrintf("[%s] AI structuring failed, using fixed text: %v\n", getElapsedTime().Round(time.Second), err)
		structured = fixed
	} else {
		vPrintf("[%s] AI structuring completed\n", getElapsedTime().Round(time.Second))
	}

	// Save output
	outputPath := calculateOutputPath(pdfPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	content := formatOutput(pdfPath, structured)
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	vPrintf("[%s] Done! Saved to %s\n", getElapsedTime().Round(time.Second), outputPath)
	vPrintf("[%s] === Extraction Complete ===\n", getElapsedTime().Round(time.Second))
	return nil
}

func calculateOutputPath(pdfPath string) string {
	// Remove .pdf extension
	base := strings.TrimSuffix(pdfPath, ".pdf")

	// Replace source directories with output directory
	replacements := []string{
		"/03-Resources/",
		"/07-Assets/PDFs/",
		"/db/pdfs/",
	}

	for _, rep := range replacements {
		if idx := strings.Index(base, rep); idx != -1 {
			base = base[:idx] + "/05-Extracted/" + base[idx+len(rep):]
			break
		}
	}

	return base + ".md"
}

func formatOutput(pdfPath, content string) string {
	return fmt.Sprintf(`---
title: %s
source: %s
date: %s
tags: [ai-ocr, extracted]
---

%s`,
		filepath.Base(pdfPath),
		pdfPath,
		time.Now().Format("2006-01-02"),
		content,
	)
}

func batchProcess(ctx context.Context, dir string) error {
	vPrintf("[%s] Scanning directory: %s\n", getElapsedTime().Round(time.Second), dir)

	var pdfs []string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".pdf") {
			pdfs = append(pdfs, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	if len(pdfs) == 0 {
		logger.Info("no PDFs found")
		return nil
	}

	vPrintf("[%s] Found %d PDFs to process\n", getElapsedTime().Round(time.Second), len(pdfs))

	// Create output directory
	if err := os.MkdirAll("05-Extracted", 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Parallel processing with semaphore
	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.Workers)
	errors := make([]error, 0)
	var errorsMu sync.Mutex

	processed := 0
	var processedMu sync.Mutex

	for i, pdf := range pdfs {
		wg.Add(1)
		sem <- struct{}{}

		go func(p string, idx int) {
			defer func() { <-sem }()
			defer wg.Done()

			if cfg.Verbose {
				vPrintf("\n[%s] Processing [%d/%d]: %s\n", getElapsedTime().Round(time.Second), idx+1, len(pdfs), filepath.Base(p))
			} else {
				logger.Info("processing", "current", idx+1, "total", len(pdfs), "path", p)
			}

			if err := extractPDF(ctx, p); err != nil {
				errorsMu.Lock()
				errors = append(errors, fmt.Errorf("%s: %w", p, err))
				errorsMu.Unlock()
				logger.Error("failed to process PDF", "path", p, "error", err)
			} else {
				processedMu.Lock()
				processed++
				if cfg.Verbose {
					vPrintf("[%s] Progress: %d/%d (%.0f%%)\n", getElapsedTime().Round(time.Second), processed, len(pdfs), float64(processed)/float64(len(pdfs))*100)
				}
				processedMu.Unlock()
			}
		}(pdf, i)
	}

	wg.Wait()

	if len(errors) > 0 {
		vPrintf("[%s] Batch completed with %d errors\n", getElapsedTime().Round(time.Second), len(errors))
		logger.Error("batch completed with errors", "failed", len(errors))
		for _, e := range errors {
			logger.Error("error", "detail", e.Error())
		}
		return fmt.Errorf("%d PDFs failed to process", len(errors))
	}

	vPrintf("[%s] Batch complete! Processed %d PDFs in %s\n", getElapsedTime().Round(time.Second), len(pdfs), getElapsedTime().Round(time.Second))
	logger.Info("batch complete", "processed", len(pdfs))
	return nil
}

func runOCR(ctx context.Context, pdfPath string) (string, int, error) {
	// Create temp directory for images
	tmpDir, err := os.MkdirTemp("", "ocr")
	if err != nil {
		return "", 0, fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Convert PDF to images using pdftoppm
	basePath := filepath.Join(tmpDir, "page")
	cmd := exec.CommandContext(ctx, "pdftoppm", "-r", "200", "-png", pdfPath, basePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", 0, fmt.Errorf("running pdftoppm: %w\n%s", err, string(out))
	}

	// Find generated images
	pattern := basePath + "*.png"
	pages, err := filepath.Glob(pattern)
	if err != nil {
		return "", 0, fmt.Errorf("finding generated images: %w", err)
	}

	if len(pages) == 0 {
		return "", 0, fmt.Errorf("no images generated from PDF")
	}

	// OCR each page in parallel
	type pageResult struct {
		index int
		text  string
		err   error
	}

	results := make(chan pageResult, len(pages))
	var wg sync.WaitGroup

	// Process pages in parallel with progress updates
	for i, page := range pages {
		wg.Add(1)
		go func(idx int, pagePath string) {
			defer wg.Done()
			text, err := ocrPage(ctx, pagePath)
			os.Remove(pagePath) // Clean up immediately
			results <- pageResult{idx, text, err}
			// Progress update
			vPrintf("[%s] Processed page %d/%d\n", getElapsedTime().Round(time.Second), idx+1, len(pages))
			printProgressBar(idx+1, len(pages))
		}(i, page)
	}

	// Close results channel when done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results in order
	var allText []string
	for result := range results {
		if result.err != nil {
			logger.Warn("failed to OCR page", "index", result.index, "error", result.err)
			continue
		}
		if len(result.text) > 50 {
			allText = append(allText, result.text)
		}
	}

	// Sort by page order
	sort.Slice(allText, func(i, j int) bool {
		return i < j
	})

	if len(allText) == 0 {
		return "", 0, fmt.Errorf("no text extracted from any page")
	}

	return strings.Join(allText, "\n\n---\n\n"), len(pages), nil
}

func ocrPage(ctx context.Context, pagePath string) (string, error) {
	cmd := exec.CommandContext(ctx, "tesseract", "-l", "ara+fra+eng", pagePath, "stdout")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("running tesseract: %w\n%s", err, string(out))
	}
	return string(out), nil
}

func fixOCR(ctx context.Context, text string) (string, error) {
	// Truncate text if too long
	maxLen := 4000
	if len(text) > maxLen {
		text = text[:maxLen]
	}

	prompt := fmt.Sprintf(`You are an expert OCR error corrector for French and Arabic educational documents.

Your task is to fix OCR errors while preserving the original meaning and structure.

IMPORTANT RULES:
1. Fix Arabic character confusion (ا-إ-آ should be corrected based on context)
2. Fix French accents (éèêë, àâä, ïîöôùûü)
3. Fix common OCR errors: rn→m, vv→w, ll→li, 0→O in context
4. Preserve mathematical notation and formulas exactly as written
5. Keep numbers and technical terms intact
6. Maintain paragraph structure

OCR Output:
%s

Return ONLY the corrected text, no explanations, no markdown formatting.`, text)

	return callOllama(ctx, prompt)
}

func structureText(ctx context.Context, text string) (string, error) {
	// Truncate text if too long
	maxLen := 5000
	if len(text) > maxLen {
		text = text[:maxLen]
	}

	prompt := fmt.Sprintf(`You are a document structurer for educational content.

Transform this text into well-organized Markdown with the following rules:

1. Add proper headings (## for major sections, ### for subsections)
2. Use bullet points (• or -) for lists
3. Preserve mathematical formulas in $...$ or $$...$$ format
4. Use > [!note] or > [!warning] callouts for important content
5. Create wikilinks between related concepts using [[concept]] notation
6. Add a brief summary at the top in a callout
7. Keep the content in French/Arabic as original

Text:
%s

Return ONLY the structured markdown, no explanations.`, text)

	return callOllama(ctx, prompt)
}

func callOllama(ctx context.Context, prompt string) (string, error) {
	// Use cloud model if configured
	model := cfg.OllamaModel
	if cfg.UseCloud {
		model = cfg.CloudModel
	}

	reqBody, err := json.Marshal(OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.OllamaURL+"/api/generate", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama returned status %d", resp.StatusCode)
	}

	var result OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}

	return result.Response, nil
}
