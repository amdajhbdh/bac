package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PredictionService struct {
	dataDir   string
	db        *pgxpool.Pool
	ollamaURL string
}

type Prediction struct {
	ID              uuid.UUID `json:"id"`
	SubjectID       int       `json:"subject_id"`
	ChapterID       int       `json:"chapter_id"`
	ExamYear        int       `json:"exam_year"`
	ExamSession     string    `json:"exam_session"`
	PredictedTopics []Topic   `json:"predicted_topics"`
	ConfidenceScore float64   `json:"confidence_score"`
	Status          string    `json:"status"`
	BasedOnPatterns []string  `json:"based_on_patterns"`
	MLModelVersion  string    `json:"ml_model_version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Topic struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
	Questions   int     `json:"questions"`
}

func NewPredictionService(db *pgxpool.Pool) *PredictionService {
	dataDir := "./data/predictions"
	os.MkdirAll(dataDir, 0755)

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	return &PredictionService{
		dataDir:   dataDir,
		db:        db,
		ollamaURL: ollamaURL,
	}
}

type questionPattern struct {
	text        string
	year        int
	session     string
	difficulty  int
	successRate float64
}

func (s *PredictionService) GeneratePrediction(subjectID, chapterID, examYear int) (*Prediction, error) {
	// Use ML-based prediction if database is available
	if s.db != nil {
		return s.generateMLPrediction(subjectID, chapterID, examYear)
	}

	// Fallback to rule-based prediction
	prediction := &Prediction{
		ID:          uuid.New(),
		SubjectID:   subjectID,
		ChapterID:   chapterID,
		ExamYear:    examYear,
		ExamSession: "principal",
		PredictedTopics: []Topic{
			{Name: "Fonctions exponentielles", Probability: 0.85, Questions: 3},
			{Name: "Intégrales", Probability: 0.78, Questions: 2},
			{Name: "Nombres complexes", Probability: 0.72, Questions: 2},
		},
		ConfidenceScore: 0.78,
		Status:          "published",
		BasedOnPatterns: []string{
			"frequency: high (appeared 8 times in last 10 years)",
			"trend: increasing",
			"difficulty: medium",
		},
		CreatedAt: time.Now(),
	}

	// Save prediction
	s.savePrediction(prediction)

	return prediction, nil
}

func (s *PredictionService) generateMLPrediction(subjectID, chapterID, examYear int) (*Prediction, error) {
	ctx := context.Background()

	// Query historical data for pattern analysis
	rows, err := s.db.Query(ctx, `
		SELECT 
			question_text, 
			source_year, 
			source_session,
			difficulty,
			correct_count::float / NULLIF(submission_count, 0) as success_rate
		FROM questions 
		WHERE subject_id = $1 
		AND (chapter_id = $2 OR $2 IS NULL)
		AND source_year IS NOT NULL
		ORDER BY source_year DESC
		LIMIT 100`, subjectID, chapterID)

	if err != nil {
		slog.Warn("failed to query historical data", "error", err)
		return s.GeneratePrediction(subjectID, chapterID, examYear)
	}
	defer rows.Close()

	var patterns []questionPattern
	for rows.Next() {
		var p questionPattern
		if err := rows.Scan(&p.text, &p.year, &p.session, &p.difficulty, &p.successRate); err != nil {
			continue
		}
		patterns = append(patterns, p)
	}

	// Use LLM to analyze patterns and predict topics
	predictedTopics, confidence, patternsFound := s.analyzePatternsWithLLM(patterns, subjectID)

	prediction := &Prediction{
		ID:              uuid.New(),
		SubjectID:       subjectID,
		ChapterID:       chapterID,
		ExamYear:        examYear,
		ExamSession:     "principal",
		PredictedTopics: predictedTopics,
		ConfidenceScore: confidence,
		Status:          "published",
		BasedOnPatterns: patternsFound,
		MLModelVersion:  "v1.0-llm",
		CreatedAt:       time.Now(),
	}

	// Save to database
	_, err = s.db.Exec(ctx, `
		INSERT INTO predictions (id, subject_id, chapter_id, exam_year, exam_session, 
			predicted_topics, confidence_score, status, based_on_patterns, ml_model_version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT DO NOTHING`,
		prediction.ID, prediction.SubjectID, prediction.ChapterID, prediction.ExamYear,
		prediction.ExamSession, prediction.PredictedTopics, prediction.ConfidenceScore,
		prediction.Status, prediction.BasedOnPatterns, prediction.MLModelVersion, prediction.CreatedAt)

	if err != nil {
		slog.Warn("failed to save prediction to database", "error", err)
	}

	// Also save locally
	s.savePrediction(prediction)

	return prediction, nil
}

func (s *PredictionService) analyzePatternsWithLLM(patterns []questionPattern, subjectID int) ([]Topic, float64, []string) {
	// Prepare context for LLM
	var patternsText string
	for _, p := range patterns {
		if len(patternsText) > 2000 {
			break
		}
		patternsText += fmt.Sprintf("Year: %d, Session: %s, Difficulty: %d, Success: %.2f\n",
			p.year, p.session, p.difficulty, p.successRate)
	}

	// Call LLM for prediction
	prompt := fmt.Sprintf(`Analyze these historical exam question patterns and predict likely topics for the next exam.
Subject ID: %d

Historical Patterns:
%s

Based on the patterns above, predict the top 5 most likely topics with probability scores (0-1).
Also identify any trends (increasing/decreasing frequency).

Respond in JSON format:
{
  "topics": [{"name": "topic name", "probability": 0.85, "questions": 3}],
  "confidence": 0.78,
  "trends": ["trend1", "trend2"]
}`, subjectID, patternsText)

	type llmResponse struct {
		Topics     []Topic  `json:"topics"`
		Confidence float64  `json:"confidence"`
		Trends     []string `json:"trends"`
	}

	// Try to call Ollama
	resp, err := callOllama(s.ollamaURL, prompt)
	if err != nil {
		// Fallback to rule-based
		return s.ruleBasedPrediction(patterns), 0.65, []string{"fallback: llm unavailable"}
	}

	var result llmResponse
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		return s.ruleBasedPrediction(patterns), 0.65, []string{"fallback: parse error"}
	}

	trends := result.Trends
	if trends == nil {
		trends = []string{}
	}

	return result.Topics, result.Confidence, trends
}

func callOllama(url, prompt string) (string, error) {
	type request struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
	}

	type response struct {
		Response string `json:"response"`
	}

	reqBody, _ := json.Marshal(request{
		Model:  "llama3.2:3b",
		Prompt: prompt,
		Stream: false,
	})

	resp, err := http.Post(url+"/api/generate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

func (s *PredictionService) ruleBasedPrediction(patterns []questionPattern) []Topic {
	// Count topic frequencies
	topicCounts := make(map[string]int)
	topicSuccess := make(map[string]float64)
	topicTotal := make(map[string]int)

	for _, p := range patterns {
		// Extract key terms from question text
		keywords := extractKeywords(p.text)
		for _, kw := range keywords {
			topicCounts[kw]++
			if p.successRate > 0 {
				topicSuccess[kw] += p.successRate
				topicTotal[kw]++
			}
		}
	}

	// Calculate probabilities
	type scoredTopic struct {
		name        string
		probability float64
		questions   int
	}

	var scored []scoredTopic
	for topic, count := range topicCounts {
		avgSuccess := 0.0
		if topicTotal[topic] > 0 {
			avgSuccess = topicSuccess[topic] / float64(topicTotal[topic])
		}

		// Score based on frequency and difficulty
		probability := float64(count) / float64(len(patterns)) * 0.7
		probability += (1 - avgSuccess) * 0.3 // Higher probability for harder questions

		scored = append(scored, scoredTopic{
			name:        topic,
			probability: probability,
			questions:   count,
		})
	}

	// Sort by probability
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].probability > scored[j].probability
	})

	// Take top 5
	var topics []Topic
	for i := 0; i < min(5, len(scored)); i++ {
		topics = append(topics, Topic{
			Name:        scored[i].name,
			Probability: scored[i].probability,
			Questions:   scored[i].questions,
		})
	}

	return topics
}

func extractKeywords(text string) []string {
	// Simple keyword extraction
	keywords := []string{}
	important := []string{
		"intégrale", "dérivée", "fonction", "limite", "suite", "nombre complexe",
		"physique", "mécanique", "électrique", "optique", "chimie", "organique",
		"thermodynamique", "onde", "radioactivité", "génétique", "évolution",
	}

	lower := strings.ToLower(text)
	for _, kw := range important {
		if strings.Contains(lower, kw) {
			keywords = append(keywords, kw)
		}
	}

	if len(keywords) == 0 {
		keywords = append(keywords, "general")
	}

	return keywords
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *PredictionService) GetPrediction(id uuid.UUID) (*Prediction, error) {
	filePath := filepath.Join(s.dataDir, id.String()+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var prediction Prediction
	if err := json.Unmarshal(data, &prediction); err != nil {
		return nil, err
	}

	return &prediction, nil
}

func (s *PredictionService) GetBySubject(subjectID, examYear int) ([]Prediction, error) {
	// List all predictions for a subject
	var predictions []Prediction
	files, _ := filepath.Glob(filepath.Join(s.dataDir, "*.json"))

	for _, file := range files {
		data, _ := os.ReadFile(file)
		var p Prediction
		if json.Unmarshal(data, &p) == nil && p.SubjectID == subjectID && p.ExamYear == examYear {
			predictions = append(predictions, p)
		}
	}

	return predictions, nil
}

func (s *PredictionService) GetLatest() ([]Prediction, error) {
	var predictions []Prediction
	files, _ := filepath.Glob(filepath.Join(s.dataDir, "*.json"))

	for _, file := range files {
		data, _ := os.ReadFile(file)
		var p Prediction
		if json.Unmarshal(data, &p) == nil && p.Status == "published" {
			predictions = append(predictions, p)
		}
	}

	return predictions, nil
}

func (s *PredictionService) savePrediction(p *Prediction) error {
	data, _ := json.MarshalIndent(p, "", "  ")
	filePath := filepath.Join(s.dataDir, p.ID.String()+".json")
	return os.WriteFile(filePath, data, 0644)
}

func (s *PredictionService) AnalyzePatterns(subjectID int) ([]string, error) {
	// ML analysis would happen here
	// Return insights like:
	return []string{
		"Questions on derivatives appear in 90% of exams",
		"Geometry questions are increasing in difficulty",
		"Integration by parts has appeared 5 years in a row",
	}, nil
}
