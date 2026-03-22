//! OCR operations on Cloud Shell
//!
//! Runs OCR on Cloud Shell using Tesseract or cloud-based OCR.

use crate::models::OcrResult;
use crate::ssh::ssh_exec;

/// Run OCR on Cloud Shell
///
/// This uploads the image to Cloud Shell and runs OCR using Tesseract.
pub async fn run_ocr(image_path: &str, language: Option<&str>) -> OcrResult {
    tracing::info!("Running OCR on Cloud Shell for: {}", image_path);

    let lang = language.unwrap_or("eng");
    let output_base = format!("/tmp/ocr_output_{}", uuid_simple());

    // Build Tesseract command
    let tesseract_cmd = format!(
        "tesseract {} {} -l {} --psm 1",
        image_path, output_base, lang
    );

    let result = ssh_exec(&tesseract_cmd, Some(120)).await;

    if !result.success {
        return OcrResult::error(
            result.error
                .or(result.stderr)
                .unwrap_or_else(|| "OCR failed".to_string()),
        );
    }

    // Read the OCR output
    let output_file = format!("{}.txt", output_base);
    let read_result = ssh_exec(&format!("cat {}", output_file), Some(10)).await;

    // Clean up temporary files
    let _ = ssh_exec(
        &format!("rm -f {} {}*.txt {}*.png 2>/dev/null", output_file, output_base, output_base),
        Some(5),
    )
    .await;

    match read_result.success {
        true => {
            let text = read_result.stdout.unwrap_or_default().trim().to_string();
            let confidence = extract_confidence(&read_result.stderr);
            OcrResult::success(text, confidence)
        }
        false => OcrResult::error(
            read_result
                .error
                .or(read_result.stderr)
                .unwrap_or_else(|| "Failed to read OCR output".to_string()),
        ),
    }
}

/// Run OCR with PDF input (for scanned documents)
pub async fn run_pdf_ocr(pdf_path: &str, language: Option<&str>) -> OcrResult {
    tracing::info!("Running PDF OCR on Cloud Shell for: {}", pdf_path);

    let lang = language.unwrap_or("eng");
    let temp_dir = format!("/tmp/pdf_ocr_{}", uuid_simple());
    let output_file = format!("{}/output", temp_dir);

    // Commands to extract PDF pages and run OCR
    let commands = vec![
        format!("mkdir -p {}", temp_dir),
        format!("pdftoppm -png {} {}/page", pdf_path, temp_dir),
        format!(
            "for page in {}/page*.png; do tesseract \"$page\" \"{}/$(basename $page)\" -l {} --psm 1; done",
            temp_dir, temp_dir, lang
        ),
        format!("cat {}/page*.txt > {}", temp_dir, output_file),
        format!("cat {}", output_file),
    ];

    let combined_cmd = commands.join(" && ");
    let result = ssh_exec(&combined_cmd, Some(300)).await;

    // Cleanup
    let _ = ssh_exec(&format!("rm -rf {}", temp_dir), Some(5)).await;

    if result.success {
        let text = result.stdout.unwrap_or_default().trim().to_string();
        OcrResult::success(text, 0.85)
    } else {
        OcrResult::error(
            result
                .error
                .or(result.stderr)
                .unwrap_or_else(|| "PDF OCR failed".to_string()),
        )
    }
}

/// Get OCR availability on Cloud Shell
pub async fn check_ocr_availability() -> bool {
    let result = ssh_exec("which tesseract && tesseract --version | head -1", Some(10)).await;
    result.success
}

/// Get available OCR languages
pub async fn get_ocr_languages() -> Vec<String> {
    let result = ssh_exec("tesseract --list-langs 2>/dev/null | tail -n +2", Some(10)).await;
    
    result
        .stdout
        .unwrap_or_default()
        .lines()
        .map(|s| s.trim().to_string())
        .filter(|s| !s.is_empty())
        .collect()
}

/// Extract confidence from Tesseract stderr output
fn extract_confidence(stderr: &Option<String>) -> f32 {
    if let Some(s) = stderr {
        // Look for "Confidence:" in Tesseract output
        for line in s.lines() {
            if line.contains("Confidence:") {
                if let Some(num_str) = line.split_whitespace().last() {
                    if let Ok(conf) = num_str.parse::<f32>() {
                        return conf / 100.0; // Convert to 0-1 range
                    }
                }
            }
        }
    }
    0.85 // Default confidence
}

/// Simple UUID generator for temp files
fn uuid_simple() -> String {
    use std::time::{SystemTime, UNIX_EPOCH};
    let nanos = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap_or_default()
        .as_nanos();
    format!("{:x}", nanos)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_uuid() {
        let id = uuid_simple();
        assert!(!id.is_empty());
        tracing::debug!("Generated UUID: {}", id);
    }

    #[test]
    fn test_confidence_parsing() {
        let stderr = "Tesseract Open Source OCR Engine v5.3.0\nFeatures used: 神经网络\nConfidence: 92.5";
        let conf = extract_confidence(&Some(stderr.to_string()));
        assert!((conf - 0.925).abs() < 0.01);
    }
}
