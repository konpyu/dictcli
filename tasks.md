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
- [ ] Create OpenAI client wrapper
  - [ ] Initialize client with API key
  - [ ] Add retry logic with exponential backoff
- [ ] Implement sentence generation
  - [ ] Create GenerateSentence function
  - [ ] Build prompt based on level/topic/word count
  - [ ] Handle API errors gracefully
- [ ] Implement text-to-speech
  - [ ] Create GenerateAudio function
  - [ ] Support multiple voices
  - [ ] Support variable speed
- [ ] Implement grading service
  - [ ] Create GradeDictation function
  - [ ] Parse JSON response
  - [ ] Calculate WER and mistakes
- [ ] Add unit tests for all service functions
- [ ] Run tests and linting

## Step 3: Storage Layer
- [ ] Implement audio cache
  - [ ] Create AudioCache struct
  - [ ] Implement SHA256-based key generation
  - [ ] Add cache hit/miss detection
  - [ ] Create save/load methods
- [ ] Implement history storage
  - [ ] Create History struct
  - [ ] Implement JSONL append
  - [ ] Add session save method
  - [ ] Add query methods for stats
- [ ] Implement cross-platform audio player
  - [ ] Detect OS and available players
  - [ ] Create play method with fallbacks
  - [ ] Handle player not found errors
- [ ] Add unit tests for storage components
- [ ] Run tests and linting

## Step 4: Basic TUI Framework
- [ ] Create main TUI model structure
  - [ ] Define state enum
  - [ ] Create Model struct
  - [ ] Initialize Bubble Tea components
- [ ] Implement state management
  - [ ] Create state transition logic
  - [ ] Handle keyboard input
  - [ ] Add quit functionality
- [ ] Create basic view rendering
  - [ ] Implement View() method
  - [ ] Add header with settings display
  - [ ] Create state-specific views
- [ ] Setup basic update cycle
  - [ ] Handle tea.Msg types
  - [ ] Route to state handlers
  - [ ] Update UI components
- [ ] Test TUI navigation flow
- [ ] Run tests and linting

## Step 5: Core Dictation Flow
- [ ] Implement sentence generation state
  - [ ] Show spinner during generation
  - [ ] Call service layer
  - [ ] Handle generation errors
- [ ] Implement audio playback state
  - [ ] Trigger audio generation
  - [ ] Play audio automatically
  - [ ] Show playback status
- [ ] Implement listening/input state
  - [ ] Create text input component
  - [ ] Handle Enter to submit
  - [ ] Add replay hotkey (R)
- [ ] Implement grading state
  - [ ] Show spinner during grading
  - [ ] Call grading service
  - [ ] Handle grading errors
- [ ] Implement result display state
  - [ ] Show score and WER
  - [ ] Display mistakes with Japanese explanations
  - [ ] Show alternative expressions
  - [ ] Add next/replay options
- [ ] Integration test full flow
- [ ] Run tests and linting

## Step 6: Settings & Configuration UI
- [ ] Create settings modal
  - [ ] Design settings layout
  - [ ] Add navigation with arrow keys
  - [ ] Implement value adjustment
- [ ] Add voice selection
  - [ ] List available voices
  - [ ] Update config on change
- [ ] Add level selection
  - [ ] Support TOEIC levels 400-990
  - [ ] Validate input range
- [ ] Add topic selection
  - [ ] Show topic list
  - [ ] Cycle through options
- [ ] Add word count adjustment
  - [ ] Support 5-30 words
  - [ ] Use +/- keys
- [ ] Persist settings changes
  - [ ] Save to config file
  - [ ] Apply immediately
- [ ] Test settings persistence
- [ ] Run tests and linting

## Step 7: CLI Commands & Flags
- [ ] Implement main command
  - [ ] Setup Cobra root command
  - [ ] Add command description
  - [ ] Initialize TUI on run
- [ ] Add CLI flags
  - [ ] --level flag
  - [ ] --words flag
  - [ ] --topic flag
  - [ ] --voice flag
  - [ ] --speed flag
  - [ ] --no-cache flag
  - [ ] --debug flag
- [ ] Implement stats command
  - [ ] Query history data
  - [ ] Calculate statistics
  - [ ] Format output nicely
- [ ] Implement config command
  - [ ] Display current configuration
  - [ ] Show config file location
- [ ] Implement clear-cache command
  - [ ] Count cached files
  - [ ] Confirm before deletion
  - [ ] Delete cache files
- [ ] Add command help text
- [ ] Test all CLI commands
- [ ] Run tests and linting

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