use anyhow::Result;
use clap::{Parser, Subcommand};
use rusqlite::Connection;
use serde_json::json;

#[derive(Parser)]
#[command(name = "vault-cli")]
#[command(about = "Study Vault CLI - Search and manage your notes")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Search notes
    Search {
        query: String,
        #[arg(short, long, default_value = "10")]
        limit: usize,
    },
    /// Show vault statistics
    Stats,
    /// Find orphaned notes
    Orphans,
    /// Search formulas
    Formulas {
        query: String,
        #[arg(short, long, default_value = "10")]
        limit: usize,
    },
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    let db_path = ".vault-index.db";
    
    match cli.command {
        Commands::Search { query, limit } => search(&query, limit, db_path)?,
        Commands::Stats => stats(db_path)?,
        Commands::Orphans => orphans(db_path)?,
        Commands::Formulas { query, limit } => formulas(&query, limit, db_path)?,
    }
    
    Ok(())
}

fn search(query: &str, limit: usize, db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    let mut stmt = conn.prepare(
        "SELECT path, title FROM fts_notes 
         WHERE fts_notes MATCH ? ORDER BY rank LIMIT ?"
    )?;
    
    let results: Vec<(String, String)> = stmt
        .query_map([query, &limit.to_string()], |row| {
            Ok((row.get(0)?, row.get(1)?))
        })?
        .collect::<Result<Vec<_>, _>>()?;
    
    println!("\n🔍 Search results for '{}':\n", query);
    for (i, (path, title)) in results.iter().enumerate() {
        println!("{}. {}", i + 1, title);
        println!("   📁 {}\n", path);
    }
    
    Ok(())
}

fn stats(db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    
    let notes: i64 = conn.query_row("SELECT COUNT(*) FROM notes", [], |row| row.get(0))?;
    let links: i64 = conn.query_row("SELECT COUNT(*) FROM links", [], |row| row.get(0))?;
    let tags: i64 = conn.query_row("SELECT COUNT(*) FROM tags", [], |row| row.get(0))?;
    let formulas: i64 = conn.query_row("SELECT COUNT(*) FROM formulas", [], |row| row.get(0))?;
    
    let mut stmt = conn.prepare("SELECT subject, COUNT(*) FROM notes WHERE subject IS NOT NULL GROUP BY subject")?;
    let by_subject: Vec<(String, i64)> = stmt
        .query_map([], |row| Ok((row.get(0)?, row.get(1)?)))?
        .collect::<Result<Vec<_>, _>>()?;
    
    println!("\n📊 Vault Statistics:\n");
    println!("  Total Notes: {}", notes);
    println!("  Total Links: {}", links);
    println!("  Total Tags: {}", tags);
    println!("  Total Formulas: {}", formulas);
    
    println!("\n📚 Notes by Subject:");
    for (subject, count) in by_subject {
        println!("  {}: {}", subject, count);
    }
    
    Ok(())
}

fn orphans(db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    let mut stmt = conn.prepare(
        "SELECT path, title FROM notes 
         WHERE path NOT IN (SELECT DISTINCT source_path FROM links)
         AND path NOT IN (SELECT DISTINCT target_path FROM links)
         AND type != 'moc'"
    )?;
    
    let orphans: Vec<(String, String)> = stmt
        .query_map([], |row| Ok((row.get(0)?, row.get(1)?)))?
        .collect::<Result<Vec<_>, _>>()?;
    
    println!("\n🔗 Orphaned Notes ({}):\n", orphans.len());
    for (path, title) in orphans {
        println!("  • {} ({})", title, path);
    }
    
    Ok(())
}

fn formulas(query: &str, limit: usize, db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    let mut stmt = conn.prepare(
        "SELECT DISTINCT path, formula FROM formulas 
         WHERE formula LIKE ? LIMIT ?"
    )?;
    
    let pattern = format!("%{}%", query);
    let results: Vec<(String, String)> = stmt
        .query_map([&pattern, &limit.to_string()], |row| {
            Ok((row.get(0)?, row.get(1)?))
        })?
        .collect::<Result<Vec<_>, _>>()?;
    
    println!("\n📐 Formulas containing '{}':\n", query);
    for (path, formula) in results {
        println!("  {}", formula);
        println!("  📁 {}\n", path);
    }
    
    Ok(())
}
