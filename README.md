# DictCLI

LLM-powered English dictation practice tool for Japanese learners.

## Prerequisites

### Required
- Go 1.22 or later
- OpenAI API key (set as `OPENAI_API_KEY` environment variable)

### Development Tools
```bash
# Install golangci-lint (required for linting)
brew install golangci-lint

# Or using go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Installation

```bash
git clone https://github.com/konpyu/dictcli.git
cd dictcli
make build
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Run everything (deps, lint, test, build)
make all

# See all available commands
make help
```

## Usage

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="your-api-key-here"

# Run the application
./dictcli
```

## Project Structure

```
dictcli/
├── cmd/dictcli/       # CLI entry point
├── internal/          # Internal packages
│   ├── config/        # Configuration management
│   ├── service/       # OpenAI service layer
│   ├── storage/       # Data persistence (TBD)
│   ├── tui/           # Terminal UI (TBD)
│   └── types/         # Shared types
├── Makefile           # Build automation
├── go.mod             # Go module definition
└── README.md          # This file
```