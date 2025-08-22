.PHONY: help build run test clean docker-build docker-up docker-down dev lint deps migrate seed docs

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
dev: ## Run the application in development mode
	@echo "Starting development server..."
	@go run ./cmd/server

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/music-app-backend ./cmd/server

run: build ## Build and run the application
	@echo "Running application..."
	@./bin/music-app-backend

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Database commands
migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run ./cmd/server --migrate-only

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	@go run ./scripts/seed

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t music-app-backend -f deploy/Dockerfile .

docker-up: ## Start development environment with Docker Compose
	@echo "Starting development environment..."
	@docker-compose -f deploy/docker-compose.yml up --build

docker-down: ## Stop development environment
	@echo "Stopping development environment..."
	@docker-compose -f deploy/docker-compose.yml down

docker-logs: ## Show Docker logs
	@docker-compose -f deploy/docker-compose.yml logs -f

# Documentation commands
docs: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g ./cmd/server/main.go -o ./docs

docs-serve: ## Serve Swagger documentation
	@echo "Starting documentation server..."
	@swagger serve -F swagger ./api/openapi.yaml

# Environment setup
setup: ## Set up development environment
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Git hooks
install-hooks: ## Install Git hooks
	@echo "Installing Git hooks..."
	@cp scripts/hooks/* .git/hooks/
	@chmod +x .git/hooks/*

# Release commands
tag: ## Create a new git tag
	@echo "Current tags:"
	@git tag -l
	@read -p "Enter new tag version (e.g., v1.0.0): " tag && \
		git tag -a $$tag -m "Release $$tag" && \
		git push origin $$tag

# Health check
health: ## Check application health
	@echo "Checking application health..."
	@curl -f http://localhost:8080/healthz || echo "Application is not running"

# Database operations
db-reset: ## Reset database (drops and recreates)
	@echo "Resetting database..."
	@docker-compose -f deploy/docker-compose.yml exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS music_app;"
	@docker-compose -f deploy/docker-compose.yml exec postgres psql -U postgres -c "CREATE DATABASE music_app;"

db-shell: ## Connect to database shell
	@echo "Connecting to database..."
	@docker-compose -f deploy/docker-compose.yml exec postgres psql -U postgres music_app
