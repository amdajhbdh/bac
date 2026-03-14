# OCR Service Test Suite

## Overview
Comprehensive test suite for the OCR service with HURL and human tests.

## Files

### Requirements & Documentation
- `hurl-requirements.md` - Full HURL requirements specification (100+ requirements)
- `human-tests.md` - Detailed human test scenarios (25+ tests)
- `implementation-summary.md` - Implementation details and changes

### Test Implementation
- `generate_test_data.py` - Python script to generate synthetic test images
- `tests/ocr-service.hurl` - HURL test file for automated API testing

### Test Execution

#### Prerequisites
```bash
# Install dependencies
pip install Pillow numpy

# Install hurl
cargo install hurl
# or
brew install hurl
```

#### Generate Test Data
```bash
python3 generate_test_data.py
```

#### Run HURL Tests
```bash
# Start the OCR service first
cd src/ocr-service
cargo run --bin ocr-service &

# Run HURL tests
hurl tests/ocr-service.hurl
```

#### Run Rust Tests
```bash
cd src/ocr-service
cargo test
```

## Test Categories

### 1. Basic Functionality (Tests 1-3)
- French text extraction
- English text extraction
- Multi-language support

### 2. Math Equations (Tests 4-6)
- Simple equations (2+2=4)
- Complex equations (quadratic formula)
- Fractions and integrals

### 3. Scientific Notation (Tests 7-8)
- Scientific constants (c, Na)
- Chemical formulas (H2O, CO2)

### 4. Complex Documents (Tests 9-11)
- Multi-page PDFs
- Tables
- Column layouts

### 5. Edge Cases (Tests 12-15)
- Noisy images
- Handwritten text
- Rotated text
- Small text

### 6. Fallback & Error Handling (Tests 16-18)
- Primary engine failure
- Multi-engine success
- All engines fail

### 7. Performance (Tests 19-21)
- Batch processing
- Large files
- Rate limiting

## Expected Results

### Success Criteria
- **Basic Tests**: 100% pass rate
- **Math Tests**: >80% pass rate
- **Scientific Tests**: >80% pass rate
- **Edge Cases**: >60% pass rate

### Performance Benchmarks
- Single image: <2 seconds
- Batch 10 images: <5 seconds
- PDF 10 pages: <10 seconds

## Manual Test Checklist

- [ ] Health check returns "healthy"
- [ ] French text extraction works
- [ ] English text extraction works
- [ ] Math equations recognized
- [ ] Scientific notation recognized
- [ ] PDF multi-page works
- [ ] Error handling works
- [ ] Rate limiting works

## CI/CD Integration

```yaml
# .github/workflows/ocr-tests.yml
name: OCR Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install Rust
        uses: actions-rs/toolchain@v1
      - name: Install hurl
        run: cargo install hurl
      - name: Generate test data
        run: python3 generate_test_data.py
      - name: Run tests
        run: hurl tests/ocr-service.hurl
## Troubleshooting

### Common```

 Issues
1. **Tesseract not found**: Install tesseract-ocr
2. **pdftoppm not found**: Install poppler-utils
3. **Python imports fail**: Install Pillow and numpy
4. **Port already in use**: Change port in config

### Debug Mode
```bash
RUST_LOG=debug cargo run --bin ocr-service
```
