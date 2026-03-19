//! Graph Tools - Binary Entry Point
//!
//! Starts the knowledge graph service on port 3005.

fn main() {
    // Initialize logging
    tracing_subscriber::fmt::init();

    tracing::info!("Starting graph-tools service on :3005");

    // Run the async service
    tokio::runtime::Runtime::new()
        .expect("Failed to create Tokio runtime")
        .block_on(graph_tools::service::run());
}
