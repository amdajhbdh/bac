package nlm

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Route struct {
	NotebookID string
	Subject    string
	Topics     []string
	QueryHash  string
}

var subjectKeywords = map[string][]string{
	"math": {
		"équation", "fonction", "dérivée", "intégrale", "limite",
		"polynôme", "factorisation", "x²", "sqrt", "racine",
		"calcul", "algèbre", "géométrie", "trigonométrie",
		"matrice", "vecteur", "probabilité", "statistique",
		"dériver", "intégrer", "résoudre", "factoriser",
	},
	"pc": {
		"physique", "mécanique", "électrique", "optique",
		"thermodynamique", "chimie", "molécule", "atome",
		"force", "vitesse", "accélération", "tension", "courant",
		"chaleur", "température", "pression", "volume",
		"onde", "lumière", "son", "électronique",
	},
	"svt": {
		"biologie", "cellule", "ADN", "génétique",
		"écologie", "évolution", "organisme", "photosynthèse",
		"respiration", "mitose", "méiose", "protéine",
		"enzyme", "gène", "chromosome", "mutation",
	},
	"philosophie": {
		"philosophie", "morale", "éthique", "existence",
		"liberté", "justice", "vérité", "connaissance",
		"conscience", "pensée", "être", "devenir",
	},
}

var notebookMap = map[string]string{
	"math":        "16b01950-5766-4353-8bed-c7f67966cb6b",
	"pc":          "16b01950-5766-4353-8bed-c7f67966cb6b",
	"svt":         "16b01950-5766-4353-8bed-c7f67966cb6b",
	"philosophie": "16b01950-5766-4353-8bed-c7f67966cb6b",
}

func ExtractRoute(problem string) Route {
	problemLower := strings.ToLower(problem)

	subject := determineSubject(problemLower)
	topics := extractTopics(problemLower)
	queryHash := GenerateQueryHash(problem)

	notebookID := SelectNotebook(subject)

	return Route{
		NotebookID: notebookID,
		Subject:    subject,
		Topics:     topics,
		QueryHash:  queryHash,
	}
}

func determineSubject(problem string) string {
	for subject, keywords := range subjectKeywords {
		for _, kw := range keywords {
			if strings.Contains(problem, kw) {
				return subject
			}
		}
	}
	return "math"
}

func extractTopics(problem string) []string {
	var topics []string

	topicKeywords := map[string][]string{
		"derivative":     {"dérivée", "dériver", "tangent", "différentiel"},
		"integral":       {"intégrale", "intégrer", "primitive", "aire"},
		"equation":       {"équation", "résoudre", "solution"},
		"function":       {"fonction", "graphe", "courbe"},
		"polynomial":     {"polynôme", "degré", "racine"},
		"matrix":         {"matrice", "déterminant", "inversible"},
		"vector":         {"vecteur", "direction", "norme"},
		"mechanics":      {"mécanique", "force", "masse", "mouvement"},
		"optics":         {"optique", "lentille", "miroir", "lumière"},
		"thermodynamics": {"thermodynamique", "chaleur", "température", "entropie"},
		"electricity":    {"électrique", "tension", "courant", "résistance"},
		"chemistry":      {"chimie", "molécule", "atome", "réaction"},
		"cell":           {"cellule", "membrane", "organite", "noyau"},
		"genetics":       {"génétique", "gène", "chromosome", "ADN"},
		"ecology":        {"écologie", "écosystème", "biodiversité"},
		"evolution":      {"évolution", "sélection", "espèce"},
		"ethics":         {"éthique", "morale", "bien", "mal"},
		"metaphysics":    {"métaphysique", "être", "existence"},
	}

	for topic, keywords := range topicKeywords {
		for _, kw := range keywords {
			if strings.Contains(problem, kw) {
				topics = append(topics, topic)
				break
			}
		}
	}

	if len(topics) == 0 {
		topics = []string{"general"}
	}

	return topics
}

func GenerateQueryHash(problem string) string {
	normalized := strings.ToLower(strings.TrimSpace(problem))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])[:12]
}

func SelectNotebook(subject string) string {
	if id, ok := notebookMap[subject]; ok {
		return id
	}
	return notebookMap["math"]
}

type NotebookConfig struct {
	Notebooks map[string]struct {
		ID     string   `yaml:"id"`
		Name   string   `yaml:"name"`
		Topics []string `yaml:"topics"`
	} `yaml:"notebooks"`
	Cache struct {
		TTLSeconds   int    `yaml:"default_ttl_seconds"`
		MaxEntries   int    `yaml:"max_entries"`
		GarageBucket string `yaml:"garage_bucket"`
	} `yaml:"cache"`
}

func LoadNotebookConfig(path string) (*NotebookConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config NotebookConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	for subject, nb := range config.Notebooks {
		notebookMap[subject] = nb.ID
	}

	return &config, nil
}
