#!/bin/bash
# Download additional OCR test images

TEST_DATA_DIR="openspec/changes/bac-ocr-service-requirements/test_data"
mkdir -p "$TEST_DATA_DIR"

echo "Downloading sample OCR test images..."

# Sample math equation images from various sources
# These are public domain / CC0 images for testing

# Image 1: Simple math (using wget/curl)
curl -s -o "$TEST_DATA_DIR/downloaded_math.png" \
	"https://upload.wikimedia.org/wikipedia/commons/thumb/e/e4/Mathematical_equations_as_png.png/640px-Mathematical_equations_as_png.png" \
	2>/dev/null && echo "Downloaded math equations image" || echo "Math equations image not available"

# Image 2: Scientific notation example
curl -s -o "$TEST_DATA_DIR/downloaded_scientific.png" \
	"https://upload.wikimedia.org/wikipedia/commons/thumb/c/c3/Einstein_equation.png/320px-Einstein_equation.png" \
	2>/dev/null && echo "Downloaded scientific notation image" || echo "Scientific notation image not available"

# Check what we have
echo ""
echo "Test data files:"
ls -la "$TEST_DATA_DIR" | head -20

echo ""
echo "Note: For more comprehensive testing, you can add your own images to the test_data directory."
