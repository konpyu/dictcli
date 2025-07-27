# DictCLI

LLM-powered English dictation practice tool for Japanese learners using Terminal UI.

## Features

DictCLI provides an interactive Terminal UI (TUI) for English dictation practice with:

- **LLM-powered content generation**: GPT-4o-mini generates sentences based on TOEIC level and topic
- **Text-to-speech audio**: OpenAI TTS-1 with multiple voice options and variable speed
- **Intelligent grading**: Automated scoring with Japanese explanations and alternative expressions
- **Audio caching**: Local caching to minimize API costs and enable offline practice
- **Session history**: Track progress and statistics over time
- **Configurable settings**: Customize voice, difficulty level, topics, word count, and speech speed

### Current Implementation Status

✅ **Completed (Steps 1-9)**:
- Foundation and project structure
- OpenAI integration (sentence generation, TTS, grading)
- Storage layer (audio cache, session history, cross-platform audio player)
- Basic TUI framework with Bubble Tea
- Core dictation flow (generate → play → listen → grade → results)
- Settings UI with full configuration management
- CLI commands and flags
- Statistics and analytics
- Polish & error handling

🚧 **In Progress (Step 10)**:
- Final testing & documentation
- Integration tests
- Performance testing
- Cross-platform verification

## Prerequisites

### Required
- Go 1.22 or later
- OpenAI API key (set as `OPENAI_API_KEY` environment variable)
- Audio player (automatically detected):
  - macOS: `afplay` (built-in)
  - Linux: `mpg123` or similar
  - Windows: PowerShell (built-in)

### Development Tools
```bash
# Install golangci-lint (required for linting)
brew install golangci-lint

# Or using go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Installation

### Option 1: From Source
```bash
git clone https://github.com/konpyu/dictcli.git
cd dictcli
make build
```

### Option 2: Direct Build
```bash
git clone https://github.com/konpyu/dictcli.git
cd dictcli
go build -o dictcli ./cmd/dictcli
```

## Quick Start

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="your-api-key-here"

# Run the application
./dictcli

# Or with specific settings
./dictcli --voice nova --level 600 --topic Business

# View statistics
./dictcli stats

# View current configuration
./dictcli config

# Clear audio cache to free space
./dictcli clear-cache
```

### Controls

- **Any key**: Start from welcome screen
- **Ctrl+R**: Replay audio (while typing)
- **Ctrl+S**: Open settings (while typing)
- **Ctrl+Q**: Quit (while typing)
- **Ctrl+C**: Force quit
- **Enter**: Submit answer / Save settings
- **Esc**: Cancel settings

### Settings Navigation

- **↑/↓**: Navigate between settings
- **←/→**: Adjust values
- **Enter**: Save and continue
- **Esc**: Cancel changes

## Configuration

Settings are automatically saved to `~/.config/dictcli/config.yaml`.

### Available Options

- **Voice**: alloy, echo, fable, onyx, nova, shimmer
- **TOEIC Level**: 400-990 (affects sentence difficulty)
- **Topic**: Business, Travel, Daily, Technology, Health
- **Word Count**: 5-30 words per sentence
- **Speech Speed**: 0.5x-2.0x playback speed

### CLI Commands

#### Main Command
```bash
dictcli [flags]                    # Launch TUI dictation practice
```

#### Additional Commands  
```bash
dictcli stats [days]               # Show statistics (default: 30 days)
dictcli config                     # Display current configuration
dictcli clear-cache                # Clear audio cache
```

#### Supported Flags
```bash
-v, --voice string       Voice selection (alloy, echo, fable, onyx, nova, shimmer)
-l, --level int          TOEIC level (400-990)
-t, --topic string       Topic (Business, Travel, Daily, Technology, Health)
-w, --words int          Word count (5-30)
-s, --speed float        Speech speed (0.5-2.0)
    --no-cache           Disable audio caching
    --debug              Enable debug logging
-h, --help               Help for dictcli
```

## Development

### Common Commands

```bash
# Build the application
make build

# Run the application
make run

# Run all tests
make test

# Run linter
make lint

# Run everything (deps, lint, test, build)
make all

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests with verbose output
go test -v ./...

# Run tests with race detector
make test-race

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test -v ./internal/tui/...

# Run integration tests (requires OPENAI_API_KEY)
go test -v ./test/...

# Run performance benchmarks
go test -bench=. ./test/...
```

### Project Architecture

DictCLI uses a three-layer architecture:

1. **TUI Layer** (`internal/tui/`): Bubble Tea-based terminal interface
2. **Service Layer** (`internal/service/`): OpenAI integration and business logic
3. **Storage Layer** (`internal/storage/`): Local data persistence and audio handling

## Project Structure

```
dictcli/
├── cmd/dictcli/           # CLI entry point and command setup
├── internal/              # Private application code
│   ├── config/           # Configuration management (Viper)
│   ├── service/          # Business logic and OpenAI integration
│   ├── storage/          # Data persistence (cache, history, audio)
│   ├── tui/              # Terminal UI (Bubble Tea components)
│   └── types/            # Shared type definitions and validation
├── test/                  # Integration and performance tests
├── CLAUDE.md             # Project instructions for Claude Code
├── Makefile              # Build automation and common tasks
├── go.mod                # Go module dependencies
├── prd.md                # Product Requirements Document
├── tech-design.md        # Technical Design Document
└── tasks.md              # Implementation task checklist
```

## Statistics and Analytics

DictCLI automatically tracks your practice sessions and provides detailed statistics:

### Session Tracking
- Complete session metadata (timing, user input, grading results)
- Per-topic performance breakdown
- Common mistakes tracking with frequency analysis
- Daily progress trends

### Statistics Display
```bash
# View last 30 days (default)
./dictcli stats

# View specific time range
./dictcli stats 7    # Last 7 days
./dictcli stats 90   # Last 90 days
```

### What's Tracked
- Total sessions and rounds completed
- Average scores and Word Error Rate (WER)
- Performance by topic (Business, Travel, etc.)
- Most common mistakes and corrections
- Recent daily progress with trends

## Privacy & Local Storage

- All user input is stored locally only
- No telemetry or usage tracking
- Audio files cached in `~/.cache/dictcli/`
- Session history stored in `~/.local/share/dictcli/`
- Only generated sentences and scores are sent to OpenAI for grading

## Dependencies

- `charmbracelet/bubbletea`: TUI framework
- `charmbracelet/bubbles`: UI components
- `sashabaranov/go-openai`: OpenAI API client
- `spf13/cobra` & `spf13/viper`: CLI and configuration
- `adrg/xdg`: Cross-platform paths