.PHONY: build run test clean setup-db dev help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=pos-api
BINARY_PATH=./cmd/api

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)

# Run the application
run: build
	@echo "🚀 Starting POS System API..."
	./$(BINARY_NAME)

# Run in development mode (with go run)
dev:
	@echo "🔧 Starting development server..."
	$(GOCMD) run $(BINARY_PATH)/main.go

# Test the application
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Setup database
setup-db:
	@echo "📦 Setting up database..."
	./scripts/setup-db.sh

# Install dependencies
deps:
	@echo "📥 Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Show help
help:
	@echo "📖 Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Build and run the application"
	@echo "  dev       - Run in development mode"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  setup-db  - Setup PostgreSQL database"
	@echo "  deps      - Install dependencies"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code (requires golangci-lint)"
	@echo "  help      - Show this help"

# Default target
all: clean deps fmt build
