use anyhow::Result;
use clap::{Parser, Subcommand};
use rusqlite::Connection;
use serde_json::{json, Value};

#[derive(Parser)]
struct Cli {
    #[command(subcommand)]
    cmd: Cmd,
}

#[derive(Subcommand)]
enum Cmd {
    Search { q: String, #[arg(short, default_value = "10")] n: usize },
    Stats,
    Orphans,
    Formulas { q: String, #[arg(short, default_value = "10")] n: usize },
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    let conn = Connection::open(".vault-index.db")?;
    
    let output = match cli.cmd {
        Cmd::Search { q, n } => search(&conn, &q, n)?,
        Cmd::Stats => stats(&conn)?,
        Cmd::Orphans => orphans(&conn)?,
        Cmd::Formulas { q, n } => formulas(&conn, &q, n)?,
    };
    
    println!("{}", serde_json::to_string_pretty(&output)?);
    Ok(())
}

fn search(conn: &Connection, q: &str, n: usize) -> Result<Value> {
    let mut stmt = conn.prepare("SELECT path, title FROM fts_notes WHERE fts_notes MATCH ? ORDER BY rank LIMIT ?")?;
    let results: Vec<Value> = stmt.query_map([q, &n.to_string()], |row| {
        Ok(json!({"path": row.get::<_, String>(0)?, "title": row.get::<_, String>(1)?}))
    })?.collect::<Result<Vec<_>, _>>()?;
    
    Ok(json!({"query": q, "count": results.len(), "results": results}))
}

fn stats(conn: &Connection) -> Result<Value> {
    let notes: i64 = conn.query_row("SELECT COUNT(*) FROM notes", [], |r| r.get(0))?;
    let links: i64 = conn.query_row("SELECT COUNT(*) FROM links", [], |r| r.get(0))?;
    let tags: i64 = conn.query_row("SELECT COUNT(*) FROM tags", [], |r| r.get(0))?;
    let formulas: i64 = conn.query_row("SELECT COUNT(*) FROM formulas", [], |r| r.get(0))?;
    
    let mut stmt = conn.prepare("SELECT subject, COUNT(*) FROM notes WHERE subject IS NOT NULL GROUP BY subject")?;
    let by_subject: Vec<Value> = stmt.query_map([], |r| {
        Ok(json!({"subject": r.get::<_, String>(0)?, "count": r.get::<_, i64>(1)?}))
    })?.collect::<Result<Vec<_>, _>>()?;
    
    Ok(json!({"notes": notes, "links": links, "tags": tags, "formulas": formulas, "by_subject": by_subject}))
}

fn orphans(conn: &Connection) -> Result<Value> {
    let mut stmt = conn.prepare(
        "SELECT path, title FROM notes WHERE path NOT IN (SELECT DISTINCT source_path FROM links) 
         AND path NOT IN (SELECT DISTINCT target_path FROM links) AND type != 'moc'"
    )?;
    
    let results: Vec<Value> = stmt.query_map([], |r| {
        Ok(json!({"path": r.get::<_, String>(0)?, "title": r.get::<_, String>(1)?}))
    })?.collect::<Result<Vec<_>, _>>()?;
    
    Ok(json!({"count": results.len(), "orphans": results}))
}

fn formulas(conn: &Connection, q: &str, n: usize) -> Result<Value> {
    let mut stmt = conn.prepare("SELECT DISTINCT path, formula FROM formulas WHERE formula LIKE ? LIMIT ?")?;
    let pattern = format!("%{}%", q);
    let results: Vec<Value> = stmt.query_map([&pattern, &n.to_string()], |r| {
        Ok(json!({"path": r.get::<_, String>(0)?, "formula": r.get::<_, String>(1)?}))
    })?.collect::<Result<Vec<_>, _>>()?;
    
    Ok(json!({"query": q, "count": results.len(), "results": results}))
}
