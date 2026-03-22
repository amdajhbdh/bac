//! Graph Tools - Data Models
//!
//! Core types for knowledge graph representation.

use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// Unique identifier for nodes and edges
pub type GraphId = String;

/// Canvas file format for Obsidian
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CanvasFile {
    pub nodes: Vec<CanvasNode>,
    pub edges: Vec<CanvasEdge>,
}

impl CanvasFile {
    /// Create a new empty canvas
    pub fn new() -> Self {
        Self {
            nodes: Vec::new(),
            edges: Vec::new(),
        }
    }

    /// Add a node to the canvas
    pub fn add_node(&mut self, node: CanvasNode) -> GraphId {
        let id = node.id.clone();
        self.nodes.push(node);
        id
    }

    /// Add an edge to the canvas
    pub fn add_edge(&mut self, edge: CanvasEdge) {
        self.edges.push(edge);
    }
}

impl Default for CanvasFile {
    fn default() -> Self {
        Self::new()
    }
}

/// Node in a canvas file
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CanvasNode {
    pub id: GraphId,
    pub x: f64,
    pub y: f64,
    pub width: f64,
    pub height: f64,
    #[serde(rename = "type")]
    pub node_type: NodeType,
    pub text: String,
    pub color: Option<String>,
}

impl CanvasNode {
    /// Create a new concept node
    pub fn concept(id: GraphId, x: f64, y: f64, text: &str) -> Self {
        Self {
            id,
            x,
            y,
            width: 400.0,
            height: 150.0,
            node_type: NodeType::Text,
            text: text.to_string(),
            color: None,
        }
    }

    /// Create a new group node
    pub fn group(id: GraphId, x: f64, y: f64, width: f64, height: f64, text: &str) -> Self {
        Self {
            id,
            x,
            y,
            width,
            height,
            node_type: NodeType::Group,
            text: text.to_string(),
            color: None,
        }
    }
}

/// Node type variants
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum NodeType {
    Text,
    Link,
    Group,
    Image,
}

/// Edge in a canvas file
#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CanvasEdge {
    pub id: GraphId,
    pub from_node: GraphId,
    pub from_side: Option<String>,
    pub to_node: GraphId,
    pub to_side: Option<String>,
    #[serde(rename = "type")]
    pub edge_type: EdgeType,
    pub label: Option<String>,
}

impl CanvasEdge {
    /// Create a new edge between nodes
    pub fn new(id: GraphId, from: GraphId, to: GraphId, edge_type: EdgeType) -> Self {
        Self {
            id,
            from_node: from,
            from_side: None,
            to_node: to,
            to_side: None,
            edge_type,
            label: None,
        }
    }

    /// Create an edge with a label
    pub fn with_label(
        id: GraphId,
        from: GraphId,
        to: GraphId,
        edge_type: EdgeType,
        label: &str,
    ) -> Self {
        Self {
            id,
            from_node: from,
            from_side: None,
            to_node: to,
            to_side: None,
            edge_type,
            label: Some(label.to_string()),
        }
    }
}

/// Edge type variants
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum EdgeType {
    #[serde(rename = "custom")]
    Custom,
    #[serde(rename = "arrow")]
    Arrow,
    #[serde(rename = "dashed")]
    Dashed,
}

// ============================================================================
// Entity Extraction Types
// ============================================================================

/// Extracted entities from notes
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Entities {
    pub concepts: Vec<Concept>,
    pub definitions: Vec<Definition>,
    pub relationships: Vec<EntityRelationship>,
}

impl Entities {
    /// Create empty entities container
    pub fn new() -> Self {
        Self {
            concepts: Vec::new(),
            definitions: Vec::new(),
            relationships: Vec::new(),
        }
    }

    /// Add a concept
    pub fn add_concept(&mut self, concept: Concept) {
        self.concepts.push(concept);
    }

    /// Add a definition
    pub fn add_definition(&mut self, definition: Definition) {
        self.definitions.push(definition);
    }

    /// Add a relationship
    pub fn add_relationship(&mut self, rel: EntityRelationship) {
        self.relationships.push(rel);
    }
}

impl Default for Entities {
    fn default() -> Self {
        Self::new()
    }
}

/// A concept extracted from notes
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Concept {
    pub id: GraphId,
    pub name: String,
    pub source_note: String,
    pub mentions: usize,
}

impl Concept {
    /// Create a new concept
    pub fn new(name: &str, source_note: &str) -> Self {
        Self {
            id: Uuid::new_v4().to_string(),
            name: name.to_string(),
            source_note: source_note.to_string(),
            mentions: 1,
        }
    }
}

/// A definition extracted from notes
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Definition {
    pub id: GraphId,
    pub term: String,
    pub definition: String,
    pub source_note: String,
}

impl Definition {
    /// Create a new definition
    pub fn new(term: &str, definition: &str, source_note: &str) -> Self {
        Self {
            id: Uuid::new_v4().to_string(),
            term: term.to_string(),
            definition: definition.to_string(),
            source_note: source_note.to_string(),
        }
    }
}

/// Relationship between entities
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct EntityRelationship {
    pub id: GraphId,
    pub source_id: GraphId,
    pub target_id: GraphId,
    pub relationship_type: RelationshipType,
    pub weight: f32,
}

impl EntityRelationship {
    /// Create a new relationship
    pub fn new(
        source_id: GraphId,
        target_id: GraphId,
        relationship_type: RelationshipType,
    ) -> Self {
        Self {
            id: Uuid::new_v4().to_string(),
            source_id,
            target_id,
            relationship_type,
            weight: 1.0,
        }
    }
}

/// Types of relationships between entities
#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum RelationshipType {
    Prerequisite,
    PartOf,
    Related,
    DependsOn,
    Example,
    Contrasts,
}

// ============================================================================
// Graph Export Types
// ============================================================================

/// Graph exported as JSON
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GraphJson {
    pub nodes: Vec<GraphNode>,
    pub edges: Vec<GraphEdge>,
    pub metadata: GraphMetadata,
}

impl GraphJson {
    /// Create graph JSON from canvas
    pub fn from_canvas(canvas: &CanvasFile, subject: &str) -> Self {
        let nodes: Vec<GraphNode> = canvas
            .nodes
            .iter()
            .map(|n| GraphNode {
                id: n.id.clone(),
                label: n.text.clone(),
                node_type: format!("{:?}", n.node_type).to_lowercase(),
                properties: serde_json::json!({
                    "x": n.x,
                    "y": n.y,
                    "width": n.width,
                    "height": n.height,
                }),
            })
            .collect();

        let edges: Vec<GraphEdge> = canvas
            .edges
            .iter()
            .map(|e| GraphEdge {
                id: e.id.clone(),
                source: e.from_node.clone(),
                target: e.to_node.clone(),
                edge_type: format!("{:?}", e.edge_type).to_lowercase(),
                label: e.label.clone(),
            })
            .collect();

        Self {
            nodes: nodes.clone(),
            edges: edges.clone(),
            metadata: GraphMetadata {
                subject: subject.to_string(),
                node_count: nodes.len(),
                edge_count: edges.len(),
            },
        }
    }
}

/// Node in exported graph
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GraphNode {
    pub id: String,
    pub label: String,
    pub node_type: String,
    pub properties: serde_json::Value,
}

/// Edge in exported graph
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GraphEdge {
    pub id: String,
    pub source: String,
    pub target: String,
    pub edge_type: String,
    pub label: Option<String>,
}

/// Graph metadata
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GraphMetadata {
    pub subject: String,
    pub node_count: usize,
    pub edge_count: usize,
}

// ============================================================================
// Visualization Types
// ============================================================================

/// Visualization URL result
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VisualizationUrl {
    pub url: String,
    pub embed_url: String,
    pub graph_data_url: Option<String>,
}

/// Health check response
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HealthResponse {
    pub status: String,
    pub version: String,
    pub service: String,
}
