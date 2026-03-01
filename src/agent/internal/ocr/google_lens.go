package ocr

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type GoogleLensClient struct {
	client    *http.Client
	apiKey    string
	useVision bool
}

type VisionRequest struct {
	Image struct {
		Content string `json:"content"`
	} `json:"image"`
	Features []struct {
		Type       string `json:"type"`
		MaxResults int    `json:"maxResults"`
	} `json:"features"`
}

type VisionResponse struct {
	Responses []struct {
		TextAnnotations []struct {
			Description string `json:"description"`
			Locale      string `json:"locale"`
			Bounds      []struct {
				Vertices []struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"vertices"`
			} `json:"boundingPoly"`
		} `json:"textAnnotations"`
		FullTextAnnotation struct {
			Text string `json:"text"`
		} `json:"fullTextAnnotation"`
	} `json:"responses"`
}

type LensResult struct {
	Text       string
	Locale     string
	Confidence float64
}

func NewGoogleLensClient(apiKey string) *GoogleLensClient {
	return &GoogleLensClient{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
		apiKey:    apiKey,
		useVision: apiKey != "",
	}
}

func (c *GoogleLensClient) OCR(ctx context.Context, imageData []byte) (*LensResult, error) {
	if c.useVision {
		return c.ocrWithVision(ctx, imageData)
	}
	return c.ocrWithLensFallback(ctx, imageData)
}

func (c *GoogleLensClient) ocrWithVision(ctx context.Context, imageData []byte) (*LensResult, error) {
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	reqBody := VisionRequest{}
	reqBody.Image.Content = base64Image
	reqBody.Features = append(reqBody.Features, struct {
		Type       string `json:"type"`
		MaxResults int    `json:"maxResults"`
	}{Type: "TEXT_DETECTION", MaxResults: 10})
	reqBody.Features = append(reqBody.Features, struct {
		Type       string `json:"type"`
		MaxResults int    `json:"maxResults"`
	}{Type: "DOCUMENT_TEXT_DETECTION", MaxResults: 1})

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%s", c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	var visionResp VisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&visionResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(visionResp.Responses) == 0 {
		return &LensResult{}, nil
	}

	result := &LensResult{}

	if visionResp.Responses[0].FullTextAnnotation.Text != "" {
		result.Text = visionResp.Responses[0].FullTextAnnotation.Text
		result.Confidence = 0.95
	} else if len(visionResp.Responses[0].TextAnnotations) > 0 {
		var texts []string
		for _, ann := range visionResp.Responses[0].TextAnnotations {
			texts = append(texts, ann.Description)
			if ann.Locale != "" {
				result.Locale = ann.Locale
			}
		}
		result.Text = joinWithNewline(texts)
		result.Confidence = 0.8
	}

	slog.Info("Google Vision OCR complete", "text_length", len(result.Text), "confidence", result.Confidence)
	return result, nil
}

func (c *GoogleLensClient) ocrWithLensFallback(ctx context.Context, imageData []byte) (*LensResult, error) {
	slog.Info("using Google Lens fallback - requires API key for production")

	result := &LensResult{
		Text:       "Google Lens fallback - configure GOOGLE_CLOUD_VISION_API_KEY for production OCR",
		Confidence: 0.0,
	}

	return result, nil
}

func (c *GoogleLensClient) DetectMath(ctx context.Context, imageData []byte) ([]*MathDetection, error) {
	if !c.useVision {
		return nil, fmt.Errorf("math detection requires Google Cloud Vision API key")
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	reqBody := map[string]interface{}{
		"image": map[string]string{"content": base64Image},
		"features": []map[string]interface{}{
			{"type": "TEXT_DETECTION"},
			{"type": "LOGO_DETECTION", "maxResults": 5},
		},
	}

	jsonBody, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%s", c.apiKey)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var visionResp VisionResponse
	json.NewDecoder(resp.Body).Decode(&visionResp)

	if len(visionResp.Responses) == 0 {
		return nil, nil
	}

	var detections []*MathDetection
	for _, ann := range visionResp.Responses[0].TextAnnotations {
		if isMathExpression(ann.Description) {
			detections = append(detections, &MathDetection{
				Text: ann.Description,
			})
		}
	}

	return detections, nil
}

type MathDetection struct {
	Text    string
	Formula string
	Bounds  []Point
}

type Point struct {
	X, Y int
}

func isMathExpression(text string) bool {
	mathIndicators := []string{"=", "+", "-", "×", "÷", "∫", "∑", "√", "π", "θ", "x²", "x^"}
	for _, ind := range mathIndicators {
		if contains(text, ind) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func joinWithNewline(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += "\n"
		}
		result += s
	}
	return result
}

func GetGoogleLensConfig() (string, string) {
	apiKey := os.Getenv("GOOGLE_CLOUD_VISION_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_LENS_API_KEY")
	}
	return apiKey, os.Getenv("GOOGLE_LENS_ENDPOINT")
}

func NewGoogleLensFromEnv() *GoogleLensClient {
	apiKey, _ := GetGoogleLensConfig()
	return NewGoogleLensClient(apiKey)
}
