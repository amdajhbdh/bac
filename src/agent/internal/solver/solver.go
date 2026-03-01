package solver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 120 * time.Second,
}

type SolveResult struct {
	Solution   string
	Steps      int
	Confidence float64
	Subject    string
	Chapter    string
	Concepts   []string
	Model      string
}

type ModelInfo struct {
	Name    string
	IsCloud bool
	IsLocal bool
	Size    int64
}

var (
	availableModels []ModelInfo
	modelPriority   = []string{
		"llama3.2:3b",
		"test:latest",
		"ministral-3:14b-cloud",
		"deepseek-v3.2:cloud",
		"deepseek-v3.1:671b-cloud",
		"minimax-m2.5:cloud",
		"kimi-k2.5:cloud",
		"glm-5:cloud",
		"qwen3-coder:480b-cloud",
	}
)

func init() {
	refreshModels()
}

func refreshModels() {
	models, err := listModels()
	if err != nil {
		slog.Warn("failed to list models", "error", err)
		availableModels = []ModelInfo{{Name: "llama3.2:3b", IsLocal: true}}
		return
	}
	availableModels = models
	slog.Info("models refreshed", "count", len(availableModels))
}

func listModels() ([]ModelInfo, error) {
	resp, err := httpClient.Get("http://127.0.0.1:11434/api/tags")
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name   string `json:"name"`
			Size   int64  `json:"size"`
			Model  string `json:"model"`
			Remote string `json:"remote_model,omitempty"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	var models []ModelInfo
	for _, m := range result.Models {
		isCloud := m.Remote != "" || strings.Contains(m.Name, "-cloud")
		models = append(models, ModelInfo{
			Name:    m.Name,
			IsCloud: isCloud,
			IsLocal: !isCloud,
			Size:    m.Size,
		})
	}

	return models, nil
}

func selectBestModel() string {
	if len(availableModels) == 0 {
		refreshModels()
		if len(availableModels) == 0 {
			return "llama3.2:3b"
		}
	}

	// Prefer local models first
	for _, m := range availableModels {
		if m.IsLocal {
			for _, preferred := range modelPriority {
				if m.Name == preferred {
					slog.Info("selected model", "model", m.Name, "cloud", m.IsCloud)
					return m.Name
				}
			}
		}
	}

	// Fall back to any available model
	selected := availableModels[0]
	slog.Info("fallback to first available model", "model", selected.Name)
	return selected.Name
}

func Solve(ctx context.Context, problem, memoryContext string) (*SolveResult, error) {
	slog.Info("solving problem", "problem", problem)

	prompt := buildPrompt(problem, memoryContext)

	result, model, err := callWithFallback(ctx, prompt)
	if err != nil {
		slog.Error("all models failed", "error", err)
		return nil, err
	}

	solveResult := parseResult(result, problem)
	solveResult.Model = model

	slog.Info("problem solved", "model", model, "steps", solveResult.Steps)
	return solveResult, nil
}

func SolveOffline(ctx context.Context, problem, memoryContext string) (*SolveResult, error) {
	slog.Info("solving problem offline", "problem", problem)

	// Check if Ollama is available
	_, err := httpClient.Get("http://127.0.0.1:11434/api/tags")
	if err != nil {
		slog.Warn("Ollama not available, using offline fallback")
		return solveOfflineFallback(problem)
	}

	// If Ollama is available, use regular solve
	return Solve(ctx, problem, memoryContext)
}

func solveOfflineFallback(problem string) (*SolveResult, error) {
	// Offline fallback using rule-based solving
	problem = strings.ToLower(problem)

	result := &SolveResult{
		Solution:   "Solution hors ligne non disponible. Connectez-vous à internet pour utiliser l'IA.",
		Confidence: 0.0,
		Model:      "offline-fallback",
	}

	// Simple pattern matching for basic math
	if strings.Contains(problem, "équation") || strings.Contains(problem, "equation") {
		result.Steps = 1
		if strings.Contains(problem, "x²") || strings.Contains(problem, "x^2") {
			result.Concepts = []string{"équations quadratiques", "discriminant"}
		}
	}

	if strings.Contains(problem, "dériv") {
		result.Steps = 1
		result.Concepts = []string{"dérivée", "calcul différentiel"}
	}

	if strings.Contains(problem, "intégr") {
		result.Steps = 1
		result.Concepts = []string{"intégrale", "calcul intégral"}
	}

	return result, nil
}

func IsOnline() bool {
	_, err := httpClient.Get("http://127.0.0.1:11434/api/tags")
	return err == nil
}

func buildPrompt(problem, memoryContext string) string {
	contextPart := ""
	if memoryContext != "" {
		contextPart = fmt.Sprintf(`Contexte de problèmes similaires:
%s

`, memoryContext)
	}

	return fmt.Sprintf(`Tu es un professeur BAC C. 
%s
Résous ce problème étape par étape en français:

Problème: %s

Format requis:
**Solution:** [réponse]
**Étapes:**
1. [étape 1]
2. [étape 2]
3. [étape 3]
**Concepts:** [liste des concepts]`, contextPart, problem)
}

func callWithFallback(ctx context.Context, prompt string) (string, string, error) {
	if len(availableModels) == 0 {
		refreshModels()
	}

	tried := map[string]bool{}
	maxAttempts := len(availableModels)
	if maxAttempts == 0 {
		maxAttempts = 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		model := selectBestModelForAttempt(attempt)
		if tried[model] {
			continue
		}
		tried[model] = true

		slog.Debug("trying model", "model", model)

		result, err := callOllama(ctx, prompt, model)
		if err != nil {
			slog.Warn("model failed, trying next", "model", model, "error", err)
			continue
		}

		return result, model, nil
	}

	return "", "", fmt.Errorf("all models failed")
}

func selectBestModelForAttempt(attempt int) string {
	if len(availableModels) == 0 {
		refreshModels()
		if len(availableModels) == 0 {
			return "llama3.2:3b"
		}
	}

	// Filter local models first
	var localModels []ModelInfo
	var cloudModels []ModelInfo
	for _, m := range availableModels {
		if m.IsLocal {
			localModels = append(localModels, m)
		} else {
			cloudModels = append(cloudModels, m)
		}
	}

	// Try local models first, then cloud
	allModels := append(localModels, cloudModels...)
	if len(allModels) == 0 {
		return "llama3.2:3b"
	}

	idx := attempt % len(allModels)
	selected := allModels[idx]
	slog.Info("selected model", "model", selected.Name, "cloud", selected.IsCloud)
	return selected.Name
}

func callOllama(ctx context.Context, prompt, model string) (string, error) {
	type GenerateRequest struct {
		Model       string  `json:"model"`
		Prompt      string  `json:"prompt"`
		Stream      bool    `json:"stream"`
		Temperature float64 `json:"temperature"`
	}

	reqBody, _ := json.Marshal(GenerateRequest{
		Model:       model,
		Prompt:      prompt,
		Stream:      false,
		Temperature: 0.3,
	})

	resp, err := httpClient.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	return response, nil
}

func parseResult(response, problem string) *SolveResult {
	stepCount := countSteps(strings.Split(response, "\n"))
	concepts := extractConcepts(response)
	subject := determineSubject(problem)

	return &SolveResult{
		Solution:   response,
		Steps:      stepCount,
		Confidence: 0.85,
		Subject:    subject,
		Chapter:    "",
		Concepts:   concepts,
	}
}

func countSteps(lines []string) int {
	count := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && (strings.HasPrefix(trimmed, "1.") ||
			strings.HasPrefix(trimmed, "2.") ||
			strings.HasPrefix(trimmed, "3.") ||
			strings.HasPrefix(trimmed, "4.") ||
			strings.HasPrefix(trimmed, "5.")) {
			count++
		}
	}
	if count == 0 {
		return 3
	}
	return count
}

func extractConcepts(response string) []string {
	var concepts []string
	lines := strings.Split(response, "\n")
	inConcepts := false

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "concept") {
			inConcepts = true
			continue
		}
		if inConcepts && strings.ContainsAny(line, "0123456789") {
			break
		}
		if inConcepts && len(strings.TrimSpace(line)) > 0 {
			concept := strings.TrimSpace(line)
			concept = strings.Trim(concept, "-*: ")
			if len(concept) > 0 {
				concepts = append(concepts, concept)
			}
		}
	}

	if len(concepts) == 0 {
		concepts = []string{"algebra", "equations"}
	}

	return concepts
}

func determineSubject(problem string) string {
	problemLower := strings.ToLower(problem)

	subjects := map[string][]string{
		"math":        {"équation", "fonction", "dérivée", "intégrale", "limite", "polynôme", "factorisation", "x²", "sqrt"},
		"pc":          {"physique", "mécanique", "électrique", "optique", "thermodynamique", "chimie", "molécule", "atome", "force", "vitesse"},
		"svt":         {"biologie", "cellule", "ADN", "génétique", "écologie", "évolution", "organisme"},
		"philosophie": {"philosophie", "morale", "éthique", "existence", "liberté", "justice"},
	}

	for subject, keywords := range subjects {
		for _, kw := range keywords {
			if strings.Contains(problemLower, kw) {
				return subject
			}
		}
	}

	return "math"
}

func FallbackSolve(problem string) *SolveResult {
	return &SolveResult{
		Solution:   fmt.Sprintf("Résolution de: %s\n\nContactez un tuteur pour une solution complète.", problem),
		Steps:      1,
		Confidence: 0.3,
		Subject:    determineSubject(problem),
		Chapter:    "",
		Concepts:   []string{},
		Model:      "fallback",
	}
}

func GetAvailableModels() []ModelInfo {
	if len(availableModels) == 0 {
		refreshModels()
	}
	return availableModels
}
