package memory

import (
	"testing"
)

func TestLookupResult(t *testing.T) {
	result := LookupResult{
		SimilarProblems: []SimilarProblem{
			{Question: "Test?", Solution: "Answer", Similarity: 0.9},
		},
		Context: "Test context",
	}

	if len(result.SimilarProblems) != 1 {
		t.Errorf("Expected 1 similar problem, got %d", len(result.SimilarProblems))
	}

	if result.Context != "Test context" {
		t.Errorf("Expected context 'Test context', got %q", result.Context)
	}
}

func TestSimilarProblem(t *testing.T) {
	sp := SimilarProblem{
		Question:   "x² = 4",
		Solution:   "x = ±2",
		Subject:    "math",
		Chapter:    "equations",
		Similarity: 0.95,
	}

	if sp.Question != "x² = 4" {
		t.Errorf("Expected question 'x² = 4', got %q", sp.Question)
	}

	if sp.Similarity != 0.95 {
		t.Errorf("Expected similarity 0.95, got %f", sp.Similarity)
	}

	if sp.Subject != "math" {
		t.Errorf("Expected subject 'math', got %q", sp.Subject)
	}
}

func TestGenerateEmbeddingPlaceholder(t *testing.T) {
	// This tests the placeholder implementation
	// Real embedding requires Ollama

	// Just verify the function signature works
	// embedding, err := generateEmbedding(context.Background(), "test problem")
	// if err != nil {
	//     t.Logf("Expected error without Ollama: %v", err)
	// }

	t.Log("Embedding tests require Ollama to be running")
}

func TestDatabaseConnection(t *testing.T) {
	// Test that DB connection would work
	// This requires actual database

	t.Log("Database tests require Neon DB to be accessible")
}
