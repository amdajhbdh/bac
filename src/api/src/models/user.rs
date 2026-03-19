use serde::{Deserialize, Serialize};

/// User model representing a platform user
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct User {
    pub id: String,
    pub username: String,
    pub email: String,
    pub created_at: String,
}
