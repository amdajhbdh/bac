# Test 1: Gateway Health
GET http://127.0.0.1:8080/health
HTTP/2 200
[Asserts]
json("$.status" == "healthy")
json("$.ocr_service" == "connected"
# Test 2: OCR Service Health
GET http://127.0.0.1:3000/health
HTTP/2 200
[Asserts]
json("$.status" == "healthy")
# Test 3: Simple Image OCR
POST http://127.0.0.1:8080/ocr
Content-Type: multipart/form-data
[Form]
file: /tmp/test_simple.png
HTTP/2 200
[Asserts]
json("$.success" == true)
json("$.data.engine" == "tesseract")
# Test 4: Math OCR
POST http://127.0.0.1:8080/ocr
Content-Type: multipart/form-data
[Form]
file: /tmp/test_math.png
HTTP/2 200
[Asserts]
json("$.success" == true)
# Test 5: French OCR
POST http://127.0.0.1:8080/ocr
Content-Type: multipart/form-data
[Form]
file: /tmp/test_french.png
HTTP/2 200
[Asserts]
json("$.success" == true)
# Test 6: PDF OCR
POST http://127.0.0.1:8080/ocr
Content-Type: multipart/form-data
[Form]
file: /tmp/test_document.pdf
HTTP/2 200
[Asserts]
json("$.success" == true)
json("$.data.engine" == "pdf")
