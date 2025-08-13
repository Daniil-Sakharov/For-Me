# Makefile for Clean Architecture URL Shortener

.PHONY: help build run test clean docker-up docker-down migrate

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the application binary
	@echo "🔨 Building application..."
	@go build -o bin/url-shortener cmd/main.go
	@echo "✅ Build complete!"

# Run the application
run: ## Run the application
	@echo "🚀 Starting application..."
	@go run cmd/main.go

# Run tests
test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "📊 Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

# Start dependencies with Docker
docker-up: ## Start PostgreSQL and Redis with Docker Compose
	@echo "🐳 Starting dependencies..."
	@docker-compose up -d
	@echo "✅ Dependencies started!"

# Stop dependencies
docker-down: ## Stop Docker dependencies
	@echo "🛑 Stopping dependencies..."
	@docker-compose down
	@echo "✅ Dependencies stopped!"

# Format code
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted!"

# Run linter
lint: ## Run golangci-lint
	@echo "🔍 Running linter..."
	@golangci-lint run
	@echo "✅ Linting complete!"

# Install dependencies
deps: ## Download dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

# Generate code (if needed)
generate: ## Run go generate
	@echo "🔧 Generating code..."
	@go generate ./...
	@echo "✅ Code generation complete!"

# Development setup
dev-setup: deps docker-up ## Setup development environment
	@echo "🛠️ Setting up development environment..."
	@echo "⏳ Waiting for database to be ready..."
	@sleep 5
	@echo "✅ Development environment ready!"

# Run all checks (format, lint, test)
check: fmt lint test ## Run all code quality checks

# Quick start for new developers
quickstart: dev-setup build ## Complete setup for new developers
	@echo ""
	@echo "🎉 Quickstart complete!"
	@echo ""
	@echo "To run the application:"
	@echo "  make run"
	@echo ""
	@echo "To run tests:"
	@echo "  make test"
	@echo ""
	@echo "To stop dependencies:"
	@echo "  make docker-down"