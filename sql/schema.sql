-- ============================================================================
-- BAC UNIFIED - National Question Bank Database Schema
-- PostgreSQL + pgvector
-- Version: 2.0
-- ============================================================================

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- ============================================================================
-- USERS & AUTHENTICATION
-- ============================================================================

CREATE TYPE USER_ROLE AS ENUM ('student', 'teacher', 'admin', 'school_admin', 'moderator');
CREATE TYPE VERIFICATION_STATUS AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE SUBMISSION_SOURCE AS ENUM ('image', 'pdf', 'url', 'manual', 'historical', 'exam');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    
    -- Profile
    full_name TEXT,
    school TEXT,
    region TEXT,
    role USER_ROLE DEFAULT 'student',
    
    -- Gamification
    points INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    streak_days INTEGER DEFAULT 0,
    total_solved INTEGER DEFAULT 0,
    accuracy_rate FLOAT DEFAULT 0.0,
    
    -- Preferences
    preferred_language TEXT DEFAULT 'fr',
    preferred_subjects TEXT[],
    
    -- Security
    mfa_enabled BOOLEAN DEFAULT FALSE,
    last_login TIMESTAMP,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- User sessions
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    device_info JSONB,
    ip_address INET,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- SUBJECTS & CURRICULUM
-- ============================================================================

CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name_fr TEXT NOT NULL,
    name_ar TEXT NOT NULL,
    icon TEXT,
    color TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    order_index INTEGER DEFAULT 0
);

-- Insert default subjects
INSERT INTO subjects (code, name_fr, name_ar, color, order_index) VALUES
('math', 'Mathématiques', 'الرياضيات', '#E91E63', 1),
('pc', 'Physique-Chimie', 'الفيزياء والكيمياء', '#2196F3', 2),
('svt', 'Sciences de la Vie et de la Terre', 'علوم الحياة والأرض', '#4CAF50', 3),
('philosophie', 'Philosophie', 'الفلسفة', '#9C27B0', 4),
('francais', 'Français', 'الفرنسية', '#FF9800', 5),
('arabe', 'Arabe', 'العربية', '#F44336', 6),
('anglais', 'Anglais', 'الإنجليزية', '#00BCD4', 7)
ON CONFLICT (code) DO NOTHING;

-- Chapters per subject
CREATE TABLE chapters (
    id SERIAL PRIMARY KEY,
    subject_id INTEGER REFERENCES subjects(id) ON DELETE CASCADE,
    name_fr TEXT NOT NULL,
    name_ar TEXT NOT NULL,
    order_index INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);

-- ============================================================================
-- QUESTIONS BANK (Core Table)
-- ============================================================================

CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Submission info
    submitted_by UUID REFERENCES users(id),
    source_type SUBMISSION_SOURCE NOT NULL,
    original_filename TEXT,
    
    -- Content
    question_text TEXT NOT NULL,
    question_latex TEXT,
    question_image BYTEA,
    
    -- Classification
    subject_id INTEGER REFERENCES subjects(id),
    chapter_id INTEGER REFERENCES chapters(id),
    topic_tags TEXT[] DEFAULT '{}',
    difficulty INTEGER CHECK (difficulty BETWEEN 1 AND 5),
    question_type TEXT, -- 'multiple_choice', 'true_false', 'short_answer', 'proof', 'calculation'
    
    -- Solution
    solution_text TEXT,
    solution_latex TEXT,
    solution_steps JSONB DEFAULT '[]',
    correct_answer TEXT,
    answer_options JSONB DEFAULT '{}',
    
    -- Source reference
    source_exam TEXT,
    source_year INTEGER,
    source_session TEXT, -- 'principal', 'controle'
    source_school TEXT,
    
    -- Validation
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMP,
    verification_status VERIFICATION_STATUS DEFAULT 'pending',
    rejection_reason TEXT,
    
    -- Statistics
    submission_count INTEGER DEFAULT 0,
    correct_count INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    success_rate FLOAT DEFAULT 0.0,
    
    -- AI Analysis
    ai_difficulty FLOAT,
    ai_concepts JSONB DEFAULT '[]',
    ai_analysis JSONB DEFAULT '{}',
    embedding vector(1536),
    
    -- Content hash for deduplication
    content_hash TEXT NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(content_hash)
);

-- Question versions (audit trail)
CREATE TABLE question_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    diff JSONB,
    changed_by UUID REFERENCES users(id),
    change_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- NOON ANIMATIONS
-- ============================================================================

CREATE TABLE animations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    
    -- Animation info
    title TEXT,
    animation_type TEXT, -- 'graph', 'geometry', 'physics', 'chemistry'
    noon_code TEXT NOT NULL,
    generation_prompt TEXT,
    
    -- Output
    video_path TEXT,
    thumbnail_path TEXT,
    duration_seconds INTEGER,
    file_size_mb FLOAT,
    
    -- Status
    generation_status TEXT DEFAULT 'pending',
    rendered_at TIMESTAMP,
    
    -- Stats
    views INTEGER DEFAULT 0,
    rating FLOAT,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- PREDICTIONS
-- ============================================================================

CREATE TABLE predictions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Context
    subject_id INTEGER REFERENCES subjects(id),
    chapter_id INTEGER REFERENCES chapters(id),
    exam_year INTEGER NOT NULL,
    exam_session TEXT DEFAULT 'principal',
    
    -- Predicted content
    predicted_topics JSONB DEFAULT '[]',
    predicted_questions JSONB DEFAULT '[]',
    confidence_score FLOAT,
    
    -- Analysis basis
    based_on_patterns JSONB DEFAULT '{}',
    historical_years JSONB DEFAULT '[]',
    ml_model_version TEXT,
    
    -- Status
    status TEXT DEFAULT 'draft',
    published_at TIMESTAMP,
    accuracy_actual JSONB DEFAULT '{}', -- Filled after exam
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- USER PROGRESS & PRACTICE
-- ============================================================================

CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    question_id UUID REFERENCES questions(id) ON DELETE CASCADE,
    
    -- Attempt details
    attempt_number INTEGER NOT NULL,
    user_answer TEXT,
    is_correct BOOLEAN,
    time_taken_seconds INTEGER,
    
    -- Help usage
    hints_used INTEGER DEFAULT 0,
    solution_viewed BOOLEAN DEFAULT FALSE,
    
    -- AI feedback
    ai_feedback TEXT,
    weak_concepts TEXT[] DEFAULT '{}',
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(user_id, question_id, attempt_number)
);

-- Study sessions
CREATE TABLE study_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    session_type TEXT, -- 'practice', 'mock_exam', 'review', 'prediction'
    subject_id INTEGER REFERENCES subjects(id),
    
    -- Time tracking
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_minutes INTEGER,
    
    -- Results
    questions_attempted INTEGER DEFAULT 0,
    questions_correct INTEGER DEFAULT 0,
    score_percentage FLOAT DEFAULT 0.0,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- GAMIFICATION
-- ============================================================================

CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name_fr TEXT NOT NULL,
    name_ar TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    criteria JSONB NOT NULL,
    rarity TEXT DEFAULT 'common'
);

-- Insert default badges
INSERT INTO badges (code, name_fr, name_ar, description, rarity) VALUES
('first_step', 'Premier Pas', 'الخطوة الأولى', 'First submission', 'common'),
('scholar', 'Erudit', 'عالم', '100 questions submitted', 'rare'),
('solver', 'Résolveur', 'حلال', '500 problems solved', 'epic'),
('streak_7', 'Semaine Royale', 'أسبوع التحدي', '7 day streak', 'common'),
('streak_30', 'Maître de la Constance', 'سيد الالتزام', '30 day streak', 'rare'),
('perfect', 'Perfection', 'الكمال', '100% on mock exam', 'epic'),
('top_contributor', 'Champion National', 'البطل الوطني', 'Top monthly contributor', 'legendary'),
('mentor', 'Mentor', 'مرشد', 'Help 100 students', 'legendary')
ON CONFLICT (code) DO NOTHING;

-- User badges
CREATE TABLE user_badges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    badge_id INTEGER REFERENCES badges(id) ON DELETE CASCADE,
    earned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY(user_id, badge_id)
);

-- Point transactions
CREATE TABLE point_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    points INTEGER NOT NULL,
    action TEXT NOT NULL,
    reference_id UUID,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Leaderboard snapshots
CREATE TABLE leaderboards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    period TEXT NOT NULL, -- 'daily', 'weekly', 'monthly', 'all_time'
    subject_id INTEGER REFERENCES subjects(id),
    
    rankings JSONB NOT NULL,
    generated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(period, subject_id)
);

-- ============================================================================
-- CONTENT STORAGE (MinIO references)
-- ============================================================================

CREATE TABLE storage_objects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    bucket TEXT NOT NULL,
    object_key TEXT NOT NULL,
    filename TEXT,
    content_type TEXT,
    file_size BIGINT,
    checksum TEXT,
    
    -- Metadata
    uploaded_by UUID REFERENCES users(id),
    is_public BOOLEAN DEFAULT FALSE,
    
    -- Lifecycle
    storage_class TEXT DEFAULT 'STANDARD',
    expires_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(bucket, object_key)
);

-- ============================================================================
-- ANALYTICS
-- ============================================================================

CREATE TABLE analytics_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    event_type TEXT NOT NULL,
    user_id UUID REFERENCES users(id),
    
    -- Context
    subject_id INTEGER REFERENCES subjects(id),
    question_id UUID REFERENCES questions(id),
    
    -- Data
    properties JSONB DEFAULT '{}',
    session_id UUID,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Materialized view for daily stats
CREATE MATERIALIZED VIEW daily_stats AS
SELECT 
    DATE(created_at) as date,
    COUNT(DISTINCT user_id) as active_users,
    COUNT(*) as total_events,
    SUM(CASE WHEN event_type = 'question_attempted' THEN 1 ELSE 0 END) as questions_attempted,
    SUM(CASE WHEN event_type = 'submission_created' THEN 1 ELSE 0 END) as submissions,
    SUM(CASE WHEN event_type = 'correct_answer' THEN 1 ELSE 0 END) as correct_answers
FROM analytics_events
GROUP BY DATE(created_at);

-- ============================================================================
-- BLOCKCHAIN AUDIT LOG
-- ============================================================================

CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Transaction
    transaction_hash TEXT UNIQUE NOT NULL,
    previous_hash TEXT NOT NULL,
    block_number BIGINT,
    
    -- Action
    action TEXT NOT NULL,
    table_name TEXT NOT NULL,
    record_id UUID NOT NULL,
    
    -- Data
    changes JSONB,
    metadata JSONB DEFAULT '{}',
    
    -- Actor
    user_id UUID REFERENCES users(id),
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Vector similarity search
CREATE INDEX idx_questions_embedding ON questions 
    USING ivfflat (embedding vector_cosine_ops) 
    WITH (lists = 100);

-- Full-text search
CREATE INDEX idx_questions_text ON questions 
    USING gin (to_tsvector('french', question_text));

-- Trigram search
CREATE INDEX idx_questions_trgm ON questions 
    USING gin (question_text gin_trgm_ops);

-- Performance indexes
CREATE INDEX idx_user_progress_user ON user_progress(user_id);
CREATE INDEX idx_user_progress_question ON user_progress(question_id);
CREATE INDEX idx_questions_subject ON questions(subject_id, chapter_id);
CREATE INDEX idx_questions_verified ON questions(verification_status);
CREATE INDEX idx_questions_difficulty ON questions(difficulty);
CREATE INDEX idx_questions_source ON questions(source_year, source_session);
CREATE INDEX idx_predictions_exam ON predictions(exam_year, exam_session);
CREATE INDEX idx_audit_created ON audit_log(created_at DESC);
CREATE INDEX idx_analytics_events ON analytics_events(event_type, created_at);
CREATE INDEX idx_point_transactions_user ON point_transactions(user_id, created_at DESC);

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER questions_updated_at
    BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Update success rate after answer
CREATE OR REPLACE FUNCTION update_question_stats()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE questions 
    SET 
        submission_count = submission_count + 1,
        correct_count = correct_count + CASE WHEN NEW.is_correct THEN 1 ELSE 0 END,
        success_rate = (correct_count + CASE WHEN NEW.is_correct THEN 1 ELSE 0 END)::float / 
                       (submission_count + 1)::float
    WHERE id = NEW.question_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER progress_inserted
    AFTER INSERT ON user_progress
    FOR EACH ROW EXECUTE FUNCTION update_question_stats();

-- Generate content hash
CREATE OR REPLACE FUNCTION generate_content_hash()
RETURNS TRIGGER AS $$
BEGIN
    NEW.content_hash = md5(
        COALESCE(NEW.question_text, '') || 
        COALESCE(NEW.subject_id::TEXT, '') ||
        COALESCE(NEW.chapter_id::TEXT, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER questions_before_insert
    BEFORE INSERT ON questions
    FOR EACH ROW EXECUTE FUNCTION generate_content_hash();

-- ============================================================================
-- VIEWS
-- ============================================================================

-- Questions ready for practice
CREATE VIEW questions_for_practice AS
SELECT 
    q.id,
    q.question_text,
    q.solution_text,
    q.difficulty,
    q.success_rate,
    s.name_fr as subject_name,
    c.name_fr as chapter_name,
    q.created_at
FROM questions q
LEFT JOIN subjects s ON q.subject_id = s.id
LEFT JOIN chapters c ON q.chapter_id = c.id
WHERE q.verification_status = 'approved'
ORDER BY q.success_rate ASC, q.difficulty ASC;

-- Top contributors
CREATE VIEW top_contributors AS
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.school,
    u.points,
    u.level,
    u.streak_days,
    COUNT(DISTINCT q.id) as questions_submitted
FROM users u
LEFT JOIN questions q ON q.submitted_by = u.id
WHERE u.is_active = TRUE
GROUP BY u.id
ORDER BY u.points DESC
LIMIT 100;

-- ============================================================================
-- SEEDS
-- ============================================================================

-- Default admin user (password: admin123 - CHANGE IN PRODUCTION!)
INSERT INTO users (username, email, password_hash, role, full_name)
VALUES ('admin', 'admin@bac.edu.mr', '$2a$10$placeholder', 'admin', 'System Administrator')
ON CONFLICT (username) DO NOTHING;

-- Sample questions for testing
INSERT INTO questions (question_text, subject_id, difficulty, solution_text, verification_status)
VALUES 
    ('Résoudre: ∫(x² + 2x)dx de 0 à 2', 1, 2, 'F(x) = x³/3 + x², Résultat = 20/3', 'approved'),
    ('Calculer la dérivée de f(x) = x³ + 2x²', 1, 2, "f''(x) = 3x² + 4x", 'approved'),
    ('Un projectile lancé à 45° avec v=30m/s. Quelle est la portée?', 2, 3, '≈ 91.8m', 'approved')
ON CONFLICT (content_hash) DO NOTHING;
