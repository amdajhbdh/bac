package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// KnowledgeGraph represents extracted entities and relationships
type KnowledgeGraph struct {
	Entities      []Entity       `json:"entities"`
	Relationships []Relationship `json:"relationships"`
	Concepts      []Concept      `json:"concepts"`
}

// Entity represents a named entity
type Entity struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // person, place, concept, formula
	Source string `json:"source"`
	Page   int    `json:"page,omitempty"`
}

// Relationship represents connections between entities
type Relationship struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    string `json:"type"` // depends_on, part_of, related_to, is_a
	Context string `json:"context,omitempty"`
}

// Concept represents a key concept
type Concept struct {
	Term       string   `json:"term"`
	Definition string   `json:"definition"`
	Subject    string   `json:"subject"`
	Keywords   []string `json:"keywords"`
}

var kgLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func extractKnowledgeGraph(ctx context.Context, text, source string) (*KnowledgeGraph, error) {
	kg := &KnowledgeGraph{
		Entities:      []Entity{},
		Relationships: []Relationship{},
		Concepts:      []Concept{},
	}

	// Use AI to extract entities and relationships
	prompt := fmt.Sprintf(`Extract knowledge graph from this educational content.

For the following text, identify:
1. ENTITIES: Key terms, definitions, formulas, people, places
2. RELATIONSHIPS: How entities connect (depends_on, part_of, related_to, is_a)
3. CONCEPTS: Key concepts with definitions

Return as JSON:
{
  "entities": [{"id": "e1", "name": "...", "type": "concept|formula|person|place", "source": "%s"}],
  "relationships": [{"from": "e1", "to": "e2", "type": "related_to", "context": "..."}],
  "concepts": [{"term": "...", "definition": "...", "subject": "math|physics|chemistry|biology", "keywords": [...]}]
}

Text:
%s

Return ONLY valid JSON.`, source, text[:min(3000, len(text))])

	response, err := callOllama(ctx, prompt)
	if err != nil {
		kgLogger.Error("AI extraction failed, using basic parsing", "error", err)
		return basicExtraction(text, source), nil
	}

	// Parse AI response
	if err := json.Unmarshal([]byte(response), kg); err != nil {
		kgLogger.Warn("Failed to parse AI response, using basic extraction", "error", err)
		return basicExtraction(text, source), nil
	}

	return kg, nil
}

func basicExtraction(text, source string) *KnowledgeGraph {
	kg := &KnowledgeGraph{
		Entities:      []Entity{},
		Relationships: []Relationship{},
		Concepts:      []Concept{},
	}

	// Extract formulas (LaTeX patterns)
	formulas := extractFormulas(text)
	for _, f := range formulas {
		kg.Entities = append(kg.Entities, Entity{
			ID:     fmt.Sprintf("formula_%d", len(kg.Entities)),
			Name:   f,
			Type:   "formula",
			Source: source,
		})
	}

	// Extract key terms (capitalized phrases)
	terms := extractKeyTerms(text)
	for _, t := range terms {
		kg.Entities = append(kg.Entities, Entity{
			ID:     fmt.Sprintf("term_%d", len(kg.Entities)),
			Name:   t,
			Type:   "concept",
			Source: source,
		})
	}

	return kg
}

func extractFormulas(text string) []string {
	var formulas []string
	// Match LaTeX display math $$...$$
	start := 0
	for {
		idx := strings.Index(text[start:], "$$")
		if idx == -1 {
			break
		}
		end := strings.Index(text[start+idx+2:], "$$")
		if end == -1 {
			break
		}
		formula := strings.TrimSpace(text[start+idx+2 : start+idx+2+end])
		if len(formula) > 2 {
			formulas = append(formulas, formula)
		}
		start = start + idx + 2
	}
	return formulas
}

func extractKeyTerms(text string) []string {
	// Simple extraction of potential key terms
	patterns := []string{
		"définition", "théorème", "propriété", "lemme",
		"corollaire", "axiome", "principe", "loi",
		"formule", "équation", "identité",
	}
	var terms []string
	lower := strings.ToLower(text)
	for _, p := range patterns {
		if strings.Contains(lower, p) {
			// Extract phrase around the pattern
			idx := strings.Index(lower, p)
			start := max(0, idx-30)
			end := min(len(text), idx+len(p)+30)
			phrase := strings.TrimSpace(text[start:end])
			if len(phrase) > 5 {
				terms = append(terms, phrase)
			}
		}
	}
	return terms
}

// SaveKnowledgeGraph saves to file
func SaveKnowledgeGraph(kg *KnowledgeGraph, pdfPath string) error {
	// Save to file
	outputDir := filepath.Join(filepath.Dir(pdfPath), "knowledge-graphs")
	os.MkdirAll(outputDir, 0755)

	outputPath := filepath.Join(outputDir, filepath.Base(pdfPath)+".kg.json")
	data, _ := json.MarshalIndent(kg, "", "  ")
	os.WriteFile(outputPath, data, 0644)

	kgLogger.Info("knowledge graph saved", "path", outputPath, "entities", len(kg.Entities))

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
