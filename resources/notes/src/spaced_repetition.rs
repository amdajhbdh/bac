use anyhow::Result;
use chrono::{DateTime, Duration, Utc};
use clap::{Parser, Subcommand};
use rusqlite::Connection;
use serde::{Deserialize, Serialize};

#[derive(Parser)]
#[command(name = "spaced-repetition")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Add { path: String },
    Review { path: String, quality: u8 },
    Due { #[arg(short, long, default_value = "20")] limit: usize },
    Stats,
}

#[derive(Serialize, Deserialize)]
struct Card {
    note_path: String,
    easiness: f64,
    interval: i32,
    repetitions: i32,
    due_date: String,
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    let db_path = ".spaced-repetition.db";
    init_db(db_path)?;
    
    match cli.command {
        Commands::Add { path } => add_card(&path, db_path)?,
        Commands::Review { path, quality } => review_card(&path, quality, db_path)?,
        Commands::Due { limit } => show_due(limit, db_path)?,
        Commands::Stats => show_stats(db_path)?,
    }
    
    Ok(())
}

fn init_db(db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    
    conn.execute(
        "CREATE TABLE IF NOT EXISTS cards (
            id INTEGER PRIMARY KEY,
            note_path TEXT UNIQUE,
            easiness REAL DEFAULT 2.5,
            interval INTEGER DEFAULT 1,
            repetitions INTEGER DEFAULT 0,
            due_date TEXT,
            last_review TEXT
        )",
        [],
    )?;
    
    conn.execute(
        "CREATE TABLE IF NOT EXISTS reviews (
            id INTEGER PRIMARY KEY,
            card_id INTEGER,
            quality INTEGER,
            timestamp TEXT,
            FOREIGN KEY(card_id) REFERENCES cards(id)
        )",
        [],
    )?;
    
    Ok(())
}

fn add_card(path: &str, db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    let now = Utc::now().to_rfc3339();
    
    conn.execute(
        "INSERT OR IGNORE INTO cards (note_path, due_date, last_review)
         VALUES (?, ?, ?)",
        [path, &now, &now],
    )?;
    
    println!("Added: {}", path);
    Ok(())
}

fn review_card(path: &str, quality: u8, db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    
    let (id, mut ef, mut interval, mut reps): (i64, f64, i32, i32) = conn.query_row(
        "SELECT id, easiness, interval, repetitions FROM cards WHERE note_path = ?",
        [path],
        |row| Ok((row.get(0)?, row.get(1)?, row.get(2)?, row.get(3)?)),
    )?;
    
    // SM-2 algorithm
    if quality >= 3 {
        interval = match reps {
            0 => 1,
            1 => 6,
            _ => (interval as f64 * ef) as i32,
        };
        reps += 1;
    } else {
        reps = 0;
        interval = 1;
    }
    
    ef = ef + (0.1 - (5 - quality as i32) as f64 * (0.08 + (5 - quality as i32) as f64 * 0.02));
    ef = ef.max(1.3);
    
    let next_review = Utc::now() + Duration::days(interval as i64);
    
    conn.execute(
        "UPDATE cards SET easiness=?, interval=?, repetitions=?, due_date=?, last_review=?
         WHERE id=?",
        [&ef.to_string(), &interval.to_string(), &reps.to_string(), 
         &next_review.to_rfc3339(), &Utc::now().to_rfc3339(), &id.to_string()],
    )?;
    
    conn.execute(
        "INSERT INTO reviews (card_id, quality, timestamp) VALUES (?, ?, ?)",
        [&id.to_string(), &quality.to_string(), &Utc::now().to_rfc3339()],
    )?;
    
    println!("✓ Next review: {} ({} days)", next_review.format("%Y-%m-%d"), interval);
    println!("  Easiness: {:.2}", ef);
    
    Ok(())
}

fn show_due(limit: usize, db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    let now = Utc::now().to_rfc3339();
    
    let mut stmt = conn.prepare(
        "SELECT note_path, due_date, interval, repetitions FROM cards 
         WHERE due_date <= ? ORDER BY due_date LIMIT ?"
    )?;
    
    let cards: Vec<(String, String, i32, i32)> = stmt
        .query_map([&now, &limit.to_string()], |row| {
            Ok((row.get(0)?, row.get(1)?, row.get(2)?, row.get(3)?))
        })?
        .collect::<Result<Vec<_>, _>>()?;
    
    println!("\n📚 {} cards due for review:\n", cards.len());
    for (path, due, interval, reps) in cards {
        println!("  • {}", path);
        println!("    Due: {} | Interval: {} days | Reps: {}\n", due, interval, reps);
    }
    
    Ok(())
}

fn show_stats(db_path: &str) -> Result<()> {
    let conn = Connection::open(db_path)?;
    
    let total: i64 = conn.query_row("SELECT COUNT(*) FROM cards", [], |row| row.get(0))?;
    let now = Utc::now().to_rfc3339();
    let due: i64 = conn.query_row(
        "SELECT COUNT(*) FROM cards WHERE due_date <= ?",
        [&now],
        |row| row.get(0),
    )?;
    
    let (avg_ef, avg_int): (f64, f64) = conn.query_row(
        "SELECT AVG(easiness), AVG(interval) FROM cards",
        [],
        |row| Ok((row.get(0).unwrap_or(2.5), row.get(1).unwrap_or(1.0))),
    )?;
    
    println!("\n📊 Spaced Repetition Stats:\n");
    println!("  Total cards: {}", total);
    println!("  Due today: {}", due);
    println!("  Avg easiness: {:.2}", avg_ef);
    println!("  Avg interval: {:.1} days", avg_int);
    
    Ok(())
}
