//! OCR Extraction CLI
//!
//! Extract text from images using OCR

use bac_ocr::{perform_ocr_on_directory, perform_ocr_on_image};
use std::env;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        println!("Usage: ocr-extract <input> [output-file] [recursive]");
        println!("  input: Image file or directory");
        println!("  output-file: Optional output file (default: stdout)");
        println!("  recursive: Scan directories recursively (true/false, default: false)");
        println!("");
        println!("Examples:");
        println!("  ocr-extract image.png");
        println!("  ocr-extract image.png output.txt");
        println!("  ocr-extract /path/to/images /path/to/output.txt true");
        return Ok(());
    }

    let input = &args[1];
    let output_file = args.get(2).map(|s| s.as_str());
    let recursive = args.get(3).map(|s| s == "true").unwrap_or(false);

    // Check if input is a file or directory
    let input_path = std::path::Path::new(input);

    if input_path.is_file() {
        // OCR single image
        println!("Performing OCR on: {}", input);
        let text = perform_ocr_on_image(input)?;

        if let Some(output) = output_file {
            std::fs::write(output, &text)?;
            println!("Text saved to: {}", output);
        } else {
            println!("\n=== OCR Result ===\n");
            println!("{}", text);
            println!("\n=== End of Result ===");
        }
    } else if input_path.is_dir() {
        // OCR directory of images
        println!("Performing OCR on directory: {}", input);
        println!("Recursive: {}", recursive);

        let results = perform_ocr_on_directory(input, recursive)?;

        println!("\nFound {} images", results.len());

        if let Some(output) = output_file {
            // Combine all results into one file
            let mut combined_text = String::new();
            for (path, text) in &results {
                combined_text.push_str(&format!("=== {} ===\n", path.display()));
                combined_text.push_str(text);
                combined_text.push_str("\n\n");
            }
            std::fs::write(output, &combined_text)?;
            println!("Combined text saved to: {}", output);
        } else {
            // Print first few results
            for (_i, (path, text)) in results.iter().enumerate().take(5) {
                println!("\n=== {} ===", path.display());
                println!("{}", text);
            }
            if results.len() > 5 {
                println!("\n... and {} more images", results.len() - 5);
            }
        }
    } else {
        eprintln!("Error: Input path does not exist: {}", input);
    }

    Ok(())
}
