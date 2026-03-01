-- BAC Unified D1 Database Schema (Simplified)

-- Questions table
CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question_text TEXT NOT NULL,
    solution_text TEXT,
    question_vector BLOB,
    topic_tags TEXT,
    difficulty INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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

-- Insert sample questions
INSERT INTO questions (question_text, solution_text, difficulty) VALUES 
    ('Résous: x + 3 = 7', 'x = 7 - 3 = 4', 1),
    ('Résous: 2x + 5 = 11', '2x = 11 - 5 = 6, x = 3', 2),
    ('Résous: x² + x - 2 = 0', 'Discriminant: 1 + 8 = 9, √9 = 3, x = (-1 ± 3)/2 = 1 ou -2', 3),
    ('Calcule la dérivée de f(x) = x³ + 2x²', "f''(x) = 3x² + 4x", 3),
    ('Résous: 3x - 7 = 2x + 1', '3x - 2x = 1 + 7, x = 8', 1);
