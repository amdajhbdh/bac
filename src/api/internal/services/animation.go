package services

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type AnimationService struct {
	noonPath  string
	outputDir string
}

type AnimationRequest struct {
	Problem     string   `json:"problem"`
	Solution    string   `json:"solution"`
	Subject     string   `json:"subject"`
	Steps       []string `json:"steps"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	FPS         int      `json:"fps"`
	Duration    float32  `json:"duration"`
	ExportVideo bool     `json:"export_video"`
}

type AnimationResult struct {
	Success  bool    `json:"success"`
	FilePath string  `json:"file_path"`
	Format   string  `json:"format"`
	Duration float32 `json:"duration"`
	Error    string  `json:"error,omitempty"`
}

func NewAnimationService() *AnimationService {
	noonPath := os.Getenv("NOON_PATH")
	if noonPath == "" {
		noonPath = "../../src/noon"
	}

	outputDir := os.Getenv("ANIMATION_OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "./animations"
	}
	os.MkdirAll(outputDir, 0755)

	return &AnimationService{
		noonPath:  noonPath,
		outputDir: outputDir,
	}
}

func (s *AnimationService) Generate(req AnimationRequest) *AnimationResult {
	// Set defaults
	if req.Width == 0 {
		req.Width = 1280
	}
	if req.Height == 0 {
		req.Height = 720
	}
	if req.FPS == 0 {
		req.FPS = 30
	}
	if req.Duration == 0 {
		req.Duration = 10.0
	}

	slog.Info("generating animation", "problem", req.Problem[:min(50, len(req.Problem))])

	// Generate Noon code
	code := s.generateNoonCode(req)

	// Create project directory
	projectDir := filepath.Join(s.outputDir, fmt.Sprintf("anim_%d", time.Now().UnixNano()))
	os.MkdirAll(filepath.Join(projectDir, "src"), 0755)

	// Write code
	codePath := filepath.Join(projectDir, "src/main.rs")
	if err := os.WriteFile(codePath, []byte(code), 0644); err != nil {
		return &AnimationResult{Success: false, Error: fmt.Sprintf("write code: %v", err)}
	}

	// Write Cargo.toml
	cargoToml := s.generateCargoToml()
	if err := os.WriteFile(filepath.Join(projectDir, "Cargo.toml"), []byte(cargoToml), 0644); err != nil {
		return &AnimationResult{Success: false, Error: fmt.Sprintf("write cargo: %v", err)}
	}

	// Try to build
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := s.build(ctx, projectDir); err != nil {
		slog.Warn("build failed, creating fallback", "error", err)
		return s.createFallback(req, projectDir)
	}

	// Run to generate frames
	if err := s.run(ctx, projectDir); err != nil {
		slog.Warn("run failed", "error", err)
	}

	// Export video if requested
	if req.ExportVideo {
		videoPath, err := s.createVideo(ctx, projectDir, req.FPS)
		if err == nil {
			return &AnimationResult{
				Success:  true,
				FilePath: videoPath,
				Format:   "video",
				Duration: req.Duration,
			}
		}
	}

	// Return frames directory
	return &AnimationResult{
		Success:  true,
		FilePath: projectDir,
		Format:   "frames",
		Duration: req.Duration,
	}
}

func (s *AnimationService) generateNoonCode(req AnimationRequest) string {
	// Extract key steps for animation
	steps := req.Steps
	if len(steps) == 0 {
		// Parse from solution text
		steps = parseSolutionSteps(req.Solution)
	}

	// Limit steps for animation
	if len(steps) > 6 {
		steps = steps[:6]
	}

	// Generate step animations
	stepsCode := ""
	for i, step := range steps {
		yPos := 2.5 - float64(i)*0.9
		cleanStep := escapeString(step)
		stepsCode += fmt.Sprintf(`
    // Step %d
    let step%d = scene.text()
        .with_text("%s")
        .with_position(-4.5, %.1f)
        .with_color(Color::WHITE)
        .with_font_size(24.0)
        .make();
    scene.play(step%d.show_creation()).run_time(0.6);
    scene.wait(0.4);`, i+1, i+1, cleanStep, yPos, i+1)
	}

	// Get subject for title
	subject := req.Subject
	if subject == "" {
		subject = "Solution"
	}

	// Subject color
	subjectColor := "Color::CYAN"
	switch strings.ToLower(subject) {
	case "math", "mathématiques":
		subjectColor = "Color::YELLOW"
	case "physique":
		subjectColor = "Color::RGB(0.2, 0.6, 1.0)"
	case "chimie":
		subjectColor = "Color::RGB(0.2, 0.9, 0.5)"
	}

	return fmt.Sprintf(`use noon::prelude::*;
use std::time::Duration;

fn main() {
    noon::app(|app| {
        let model = Model {
            scene: create_scene(app.window_rect()),
        };
        
        app.set_loop_mode(LoopMode::loop_once());
        model
    })
    .update(|_app, _model, _update| {})
    .view(|app, model, frame| {
        model.scene.draw(frame);
    })
    .run();
}

struct Model {
    scene: Scene,
}

fn create_scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Background with gradient effect
    let bg = scene.rectangle()
        .with_position(0.0, 0.0)
        .with_size(20.0, 12.0)
        .with_color(Color::rgb(0.08, 0.08, 0.12))
        .make();
    
    // Title
    let title = scene.text()
        .with_text("BAC - %s")
        .with_position(-3.5, 4.5)
        .with_color(%s)
        .with_font_size(36.0)
        .make();
    scene.play(title.show_creation()).run_time(0.8);
    scene.wait(0.3);
    
    // Problem statement
    let problem = scene.text()
        .with_text("%s")
        .with_position(-4.5, 3.5)
        .with_color(Color::rgb(0.9, 0.9, 0.95))
        .with_font_size(20.0)
        .make();
    scene.play(problem.show_creation()).run_time(0.6);
    scene.wait(0.5);%s
    
    // Final answer
    let answer = scene.text()
        .with_text("✓ Solution Complete")
        .with_position(-2.5, -4.0)
        .with_color(Color::rgb(0.2, 0.9, 0.4))
        .with_font_size(28.0)
        .make();
    scene.play(answer.show_creation()).run_time(0.5);
    scene.wait(1.0);
    
    // Fade out all
    scene.play(scene.fade_out()).run_time(0.5);
    
    scene
}
`, subject, subjectColor, escapeString(req.Problem[:min(100, len(req.Problem))]), stepsCode)
}

func (s *AnimationService) generateCargoToml() string {
	return `[package]
name = "bac-animation"
version = "0.1.0"
edition = "2021"

[dependencies]
noon = { path = "` + s.noonPath + `/noon" }
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

[profile.release]
opt-level = 3
lto = true
`
}

func (s *AnimationService) build(ctx context.Context, projectDir string) error {
	slog.Info("building animation", "dir", projectDir)

	cmd := exec.CommandContext(ctx, "cargo", "build", "--release")
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		slog.Debug("build output", "stdout", stdout.String(), "stderr", stderr.String())
		return fmt.Errorf("build failed: %w - %s", err, stderr.String())
	}

	return nil
}

func (s *AnimationService) run(ctx context.Context, projectDir string) error {
	slog.Info("running animation", "dir", projectDir)

	binPath := filepath.Join(projectDir, "target/release/bac-animation")

	cmd := exec.CommandContext(ctx, binPath)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "RENDER_FRAMES=true")

	timeout := 15 * time.Second
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		cmd.Process.Kill()
		return nil // Timeout is OK - we just need frames
	}
}

func (s *AnimationService) createVideo(ctx context.Context, projectDir string, fps int) (string, error) {
	slog.Info("creating video", "dir", projectDir)

	framesDir := filepath.Join(projectDir, "frames")
	if _, err := os.Stat(framesDir); os.IsNotExist(err) {
		// Try alternate location
		framesDir = projectDir
	}

	videoPath := filepath.Join(projectDir, "animation.mp4")

	// Check if ffmpeg is available
	cmd := exec.CommandContext(ctx, "ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg not available: %w", err)
	}

	// Create video from frames
	ffmpegCmd := exec.CommandContext(ctx, "ffmpeg",
		"-framerate", fmt.Sprintf("%d", fps),
		"-i", filepath.Join(framesDir, "%05d.png"),
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-crf", "23",
		"-y", videoPath)

	if err := ffmpegCmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w", err)
	}

	return videoPath, nil
}

func (s *AnimationService) createFallback(req AnimationRequest, projectDir string) *AnimationResult {
	// Create HTML fallback
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>BAC Animation - %s</title>
    <style>
        body { 
            background: linear-gradient(135deg, #1a1a2e 0%%, #16213e 100%%);
            color: white; 
            font-family: 'Courier New', monospace;
            padding: 40px;
            min-height: 100vh;
        }
        .container { max-width: 900px; margin: 0 auto; }
        h1 { color: #f1c40f; text-align: center; }
        .problem { 
            background: rgba(255,255,255,0.1); 
            padding: 20px; 
            border-radius: 10px;
            margin: 20px 0;
        }
        .steps { 
            background: rgba(0,0,0,0.3); 
            padding: 20px; 
            border-radius: 10px;
        }
        .step { 
            margin: 15px 0; 
            padding: 10px;
            border-left: 3px solid #3498db;
            background: rgba(52, 152, 219, 0.1);
        }
        .answer {
            text-align: center;
            color: #2ecc71;
            font-size: 1.5em;
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>BAC - %s</h1>
        <div class="problem">
            <h2>Problème:</h2>
            <p>%s</p>
        </div>
        <div class="steps">
            <h2>Solution:</h2>
            %s
        </div>
        <div class="answer">✓ Solution Complete</div>
    </div>
</body>
</html>`, req.Subject, req.Subject, escapeHTML(req.Problem), formatStepsHTML(req.Steps))

	htmlPath := filepath.Join(projectDir, "animation.html")
	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		return &AnimationResult{Success: false, Error: err.Error()}
	}

	return &AnimationResult{
		Success:  true,
		FilePath: htmlPath,
		Format:   "html",
		Duration: float32(len(req.Steps)) * 2.0,
	}
}

func parseSolutionSteps(solution string) []string {
	var steps []string
	lines := strings.Split(solution, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 10 && len(trimmed) < 200 {
			// Remove common prefixes
			clean := strings.TrimLeft(trimmed, "0123456789.)→>- ")
			if len(clean) > 5 {
				steps = append(steps, clean)
			}
		}
	}

	if len(steps) == 0 && len(solution) > 0 {
		// If no clear steps, just use the whole solution
		steps = append(steps, solution[:min(200, len(solution))])
	}

	return steps
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	return s
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

func formatStepsHTML(steps []string) string {
	var html string
	for i, step := range steps {
		html += fmt.Sprintf(`<div class="step"><strong>Étape %d:</strong> %s</div>`, i+1, escapeHTML(step))
	}
	return html
}
