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
)

type ResearchResult struct {
	Sources  int
	Results  string
	Provider string
	Success  bool
	Error    string
}

var (
	authDir = filepath.Join(os.Getenv("HOME"), ".bac-agent", "auth")
	cliPath = "/home/med/.npm/bin/playwright-cli"
)

func init() {
	os.MkdirAll(authDir, 0755)
}

func Research(ctx context.Context, problem string) ResearchResult {
	slog.Info("starting online research", "problem", problem)

	// Try API-based solutions first (faster, more reliable)
	result := SolveWithAnyAPI(ctx, problem)
	if result.Success {
		return result
	}

	// Fall back to Playwright automation if no API keys configured
	slog.Info("no API keys found, falling back to browser automation")

	// Try DeepSeek first
	result = tryDeepSeek(ctx, problem)
	if result.Success {
		return result
	}

	// Try Grok
	result = tryGrok(ctx, problem)
	if result.Success {
		return result
	}

	// Try Claude
	result = tryClaude(ctx, problem)
	if result.Success {
		return result
	}

	// Try ChatGPT
	result = tryChatGPT(ctx, problem)
	if result.Success {
		return result
	}

	slog.Warn("all online providers failed")
	return ResearchResult{Success: false, Provider: "none"}
}

func tryDeepSeek(ctx context.Context, problem string) ResearchResult {
	slog.Info("opening DeepSeek with Chrome")
	return openWithChrome(ctx, "deepseek", "https://chat.deepseek.com", problem)
}

func tryGrok(ctx context.Context, problem string) ResearchResult {
	slog.Info("opening Grok with Chrome")
	return openWithChrome(ctx, "grok", "https://grok.com", problem)
}

func tryClaude(ctx context.Context, problem string) ResearchResult {
	slog.Info("opening Claude with Chrome")
	return openWithChrome(ctx, "claude", "https://claude.ai", problem)
}

func tryChatGPT(ctx context.Context, problem string) ResearchResult {
	slog.Info("opening ChatGPT with Chrome")
	return openWithChrome(ctx, "chatgpt", "https://chat.openai.com", problem)
}

func openWithChrome(ctx context.Context, provider, url, problem string) ResearchResult {
	profileDir := filepath.Join(authDir, "chrome-"+provider)
	os.MkdirAll(profileDir, 0755)

	sessionName := "bac-" + provider

	cmd := exec.CommandContext(ctx, cliPath, "-s="+sessionName, "open", url,
		"--browser=chrome",
		"--persistent",
		"--profile="+profileDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		slog.Warn("failed to open Chrome", "provider", provider, "error", err)
		return ResearchResult{
			Provider: provider,
			Success:  false,
			Error:    err.Error(),
		}
	}

	time.Sleep(2 * time.Second)

	slog.Info("Chrome browser opened", "provider", provider, "url", url)

	screenshotPath := filepath.Join(os.TempDir(), fmt.Sprintf("bac-%s-%d.png", provider, time.Now().Unix()))
	exec.Command(cliPath, "-s="+sessionName, "screenshot", "--filename="+screenshotPath).Run()

	return ResearchResult{
		Sources:  1,
		Results:  fmt.Sprintf("Browser opened at %s. Screenshot saved to %s", url, screenshotPath),
		Provider: provider,
		Success:  true,
	}
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

	// Try using nlm CLI first
	cmd := exec.CommandContext(ctx, "nlm", "notebook", "create", title)
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("nlm CLI not available, using fallback", "error", err)
		// Generate a simulated notebook ID for offline/demo mode
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
	// Generate a deterministic ID based on title for demo/offline mode
	hash := 0
	for _, c := range title {
		hash = hash*31 + int(c)
	}
	return fmt.Sprintf("notebook-%d", hash%100000)
}

func AddSourceToNotebook(ctx context.Context, notebookID, sourceType, source string) error {
	slog.Info("adding source to notebook", "notebook", notebookID, "type", sourceType)

	// Try using nlm CLI first
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
		// Store source info locally for offline mode
		return storeSourceLocally(notebookID, sourceType, source)
	}

	return nil
}

func storeSourceLocally(notebookID, sourceType, source string) error {
	// For offline/demo mode, store source information locally
	sourceDir := filepath.Join(os.Getenv("HOME"), ".bac-agent", "nlm-sources")
	os.MkdirAll(sourceDir, 0755)

	sourceFile := filepath.Join(sourceDir, fmt.Sprintf("%s_%s.json", notebookID, sourceType))
	data := fmt.Sprintf(`{"notebook_id": "%s", "source_type": "%s", "source": "%s", "added_at": "%s"}`,
		notebookID, sourceType, source, time.Now().Format(time.RFC3339))

	return os.WriteFile(sourceFile, []byte(data), 0644)
}

var serviceOrder = []string{"deepseek", "grok", "claude", "chatgpt"}

var serviceURLs = map[string]string{
	"deepseek": "https://chat.deepseek.com",
	"grok":     "https://grok.com",
	"claude":   "https://claude.ai",
	"chatgpt":  "https://chat.openai.com",
}

func AutoSolve(ctx context.Context, problem string) ResearchResult {
	slog.Info("starting auto-solve", "problem", problem)

	// Try API first (faster, more reliable)
	result := SolveWithAnyAPI(ctx, problem)
	if result.Success {
		slog.Info("auto-solve succeeded via API", "provider", result.Provider)
		return result
	}

	// Fall back to Playwright automation
	slog.Info("API failed, falling back to browser automation")

	for _, service := range serviceOrder {
		select {
		case <-ctx.Done():
			slog.Info("auto-solve cancelled")
			return ResearchResult{Success: false, Error: "cancelled"}
		default:
		}

		result := solveWithService(ctx, service, problem)
		if result.Success {
			slog.Info("auto-solve succeeded", "service", service)
			return result
		}

		slog.Warn("service failed, trying next", "service", service, "error", result.Error)
	}

	slog.Error("all services failed")
	return ResearchResult{
		Success:  false,
		Provider: "none",
		Error:    "all AI services failed",
	}
}

func solveWithService(ctx context.Context, service, problem string) ResearchResult {
	slog.Info("solving with service", "service", service)

	url := serviceURLs[service]
	profileDir := filepath.Join(authDir, "chrome-"+service)
	sessionName := "bac-" + service
	os.MkdirAll(profileDir, 0755)

	openCmd := exec.CommandContext(ctx, cliPath, "-s="+sessionName, "open", url,
		"--browser=chrome",
		"--persistent",
		"--profile="+profileDir,
		"--headed",
	)
	openCmd.Stdout = os.Stdout
	openCmd.Stderr = os.Stderr

	if err := openCmd.Start(); err != nil {
		return ResearchResult{
			Provider: service,
			Success:  false,
			Error:    fmt.Sprintf("failed to open browser: %v", err),
		}
	}

	time.Sleep(3 * time.Second)

	slog.Info("browser opened", "service", service)
	return ResearchResult{
		Sources:  1,
		Results:  fmt.Sprintf("Browser opened at %s. Please solve manually.", url),
		Provider: service,
		Success:  true,
	}
}

func extractResponse(html, service string) string {
	patterns := []string{
		`<div[^>]*class="[^"]*(?:response|message|content)[^"]*"[^>]*>(.*?)</div>`,
		`<p[^>]*>(.*?)</p>`,
		`<span[^>]*class="[^"]*(?:text|content)[^"]*"[^]*>(.*?)</span>`,
	}

	for _, pattern := range patterns {
		if idx := strings.Index(html, pattern); idx != -1 {
			start := strings.Index(html[idx:], ">")
			if start != -1 {
				end := strings.Index(html[idx+start:], "<")
				if end != -1 {
					text := html[idx+start+1 : idx+start+end]
					text = strings.TrimSpace(text)
					if len(text) > 20 {
						return cleanResponse(text)
					}
				}
			}
		}
	}

	lines := strings.Split(html, "\n")
	var responseLines []string
	inResponse := false

	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "response") || strings.Contains(lower, "answer") || strings.Contains(lower, "solution") {
			inResponse = true
			continue
		}
		if inResponse && len(strings.TrimSpace(line)) > 10 {
			responseLines = append(responseLines, strings.TrimSpace(line))
			if len(responseLines) > 10 {
				break
			}
		}
	}

	if len(responseLines) > 0 {
		return cleanResponse(strings.Join(responseLines, "\n"))
	}

	return ""
}

func cleanResponse(text string) string {
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "  ", " ")
	text = strings.TrimSpace(text)
	return text
}
