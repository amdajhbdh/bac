package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type OpenStaxClient struct {
	client  *http.Client
	baseURL string
}

type OpenStaxBook struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Subject     string   `json:"subject"`
	Grades      []string `json:"grades"`
	CoverURL    string   `json:"cover_url"`
}

type OpenStaxSearchResult struct {
	Books      []OpenStaxBook `json:"items"`
	TotalCount int            `json:"total_count"`
}

func NewOpenStaxClient() *OpenStaxClient {
	return &OpenStaxClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://openstax.atlas.icloud.com/api/v1",
	}
}

func (c *OpenStaxClient) SearchBooks(ctx context.Context, query string, subject string) (*OpenStaxSearchResult, error) {
	slog.Info("searching OpenStax", "query", query, "subject", subject)

	url := fmt.Sprintf("%s/books?query=%s", c.baseURL, query)
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

	var result OpenStaxSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Info("OpenStax search complete", "found", len(result.Books))
	return &result, nil
}

func (c *OpenStaxClient) GetBook(ctx context.Context, bookID string) (*OpenStaxBook, error) {
	url := fmt.Sprintf("%s/books/%s", c.baseURL, bookID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var book OpenStaxBook
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &book, nil
}
