.PHONY: build test lint clean install run help

BINARY_NAME=dictcli
GO=go
GOFLAGS=-v
MAIN_PKG=./cmd/dictcli

help:
	@echo "Available targets:"
	@echo "  build     - Build the binary"
	@echo "  test      - Run tests"
	@echo "  lint      - Run golangci-lint"
	@echo "  clean     - Clean build artifacts"
	@echo "  install   - Install the binary"
	@echo "  run       - Run the application"

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PKG)

test:
	$(GO) test -v ./...

test-coverage:
	$(GO) test -v -cover ./...

test-race:
	$(GO) test -race ./...

lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, installing..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

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