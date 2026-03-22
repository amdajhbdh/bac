//! Export Module
//!
//! Exports knowledge graphs in various formats.

use crate::models::{CanvasFile, GraphJson};
use serde_json;

/// Export graph as JSON
///
/// Converts the current canvas/graph to JSON format
///
/// # Arguments
/// * `canvas` - The canvas to export
/// * `subject` - The subject for metadata
///
/// # Returns
/// JSON string representation of the graph
pub fn export_json(canvas: &CanvasFile, subject: &str) -> String {
    let graph = GraphJson::from_canvas(canvas, subject);
    serde_json::to_string_pretty(&graph).unwrap_or_default()
}

/// Export to compact JSON (no pretty printing)
pub fn export_json_compact(canvas: &CanvasFile, subject: &str) -> String {
    let graph = GraphJson::from_canvas(canvas, subject);
    serde_json::to_string(&graph).unwrap_or_default()
}

/// Export entities to JSON
pub fn export_entities_json(entities: &crate::models::Entities) -> String {
    serde_json::to_string_pretty(entities).unwrap_or_default()
}

/// Parse JSON back to GraphJson
///
/// # Arguments
/// * `json` - JSON string
///
/// # Returns
/// Parsed GraphJson or error
pub fn parse_graph_json(json: &str) -> Result<GraphJson, serde_json::Error> {
    serde_json::from_str(json)
}

/// Parse JSON to CanvasFile
///
/// # Arguments
/// * `json` - JSON string
///
/// # Returns
/// Parsed CanvasFile or error
pub fn parse_canvas_json(json: &str) -> Result<CanvasFile, serde_json::Error> {
    serde_json::from_str(json)
}

/// Export as DOT format for Graphviz
///
/// # Arguments
/// * `canvas` - The canvas to export
/// * `subject` - Graph name
///
/// # Returns
/// DOT format string
pub fn export_dot(canvas: &CanvasFile, subject: &str) -> String {
    let mut dot = format!("digraph {} {{\n", escape_dot_id(subject));
    dot.push_str("    rankdir=LR;\n");
    dot.push_str("    node [shape=box];\n\n");

    // Nodes
    for node in &canvas.nodes {
        let label = escape_dot_label(&node.text);
        dot.push_str(&format!(
            "    {} [label=\"{}\"];\n",
            escape_dot_id(&node.id),
            label
        ));
    }

    // Edges
    for edge in &canvas.edges {
        let edge_type = match edge.edge_type {
            crate::models::EdgeType::Arrow => "[arrowhead=normal]",
            crate::models::EdgeType::Dashed => "[style=dashed]",
            crate::models::EdgeType::Custom => "",
        };

        if let Some(label) = &edge.label {
            dot.push_str(&format!(
                "    {} -> {} {} [label=\"{}\"];\n",
                escape_dot_id(&edge.from_node),
                escape_dot_id(&edge.to_node),
                edge_type,
                escape_dot_label(label)
            ));
        } else {
            dot.push_str(&format!(
                "    {} -> {} {};\n",
                escape_dot_id(&edge.from_node),
                escape_dot_id(&edge.to_node),
                edge_type
            ));
        }
    }

    dot.push_str("}\n");
    dot
}

/// Escape string for DOT ID
fn escape_dot_id(s: &str) -> String {
    format!("n_{}", s.replace("-", "_").replace(":", "_"))
}

/// Escape string for DOT label
fn escape_dot_label(s: &str) -> String {
    s.replace("\\", "\\\\")
        .replace("\"", "\\\"")
        .replace("\n", "\\n")
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::canvas::{add_concept_node, add_edge};
    use crate::models::EdgeType;

    #[test]
    fn test_export_json() {
        let canvas = crate::canvas::generate_canvas("Test Subject");
        let json = export_json(&canvas, "Test Subject");

        assert!(!json.is_empty());
        assert!(json.contains("Test Subject"));
        assert!(json.contains("nodes"));
        assert!(json.contains("edges"));
    }

    #[test]
    fn test_export_dot() {
        let mut canvas = crate::canvas::generate_canvas("A");
        let a_id = canvas.nodes[0].id.clone();
        let b_id = add_concept_node(&mut canvas, "B", 100.0, 0.0, None);
        add_edge(&mut canvas, a_id, b_id, EdgeType::Arrow);

        let dot = export_dot(&canvas, "Test");

        assert!(dot.contains("digraph n_Test"));
        assert!(dot.contains("A"));
        assert!(dot.contains("B"));
        assert!(dot.contains("->"));
    }

    #[test]
    fn test_parse_roundtrip() {
        let canvas = crate::canvas::generate_canvas("Roundtrip");
        let json = export_json(&canvas, "Roundtrip");

        let parsed: GraphJson = parse_graph_json(&json).unwrap();
        assert_eq!(parsed.metadata.subject, "Roundtrip");
        assert_eq!(parsed.nodes.len(), 1);
    }
}
