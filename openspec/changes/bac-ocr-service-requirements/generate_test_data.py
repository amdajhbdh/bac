#!/usr/bin/env python3
"""
OCR Test Data Generator

Generates synthetic test images and PDFs for OCR testing including:
- Simple text images
- Mathematical equations
- Scientific notation
- Complex layouts (tables, columns)
- Edge cases (noisy, handwritten style)
"""

import os
import sys
from pathlib import Path

try:
    from PIL import Image, ImageDraw, ImageFont
    import numpy as np
except ImportError:
    print("Error: Pillow and numpy required. Install with:")
    print("  pip install Pillow numpy")
    sys.exit(1)

OUTPUT_DIR = Path(__file__).parent / "test_data"
OUTPUT_DIR.mkdir(exist_ok=True)

FONTS = {
    "regular": "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
    "bold": "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
    "mono": "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
}


def get_font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont:
    """Get font with fallback"""
    font_path = FONTS["bold"] if bold else FONTS["regular"]
    try:
        return ImageFont.truetype(font_path, size)
    except:
        return ImageFont.load_default()


def create_simple_text_image(text: str, output_path: Path, size: tuple = (800, 200)):
    """Create simple text image"""
    img = Image.new("L", size, color=255)
    draw = ImageDraw.Draw(img)
    font = get_font(32)

    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]

    x = (size[0] - text_width) // 2
    y = (size[1] - text_height) // 2

    draw.text((x, y), text, fill=0, font=font)
    img.save(output_path)
    print(f"Created: {output_path}")


def create_math_equation_image(equation: str, output_path: Path):
    """Create image with mathematical equation"""
    img = Image.new("L", (800, 200), color=255)
    draw = ImageDraw.Draw(img)
    font = get_font(36, bold=True)

    bbox = draw.textbbox((0, 0), equation, font=font)
    text_width = bbox[2] - bbox[0]

    x = (800 - text_width) // 2
    y = 80

    draw.text((x, y), equation, fill=0, font=font)
    img.save(output_path)
    print(f"Created: {output_path}")


def create_scientific_notation_image(text: str, output_path: Path):
    """Create image with scientific notation"""
    img = Image.new("L", (900, 200), color=255)
    draw = ImageDraw.Draw(img)
    font = get_font(32, bold=True)

    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]

    x = (900 - text_width) // 2
    y = 80

    draw.text((x, y), text, fill=0, font=font)
    img.save(output_path)
    print(f"Created: {output_path}")


def create_table_image(data: list, output_path: Path):
    """Create image with table"""
    rows = len(data)
    cols = len(data[0]) if data else 0

    cell_width = 200
    cell_height = 50
    padding = 10

    width = cols * cell_width + padding * 2
    height = rows * cell_height + padding * 2

    img = Image.new("L", (width, height), color=255)
    draw = ImageDraw.Draw(img)
    font = get_font(20)

    for row_idx, row in enumerate(data):
        for col_idx, cell in enumerate(row):
            x = padding + col_idx * cell_width
            y = padding + row_idx * cell_height

            draw.rectangle([x, y, x + cell_width, y + cell_height], outline=0)
            draw.text((x + 10, y + 15), cell, fill=0, font=font)

    img.save(output_path)
    print(f"Created: {output_path}")


def create_noisy_image(text: str, output_path: Path, noise_level: float = 0.1):
    """Create image with noise (simulating scanned document)"""
    img = Image.new("L", (800, 200), color=255)
    draw = ImageDraw.Draw(img)
    font = get_font(32)

    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]

    x = (800 - text_width) // 2
    y = (200 - text_height) // 2

    draw.text((x, y), text, fill=0, font=font)

    img_array = np.array(img)
    noise = np.random.normal(0, 255 * noise_level, img_array.shape)
    img_array = np.clip(img_array + noise, 0, 255).astype(np.uint8)

    img = Image.fromarray(img_array)
    img.save(output_path)
    print(f"Created: {output_path}")


def create_multipage_pdf(pages: list, output_path: Path):
    """Create multi-page PDF"""
    try:
        from PIL import Image as PILImage, ImageDraw as PILImageDraw

        images = []
        for text in pages:
            img = PILImage.new("L", (800, 1000), color=255)
            draw = PILImageDraw.Draw(img)
            font = get_font(24)

            y = 50
            for line in text.split("\n"):
                draw.text((50, y), line, fill=0, font=font)
                y += 40

            images.append(img)

        if images:
            images[0].save(output_path, save_all=True, append_images=images[1:])
            print(f"Created: {output_path}")
    except Exception as e:
        print(f"Error creating PDF: {e}")


def main():
    """Generate all test data"""
    print(f"Generating test data in: {OUTPUT_DIR}")

    basic_texts = [
        ("basic_french.txt", "Bonjour le monde!"),
        ("basic_english.txt", "Hello world!"),
        ("basic_mixed.txt", "Bonjour + Hello = BonjourHello"),
    ]

    math_equations = [
        ("math_simple.txt", "2 + 2 = 4"),
        ("math_quadratic.txt", "x = (-b ± √(b²-4ac)) / 2a"),
        ("math_fraction.txt", "a/b + c/d = (ad+bc)/bd"),
        ("math_integral.txt", "∫f(x)dx = F(x) + C"),
    ]

    scientific_notations = [
        ("sci_constants.txt", "c = 3.0×10⁸ m/s"),
        ("sci_avogadro.txt", "Nₐ = 6.022×10²³ mol⁻¹"),
        ("sci_formula_h2o.txt", "H₂O"),
        ("sci_formula_co2.txt", "CO₂"),
    ]

    for filename, text in basic_texts:
        create_simple_text_image(
            text, OUTPUT_DIR / f"{filename.replace('.txt', '.png')}"
        )

    for filename, equation in math_equations:
        create_math_equation_image(
            equation, OUTPUT_DIR / f"{filename.replace('.txt', '.png')}"
        )

    for filename, notation in scientific_notations:
        create_scientific_notation_image(
            notation, OUTPUT_DIR / f"{filename.replace('.txt', '.png')}"
        )

    table_data = [
        ["Item", "Price", "Qty"],
        ["Apple", "$1.50", "10"],
        ["Banana", "$0.75", "20"],
    ]
    create_table_image(table_data, OUTPUT_DIR / "table_simple.png")

    create_noisy_image(
        "Noisy text here", OUTPUT_DIR / "noisy_text.png", noise_level=0.15
    )
    create_noisy_image(
        "Handwritten style", OUTPUT_DIR / "handwritten_style.png", noise_level=0.25
    )

    pdf_pages = [
        "Page 1: Introduction\nThis is a test document.",
        "Page 2: Mathematics\n2 + 2 = 4\nE = mc²",
        "Page 3: Conclusion\nThank you for reading.",
    ]
    create_multipage_pdf(pdf_pages, OUTPUT_DIR / "multipage_test.pdf")

    print(f"\nTest data generation complete!")
    print(f"Output directory: {OUTPUT_DIR}")
    print(f"\nGenerated files:")
    for f in sorted(OUTPUT_DIR.iterdir()):
        print(f"  - {f.name}")


if __name__ == "__main__":
    main()
