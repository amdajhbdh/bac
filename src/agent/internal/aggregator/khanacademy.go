package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type KhanAcademyClient struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

type KhanAcademyTopic struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	NodeSlug    string `json:"node_slug"`
}

type KhanAcademyVideo struct {
	ID          string `json:"youtube_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Duration    int    `json:"duration"`
}

type KhanAcademySearchResult struct {
	Topics []KhanAcademyTopic `json:"topics"`
	Videos []KhanAcademyVideo `json:"videos"`
}

func NewKhanAcademyClient(apiKey string) *KhanAcademyClient {
	return &KhanAcademyClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://www.khanacademy.org/api",
		apiKey:  apiKey,
	}
}

func (c *KhanAcademyClient) Search(ctx context.Context, query string, subject string) (*KhanAcademySearchResult, error) {
	slog.Info("searching Khan Academy", "query", query, "subject", subject)

	url := fmt.Sprintf("%s/videos?format=search&query=%s", c.baseURL, query)
	if subject != "" {
		url += "&subject=" + subject
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result KhanAcademySearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Info("Khan Academy search complete", "topics", len(result.Topics), "videos", len(result.Videos))
	return &result, nil
}

func (c *KhanAcademyClient) GetTopics(ctx context.Context, subject string) ([]KhanAcademyTopic, error) {
	url := fmt.Sprintf("%s/topics/%s", c.baseURL, subject)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var topics []KhanAcademyTopic
	if err := json.NewDecoder(resp.Body).Decode(&topics); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return topics, nil
}
