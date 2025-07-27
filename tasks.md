# DictCLI Implementation Tasks

## Step 1: Project Setup & Foundation
- [x] Initialize Go project structure
  - [x] Create directory structure (cmd/, internal/, etc.)
  - [x] Initialize go.mod with module name
  - [x] Add .gitignore for Go projects
- [x] Setup core dependencies
  - [x] Add Bubble Tea framework
  - [x] Add OpenAI SDK
  - [x] Add Cobra/Viper for CLI
  - [x] Add XDG for cross-platform paths
- [x] Create basic type definitions
  - [x] Define Config struct
  - [x] Define DictationSession struct
  - [x] Define Grade struct
  - [x] Define Mistake struct
- [x] Implement configuration management
  - [x] Create config loader with Viper
  - [x] Setup default values
  - [x] Add environment variable support
- [x] Setup project Makefile
  - [x] Add build target
  - [x] Add test target
  - [x] Add lint target (golangci-lint)
- [x] Run tests and linting

## Step 2: Service Layer - OpenAI Integration
- [x] Create OpenAI client wrapper
  - [x] Initialize client with API key
  - [x] Add retry logic with exponential backoff
- [x] Implement sentence generation
  - [x] Create GenerateSentence function
  - [x] Build prompt based on level/topic/word count
  - [x] Handle API errors gracefully
- [x] Implement text-to-speech
  - [x] Create GenerateAudio function
  - [x] Support multiple voices
  - [x] Support variable speed
- [x] Implement grading service
  - [x] Create GradeDictation function
  - [x] Parse JSON response
  - [x] Calculate WER and mistakes
- [x] Add unit tests for all service functions
- [x] Run tests and linting

## Step 3: Storage Layer
- [x] Implement audio cache
  - [x] Create AudioCache struct
  - [x] Implement SHA256-based key generation
  - [x] Add cache hit/miss detection
  - [x] Create save/load methods
- [x] Implement history storage
  - [x] Create History struct
  - [x] Implement JSONL append
  - [x] Add session save method
  - [x] Add query methods for stats
- [x] Implement cross-platform audio player
  - [x] Detect OS and available players
  - [x] Create play method with fallbacks
  - [x] Handle player not found errors
- [x] Add unit tests for storage components
- [x] Run tests and linting

## Step 4: Basic TUI Framework
- [x] Create main TUI model structure
  - [x] Define state enum
  - [x] Create Model struct
  - [x] Initialize Bubble Tea components
- [x] Implement state management
  - [x] Create state transition logic
  - [x] Handle keyboard input
  - [x] Add quit functionality
- [x] Create basic view rendering
  - [x] Implement View() method
  - [x] Add header with settings display
  - [x] Create state-specific views
- [x] Setup basic update cycle
  - [x] Handle tea.Msg types
  - [x] Route to state handlers
  - [x] Update UI components
- [x] Test TUI navigation flow
- [x] Run tests and linting

## Step 5: Core Dictation Flow
- [x] Implement sentence generation state
  - [x] Show spinner during generation
  - [x] Call service layer
  - [x] Handle generation errors
- [x] Implement audio playback state
  - [x] Trigger audio generation
  - [x] Play audio automatically
  - [x] Show playback status
- [x] Implement listening/input state
  - [x] Create text input component
  - [x] Handle Enter to submit
  - [x] Add replay hotkey (R)
- [x] Implement grading state
  - [x] Show spinner during grading
  - [x] Call grading service
  - [x] Handle grading errors
- [x] Implement result display state
  - [x] Show score and WER
  - [x] Display mistakes with Japanese explanations
  - [x] Show alternative expressions
  - [x] Add next/replay options
- [x] Integration test full flow
- [x] Run tests and linting

## Step 6: Settings & Configuration UI
- [x] Create settings modal
  - [x] Design settings layout
  - [x] Add navigation with arrow keys
  - [x] Implement value adjustment
- [x] Add voice selection
  - [x] List available voices
  - [x] Update config on change
- [x] Add level selection
  - [x] Support TOEIC levels 400-990
  - [x] Validate input range
- [x] Add topic selection
  - [x] Show topic list
  - [x] Cycle through options
- [x] Add word count adjustment
  - [x] Support 5-30 words
  - [x] Use +/- keys
- [x] Persist settings changes
  - [x] Save to config file
  - [x] Apply immediately
- [x] Test settings persistence
- [x] Run tests and linting

## Step 7: CLI Commands & Flags
- [x] Implement main command
  - [x] Setup Cobra root command
  - [x] Add command description
  - [x] Initialize TUI on run
- [x] Add CLI flags
  - [x] --level flag
  - [x] --words flag
  - [x] --topic flag
  - [x] --voice flag
  - [x] --speed flag
  - [x] --no-cache flag
  - [x] --debug flag
- [x] Implement stats command
  - [x] Query history data
  - [x] Calculate statistics
  - [x] Format output nicely
- [x] Implement config command
  - [x] Display current configuration
  - [x] Show config file location
- [x] Implement clear-cache command
  - [x] Count cached files
  - [x] Confirm before deletion
  - [x] Delete cache files
- [x] Add command help text
- [x] Test all CLI commands
- [x] Run tests and linting

## Step 8: History & Statistics
- [ ] Implement session saving
  - [ ] Save after each round
  - [ ] Include all session data
  - [ ] Handle save errors
- [ ] Create statistics calculator
  - [ ] Calculate average scores
  - [ ] Track progress over time
  - [ ] Find common mistakes
- [ ] Implement stats display
  - [ ] Overall statistics
  - [ ] Topic breakdown
  - [ ] Common mistakes
  - [ ] Recent progress
- [ ] Add date range filtering
  - [ ] Support custom day ranges
  - [ ] Default to 30 days
- [ ] Test statistics accuracy
- [ ] Run tests and linting

## Step 9: Polish & Error Handling
- [ ] Add comprehensive error handling
  - [ ] Network timeouts
  - [ ] API rate limits
  - [ ] Invalid API responses
  - [ ] File system errors
- [ ] Improve user feedback
  - [ ] Clear error messages
  - [ ] Loading indicators
  - [ ] Success confirmations
- [ ] Add keyboard shortcuts help
  - [ ] Show available hotkeys
  - [ ] Context-sensitive help
- [ ] Optimize performance
  - [ ] Minimize API calls
  - [ ] Efficient cache usage
  - [ ] Fast UI updates
- [ ] Add debug logging
  - [ ] Log API calls when --debug
  - [ ] Log cache hits/misses
  - [ ] Log errors with context
- [ ] Run comprehensive tests
- [ ] Run linting and fix issues

## Step 10: Final Testing & Documentation
- [ ] Create integration tests
  - [ ] Test full user flows
  - [ ] Test error scenarios
  - [ ] Test edge cases
- [ ] Performance testing
  - [ ] Measure startup time (<500ms)
  - [ ] Check response times (<100ms)
  - [ ] Verify memory usage
- [ ] Cross-platform testing
  - [ ] Test on macOS
  - [ ] Test on Linux
  - [ ] Test on Windows
- [ ] Create README.md
  - [ ] Installation instructions
  - [ ] Usage examples
  - [ ] Configuration guide
- [ ] Add inline code documentation
  - [ ] Document public APIs
  - [ ] Add package descriptions
- [ ] Create release build
  - [ ] Add version information
  - [ ] Create binaries for each platform
- [ ] Final linting and test run

## Testing & Linting Commands

Each step should end with:
```bash
# Run tests
go test ./...

# Run linting (requires golangci-lint)
golangci-lint run

# Run specific package tests
go test -v ./internal/service/...

# Run with race detector
go test -race ./...

# Check test coverage
go test -cover ./...
```