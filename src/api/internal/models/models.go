package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Username     string     `gorm:"uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex" json:"email"`
	PasswordHash string     `gorm:"not null" json:"-"`
	FullName     string     `json:"full_name"`
	School       string     `json:"school"`
	Region       string     `json:"region"`
	Role         string     `gorm:"default:'student'" json:"role"`
	Points       int        `gorm:"default:0" json:"points"`
	Level        int        `gorm:"default:1" json:"level"`
	StreakDays   int        `gorm:"default:0" json:"streak_days"`
	LastActive   *time.Time `json:"last_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Question struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	SubmittedBy        *uuid.UUID `gorm:"type:uuid" json:"submitted_by"`
	SourceType         string     `gorm:"not null" json:"source_type"`
	QuestionText       string     `json:"question_text"`
	QuestionLatex      string     `json:"question_latex"`
	SubjectID          int        `json:"subject_id"`
	ChapterID          int        `json:"chapter_id"`
	TopicTags          []string   `gorm:"-" json:"topic_tags"`
	Difficulty         int        `json:"difficulty"`
	QuestionType       string     `json:"question_type"`
	SolutionText       string     `json:"solution_text"`
	SolutionLatex      string     `json:"solution_latex"`
	SolutionSteps      string     `gorm:"type:text" json:"solution_steps"`
	CorrectAnswer      string     `json:"correct_answer"`
	AnswerOptions      string     `gorm:"type:text" json:"answer_options"`
	SourceExam         string     `json:"source_exam"`
	SourceYear         int        `json:"source_year"`
	VerifiedBy         *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	VerificationStatus string     `gorm:"default:'pending'" json:"verification_status"`
	SubmissionCount    int        `gorm:"default:0" json:"submission_count"`
	CorrectCount       int        `gorm:"default:0" json:"correct_count"`
	ViewCount          int        `gorm:"default:0" json:"view_count"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type ManimVideo struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	QuestionID       uuid.UUID `gorm:"type:uuid" json:"question_id"`
	Title            string    `json:"title"`
	AnimationType    string    `json:"animation_type"`
	DurationSeconds  int       `json:"duration_seconds"`
	FileSizeMB       float64   `json:"file_size_mb"`
	VideoURL         string    `json:"video_url"`
	ThumbnailURL     string    `json:"thumbnail_url"`
	GenerationPrompt string    `json:"generation_prompt"`
	GenerationStatus string    `json:"generation_status"`
	Views            int       `json:"views"`
	CreatedAt        time.Time `json:"created_at"`
}

type Prediction struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	SubjectID          int       `json:"subject_id"`
	ChapterID          int       `json:"chapter_id"`
	ExamYear           int       `json:"exam_year"`
	ExamSession        string    `json:"exam_session"`
	PredictedTopics    string    `gorm:"type:text" json:"predicted_topics"`
	PredictedQuestions string    `gorm:"type:text" json:"predicted_questions"`
	ConfidenceScore    float64   `json:"confidence_score"`
	Status             string    `gorm:"default:'draft'" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
}

type UserProgress struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	QuestionID       uuid.UUID `gorm:"type:uuid;index" json:"question_id"`
	AttemptNumber    int       `json:"attempt_number"`
	UserAnswer       string    `json:"user_answer"`
	IsCorrect        bool      `json:"is_correct"`
	TimeTakenSeconds int       `json:"time_taken_seconds"`
	HintsUsed        int       `json:"hints_used"`
	SolutionViewed   bool      `json:"solution_viewed"`
	CreatedAt        time.Time `json:"created_at"`
}

type Badge struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"uniqueIndex" json:"code"`
	NameFR      string `json:"name_fr"`
	NameAR      string `json:"name_ar"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Criteria    string `gorm:"type:text" json:"criteria"`
	Rarity      string `json:"rarity"`
}

type UserBadge struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid" json:"user_id"`
	BadgeID  int       `json:"badge_id"`
	EarnedAt time.Time `json:"earned_at"`
}

type Subject struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Code   string `gorm:"uniqueIndex" json:"code"`
	NameFR string `json:"name_fr"`
	NameAR string `json:"name_ar"`
	Icon   string `json:"icon"`
	Color  string `json:"color"`
}
