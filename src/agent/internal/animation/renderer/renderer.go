package renderer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bac-unified/agent/internal/animation/bridge"
)

type QualityPreset struct {
	Name       string
	Width      int
	Height     int
	FPS        int
	Bitrate    string
	Codec      string
	RenderArgs []string
}

var QualityPresets = map[string]QualityPreset{
	"preview": {
		Name:       "preview",
		Width:      854,
		Height:     480,
		FPS:        15,
		Bitrate:    "1M",
		Codec:      "libx264",
		RenderArgs: []string{"-ql"},
	},
	"low": {
		Name:       "low",
		Width:      854,
		Height:     480,
		FPS:        30,
		Bitrate:    "2M",
		Codec:      "libx264",
		RenderArgs: []string{"-ql"},
	},
	"medium": {
		Name:       "medium",
		Width:      1280,
		Height:     720,
		FPS:        30,
		Bitrate:    "5M",
		Codec:      "libx264",
		RenderArgs: []string{"-qm"},
	},
	"high": {
		Name:       "high",
		Width:      1920,
		Height:     1080,
		FPS:        60,
		Bitrate:    "10M",
		Codec:      "libx264",
		RenderArgs: []string{"-qh"},
	},
	"production": {
		Name:       "production",
		Width:      2560,
		Height:     1440,
		FPS:        60,
		Bitrate:    "20M",
		Codec:      "libx265",
		RenderArgs: []string{"-qp"},
	},
}

const DefaultQuality = "medium"

type RenderRequest struct {
	SceneCode string
	Quality   string
	Width     int
	Height    int
	FPS       int
	Timeout   int
}

type RenderResult struct {
	Success    bool
	VideoPath  string
	FramesDir  string
	Duration   float32
	Quality    string
	Resolution string
	HasAudio   bool
	Error      string
	RenderTime time.Duration
}

type Service struct {
	bridge *bridge.Renderer
}

func NewService() *Service {
	return &Service{
		bridge: bridge.NewRenderer(),
	}
}

func (s *Service) Render(ctx context.Context, req RenderRequest) (*RenderResult, error) {
	startTime := time.Now()

	// Validate and get quality preset
	quality := req.Quality
	if quality == "" {
		quality = DefaultQuality
	}

	preset, exists := QualityPresets[quality]
	if !exists {
		slog.Warn("unknown quality preset, using default", "requested", quality, "default", DefaultQuality)
		preset = QualityPresets[DefaultQuality]
		quality = DefaultQuality
	}

	// Override with custom values if provided
	width := preset.Width
	height := preset.Height
	fps := preset.FPS

	if req.Width > 0 {
		width = req.Width
	}
	if req.Height > 0 {
		height = req.Height
	}
	if req.FPS > 0 {
		fps = req.FPS
	}

	slog.Info("rendering animation",
		"quality", quality,
		"resolution", fmt.Sprintf("%dx%d", width, height),
		"fps", fps)

	// Convert to bridge request
	bridgeReq := bridge.RenderRequest{
		SceneCode: req.SceneCode,
		Quality:   quality,
		Width:     width,
		Height:    height,
		FPS:       fps,
		Timeout:   req.Timeout,
	}

	// Render
	resp, err := s.bridge.Render(ctx, bridgeReq)
	if err != nil {
		return &RenderResult{
			Success:    false,
			Error:      fmt.Sprintf("render error: %v", err),
			Quality:    quality,
			RenderTime: time.Since(startTime),
		}, nil
	}

	if !resp.Success {
		return &RenderResult{
			Success:    false,
			Error:      resp.Error,
			Quality:    quality,
			RenderTime: time.Since(startTime),
		}, nil
	}

	result := &RenderResult{
		Success:    true,
		VideoPath:  resp.VideoPath,
		FramesDir:  resp.FramesDir,
		Quality:    quality,
		Resolution: fmt.Sprintf("%dx%d", width, height),
		RenderTime: time.Since(startTime),
	}

	slog.Info("render complete",
		"video_path", result.VideoPath,
		"duration", result.RenderTime)

	return result, nil
}

func (s *Service) ValidateSetup(ctx context.Context) error {
	return s.bridge.ValidateSetup(ctx)
}

func (s *Service) Test(ctx context.Context) (*RenderResult, error) {
	slog.Info("running test render")

	// Simple test scene
	testCode := `from manim import *

class AnimatedScene(Scene):
    def construct(self):
        circle = Circle()
        square = Square()
        self.play(Create(circle))
        self.play(Transform(circle, square))
`

	req := RenderRequest{
		SceneCode: testCode,
		Quality:   "preview",
		Timeout:   120,
	}

	return s.Render(ctx, req)
}

func (s *Service) GetPreset(name string) (QualityPreset, bool) {
	preset, ok := QualityPresets[name]
	return preset, ok
}

func (s *Service) ListPresets() []QualityPreset {
	presets := make([]QualityPreset, 0, len(QualityPresets))
	for _, preset := range QualityPresets {
		presets = append(presets, preset)
	}
	return presets
}
