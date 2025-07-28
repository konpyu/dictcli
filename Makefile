.PHONY: build test lint run clean install help

# Variables
BINARY_NAME=dictcli
BINARY_DIR=bin
GO=go
GOFLAGS=
LDFLAGS=-ldflags "-w -s"

# Default target
all: lint test build

# Build the application
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BINARY_DIR}
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o ${BINARY_DIR}/${BINARY_NAME} ./cmd/dictcli

# Run the application
run: build
	@echo "Running ${BINARY_NAME}..."
	./${BINARY_DIR}/${BINARY_NAME}

# Run with debug flag
debug: build
	@echo "Running ${BINARY_NAME} with debug..."
	./${BINARY_DIR}/${BINARY_NAME} --debug

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Please install it." && exit 1)
	golangci-lint run ./...

# Install golangci-lint
install-linter:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.63.0

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ${BINARY_DIR}
	@rm -f coverage.out coverage.html
	@rm -rf logs/*.log

# Install the binary to GOPATH/bin
install: build
	@echo "Installing ${BINARY_NAME}..."
	$(GO) install ./cmd/dictcli

# Update dependencies
deps:
	@echo "Updating dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GO) mod verify

# Run the application with different settings
run-business:
	./${BINARY_DIR}/${BINARY_NAME} --topic Business --level 700

run-travel:
	./${BINARY_DIR}/${BINARY_NAME} --topic Travel --level 600

run-daily:
	./${BINARY_DIR}/${BINARY_NAME} --topic Daily --level 500

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Run lint, test, and build"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  debug         - Run with debug logging"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run golangci-lint"
	@echo "  install-linter- Install golangci-lint"
	@echo "  fmt           - Format code"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  deps          - Update dependencies"
	@echo "  verify        - Verify dependencies"
	@echo "  help          - Show this help message"