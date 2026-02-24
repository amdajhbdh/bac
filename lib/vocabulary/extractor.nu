# ============================================================================
# FRENCH VOCABULARY EXTRACTOR
# Extracts French words from BAC exam PDFs for Arabic translation
# ============================================================================

# ============================================================================
# MAIN VOCABULARY EXTRACTOR
# ============================================================================

export def extract-vocabulary [
    --source: string = "pdfs"  # pdfs, database, text
    --subject: string = ""      # math, pc, svt, chimie
    --min-frequency: int = 2   # minimum occurrences
] {
    print "%(char-esc)[1;36m📚 Extracting French Vocabulary...%(char-reset)"
    
    let words = match $source {
        "pdfs" => (extract-from-pdfs $subject $min-frequency)
        "database" => (extract-from-database $subject $min-frequency)
        _ => (extract-from-text $source $min-frequency)
    }
    
    print $"Extracted: ($words | length) unique words"
    
    $words
}

# ============================================================================
# EXTRACT FROM PDFS
# ============================================================================

def extract-from-pdfs [ min_freq: intsubject: string,] {
    let pdf_dir = "db/pdfs"
    
    # Filter PDFs by subject if specified
    let subject_filter = if ($subject != "") {
        ls $pdf_dir | where name =~ $subject
    } else {
        ls $pdf_dir
    }
    
    let all_text = ""
    
    # Process each PDF
    for pdf in $subject_filter {
        print $"Processing: ($pdf.name)"
        let text = (pdf-to-text $pdf.name)
        $all_text = $all_text + " " + $text
    }
    
    # Extract French words
    extract-french-words $all_text $min_freq
}

# ============================================================================
# PDF TO TEXT
# ============================================================================

def pdf-to-text [file: string] {
    let filepath = ($file | path join "db/pdfs" $file)
    
    if not ($filepath | path exists) {
        return ""
    }
    
    let output = $"/tmp/($file | path parse | get stem).txt"
    
    let result = (^pdftotext -layout $filepath $output 2>/dev/null)
    
    if ($output | path exists) {
        open $output
    } else {
        ""
    }
}

# ============================================================================
# EXTRACT FRENCH WORDS
# ============================================================================

def extract-french-words [text: string, min_freq: int] {
    # Common French words to filter out
    let french_stopwords = [
        "le", "la", "les", "un", "une", "des", "de", "du", "des"
        "et", "est", "sont", "etait", "etaient", "etre", "avoir"
        "je", "tu", "il", "elle", "nous", "vous", "ils", "elles"
        "que", "qui", "quoi", "dont", "où", "quand", "comment"
        "pour", "dans", "sur", "avec", "sans", "par", "chez"
        "ce", "cette", "ces", "celui", "celle", "ceux", "celles"
        "plus", "moins", "tres", "bien", "mal", "peut", "doit"
        "fait", "faire", "fait", "fait", "dire", "dit", "dit"
        "a", "au", "aux", "ou", "mais", "donc", "car", "ni"
        "ne", "pas", "plus", "jamais", "toujours", "aussi"
        "lorsque", "apres", "avant", "pendant", "entre", "sous"
        "se", "soi", "lui", "leur", "moi", "toi", "nous", "vous"
    ]
    
    # Extract words (French/academic vocabulary)
    let words = ($text | split row -r '\W+')
    
    # Count frequencies
    let word_counts = {}
    for word in $words {
        let w = ($word | str downcase | str trim)
        if ($w | str length > 3 and not ($french_stopwords | any {|s| $s == $w })) {
            if $w in $word_counts {
                $word_counts[$w] = ($word_counts[$w] + 1)
            } else {
                $word_counts[$w] = 1
            }
        }
    }
    
    # Filter by frequency and sort
    let vocab = ($word_counts | items | where {|_, count| $count >= $min_freq} | sort-by -r {|_, b| b} | take 100)
    
    $vocab
}

# ============================================================================
# EXTRACT FROM DATABASE
# ============================================================================

def extract-from-database [subject: string, min_freq: int] {
    let db_url = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"
    
    let subject_filter = if ($subject != "") {
        $"WHERE subject = '($subject)'"
    } else {
        ""
    }
    
    let sql = $"SELECT content FROM problems ($subject_filter)"
    
    let result = (^psql $db_url -t -c $sql 2>/dev/null)
    
    if ($result == null) {
        return []
    }
    
    let all_text = ($result | str join " ")
    extract-french-words $all_text $min_freq
}

# ============================================================================
# EXTRACT FROM TEXT
# ============================================================================

def extract-from-text [text: string, min_freq: int] {
    extract-french-words $text $min_freq
}

# ============================================================================
# TRANSLATE TO ARABIC (Google Translate)
# ============================================================================

export def translate-arabic [word: string] {
    # Use Google Translate web (free, no API key)
    let encoded = ($word | str replace -a " " "+")
    let url = $"https://translate.google.com/?sl=fr&tl=ar&text=($encoded)&op=translate"
    
    # Would need Playwright for actual translation
    # For now, return placeholder
    {
        french: $word
        arabic: "[translation requires Playwright]"
        source: "google-translate"
    }
}

# ============================================================================
# GET VOCABULARY WITH TRANSLATIONS
# ============================================================================

export def get-vocabulary [
    --subject: string = ""
    --limit: int = 50
] {
    # Extract words
    let words = (extract-vocabulary --source "pdfs" --subject $subject --min-frequency 2)
    
    # Take top words
    let top_words = ($words | take $limit)
    
    # Return with placeholders (translation would use Playwright)
    $top_words | each {|item|
        {
            french: $item.0
            frequency: $item.1
            arabic: "[translate]"
        }
    }
}

# ============================================================================
# SAVE VOCABULARY TO DATABASE
# ============================================================================

export def save-vocabulary [words: list] {
    let db_url = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb"
    
    # Create table if not exists
    let create_sql = "CREATE TABLE IF NOT EXISTS vocabulary (
        id SERIAL PRIMARY KEY,
        french_word VARCHAR(255) NOT NULL,
        arabic_translation TEXT,
        subject VARCHAR(50),
        frequency INT DEFAULT 1,
        created_at TIMESTAMP DEFAULT NOW()
    )"
    
    ^psql $db_url -c $create_sql 2>/dev/null
    
    # Insert words
    for word in $words {
        let french = $word.french
        let arabic = $word.arabic
        let freq = $word.frequency
        
        let sql = $"INSERT INTO vocabulary (french_word, arabic_translation, frequency) 
                   VALUES ('($french)', '($arabic)', ($freq)) 
                   ON CONFLICT (french_word) DO UPDATE SET frequency = vocabulary.frequency + 1"
        
        ^psql $db_url -c $sql 2>/dev/null
    }
    
    print $"Saved ($words | length) words to vocabulary table"
}
