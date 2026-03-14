-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetNotesByUser :many
SELECT * FROM notes 
WHERE user_id = $1 
ORDER BY updated_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateNote :one
INSERT INTO notes (user_id, title, content, subject, chapter, tags)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateNote :one
UPDATE notes 
SET title = $2, content = $3, subject = $4, chapter = $5, tags = $6, updated_at = NOW()
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1 AND user_id = $2;

-- name: GetCardsByUser :many
SELECT * FROM cards 
WHERE user_id = $1 AND next_review <= NOW()
ORDER BY next_review ASC
LIMIT 20;

-- name: CreateCard :one
INSERT INTO cards (user_id, note_id, front, back)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCardReview :exec
UPDATE cards 
SET ease_factor = $2,
    interval_days = $3,
    next_review = NOW() + (interval '1 day' * $4),
    reviews = reviews + 1
WHERE id = $1;

-- name: GetResourcesBySubject :many
SELECT * FROM resources 
WHERE subject = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: CreateStudySession :one
INSERT INTO study_sessions (user_id, subject)
VALUES ($1, $2)
RETURNING *;

-- name: EndStudySession :one
UPDATE study_sessions
SET ended_at = NOW(),
    duration_minutes = EXTRACT(EPOCH FROM (NOW() - started_at))::INTEGER / 60,
    notes_studied = $2
WHERE id = $1
RETURNING *;
