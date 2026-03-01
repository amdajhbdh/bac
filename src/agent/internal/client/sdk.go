package bac

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Config struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

type Client struct {
	config      Config
	httpClient  *http.Client
	auth        *AuthService
	questions   *QuestionsService
	solver      *SolverService
	predictions *PredictionsService
}

func NewClient(config Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	if config.BaseURL == "" {
		config.BaseURL = os.Getenv("BAC_API_URL")
	}
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080/api/v1"
	}

	c := &Client{
		config:     config,
		httpClient: config.HTTPClient,
	}

	c.auth = &AuthService{client: c}
	c.questions = &QuestionsService{client: c}
	c.solver = &SolverService{client: c}
	c.predictions = &PredictionsService{client: c}

	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.config.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(data))
	}

	return data, nil
}

type AuthService struct {
	client *Client
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	User      User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Points   int    `json:"points"`
	Level    int    `json:"level"`
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	data, err := s.client.doRequest(ctx, "POST", "/auth/register", req)
	if err != nil {
		return nil, err
	}

	var resp AuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	data, err := s.client.doRequest(ctx, "POST", "/auth/login", req)
	if err != nil {
		return nil, err
	}

	var resp AuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

type QuestionsService struct {
	client *Client
}

type Question struct {
	ID           string `json:"id"`
	QuestionText string `json:"question_text"`
	SolutionText string `json:"solution_text,omitempty"`
	Subject      string `json:"subject"`
	Chapter      string `json:"chapter"`
	Difficulty   int    `json:"difficulty"`
}

type ListQuestionsResponse struct {
	Questions []Question `json:"questions"`
	Total     int        `json:"total"`
}

func (s *QuestionsService) List(ctx context.Context, page, limit int) (*ListQuestionsResponse, error) {
	path := fmt.Sprintf("/questions?page=%d&limit=%d", page, limit)
	data, err := s.client.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp ListQuestionsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

func (s *QuestionsService) Get(ctx context.Context, id string) (*Question, error) {
	path := fmt.Sprintf("/questions/%s", id)
	data, err := s.client.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp Question
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

func (s *QuestionsService) Create(ctx context.Context, q Question) (*Question, error) {
	data, err := s.client.doRequest(ctx, "POST", "/questions", q)
	if err != nil {
		return nil, err
	}

	var resp Question
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

type SolverService struct {
	client *Client
}

type SolveRequest struct {
	Problem string `json:"problem"`
}

type SolveResponse struct {
	ID         string   `json:"id"`
	Solution   string   `json:"solution"`
	Steps      int      `json:"steps"`
	Confidence float64  `json:"confidence"`
	Subject    string   `json:"subject"`
	Chapter    string   `json:"chapter"`
	Concepts   []string `json:"concepts"`
	Model      string   `json:"model"`
}

func (s *SolverService) Solve(ctx context.Context, req SolveRequest) (*SolveResponse, error) {
	data, err := s.client.doRequest(ctx, "POST", "/solve", req)
	if err != nil {
		return nil, err
	}

	var resp SolveResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

func (s *SolverService) GetSteps(ctx context.Context, id string) ([]string, error) {
	path := fmt.Sprintf("/solve/%s/steps", id)
	data, err := s.client.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var steps []string
	if err := json.Unmarshal(data, &steps); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return steps, nil
}

type PredictionsService struct {
	client *Client
}

type Prediction struct {
	ID              string  `json:"id"`
	Subject         string  `json:"subject"`
	Chapter         string  `json:"chapter"`
	Topic           string  `json:"topic"`
	Probability     float64 `json:"probability"`
	ConfidenceScore float64 `json:"confidence_score"`
	Year            int     `json:"year"`
}

func (s *PredictionsService) List(ctx context.Context) ([]Prediction, error) {
	data, err := s.client.doRequest(ctx, "GET", "/predictions", nil)
	if err != nil {
		return nil, err
	}

	var predictions []Prediction
	if err := json.Unmarshal(data, &predictions); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return predictions, nil
}

func (s *PredictionsService) GetBySubject(ctx context.Context, subject string) ([]Prediction, error) {
	path := fmt.Sprintf("/predictions/subject/%s", subject)
	data, err := s.client.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var predictions []Prediction
	if err := json.Unmarshal(data, &predictions); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return predictions, nil
}

func init() {
	slog.Info("BAC client SDK initialized")
}
