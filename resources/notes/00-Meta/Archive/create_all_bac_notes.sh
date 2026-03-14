#!/bin/bash

VAULT="/home/med/Documents/bac/resources/notes/Study-Vault"
SOURCE="$VAULT/Extracted/General/7c_bac_pc_كامل.md"

# Array of line ranges for each exam
declare -A exams=(
    ["2002-Normal"]="65:507"
    ["2002-Complementaire"]="508:831"
    ["2003-Normal"]="832:1333"
    ["2003-Complementaire"]="1334:1756"
    ["2004-Normal"]="1757:2197"
    ["2004-Complementaire"]="2198:2613"
    ["2005-Normal"]="2614:3121"
    ["2005-Complementaire"]="3122:3515"
    ["2006-Normal"]="3516:3871"
    ["2006-Complementaire"]="3872:4261"
    ["2007-Normal"]="4262:4634"
    ["2007-Complementaire"]="4635:5035"
    ["2008-Normal"]="5036:5368"
    ["2008-Complementaire"]="5369:5748"
    ["2009-Normal"]="5749:6124"
    ["2009-Complementaire"]="6125:6474"
    ["2010-Normal"]="6475:6790"
    ["2010-Complementaire"]="6791:7147"
    ["2011-Normal"]="7148:7420"
    ["2011-Complementaire"]="7421:7773"
    ["2012-Normal"]="7774:8082"
    ["2012-Complementaire"]="8083:8346"
)

for exam in "${!exams[@]}"; do
    year="${exam%-*}"
    session="${exam#*-}"
    session_tag=$(echo "$session" | tr '[:upper:]' '[:lower:]')
    range="${exams[$exam]}"
    start="${range%:*}"
    end="${range#*:}"
    
    filename="$VAULT/BAC-${exam}.md"
    
    cat > "$filename" << HEADER
---
tags: [bac, mauritania, $year, ${session_tag}-session, physics, chemistry, exam, 7C]
year: $year
session: $session
subjects: [physics-chemistry]
duration: 4h
coefficient: 8/4
---

# BAC $year - Session $session - Sciences Physiques

**Série:** Mathématiques/T.M.G.M  
**Durée:** 4 heures  
**Coefficient:** 8/4

---

HEADER

    # Extract content
    sed -n "${start},${end}p" "$SOURCE" >> "$filename"
    
    # Add footer with links
    cat >> "$filename" << FOOTER

---

## Liens connexes

FOOTER

    # Add year-based links
    prev_year=$((year - 1))
    next_year=$((year + 1))
    
    if [ "$session" = "Normal" ]; then
        echo "- [[BAC-${year}-Complementaire]] - Session complémentaire" >> "$filename"
    else
        echo "- [[BAC-${year}-Normal]] - Session normale" >> "$filename"
    fi
    
    echo "- [[BAC-${prev_year}-${session}]] - Année précédente" >> "$filename"
    echo "- [[BAC-${next_year}-${session}]] - Année suivante" >> "$filename"
    
    cat >> "$filename" << FOOTER2

## Tags thématiques

#redox #cinétique #acide-base #pH #organique #mécanique #énergie #oscillations #électromagnétisme #induction #circuits-AC #RLC #résonance

---

**Source:** Institut Pédagogique National (IPN) - Mauritanie
FOOTER2

    echo "Created: $filename"
done

echo ""
echo "All 22 BAC notes created successfully!"
