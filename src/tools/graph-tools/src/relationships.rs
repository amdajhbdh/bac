//! Relationships Module
//!
//! Builds relationship graphs between extracted entities.

use crate::models::{
    CanvasEdge, CanvasFile, CanvasNode, EdgeType, Entities, EntityRelationship, GraphId,
    RelationshipType,
};
use std::collections::HashMap;
use uuid::Uuid;

/// Build relationships from extracted entities
///
/// Creates a canvas file with nodes for concepts and edges for relationships
///
/// # Returns
/// A `CanvasFile` representing the relationship graph
pub fn build_relationships() -> CanvasFile {
    // Get entities first
    let entities = crate::extract::extract_entities_from_notes();
    build_from_entities(&entities)
}

/// Build relationship canvas from entities
///
/// # Arguments
/// * `entities` - Extracted entities to visualize
///
/// # Returns
/// Canvas file with nodes and edges
pub fn build_from_entities(entities: &Entities) -> CanvasFile {
    let mut canvas = CanvasFile::new();
    let mut node_ids: HashMap<String, GraphId> = HashMap::new();

    // Add concept nodes in a grid layout
    let cols = 3;
    let spacing_x = 450.0;
    let spacing_y = 200.0;
    let start_x = 100.0;
    let start_y = 100.0;

    for (i, concept) in entities.concepts.iter().enumerate() {
        let x = start_x + (i as f64 % cols as f64) * spacing_x;
        let y = start_y + (i as f64 / cols as f64) * spacing_y;

        let node = CanvasNode::concept(concept.id.clone(), x, y, &concept.name);
        node_ids.insert(concept.name.clone(), concept.id.clone());
        canvas.add_node(node);
    }

    // Add definition nodes
    let def_start_y =
        start_y + ((entities.concepts.len() as f64 / cols as f64).ceil() + 1.0) * spacing_y;
    for (i, def) in entities.definitions.iter().enumerate() {
        let x = start_x + (i as f64 % cols as f64) * spacing_x;
        let y = def_start_y + (i as f64 / cols as f64) * spacing_y;

        let node = CanvasNode::concept(
            def.id.clone(),
            x,
            y,
            &format!("{}: {}", def.term, def.definition),
        );
        node_ids.insert(def.term.clone(), def.id.clone());
        canvas.add_node(node);
    }

    // Add edges from relationships
    for rel in &entities.relationships {
        let edge_type = match rel.relationship_type {
            RelationshipType::Prerequisite => EdgeType::Arrow,
            RelationshipType::PartOf => EdgeType::Dashed,
            RelationshipType::Related => EdgeType::Arrow,
            RelationshipType::DependsOn => EdgeType::Arrow,
            RelationshipType::Example => EdgeType::Dashed,
            RelationshipType::Contrasts => EdgeType::Dashed,
        };

        let canvas_edge = CanvasEdge::new(
            Uuid::new_v4().to_string(),
            rel.source_id.clone(),
            rel.target_id.clone(),
            edge_type,
        );
        canvas.add_edge(canvas_edge);
    }

    canvas
}

/// Infer relationships between concepts based on text analysis
///
/// Analyzes concept co-occurrence and proximity to infer relationships
pub fn infer_relationships(entities: &Entities) -> Vec<EntityRelationship> {
    let mut inferred: Vec<EntityRelationship> = Vec::new();

    // Check for prerequisite patterns
    for def in &entities.definitions {
        // If definition mentions another concept, it's a prerequisite
        for concept in &entities.concepts {
            if def.definition.contains(&concept.name) && def.term != concept.name {
                inferred.push(EntityRelationship::new(
                    concept.id.clone(),
                    def.id.clone(),
                    RelationshipType::Prerequisite,
                ));
            }
        }
    }

    inferred
}

/// Calculate relationship weights based on frequency
pub fn calculate_weights(relationships: &mut [EntityRelationship]) {
    let mut counts: HashMap<(String, String), usize> = HashMap::new();

    for rel in relationships.iter() {
        let key = (rel.source_id.clone(), rel.target_id.clone());
        let count = counts.entry(key).or_insert(0);
        *count += 1;
    }

    for rel in relationships.iter_mut() {
        let key = (rel.source_id.clone(), rel.target_id.clone());
        if let Some(&count) = counts.get(&key) {
            rel.weight = count as f32;
        }
    }
}

/// Get all unique relationship types in use
pub fn get_relationship_types(entities: &Entities) -> Vec<RelationshipType> {
    let mut types: Vec<RelationshipType> = entities
        .relationships
        .iter()
        .map(|r| r.relationship_type.clone())
        .collect();

    types.sort_by(|a, b| format!("{:?}", a).cmp(&format!("{:?}", b)));
    types.dedup();
    types
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_build_from_empty_entities() {
        let entities = Entities::new();
        let canvas = build_from_entities(&entities);

        assert!(canvas.nodes.is_empty());
        assert!(canvas.edges.is_empty());
    }

    #[test]
    fn test_calculate_weights() {
        let mut rels = vec![
            EntityRelationship::new("a".to_string(), "b".to_string(), RelationshipType::Related),
            EntityRelationship::new("a".to_string(), "b".to_string(), RelationshipType::Related),
        ];

        calculate_weights(&mut rels);

        assert_eq!(rels[0].weight, 2.0);
        assert_eq!(rels[1].weight, 2.0);
    }
}
