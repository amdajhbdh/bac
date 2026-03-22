//! Vector Tools Library
//!
//! Provides pgvector operations functionality with HNSW index support.

pub mod client;
pub mod models;
pub mod search;
pub mod insert;
pub mod batch;
pub mod delete;
pub mod error;
pub mod service;

pub use client::PgVectorClient;
pub use error::VectorError;
pub use models::*;
pub use service::run;

/// Initialize the vector-tools library.
pub fn init() {
    tracing::info!("vector-tools initialized");
}
