//! Visualization Module
//!
//! Generates visualization URLs for knowledge graphs.

use crate::models::{CanvasFile, VisualizationUrl};

/// Generate visualization URL for a subject
///
/// Creates URLs to various graph visualization services
///
/// # Arguments
/// * `subject` - The subject to visualize
///
/// # Returns
/// Visualization URLs for different services
pub fn generate_visualization(subject: &str) -> VisualizationUrl {
    // Encode subject for URL
    let encoded = urlencoding::encode(subject);

    // Mermaid Live Editor URL
    let mermaid_url = format!("https://mermaid.live/edit#gist/{}", encoded);

    // Graphviz Online URL
    let graphviz_url = format!("https://dreampuf.github.io/GraphvizOnline/#{}", encoded);

    VisualizationUrl {
        url: format!("{}?v=1", mermaid_url),
        embed_url: format!("https://mermaid.live/embed/{}", encoded),
        graph_data_url: Some(graphviz_url),
    }
}

/// Generate Mermaid diagram code from canvas
///
/// Converts canvas to Mermaid flowchart syntax
pub fn to_mermaid(canvas: &CanvasFile) -> String {
    let mut mermaid = String::from("flowchart TD\n");

    // Add nodes
    for node in &canvas.nodes {
        let node_type = match node.node_type {
            crate::models::NodeType::Group => ":::group",
            crate::models::NodeType::Image => ":::image",
            _ => "",
        };
        let safe_id = sanitize_mermaid_id(&node.id);
        let safe_text = sanitize_mermaid_text(&node.text);
        mermaid.push_str(&format!(
            "    {}{}[\"{}\"]\n",
            safe_id, node_type, safe_text
        ));
    }

    // Add edges
    for edge in &canvas.edges {
        let safe_from = sanitize_mermaid_id(&edge.from_node);
        let safe_to = sanitize_mermaid_id(&edge.to_node);

        let arrow = match edge.edge_type {
            crate::models::EdgeType::Arrow => "-->",
            crate::models::EdgeType::Dashed => "-.->",
            crate::models::EdgeType::Custom => "---",
        };

        if let Some(label) = &edge.label {
            let safe_label = sanitize_mermaid_text(label);
            mermaid.push_str(&format!(
                "    {} {} {} |\"{}\"| {}\n",
                safe_from, arrow, safe_to, safe_label, safe_to
            ));
        } else {
            mermaid.push_str(&format!("    {} {} {}\n", safe_from, arrow, safe_to));
        }
    }

    // Add styling
    mermaid.push_str("\n    classDef concept fill:#f9f,stroke:#333,stroke-width:2px;\n");
    mermaid.push_str("    classDef group fill:#bbf,stroke:#333,stroke-width:4px;\n");

    mermaid
}

/// Sanitize string for Mermaid ID
fn sanitize_mermaid_id(s: &str) -> String {
    s.replace("-", "_")
        .replace(":", "_")
        .replace(".", "_")
        .chars()
        .filter(|c| c.is_alphanumeric() || *c == '_')
        .collect()
}

/// Sanitize string for Mermaid text
fn sanitize_mermaid_text(s: &str) -> String {
    s.replace("\\", "\\\\")
        .replace("\"", "'")
        .replace("\n", "<br />")
}

/// Generate D3.js compatible JSON for custom visualization
pub fn to_d3_json(canvas: &CanvasFile) -> String {
    let nodes: Vec<serde_json::Value> = canvas
        .nodes
        .iter()
        .map(|n| {
            serde_json::json!({
                "id": n.id,
                "name": n.text,
                "group": match n.node_type {
                    crate::models::NodeType::Group => "group",
                    crate::models::NodeType::Image => "image",
                    _ => "concept"
                }
            })
        })
        .collect();

    let links: Vec<serde_json::Value> = canvas
        .edges
        .iter()
        .map(|e| {
            serde_json::json!({
                "source": e.from_node,
                "target": e.to_node,
                "type": format!("{:?}", e.edge_type).to_lowercase()
            })
        })
        .collect();

    let graph = serde_json::json!({
        "nodes": nodes,
        "links": links
    });

    serde_json::to_string_pretty(&graph).unwrap_or_default()
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::canvas::{add_concept_node, add_edge};
    use crate::models::EdgeType;

    #[test]
    fn test_generate_visualization() {
        let viz = generate_visualization("Physics");

        assert!(viz.url.contains("mermaid"));
        assert!(viz.url.contains("Physics"));
        assert!(viz.embed_url.contains("mermaid"));
        assert!(viz.graph_data_url.is_some());
    }

    #[test]
    fn test_to_mermaid() {
        let mut canvas = crate::canvas::generate_canvas("Main");
        let main_id = canvas.nodes[0].id.clone();
        let sub_id = add_concept_node(&mut canvas, "Sub", 0.0, 0.0, None);
        add_edge(&mut canvas, main_id, sub_id, EdgeType::Arrow);

        let mermaid = to_mermaid(&canvas);

        assert!(mermaid.contains("flowchart"));
        assert!(mermaid.contains("Main"));
        assert!(mermaid.contains("Sub"));
        assert!(mermaid.contains("-->"));
    }

    #[test]
    fn test_sanitize() {
        let result = sanitize_mermaid_id("abc-123:def");
        assert_eq!(result, "abc_123_def");
    }
}
