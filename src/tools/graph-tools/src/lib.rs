//! Graph Tools Library
//!
//! Provides knowledge graph operations for Obsidian vault:
//! - Generate .canvas files from subjects
//! - Extract entities from notes
//! - Build relationship graphs
//! - Export graphs as JSON
//! - Generate visualization URLs

pub mod canvas;
pub mod extract;
pub mod relationships;
pub mod export;
pub mod visualize;
pub mod models;
pub mod service;

// Re-export models for external use
pub use models::{CanvasEdge, CanvasFile, CanvasNode, EdgeType, GraphId, NodeType};

// Re-export main functions for convenience
pub use canvas::generate_canvas;
pub use extract::extract_entities_from_notes;
pub use relationships::build_relationships;
pub use export::export_json;
pub use visualize::generate_visualization;
pub use models::*;

/// Initialize the graph-tools library.
pub fn init() {
    tracing::info!("graph-tools initialized");
}
