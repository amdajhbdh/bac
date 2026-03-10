# OCR Service HURL Tests
# Run with: hurl tests/ocr-service.hurl

# Health Check
GET http://127.0.0.1:3000/health
HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"

# Test 1: Basic French Text Extraction
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/basic_french.png
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.text" contains "Bonjour"
jsonpath "$.data.confidence" >= 0.7

# Test 2: Basic English Text Extraction
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/basic_english.png
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.text" contains "Hello"

# Test 3: Mixed Language Text
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/basic_mixed.png
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test 4: Simple Math Equation
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/math_simple.png
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test 5: Complex Math Equation
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/math_quadratic.png
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test 6: Scientific Notation
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/sci_constants.png
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test 7: Chemical Formula
POST http://127.0.0.1:3000/ocr
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/sci_formula_h2o.png
HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test 8: PDF Processing
POST http://127.0.0.1:3000/pdf
[Form]
file @openspec/changes/bac-ocr-service-requirements/test_data/multipage_test.pdf
HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.total_pages" >= 1

# Test 9: Image from URL
POST http://127.0.0.1:3000/ocr/url
Content-Type: application/json
{
  "url": "https://example.com/test_image.png"
}
HTTP 400
