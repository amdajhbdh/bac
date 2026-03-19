//! Batch PDF Extraction CLI
//!
//! Extract text from all PDFs in a directory

use bac_ocr::{extract_text_from_pdf, pdf_to_markdown, validate_pdf};
use std::env;
use std::fs;
use std::path::Path;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        println!("Usage: batch-extract <directory> [output-format]");
        println!("  output-format: text (default) | markdown");
        println!("Examples:");
        println!("  batch-extract /path/to/pdfs");
        println!("  batch-extract /path/to/pdfs markdown");
        return Ok(());
    }

    let dir_path = &args[1];
    let output_format = args.get(2).map(|s| s.as_str()).unwrap_or("text");

    // Check if directory exists
    if !Path::new(dir_path).is_dir() {
        eprintln!("Error: Directory does not exist: {}", dir_path);
        return Ok(());
    }

    println!("Scanning directory: {}", dir_path);
    println!("Output format: {}", output_format);

    // Find all PDF files
    let mut pdf_files = Vec::new();
    for entry in fs::read_dir(dir_path)? {
        let entry = entry?;
        let path = entry.path();

        if path.is_file() {
            let ext = path
                .extension()
                .and_then(|s| s.to_str())
                .unwrap_or("")
                .to_lowercase();
            if ext == "pdf" {
                pdf_files.push(path);
            }
        }
    }

    println!("Found {} PDF files", pdf_files.len());

    // Process each PDF
    for (i, pdf_path) in pdf_files.iter().enumerate() {
        let filename = pdf_path.file_name().unwrap().to_string_lossy();
        println!("[{}/{}] Processing: {}", i + 1, pdf_files.len(), filename);

        // Validate PDF
        if !validate_pdf(pdf_path) {
            println!("  ⚠️  Skipping: Not a valid PDF");
            continue;
        }

        // Extract text
        let result = match output_format {
            "markdown" | "md" => pdf_to_markdown(pdf_path),
            "text" | "txt" => extract_text_from_pdf(pdf_path),
            _ => {
                println!("  ⚠️  Skipping: Unknown format '{}'", output_format);
                continue;
            }
        };

        match result {
            Ok(content) => {
                // Create output filename
                let mut output_path = pdf_path.clone();
                output_path.set_extension(match output_format {
                    "markdown" | "md" => "md",
                    _ => "txt",
                });

                // Write to file
                match fs::write(&output_path, &content) {
                    Ok(_) => {
                        println!(
                            "  ✅ Extracted: {} chars -> {}",
                            content.len(),
                            output_path.file_name().unwrap().to_string_lossy()
                        );
                    }
                    Err(e) => {
                        println!("  ❌ Error writing file: {}", e);
                    }
                }
            }
            Err(e) => {
                println!("  ❌ Error extracting: {}", e);
            }
        }
    }

    println!("Batch extraction complete!");

    Ok(())
}
