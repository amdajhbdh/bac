---
source: Context7 API
library: tesseract-rs
package: tesseract-rs
topic: Tesseract OCR usage in Rust
fetched: 2026-03-18T10:30:00Z
official_docs: https://docs.rs/tesseract-rs/latest/tesseract_rs/
---

# Tesseract OCR Usage in Rust

## Overview
`tesseract-rs` is a Rust binding for Tesseract OCR with built-in compilation of Tesseract and Leptonica libraries. It provides a safe and idiomatic Rust interface to Tesseract's functionality.

## Installation
Add this to your `Cargo.toml`:

```toml
[dependencies]
tesseract-rs = { version = "0.1.20", features = ["build-tesseract"] }
```

For development and testing:
```toml
[dev-dependencies]
image = "0.25.5"
imageproc = "0.25.0"
```

## System Requirements
- C++ compiler (gcc, clang)
- CMake
- Internet connection (for downloading Tesseract training data)
- Rust 1.83.0 or later

## Basic Usage Example

```rust
use std::path::PathBuf;
use std::error::Error;
use tesseract_rs::TesseractAPI;

fn get_default_tessdata_dir() -> PathBuf {
    if cfg!(target_os = "macos") {
        let home_dir = std::env::var("HOME").expect("HOME environment variable not set");
        PathBuf::from(home_dir)
            .join("Library")
            .join("Application Support")
            .join("tesseract-rs")
            .join("tessdata")
    } else if cfg!(target_os = "linux") {
        let home_dir = std::env::var("HOME").expect("HOME environment variable not set");
        PathBuf::from(home_dir)
            .join(".tesseract-rs")
            .join("tessdata")
    } else if cfg!(target_os = "windows") {
        PathBuf::from(std::env::var("APPDATA").expect("APPDATA environment variable not set"))
            .join("tesseract-rs")
            .join("tessdata")
    } else {
        panic!("Unsupported operating system");
    }
}

fn main() -> Result<(), Box<dyn Error>> {
    let api = TesseractAPI::new()?;
    let tessdata_dir = get_default_tessdata_dir();
    api.init(tessdata_dir.to_str().unwrap(), "eng")?;

    let width = 24;
    let height = 24;
    let bytes_per_pixel = 1;
    let bytes_per_line = width * bytes_per_pixel;

    // Initialize image data with all white pixels
    let mut image_data = vec![255u8; width * height];

    // Set the image data
    api.set_image(
        &image_data,
        width.try_into().unwrap(),
        height.try_into().unwrap(),
        bytes_per_pixel.try_into().unwrap(),
        bytes_per_line.try_into().unwrap(),
    )?;

    // Set PSM mode
    api.set_variable("tessedit_pageseg_mode", "1")?;

    // Get the recognized text
    let text = api.get_utf8_text()?;
    println!("Recognized text: {}", text.trim());

    Ok(())
}
```

## Advanced Usage: Multi-language Support

For Arabic and English support, you need to load the appropriate traineddata files:

```rust
use tesseract_rs::TesseractAPI;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let api = TesseractAPI::new()?;
    
    // For Arabic + English support
    // Note: You need to have the appropriate traineddata files
    // Download from: https://github.com/tesseract-ocr/tessdata
    api.init("/path/to/tessdata", "ara+eng")?;
    
    // Set image (assuming you have image data)
    // api.set_image(...)?
    
    let text = api.get_utf8_text()?;
    println!("Recognized text: {}", text);
    
    Ok(())
}
```

## Configuration Variables

### Page Segmentation Modes (PSM)
- `0`: Orientation and script detection (OSD) only
- `1`: Automatic page segmentation with OSD
- `2`: Automatic page segmentation, but no OSD (OCR only)
- `3`: Fully automatic page segmentation, but no OSD
- `4`: Assume a single column of text of variable sizes
- `5`: Assume a single uniform block of vertically aligned text
- `6`: Assume a single uniform block of text
- `7`: Treat the image as a single text line
- `8`: Treat the image as a single word
- `9`: Treat the image as a single word in a circle
- `10`: Treat the image as a single character
- `11`: Sparse text. Find as much text as possible in no particular order
- `12`: Sparse text with OSD
- `13`: Raw line. Treat the image as a single text line, bypassing hacks

### Character Whitelisting
```rust
api.set_variable("tessedit_char_whitelist", "0123456789")?;
```

## Thread Safety
The API supports thread-safe operations:

```rust
use tesseract_rs::TesseractAPI;
use std::sync::Arc;
use std::thread;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let api = TesseractAPI::new()?;
    api.init("/path/to/tessdata", "eng")?;
    
    let api_arc = Arc::new(api);
    let mut handles = vec![];
    
    for _ in 0..3 {
        let api_clone = Arc::clone(&api_arc);
        let handle = thread::spawn(move || {
            // Perform OCR in parallel
            // Note: Each thread needs its own image data
        });
        handles.push(handle);
    }
    
    for handle in handles {
        handle.join().unwrap();
    }
    
    Ok(())
}
```

## Cache and Data Directories
The crate uses the following directory structure:
- macOS: `~/Library/Application Support/tesseract-rs`
- Linux: `~/.tesseract-rs`
- Windows: `%APPDATA%/tesseract-rs`

The cache includes compiled Tesseract and Leptonica libraries and downloaded training data.

## Alternative: kreuzberg-tesseract
For a more maintained version with C++17 support and cross-compilation fixes:
```toml
[dependencies]
kreuzberg-tesseract = "4.4.6"
```

## References
- [tesseract-rs Documentation](https://docs.rs/tesseract-rs/latest/tesseract_rs/)
- [Tesseract OCR GitHub](https://github.com/tesseract-ocr/tesseract)
- [Trained Data Files](https://github.com/tesseract-ocr/tessdata)
