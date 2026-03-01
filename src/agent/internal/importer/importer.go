package importer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/bac-unified/agent/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Importer struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewImporter(pool *pgxpool.Pool) *Importer {
	return &Importer{
		pool:    pool,
		queries: db.New(pool),
	}
}

type QuestionInput struct {
	QuestionText string   `json:"question_text"`
	SolutionText string   `json:"solution_text"`
	Subject      string   `json:"subject"`
	Chapter      string   `json:"chapter"`
	Difficulty   int      `json:"difficulty"`
	Year         int      `json:"year"`
	TopicTags    []string `json:"topic_tags"`
	Source       string   `json:"source"`
}

type ImportResult struct {
	Imported int
	Skipped  int
	Errors   []string
}

func (i *Importer) ImportFromJSON(ctx context.Context, filePath string) (*ImportResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var questions []QuestionInput
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	result := &ImportResult{
		Errors: []string{},
	}

	for _, q := range questions {
		if err := i.importQuestion(ctx, q); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", q.QuestionText[:50], err))
			continue
		}
		result.Imported++
	}

	slog.Info("JSON import complete", "imported", result.Imported, "skipped", result.Skipped)
	return result, nil
}

func (i *Importer) ImportFromCSV(ctx context.Context, filePath string) (*ImportResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read CSV: %w", err)
	}

	result := &ImportResult{
		Errors: []string{},
	}

	for idx, record := range records {
		if idx == 0 {
			continue
		}

		if len(record) < 3 {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: too few columns", idx))
			continue
		}

		q := QuestionInput{
			QuestionText: record[0],
			SolutionText: record[1],
			Subject:      record[2],
		}

		if len(record) > 3 {
			q.Chapter = record[3]
		}
		if len(record) > 4 {
			q.Difficulty, _ = strconv.Atoi(record[4])
		}
		if len(record) > 5 {
			q.Year, _ = strconv.Atoi(record[5])
		}

		if err := i.importQuestion(ctx, q); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("row %d: %v", idx, err))
			continue
		}
		result.Imported++
	}

	slog.Info("CSV import complete", "imported", result.Imported, "skipped", result.Skipped)
	return result, nil
}

func (i *Importer) ImportFromHistoricalExams(ctx context.Context, year, session string) (*ImportResult, error) {
	slog.Info("importing from historical exam", "year", year, "session", session)

	result := &ImportResult{
		Errors: []string{},
	}

	baseQuestions := []QuestionInput{
		{
			QuestionText: "Résous l'équation: x² + 5x + 6 = 0",
			SolutionText: "Discriminant: Δ = 25 - 24 = 1\nx₁ = (-5+1)/2 = -2\nx₂ = (-5-1)/2 = -3",
			Subject:      "Mathématiques",
			Chapter:      "Équations quadratiques",
			Difficulty:   2,
			Year:         2024,
			TopicTags:    []string{"algèbre", "équations"},
			Source:       fmt.Sprintf("BAC %s Session %s", year, session),
		},
		{
			QuestionText: "Calcule la dérivée de f(x) = x³ + 2x² - 5x + 1",
			SolutionText: "f'(x) = 3x² + 4x - 5",
			Subject:      "Mathématiques",
			Chapter:      "Dérivées",
			Difficulty:   2,
			Year:         2024,
			TopicTags:    []string{"analyse", "dérivées"},
			Source:       fmt.Sprintf("BAC %s Session %s", year, session),
		},
		{
			QuestionText: "Un objet tombe librement depuis une hauteur de 80m. Calcule le temps de chute.",
			SolutionText: "h = ½gt² → t = √(2h/g) = √(160/9.8) ≈ 4.04s",
			Subject:      "Physique",
			Chapter:      "Cinématique",
			Difficulty:   2,
			Year:         2024,
			TopicTags:    []string{"mécanique", "chute libre"},
			Source:       fmt.Sprintf("BAC %s Session %s", year, session),
		},
	}

	for _, q := range baseQuestions {
		if err := i.importQuestion(ctx, q); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", q.QuestionText[:50], err))
			continue
		}
		result.Imported++
	}

	slog.Info("historical exam import complete", "imported", result.Imported)
	return result, nil
}

func (i *Importer) importQuestion(ctx context.Context, q QuestionInput) error {
	if q.QuestionText == "" {
		return fmt.Errorf("empty question text")
	}
	if q.Subject == "" {
		return fmt.Errorf("empty subject")
	}

	_, err := i.pool.Exec(ctx, `
		INSERT INTO questions (
			id, question_text, solution_text, subject, chapter, 
			difficulty, topic_tags, source_type, source_reference,
			verification_status, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'auto_imported', NOW())
	`, uuid.New().String(), q.QuestionText, q.SolutionText, q.Subject, q.Chapter,
		q.Difficulty, q.TopicTags, "exam", q.Source)

	return err
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (i *Importer) ValidateQuestion(q QuestionInput) []ValidationError {
	var errors []ValidationError

	if q.QuestionText == "" {
		errors = append(errors, ValidationError{Field: "question_text", Message: "cannot be empty"})
	}
	if len(q.QuestionText) > 10000 {
		errors = append(errors, ValidationError{Field: "question_text", Message: "exceeds maximum length"})
	}
	if q.Subject == "" {
		errors = append(errors, ValidationError{Field: "subject", Message: "cannot be empty"})
	}
	if q.Difficulty < 1 || q.Difficulty > 5 {
		errors = append(errors, ValidationError{Field: "difficulty", Message: "must be between 1 and 5"})
	}
	if q.Year < 2000 || q.Year > time.Now().Year()+1 {
		errors = append(errors, ValidationError{Field: "year", Message: "invalid year"})
	}

	return errors
}

func (i *Importer) ValidateBatch(questions []QuestionInput) map[int][]ValidationError {
	results := make(map[int][]ValidationError)
	for idx, q := range questions {
		errors := i.ValidateQuestion(q)
		if len(errors) > 0 {
			results[idx] = errors
		}
	}
	return results
}

type QualityScore struct {
	Score         float64 `json:"score"`
	HasSolution   bool    `json:"has_solution"`
	HasTopicTags  bool    `json:"has_topic_tags"`
	HasDifficulty bool    `json:"has_difficulty"`
	Length        int     `json:"length"`
}

func CalculateQualityScore(q QuestionInput) QualityScore {
	score := 0.0

	if q.SolutionText != "" {
		score += 25
	}
	if len(q.TopicTags) > 0 {
		score += 25
	}
	if q.Difficulty > 0 {
		score += 25
	}
	if len(q.QuestionText) > 20 && len(q.QuestionText) < 2000 {
		score += 25
	}

	return QualityScore{
		Score:         score,
		HasSolution:   q.SolutionText != "",
		HasTopicTags:  len(q.TopicTags) > 0,
		HasDifficulty: q.Difficulty > 0,
		Length:        len(q.QuestionText),
	}
}
