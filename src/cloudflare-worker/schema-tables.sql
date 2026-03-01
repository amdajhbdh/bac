-- Add missing tables

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    points INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Predictions table
CREATE TABLE IF NOT EXISTS predictions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_id INTEGER,
    probability REAL,
    confidence REAL,
    predicted_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User progress table
CREATE TABLE IF NOT EXISTS user_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    question_id INTEGER,
    attempted INTEGER DEFAULT 0,
    correct INTEGER DEFAULT 0,
    last_attempt TIMESTAMP,
    time_spent_seconds INTEGER
);

-- Sample users
INSERT OR IGNORE INTO users (name, email, points) VALUES 
    ('Ahmed', 'ahmed@bac.mr', 150),
    ('Fatima', 'fatima@bac.mr', 120),
    ('Mohamed', 'mohamed@bac.mr', 100);
