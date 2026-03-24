//! BAC RAG Pipeline Binary

use bac_rag::{QueryRequest, QueryResponse, RagEngine};
use anyhow::Result;
use axum::{
    extract::State,
    http::StatusCode,
    response::IntoResponse,
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use std::sync::Arc;
use tokio::sync::RwLock;
use tower_http::cors::CorsLayer;
use tracing::{error, info};
use tracing_subscriber;

#[derive(Clone)]
pub struct AppState {
    engine: Arc<RwLock<Option<RagEngine>>>,
}

unsafe impl Send for AppState {}
unsafe impl Sync for AppState {}

impl AppState {
    pub fn new() -> Self {
        Self {
            engine: Arc::new(RwLock::new(None)),
        }
    }

    pub async fn init(&self, data_dir: PathBuf) -> Result<()> {
        let engine = RagEngine::new(data_dir).await?;
        *self.engine.write().await = Some(engine);
        Ok(())
    }
    
    pub fn engine(&self) -> Arc<RwLock<Option<RagEngine>>> {
        self.engine.clone()
    }
}

async fn query_handler(
    State(state): State<AppState>,
    Json(request): Json<QueryRequest>,
) -> impl IntoResponse {
    let engine = state.engine.read().await;
    let engine = match engine.as_ref() {
        Some(e) => e,
        None => return StatusCode::SERVICE_UNAVAILABLE.into_response(),
    };
    
    match engine.query(&request).await {
        Ok(response) => Json(response).into_response(),
        Err(e) => {
            error!("Query error: {}", e);
            StatusCode::INTERNAL_SERVER_ERROR.into_response()
        }
    }
}

async fn index_handler(
    State(state): State<AppState>,
    Json(request): Json<IndexRequest>,
) -> impl IntoResponse {
    let engine = state.engine.read().await;
    let engine = match engine.as_ref() {
        Some(e) => e,
        None => return StatusCode::SERVICE_UNAVAILABLE.into_response(),
    };
    
    let content = match std::fs::read_to_string(&request.path) {
        Ok(c) => c,
        Err(_) => return StatusCode::NOT_FOUND.into_response(),
    };
    
    let doc_type = request.doc_type.unwrap_or_else(|| "default".to_string());
    
    match engine.index_documents(vec![(request.path, content, doc_type)]).await {
        Ok(count) => Json(IndexResponse {
            indexed: count,
            status: "indexed".to_string(),
        }).into_response(),
        Err(e) => {
            error!("Index error: {}", e);
            StatusCode::INTERNAL_SERVER_ERROR.into_response()
        }
    }
}

async fn health_handler(State(state): State<AppState>) -> impl IntoResponse {
    let engine = state.engine.read().await;
    let (status, indexed, cache_size) = if let Some(e) = engine.as_ref() {
        ("ok".to_string(), e.indexed_count(), e.cache_len())
    } else {
        ("initializing".to_string(), 0, 0)
    };
    
    Json(HealthResponse {
        status,
        indexed_docs: indexed,
        cache_size,
    })
}

pub fn routes(state: AppState) -> Router {
    let cors = CorsLayer::permissive();
    
    Router::new()
        .route("/query", post(query_handler))
        .route("/index", post(index_handler))
        .route("/health", get(health_handler))
        .layer(cors)
        .with_state(state)
}

#[derive(Debug, Deserialize)]
pub struct IndexRequest {
    pub path: String,
    pub doc_type: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct IndexResponse {
    pub indexed: usize,
    pub status: String,
}

#[derive(Debug, Serialize)]
pub struct HealthResponse {
    pub status: String,
    pub indexed_docs: usize,
    pub cache_size: usize,
}

// CLI
use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "bac-rag")]
#[command(about = "RAG pipeline with hybrid search", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Serve {
        #[arg(long, default_value = "5000")]
        port: u16,
        #[arg(long, default_value = "./data")]
        data_dir: String,
    },
    Query {
        query: String,
        #[arg(long, default_value = "http://127.0.0.1:5000")]
        url: String,
    },
}

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt::init();
    
    let cli = Cli::parse();
    
    match cli.command {
        Commands::Serve { port, data_dir } => {
            info!("Starting RAG server on port {}", port);
            
            let state = AppState::new();
            state.init(PathBuf::from(data_dir)).await?;
            
            let app = routes(state);
            
            let addr = format!("0.0.0.0:{}", port);
            let listener = tokio::net::TcpListener::bind(&addr).await?;
            
            axum::serve(listener, app).await?;
        }
        Commands::Query { query, url } => {
            let client = reqwest::Client::new();
            let response = client.post(format!("{}/query", url))
                .json(&QueryRequest {
                    query,
                    k: Some(5),
                    use_cache: Some(true),
                    doc_filter: None,
                })
                .send()
                .await?;
            
            let result: QueryResponse = response.json().await?;
            
            println!("Query: {}", result.query);
            println!("Cached: {}", result.cached);
            println!("Latency: {}ms", result.latency_ms);
            println!("\nResults:");
            
            for (i, r) in result.results.iter().enumerate() {
                println!("{}. [{}] (score: {:.3})", i + 1, r.source, r.score);
                println!("   {}", r.content.chars().take(200).collect::<String>());
                println!();
            }
        }
    }
    
    Ok(())
}