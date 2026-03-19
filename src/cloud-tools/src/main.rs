//! Cloud Tools - Cloud Shell integration service (:3004)

use cloud_tools::run;

#[tokio::main]
async fn main() {
    run().await;
}
