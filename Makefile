.PHONY: help install test clean build run-all run-api run-rag run-ocr run-ocr-rust stop deploy termux

# Colors
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m

help:
	@echo "BAC Study - Available Commands"
	@echo "=============================="
	@echo "make install     - Install Python dependencies"
	@echo "make build       - Build all Rust services"
	@echo "make run-all     - Start all services"
	@echo "make stop       - Stop all services"
	@echo "make run-api    - Run Go API server"
	@echo "make run-rag    - Run Rust RAG pipeline"
	@echo "make run-ocr    - Run Python OCR server"
	@echo "make run-ocr-rust - Run Rust OCR service"
	@echo "make deploy     - Deploy to Cloudflare"
	@echo "make clean      - Clean build artifacts"
	@echo "make termux     - Setup for Termux"

install:
	@echo "$(GREEN)Installing Python dependencies...$(NC)"
	pip install -r scripts/requirements.txt --break-system-packages || pip install -r scripts/requirements.txt

termux:
	@echo "$(GREEN)Running Termux setup...$(NC)"
	bash scripts/setup-termux.sh

build:
	@echo "$(GREEN)Building Rust services...$(NC)"
	cargo build --release -p bac-rag || cargo build -p bac-rag
	@echo "$(GREEN)Building Go API...$(NC)"
	cd src/api && go build -o ../../bin/api . 2>/dev/null || true

run-all:
	@echo "$(GREEN)Starting all services...$(NC)"
	bash scripts/start-services.sh

stop:
	@echo "$(GREEN)Stopping all services...$(NC)"
	bash scripts/stop-services.sh

run-api:
	@echo "$(GREEN)Starting API server...$(NC)"
	cd src/api && go run main.go

run-rag:
	@echo "$(GREEN)Starting RAG pipeline (Rust)...$(NC)"
	cargo run -p bac-rag -- serve --port 5000

run-ocr:
	@echo "$(GREEN)Starting OCR server (Python)...$(NC)"
	python3 scripts/ocr/local-ocr-server.py

run-ocr-rust:
	@echo "$(GREEN)Starting OCR service (Rust)...$(NC)"
	cargo run -p bac-ocr

test:
	@echo "$(GREEN)Running tests...$(NC)"
	cd src/api && go test ./... || true
	cargo test || true

deploy:
	@echo "$(GREEN)Deploying to Cloudflare...$(NC)"
	cd src/cloudflare && wrangler deploy

clean:
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	cargo clean
	rm -rf bin/ data/
	rm -rf cache/ chroma_db/ __pycache__/
	find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	find . -type f -name "*.pyc" -delete 2>/dev/null || true

# Development
dev: build run-all
