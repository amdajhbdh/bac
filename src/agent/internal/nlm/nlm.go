package nlm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type ResearchResult struct {
	NotebookID string
	Results    string
	Success    bool
	Error      string
}

type NotebookInfo struct {
	ID    string
	Title string
}

var defaultNotebookID = "16b01950-5766-4353-8bed-c7f67966cb6b"

func GetNotebooks(ctx context.Context) ([]NotebookInfo, error) {
	slog.Info("fetching NLM notebooks")

	cmd := exec.CommandContext(ctx, "nlm", "notebook", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("failed to list notebooks", "error", err)
		return nil, err
	}

	var notebooks []NotebookInfo
	if err := json.Unmarshal(output, &notebooks); err != nil {
		slog.Warn("failed to parse notebooks", "error", err)
		return nil, err
	}

	return notebooks, nil
}

func Query(ctx context.Context, notebookID, question string) ResearchResult {
	slog.Info("querying NLM", "notebook", notebookID, "question", question)

	if notebookID == "" {
		notebookID = defaultNotebookID
	}

	// Try NLM query
	cmd := exec.CommandContext(ctx, "nlm", "notebook", "query", notebookID, question)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Warn("NLM query failed", "error", err, "stderr", stderr.String())
		return ResearchResult{
			NotebookID: notebookID,
			Success:    false,
			Error:      err.Error(),
		}
	}

	result := strings.TrimSpace(stdout.String())
	slog.Info("NLM query successful", "result_len", len(result))

	return ResearchResult{
		NotebookID: notebookID,
		Results:    result,
		Success:    true,
	}
}

func Research(ctx context.Context, problem string) ResearchResult {
	slog.Info("NLM research started", "problem", problem)

	// Try to find relevant notebooks
	notebooks, err := GetNotebooks(ctx)
	if err != nil || len(notebooks) == 0 {
		slog.Warn("no notebooks found, using default")
		return Query(ctx, defaultNotebookID, problem)
	}

	// Query the first available notebook
	notebookID := notebooks[0].ID
	if notebookID == "" {
		notebookID = defaultNotebookID
	}

	slog.Info("using notebook", "id", notebookID, "title", notebooks[0].Title)
	return Query(ctx, notebookID, problem)
}

func CreateNotebook(ctx context.Context, title string) (string, error) {
	slog.Info("creating NLM notebook", "title", title)

	cmd := exec.CommandContext(ctx, "nlm", "notebook", "create", title)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("create notebook: %w", err)
	}

	// Parse output to get ID
	outputStr := strings.TrimSpace(string(output))
	parts := strings.Split(outputStr, ":")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[1]), nil
	}

	return outputStr, nil
}

func AddSource(ctx context.Context, notebookID, sourceType, source string) error {
	slog.Info("adding source to notebook", "notebook", notebookID, "type", sourceType)

	var cmd *exec.Cmd
	switch sourceType {
	case "url":
		cmd = exec.CommandContext(ctx, "nlm", "source", "add", notebookID, "--url", source)
	case "text":
		cmd = exec.CommandContext(ctx, "nlm", "source", "add", notebookID, "--text", source)
	default:
		return fmt.Errorf("unsupported source type: %s", sourceType)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("add source: %w", err)
	}

	return nil
}

func GenerateQuiz(ctx context.Context, notebookID string) (string, error) {
	slog.Info("generating quiz", "notebook", notebookID)

	cmd := exec.CommandContext(ctx, "nlm", "quiz", "create", notebookID, "--confirm")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("generate quiz: %w", err)
	}

	return string(output), nil
}

func GenerateAudio(ctx context.Context, notebookID string) (string, error) {
	slog.Info("generating audio", "notebook", notebookID)

	cmd := exec.CommandContext(ctx, "nlm", "audio", "create", notebookID, "--confirm")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("generate audio: %w", err)
	}

	return string(output), nil
}

type AnimationContext struct {
	Concepts       []string
	VisualElements []string
	AudioCues      []string
	Explanation    string
}

func GetAnimationContext(ctx context.Context, problem string) AnimationContext {
	slog.Info("getting animation context from NLM", "problem", problem)

	result := Research(ctx, problem)
	if !result.Success {
		slog.Warn("NLM fallback failed, using defaults")
		return AnimationContext{
			Concepts:       []string{"Math"},
			VisualElements: []string{"text", "shapes"},
			AudioCues:      []string{"step"},
			Explanation:    problem,
		}
	}

	concepts := extractConcepts(result.Results)
	visuals := extractVisualElements(result.Results)

	return AnimationContext{
		Concepts:       concepts,
		VisualElements: visuals,
		AudioCues:      []string{"step", "complete"},
		Explanation:    result.Results,
	}
}

func extractConcepts(text string) []string {
	keywords := []string{"equation", "function", "derivative", "integral", "matrix", "vector", "probability", "statistics", "geometry", "algebra", "trigonometry", "calculus"}
	var found []string
	lower := strings.ToLower(text)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			found = append(found, kw)
		}
	}
	if len(found) == 0 {
		found = []string{"math"}
	}
	return found
}

func extractVisualElements(text string) []string {
	var elements []string
	lower := strings.ToLower(text)

	if strings.Contains(lower, "graph") || strings.Contains(lower, "plot") || strings.Contains(lower, "curve") {
		elements = append(elements, "graph")
	}
	if strings.Contains(lower, "shape") || strings.Contains(lower, "circle") || strings.Contains(lower, "triangle") || strings.Contains(lower, "square") {
		elements = append(elements, "shapes")
	}
	if strings.Contains(lower, "number") || strings.Contains(lower, "equation") || strings.Contains(lower, "formula") {
		elements = append(elements, "text")
	}
	if strings.Contains(lower, "animate") || strings.Contains(lower, "motion") || strings.Contains(lower, "move") {
		elements = append(elements, "animation")
	}

	if len(elements) == 0 {
		elements = []string{"text", "shapes"}
	}
	return elements
}
