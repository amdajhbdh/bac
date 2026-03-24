//! BAC RAG Pipeline - Hybrid BM25 + Vector Search

pub mod embeddings;
pub mod cache;

use std::path::PathBuf;
use std::sync::Arc;

use anyhow::Result;
use parking_lot::RwLock;
use serde::{Deserialize, Serialize};
use tantivy::collector::TopDocs;
use tantivy::query::QueryParser;
use tantivy::schema::*;
use tantivy::{doc, Index, IndexWriter, ReloadPolicy};

// ============================================================================
// Data Types
// ============================================================================

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RetrievalResult {
    pub content: String,
    pub score: f32,
    pub source: String,
    pub doc_type: String,
    pub path: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QueryRequest {
    pub query: String,
    pub k: Option<usize>,
    pub use_cache: Option<bool>,
    pub doc_filter: Option<Vec<String>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QueryResponse {
    pub results: Vec<RetrievalResult>,
    pub query: String,
    pub cached: bool,
    pub latency_ms: u64,
}

// ============================================================================
// RAG Engine
// ============================================================================

pub struct RagEngine {
    index: Index,
    schema: Schema,
    embedder: embeddings::SentenceEmbeddings,
    cache: Arc<cache::QueryCache>,
    content_store: Arc<RwLock<Vec<DocumentStore>>>,
}

struct DocumentStore {
    path: String,
    content: String,
    doc_type: String,
}

impl RagEngine {
    pub async fn new(data_dir: PathBuf) -> Result<Self> {
        let embedder = embeddings::SentenceEmbeddings::new("all-MiniLM-L6-v2");
        
        let mut schema_builder = Schema::builder();
        schema_builder.add_text_field("content", TEXT | STORED);
        schema_builder.add_text_field("path", TEXT | STORED);
        schema_builder.add_text_field("doc_type", TEXT | STORED);
        let schema = schema_builder.build();

        let index = if data_dir.join("index").exists() {
            Index::open_in_dir(data_dir.join("index"))?
        } else {
            std::fs::create_dir_all(data_dir.join("index"))?;
            Index::create_in_dir(data_dir.join("index"), schema.clone())?
        };

        let cache = Arc::new(cache::QueryCache::new(data_dir.join("cache")));
        let content_store = Arc::new(RwLock::new(Vec::new()));

        Ok(Self {
            index,
            schema,
            embedder,
            cache,
            content_store,
        })
    }

    pub async fn index_documents(&self, documents: Vec<(String, String, String)>) -> Result<usize> {
        let mut index_writer: IndexWriter = self.index.writer(50_000_000)?;
        
        let content_field = self.schema.get_field("content")?;
        let path_field = self.schema.get_field("path")?;
        let doc_type_field = self.schema.get_field("doc_type")?;

        let mut indexed = 0;
        
        for (path, content, doc_type) in documents {
            let chunks = self.chunk_text(&content, 1000, 200);
            
            for chunk in chunks {
                index_writer.add_document(doc!(
                    content_field => chunk,
                    path_field => path.clone(),
                    doc_type_field => doc_type.clone(),
                ))?;
            }
            
            indexed += 1;
            
            self.content_store.write().push(DocumentStore {
                path,
                content,
                doc_type,
            });
        }

        index_writer.commit()?;
        Ok(indexed)
    }

    fn chunk_text(&self, text: &str, chunk_size: usize, overlap: usize) -> Vec<String> {
        let mut chunks = Vec::new();
        let chars: Vec<char> = text.chars().collect();
        
        if chars.len() <= chunk_size {
            chunks.push(text.to_string());
            return chunks;
        }

        let mut start = 0;
        while start < chars.len() {
            let end = (start + chunk_size).min(chars.len());
            let chunk: String = chars[start..end].iter().collect();
            chunks.push(chunk);
            start += chunk_size - overlap;
        }

        chunks
    }

    pub async fn query(&self, request: &QueryRequest) -> Result<QueryResponse> {
        let start = std::time::Instant::now();
        let k = request.k.unwrap_or(5);
        let use_cache = request.use_cache.unwrap_or(true);

        if use_cache {
            if let Some(cached) = self.cache.get(&request.query) {
                return Ok(QueryResponse {
                    results: cached,
                    query: request.query.clone(),
                    cached: true,
                    latency_ms: start.elapsed().as_millis() as u64,
                });
            }
        }

        let results = self.hybrid_search(&request.query, k).await?;

        if use_cache {
            self.cache.set(&request.query, results.clone()).await;
        }

        Ok(QueryResponse {
            results,
            query: request.query.clone(),
            cached: false,
            latency_ms: start.elapsed().as_millis() as u64,
        })
    }

    async fn hybrid_search(&self, query: &str, k: usize) -> Result<Vec<RetrievalResult>> {
        let vector_results = self.vector_search(query, k * 2).await?;
        let bm25_results = self.bm25_search(query, k * 2);

        let mut merged: std::collections::HashMap<String, RetrievalResult> = std::collections::HashMap::new();
        
        for r in vector_results {
            let key = r.content.chars().take(100).collect::<String>();
            merged.insert(key, RetrievalResult {
                content: r.content,
                score: r.score * 0.5,
                source: "vector".to_string(),
                doc_type: r.doc_type,
                path: r.path,
            });
        }

        for r in bm25_results {
            let key = r.content.chars().take(100).collect::<String>();
            if let Some(existing) = merged.get_mut(&key) {
                existing.score += r.score * 0.5;
                if existing.source == "vector" {
                    existing.source = "hybrid".to_string();
                }
            } else {
                merged.insert(key, r);
            }
        }

        let mut results: Vec<RetrievalResult> = merged.into_values().collect();
        results.sort_by(|a, b| b.score.partial_cmp(&a.score).unwrap_or(std::cmp::Ordering::Equal));
        results.truncate(k);
        
        Ok(results)
    }

    async fn vector_search(&self, query: &str, k: usize) -> Result<Vec<RetrievalResult>> {
        let _query_embedding = self.embedder.embed(query).await?;
        
        let reader = self.index
            .reader_builder()
            .reload_policy(ReloadPolicy::OnCommitWithDelay)
            .try_into()?;
        
        let searcher = reader.searcher();
        let content_field = self.schema.get_field("content")?;
        let path_field = self.schema.get_field("path")?;
        let doc_type_field = self.schema.get_field("doc_type")?;

        let query_parser = QueryParser::for_index(&self.index, vec![content_field]);
        let tantivy_query = query_parser.parse_query(query)?;
        
        let top_docs = searcher.search(&tantivy_query, &TopDocs::with_limit(k))?;
        
        let mut results = Vec::new();
        for (_score, doc_address) in top_docs {
            let retrieved_doc: TantivyDocument = searcher.doc(doc_address)?;
            let content = retrieved_doc.get_first(content_field)
                .and_then(|v| v.as_str())
                .unwrap_or("")
                .to_string();
            let path = retrieved_doc.get_first(path_field)
                .and_then(|v| v.as_str())
                .unwrap_or("")
                .to_string();
            let doc_type = retrieved_doc.get_first(doc_type_field)
                .and_then(|v| v.as_str())
                .unwrap_or("default")
                .to_string();
            
            results.push(RetrievalResult {
                content,
                score: 1.0,
                source: "vector".to_string(),
                doc_type,
                path,
            });
        }

        Ok(results)
    }

    fn bm25_search(&self, query: &str, k: usize) -> Vec<RetrievalResult> {
        use bm25::{SearchEngineBuilder, Language};
        
        let documents = self.content_store.read();
        if documents.is_empty() {
            return Vec::new();
        }

        let corpus: Vec<String> = documents.iter()
            .map(|d| d.content.clone())
            .collect();
        
        let search_engine = SearchEngineBuilder::<u32>::with_corpus(Language::English, &corpus).build();
        let bm25_results = search_engine.search(query, k);
        
        let mut results: Vec<RetrievalResult> = bm25_results.into_iter()
            .take(k)
            .map(|r| {
                let doc_idx = r.document.id;
                let doc = &documents[doc_idx as usize];
                RetrievalResult {
                    content: doc.content.chars().take(500).collect(),
                    score: r.score,
                    source: "bm25".to_string(),
                    doc_type: doc.doc_type.clone(),
                    path: doc.path.clone(),
                }
            })
            .collect();
        
        results.sort_by(|a, b| b.score.partial_cmp(&a.score).unwrap_or(std::cmp::Ordering::Equal));
        
        results
    }

    pub fn indexed_count(&self) -> usize {
        self.content_store.read().len()
    }

    pub fn cache_len(&self) -> usize {
        self.cache.len()
    }
}

unsafe impl Send for RagEngine {}
unsafe impl Sync for RagEngine {}
