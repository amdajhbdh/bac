package online

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bac-unified/agent/internal/nlm"
)

type ResearchResult struct {
	Sources  int
	Results  string
	Provider string
	Success  bool
	Error    string
}

var authDir = filepath.Join(os.Getenv("HOME"), ".bac-agent", "auth")

func init() {
	os.MkdirAll(authDir, 0755)
}

func Research(ctx context.Context, problem string) ResearchResult {
	slog.Info("starting online research", "problem", problem)

	// Priority 1: Try API-based solutions (faster, more reliable)
	result := SolveWithAnyAPI(ctx, problem)
	if result.Success {
		return result
	}

	// Priority 2: Fall back to NLM CLI (with caching)
	slog.Info("falling back to NLM research")
	nlmResult := nlm.ResearchWithCache(ctx, problem)
	if nlmResult.Success {
		return ResearchResult{
			Results:  nlmResult.Results,
			Provider: "nlm",
			Success:  true,
		}
	}

	slog.Warn("all online providers failed")
	return ResearchResult{Success: false, Provider: "none"}
}

func ResearchWithNLM(ctx context.Context, problem string, notebookID string) ResearchResult {
	slog.Info("researching with NLM", "notebook", notebookID, "problem", problem)

	cmd := exec.CommandContext(ctx, "nlm", "notebook", "query", notebookID, problem)
	output, err := cmd.Output()
	if err != nil {
		slog.Error("NLM query failed", "error", err)
		return ResearchResult{}
	}

	return ResearchResult{
		Sources:  1,
		Results:  strings.TrimSpace(string(output)),
		Provider: "notebooklm",
		Success:  true,
	}
}

func CreateNotebookForResearch(ctx context.Context, title string) (string, error) {
	slog.Info("creating NLM notebook", "title", title)

	cmd := exec.CommandContext(ctx, "nlm", "notebook", "create", title)
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("nlm CLI not available, using fallback", "error", err)
		return generateNotebookID(title), nil
	}

	outputStr := string(output)
	parts := strings.Split(outputStr, ":")
	if len(parts) < 2 {
		return strings.TrimSpace(outputStr), nil
	}

	return strings.TrimSpace(parts[1]), nil
}

func generateNotebookID(title string) string {
	hash := 0
	for _, c := range title {
		hash = hash*31 + int(c)
	}
	return fmt.Sprintf("notebook-%d", hash%100000)
}

func AddSourceToNotebook(ctx context.Context, notebookID, sourceType, source string) error {
	slog.Info("adding source to notebook", "notebook", notebookID, "type", sourceType)

	var cmd *exec.Cmd
	switch sourceType {
	case "url":
		cmd = exec.CommandContext(ctx, "nlm", "source", "add", notebookID, "--url", source)
	case "text":
		cmd = exec.CommandContext(ctx, "nlm", "source", "add", notebookID, "--text", source)
	case "youtube":
		cmd = exec.CommandContext(ctx, "nlm", "source", "add", notebookID, "--youtube", source)
	case "gdocs":
		cmd = exec.CommandContext(ctx, "nlm", "source", "add", notebookID, "--gdocs", source)
	default:
		return fmt.Errorf("unsupported source type: %s", sourceType)
	}

	if err := cmd.Run(); err != nil {
		slog.Warn("nlm CLI failed, source will be stored locally", "error", err)
		return storeSourceLocally(notebookID, sourceType, source)
	}

	return nil
}

func storeSourceLocally(notebookID, sourceType, source string) error {
	sourceDir := filepath.Join(os.Getenv("HOME"), ".bac-agent", "nlm-sources")
	os.MkdirAll(sourceDir, 0755)

	sourceFile := filepath.Join(sourceDir, fmt.Sprintf("%s_%s.json", notebookID, sourceType))
	data := fmt.Sprintf(`{"notebook_id": "%s", "source_type": "%s", "source": "%s", "added_at": "%s"}`,
		notebookID, sourceType, source, time.Now().Format(time.RFC3339))

	return os.WriteFile(sourceFile, []byte(data), 0644)
}

func AutoSolve(ctx context.Context, problem string) ResearchResult {
	slog.Info("starting auto-solve", "problem", problem)

	// Priority 1: Try API first (faster, more reliable)
	result := SolveWithAnyAPI(ctx, problem)
	if result.Success {
		slog.Info("auto-solve succeeded via API", "provider", result.Provider)
		return result
	}

	// Priority 2: Fall back to NLM (with caching)
	slog.Info("API failed, falling back to NLM research")
	nlmResult := nlm.ResearchWithCache(ctx, problem)
	if nlmResult.Success {
		slog.Info("auto-solve succeeded via NLM", "notebook", nlmResult.NotebookID)
		return ResearchResult{
			Results:  nlmResult.Results,
			Provider: "nlm",
			Success:  true,
		}
	}

	slog.Error("all services failed")
	return ResearchResult{
		Success:  false,
		Provider: "none",
		Error:    "all AI services failed",
	}
}
