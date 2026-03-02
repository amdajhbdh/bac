package bridge

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type RenderRequest struct {
	SceneCode string `json:"scene_code"`
	Quality   string `json:"quality"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	FPS       int    `json:"fps,omitempty"`
	Timeout   int    `json:"timeout,omitempty"`
}

type RenderResponse struct {
	Success   bool   `json:"success"`
	VideoPath string `json:"video_path,omitempty"`
	FramesDir string `json:"frames_dir,omitempty"`
	Error     string `json:"error,omitempty"`
	Logs      string `json:"logs,omitempty"`
}

type Renderer struct {
	containerImage string
	containerName  string
	timeout        time.Duration
	workDir        string
}

func NewRenderer() *Renderer {
	containerImage := os.Getenv("MANIM_CONTAINER_IMAGE")
	if containerImage == "" {
		containerImage = "docker.io/manimcommunity/manim:latest"
	}

	containerName := os.Getenv("MANIM_CONTAINER_NAME")
	if containerName == "" {
		containerName = "bac-animator"
	}

	timeout := 5 * time.Minute
	if t := os.Getenv("ANIMATION_TIMEOUT"); t != "" {
		if seconds, err := time.ParseDuration(t + "s"); err == nil {
			timeout = seconds
		}
	}

	workDir := os.Getenv("ANIMATION_OUTPUT_DIR")
	if workDir == "" {
		workDir = "/tmp/bac-animations"
	}

	return &Renderer{
		containerImage: containerImage,
		containerName:  containerName,
		timeout:        timeout,
		workDir:        workDir,
	}
}

func (r *Renderer) Render(ctx context.Context, req RenderRequest) (*RenderResponse, error) {
	slog.Info("rendering animation with Podman",
		"image", r.containerImage,
		"quality", req.Quality,
		"timeout", r.timeout)

	// Ensure work directory exists
	os.MkdirAll(r.workDir, 0755)

	// Write scene code to file
	sceneFile := filepath.Join(r.workDir, "scene.py")
	if err := os.WriteFile(sceneFile, []byte(req.SceneCode), 0644); err != nil {
		return nil, fmt.Errorf("write scene file: %w", err)
	}

	// Get quality flags
	qualityFlags := getQualityFlags(req.Quality)

	// Build podman command
	cmd := r.buildPodmanCommand(sceneFile, qualityFlags, req.Timeout)

	// Run command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &RenderResponse{
			Success: false,
			Error:   fmt.Sprintf("render failed: %s - %s", err, string(output)),
			Logs:    string(output),
		}, nil
	}

	// Find output video
	videoPath := r.findOutputVideo(req.Quality)

	return &RenderResponse{
		Success:   true,
		VideoPath: videoPath,
		FramesDir: r.workDir,
		Logs:      string(output),
	}, nil
}

func (r *Renderer) buildPodmanCommand(sceneFile string, qualityFlags []string, timeout int) *exec.Cmd {
	timeoutSeconds := int(r.timeout.Seconds())
	if timeout > 0 {
		timeoutSeconds = timeout
	}

	workDirAbs, _ := filepath.Abs(r.workDir)

	runner := "podman"
	if _, err := exec.LookPath("podman"); err != nil {
		runner = "docker"
		slog.Warn("podman not found, falling back to docker")
	}

	containerName := fmt.Sprintf("%s-%d", r.containerName, time.Now().UnixNano())

	args := []string{
		"run",
		"--rm",
		"--name", containerName,
		"--network", "none",
		"-v", workDirAbs + ":/workspace:z",
		"-w", "/workspace",
		"--ulimits", "memlock=512000:512000",
		"-e", "MANIM_PREVIEW=0",
		r.containerImage,
		"manim",
	}
	args = append(args, qualityFlags...)
	args = append(args, []string{
		filepath.Base(sceneFile),
		"AnimatedScene",
		"--output_file", "/workspace/output",
		"--disable_caching",
	}...)

	cmd := exec.Command(runner, args...)
	_ = timeoutSeconds // Used for future timeout handling
	return cmd
}

func getQualityFlags(quality string) []string {
	flags := map[string][]string{
		"preview":    {"-ql"},
		"low":        {"-ql"},
		"medium":     {"-qm"},
		"high":       {"-qh"},
		"production": {"-qp"},
	}
	if f, ok := flags[quality]; ok {
		return f
	}
	return flags["medium"]
}

func (r *Renderer) findOutputVideo(quality string) string {
	qualityDir := map[string]string{
		"preview":    "480p15",
		"low":        "480p30",
		"medium":     "720p30",
		"high":       "1080p60",
		"production": "1440p60",
	}

	dir := qualityDir[quality]
	if dir == "" {
		dir = "720p30"
	}

	possiblePaths := []string{
		filepath.Join(r.workDir, "media", "videos", "scene", dir, "output.mp4"),
		filepath.Join(r.workDir, "media", "videos", "scene", "480p15", "output.mp4"),
		filepath.Join(r.workDir, "output.mp4"),
		filepath.Join(r.workDir, "scene.mp4"),
	}

	for _, p := range possiblePaths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return possiblePaths[0]
}

func (r *Renderer) Test(ctx context.Context) (*RenderResponse, error) {
	slog.Info("running test render with Podman")

	testCode := `from manim import *

class AnimatedScene(Scene):
    def construct(self):
        circle = Circle()
        self.play(Create(circle))
        self.play(Transform(circle, Square()))
`

	req := RenderRequest{
		SceneCode: testCode,
		Quality:   "preview",
		Timeout:   120,
	}

	return r.Render(ctx, req)
}

func (r *Renderer) ValidateSetup(ctx context.Context) error {
	runner := "podman"
	if _, err := exec.LookPath("podman"); err != nil {
		runner = "docker"
		if _, err := exec.LookPath("docker"); err != nil {
			return fmt.Errorf("neither podman nor docker found")
		}
	}
	_ = runner // Runner found, validation passed

	os.MkdirAll(r.workDir, 0755)
	testFile := filepath.Join(r.workDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("work directory not writable: %w", err)
	}
	os.Remove(testFile)

	return nil
}

func (r *Renderer) PullImage(ctx context.Context) error {
	runner := "podman"
	if _, err := exec.LookPath("podman"); err != nil {
		runner = "docker"
	}

	slog.Info("pulling container image", "image", r.containerImage, "runner", runner)

	cmd := exec.CommandContext(ctx, runner, "pull", r.containerImage)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pull failed: %w - output: %s", err, string(output))
	}

	slog.Info("image pulled successfully")
	return nil
}
