//! Canvas Generation Module
//!
//! Generates Obsidian .canvas files from subjects and concepts.

use crate::models::{CanvasEdge, CanvasFile, CanvasNode, EdgeType, GraphId};
use uuid::Uuid;

/// Generate a canvas file for a given subject
///
/// # Arguments
/// * `subject` - The main subject/topic for the canvas
///
/// # Returns
/// A new `CanvasFile` with basic structure
pub fn generate_canvas(subject: &str) -> CanvasFile {
    let mut canvas = CanvasFile::new();

    // Add subject as central node
    let subject_id = Uuid::new_v4().to_string();
    let subject_node = CanvasNode::concept(subject_id.clone(), 500.0, 300.0, subject);
    canvas.add_node(subject_node);

    canvas
}

/// Add a concept node to an existing canvas
///
/// # Arguments
/// * `canvas` - The canvas to add to
/// * `concept` - The concept text
/// * `x`, `y` - Position coordinates
/// * `color` - Optional color
///
/// # Returns
/// The ID of the created node
pub fn add_concept_node(
    canvas: &mut CanvasFile,
    concept: &str,
    x: f64,
    y: f64,
    color: Option<&str>,
) -> GraphId {
    let id = Uuid::new_v4().to_string();
    let mut node = CanvasNode::concept(id.clone(), x, y, concept);
    if let Some(c) = color {
        node.color = Some(c.to_string());
    }
    canvas.add_node(node)
}

/// Add an edge connecting two nodes
///
/// # Arguments
/// * `canvas` - The canvas to add to
/// * `from_id` - Source node ID
/// * `to_id` - Target node ID
/// * `edge_type` - Type of edge
pub fn add_edge(canvas: &mut CanvasFile, from_id: GraphId, to_id: GraphId, edge_type: EdgeType) {
    let edge = CanvasEdge::new(Uuid::new_v4().to_string(), from_id, to_id, edge_type);
    canvas.add_edge(edge);
}

/// Add an edge with label
///
/// # Arguments
/// * `canvas` - The canvas to add to
/// * `from_id` - Source node ID
/// * `to_id` - Target node ID
/// * `edge_type` - Type of edge
/// * `label` - Edge label text
pub fn add_labeled_edge(
    canvas: &mut CanvasFile,
    from_id: GraphId,
    to_id: GraphId,
    edge_type: EdgeType,
    label: &str,
) {
    let edge = CanvasEdge::with_label(Uuid::new_v4().to_string(), from_id, to_id, edge_type, label);
    canvas.add_edge(edge);
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_generate_canvas() {
        let canvas = generate_canvas("Physics");
        assert_eq!(canvas.nodes.len(), 1);
        assert_eq!(canvas.edges.len(), 0);
        assert_eq!(canvas.nodes[0].text, "Physics");
    }

    #[test]
    fn test_add_concept_node() {
        let mut canvas = generate_canvas("Test");
        let id = add_concept_node(&mut canvas, "New Concept", 100.0, 200.0, None);

        assert_eq!(canvas.nodes.len(), 2);
        assert_eq!(canvas.nodes[1].text, "New Concept");
        assert!(canvas.nodes[1].id.starts_with(&id[..]));
    }

    #[test]
    fn test_add_edge() {
        let mut canvas = generate_canvas("Test");
        let subject_id = canvas.nodes[0].id.clone();
        let concept_id = add_concept_node(&mut canvas, "Child", 100.0, 200.0, None);

        add_edge(&mut canvas, subject_id, concept_id, EdgeType::Arrow);

        assert_eq!(canvas.edges.len(), 1);
        assert_eq!(canvas.edges[0].from_node, subject_id);
        assert_eq!(canvas.edges[0].to_node, concept_id);
    }
}
