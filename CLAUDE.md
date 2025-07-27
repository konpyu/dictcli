# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Build and Run
```bash
# Build the application
make build

# Run the application  
make run
# or directly: ./dictcli

# Clean build artifacts
make clean
```

### Testing
```bash
# Run all tests
make test
# or: go test ./...

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Run tests for specific package
go test -v ./internal/service/...
```

### Linting and Code Quality
```bash
# Run linter (requires golangci-lint)
make lint

# Install golangci-lint if not available
make lint-install

# Run everything (deps, lint, test, build)
make all
```

### Dependencies
```bash
# Update dependencies
make deps
```

## Project Architecture

DictCLI is a Go-based Terminal UI (TUI) application for English dictation practice, targeting Japanese learners. It uses a three-layer architecture:

### Core Components

1. **TUI Layer** (`internal/tui/`)
   - Built with Bubble Tea framework for interactive terminal UI
   - State-based UI management with clear state transitions
   - Main states: Welcome → Generating → Playing → Listening → Grading → ShowingResult → Settings
   - Uses `charmbracelet/bubbles` for UI components (textinput, spinner)

2. **Service Layer** (`internal/service/`)
   - `DictationService`: Main orchestrator for dictation flow
   - `OpenAIService`: Handles OpenAI API integration (GPT-4o-mini for text, TTS-1 for audio)
   - Provides sentence generation, text-to-speech, and grading with Japanese explanations

3. **Storage Layer** (`internal/storage/`)
   - `AudioCache`: SHA256-based caching for TTS audio files (XDG-compliant paths)
   - `History`: JSONL-based session logging for statistics and progress tracking with comprehensive statistics calculation
   - `AudioPlayer`: Cross-platform audio playback (afplay/mpg123/powershell)

### Key Data Structures (`internal/types/`)

- `DictationSession`: Complete session metadata including timing, grading, and user input
- `Config`: User preferences (voice, TOEIC level 400-990, topic, word count 5-30, speed 0.5-2.0)
- `Grade`: Grading results with WER (Word Error Rate), mistakes array, Japanese explanations

### Configuration Management (`internal/config/`)
- Uses Viper for YAML configuration with XDG-compliant storage
- Environment variable support with `DICTCLI_` prefix  
- CLI flag overrides supported
- Default config path: `~/.config/dictcli/config.yaml`

## OpenAI Integration

### API Requirements
- Set `OPENAI_API_KEY` environment variable
- Uses GPT-4o-mini for sentence generation and grading (cost-optimized)
- TTS-1 model for audio generation with voice options (alloy, echo, fable, onyx, nova, shimmer)

### Prompt Strategy
- Sentence generation includes TOEIC level, topic, and word count constraints
- Grading provides structured JSON with WER calculation and Japanese explanations
- Temperature: 0.7 for generation, 0 for grading (consistency)

## Development Patterns

### State Management
The TUI uses a finite state machine pattern. Each state has specific responsibilities:
- State transitions are explicit and documented in `internal/tui/state.go`
- Update logic is centralized in `internal/tui/update.go` with state-specific handlers
- View rendering is state-dependent in `internal/tui/view.go`

### Error Handling
- All service calls use context with timeouts (30s for API calls)
- Retry logic with exponential backoff for API failures
- Graceful degradation when cache/history operations fail
- User-friendly error messages in the TUI

### Testing Strategy
- Unit tests for all service layer components
- Mock OpenAI client for testing without API calls
- Test files use `_test.go` suffix and are co-located with implementation
- Use `go test -race` to detect race conditions

### Performance Considerations
- Audio caching to minimize TTS API calls and costs
- Efficient JSONL append for history (no full file rewrites)
- Minimal UI updates to maintain responsive interface
- Background audio generation when possible

## CLI Commands and Flags

### Main Command
```bash
dictcli [flags]  # Launch TUI dictation practice
```

### Additional Commands
```bash
dictcli stats [days]               # Show statistics (default: 30 days)
dictcli config                     # Display current configuration
dictcli clear-cache                # Clear audio cache
```

### Supported Flags
- `-v, --voice`: Voice selection (alloy, echo, fable, onyx, nova, shimmer)
- `-l, --level`: TOEIC level (400-990)  
- `-t, --topic`: Topic (Business, Travel, Daily, Technology, Health)
- `-w, --words`: Word count (5-30)
- `-s, --speed`: Speech speed (0.5-2.0)
- `--no-cache`: Disable audio caching
- `--debug`: Enable debug logging

## File Organization

```
dictcli/
├── cmd/dictcli/           # CLI entry point and command setup
├── internal/              # Private application code
│   ├── config/           # Configuration management (Viper)
│   ├── service/          # Business logic and OpenAI integration  
│   ├── storage/          # Data persistence (cache, history, audio)
│   ├── tui/              # Terminal UI (Bubble Tea components)
│   └── types/            # Shared type definitions and validation
├── Makefile              # Build automation and common tasks
└── go.mod                # Go module dependencies
```

## Privacy and Local Storage

- All user input and session data stored locally only
- No telemetry or usage tracking
- Audio files cached in XDG cache directory (`~/.cache/dictcli/`)
- History stored in XDG data directory (`~/.local/share/dictcli/`)
- Only generated sentences and user scores sent to OpenAI for grading

## Current Implementation Status

The project has completed Steps 1-8 from the implementation plan:
- Foundation, OpenAI integration, storage layer, and basic TUI framework are complete
- Core dictation flow, settings UI, CLI commands, and history/statistics are complete
- Session saving, statistics calculation, and CLI stats command are implemented
- Ready for Step 9 (Polish & Error Handling) and Step 10 (Final Testing & Documentation)
- See `tasks.md` for detailed implementation checklist

## Statistics and History

### Session Tracking
- Automatic saving of each dictation round with complete session metadata
- Timing data (start/end time, duration), user input, grading results
- Session data stored in JSONL format for efficient querying

### Statistics Features
- Overall statistics: total sessions/rounds, average score and WER
- Topic breakdown: per-topic performance metrics
- Common mistakes tracking: frequency analysis of repeated errors
- Date range filtering: customizable lookback periods (default 30 days)
- Recent progress: daily performance trends