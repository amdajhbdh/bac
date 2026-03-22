//! Vault Tools Library
//!
//! Provides Obsidian vault operations for reading, writing, linking, and searching.
//!
//! # Features
//! - Read/write notes with frontmatter support
//! - Wikilink creation between notes
//! - Map of Content (MOC) updates
//! - Full-text search in vault

pub mod read;
pub mod write;
pub mod link;
pub mod moc;
pub mod search;
pub mod service;
pub mod models;

pub use models::*;
pub use service::run;

/// Initialize the vault-tools library.
pub fn init() {
    tracing::info!("vault-tools initialized");
}
