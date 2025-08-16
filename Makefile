# a3s Makefile

# Variables
BINARY_NAME=a3s
BUILD_DIR=bin
MAIN_PATH=cmd/a3s/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	go run $(MAIN_PATH)

# Run with specific profile and region
.PHONY: run-dev
run-dev:
	go run $(MAIN_PATH) -profile dev -region us-west-2

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run tests
.PHONY: test
test:
	go test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install the binary to GOPATH/bin
.PHONY: install
install:
	go install $(MAIN_PATH)

# Development workflow - format, test, and build
.PHONY: dev
dev: fmt test build

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        Build the application"
	@echo "  run          Run the application"
	@echo "  run-dev      Run with dev profile and us-west-2 region"
	@echo "  deps         Download and tidy dependencies"
	@echo "  fmt          Format code"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  lint         Run linter (requires golangci-lint)"
	@echo "  clean        Clean build artifacts"
	@echo "  install      Install binary to GOPATH/bin"
	@echo "  dev          Format, test, and build"
	@echo "  help         Show this help message"