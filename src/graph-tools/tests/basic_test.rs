//! Graph Tools - Basic Integration Tests
//!
//! Tests for all public API functions.

use graph_tools::{
    canvas, export, extract,
    models::{CanvasFile, Entities, RelationshipType},
    relationships, visualize,
};

/// Test canvas generation
#[test]
fn test_generate_canvas() {
    let canvas = canvas::generate_canvas("Physics");

    assert!(!canvas.nodes.is_empty());
    assert_eq!(canvas.nodes[0].text, "Physics");
    assert_eq!(canvas.edges.len(), 0);
}

/// Test adding concept nodes
#[test]
fn test_add_concept_node() {
    let mut canvas = canvas::generate_canvas("Test");
    let id = canvas::add_concept_node(&mut canvas, "New Concept", 100.0, 200.0, None);

    assert_eq!(canvas.nodes.len(), 2);
    assert!(canvas.nodes[1]
        .id
        .starts_with(&id[..id - 1.min(id.len() - 1)]));
}

/// Test adding edges
#[test]
fn test_add_edge() {
    let mut canvas = canvas::generate_canvas("Test");
    let subject_id = canvas.nodes[0].id.clone();
    let concept_id = canvas::add_concept_node(&mut canvas, "Child", 100.0, 200.0, None);

    canvas::add_edge(
        &mut canvas,
        subject_id,
        concept_id,
        graph_tools::models::EdgeType::Arrow,
    );

    assert_eq!(canvas.edges.len(), 1);
    assert_eq!(canvas.edges[0].from_node, subject_id);
    assert_eq!(canvas.edges[0].to_node, concept_id);
}

/// Test entity extraction from content
#[test]
fn test_extract_from_content() {
    let content = r#"## Electric Field
Force:: A push or pull on an object
### Coulomb's Law"#;

    let entities = extract::extract_from_content(content, "physics.md");

    assert!(entities.concepts.len() >= 2); // "Electric Field" and "Coulomb's Law"
    assert!(entities.definitions.len() >= 1); // Force definition
}

/// Test building relationships from entities
#[test]
fn test_build_relationships() {
    let entities = Entities::new();
    entities.add_concept(graph_tools::models::Concept::new("Concept A", "note.md"));
    entities.add_concept(graph_tools::models::Concept::new("Concept B", "note.md"));

    let canvas = relationships::build_from_entities(&entities);

    assert_eq!(canvas.nodes.len(), 2);
}

/// Test JSON export
#[test]
fn test_export_json() {
    let canvas = canvas::generate_canvas("Test Subject");
    let json = export::export_json(&canvas, "Test Subject");

    assert!(!json.is_empty());
    assert!(json.contains("Test Subject"));
    assert!(json.contains("nodes"));
    assert!(json.contains("edges"));
    assert!(json.contains("metadata"));
}

/// Test compact JSON export
#[test]
fn test_export_json_compact() {
    let canvas = canvas::generate_canvas("Compact Test");
    let json = export::export_json_compact(&canvas, "Compact Test");

    assert!(!json.contains("\n"));
    assert!(json.contains("Compact Test"));
}

/// Test DOT export
#[test]
fn test_export_dot() {
    let mut canvas = canvas::generate_canvas("A");
    let a_id = canvas.nodes[0].id.clone();
    let b_id = canvas::add_concept_node(&mut canvas, "B", 100.0, 0.0, None);
    canvas::add_edge(
        &mut canvas,
        a_id,
        b_id,
        graph_tools::models::EdgeType::Arrow,
    );

    let dot = export::export_dot(&canvas, "TestGraph");

    assert!(dot.contains("digraph TestGraph"));
    assert!(dot.contains("A"));
    assert!(dot.contains("B"));
    assert!(dot.contains("->"));
}

/// Test visualization URL generation
#[test]
fn test_generate_visualization() {
    let viz = visualize::generate_visualization("Physics");

    assert!(viz.url.contains("mermaid"));
    assert!(viz.embed_url.contains("mermaid"));
    assert!(viz.graph_data_url.is_some());
}

/// Test Mermaid diagram generation
#[test]
fn test_to_mermaid() {
    let mut canvas = canvas::generate_canvas("Main");
    let main_id = canvas.nodes[0].id.clone();
    let sub_id = canvas::add_concept_node(&mut canvas, "Sub", 0.0, 0.0, None);
    canvas::add_edge(
        &mut canvas,
        main_id,
        sub_id,
        graph_tools::models::EdgeType::Arrow,
    );

    let mermaid = visualize::to_mermaid(&canvas);

    assert!(mermaid.contains("flowchart"));
    assert!(mermaid.contains("Main"));
    assert!(mermaid.contains("Sub"));
    assert!(mermaid.contains("-->"));
}

/// Test D3 JSON export
#[test]
fn test_to_d3_json() {
    let canvas = canvas::generate_canvas("D3 Test");
    let d3 = visualize::to_d3_json(&canvas);

    assert!(d3.contains("nodes"));
    assert!(d3.contains("links"));
    assert!(d3.contains("D3 Test"));
}

/// Test parsing JSON back to GraphJson
#[test]
fn test_parse_roundtrip() {
    let canvas = canvas::generate_canvas("Roundtrip");
    let json = export::export_json(&canvas, "Roundtrip");

    let parsed = export::parse_graph_json(&json);
    assert!(parsed.is_ok());

    let graph = parsed.unwrap();
    assert_eq!(graph.metadata.subject, "Roundtrip");
    assert_eq!(graph.nodes.len(), 1);
}

/// Test library initialization
#[test]
fn test_init() {
    graph_tools::init();
    // If we get here without panicking, init works
}

/// Test canvas file operations
#[test]
fn test_canvas_operations() {
    let mut canvas = CanvasFile::new();

    // Add nodes
    let id1 = canvas::add_concept_node(&mut canvas, "Node 1", 0.0, 0.0, None);
    let id2 = canvas::add_concept_node(&mut canvas, "Node 2", 100.0, 100.0, None);

    assert_eq!(canvas.nodes.len(), 2);

    // Add labeled edge
    canvas::add_labeled_edge(
        &mut canvas,
        id1.clone(),
        id2.clone(),
        graph_tools::models::EdgeType::Arrow,
        "depends on",
    );

    assert_eq!(canvas.edges.len(), 1);
    assert_eq!(canvas.edges[0].label.as_deref(), Some("depends on"));
}

/// Test entity creation
#[test]
fn test_entity_creation() {
    let concept = graph_tools::models::Concept::new("Electric Field", "physics.md");
    assert_eq!(concept.name, "Electric Field");
    assert_eq!(concept.source_note, "physics.md");
    assert_eq!(concept.mentions, 1);

    let definition = graph_tools::models::Definition::new("Force", "A push or pull", "physics.md");
    assert_eq!(definition.term, "Force");
    assert_eq!(definition.definition, "A push or pull");
}

/// Test relationship creation
#[test]
fn test_relationship_creation() {
    let rel = graph_tools::models::EntityRelationship::new(
        "id1".to_string(),
        "id2".to_string(),
        RelationshipType::Prerequisite,
    );
    assert_eq!(rel.source_id, "id1");
    assert_eq!(rel.target_id, "id2");
    assert_eq!(rel.weight, 1.0);
}

/// Test canvas JSON parsing
#[test]
fn test_parse_canvas_json() {
    let canvas = canvas::generate_canvas("Parse Test");
    let json = serde_json::to_string(&canvas).unwrap();

    let parsed = export::parse_canvas_json(&json);
    assert!(parsed.is_ok());
    assert_eq!(parsed.unwrap().nodes.len(), 1);
}

/// Test weight calculation
#[test]
fn test_calculate_weights() {
    let mut rels = vec![
        graph_tools::models::EntityRelationship::new(
            "a".to_string(),
            "b".to_string(),
            RelationshipType::Related,
        ),
        graph_tools::models::EntityRelationship::new(
            "a".to_string(),
            "b".to_string(),
            RelationshipType::Related,
        ),
    ];

    relationships::calculate_weights(&mut rels);

    assert_eq!(rels[0].weight, 2.0);
    assert_eq!(rels[1].weight, 2.0);
}

/// Test empty entities handling
#[test]
fn test_empty_entities() {
    let entities = Entities::new();
    let canvas = relationships::build_from_entities(&entities);

    assert!(canvas.nodes.is_empty());
    assert!(canvas.edges.is_empty());
}
