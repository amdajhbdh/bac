package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Client struct {
	client    *elasticsearch.Client
	indexName string
}

type Document struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Subject   string                 `json:"subject"`
	Chapter   string                 `json:"chapter"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Tags      []string               `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
}

type SearchResult struct {
	Total  int64              `json:"total"`
	Hits   []Hit              `json:"hits"`
	Facets map[string][]Facet `json:"facets"`
}

type Hit struct {
	ID     string   `json:"_id"`
	Score  float64  `json:"_score"`
	Source Document `json:"_source"`
}

type Facet struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

func NewClient() (*Client, error) {
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://localhost:9200"
	}

	indexName := os.Getenv("ELASTICSEARCH_INDEX")
	if indexName == "" {
		indexName = "bac-resources"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	slog.Info("elasticsearch client created", "url", url, "index", indexName)

	return &Client{
		client:    client,
		indexName: indexName,
	}, nil
}

func (c *Client) CreateIndex(ctx context.Context) error {
	exists, err := c.IndexExists(ctx)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"french_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter":    []string{"lowercase", "asciifolding", "french_stemmer"},
					},
				},
				"filter": map[string]interface{}{
					"french_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "french",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":         map[string]string{"type": "keyword"},
				"type":       map[string]string{"type": "keyword"},
				"subject":    map[string]string{"type": "keyword"},
				"chapter":    map[string]string{"type": "keyword"},
				"title":      map[string]string{"type": "text", "analyzer": "french_analyzer"},
				"content":    map[string]string{"type": "text", "analyzer": "french_analyzer"},
				"tags":       map[string]string{"type": "keyword"},
				"metadata":   map[string]interface{}{"type": "object", "enabled": false},
				"created_at": map[string]string{"type": "date"},
			},
		},
	}

	body, _ := json.Marshal(mapping)

	req := esapi.IndicesCreateRequest{
		Index: c.indexName,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index error: %s", res.String())
	}

	slog.Info("elasticsearch index created", "index", c.indexName)
	return nil
}

func (c *Client) IndexExists(ctx context.Context) (bool, error) {
	res, err := c.client.Indices.Exists(
		[]string{c.indexName},
		c.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

func (c *Client) Index(ctx context.Context, doc Document) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      c.indexName,
		DocumentID: doc.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}

	return nil
}

func (c *Client) Search(ctx context.Context, query string, subject string, chapter string, from int, size int) (*SearchResult, error) {
	must := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^2", "content", "tags"},
				"type":   "best_fields",
			},
		},
	}

	if subject != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]string{"subject": subject},
		})
	}

	if chapter != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]string{"chapter": chapter},
		})
	}

	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"from": from,
		"size": size,
		"aggs": map[string]interface{}{
			"subjects": map[string]interface{}{
				"terms": map[string]string{"field": "subject"},
			},
			"chapters": map[string]interface{}{
				"terms": map[string]string{"field": "chapter"},
			},
		},
	}

	body, _ := json.Marshal(searchBody)

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.indexName),
		c.client.Search.WithBody(bytes.NewReader(body)),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string          `json:"_id"`
				Score  float64         `json:"_score"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
		Aggregations map[string]json.RawMessage `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	searchResult := &SearchResult{
		Total:  result.Hits.Total.Value,
		Hits:   make([]Hit, 0, len(result.Hits.Hits)),
		Facets: make(map[string][]Facet),
	}

	for _, hit := range result.Hits.Hits {
		var doc Document
		json.Unmarshal(hit.Source, &doc)
		searchResult.Hits = append(searchResult.Hits, Hit{
			ID:     hit.ID,
			Score:  hit.Score,
			Source: doc,
		})
	}

	for aggName, aggData := range result.Aggregations {
		var aggResult struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int64  `json:"doc_count"`
			} `json:"buckets"`
		}
		json.Unmarshal(aggData, &aggResult)

		facets := make([]Facet, 0, len(aggResult.Buckets))
		for _, bucket := range aggResult.Buckets {
			facets = append(facets, Facet{
				Key:   bucket.Key,
				Count: bucket.DocCount,
			})
		}
		searchResult.Facets[aggName] = facets
	}

	return searchResult, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      c.indexName,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && !strings.Contains(res.String(), "404") {
		return fmt.Errorf("delete error: %s", res.String())
	}

	return nil
}

func (c *Client) BulkIndex(ctx context.Context, docs []Document) error {
	if len(docs) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, doc := range docs {
		meta := map[string]interface{}{
			"index": map[string]string{
				"_index": c.indexName,
				"_id":    doc.ID,
			},
		}
		metaBytes, _ := json.Marshal(meta)
		buf.Write(metaBytes)
		buf.WriteByte('\n')

		docBytes, _ := json.Marshal(doc)
		buf.Write(docBytes)
		buf.WriteByte('\n')
	}

	req := esapi.BulkRequest{
		Body: bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("bulk index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index error: %s", res.String())
	}

	return nil
}
