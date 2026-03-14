package nlm

import "testing"

func TestExtractRoute(t *testing.T) {
	tests := []struct {
		name    string
		problem string
		want    string
	}{
		{"math derivative", "Calcule la dérivée de x³", "math"},
		{"pc force", "Quelle est la force?", "pc"},
		{"svt cell", "Explique la structure de la cellule", "svt"},
		{"philosophy ethics", "Qu'est-ce que la liberté?", "philosophie"},
		{"unknown defaults to math", "random text here", "math"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractRoute(tt.problem)
			if got.Subject != tt.want {
				t.Errorf("ExtractRoute(%q) subject = %q, want %q", tt.problem, got.Subject, tt.want)
			}
		})
	}
}

func TestGenerateQueryHash(t *testing.T) {
	hash1 := GenerateQueryHash("test problem")
	hash2 := GenerateQueryHash("test problem")
	hash3 := GenerateQueryHash("different problem")

	if hash1 != hash2 {
		t.Errorf("Same input should produce same hash: %q != %q", hash1, hash2)
	}

	if hash1 == hash3 {
		t.Errorf("Different input should produce different hash")
	}

	if len(hash1) != 12 {
		t.Errorf("Hash should be 12 chars, got %d", len(hash1))
	}
}

func TestSelectNotebook(t *testing.T) {
	mathNotebook := SelectNotebook("math")
	pcNotebook := SelectNotebook("pc")
	unknownNotebook := SelectNotebook("unknown")

	if mathNotebook == "" {
		t.Error("math notebook should not be empty")
	}

	if pcNotebook == "" {
		t.Error("pc notebook should not be empty")
	}

	if unknownNotebook != mathNotebook {
		t.Errorf("unknown should default to math notebook")
	}
}

func TestExtractTopics(t *testing.T) {
	route := ExtractRoute("Calcule la dérivée de x² + 3x")

	if len(route.Topics) == 0 {
		t.Error("should extract at least one topic")
	}

	found := false
	for _, topic := range route.Topics {
		if topic == "derivative" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("should find 'derivative' topic, got %v", route.Topics)
	}
}
