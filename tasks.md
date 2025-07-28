# DictCLI Implementation Tasks

## Step 1: Golang Environment Setup

- [x] Initialize Go module (`go mod init github.com/yourusername/dictcli`)
- [x] Create project directory structure
  - [x] `cmd/dictcli/` - CLI entry points
  - [x] `internal/config/` - Configuration management
  - [x] `internal/service/` - Business logic
  - [x] `internal/storage/` - Data persistence
  - [x] `internal/tui/` - Terminal UI
  - [x] `internal/types/` - Shared types
  - [x] `internal/logger/` - Logging infrastructure
- [x] Install core dependencies
  - [x] `go get github.com/charmbracelet/bubbletea`
  - [x] `go get github.com/sashabaranov/go-openai`
  - [x] `go get github.com/spf13/cobra`
  - [x] `go get github.com/spf13/viper`
  - [x] `go get github.com/adrg/xdg`
- [x] Set up golangci-lint (v2.3.0) configuration(with gosec)
- [x] Create Makefile with common tasks (build, test, lint, run)
- [x] reasonable .gitignore (with logs/)

## Step 2: Logging Infrastructure

- [x] Create custom logger package (`internal/logger/`)
  - [x] Implement file-based logging with rotation
  - [x] Add log levels (DEBUG, INFO, WARN, ERROR)
  - [x] Create structured logging format for TUI debugging
  - [x] Implement conditional logging based on --debug flag
- [x] Set up logs directory structure
  - [x] Create logs/ directory in project root
  - [x] Implement log file naming with timestamps
- [x] Add logging points for TUI state transitions
- [x] Test logging with sample debug messages

## Step 3: Core Types and Interfaces

- [x] Define core data structures (`internal/types/`)
  - [x] DictationSession struct
  - [x] Grade struct with Japanese explanation fields
  - [x] Config struct with all settings
  - [x] Mistake struct for error tracking
- [x] Create service interfaces
  - [x] DictationService interface
  - [x] AudioPlayer interface
  - [x] AudioCache interface
  - [x] Storage interface

## Step 4: TUI Implementation (Priority)

- [ ] Create TUI package structure (`internal/tui/`)
  - [ ] `model.go` - Main model with state management
  - [ ] `update.go` - Message handling and state transitions
  - [ ] `view.go` - Rendering logic for each state
  - [ ] `styles.go` - Consistent styling definitions
  - [ ] `messages.go` - Message types for bubble tea
- [ ] Implement all TUI states
  - [ ] StateWelcome - Welcome screen with tips
  - [ ] StateGenerating - Loading animation while generating
  - [ ] StatePlaying - Audio playback indicator
  - [ ] StateListening - Input field for dictation
  - [ ] StateGrading - Loading while grading
  - [ ] StateShowingResult - Result display with Japanese feedback
  - [ ] StateSettings - Settings modal with arrow key navigation
- [ ] Implement keyboard shortcuts
  - [ ] R - Replay audio
  - [ ] S - Open settings
  - [ ] Q - Quit application
  - [ ] N - Next round (and Show answer)
  - [ ] Enter - Submit input
  - [ ] Esc - Cancel/back
- [ ] Add proper error handling and user feedback
- [ ] Test all state transitions thoroughly

## Step 5: Mock Service Layer

- [ ] Create mock implementations (`internal/service/mock/`)
  - [ ] MockDictationService
    - [ ] GenerateSentence - Return predefined sentences
    - [ ] GenerateAudio - Return dummy audio file path
    - [ ] GradeDictation - Return fake grading results with Japanese
  - [ ] MockAudioPlayer - Log play commands without actual playback
  - [ ] MockAudioCache - Simulate cache hits/misses
- [ ] Create test data
  - [ ] Sample sentences for different levels/topics
  - [ ] Sample grading results with various mistake types
  - [ ] Japanese explanations for common errors
- [ ] Wire mock services into TUI for testing

## Step 6: Configuration Management

- [ ] Implement config package (`internal/config/`)
  - [ ] Default configuration values
  - [ ] Viper integration for file/env/flag support
  - [ ] Config validation
- [ ] Create config file structure
  - [ ] YAML configuration schema
  - [ ] XDG-compliant config location
- [ ] Implement settings persistence
  - [ ] Save settings when modified in TUI
  - [ ] Load settings on startup

## Step 7: Storage Layer

- [ ] Implement history storage (`internal/storage/`)
  - [ ] JSONL file handling
  - [ ] Session serialization/deserialization
  - [ ] Append-only write operations
- [ ] Create statistics calculator
  - [ ] WER trends over time
  - [ ] Common mistakes aggregation
  - [ ] Progress by topic/level
- [ ] Implement cache directory management
  - [ ] XDG cache paths
  - [ ] Cache cleanup utilities

## Step 8: CLI Setup

- [ ] Create main command (`cmd/dictcli/main.go`)
  - [ ] Cobra command structure
  - [ ] Flag definitions matching spec
  - [ ] Environment variable support
- [ ] Implement subcommands
  - [ ] `stats` - Show learning statistics
  - [ ] `config` - Display current configuration
  - [ ] `clear-cache` - Clear audio cache
- [ ] Add version information and help text

## Step 9: Audio Player Integration

- [ ] Implement cross-platform audio player (`internal/service/audio/`)
  - [ ] macOS - afplay command
  - [ ] Linux - mpg123 command
  - [ ] Windows - PowerShell audio playback
- [ ] Add audio player detection
  - [ ] Check available commands
  - [ ] Fallback handling
- [ ] Implement background playback
  - [ ] Non-blocking audio play
  - [ ] Playback status tracking

## Step 10: Real Service Implementation

- [ ] Replace mock with real DictationService
  - [ ] OpenAI client setup with API key
  - [ ] Implement GenerateSentence with GPT-4o-mini
    - [ ] Prompt engineering for level/topic/length
    - [ ] Temperature and parameter tuning
  - [ ] Implement GenerateAudio with TTS-1
    - [ ] Voice selection support
    - [ ] Speed adjustment
    - [ ] MP3 file generation
  - [ ] Implement GradeDictation
    - [ ] WER calculation
    - [ ] Japanese explanation generation
    - [ ] Alternative expression suggestions
- [ ] Add retry logic with exponential backoff
- [ ] Implement rate limit handling
- [ ] Add proper error messages for API failures

## Step 11: Audio Caching

- [ ] Implement audio cache service
  - [ ] SHA256-based cache keys
  - [ ] Cache hit detection
  - [ ] Cache miss handling
  - [ ] File system operations
- [ ] Add cache management features
  - [ ] Size tracking
  - [ ] Age-based cleanup
  - [ ] Manual cache clearing

## Step 12: Testing and Validation

- [ ] Unit tests for core components
  - [ ] Service layer tests
  - [ ] Storage layer tests
  - [ ] Config validation tests
- [ ] Integration tests
  - [ ] TUI state machine tests
  - [ ] End-to-end flow tests
- [ ] Manual testing checklist
  - [ ] All keyboard shortcuts
  - [ ] Error scenarios
  - [ ] Cross-platform compatibility

## Step 13: Polish and Optimization

- [ ] Performance optimization
  - [ ] Reduce API latency perception
  - [ ] Optimize TUI rendering
  - [ ] Memory usage profiling
- [ ] UX improvements
  - [ ] Loading animations
  - [ ] Smooth transitions
  - [ ] Clear error messages
- [ ] Documentation
  - [ ] README.md with usage instructions
  - [ ] API key setup guide
  - [ ] Troubleshooting section

## Step 14: Release Preparation

- [ ] Build scripts for multiple platforms
  - [ ] macOS (Intel/ARM)
  - [ ] Linux (x64/ARM)
  - [ ] Windows (x64)
- [ ] Create installation instructions
- [ ] Set up GitHub releases workflow
- [ ] Create demo video/screenshots

## Priority Notes

1. **TUI First**: Complete TUI with mocked services to ensure smooth user experience
2. **Logging Early**: Set up logging before TUI to aid in debugging
3. **Mock Everything**: Use mocks for all external dependencies initially
4. **Iterate**: Replace mocks with real implementations one by one
5. **Test Continuously**: Validate each component as it's built

## Success Criteria

- [ ] TUI responds to all inputs within 100ms (excluding API calls)
- [ ] All states transition smoothly without flicker
- [ ] Japanese feedback is clear and helpful
- [ ] Audio playback works on all platforms
- [ ] Settings persist between sessions
- [ ] History tracking works accurately
- [ ] API errors are handled gracefully

## Important!

- After completing each STEP, ensure that `make all` passes successfully. Passing is the condition for completion.