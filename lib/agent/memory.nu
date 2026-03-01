# ============================================================================
# MEMORY AGENT - Vector DB Lookup & Storage
# Uses Neon PostgreSQL with pgvector for semantic search
# ============================================================================

use config/database.yaml

# ============================================================================
# PUBLIC: Memory Lookup
# ============================================================================

export def memory-lookup [text: string, subject: string = "", top_k: int = 3] {
    print $"  🔍 Looking up similar problems in memory..."
    
    # Generate embedding for the query
    let embedding = (generate-embedding $text)
    
    if ($embedding == null) {
        print "  ⚠️ Could not generate embedding, using fallback"
        return (fallback-lookup $text $subject $top_k)
    }
    
    # Search in vector DB
    let results = (vector-search $embedding $top_k)
    
    if ($results == null or ($results | is-empty)) {
        print $"  📭 No similar problems found in memory"
        return {
            similar_problems: []
            similar_count: 0
            context: ""
        }
    }
    
    # Format results
    let similar = ($results | each {|r| {
        question: $r.question_text
        solution: $r.solution_text
        subject: $r.subject_name
        chapter: $r.chapter_name
        similarity: $r.similarity
    }})
    
    let context = ($similar | each {|p| 
        $"Problème similaire: " + $p.question + "\nSolution: " + $p.solution
    } | str join "\n\n")
    
    print $"  ✓ Found ($similar | length) similar problems"
    
    {
        similar_problems: $similar
        similar_count: ($similar | length)
        context: $context
    }
}

# ============================================================================
# PUBLIC: Store Problem in Memory
# ============================================================================

export def memory-store [problem: string, solution: string, subject: string, chapter: string = "", concepts: list<string> = []] {
    print $"  💾 Storing problem in memory..."
    
    # Generate embedding
    let embedding = (generate-embedding $problem + " " + $solution)
    
    if ($embedding == null) {
        print "  ⚠️ Could not generate embedding, skipping storage"
        return { stored: false }
    }
    
    # Store in database
    let result = (store-in-db $problem $solution $subject $chapter $concepts $embedding)
    
    if $result {
        print "  ✓ Problem stored in vector DB"
    } else {
        print "  ⚠️ Storage failed"
    }
    
    { stored: $result }
}

# ============================================================================
# HELPER: Generate Embedding via Ollama
# ============================================================================

def generate-embedding [text: string] {
    let truncated = if (($text | str length) > 8000) {
        ($text | str substring 0..8000)
    } else { $text }
    
    let body = $"{\"model\":\"nomic-embed-text\",\"prompt\":\"($truncated)\",\"stream\":false}"
    
    let out = (^curl -s -X POST "http://localhost:11434/api/embeddings" 
        -H "Content-Type:application/json" 
        -d $body)
    
    try {
        let j = ($out | from json)
        $j.embedding
    } catch {
        null
    }
}

# ============================================================================
# HELPER: Vector Similarity Search
# ============================================================================

def vector-search [embedding: list<float>, top_k: int] {
    # Convert embedding to PostgreSQL array format
    let embed_str = ($embedding | each {|v| $v | into string} | str join ",")
    let query = $"SELECT q.question_text, q.solution_text, s.name_fr as subject_name, 
        c.name_fr as chapter_name,
        1 - (q.embedding <=> '[($embed_str)]'::vector) as similarity
        FROM questions q
        LEFT JOIN subjects s ON q.subject_id = s.id
        LEFT JOIN chapters c ON q.chapter_id = c.id
        WHERE q.embedding IS NOT NULL AND q.verification_status = 'approved'
        ORDER BY q.embedding <=> '[($embed_str)]'::vector
        LIMIT ($top_k)"
    
    let result = (query-db $query)
    
    if ($result == null) {
        []
    } else { $result }
}

# ============================================================================
# HELPER: Store in Database
# ============================================================================

def store-in-db [problem: string, solution: string, subject: string, chapter: string, concepts: list<string>, embedding: list<float>] {
    # Get subject ID
    let subject_id = (get-subject-id $subject)
    
    # Get chapter ID (create if not exists)
    let chapter_id = (get-or-create-chapter $subject_id $chapter)
    
    # Convert embedding to array string
    let embed_str = ($embedding | each {|v| $v | into string} | str join ",")
    
    let query = $"INSERT INTO questions 
        (question_text, solution_text, subject_id, chapter_id, topic_tags, embedding, verification_status)
        VALUES ('" + ($problem | str escape) + "', '" + ($solution | str escape) + "', " + 
        $subject_id + ", " + $chapter_id + ", '{" + ($concepts | str join ",") + "}', 
        '[($embed_str)]'::vector, 'approved')
        ON CONFLICT (content_hash) DO NOTHING
        RETURNING id"
    
    let result = (query-db $query)
    $result != null
}

# ============================================================================
# HELPER: Fallback Lookup (text search)
# ============================================================================

def fallback-lookup [text: string, subject: string, top_k: int] {
    let subject_filter = if ($subject != "") {
        $" AND s.name_fr ILIKE '%" + $subject + "%'"
    } else { "" }
    
    let query = $"SELECT q.question_text, q.solution_text, s.name_fr as subject_name,
        c.name_fr as chapter_name
        FROM questions q
        LEFT JOIN subjects s ON q.subject_id = s.id
        LEFT JOIN chapters c ON q.chapter_id = c.id
        WHERE q.question_text ILIKE '%" + ($text | str escape) + "%' " + $subject_filter + "
        AND q.verification_status = 'approved'
        LIMIT " + ($top_k | into string)
    
    let results = (query-db $query)
    
    if ($results == null or ($results | is-empty)) {
        return {
            similar_problems: []
            similar_count: 0
            context: ""
        }
    }
    
    let similar = ($results | each {|r| {
        question: $r.question_text
        solution: $r.solution_text
        subject: $r.subject_name
        chapter: $r.chapter_name
        similarity: 0.5
    }})
    
    let context = ($similar | each {|p| 
        $"Problème similaire: " + $p.question + "\nSolution: " + $p.solution
    } | str join "\n\n")
    
    {
        similar_problems: $similar
        similar_count: ($similar | length)
        context: $context
    }
}

# ============================================================================
# HELPER: Get Subject ID
# ============================================================================

def get-subject-id [subject_code: string] {
    let code = match $subject_code {
        "math" => "math"
        "pc" | "physique" => "pc"
        "chimie" => "pc"
        "svt" => "svt"
        "philosophie" => "philosophie"
        "francais" | "français" => "francais"
        "arabe" => "arabe"
        "anglais" => "anglais"
        _ => "math"
    }
    
    let query = $"SELECT id FROM subjects WHERE code = '" + $code + "'"
    let result = (query-db $query)
    
    if ($result != null and ($result | length) > 0) {
        ($result.0.id | into string)
    } else { "1" }
}

# ============================================================================
# HELPER: Get or Create Chapter
# ============================================================================

def get-or-create-chapter [subject_id: int, chapter_name: string] {
    if ($chapter_name == "") {
        return "null"
    }
    
    let query = $"SELECT id FROM chapters WHERE subject_id = " + ($subject_id | into string) + " AND name_fr ILIKE '%" + ($chapter_name | str escape) + "%'"
    let result = (query-db $query)
    
    if ($result != null and ($result | length) > 0) {
        ($result.0.id | into string)
    } else {
        let insert = $"INSERT INTO chapters (subject_id, name_fr) VALUES (" + ($subject_id | into string) + ", '" + ($chapter_name | str escape) + "') RETURNING id"
        let insert_result = (query-db $insert)
        if ($insert_result != null and ($insert_result | length) > 0) {
            ($insert_result.0.id | into string)
        } else { "null" }
    }
}

# ============================================================================
# HELPER: Query Database
# ============================================================================

def query-db [query: string] {
    let db_url = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb?sslmode=require"
    
    let result = (^sh -c $"psql \""($db_url)\" -t -c \""($query)"\" 2>/dev/null" | complete)
    
    if $result.exit_code != 0 {
        return null
    }
    
    let output = $result.stdout
    
    if ($output | str trim) == "" {
        return []
    }
    
    # Parse PSQL output - it's pipe-separated
    let rows = ($output | split row "\n" | where {|r| ($r | str trim) != "" })
    
    $rows | each {|row|
        let cols = ($row | split column "|" | each {|c| $c | str trim })
        if ($cols | length) >= 2 {
            if ($cols.0 != "" and $cols.1 != "") {
                # This is a SELECT result - need to handle different column counts
                if ($cols | length) >= 4 {
                    {
                        question_text: $cols.0
                        solution_text: $cols.1
                        subject_name: (if ($cols | length) > 2 { $cols.2 } else { "" })
                        chapter_name: (if ($cols | length) > 3 { $cols.3 } else { "" })
                    }
                } else { null }
            } else { null }
        } else { null }
    } | where {|r| $r != null}
}
