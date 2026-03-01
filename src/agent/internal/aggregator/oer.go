package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type OERClient struct {
	client  *http.Client
	baseURL string
}

type OERResource struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Subject     string `json:"subject"`
	GradeLevel  string `json:"grade_level"`
	Format      string `json:"format"`
	License     string `json:"license"`
}

type OERSearchResult struct {
	Resources []OERResource `json:"resources"`
	Total     int           `json:"total_count"`
	Page      int           `json:"page"`
}

func NewOERClient() *OERClient {
	return &OERClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://www.oercommons.org/api/v3",
	}
}

func (c *OERClient) Search(ctx context.Context, query string, subject string, limit int) (*OERSearchResult, error) {
	slog.Info("searching OER Commons", "query", query, "subject", subject)

	url := fmt.Sprintf("%s/search?q=%s&limit=%d", c.baseURL, query, limit)
	if subject != "" {
		url += "&subject=" + subject
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var result OERSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Info("OER search complete", "found", len(result.Resources))
	return &result, nil
}

func (c *OERClient) GetResource(ctx context.Context, id string) (*OERResource, error) {
	url := fmt.Sprintf("%s/resources/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var resource OERResource
	if err := json.NewDecoder(resp.Body).Decode(&resource); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &resource, nil
}
