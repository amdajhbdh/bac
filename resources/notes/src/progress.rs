use anyhow::Result;
use chrono::Utc;
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
    Session { subject: String, minutes: i32, correct: i32, incorrect: i32 },
    Stats { #[arg(short, default_value = "7")] days: i32 },
    Streak,
}

fn main() -> Result<()> {
    let cli = Cli::parse();
    let conn = Connection::open(".progress.db")?;
    init(&conn)?;
    
    let output = match cli.cmd {
        Cmd::Session { subject, minutes, correct, incorrect } => session(&conn, &subject, minutes, correct, incorrect)?,
        Cmd::Stats { days } => stats(&conn, days)?,
        Cmd::Streak => streak(&conn)?,
    };
    
    println!("{}", serde_json::to_string_pretty(&output)?);
    Ok(())
}

fn init(conn: &Connection) -> Result<()> {
    conn.execute(
        "CREATE TABLE IF NOT EXISTS sessions (
            id INTEGER PRIMARY KEY, subject TEXT, date TEXT, duration_minutes INTEGER,
            correct INTEGER, incorrect INTEGER
        )", [])?;
    Ok(())
}

fn session(conn: &Connection, subject: &str, minutes: i32, correct: i32, incorrect: i32) -> Result<Value> {
    let now = Utc::now().to_rfc3339();
    conn.execute("INSERT INTO sessions (subject, date, duration_minutes, correct, incorrect) VALUES (?, ?, ?, ?, ?)",
        [subject, &now, &minutes.to_string(), &correct.to_string(), &incorrect.to_string()])?;
    
    let accuracy = if correct + incorrect > 0 { (correct as f64 / (correct + incorrect) as f64) * 100.0 } else { 0.0 };
    Ok(json!({"subject": subject, "minutes": minutes, "correct": correct, "incorrect": incorrect, "accuracy": accuracy}))
}

fn stats(conn: &Connection, days: i32) -> Result<Value> {
    let cutoff = (Utc::now() - chrono::Duration::days(days as i64)).to_rfc3339();
    
    let mut stmt = conn.prepare("SELECT subject, SUM(duration_minutes), SUM(correct), SUM(incorrect) FROM sessions WHERE date >= ? GROUP BY subject")?;
    let by_subject: Vec<Value> = stmt.query_map([&cutoff], |r| {
        let correct: i32 = r.get(2)?;
        let incorrect: i32 = r.get(3)?;
        let accuracy = if correct + incorrect > 0 { (correct as f64 / (correct + incorrect) as f64) * 100.0 } else { 0.0 };
        Ok(json!({"subject": r.get::<_, String>(0)?, "minutes": r.get::<_, i32>(1)?, "correct": correct, "incorrect": incorrect, "accuracy": accuracy}))
    })?.collect::<Result<Vec<_>, _>>()?;
    
    let total_mins: i32 = conn.query_row("SELECT SUM(duration_minutes) FROM sessions WHERE date >= ?", [&cutoff], |r| r.get(0))?;
    Ok(json!({"period_days": days, "total_minutes": total_mins, "by_subject": by_subject}))
}

fn streak(conn: &Connection) -> Result<Value> {
    let mut stmt = conn.prepare("SELECT DISTINCT DATE(date) as day FROM sessions ORDER BY day DESC")?;
    let days: Vec<String> = stmt.query_map([], |r| r.get(0))?.collect::<Result<Vec<_>, _>>()?;
    
    let mut streak = 0;
    let today = Utc::now().format("%Y-%m-%d").to_string();
    
    for (i, day) in days.iter().enumerate() {
        let expected = (Utc::now() - chrono::Duration::days(i as i64)).format("%Y-%m-%d").to_string();
        if day == &expected { streak += 1; } else { break; }
    }
    
    Ok(json!({"current_streak": streak, "last_study": days.first()}))
}
