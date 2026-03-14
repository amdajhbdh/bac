#!/bin/bash

EXTRACTED_DIR="/home/med/Documents/bac/resources/notes/Study-Vault/Extracted"
ANALYSIS_DIR="/home/med/Documents/bac/resources/notes/Study-Vault/Analysis"

mkdir -p "$ANALYSIS_DIR"/{patterns,topics,schools,exercises,formulas}

echo "=== PATTERN EXTRACTION AGENT ==="
echo "Analyzing 63 extracted files..."

# Agent 1: Topic Extractor
echo "[Agent 1] Extracting topics and keywords..."
for category in Math Physics Chemistry Biology French General; do
    if [ -d "$EXTRACTED_DIR/$category" ]; then
        grep -h "^#" "$EXTRACTED_DIR/$category"/*.md 2>/dev/null | \
        grep -v "^# " | sort | uniq -c | sort -rn > "$ANALYSIS_DIR/topics/${category}_topics.txt"
    fi
done

# Agent 2: School & Year Analyzer
echo "[Agent 2] Analyzing schools and academic years..."
grep -rh "LYCÉE\|E\.P\|lycée" "$EXTRACTED_DIR" --include="*.md" | \
    grep -oE "(LYCÉE|E\.P|lycée)[^0-9]*[0-9]{4}" | \
    sort | uniq -c | sort -rn > "$ANALYSIS_DIR/schools/schools_frequency.txt"

# Agent 3: Exercise Pattern Detector
echo "[Agent 3] Detecting exercise patterns..."
{
    echo "# Exercise Patterns Detected"
    echo ""
    echo "## Exercise Keywords:"
    grep -rh "exercice\|Exercice\|EXERCICE" "$EXTRACTED_DIR" --include="*.md" | wc -l
    echo ""
    echo "## Problem Types:"
    grep -rh "Calculer\|Démontrer\|Déterminer\|Résoudre\|Montrer" "$EXTRACTED_DIR" --include="*.md" | \
        grep -oE "Calculer|Démontrer|Déterminer|Résoudre|Montrer" | sort | uniq -c | sort -rn
} > "$ANALYSIS_DIR/exercises/exercise_patterns.txt"

# Agent 4: Formula Extractor
echo "[Agent 4] Extracting mathematical formulas and equations..."
{
    echo "# Mathematical Patterns"
    echo ""
    echo "## Common Symbols:"
    grep -roh "[∫∑∏√∂∆∇≈≠≤≥±×÷∞]" "$EXTRACTED_DIR" --include="*.md" | sort | uniq -c | sort -rn
    echo ""
    echo "## Variables & Constants:"
    grep -roh "\b[a-z]\b\|[α-ω]" "$EXTRACTED_DIR" --include="*.md" | sort | uniq -c | sort -rn | head -30
} > "$ANALYSIS_DIR/formulas/math_symbols.txt"

# Agent 5: Subject Coverage Analyzer
echo "[Agent 5] Analyzing subject coverage..."
{
    echo "# Subject Coverage Analysis"
    echo ""
    for category in Math Physics Chemistry Biology; do
        echo "## $category"
        count=$(find "$EXTRACTED_DIR/$category" -name "*.md" 2>/dev/null | wc -l)
        echo "Files: $count"
        
        # Extract key topics
        case $category in
            Math)
                echo "Topics detected:"
                grep -rh "nombre complexe\|intégrale\|matrice\|arithmétique\|suite" "$EXTRACTED_DIR/$category" --include="*.md" | \
                    grep -oE "nombre complexe|intégrale|matrice|arithmétique|suite" | sort | uniq -c | sort -rn
                ;;
            Physics)
                echo "Topics detected:"
                grep -rh "dynamique\|projectile\|ressort\|champ électrique\|induction" "$EXTRACTED_DIR/$category" --include="*.md" | \
                    grep -oE "dynamique|projectile|ressort|champ électrique|induction" | sort | uniq -c | sort -rn
                ;;
            Chemistry)
                echo "Topics detected:"
                grep -rh "organique\|cinétique\|alcool\|réaction" "$EXTRACTED_DIR/$category" --include="*.md" | \
                    grep -oE "organique|cinétique|alcool|réaction" | sort | uniq -c | sort -rn
                ;;
            Biology)
                echo "Topics detected:"
                grep -rh "génétique\|hormone\|nerveux\|brassage" "$EXTRACTED_DIR/$category" --include="*.md" | \
                    grep -oE "génétique|hormone|nerveux|brassage" | sort | uniq -c | sort -rn
                ;;
        esac
        echo ""
    done
} > "$ANALYSIS_DIR/patterns/subject_coverage.txt"

# Agent 6: Difficulty Level Detector
echo "[Agent 6] Detecting difficulty levels..."
{
    echo "# Difficulty Indicators"
    echo ""
    echo "## By Class Level:"
    grep -roh "7[CD]" "$EXTRACTED_DIR" --include="*.md" | sort | uniq -c | sort -rn
    echo ""
    echo "## Complexity Markers:"
    echo "Devoir (Tests): $(grep -rc "devoir" "$EXTRACTED_DIR" --include="*.md" | awk -F: '{sum+=$2} END {print sum}')"
    echo "Série (Practice): $(grep -rc "série" "$EXTRACTED_DIR" --include="*.md" | awk -F: '{sum+=$2} END {print sum}')"
    echo "Corrigé (Solutions): $(grep -rc "corrigé\|CORRIGÉ" "$EXTRACTED_DIR" --include="*.md" | awk -F: '{sum+=$2} END {print sum}')"
} > "$ANALYSIS_DIR/patterns/difficulty_levels.txt"

echo ""
echo "=== ANALYSIS COMPLETE ==="
echo "Results saved to: $ANALYSIS_DIR"
