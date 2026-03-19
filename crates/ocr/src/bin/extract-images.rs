//! Image Extraction CLI
//!
//! Extract images from PDF files or process standalone images

use bac_ocr::{extract_images_from_pdf, process_standalone_images};
use std::env;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();

    if args.len() < 3 {
        println!("Usage: extract-images <input> <output> [mode]");
        println!("  mode: pdf (default) | standalone");
        println!("Examples:");
        println!("  extract-images document.pdf images/ pdf");
        println!("  extract-images /path/to/images /path/to/output standalone");
        return Ok(());
    }

    let input = &args[1];
    let output = &args[2];
    let mode = args.get(3).map(|s| s.as_str()).unwrap_or("pdf");

    match mode {
        "pdf" => {
            println!("Extracting images from PDF: {}", input);
            let images = extract_images_from_pdf(input, output)?;
            println!("Extracted {} images to {}", images.len(), output);
            for img in images {
                println!("  - {}", img.display());
            }
        }
        "standalone" => {
            println!("Processing standalone images from: {}", input);
            let images = process_standalone_images(input, output)?;
            println!("Processed {} images to {}", images.len(), output);
            for img in images {
                println!("  - {}", img.display());
            }
        }
        _ => {
            eprintln!("Error: Unknown mode '{}'", mode);
            eprintln!("Valid modes: pdf, standalone");
        }
    }

    Ok(())
}
