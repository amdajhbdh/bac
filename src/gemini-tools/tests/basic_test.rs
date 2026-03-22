//! Basic Integration Tests for gemini-tools
//!
//! These tests verify the core functionality of gemini-tools.
//! Note: Most tests require a valid GEMINI_API_KEY.

use gemini_tools::{
    analyze::{self, AnalysisResult},
    correct::{self, CorrectedText},
    embed::{self, Embedding},
    extract::{self, ExtractedEntities},
    generate::{self, GeneratedNote},
    GeminiClient,
};

/// Create a test client (requires GEMINI_API_KEY)
fn create_test_client() -> Option<GeminiClient> {
    std::env::var("GEMINI_API_KEY")
        .ok()
        .map(GeminiClient::new)
}

/// Test JSON parsing for AnalysisResult
#[test]
fn test_analysis_result_serde() {
    let json = r#"{
        "summary": "Test summary",
        "difficulty": "beginner",
        "key_points": ["Point 1", "Point 2"],
        "prerequisites": [],
        "learning_objectives": ["Objective 1"]
    }"#;

    let result: AnalysisResult = serde_json::from_str(json).unwrap();
    assert_eq!(result.summary, "Test summary");
    assert_eq!(result.difficulty, "beginner");
    assert_eq!(result.key_points.len(), 2);
}

/// Test JSON parsing for ExtractedEntities
#[test]
fn test_extracted_entities_serde() {
    let json = r#"{
        "concepts": ["Concept A", "Concept B"],
        "formulas": ["F = ma"],
        "topics": ["Physics"],
        "definitions": [
            {"term": "Force", "definition": "Push or pull"}
        ],
        "examples": ["Example 1"]
    }"#;

    let result: ExtractedEntities = serde_json::from_str(json).unwrap();
    assert_eq!(result.concepts.len(), 2);
    assert_eq!(result.formulas[0], "F = ma");
    assert_eq!(result.definitions.len(), 1);
}

/// Test JSON parsing for GeneratedNote
#[test]
fn test_generated_note_serde() {
    let json = "{\"title\": \"Newton's Laws\", \"content\": \"# Introduction\\n\\nContent here\", \"format\": \"markdown\"}";
    let result: GeneratedNote = serde_json::from_str(json).unwrap();
    assert_eq!(result.title, "Newton's Laws");
    assert!(result.content.contains("Introduction"));
}

/// Test JSON parsing for Embedding
#[test]
fn test_embedding_serde() {
    let json = r#"{
        "values": [0.1, 0.2, 0.3],
        "model": "gemini-embedding",
        "task_type": "RETRIEVAL_DOCUMENT"
    }"#;

    let result: Embedding = serde_json::from_str(json).unwrap();
    assert_eq!(result.values.len(), 3);
    assert_eq!(result.values[0], 0.1);
    assert_eq!(result.model, "gemini-embedding");
}

/// Test JSON parsing for CorrectedText
#[test]
fn test_corrected_text_serde() {
    let json = r#"{
        "corrected": "Hello World",
        "corrections": [
            {"original": "He110", "corrected": "Hello", "reason": "OCR error"}
        ]
    }"#;

    let result: CorrectedText = serde_json::from_str(json).unwrap();
    assert_eq!(result.corrected, "Hello World");
    assert_eq!(result.corrections.len(), 1);
    assert_eq!(result.corrections[0].original, "He110");
}

/// Test analyze::extract_json_from_response
#[test]
fn test_analyze_extract_json() {
    let response = r#"Some text before
{
    "summary": "Test",
    "difficulty": "beginner",
    "key_points": [],
    "prerequisites": [],
    "learning_objectives": []
}
Some text after"#;

    let result = analyze::extract_json_from_response(response);
    assert!(result.is_ok());
    assert_eq!(result.unwrap().summary, "Test");
}

/// Test extract::extract_json_from_response
#[test]
fn test_extract_parsed_json() {
    let response = r#"{
        "concepts": ["Test"],
        "formulas": [],
        "topics": [],
        "definitions": [],
        "examples": []
    }"#;

    let result = extract::extract_json_from_response(response);
    assert!(result.is_ok());
}

/// Test correct::extract_json_from_response
#[test]
fn test_correct_extract_json() {
    let response = r#"Analysis complete:
{
    "corrected": "Fixed text",
    "corrections": []
}"#;

    let result = correct::extract_json_from_response(response);
    assert!(result.is_ok());
}

/// Test generate::extract_json_content
#[test]
fn test_generate_extract_content() {
    let response = r#"{
        "title": "Test",
        "content": "Actual content",
        "format": "markdown"
    }"#;

    let content = generate::extract_json_content(response);
    assert!(content.is_some());
    assert_eq!(content.unwrap(), "Actual content");
}

/// Test client creation from string
#[test]
fn test_client_creation() {
    let client = GeminiClient::new("test-api-key");
    // Just verify it can be created without panic
    let _ = client;
}

/// Test client creation with model
#[test]
fn test_client_with_model() {
    let client = GeminiClient::new("test-api-key").with_model("gemini-pro");
    let _ = client;
}

// Integration tests (skip if no API key)
#[test]
#[ignore]
fn test_analyze_integration() {
    let client = create_test_client().expect("GEMINI_API_KEY not set");
    let rt = tokio::runtime::Runtime::new().unwrap();
    
    rt.block_on(async {
        let result = analyze::analyze(
            &client,
            "The equation F = ma describes Newton's second law.",
            Some("Physics"),
        )
        .await;
        assert!(result.is_ok());
    });
}

#[test]
#[ignore]
fn test_extract_integration() {
    let client = create_test_client().expect("GEMINI_API_KEY not set");
    let rt = tokio::runtime::Runtime::new().unwrap();
    
    rt.block_on(async {
        let result = extract::extract(&client, "F = ma is Newton's second law.")
            .await;
        assert!(result.is_ok());
    });
}
