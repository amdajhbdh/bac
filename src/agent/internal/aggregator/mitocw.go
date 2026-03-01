package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type MITOCWClient struct {
	client  *http.Client
	baseURL string
}

type MITOCWCourse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Department  string `json:"department"`
	Level       string `json:"level"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail"`
}

type MITOCWSearchResult struct {
	Courses []MITOCWCourse `json:"courses"`
	Total   int            `json:"total_count"`
}

func NewMITOCWClient() *MITOCWClient {
	return &MITOCWClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://ocw.mit.edu/api",
	}
}

func (c *MITOCWClient) Search(ctx context.Context, query string, department string) (*MITOCWSearchResult, error) {
	slog.Info("searching MIT OCW", "query", query, "department", department)

	url := fmt.Sprintf("%s/courses", c.baseURL)
	if query != "" {
		url += "?q=" + query
	}
	if department != "" {
		if query != "" {
			url += "&"
		} else {
			url += "?"
		}
		url += "department=" + department
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

	var result MITOCWSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Info("MIT OCW search complete", "found", len(result.Courses))
	return &result, nil
}

func (c *MITOCWClient) GetCourse(ctx context.Context, courseID string) (*MITOCWCourse, error) {
	url := fmt.Sprintf("%s/courses/%s", c.baseURL, courseID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var course MITOCWCourse
	if err := json.NewDecoder(resp.Body).Decode(&course); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &course, nil
}
