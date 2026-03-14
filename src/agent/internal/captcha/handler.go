package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Solver interface {
	Solve(ctx context.Context, imageURL string) (string, error)
	CanSolve(captchaType string) bool
}

type Config struct {
	APIKey     string
	Provider   string
	Timeout    time.Duration
	MaxRetries int
}

type Result struct {
	Solution string
	ID       string
	Solved   bool
}

type CaptchaHandler struct {
	client  *http.Client
	config  Config
	solvers []Solver
}

func NewCaptchaHandler(config Config) *CaptchaHandler {
	h := &CaptchaHandler{
		client:  &http.Client{Timeout: config.Timeout},
		config:  config,
		solvers: []Solver{},
	}

	if config.APIKey != "" {
		h.solvers = append(h.solvers, NewTwoCaptchaSolver(config.APIKey))
	}

	h.solvers = append(h.solvers, &BasicImageSolver{})

	return h
}

func (h *CaptchaHandler) DetectCaptcha(pageText string) string {
	indicators := []struct {
		pattern string
		captcha string
	}{
		{"captcha", "image"},
		{"recaptcha", "recaptcha"},
		{"hcaptcha", "hcaptcha"},
		{"turnstile", "cloudflare"},
		{"verify you are human", "generic"},
		{"i'm not a robot", "recaptcha"},
		{"selenium", "detection"},
	}

	pageText = strings.ToLower(pageText)
	for _, ind := range indicators {
		if strings.Contains(pageText, ind.pattern) {
			return ind.captcha
		}
	}

	return ""
}

func (h *CaptchaHandler) Solve(ctx context.Context, captchaType, imageURL string) (*Result, error) {
	slog.Info("solving captcha", "type", captchaType)

	for _, solver := range h.solvers {
		if solver.CanSolve(captchaType) {
			solution, err := solver.Solve(ctx, imageURL)
			if err != nil {
				slog.Warn("solver failed", "solver", fmt.Sprintf("%T", solver), "error", err)
				continue
			}
			return &Result{
				Solution: solution,
				Solved:   true,
			}, nil
		}
	}

	return &Result{Solved: false}, fmt.Errorf("no solver available for %s", captchaType)
}

type BasicImageSolver struct{}

func (s *BasicImageSolver) CanSolve(captchaType string) bool {
	return captchaType == "image" || captchaType == ""
}

func (s *BasicImageSolver) Solve(ctx context.Context, imageURL string) (string, error) {
	return "", fmt.Errorf("basic solver cannot solve captchas - requires API key")
}

type TwoCaptchaSolver struct {
	apiKey string
	client *http.Client
}

func NewTwoCaptchaSolver(apiKey string) *TwoCaptchaSolver {
	return &TwoCaptchaSolver{
		apiKey: apiKey,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (s *TwoCaptchaSolver) CanSolve(captchaType string) bool {
	return s.apiKey != ""
}

func (s *TwoCaptchaSolver) Solve(ctx context.Context, imageURL string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("no API key configured")
	}

	submitURL := fmt.Sprintf("http://2captcha.com/in.php?key=%s&method=userrecaptcha&googlekey=%s&pageurl=%s",
		s.apiKey, "", imageURL)

	resp, err := s.client.Get(submitURL)
	if err != nil {
		return "", fmt.Errorf("submit failed: %w", err)
	}
	defer resp.Body.Close()

	var submitResp struct {
		Status  string `json:"status"`
		Request string `json:"request"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&submitResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if submitResp.Status != "1" {
		return "", fmt.Errorf("submit failed: %s", submitResp.Request)
	}

	for i := 0; i < 30; i++ {
		time.Sleep(5 * time.Second)

		resultURL := fmt.Sprintf("http://2captcha.com/res.php?key=%s&action=get&id=%s",
			s.apiKey, submitResp.Request)

		resp, err := s.client.Get(resultURL)
		if err != nil {
			continue
		}

		var resultResp struct {
			Status  string `json:"status"`
			Request string `json:"request"`
		}
		json.NewDecoder(resp.Body).Decode(&resultResp)
		resp.Body.Close()

		if resultResp.Status == "1" {
			return resultResp.Request, nil
		}
	}

	return "", fmt.Errorf("captcha solving timeout")
}

func DetectAndWait(pageText string, waitFunc func() error) error {
	indicators := []string{"captcha", "recaptcha", "hcaptcha", "verify"}
	lower := strings.ToLower(pageText)

	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			slog.Info("captcha detected, waiting for user", "indicator", ind)
			return waitFunc()
		}
	}

	return nil
}

func GetEnvConfig() Config {
	return Config{
		APIKey:     os.Getenv("CAPTCHA_API_KEY"),
		Provider:   os.Getenv("CAPTCHA_PROVIDER"),
		Timeout:    120 * time.Second,
		MaxRetries: 3,
	}
}
