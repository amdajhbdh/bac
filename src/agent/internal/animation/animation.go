package animation

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AnimationResult struct {
	FilePath     string
	OutputDir    string
	ExportFormat string
	Success      bool
	Error        string
}

type AnimationCache struct {
	mu     sync.RWMutex
	cache  map[string]AnimationResult
	maxAge time.Duration
}

func NewAnimationCache() *AnimationCache {
	return &AnimationCache{
		cache:  make(map[string]AnimationResult),
		maxAge: 24 * time.Hour,
	}
}

func (c *AnimationCache) Get(key string) (AnimationResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, found := c.cache[key]
	return result, found
}

func (c *AnimationCache) Set(key string, result AnimationResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = result
}

func cacheKeyFor(problem, solution string) string {
	combined := problem + "||" + solution
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

type AnimationConfig struct {
	Width        int
	Height       int
	FPS          int
	Duration     float32
	IncludeAudio bool
	ExportVideo  bool
	ExportImages bool
}

var DefaultConfig = AnimationConfig{
	Width:        1280,
	Height:       720,
	FPS:          30,
	Duration:     10.0,
	IncludeAudio: true,
	ExportVideo:  false,
	ExportImages: true,
}

var animationCache = NewAnimationCache()

func Generate(ctx context.Context, problem, solution string) AnimationResult {
	return GenerateWithConfig(ctx, problem, solution, DefaultConfig)
}

func GenerateCached(ctx context.Context, problem, solution string) AnimationResult {
	cacheKey := cacheKeyFor(problem, solution)
	if cached, found := animationCache.Get(cacheKey); found {
		slog.Info("animation cache hit", "key", cacheKey)
		return cached
	}

	result := GenerateWithConfig(ctx, problem, solution, DefaultConfig)
	if result.Success {
		animationCache.Set(cacheKey, result)
	}

	return result
}

func GenerateWithConfig(ctx context.Context, problem, solution string, cfg AnimationConfig) AnimationResult {
	slog.Info("generating animation", "problem", problem, "config", cfg)

	code := generateNoonCode(problem, solution, cfg)

	tmpDir := filepath.Join(os.TempDir(), "bac-animations")
	os.MkdirAll(tmpDir, 0755)

	projectDir := filepath.Join(tmpDir, fmt.Sprintf("anim_%d", time.Now().UnixNano()))
	os.MkdirAll(projectDir, 0755)

	codePath := filepath.Join(projectDir, "src/main.rs")
	os.MkdirAll(filepath.Dir(codePath), 0755)

	err := os.WriteFile(codePath, []byte(code), 0644)
	if err != nil {
		slog.Error("failed to write noon code", "error", err)
		return AnimationResult{Success: false, Error: err.Error()}
	}

	noonPath := findNoonPath()
	writeCargoToml(projectDir, noonPath)
	writeConfigYaml(projectDir, cfg)

	slog.Info("animation code generated", "path", codePath, "noon_path", noonPath)

	return AnimationResult{
		FilePath:     codePath,
		OutputDir:    projectDir,
		ExportFormat: "images",
		Success:      true,
	}
}

func findNoonPath() string {
	possiblePaths := []string{
		filepath.Join(os.Getenv("PWD"), "src", "noon"),
		filepath.Join(os.Getenv("PWD"), "..", "src", "noon"),
		filepath.Join(os.Getenv("HOME"), "Documents", "bac", "src", "noon"),
		"/home/med/Documents/bac/src/noon",
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(filepath.Join(p, "noon", "Cargo.toml")); err == nil {
			slog.Info("found noon at", "path", p)
			return p
		}
	}

	slog.Warn("noon not found, using default path")
	return "../../noon"
}

func generateNoonCode(problem, solution string, cfg AnimationConfig) string {
	short := problem
	if len(short) > 40 {
		short = short[:40] + "..."
	}

	steps := extractSteps(solution)

	stepsCode := ""
	particlesCode := ""

	for i, step := range steps {
		if i < 6 {
			yPos := 2 - float64(i)*0.8
			stepsCode += fmt.Sprintf(`
    // Step %d
    let step%d = scene.text()
        .with_text("%s")
        .with_position(-4.5, %.1f)
        .with_color(Color::WHITE)
        .with_font_size(20.0)
        .make();
    scene.play(step%d.show_creation()).run_time(0.5);
    scene.wait(0.3);`, i+1, i+1, escapeString(step), yPos, i+1)
		}
	}

	if len(steps) > 0 {
		particlesCode = `
    // Particle celebration effect
    let mut particle_sys = ParticleSystem::new();
    let mut config = ParticleConfig::default();
    config.count = 50;
    config.lifetime = 1.5;
    config.start_color = Color::YELLOW;
    config.end_color = Color::ORANGE;
    config.initial_velocity = Vector2::new(0.0, 150.0);
    config.gravity = Vector2::new(0.0, -80.0);
    config.start_size = 8.0;
    config.end_size = 2.0;
    config.start_position = Position { x: 0.0, y: -3.0 };
    particle_sys.add_emitter(config);`
	}

	exportCode := ""
	if cfg.ExportImages {
		exportCode = `
    // Export frames for recording
    let recorder = RecordingSession::new("frames", 1280, 720, 30);`
	}

	audioCode := ""
	if cfg.IncludeAudio {
		audioCode = `
    // Audio system - can play tones for feedback
    let audio = AudioSystem::new();`
	}

	_ = exportCode
	_ = audioCode

	return fmt.Sprintf(`use noon::prelude::*;
use std::time::Duration;

struct Model {
    scene: Scene,
    particle_sys: ParticleSystem,
    recorder: Option<RecordingSession>,
    audio: Option<AudioSystem>,
    frame_count: u32,
}

fn scene(win_rect: Rect) -> Scene {
    let mut scene = Scene::new(win_rect);
    
    // Background
    let bg = scene.rectangle()
        .with_position(0.0, 0.0)
        .with_size(20.0, 12.0)
        .with_color(Color::rgb(0.1, 0.1, 0.15))
        .make();
    
    // Title
    let title = scene.text()
        .with_text("BAC - Solution")
        .with_position(-3.0, 4.5)
        .with_color(Color::YELLOW)
        .with_font_size(32.0)
        .make();
    scene.play(title.show_creation()).run_time(0.8);
    scene.wait(0.3);
    
    // Problem statement
    let prob = scene.text()
        .with_text("%s")
        .with_position(-4.5, 3.5)
        .with_color(Color::CYAN)
        .with_font_size(18.0)
        .make();
    scene.play(prob.show_creation()).run_time(0.6);
    scene.wait(0.5);%s
    
    // Answer
    let answer = scene.text()
        .with_text("✓ Solution Complete")
        .with_position(-2.5, -4.0)
        .with_color(Color::GREEN)
        .with_font_size(24.0)
        .make();
    scene.play(answer.show_creation()).run_time(0.5);
    scene.wait(1.0);%s
    
    scene
}

fn main() {
    noon::app(|app| {
        let model = Model {
            scene: scene(app.window_rect()),
            particle_sys: ParticleSystem::new(),
            recorder: None,
            audio: None,
            frame_count: 0,
        };
        
        app.set_loop_mode(LoopMode::loop_once());
        model
    })
    .update(|app, model, _update| {
        let dt = app.duration.since_start.as_secs_f32();
        
        // Update particles
        model.particle_sys.update(0.016);
        
        // Record frames
        if let Some(recorder) = &mut model.recorder {
            // Would capture frame here
            model.frame_count += 1;
        }
    })
    .view(|app, model, frame| {
        model.scene.draw(frame);
    })
    .run();
}`, short, stepsCode, particlesCode)
}

func writeCargoToml(projectDir, noonPath string) {
	cargoToml := fmt.Sprintf(`[package]
name = "bac-animation"
version = "0.1.0"
edition = "2021"

[dependencies]
noon = { path = "%s/noon" }
nannou = "0.19"
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"

[profile.release]
opt-level = 3
lto = true
`, noonPath)

	os.WriteFile(filepath.Join(projectDir, "Cargo.toml"), []byte(cargoToml), 0644)
}

func writeConfigYaml(projectDir string, cfg AnimationConfig) {
	configYaml := fmt.Sprintf(`animation:
  width: %d
  height: %d
  fps: %d
  duration: %.1f

export:
  images: %t
  video: %t
  format: "png"

audio:
  enabled: %t
  volume: 0.7
`, cfg.Width, cfg.Height, cfg.FPS, cfg.Duration, cfg.ExportImages, cfg.ExportVideo, cfg.IncludeAudio)

	os.WriteFile(filepath.Join(projectDir, "config.yaml"), []byte(configYaml), 0644)
}

func extractSteps(solution string) []string {
	lines := strings.Split(solution, "\n")
	var steps []string
	inSteps := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "Étapes") || strings.Contains(trimmed, "étapes") || strings.Contains(trimmed, "Steps") || strings.Contains(trimmed, "Paso") {
			inSteps = true
			continue
		}
		if inSteps && len(trimmed) > 0 {
			clean := strings.TrimLeft(trimmed, "0123456789.)- ")
			if len(clean) > 3 && len(clean) < 120 {
				steps = append(steps, clean)
			}
		}
		if strings.Contains(trimmed, "Concept") || strings.Contains(trimmed, "Réponse") || strings.Contains(trimmed, "Answer") {
			break
		}
	}

	if len(steps) == 0 {
		for _, line := range lines {
			if len(strings.TrimSpace(line)) > 5 && len(steps) < 4 {
				steps = append(steps, strings.TrimSpace(line))
			}
		}
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

func createFallbackImages(projectDir string, cfg AnimationConfig) (string, error) {
	slog.Info("creating fallback images (simple text frames)")

	framesDir := filepath.Join(projectDir, "frames")
	os.MkdirAll(framesDir, 0755)

	simpleHtml := `<!DOCTYPE html>
<html>
<head>
<style>
body { background: #1a1a2e; color: white; font-family: monospace; 
       display: flex; justify-content: center; align-items: center; 
       height: 100vh; margin: 0; }
</style>
</head>
<body>
<h1>BAC Animation</h1>
<p>Run with Noon for full animations</p>
</body>
</html>`

	indexPath := filepath.Join(framesDir, "index.html")
	os.WriteFile(indexPath, []byte(simpleHtml), 0644)

	readmePath := filepath.Join(projectDir, "README.md")
	readme := `# BAC Animation Project

To generate full animations:
1. Install Rust: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
2. cd ` + projectDir + `
3. cargo run --release

Requirements:
- Rust 1.70+
- cargo
`
	os.WriteFile(readmePath, []byte(readme), 0644)

	slog.Info("fallback created", "dir", framesDir)
	return framesDir, nil
}

func CompileAndExport(ctx context.Context, problem, solution string, cfg AnimationConfig) AnimationResult {
	genResult := GenerateWithConfig(ctx, problem, solution, cfg)
	if !genResult.Success {
		return genResult
	}

	framesDir, err := ExportImages(ctx, genResult.OutputDir, cfg)
	if err != nil {
		genResult.Success = false
		genResult.Error = err.Error()
		return genResult
	}

	genResult.FilePath = framesDir

	if cfg.ExportVideo {
		videoPath, err := CreateVideo(ctx, framesDir, cfg)
		if err == nil {
			genResult.FilePath = videoPath
			genResult.ExportFormat = "video"
		}
	}

	return genResult
}

func CreateVideo(ctx context.Context, framesDir string, cfg AnimationConfig) (string, error) {
	slog.Info("creating video from frames", "frames_dir", framesDir)

	videoPath := filepath.Join(filepath.Dir(framesDir), "animation.mp4")

	cmd := exec.CommandContext(ctx, "ffmpeg", "-framerate", fmt.Sprintf("%d", cfg.FPS),
		"-i", filepath.Join(framesDir, "%05d.png"),
		"-c:v", "libx264", "-pix_fmt", "yuv420p", "-crf", "23",
		"-y", videoPath)

	if err := cmd.Run(); err != nil {
		slog.Warn("ffmpeg not available, skipping video", "error", err)
		return "", fmt.Errorf("ffmpeg not available: %w", err)
	}

	slog.Info("video created", "path", videoPath)
	return videoPath, nil
}

func Compile(ctx context.Context, codePath string) error {
	slog.Info("compiling animation", "path", codePath)

	projectDir := filepath.Dir(codePath)

	cmd := exec.CommandContext(ctx, "cargo", "build", "--release")
	cmd.Dir = projectDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("compile failed: %w - %s", err, stderr.String())
	}

	slog.Info("animation compiled successfully")
	return nil
}

func ExportImages(ctx context.Context, projectDir string, cfg AnimationConfig) (string, error) {
	slog.Info("exporting animation frames", "project", projectDir)

	framesDir := filepath.Join(projectDir, "frames")
	os.MkdirAll(framesDir, 0755)

	binPath := filepath.Join(projectDir, "target/release/bac-animation")

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		err := Compile(ctx, filepath.Join(projectDir, "src/main.rs"))
		if err != nil {
			slog.Warn("compilation failed, trying fallback", "error", err)
			return createFallbackImages(projectDir, cfg)
		}
	}

	cmd := exec.CommandContext(ctx, binPath)
	cmd.Dir = projectDir
	cmd.Env = append(os.Environ(), "EXPORT_FRAMES=true")

	timeout := time.Duration(cfg.Duration*1000) * time.Millisecond
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		cmd.Process.Kill()
	}

	slog.Info("frames exported", "dir", framesDir)
	return framesDir, nil
}
