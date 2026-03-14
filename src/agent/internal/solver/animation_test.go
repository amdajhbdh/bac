package solver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSolverGeneratesAnimation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/animation" {
			t.Errorf("Expected /animation, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		var req AnimationRequest
		json.NewDecoder(r.Body).Decode(&req)

		if !strings.Contains(req.Code, "from manim import") {
			t.Error("Expected Manim code")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AnimationResponse{
			JobID:  "test-job-123",
			Status: "queued",
		})
	}))
	defer server.Close()

	solveResult := &SolveResult{
		Solution:   "x = 5",
		Steps:      2,
		Confidence: 0.95,
		Subject:    "mathematics",
		Chapter:    "equations",
		Concepts:   []string{"linear equations", "algebra"},
		Model:      "llama3.2:3b",
	}

	manimCode := generateManimCode("solve x + 3 = 8", solveResult.Solution)

	client := &http.Client{Timeout: 5 * time.Second}
	animReq := AnimationRequest{
		Code:    manimCode,
		Quality: "low",
	}
	body, _ := json.Marshal(animReq)

	resp, err := client.Post(server.URL+"/animation", "application/json", strings.NewReader(string(body)))
	if err != nil {
		t.Fatalf("Failed to call animation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	var animResp AnimationResponse
	json.NewDecoder(resp.Body).Decode(&animResp)

	if animResp.JobID == "" {
		t.Error("Expected job ID")
	}

	if animResp.Status != "queued" {
		t.Errorf("Expected queued, got %s", animResp.Status)
	}

	t.Logf("Generated Manim code:\n%s", manimCode)
}

func TestGenerateManimCode(t *testing.T) {
	problem := "Résoudre l'équation x² + 2x + 1 = 0"
	solution := `x² + 2x + 1 = 0
(x + 1)² = 0
x = -1`

	code := generateManimCode(problem, solution)

	expectedContains := []string{
		"from manim import",
		"class Solution(Scene):",
		"MathTex",
		"Write(equation)",
	}

	for _, exp := range expectedContains {
		if !strings.Contains(code, exp) {
			t.Errorf("Expected code to contain %s", exp)
		}
	}

	if !strings.Contains(code, problem) {
		t.Error("Expected problem comment in code")
	}
}

func TestSolveWithAnimation(t *testing.T) {
	problem := "Résoudre: 2x + 5 = 13"

	ctx := context.Background()

	result, err := Solve(ctx, problem, "")
	if err != nil {
		t.Logf("Solve returned error (expected if Ollama not running): %v", err)
	}

	if result != nil {
		t.Logf("Solution: %s", result.Solution)
		t.Logf("Steps: %d", result.Steps)
		t.Logf("Confidence: %.2f", result.Confidence)
		t.Logf("Model: %s", result.Model)

		manimCode := generateManimCode(problem, result.Solution)
		t.Logf("Manim code generated: %d bytes", len(manimCode))

		if len(manimCode) == 0 {
			t.Error("Expected non-empty Manim code")
		}
	}
}
