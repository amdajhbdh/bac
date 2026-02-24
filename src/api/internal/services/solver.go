package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type SolverService struct {
	ollamaURL string
	cloudURL  string
	noonPath  string
}

type SolveRequest struct {
	Question      string `json:"question"`
	Subject       string `json:"subject"`
	GenerateVideo bool   `json:"generate_video"`
}

type SolveResponse struct {
	ID             uuid.UUID `json:"id"`
	Question       string    `json:"question"`
	Solution       string    `json:"solution"`
	SolutionLatex  string    `json:"solution_latex"`
	Steps          []string  `json:"steps"`
	AnswerOptions  []string  `json:"answer_options"`
	CorrectAnswer  string    `json:"correct_answer"`
	Concepts       []string  `json:"concepts"`
	Difficulty     int       `json:"difficulty"`
	AnimationVideo string    `json:"animation_video,omitempty"`
}

func NewSolverService() *SolverService {
	return &SolverService{
		ollamaURL: os.Getenv("OLLAMA_URL"),
		cloudURL:  os.Getenv("CLOUD_API_URL"),
		noonPath:  "../../src/noon",
	}
}

func (s *SolverService) Solve(req SolveRequest) (*SolveResponse, error) {
	// Try local Ollama first
	resp, err := s.solveWithOllama(req.Question)
	if err == nil {
		// Generate noon animation if requested
		if req.GenerateVideo {
			videoURL, err := s.generateNoonAnimation(resp.Steps)
			if err == nil {
				resp.AnimationVideo = videoURL
			}
		}
		return resp, nil
	}

	// Fallback to cloud API
	return s.solveWithCloud(req)
}

func (s *SolverService) solveWithOllama(question string) (*SolveResponse, error) {
	prompt := fmt.Sprintf(`You are an expert BAC examiner. 
Given the following question, provide:
1. Step-by-step solution
2. All possible answer choices (A, B, C, D)
3. Correct answer
4. Key concepts involved
5. Difficulty rating (1-5)

Format your response as JSON:
{
  "solution": "...",
  "solution_latex": "...",
  "steps": ["step1", "step2", ...],
  "answer_options": ["A: ...", "B: ...", "C: ...", "D: ..."],
  "correct_answer": "A",
  "concepts": ["concept1", "concept2"],
  "difficulty": 3
}

Question: %s`, question)

	body, _ := json.Marshal(map[string]string{
		"model":  "llama3.2:3b",
		"prompt": prompt,
		"stream": "false",
	})

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Parse the response
	responseText := result["response"].(string)
	var solveResp SolveResponse
	if err := json.Unmarshal([]byte(responseText), &solveResp); err != nil {
		// Return raw response if JSON parsing fails
		solveResp = SolveResponse{
			ID:         uuid.New(),
			Question:   question,
			Solution:   responseText,
			Difficulty: 3,
		}
	}

	return &solveResp, nil
}

func (s *SolverService) solveWithCloud(req SolveRequest) (*SolveResponse, error) {
	// Placeholder for cloud API fallback
	return nil, fmt.Errorf("cloud solver not implemented")
}

func (s *SolverService) generateNoonAnimation(steps []string) (string, error) {
	// Generate noon code from solution steps
	noonCode := s.generateNoonCode(steps)

	// Save the code
	filename := fmt.Sprintf("animation_%d.go", time.Now().Unix())
	os.WriteFile(s.noonPath+"/examples/"+filename, []byte(noonCode), 0644)

	// Build and run
	// In production, this would be queued and processed asynchronously

	return "/animations/" + filename + ".mp4", nil
}

func (s *SolverService) generateNoonCode(steps []string) string {
	var buf bytes.Buffer
	buf.WriteString(`// Auto-generated noon animation
package main

import (
	"noon"
	"noon/ease"
)

func main() {
	noon.Run(func(scene *noon.Scene) {
`)

	for i, step := range steps {
		buf.WriteString(fmt.Sprintf(`
		// Step %d: %s
		text%d := scene.Text("%s").Make()
		scene.Play(text%d.ShowCreation()).RunTime(1.0)
		scene.Wait(0.5)
`, i+1, step, i+1, step, i+1))
	}

	buf.WriteString(`
	})
}
`)
	return buf.String()
}

func (s *SolverService) GetSteps(id uuid.UUID) (*SolveResponse, error) {
	// Fetch from cache/database
	return nil, nil
}

func (s *SolverService) GetAnimation(id uuid.UUID) (string, error) {
	return "", nil
}

// SolveWithImage processes an image question
func (s *SolverService) SolveWithImage(imageData []byte, subject string) (*SolveResponse, error) {
	// OCR would be done by submission service first
	// Then call solve with extracted text
	return nil, nil
}

func (c *SolverService) makeRequest(url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	return client.Post(url, "application/json", body)
}
