.PHONY: help install test clean run-api run-ocr run-rag deploy termux

# Colors
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

help:
	@echo "BAC Study - Available Commands"
	@echo "=============================="
	@echo "make install     - Install all dependencies"
	@echo "make termux      - Setup for Termux"
	@echo "make run-api     - Run Go API server"
	@echo "make run-ocr     - Run OCR server"
	@echo "make run-rag     - Run RAG pipeline"
	@echo "make deploy      - Deploy to Cloudflare"
	@echo "make clean       - Clean build artifacts"

install:
	@echo "$(GREEN)Installing dependencies...$(NC)"
	pip install -r scripts/requirements.txt --break-system-packages || pip install -r scripts/requirements.txt

termux:
	@echo "$(GREEN)Running Termux setup...$(NC)"
	bash scripts/setup-termux.sh

run-api:
	@echo "$(GREEN)Starting API server...$(NC)"
	cd src/api && go run main.go

run-ocr:
	@echo "$(GREEN)Starting OCR server...$(NC)"
	python scripts/ocr/local-ocr-server.py

run-rag:
	@echo "$(GREEN)Starting RAG pipeline...$(NC)"
	python scripts/rag_pipeline.py

test:
	@echo "$(GREEN)Running tests...$(NC)"
	cd src/api && go test ./...
	python -m pytest tests/ || true

deploy:
	@echo "$(GREEN)Deploying to Cloudflare...$(NC)"
	cd src/cloudflare && wrangler deploy

clean:
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -rf cache/ chroma_db/ __pycache__/
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete 2>/dev/null || true

# Build all services
build: build-api

build-api:
	@echo "$(GREEN)Building Go API...$(NC)"
	cd src/api && go build -o bin/api main.go
