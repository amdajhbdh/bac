//! PDF Text Extraction CLI
//!
//! Extract text from PDF files using unpdf

use bac_ocr::{detect_pdf_format, extract_text_from_pdf, pdf_to_markdown, validate_pdf};
use std::env;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        println!("Usage: extract-pdf <pdf-file> [output-format]");
        println!("  output-format: text (default) | markdown");
        println!("Examples:");
        println!("  extract-pdf document.pdf");
        println!("  extract-pdf document.pdf markdown");
        return Ok(());
    }

    let pdf_path = &args[1];
    let output_format = args.get(2).map(|s| s.as_str()).unwrap_or("text");

    // Validate PDF
    if !validate_pdf(pdf_path) {
        eprintln!("Error: File is not a valid PDF: {}", pdf_path);
        return Ok(());
    }

    // Detect format
    match detect_pdf_format(pdf_path) {
        Ok(format) => {
            println!("Detected PDF format: {:?}", format);
        }
        Err(e) => {
            eprintln!("Warning: Could not detect PDF format: {}", e);
        }
    }

    // Extract text based on format
    let result = match output_format {
        "markdown" | "md" => {
            println!("Extracting as Markdown...");
            pdf_to_markdown(pdf_path)
        }
        "text" | "txt" => {
            println!("Extracting as plain text...");
            extract_text_from_pdf(pdf_path)
        }
        _ => {
            eprintln!("Error: Unknown output format '{}'", output_format);
            return Ok(());
        }
    };

    match result {
        Ok(content) => {
            println!("\n=== Extracted Content ===\n");
            println!("{}", content);
            println!("\n=== End of Content ===");
            println!("\nTotal characters: {}", content.len());
            println!("Total lines: {}", content.lines().count());
        }
        Err(e) => {
            eprintln!("Error extracting PDF: {}", e);
        }
    }

    Ok(())
}
