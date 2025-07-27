.PHONY: build test lint clean install run help

BINARY_NAME=dictcli
GO=go
GOFLAGS=-v
MAIN_PKG=./cmd/dictcli

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  test         - Run tests"
	@echo "  lint         - Run golangci-lint (requires golangci-lint installed)"
	@echo "  lint-install - Install golangci-lint using go install"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install the binary"
	@echo "  run          - Run the application"

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PKG)

test:
	$(GO) test -v ./...

test-coverage:
	$(GO) test -v -cover ./...

test-race:
	$(GO) test -race ./...

lint:
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint not found. Please install it with: brew install golangci-lint" && exit 1)
	golangci-lint run --enable=gosec --enable=staticcheck

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "golangci-lint installed to $$GOPATH/bin"

clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

install:
	$(GO) install $(GOFLAGS) $(MAIN_PKG)

run: build
	./$(BINARY_NAME)

deps:
	$(GO) mod tidy
	$(GO) mod download

all: deps lint test build