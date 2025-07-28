# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DictCLI is an LLM-powered English dictation practice TUI application designed for Japanese learners. It uses OpenAI APIs for sentence generation, text-to-speech, and grading with Japanese feedback.

## Key Commands

### Development Commands
- `make all` - Run lint, test, and build (MUST pass before completing any step)
- `make build` - Build the application
- `make run` - Build and run the application
- `make debug` - Run with debug logging enabled
- `make test` - Run tests with race detector and coverage
- `make lint` - Run golangci-lint (includes gosec)
- `make fmt` - Format code

### Testing Different Configurations
- `make run-business` - Run with Business topic, level 700
- `make run-travel` - Run with Travel topic, level 600
- `make run-daily` - Run with Daily topic, level 500

## Architecture Overview

### Three-layer Architecture
1. **TUI Layer** (`internal/tui/`) - Bubble Tea-based terminal UI with state management
2. **Service Layer** (`internal/service/`) - Business logic, OpenAI integration, audio handling
3. **Storage Layer** (`internal/storage/`) - Configuration (Viper), history (JSONL), cache (XDG)

### Key Components
- **TUI States**: Welcome → Generating → Playing → Listening → Grading → ShowingResult
- **Services**: DictationService (OpenAI), AudioPlayer, AudioCache
- **Data Types**: DictationSession, Grade (with Japanese explanations), Config

### Important Design Decisions
1. **TUI-First Development**: UI functionality was prioritized over service implementation
2. **Mock Services**: Initial development used mocks to test TUI independently
3. **Japanese-Focused UX**: All feedback and explanations are in Japanese
4. **Privacy-Conscious**: User input stays local, only prompts sent to API

## Current Implementation Status

Based on tasks.md:
- Steps 1-4 completed: Environment setup, logging, types, and TUI implementation
- Step 5 in progress: Mock service layer implementation
- Remaining: Real OpenAI integration, audio caching, statistics, and polish

## Testing Requirements

Every step completion requires `make all` to pass successfully. This includes:
- golangci-lint with gosec security checks
- Unit tests with race detection
- Build verification

## Keyboard Shortcuts
- `R` - Replay audio
- `S` - Open settings
- `Q` - Quit application  
- `N` - Next round / Show answer
- `Enter` - Submit input
- `Esc` - Cancel/back
- Arrow keys - Navigate in settings

## Logging and Debugging
- Custom logger in `internal/logger/` writes to `logs/` directory
- Use `--debug` flag for detailed logging of state transitions
- Log files include timestamps and are gitignored