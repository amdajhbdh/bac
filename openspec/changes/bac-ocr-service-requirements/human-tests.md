# OCR Service Human Tests

## Overview
Comprehensive human tests for OCR service with focus on complex scenarios including math equations, scientific notation, and multi-layer fallback.

## Test Categories

### Basic Functionality Tests

#### Test 1: Basic Text Extraction
**Scenario**: Extract simple French text from clean image
- **Input**: Clean PNG image with "Bonjour le monde! 2+2=4"
- **Expected**: 
  - Success: true
  - Text contains "Bonjour le monde"
  - Confidence > 0.8
  - Engine: tesseract

#### Test 2: English Text Extraction
**Scenario**: Extract simple English text from clean image
- **Input**: Clean PNG image with "Hello world! E=mc²"
- **Expected**: 
  - Success: true
  - Text contains "Hello world"
  - Confidence > 0.8
  - Engine: tesseract

#### Test 3: Multi-language Support
**Scenario**: Extract mixed language text
- **Input**: Clean PNG image with "Bonjour + Hello = BonjourHello"
- **Expected**: 
  - Success: true
  - Text contains both "Bonjour" and "Hello"
  - Confidence > 0.7
  - Engine: tesseract

### Math Equations Tests

#### Test 4: Simple Math Equation
**Scenario**: OCR simple mathematical equation
- **Input**: Image with "2 + 2 = 4" in clear font
- **Expected**: 
  - Success: true
  - Text contains "2 + 2 = 4"
  - Confidence > 0.7
  - Engine: tesseract

#### Test 5: Complex Math Equation
**Scenario**: OCR complex mathematical expression
- **Input**: Image with "√(b²-4ac) / 2a"
- **Expected**: 
  - Success: true
  - Text contains square root symbol and expression
  - Confidence > 0.6
  - Engine: tesseract (or surya if tesseract fails)

#### Test 6: Matrix Equation
**Scenario**: OCR matrix or system of equations
- **Input**: Image with:
  ```
  | 2  3 |   | x |   | 5 |
  | 1 -1 | x | y | = | 2 |
  ```
- **Expected**: 
  - Success: true
  - Text contains matrix structure
  - Confidence > 0.5
  - May require surya for complex layout

### Scientific Notation Tests

#### Test 7: Scientific Notation
**Scenario**: OCR scientific notation numbers
- **Input**: Image with "6.022×10²³ mol⁻¹" and "3.0×10⁸ m/s"
- **Expected**: 
  - Success: true
  - Text contains "6.022×10²³" and "3.0×10⁸"
  - Confidence > 0.7
  - Engine: tesseract

#### Test 8: Chemical Formulas
**Scenario**: OCR chemical formulas with subscripts
- **Input**: Image with "H₂O", "CO₂", "C₁₆H₁₀O₂"
- **Expected**: 
  - Success: true
  - Text contains correct chemical formulas
  - Confidence > 0.7
  - May need surya for subscript handling

### Complex Document Tests

#### Test 9: Multi-page PDF
**Scenario**: OCR multi-page PDF document
- **Input**: PDF with 5 pages of mixed content
- **Expected**: 
  - Success: true
  - All 5 pages processed
  - Combined text contains content from all pages
  - Average confidence reported

#### Test 10: Table Recognition
**Scenario**: OCR table with numerical data
- **Input**: Image with table:
  ```
  | Item     | Price | Quantity |
  |----------|-------|----------|
  | Apple    | $1.50 | 10       |
  | Banana   | $0.75 | 20       |
  ```
- **Expected**: 
  - Success: true
  - Text contains table structure
  - Confidence > 0.6
  - May require surya for table layout

#### Test 11: Column Layout
**Scenario**: OCR document with multiple columns
- **Input**: Image with two-column newspaper layout
- **Expected**: 
  - Success: true
  - Text preserves column structure
  - Confidence > 0.6
  - May need surya for column detection

### Edge Case Tests

#### Test 12: Handwritten Text
**Scenario**: OCR handwritten notes
- **Input**: Image with handwritten "Meet at 3 PM tomorrow"
- **Expected**: 
  - Success: true
  - Text contains handwritten content
  - Confidence > 0.4 (lower due to handwriting)
  - May require surya for better results

#### Test 13: Noisy Scanned Document
**Scenario**: OCR noisy scanned document
- **Input**: Low-quality scan with background noise
- **Expected**: 
  - Success: true
  - Text extracted despite noise
  - Confidence > 0.5
  - Multi-layer fallback should help

#### Test 14: Rotated Text
**Scenario**: OCR text at various angles
- **Input**: Image with text rotated 45 degrees
- **Expected**: 
  - Success: true
  - Text correctly oriented in output
  - Confidence > 0.6
  - May need preprocessing

#### Test 15: Small Text
**Scenario**: OCR very small text
- **Input**: Image with tiny text (6pt font)
- **Expected**: 
  - Success: true
  - Text extracted despite small size
  - Confidence > 0.5
  - May need DPI adjustment

### Fallback Mechanism Tests

#### Test 16: Tesseract Failure Fallback
**Scenario**: Primary OCR engine fails, fallback succeeds
- **Input**: Corrupted image that breaks tesseract
- **Expected**: 
  - Success: true
  - Engine: surya (if enabled)
  - Confidence reported
  - Attempts array shows tesseract failure

#### Test 17: Multi-engine Success
**Scenario**: Multiple engines succeed, best result chosen
- **Input**: Complex image suitable for both engines
- **Expected**: 
  - Success: true
  - Highest confidence engine selected
  - All attempts recorded
  - Confidence meets threshold

#### Test 18: All Engines Fail
**Scenario**: All OCR engines fail gracefully
- **Input**: Completely unreadable image
- **Expected**: 
  - Success: false
  - Error: All engines failed
  - Attempts array with all errors
  - Proper error message

### Performance Tests

#### Test 19: Batch Processing
**Scenario**: Process batch of 10 images
- **Input**: 10 different images
- **Expected**: 
  - Success: true
  - All 10 results returned
  - Total time < 5 seconds
  - Individual confidences reported

#### Test 20: Large File Handling
**Scenario**: Process large image (49MB)
- **Input**: 49MB image file
- **Expected**: 
  - Success: true
  - Processed without memory issues
  - Confidence reported
  - Processing time reasonable

#### Test 21: Rate Limiting
**Scenario**: Test rate limiting behavior
- **Input**: 101 requests in 1 minute
- **Expected**: 
  - First 100 succeed
  - 101st fails with rate limit error
  - Proper HTTP status code (429)

### API Endpoint Tests

#### Test 22: REST API - Single Image
**Scenario**: Test /ocr endpoint with valid image
- **Input**: POST /ocr with image file
- **Expected**: 
  - HTTP 200 OK
  - JSON response with success: true
  - Data contains text and confidence

#### Test 23: REST API - Invalid Input
**Scenario**: Test /ocr with invalid data
- **Input**: POST /ocr with corrupted file
- **Expected**: 
  - HTTP 400 Bad Request
  - JSON response with success: false
  - Error message explaining issue

#### Test 24: REST API - Batch Processing
**Scenario**: Test /ocr/batch endpoint
- **Input**: POST /ocr/batch with JSON array of images
- **Expected**: 
  - HTTP 200 OK
  - JSON response with batch results
  - Success/failure counts

#### Test 25: REST API - URL Processing
**Scenario**: Test /ocr/url endpoint
- **Input**: POST /ocr/url with image URL
- **Expected**: 
  - HTTP 200 OK
  - JSON response with OCR result
  - Handles network errors gracefully

## Test Execution Guidelines

### Prerequisites
1. OCR service must be running on port 50051
2. All dependencies installed (tesseract, pdftoppm, python/surya if enabled)
3. Test images available in test_data/ directory

### Test Data Requirements
- **Basic Tests**: Clean PNG images with simple text
- **Math Tests**: Images with mathematical equations
- **Scientific Tests**: Images with scientific notation and formulas
- **Complex Tests**: Multi-page PDFs, tables, column layouts
- **Edge Cases**: Noisy scans, handwritten text, rotated content

### Success Criteria
- **High Priority**: All basic functionality tests pass
- **Medium Priority**: Math and scientific notation tests pass
- **Low Priority**: Edge case tests pass with reasonable confidence

### Failure Analysis
For each failed test:
1. Check engine used (tesseract vs surya)
2. Verify confidence scores
3. Review error messages
4. Check system resources (CPU, memory)
5. Verify dependencies availability

## Automation Considerations

### Test Scripts
- **setup_tests.sh**: Prepare test environment
- **run_tests.sh**: Execute all human tests
- **generate_test_data.py**: Create synthetic test images
- **analyze_results.py**: Evaluate test outcomes

### Continuous Integration
- **GitHub Actions**: Run tests on each commit
- **Test Matrix**: Different OS, Python versions
- **Coverage Reporting**: Track test coverage over time

### Performance Monitoring
- **Benchmark Suite**: Track performance regressions
- **Memory Profiling**: Detect memory leaks
- **Load Testing**: Validate concurrent request handling

## Reporting

### Test Reports
- **HTML Report**: Detailed test results with screenshots
- **JSON Summary**: Machine-readable results
- **Coverage Report**: Test coverage statistics
- **Performance Report**: Timing and resource usage

### Failure Reports
- **Error Logs**: Detailed error information
- **Screenshots**: Visual context for failures
- **System Info**: Environment details
- **Dependencies**: Version information

## Maintenance

### Test Data Updates
- **Monthly**: Add new test cases
- **Quarterly**: Review test effectiveness
- **Annually**: Major test suite revision

### Test Suite Evolution
- **New Features**: Add tests for new functionality
- **Bug Fixes**: Add regression tests
- **Performance**: Update benchmarks
- **Security**: Add security-focused tests

## Quality Gates

### Release Criteria
- **All High Priority Tests Pass**
- **Coverage ≥ 90%**
- **Performance Benchmarks Met**
- **No Critical Security Issues**

### Daily Build
- **Quick Test Suite**: Basic functionality
- **Smoke Tests**: Core features
- **Health Checks**: Service availability

### Pre-release
- **Full Test Suite**: All test categories
- **Integration Tests**: End-to-end scenarios
- **Load Tests**: Performance under stress
- **Security Scan**: Vulnerability assessment