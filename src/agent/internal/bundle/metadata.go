package bundle

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode"
)

type DublinCore struct {
	Title       string `json:"title"`
	Creator     string `json:"creator"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
	Contributor string `json:"contributor"`
	Date        string `json:"date"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Identifier  string `json:"identifier"`
	Source      string `json:"source"`
	Language    string `json:"language"`
	Relation    string `json:"relation"`
	Coverage    string `json:"coverage"`
	Rights      string `json:"rights"`
}

type LOM struct {
	General struct {
		Title       string   `json:"title"`
		Language    string   `json:"language"`
		Description string   `json:"description"`
		Keywords    []string `json:"keywords"`
		Coverage    string   `json:"coverage"`
	} `json:"general"`
	LifeCycle struct {
		Version    string `json:"version"`
		Status     string `json:"status"`
		Contribute []struct {
			Role   string `json:"role"`
			Entity string `json:"entity"`
			Date   string `json:"date"`
		} `json:"contribute"`
	} `json:"lifeCycle"`
	Technical struct {
		Format      string `json:"format"`
		Size        int64  `json:"size"`
		Location    string `json:"location"`
		Requirement []struct {
			Type           string `json:"type"`
			Name           string `json:"name"`
			MinimumVersion string `json:"minimumVersion"`
			MaximumVersion string `json:"maximumVersion"`
		} `json:"requirement"`
	} `json:"technical"`
	Educational struct {
		InteractivityType    string   `json:"interactivityType"`
		LearningResourceType []string `json:"learningResourceType"`
		IntendedEndUserRole  []string `json:"intendedEndUserRole"`
		Context              []string `json:"context"`
		Difficulty           string   `json:"difficulty"`
		TypicalLearningTime  string   `json:"typicalLearningTime"`
	} `json:"educational"`
	Rights struct {
		Cost      bool   `json:"cost"`
		Copyright string `json:"copyright"`
	} `json:"rights"`
}

func MetadataToDublinCore(meta Metadata) DublinCore {
	return DublinCore{
		Title:       meta.Title,
		Subject:     meta.Subject,
		Description: meta.Description,
		Publisher:   "BAC Unified",
		Date:        meta.CreatedAt.Format("2006-01-02"),
		Format:      meta.Format,
		Language:    meta.Language,
		Rights:      meta.License,
	}
}

func DublinCoreToMetadata(dc DublinCore) Metadata {
	return Metadata{
		Title:       dc.Title,
		Subject:     dc.Subject,
		Description: dc.Description,
		Language:    dc.Language,
		License:     dc.Rights,
	}
}

func MetadataToLOM(meta Metadata) LOM {
	lom := LOM{}
	lom.General.Title = meta.Title
	lom.General.Language = meta.Language
	lom.General.Description = meta.Description
	lom.General.Keywords = meta.Keywords
	lom.General.Coverage = meta.Chapter

	lom.Technical.Format = meta.Format
	lom.Technical.Size = meta.Size

	lom.LifeCycle.Status = "final"
	lom.LifeCycle.Contribute = []struct {
		Role   string `json:"role"`
		Entity string `json:"entity"`
		Date   string `json:"date"`
	}{
		{
			Role:   "author",
			Entity: strings.Join(meta.Authors, ", "),
			Date:   meta.CreatedAt.Format("2006-01-02"),
		},
	}

	lom.Educational.LearningResourceType = []string{"question bank"}
	lom.Educational.Context = []string{"higher education"}

	return lom
}

func LOMToMetadata(lom LOM) Metadata {
	meta := Metadata{
		Title:       lom.General.Title,
		Language:    lom.General.Language,
		Description: lom.General.Description,
		Keywords:    lom.General.Keywords,
		Chapter:     lom.General.Coverage,
		Format:      lom.Technical.Format,
		Size:        lom.Technical.Size,
	}

	if len(lom.LifeCycle.Contribute) > 0 {
		meta.Authors = strings.Split(lom.LifeCycle.Contribute[0].Entity, ", ")
	}

	return meta
}

func ClassifyContent(title, content string) (string, []string) {
	keywords := extractKeywords(title + " " + content)

	subjects := map[string][]string{
		"math":       {"mathématiques", "algèbre", "géométrie", "calcul", "équation", "dérivée", "intégrale", "fonction", "limite", "matrix", "vecteur"},
		"physics":    {"physique", "mécanique", "électricité", "optique", "thermodynamique", "force", "vitesse", "accélération", "masse", "énergie"},
		"chemistry":  {"chimie", "atome", "molécule", "réaction", "acide", "base", "oxydation", "organique", "minérale", "stœchiométrie"},
		"biology":    {"biologie", "cellule", "ADN", "génétique", "évolution", "écologie", "anatomie", "physiologie", "photosynthèse"},
		"history":    {"histoire", "guerre", "empire", "dynastie", "révolution", "colonialisme", "indépendance"},
		"geography":  {"géographie", "continent", "pays", "océan", "montagne", "climat", "population", "économie"},
		"philosophy": {"philosophie", "justice", "liberté", "morale", "éthique", "existence", "connaissance", "vérité"},
		"french":     {"français", "littérature", "poésie", "roman", "théâtre", "grammaire", "linguistique"},
		"arabic":     {"عربي", "أدب", "شعر", "نثر", "قواعد", "بلاغة"},
		"english":    {"english", "literature", "grammar", "poetry", "novel"},
	}

	contentLower := strings.ToLower(title + " " + content)

	for subject, terms := range subjects {
		for _, term := range terms {
			if strings.Contains(contentLower, term) {
				return subject, keywords
			}
		}
	}

	return "general", keywords
}

func extractKeywords(text string) []string {
	text = strings.ToLower(text)

	re := regexp.MustCompile(`[^\w\s]`)
	text = re.ReplaceAllString(text, "")

	words := strings.Fields(text)

	stopWords := map[string]bool{
		"le": true, "la": true, "les": true, "un": true, "une": true, "des": true,
		"et": true, "ou": true, "mais": true, "donc": true, "car": true,
		"est": true, "sont": true, "était": true, "être": true, "avoir": true,
		"qui": true, "que": true, "quoi": true, "dont": true, "ce": true,
		"this": true, "that": true, "the": true, "and": true, "or": true,
	}

	wordCount := make(map[string]int)
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 3 && !stopWords[word] {
			wordCount[word]++
		}
	}

	type kw struct {
		word  string
		count int
	}

	var sorted []kw
	for w, c := range wordCount {
		sorted = append(sorted, kw{w, c})
	}

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].count > sorted[i].count {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	keywords := make([]string, 0, 10)
	for i := 0; i < len(sorted) && i < 10; i++ {
		keywords = append(keywords, sorted[i].word)
	}

	return keywords
}

func DetectLanguage(text string) string {
	arabicPattern := regexp.MustCompile(`[\u0600-\u06FF]`)
	arabicCount := len(arabicPattern.FindAllString(text, -1))

	if float64(arabicCount) > float64(len(text))*0.3 {
		return "ar"
	}

	frenchPatterns := []string{
		"le ", "la ", "les ", "un ", "une ", "des ",
		"est ", "sont ", "été ", "être ", "avoir ", "fait ",
		"pour", "avec", "dans", "sur", "ce ", "cette",
	}

	frenchCount := 0
	for _, p := range frenchPatterns {
		frenchCount += strings.Count(strings.ToLower(text), p)
	}

	if frenchCount > 5 {
		return "fr"
	}

	englishPatterns := []string{
		" the ", " is ", " are ", " was ", " were ", " have ", " has ",
		" for ", " with ", " this ", " that ",
	}

	englishCount := 0
	for _, p := range englishPatterns {
		englishCount += strings.Count(strings.ToLower(text), p)
	}

	if englishCount > 5 {
		return "en"
	}

	return "fr"
}

func NormalizeText(text string) string {
	text = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, text)

	text = strings.Join(strings.Fields(text), " ")

	return text
}

func ToDublinCoreJSON(meta Metadata) (string, error) {
	dc := MetadataToDublinCore(meta)
	data, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal dublin core: %w", err)
	}
	slog.Debug("converted to dublin core", "language", dc.Language)
	return string(data), nil
}

func ToLOMJSON(meta Metadata) (string, error) {
	lom := MetadataToLOM(meta)
	data, err := json.MarshalIndent(lom, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal LOM: %w", err)
	}
	slog.Debug("converted to LOM")
	return string(data), nil
}
