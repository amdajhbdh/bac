//! Entity Extraction Module
//!
//! Extracts entities (concepts, definitions, relationships) from notes.

use crate::models::{Concept, Definition, Entities, EntityRelationship, RelationshipType};
use std::collections::HashMap;

/// Extract entities from all notes
///
/// Scans note content and extracts:
/// - Concepts (capitalized terms, headers)
/// - Definitions (term: definition patterns)
/// - Implicit relationships (proximity-based)
///
/// # Returns
/// A collection of extracted `Entities`
pub fn extract_entities_from_notes() -> Entities {
    // This is a simplified implementation
    // In production, this would read from the vault
    Entities::new()
}

/// Extract entities from a single note
///
/// # Arguments
/// * `content` - The note content
/// * `source_note` - The source file name
///
/// # Returns
/// Extracted entities from the note
pub fn extract_from_content(content: &str, source_note: &str) -> Entities {
    let mut entities = Entities::new();

    // Extract concepts from # headers (## and ###)
    for line in content.lines() {
        let trimmed = line.trim();
        if trimmed.starts_with("## ") {
            let concept_name = trimmed.trim_start_matches("## ").trim();
            entities.add_concept(Concept::new(concept_name, source_note));
        } else if trimmed.starts_with("### ") {
            let concept_name = trimmed.trim_start_matches("### ").trim();
            entities.add_concept(Concept::new(concept_name, source_note));
        }
    }

    // Extract definitions (term:: definition pattern)
    for line in content.lines() {
        if let Some((term, def)) = extract_definition(line) {
            entities.add_definition(Definition::new(&term, &def, source_note));
        }
    }

    entities
}

/// Extract definition from a line
///
/// Matches patterns like: `term:: definition text`
fn extract_definition(line: &str) -> Option<(String, String)> {
    let trimmed = line.trim();
    if let Some(pos) = trimmed.find("::") {
        let term = trimmed[..pos].trim().to_string();
        let def = trimmed[pos + 2..].trim().to_string();
        if !term.is_empty() && !def.is_empty() {
            return Some((term, def));
        }
    }
    None
}

/// Extract concepts using pattern matching
///
/// Matches capitalized multi-word terms
pub fn extract_concepts(content: &str, source_note: &str) -> Vec<Concept> {
    let mut concepts: Vec<Concept> = Vec::new();
    let mut seen: HashMap<String, usize> = HashMap::new();

    // Pattern for capitalized terms (at least 2 words)
    let pattern = regex::Regex::new(r"\b([A-Z][a-z]+(?:\s+[A-Z][a-z]+)+)\b").unwrap();

    for cap in pattern.find_iter(content) {
        let term = cap.as_str();
        let count = seen.entry(term.to_string()).or_insert(0);
        *count += 1;
    }

    for (term, mentions) in seen {
        let mut concept = Concept::new(&term, source_note);
        concept.mentions = mentions;
        concepts.push(concept);
    }

    concepts
}

/// Build relationship map from concepts
///
/// Finds related concepts based on co-occurrence in same section
pub fn find_related_concepts(entities: &Entities) -> Vec<EntityRelationship> {
    let mut relationships = Vec::new();

    // Group concepts by source note
    let mut by_source: HashMap<&str, Vec<&Concept>> = HashMap::new();
    for concept in &entities.concepts {
        by_source
            .entry(&concept.source_note)
            .or_default()
            .push(concept);
    }

    // Create relationships between concepts in same note
    for concepts in by_source.values() {
        if concepts.len() >= 2 {
            for i in 0..concepts.len() {
                for j in (i + 1)..concepts.len() {
                    relationships.push(EntityRelationship::new(
                        concepts[i].id.clone(),
                        concepts[j].id.clone(),
                        RelationshipType::Related,
                    ));
                }
            }
        }
    }

    relationships
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_extract_definition() {
        let result = extract_definition("Electric Field:: A region around a charged particle");
        assert!(result.is_some());
        let (term, def) = result.unwrap();
        assert_eq!(term, "Electric Field");
        assert_eq!(def, "A region around a charged particle");
    }

    #[test]
    fn test_extract_from_content() {
        let content = r#"## Electric Field
Force:: A push or pull on an object
### Coulomb's Law"#;
        let entities = extract_from_content(content, "physics.md");

        assert_eq!(entities.concepts.len(), 2);
        assert_eq!(entities.definitions.len(), 1);
    }
}
