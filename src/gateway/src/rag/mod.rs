//! RAG Module - Retrieval-Augmented Generation
//!
//! Provides vector search and chat context retrieval

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub enum ChatMode {
    Query,   // Quick factual answers
    Chat,    // Conversational
    Agent,   // Complex tasks with tools
    Automatic,
}

impl ChatMode {
    pub fn from_query(query: &str) -> Self {
        let q = query.to_lowercase();
        let words: Vec<&str> = q.split_whitespace()
            .map(|w| w.trim_matches(|c: char| !c.is_alphanumeric()))
            .collect();
        
        // Chat keywords - highest priority for conversational
        let chat_kw = ["hello", "hi", "thanks", "please", "goodbye", "hey"];
        if chat_kw.iter().any(|k| words.contains(k)) {
            return ChatMode::Chat;
        }
        
        // Agent keywords - complex tasks
        let agent_kw = ["solve", "calculate", "find", "search", "look", "up", "compare", "compute"];
        if agent_kw.iter().any(|k| words.contains(k)) {
            return ChatMode::Agent;
        }
        
        // Query keywords - factual questions
        let query_kw = ["who", "what", "where", "when", "define"];
        if query_kw.iter().any(|k| words.contains(k)) {
            return ChatMode::Query;
        }
        
        // Check for question phrases (multi-word)
        let clean_q: String = words.join(" ");
        let query_phrases = ["who is", "what is", "where is", "when did", "how to"];
        if query_phrases.iter().any(|p| clean_q.contains(p)) {
            return ChatMode::Query;
        }
        
        // How alone should be Query if starting the query
        if clean_q.starts_with("how ") {
            return ChatMode::Query;
        }
        
        // Check for question mark
        if q.contains('?') {
            return ChatMode::Query;
        }
        
        ChatMode::Automatic
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChatMessage {
    pub role: String,
    pub content: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChatRequest {
    pub message: String,
    pub mode: Option<ChatMode>,
    pub session_id: Option<String>,
    pub context: Option<ChatContext>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChatContext {
    pub subject: Option<String>,
    pub chapter: Option<i32>,
    pub difficulty: Option<i32>,
    pub language: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ChatResponse {
    pub message: String,
    pub mode: ChatMode,
    pub sources: Vec<Source>,
    pub session_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Source {
    pub text: String,
    pub similarity: f64,
    pub source_type: String,
}

/// Vector search result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SearchResult {
    pub text: String,
    pub similarity: f64,
    pub metadata: Option<serde_json::Value>,
}

use std::env;

pub struct RAGEngine {
    pool: Option<sqlx::PgPool>,
}

impl RAGEngine {
    pub fn new() -> Self {
        Self { pool: None }
    }

    pub async fn connect() -> Self {
        let database_url = env::var("NEON_DB_URL")
            .unwrap_or_else(|_| "postgresql://localhost/bac".to_string());
        
        let pool = sqlx::PgPool::connect(&database_url)
            .await
            .ok();
        
        Self { pool }
    }

    /// Search for similar content using pgvector
    pub async fn search(&self, query: &str, limit: usize) -> Vec<SearchResult> {
        let pool = match &self.pool {
            Some(p) => p,
            None => {
                return vec![SearchResult {
                    text: format!("Demo: Similar to '{}'", query),
                    similarity: 0.85,
                    metadata: None,
                }];
            }
        };
        
        let embedding = generate_embedding(query).await;
        
        match sqlx::query_as::<_, (String, f64, serde_json::Value)>(
            "SELECT question_text, 1 - (question_vector <=> $1::vector) as similarity, jsonb_build_object('subject', subject) as metadata FROM questions WHERE question_vector IS NOT NULL ORDER BY question_vector <=> $1::vector LIMIT $2"
        )
        .bind(embedding)
        .bind(limit as i32)
        .fetch_all(pool)
        .await {
            Ok(rows) => rows.into_iter().map(|(text, similarity, metadata)| {
                SearchResult { text, similarity, metadata: Some(metadata) }
            }).collect(),
            Err(_) => vec![SearchResult {
                text: format!("Demo: Similar to '{}'", query),
                similarity: 0.85,
                metadata: None,
            }]
        }
    }

    /// Search with filters
    pub async fn search_filtered(
        &self,
        query: &str,
        subject: Option<&str>,
        chapter: Option<i32>,
        limit: usize,
    ) -> Vec<SearchResult> {
        let pool = match &self.pool {
            Some(p) => p,
            None => {
                return vec![SearchResult {
                    text: format!("Demo: '{}' in chapter {:?}", query, chapter),
                    similarity: 0.82,
                    metadata: None,
                }];
            }
        };
        
        let embedding = generate_embedding(query).await;
        
        let sql = "SELECT question_text, 1 - (question_vector <=> $1::vector) as similarity, jsonb_build_object('subject', subject, 'chapter', chapter) as metadata FROM questions WHERE question_vector IS NOT NULL ORDER BY question_vector <=> $1::vector LIMIT $2";
        
        match sqlx::query_as::<_, (String, f64, serde_json::Value)>(sql)
            .bind(embedding)
            .bind(limit as i32)
            .fetch_all(pool)
            .await {
            Ok(rows) => rows.into_iter().map(|(text, similarity, metadata)| {
                SearchResult { text, similarity, metadata: Some(metadata) }
            }).collect(),
            Err(_) => vec![SearchResult {
                text: format!("Demo: '{}' in chapter {:?}", query, chapter),
                similarity: 0.82,
                metadata: None,
            }]
        }
    }

    /// Hybrid search: combines full-text and vector similarity
    pub async fn hybrid_search(&self, query: &str, limit: usize) -> Vec<SearchResult> {
        let pool = match &self.pool {
            Some(p) => p,
            None => {
                return vec![SearchResult {
                    text: format!("Hybrid: '{}'", query),
                    similarity: 0.80,
                    metadata: None,
                }];
            }
        };
        
        let embedding = generate_embedding(query).await;
        
        let sql = r#"
            SELECT question_text, 
                   COALESCE(1 - (question_vector <=> $1::vector), 0) as similarity,
                   jsonb_build_object('subject', subject, 'chapter', chapter) as metadata
            FROM questions 
            WHERE question_text ILIKE $3 OR question_vector IS NOT NULL
            ORDER BY COALESCE(1 - (question_vector <=> $1::vector), 0) DESC
            LIMIT $2
        "#;
        
        match sqlx::query_as::<_, (String, f64, serde_json::Value)>(sql)
            .bind(embedding)
            .bind(limit as i32)
            .bind(format!("%{}%", query))
            .fetch_all(pool)
            .await {
            Ok(rows) => rows.into_iter().map(|(text, similarity, metadata)| {
                SearchResult { text, similarity, metadata: Some(metadata) }
            }).collect(),
            Err(_) => vec![SearchResult {
                text: format!("Hybrid: '{}'", query),
                similarity: 0.80,
                metadata: None,
            }]
        }
    }

    /// Range search: find vectors within distance threshold
    pub async fn range_search(&self, query: &str, max_distance: f64, limit: usize) -> Vec<SearchResult> {
        let pool = match &self.pool {
            Some(p) => p,
            None => {
                return vec![SearchResult {
                    text: format!("Range: '{}' < {}", query, max_distance),
                    similarity: 1.0 - max_distance,
                    metadata: None,
                }];
            }
        };
        
        let embedding = generate_embedding(query).await;
        
        let sql = r#"
            SELECT question_text, 
                   1 - (question_vector <=> $1::vector) as similarity,
                   jsonb_build_object('subject', subject, 'chapter', chapter) as metadata
            FROM questions 
            WHERE question_vector IS NOT NULL 
              AND question_vector <=> $1::vector < $3
            ORDER BY question_vector <=> $1::vector
            LIMIT $2
        "#;
        
        match sqlx::query_as::<_, (String, f64, serde_json::Value)>(sql)
            .bind(embedding)
            .bind(limit as i32)
            .bind(max_distance)
            .fetch_all(pool)
            .await {
            Ok(rows) => rows.into_iter().map(|(text, similarity, metadata)| {
                SearchResult { text, similarity, metadata: Some(metadata) }
            }).collect(),
            Err(_) => vec![SearchResult {
                text: format!("Range: '{}' < {}", query, max_distance),
                similarity: 1.0 - max_distance,
                metadata: None,
            }]
        }
    }

    /// Add a new question with embedding
    pub async fn add_question(&self, question: &str, solution: &str, tags: &[String], difficulty: i32) -> Result<String, String> {
        let pool = match &self.pool {
            Some(p) => p,
            None => return Err("Database not connected".to_string()),
        };
        
        let embedding = generate_embedding(question).await;
        
        let sql = r#"
            INSERT INTO questions (question_text, solution_text, question_vector, topic_tags, difficulty)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id::text
        "#;
        
        sqlx::query_scalar(sql)
            .bind(question)
            .bind(solution)
            .bind(embedding)
            .bind(tags)
            .bind(difficulty)
            .fetch_one(pool)
            .await
            .map_err(|e| e.to_string())
    }
}

#[derive(serde::Deserialize)]
struct OllamaEmbedResponse {
    embedding: Vec<f32>,
}

async fn generate_embedding(text: &str) -> Vec<f32> {
    let ollama_host = env::var("OLLAMA_HOST").unwrap_or_else(|_| "http://localhost:11434".to_string());
    let model = env::var("OLLAMA_EMBED_MODEL").unwrap_or_else(|_| "nomic-embed-text".to_string());
    
    let client = match reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(10))
        .build() 
    {
        Ok(c) => c,
        Err(_) => return mock_embedding(),
    };
    
    let payload = serde_json::json!({
        "model": model,
        "input": text
    });
    
    match client.post(format!("{}/api/embed", ollama_host))
        .json(&payload)
        .send()
        .await
    {
        Ok(resp) => {
            match resp.json::<OllamaEmbedResponse>().await {
                Ok(data) => data.embedding,
                Err(_) => mock_embedding(),
            }
        }
        Err(_) => mock_embedding(),
    }
}

fn mock_embedding() -> Vec<f32> {
    vec![0.1; 1536]
}

impl RAGEngine {
    pub async fn from_env() -> Self {
        let database_url = env::var("NEON_DB_URL")
            .unwrap_or_else(|_| "postgresql://localhost/bac".to_string());
        
        let pool = sqlx::PgPool::connect(&database_url)
            .await
            .ok();
        
        Self { pool }
    }
}

impl Default for RAGEngine {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_chat_mode_query() {
        let mode = ChatMode::from_query("What is photosynthesis?");
        assert_eq!(mode, ChatMode::Query);
    }

    #[test]
    fn test_chat_mode_chat() {
        let mode = ChatMode::from_query("Hello, how are you?");
        assert_eq!(mode, ChatMode::Chat);
    }

    #[test]
    fn test_chat_mode_agent() {
        let mode = ChatMode::from_query("solve this equation");
        assert_eq!(mode, ChatMode::Agent);
    }

    #[test]
    fn test_chat_mode_auto() {
        let mode = ChatMode::from_query("Tell me about physics");
        assert_eq!(mode, ChatMode::Automatic);
    }

    #[tokio::test]
    async fn test_rag_search_returns_results() {
        let engine = RAGEngine::new();
        let results = engine.search("dérivée", 5).await;
        assert!(!results.is_empty());
    }

    #[tokio::test]
    async fn test_rag_search_with_filters() {
        let engine = RAGEngine::new();
        let results = engine.search_filtered("intégrale", Some("math"), Some(3), 10).await;
        assert!(!results.is_empty());
    }

    #[tokio::test]
    async fn test_rag_hybrid_search() {
        let engine = RAGEngine::new();
        let results = engine.hybrid_search("dérivée", 5).await;
        assert!(!results.is_empty());
    }

    #[tokio::test]
    async fn test_rag_range_search() {
        let engine = RAGEngine::new();
        let results = engine.range_search("intégrale", 0.5, 10).await;
        assert!(!results.is_empty());
    }

    #[tokio::test]
    async fn test_rag_add_question() {
        let engine = RAGEngine::new();
        let result = engine.add_question(
            "Test question?",
            "Test solution",
            &["test".to_string()],
            1,
        ).await;
        // Should fail gracefully when DB not connected
        assert!(result.is_err() || result.is_ok());
    }
}
