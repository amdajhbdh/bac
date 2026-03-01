-- name: GetSimilarQuestions :many
-- Get similar questions based on vector embedding similarity
SELECT 
    question_text, 
    solution_text, 
    difficulty,
    1 - (question_vector <=> $2::vector) as similarity
FROM questions 
WHERE question_vector IS NOT NULL 
ORDER BY question_vector <=> $2::vector 
LIMIT $1;

-- name: GetSimilarQuestionsFiltered :many
-- Get similar questions with subject/chapter filters
SELECT 
    question_text, 
    solution_text, 
    difficulty,
    subject,
    chapter,
    1 - (question_vector <=> $2::vector) as similarity
FROM questions 
WHERE question_vector IS NOT NULL 
    AND ($3::text IS NULL OR subject = $3)
    AND ($4::text IS NULL OR chapter = $4)
ORDER BY question_vector <=> $2::vector 
LIMIT $1;

-- name: InsertQuestion :exec
-- Insert a new question with its solution and embedding
INSERT INTO questions (
    question_text, 
    solution_text, 
    question_vector, 
    topic_tags, 
    difficulty
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: SearchQuestionsByText :many
-- Search questions by text similarity
SELECT 
    question_text, 
    solution_text, 
    difficulty
FROM questions 
WHERE question_text ILIKE $1 
LIMIT $2;