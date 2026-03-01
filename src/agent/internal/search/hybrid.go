package search

import (
	"context"
	"log/slog"
)

type HybridClient struct {
	es        *Client
	vectorCol VectorCollection
}

type VectorCollection interface {
	Search(ctx context.Context, query []float32, subject string, limit int) ([]VectorResult, error)
}

type VectorResult struct {
	ID      string
	Score   float64
	Subject string
	Chapter string
}

func NewHybridClient(es *Client, vectorCol VectorCollection) *HybridClient {
	return &HybridClient{
		es:        es,
		vectorCol: vectorCol,
	}
}

type HybridResult struct {
	ID            string
	Title         string
	Content       string
	Subject       string
	Chapter       string
	ESScore       float64
	VectorScore   float64
	CombinedScore float64
	Source        string
}

func (h *HybridClient) Search(ctx context.Context, query string, queryVector []float32, subject string, chapter string, limit int) ([]HybridResult, error) {
	results := make(chan []HybridResult, 2)
	errs := make(chan error, 2)

	go func() {
		esResults, err := h.es.Search(ctx, query, subject, chapter, 0, limit*2)
		if err != nil {
			errs <- err
			return
		}

		hybrid := make([]HybridResult, 0, len(esResults.Hits))
		for _, hit := range esResults.Hits {
			hybrid = append(hybrid, HybridResult{
				ID:            hit.Source.ID,
				Title:         hit.Source.Title,
				Content:       hit.Source.Content,
				Subject:       hit.Source.Subject,
				Chapter:       hit.Source.Chapter,
				ESScore:       hit.Score,
				CombinedScore: hit.Score,
				Source:        "elasticsearch",
			})
		}
		results <- hybrid
	}()

	go func() {
		if queryVector == nil {
			results <- []HybridResult{}
			return
		}

		vectorResults, err := h.vectorCol.Search(ctx, queryVector, subject, limit*2)
		if err != nil {
			errs <- err
			return
		}

		hybrid := make([]HybridResult, 0, len(vectorResults))
		for _, vr := range vectorResults {
			hybrid = append(hybrid, HybridResult{
				ID:            vr.ID,
				Subject:       vr.Subject,
				Chapter:       vr.Chapter,
				VectorScore:   vr.Score,
				CombinedScore: vr.Score,
				Source:        "vector",
			})
		}
		results <- hybrid
	}()

	var esResults, vectorResults []HybridResult
	for i := 0; i < 2; i++ {
		select {
		case r := <-results:
			if esResults == nil {
				esResults = r
			} else {
				vectorResults = r
			}
		case err := <-errs:
			slog.Warn("hybrid search partial error", "error", err)
		}
	}

	combined := h.rrfMerge(esResults, vectorResults, limit)

	slog.Info("hybrid search complete", "es_results", len(esResults), "vector_results", len(vectorResults), "combined", len(combined))

	return combined, nil
}

func (h *HybridClient) rrfMerge(esResults, vectorResults []HybridResult, k int) []HybridResult {
	scores := make(map[string]float64)
	items := make(map[string]*HybridResult)

	for _, r := range esResults {
		rrfScore := 1.0 / (60.0 + r.CombinedScore)
		scores[r.ID] += rrfScore
		normalized := r.CombinedScore
		if normalized == 0 {
			normalized = 1.0
		}
		items[r.ID] = &HybridResult{
			ID:            r.ID,
			Title:         r.Title,
			Content:       r.Content,
			Subject:       r.Subject,
			Chapter:       r.Chapter,
			ESScore:       r.ESScore,
			VectorScore:   0,
			CombinedScore: normalized,
			Source:        "elasticsearch",
		}
	}

	for _, r := range vectorResults {
		rrfScore := 1.0 / (60.0 + r.CombinedScore)
		scores[r.ID] += rrfScore
		if existing, ok := items[r.ID]; ok {
			existing.VectorScore = r.VectorScore
			existing.Source = "hybrid"
		} else {
			normalized := r.CombinedScore
			if normalized == 0 {
				normalized = 1.0
			}
			items[r.ID] = &HybridResult{
				ID:            r.ID,
				Subject:       r.Subject,
				Chapter:       r.Chapter,
				VectorScore:   r.VectorScore,
				CombinedScore: normalized,
				Source:        "vector",
			}
		}
	}

	type scoredItem struct {
		id    string
		score float64
	}
	ranked := make([]scoredItem, 0, len(scores))
	for id, score := range scores {
		ranked = append(ranked, scoredItem{id: id, score: score})
	}

	for i := 0; i < len(ranked); i++ {
		for j := i + 1; j < len(ranked); j++ {
			if ranked[j].score > ranked[i].score {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}

	result := make([]HybridResult, 0, k)
	for i := 0; i < len(ranked) && i < k; i++ {
		if item, ok := items[ranked[i].id]; ok {
			item.CombinedScore = ranked[i].score
			result = append(result, *item)
		}
	}

	return result
}
