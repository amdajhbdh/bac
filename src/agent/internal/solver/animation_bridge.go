package solver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type AnimationBridge struct {
	gatewayURL string
	client     *http.Client
	enabled    bool
}

func NewAnimationBridge() *AnimationBridge {
	gatewayURL := os.Getenv("GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "http://localhost:8081"
	}

	enabled := os.Getenv("ANIMATION_ENABLED") == "true"

	return &AnimationBridge{
		gatewayURL: gatewayURL,
		client:     &http.Client{Timeout: 10 * time.Second},
		enabled:    enabled,
	}
}

func (a *AnimationBridge) GenerateAnimation(ctx context.Context, problem, solution string) (string, error) {
	if !a.enabled {
		slog.Debug("animation disabled")
		return "", nil
	}

	code := generateManimCode(problem, solution)

	req := AnimationRequest{
		Code:    code,
		Quality: "low",
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.gatewayURL+"/animation", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("call gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gateway returned %d", resp.StatusCode)
	}

	var animResp AnimationResponse
	if err := json.NewDecoder(resp.Body).Decode(&animResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	slog.Info("animation queued", "job_id", animResp.JobID, "status", animResp.Status)

	return animResp.JobID, nil
}

func generateManimCode(problem, solution string) string {
	var buf strings.Builder
	buf.WriteString("from manim import *\n\n")
	buf.WriteString("class Solution(Scene):\n")
	buf.WriteString("    def construct(self):\n")
	buf.WriteString("        # Problem: " + problem + "\n")
	buf.WriteString("        # Solution\n")

	lines := strings.Split(solution, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		buf.WriteString("        equation = MathTex(r\"" + strings.ReplaceAll(line, "\"", "\\\"") + "\")\n")
		buf.WriteString("        self.play(Write(equation))\n")
		buf.WriteString("        self.wait()\n")
	}

	return buf.String()
}
