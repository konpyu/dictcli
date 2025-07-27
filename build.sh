#!/bin/bash

# Build script for DictCLI with version information

set -e

# Default values
VERSION=${VERSION:-"dev"}
OUTPUT_DIR=${OUTPUT_DIR:-"./dist"}

# Get build information
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
GO_VERSION=$(go version | awk '{print $3}')

# Ensure output directory exists
mkdir -p "$OUTPUT_DIR"

# Build flags with version information
LDFLAGS="-s -w \
  -X 'github.com/konpyu/dictcli/internal/version.Version=${VERSION}' \
  -X 'github.com/konpyu/dictcli/internal/version.GitCommit=${GIT_COMMIT}' \
  -X 'github.com/konpyu/dictcli/internal/version.BuildTime=${BUILD_TIME}'"

echo "Building DictCLI..."
echo "Version: $VERSION"
echo "Git Commit: $GIT_COMMIT"
echo "Build Time: $BUILD_TIME"
echo "Go Version: $GO_VERSION"
echo ""

# Build for current platform
echo "Building for current platform..."
go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli" ./cmd/dictcli

# Cross-compile for multiple platforms if requested
if [[ "$1" == "all" ]]; then
    echo "Cross-compiling for multiple platforms..."
    
    # Linux amd64
    echo "Building for Linux amd64..."
    GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli-linux-amd64" ./cmd/dictcli
    
    # Linux arm64
    echo "Building for Linux arm64..."
    GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli-linux-arm64" ./cmd/dictcli
    
    # macOS amd64
    echo "Building for macOS amd64..."
    GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli-darwin-amd64" ./cmd/dictcli
    
    # macOS arm64 (Apple Silicon)
    echo "Building for macOS arm64..."
    GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli-darwin-arm64" ./cmd/dictcli
    
    # Windows amd64
    echo "Building for Windows amd64..."
    GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli-windows-amd64.exe" ./cmd/dictcli
    
    # Windows arm64
    echo "Building for Windows arm64..."
    GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -o "$OUTPUT_DIR/dictcli-windows-arm64.exe" ./cmd/dictcli
    
    echo ""
    echo "Cross-compilation complete. Binaries available in $OUTPUT_DIR:"
    ls -la "$OUTPUT_DIR"
fi

echo ""
echo "Build complete!"
echo "Binary location: $OUTPUT_DIR/dictcli"

# Test the binary
if [[ -x "$OUTPUT_DIR/dictcli" ]]; then
    echo ""
    echo "Version information:"
    "$OUTPUT_DIR/dictcli" version
fi