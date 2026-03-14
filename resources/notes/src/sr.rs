use anyhow::Result;
use chrono::{DateTime, Duration, Utc};
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
    Add { path: String },
    Review { path: String, quality: u8 },
    Due { #[arg(short, default_value = "20")] n: usize },
    Stats,
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    let conn = Connection::open(".spaced-repetition.db")?;
    init(&conn)?;
    
    let output = match cli.cmd {
        Cmd::Add { path } => add(&conn, &path)?,
        Cmd::Review { path, quality } => review(&conn, &path, quality)?,
        Cmd::Due { n } => due(&conn, n)?,
        Cmd::Stats => stats(&conn)?,
    };
    
    println!("{}", serde_json::to_string_pretty(&output)?);
    Ok(())
}

fn init(conn: &Connection) -> Result<()> {
    conn.execute(
        "CREATE TABLE IF NOT EXISTS cards (
            id INTEGER PRIMARY KEY, note_path TEXT UNIQUE, easiness REAL DEFAULT 2.5,
            interval INTEGER DEFAULT 1, repetitions INTEGER DEFAULT 0,
            due_date TEXT, last_review TEXT
        )", [])?;
    
    conn.execute(
        "CREATE TABLE IF NOT EXISTS reviews (
            id INTEGER PRIMARY KEY, card_id INTEGER, quality INTEGER, timestamp TEXT,
            FOREIGN KEY(card_id) REFERENCES cards(id)
        )", [])?;
    Ok(())
}

fn add(conn: &Connection, path: &str) -> Result<Value> {
    let now = Utc::now().to_rfc3339();
    conn.execute("INSERT OR IGNORE INTO cards (note_path, due_date, last_review) VALUES (?, ?, ?)",
        [path, &now, &now])?;
    Ok(json!({"action": "added", "path": path}))
}

fn review(conn: &Connection, path: &str, quality: u8) -> Result<Value> {
    let (id, mut ef, mut interval, mut reps): (i64, f64, i32, i32) = conn.query_row(
        "SELECT id, easiness, interval, repetitions FROM cards WHERE note_path = ?",
        [path], |r| Ok((r.get(0)?, r.get(1)?, r.get(2)?, r.get(3)?)))?;
    
    if quality >= 3 {
        interval = match reps { 0 => 1, 1 => 6, _ => (interval as f64 * ef) as i32 };
        reps += 1;
    } else {
        reps = 0;
        interval = 1;
    }
    
    ef = (ef + (0.1 - (5 - quality as i32) as f64 * (0.08 + (5 - quality as i32) as f64 * 0.02))).max(1.3);
    let next = Utc::now() + Duration::days(interval as i64);
    
    conn.execute("UPDATE cards SET easiness=?, interval=?, repetitions=?, due_date=?, last_review=? WHERE id=?",
        [&ef.to_string(), &interval.to_string(), &reps.to_string(), &next.to_rfc3339(), &Utc::now().to_rfc3339(), &id.to_string()])?;
    
    conn.execute("INSERT INTO reviews (card_id, quality, timestamp) VALUES (?, ?, ?)",
        [&id.to_string(), &quality.to_string(), &Utc::now().to_rfc3339()])?;
    
    Ok(json!({"path": path, "next_review": next.to_rfc3339(), "interval_days": interval, "easiness": ef, "repetitions": reps}))
}

fn due(conn: &Connection, n: usize) -> Result<Value> {
    let mut stmt = conn.prepare("SELECT note_path, due_date, interval, repetitions FROM cards WHERE due_date <= ? ORDER BY due_date LIMIT ?")?;
    let cards: Vec<Value> = stmt.query_map([&Utc::now().to_rfc3339(), &n.to_string()], |r| {
        Ok(json!({"path": r.get::<_, String>(0)?, "due": r.get::<_, String>(1)?, "interval": r.get::<_, i32>(2)?, "reps": r.get::<_, i32>(3)?}))
    })?.collect::<Result<Vec<_>, _>>()?;
    
    Ok(json!({"count": cards.len(), "cards": cards}))
}

fn stats(conn: &Connection) -> Result<Value> {
    let total: i64 = conn.query_row("SELECT COUNT(*) FROM cards", [], |r| r.get(0))?;
    let due: i64 = conn.query_row("SELECT COUNT(*) FROM cards WHERE due_date <= ?", [&Utc::now().to_rfc3339()], |r| r.get(0))?;
    let (avg_ef, avg_int): (f64, f64) = conn.query_row("SELECT AVG(easiness), AVG(interval) FROM cards", [], 
        |r| Ok((r.get(0).unwrap_or(2.5), r.get(1).unwrap_or(1.0))))?;
    
    Ok(json!({"total": total, "due": due, "avg_easiness": avg_ef, "avg_interval_days": avg_int}))
}
