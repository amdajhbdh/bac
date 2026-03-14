package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/bac-unified/agent/internal/animation"
	"github.com/bac-unified/agent/internal/memory"
	"github.com/bac-unified/agent/internal/nlm"
	"github.com/bac-unified/agent/internal/online"
	"github.com/bac-unified/agent/internal/solver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Parse flags
	animate := flag.Bool("animate", false, "Generate animation")
	forceOnline := flag.Bool("online", false, "Force online auto-solve")
	serverMode := flag.Bool("server", false, "Run as HTTP server")
	port := flag.String("port", "8081", "Server port")
	flag.Parse()

	// Server mode
	if *serverMode {
		slog.Info("starting agent server", "port", *port)
		http.HandleFunc("/solve", handleSolve)
		http.HandleFunc("/health", handleHealth)
		err := http.ListenAndServe(":"+*port, nil)
		if err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
		return
	}

	// CLI mode
	if flag.NArg() < 1 {
		logger.Error("no problem provided", "usage", "bac-agent [-animate] [-server] [-port 8081] <problem>")
		os.Exit(1)
	}

	problem := strings.Join(flag.Args(), " ")
	if problem == "" {
		logger.Error("empty problem")
		os.Exit(1)
	}

	ctx := context.Background()

	slog.Info("agent started", "problem", problem)

	// Stage 1: Memory lookup
	slog.Info("stage: memory-lookup")
	memoryResult := memory.Lookup(ctx, problem, 3)
	if len(memoryResult.SimilarProblems) > 0 {
		slog.Info("found similar problems", "count", len(memoryResult.SimilarProblems))
	}

	// Stage 2: Solve
	slog.Info("stage: solving")
	solveResult, err := solver.Solve(ctx, problem, memoryResult.Context)
	if err != nil {
		slog.Error("solver failed", "error", err)
		solveResult = solver.FallbackSolve(problem)
	}

	slog.Info("solution generated", "steps", solveResult.Steps, "model", solveResult.Model)

	// Stage 3: Online auto-solve (low confidence or forced)
	var onlineProvider, onlineResultText string
	if solveResult.Confidence < 0.7 || *forceOnline {
		reason := "low-confidence"
		if *forceOnline {
			reason = "forced"
		}
		slog.Info("stage: online-auto-solve", "reason", reason, "confidence", solveResult.Confidence)

		// Try auto-solve with online AI services
		autoSolveResult := online.AutoSolve(ctx, problem)
		if autoSolveResult.Success {
			slog.Info("auto-solve succeeded", "provider", autoSolveResult.Provider, "result_len", len(autoSolveResult.Results))
			onlineProvider = autoSolveResult.Provider
			onlineResultText = autoSolveResult.Results

			// If online AI provided a better solution, use it
			if autoSolveResult.Results != "" {
				solveResult.Solution = autoSolveResult.Results
				solveResult.Confidence = 1.0
				solveResult.Model = autoSolveResult.Provider
			}
		} else {
			slog.Warn("auto-solve failed", "error", autoSolveResult.Error)

			// Fallback: open browser for manual interaction
			slog.Info("stage: online-research-fallback")
			onlineRes := online.Research(ctx, problem)
			slog.Info("online research complete", "provider", onlineRes.Provider)
			onlineProvider = onlineRes.Provider
		}

		// Also try NLM research with caching
		slog.Info("stage: nlm-research")
		nlmResult := nlm.ResearchWithCache(ctx, problem)
		if nlmResult.Success {
			slog.Info("NLM research complete", "notebook", nlmResult.NotebookID, "result_len", len(nlmResult.Results))
		}
	}

	// Stage 4: Generate animation (if requested)
	var animPath string
	if *animate {
		slog.Info("stage: animation")
		animResult := animation.CompileAndExport(ctx, problem, solveResult.Solution, animation.DefaultConfig)
		if animResult.Success {
			animPath = animResult.FilePath
			slog.Info("animation generated", "path", animPath, "format", animResult.ExportFormat)
		} else {
			slog.Warn("animation failed", "error", animResult.Error)
		}
	}

	// Stage 5: Store in memory
	slog.Info("stage: learning")
	err = memory.Store(ctx, problem, solveResult.Solution, solveResult.Subject, solveResult.Chapter, solveResult.Concepts)
	if err != nil {
		slog.Warn("storage failed", "error", err)
	}

	// Output solution as structured data
	type Output struct {
		Problem        string   `json:"problem"`
		Solution       string   `json:"solution"`
		Subject        string   `json:"subject"`
		Concepts       []string `json:"concepts"`
		Steps          int      `json:"steps"`
		Confidence     float64  `json:"confidence"`
		Model          string   `json:"model"`
		SimilarLen     int      `json:"similar_found"`
		Animation      string   `json:"animation,omitempty"`
		OnlineProvider string   `json:"online_provider,omitempty"`
		OnlineResult   string   `json:"online_result,omitempty"`
	}

	out := Output{
		Problem:        problem,
		Solution:       solveResult.Solution,
		Subject:        solveResult.Subject,
		Concepts:       solveResult.Concepts,
		Steps:          solveResult.Steps,
		Confidence:     solveResult.Confidence,
		Model:          solveResult.Model,
		SimilarLen:     len(memoryResult.SimilarProblems),
		Animation:      animPath,
		OnlineProvider: onlineProvider,
		OnlineResult:   onlineResultText,
	}

	slog.Info("agent completed", "result", out)
}

func handleSolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	problem := r.FormValue("problem")
	if problem == "" {
		http.Error(w, "problem required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	slog.Info("solving problem via HTTP", "problem", problem)

	memoryResult := memory.Lookup(ctx, problem, 3)
	solveResult, err := solver.Solve(ctx, problem, memoryResult.Context)
	if err != nil {
		slog.Error("solver failed", "error", err)
		solveResult = solver.FallbackSolve(problem)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"problem": "%s", "solution": "%s", "confidence": %f}`, problem, solveResult.Solution, solveResult.Confidence)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
