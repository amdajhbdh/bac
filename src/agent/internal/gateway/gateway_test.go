package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type GatewayRequest struct {
	Message string `json:"message"`
	Mode    string `json:"mode,omitempty"`
}

type GatewayResponse struct {
	Message   string   `json:"message"`
	Mode      string   `json:"mode"`
	Sources   []Source `json:"sources"`
	SessionID string   `json:"session_id"`
}

type Source struct {
	Text       string  `json:"text"`
	Similarity float64 `json:"similarity"`
	SourceType string  `json:"source_type"`
}

func TestGatewayChatIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat" {
			t.Errorf("Expected /chat, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		var req GatewayRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Message != "Hello" {
			t.Errorf("Expected Hello, got %s", req.Message)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GatewayResponse{
			Message:   "Hi there!",
			Mode:      "Chat",
			Sources:   []Source{},
			SessionID: "test-session",
		})
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	reqBody := GatewayRequest{Message: "Hello", Mode: "Chat"}
	body, _ := json.Marshal(reqBody)

	resp, err := client.Post(server.URL+"/chat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to call gateway: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	var gatewayResp GatewayResponse
	json.NewDecoder(resp.Body).Decode(&gatewayResp)

	if gatewayResp.Message != "Hi there!" {
		t.Errorf("Expected Hi there!, got %s", gatewayResp.Message)
	}

	if gatewayResp.Mode != "Chat" {
		t.Errorf("Expected Chat mode, got %s", gatewayResp.Mode)
	}
}

func TestGatewayOCRIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ocr" {
			t.Errorf("Expected /ocr, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"text":               "Test OCR result",
				"confidence":         0.95,
				"engine":             "tesseract",
				"processing_time_ms": 100,
			},
		})
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	formData := []byte(`file=test`)
	resp, err := client.Post(server.URL+"/ocr", "application/octet-stream", bytes.NewBuffer(formData))
	if err != nil {
		t.Fatalf("Failed to call OCR: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestGatewayAnimationIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/animation" {
			t.Errorf("Expected /animation, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"job_id": "test-job-123",
			"status": "queued",
		})
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	animationReq := map[string]string{
		"code":    "from manim import *\nclass Test(Scene):\n    pass",
		"quality": "low",
	}
	body, _ := json.Marshal(animationReq)

	resp, err := client.Post(server.URL+"/animation", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to call animation: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestGatewayHealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected /health, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":          "healthy",
			"gateway_version": "0.1.0",
		})
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to call health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
