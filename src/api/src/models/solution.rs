use serde::{Deserialize, Serialize};

/// Solution model representing a problem solution
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Solution {
    pub id: String,
    pub problem_id: String,
    pub content: String,
    pub author: String,
    pub created_at: String,
}
