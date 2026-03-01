package online

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage,omitempty"`
}

type Choice struct {
	Message ChatMessage `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ClaudeRequest struct {
	Model       string          `json:"model"`
	Messages    []ClaudeMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
	Usage   Usage           `json:"usage,omitempty"`
}

type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type DeepSeekClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	return &DeepSeekClient{
		apiKey:     apiKey,
		baseURL:    "https://api.deepseek.com",
		model:      "deepseek-chat",
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *DeepSeekClient) Solve(ctx context.Context, problem string) (string, error) {
	url := c.baseURL + "/chat/completions"

	reqBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a math tutor for BAC C exam students in Mauritania. Provide clear, step-by-step solutions."},
			{Role: "user", Content: problem},
		},
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	slog.Info("deepseek solved", "tokens", chatResp.Usage.TotalTokens)
	return chatResp.Choices[0].Message.Content, nil
}

type GrokClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewGrokClient(apiKey string) *GrokClient {
	return &GrokClient{
		apiKey:     apiKey,
		baseURL:    "https://api.x.ai/v1",
		model:      "grok-2-1212",
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *GrokClient) Solve(ctx context.Context, problem string) (string, error) {
	url := c.baseURL + "/chat/completions"

	reqBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a math tutor for BAC C exam students in Mauritania. Provide clear, step-by-step solutions."},
			{Role: "user", Content: problem},
		},
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	slog.Info("grok solved", "tokens", chatResp.Usage.TotalTokens)
	return chatResp.Choices[0].Message.Content, nil
}

type ClaudeClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewClaudeClient(apiKey string) *ClaudeClient {
	return &ClaudeClient{
		apiKey:     apiKey,
		baseURL:    "https://api.anthropic.com/v1",
		model:      "claude-3-sonnet-20240229",
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *ClaudeClient) Solve(ctx context.Context, problem string) (string, error) {
	url := c.baseURL + "/messages"

	reqBody := ClaudeRequest{
		Model: c.model,
		Messages: []ClaudeMessage{
			{Role: "user", Content: "You are a math tutor for BAC C exam students in Mauritania. Provide clear, step-by-step solutions.\n\n" + problem},
		},
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("no response content")
	}

	var text strings.Builder
	for _, c := range claudeResp.Content {
		text.WriteString(c.Text)
	}

	slog.Info("claude solved", "tokens", claudeResp.Usage.TotalTokens)
	return text.String(), nil
}

type OpenAIClient struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey:     apiKey,
		baseURL:    "https://api.openai.com/v1",
		model:      "gpt-4o",
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *OpenAIClient) Solve(ctx context.Context, problem string) (string, error) {
	url := c.baseURL + "/chat/completions"

	reqBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a math tutor for BAC C exam students in Mauritania. Provide clear, step-by-step solutions."},
			{Role: "user", Content: problem},
		},
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d - %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	slog.Info("openai solved", "tokens", chatResp.Usage.TotalTokens)
	return chatResp.Choices[0].Message.Content, nil
}

type CloudSolver interface {
	Solve(ctx context.Context, problem string) (string, error)
}

func GetCloudClient(provider, apiKey string) CloudSolver {
	switch provider {
	case "deepseek":
		return NewDeepSeekClient(apiKey)
	case "grok":
		return NewGrokClient(apiKey)
	case "claude":
		return NewClaudeClient(apiKey)
	case "chatgpt":
		return NewOpenAIClient(apiKey)
	default:
		return NewDeepSeekClient(apiKey)
	}
}

func SolveWithAPI(ctx context.Context, provider, problem string) ResearchResult {
	apiKey := os.Getenv(strings.ToUpper(provider) + "_API_KEY")
	if apiKey == "" {
		envKey := provider + "_API_KEY"
		slog.Debug("API key not found", "provider", provider, "env", envKey)
		return ResearchResult{
			Provider: provider,
			Success:  false,
			Error:    "API key not configured. Set " + envKey,
		}
	}

	client := GetCloudClient(provider, apiKey)

	result, err := client.Solve(ctx, problem)
	if err != nil {
		return ResearchResult{
			Provider: provider,
			Success:  false,
			Error:    err.Error(),
		}
	}

	return ResearchResult{
		Sources:  1,
		Results:  result,
		Provider: provider,
		Success:  true,
	}
}

func SolveWithAnyAPI(ctx context.Context, problem string) ResearchResult {
	providers := []string{"deepseek", "grok", "claude", "chatgpt"}

	for _, provider := range providers {
		select {
		case <-ctx.Done():
			return ResearchResult{Success: false, Error: "cancelled"}
		default:
		}

		result := SolveWithAPI(ctx, provider, problem)
		if result.Success {
			return result
		}

		slog.Warn("API failed, trying next", "provider", provider, "error", result.Error)
	}

	return ResearchResult{
		Provider: "none",
		Success:  false,
		Error:    "all API providers failed",
	}
}
