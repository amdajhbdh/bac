---
source: Context7 API
library: imageproc / OpenCV
package: imageproc
topic: Image preprocessing techniques for OCR
fetched: 2026-03-18T10:30:00Z
official_docs: https://docs.rs/imageproc/latest/imageproc/
---

# Image Preprocessing Techniques for OCR

## Overview
Image preprocessing is crucial for improving OCR accuracy. This document covers techniques for binarization, noise removal, and document preparation specifically for scanned exam papers in Arabic and French.

## Common Preprocessing Steps

### 1. Grayscale Conversion
Convert color images to grayscale to reduce complexity and improve processing speed.

```rust
use image::{DynamicImage, GrayImage};
use imageproc::contrast::threshold;

fn convert_to_grayscale(img: &DynamicImage) -> GrayImage {
    img.to_luma8()
}
```

### 2. Binarization (Thresholding)
Convert grayscale images to black and white (binary) images.

#### Global Thresholding
```rust
use imageproc::contrast::threshold;

fn global_threshold(img: &GrayImage, threshold: u8) -> GrayImage {
    threshold(img, threshold)
}
```

#### Adaptive Thresholding (Better for uneven lighting)
```rust
use imageproc::contrast::adaptive_threshold;

fn adaptive_threshold(img: &GrayImage, block_size: u32) -> GrayImage {
    adaptive_threshold(img, block_size)
}
```

### 3. Noise Removal
Remove noise from scanned documents using various filtering techniques.

#### Median Filtering (Good for salt-and-pepper noise)
```rust
use imageproc::filter::median_filter;

fn remove_noise_median(img: &GrayImage) -> GrayImage {
    // 3x3 kernel
    median_filter(img, 1)
}
```

#### Gaussian Blur (Good for general noise reduction)
```rust
use imageproc::filter::gaussian_blur_f32;

fn blur_image(img: &GrayImage, sigma: f32) -> GrayImage {
    gaussian_blur_f32(img, sigma)
}
```

### 4. Deskewing (Straightening Skewed Documents)
Straighten tilted documents for better OCR accuracy.

```rust
use imageproc::geometric_transformations::{rotate_about_center, Interpolation};

fn deskew_image(img: &GrayImage, angle: f32) -> GrayImage {
    rotate_about_center(
        img,
        angle,
        Interpolation::Bilinear,
        [255; 3] // Background color
    )
}
```

### 5. Contrast Enhancement
Improve text visibility by enhancing contrast.

```rust
use imageproc::contrast::stretch_contrast;

fn enhance_contrast(img: &GrayImage) -> GrayImage {
    stretch_contrast(img)
}
```

## Complete Preprocessing Pipeline for Scanned Documents

```rust
use image::{DynamicImage, GrayImage};
use imageproc::{
    contrast::{adaptive_threshold, stretch_contrast},
    filter::median_filter,
    geometric_transformations::{rotate_about_center, Interpolation},
};

pub struct PreprocessingConfig {
    pub threshold_block_size: u32,
    pub median_kernel: u32,
    pub deskew_angle: f32,
    pub enhance_contrast: bool,
}

impl Default for PreprocessingConfig {
    fn default() -> Self {
        Self {
            threshold_block_size: 41,
            median_kernel: 1,
            deskew_angle: 0.0,
            enhance_contrast: true,
        }
    }
}

pub fn preprocess_document(
    img: &DynamicImage,
    config: &PreprocessingConfig,
) -> GrayImage {
    // Step 1: Convert to grayscale
    let mut processed = img.to_luma8();
    
    // Step 2: Remove noise
    if config.median_kernel > 0 {
        processed = median_filter(&processed, config.median_kernel);
    }
    
    // Step 3: Enhance contrast (optional)
    if config.enhance_contrast {
        processed = stretch_contrast(&processed);
    }
    
    // Step 4: Deskew if angle is provided
    if config.deskew_angle.abs() > 0.001 {
        processed = rotate_about_center(
            &processed,
            config.deskew_angle,
            Interpolation::Bilinear,
            [255; 3]
        );
    }
    
    // Step 5: Binarization (adaptive thresholding)
    processed = adaptive_threshold(&processed, config.threshold_block_size);
    
    processed
}
```

## Rust Crates for Image Processing

### imageproc
- **Purpose**: Image processing operations
- **Features**: Filtering, contrast, geometric transformations
- **Link**: https://docs.rs/imageproc

### image
- **Purpose**: Image loading and basic manipulation
- **Features**: Support for multiple formats (PNG, JPEG, TIFF)
- **Link**: https://docs.rs/image

### opencv (Rust bindings)
- **Purpose**: Comprehensive computer vision library
- **Features**: Advanced preprocessing, deskewing, noise removal
- **Link**: https://docs.rs/opencv

## Best Practices for Scanned Documents

### 1. Resolution
- Scan at minimum 300 DPI for optimal OCR accuracy
- Higher resolution (600 DPI) for small fonts or poor quality originals

### 2. Lighting
- Ensure even lighting across the document
- Avoid shadows and glare
- Use consistent lighting conditions

### 3. Document Preparation
- Remove staples, folds, and creases before scanning
- Flatten documents to avoid distortion
- Clean scanner glass to avoid artifacts

### 4. File Format
- Use uncompressed formats (TIFF, PNG) for preprocessing
- Avoid JPEG compression artifacts
- Save final processed images as PNG for OCR

## Arabic-Specific Considerations

### Right-to-Left (RTL) Text
- Arabic text flows right-to-left
- Ensure preprocessing doesn't reverse text direction
- Use appropriate page segmentation modes in Tesseract

### Connected Characters
- Arabic script has connected characters
- Avoid aggressive noise removal that might break connections
- Use moderate binarization thresholds

### Diacritical Marks
- Arabic has many diacritical marks (harakat)
- Ensure preprocessing preserves small details
- Avoid excessive blurring that might merge marks with letters

## Example: Complete Pipeline for Arabic/English Documents

```rust
use image::{DynamicImage, GrayImage};
use imageproc::{
    contrast::adaptive_threshold,
    filter::median_filter,
};

pub fn preprocess_arabic_english_document(img: &DynamicImage) -> GrayImage {
    // Convert to grayscale
    let mut processed = img.to_luma8();
    
    // Remove noise (moderate kernel to preserve Arabic connections)
    processed = median_filter(&processed, 1);
    
    // Adaptive thresholding with larger block size for uneven lighting
    processed = adaptive_threshold(&processed, 51);
    
    processed
}
```

## References
- [imageproc Documentation](https://docs.rs/imageproc/latest/imageproc/)
- [OpenCV Image Processing](https://docs.rs/opencv/latest/opencv/imgproc/index.html)
- [Tesseract OCR Image Quality Guide](https://guides.nyu.edu/tesseract/image-quality)
