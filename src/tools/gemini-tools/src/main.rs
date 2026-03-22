//! Gemini Tools - CLI and HTTP Service
//!
//! Provides Gemini API wrapper with HTTP server on :3001

use anyhow::Result;
use clap::{Parser, Subcommand};
use gemini_tools::{analyze, correct, embed, extract, generate, GeminiClient};
use tracing_subscriber;

// CLI Commands
#[derive(Parser)]
#[command(name = "gemini-tools")]
#[command(about = "Gemini API wrapper tool", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Analyze educational content
    Analyze {
        /// Content to analyze
        #[arg(long)]
        content: String,
        /// Subject (optional)
        #[arg(long)]
        subject: Option<String>,
    },
    /// Extract entities from content
    Extract {
        /// Content to extract from
        #[arg(long)]
        content: String,
    },
    /// Generate study notes
    Generate {
        /// Topic to generate notes for
        #[arg(long)]
        topic: String,
        /// Subject (optional)
        #[arg(long)]
        subject: Option<String>,
        /// Format: markdown or latex
        #[arg(long, default_value = "markdown")]
        format: String,
    },
    /// Create embeddings
    Embed {
        /// Text to embed
        #[arg(long)]
        text: String,
        /// Task type (optional)
        #[arg(long)]
        task_type: Option<String>,
    },
    /// Correct OCR text
    Correct {
        /// OCR text to correct
        #[arg(long)]
        ocr_text: String,
        /// Language (optional)
        #[arg(long)]
        language: Option<String>,
    },
    /// Start HTTP server
    Serve {
        /// Port to listen on
        #[arg(long, default_value = "3001")]
        port: u16,
    },
}

#[tokio::main]
async fn main() -> Result<()> {
    // Load .env file if present
    let _ = dotenvy::dotenv();

    // Initialize logging
    tracing_subscriber::fmt::init();

    let cli = Cli::parse();

    match cli.command {
        Commands::Analyze { content, subject } => {
            run_analyze(&content, subject.as_deref()).await?;
        }
        Commands::Extract { content } => {
            run_extract(&content).await?;
        }
        Commands::Generate { topic, subject, format } => {
            run_generate(&topic, subject.as_deref(), Some(&format)).await?;
        }
        Commands::Embed { text, task_type } => {
            run_embed(&text, task_type.as_deref()).await?;
        }
        Commands::Correct { ocr_text, language } => {
            run_correct(&ocr_text, language.as_deref()).await?;
        }
        Commands::Serve { port } => {
            run_server(port).await?;
        }
    }

    Ok(())
}

async fn get_client() -> Result<GeminiClient> {
    GeminiClient::from_env()
}

async fn run_analyze(content: &str, subject: Option<&str>) -> Result<()> {
    let client = get_client().await?;
    let result = analyze::analyze(&client, content, subject).await?;

    println!("{}", serde_json::to_string_pretty(&result)?);
    Ok(())
}

async fn run_extract(content: &str) -> Result<()> {
    let client = get_client().await?;
    let result = extract::extract(&client, content).await?;

    println!("{}", serde_json::to_string_pretty(&result)?);
    Ok(())
}

async fn run_generate(topic: &str, subject: Option<&str>, format: Option<&str>) -> Result<()> {
    let client = get_client().await?;
    let result = generate::generate(&client, topic, subject, format).await?;

    println!("{}", serde_json::to_string_pretty(&result)?);
    Ok(())
}

async fn run_embed(text: &str, task_type: Option<&str>) -> Result<()> {
    let client = get_client().await?;
    let result = embed::embed(&client, text, task_type).await?;

    println!("{}", serde_json::to_string_pretty(&result)?);
    Ok(())
}

async fn run_correct(ocr_text: &str, language: Option<&str>) -> Result<()> {
    let client = get_client().await?;
    let result = correct::correct(&client, ocr_text, language).await?;

    println!("{}", serde_json::to_string_pretty(&result)?);
    Ok(())
}

async fn run_server(port: u16) -> Result<()> {
    let client = get_client().await?;
    gemini_tools::run_server(port, client).await?;
    Ok(())
}
