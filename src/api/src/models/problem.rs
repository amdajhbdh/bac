use serde::{Deserialize, Serialize};

/// Problem model representing a study problem/question
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Problem {
    pub id: String,
    pub title: String,
    pub description: String,
    pub subject: Option<String>,
    pub created_at: String,
}
