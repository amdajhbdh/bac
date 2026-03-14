package dedup

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/bac-unified/agent/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Deduplicator struct {
	pool      *pgxpool.Pool
	queries   *db.Queries
	threshold float64
}

func NewDeduplicator(pool *pgxpool.Pool) *Deduplicator {
	return &Deduplicator{
		pool:      pool,
		queries:   db.New(pool),
		threshold: 0.85,
	}
}

type Duplicate struct {
	ExistingID   string
	Similarity   float64
	ExistingText string
}

func (d *Deduplicator) CheckDuplicate(ctx context.Context, questionText string) (*Duplicate, error) {
	// First check exact hash match
	hash := hashText(questionText)
	exact, err := d.checkExactMatch(ctx, hash)
	if err != nil {
		return nil, err
	}
	if exact != nil {
		return exact, nil
	}

	// Then check fuzzy similarity
	fuzzy, err := d.checkFuzzyMatch(ctx, questionText)
	if err != nil {
		return nil, err
	}

	return fuzzy, nil
}

func (d *Deduplicator) checkExactMatch(ctx context.Context, hash string) (*Duplicate, error) {
	rows, err := d.pool.Query(ctx, `
		SELECT id, question_text, 1.0 as similarity 
		FROM questions 
		WHERE question_hash = $1
		LIMIT 1
	`, hash)
	if err != nil {
		return nil, fmt.Errorf("exact match query: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var id string
		var text string
		var sim float64
		if err := rows.Scan(&id, &text, &sim); err != nil {
			return nil, err
		}
		return &Duplicate{
			ExistingID:   id,
			Similarity:   sim,
			ExistingText: text,
		}, nil
	}

	return nil, nil
}

func (d *Deduplicator) checkFuzzyMatch(ctx context.Context, questionText string) (*Duplicate, error) {
	// Get recent questions for comparison
	rows, err := d.pool.Query(ctx, `
		SELECT id, question_text 
		FROM questions 
		WHERE created_at > NOW() - INTERVAL '30 days'
		LIMIT 100
	`)
	if err != nil {
		return nil, fmt.Errorf("fuzzy match query: %w", err)
	}
	defer rows.Close()

	var bestMatch *Duplicate
	normalizedNew := normalizeText(questionText)

	for rows.Next() {
		var id string
		var existingText string
		if err := rows.Scan(&id, &existingText); err != nil {
			continue
		}

		normalizedExisting := normalizeText(existingText)
		sim := similarity(normalizedNew, normalizedExisting)

		if sim >= d.threshold && (bestMatch == nil || sim > bestMatch.Similarity) {
			bestMatch = &Duplicate{
				ExistingID:   id,
				Similarity:   sim,
				ExistingText: existingText,
			}
		}
	}

	return bestMatch, nil
}

func (d *Deduplicator) MarkAsDuplicate(ctx context.Context, originalID, duplicateID string) error {
	_, err := d.pool.Exec(ctx, `
		UPDATE questions 
		SET verification_status = 'duplicate', 
		    parent_question_id = $1
		WHERE id = $2
	`, originalID, duplicateID)
	return err
}

func (d *Deduplicator) SetThreshold(threshold float64) {
	d.threshold = threshold
}

func hashText(text string) string {
	normalized := normalizeText(text)
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

func normalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	// Remove extra whitespace
	text = strings.Join(strings.Fields(text), " ")
	// Remove common math notation variations
	text = strings.ReplaceAll(text, "×", "*")
	text = strings.ReplaceAll(text, "÷", "/")
	text = strings.ReplaceAll(text, "−", "-")
	return text
}

func similarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	// Simple word-based similarity
	wordsA := strings.Fields(a)
	wordsB := strings.Fields(b)

	if len(wordsA) == 0 || len(wordsB) == 0 {
		return 0.0
	}

	matchCount := 0
	for _, wA := range wordsA {
		for _, wB := range wordsB {
			if wA == wB {
				matchCount++
				break
			}
		}
	}

	// Jaccard-like similarity
	union := len(wordsA) + len(wordsB) - matchCount
	if union == 0 {
		return 0.0
	}

	return float64(matchCount) / float64(union)
}
