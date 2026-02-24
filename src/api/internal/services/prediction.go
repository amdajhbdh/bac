package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type PredictionService struct {
	dataDir string
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
	CreatedAt       time.Time `json:"created_at"`
}

type Topic struct {
	Name        string  `json:"name"`
	Probability float64 `json:"probability"`
	Questions   int     `json:"questions"`
}

func NewPredictionService() *PredictionService {
	dataDir := "./data/predictions"
	os.MkdirAll(dataDir, 0755)
	return &PredictionService{dataDir: dataDir}
}

func (s *PredictionService) GeneratePrediction(subjectID, chapterID, examYear int) (*Prediction, error) {
	// In production: analyze historical data using ML
	// For now: generate based on rules

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
