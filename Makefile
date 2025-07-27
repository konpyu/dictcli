.PHONY: build build-all test lint clean install run help release

BINARY_NAME=dictcli
GO=go
GOFLAGS=-v
MAIN_PKG=./cmd/dictcli
DIST_DIR=./dist

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-all    - Build binaries for all supported platforms"
	@echo "  test         - Run tests"
	@echo "  lint         - Run golangci-lint (requires golangci-lint installed)"
	@echo "  lint-install - Install golangci-lint using go install"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install the binary"
	@echo "  run          - Run the application"
	@echo "  release      - Build release binaries with version info"

build:
	./build.sh

build-all:
	./build.sh all

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
	rm -rf $(DIST_DIR)
	$(GO) clean

install:
	$(GO) install $(GOFLAGS) $(MAIN_PKG)

run: build
	$(DIST_DIR)/$(BINARY_NAME)

release: clean deps lint test build-all
	@echo "Release build complete. Binaries available in $(DIST_DIR)/"

deps:
	$(GO) mod tidy
	$(GO) mod download

all: deps lint test build