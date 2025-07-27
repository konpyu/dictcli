# DictCLI Technical Design (Simplified)

> **Date:** 2025-07-27  
> **Go Version:** 1.22+  
> **golangci-lint:** 2.3.0
> **Architecture:** Simple Three-layer Architecture  

## 1. Architecture Overview

```
┌─────────────────────────────┐
│   TUI Layer (Bubble Tea)    │ - State management, UI rendering
├─────────────────────────────┤
│      Service Layer          │ - OpenAI API, Audio, Orchestration  
├─────────────────────────────┤
│      Storage Layer          │ - Config (Viper), History (JSONL), Cache (XDG)
└─────────────────────────────┘
```

## 2. Core Dependencies

- `charmbracelet/bubbletea` - TUI framework
- `sashabaranov/go-openai` - OpenAI API client
- `spf13/cobra` & `spf13/viper` - CLI and configuration
- `adrg/xdg` - Cross-platform paths

## 3. Key Data Structures

### DictationSession
- Session metadata (ID, timestamp, config)
- Generated content (sentence, audio path)
- User interaction (input, timing, replay count)
- Grading result (WER, score, mistakes, explanation)

### Grade
- WER (Word Error Rate: 0.0-1.0)
- Score (0-100)
- Mistakes array (position, expected, actual, type)
- Japanese explanation
- Alternative expressions

### Config
- Voice selection (alloy/echo/fable/onyx/nova/shimmer)
- TOEIC level (400-990)
- Topic (Business/Travel/Daily/Technology/Health)
- Word count (5-30)
- Speech speed (0.5-2.0)

## 4. Service Layer Components

### DictationService
- **GenerateSentence**: Creates English sentences via GPT-4o-mini
- **GenerateAudio**: Converts text to speech, with caching
- **GradeDictation**: Evaluates user input and provides feedback

### AudioCache
- SHA256-based cache keys
- XDG-compliant storage location
- Cache hit/miss detection

### AudioPlayer
- Cross-platform support (afplay/mpg123/powershell)
- Background playback

## 5. TUI States

1. **StateWelcome** - Initial welcome screen (shown once)
2. **StateGenerating** - Creating new sentence
3. **StatePlaying** - Playing audio
4. **StateListening** - Accepting user input
5. **StateGrading** - Evaluating response
6. **StateShowingResult** - Displaying feedback
7. **StateSettings** - Configuration menu

### Welcome UI
```
╭─────────────────────────────────────────────────────╮
│ ✻ Welcome to DictCLI!                               │
│                                                     │
│   LLM-powered English dictation practice            │
│   for Japanese learners                             │
│                                                     │
│   Press any key to start...                         │
╰─────────────────────────────────────────────────────╯

 Tips for getting started:

 • Listen carefully to the audio
 • Type what you hear
 • Get instant feedback in Japanese
 • Press 'R' to replay audio anytime
```

## 6. Storage

### History (JSONL)
- Session-by-session append-only log
- Statistics calculation (by period, topic)
- Common mistake tracking
- Progress visualization

### Configuration
- YAML-based config file
- Environment variable support
- CLI flag overrides

### Audio Cache
- MP3 files in XDG cache directory
- Content-based naming (SHA256)

## 7. OpenAI Integration

### Sentence Generation
- Model: GPT-4o-mini
- Temperature: 0.7
- Prompt includes: topic, level, word count

### Grading
- Model: GPT-4o-mini
- Temperature: 0.2
- JSON response format
- Evaluates: accuracy, provides Japanese explanations

### Text-to-Speech
- Model: TTS-1
- Multiple voice options
- Variable speed support

## 8. CLI Commands

```bash
dictcli [flags]                    # Main practice mode
dictcli stats [days]               # Show statistics
dictcli config                     # Display configuration
dictcli clear-cache                # Clear audio cache
```

### Flags
- `-l, --level` - TOEIC level(default: 600)
- `-w, --words` - Word count(default: 7)
- `-t, --topic` - Topic selection(default: Business)
- `-v, --voice` - Voice selection(default: alloy)
- `-s, --speed` - Speech speed
- `--no-cache` - Disable caching
- `--debug` - Debug logging

## 9. User Flow

1. Generate sentence based on config
2. Generate/retrieve audio
3. Play audio (with replay option)
4. User types what they heard
5. Grade the input
6. Show results with Japanese feedback
7. Save session to history
8. Next round or quit

## 10. Error Handling

- Retry with exponential backoff for API calls
- Rate limit detection and longer waits
- Graceful fallbacks for missing audio players

## 11. Key Design Principles

1. **Simplicity First** - Avoid over-engineering
2. **Japanese-focused UX** - Native language feedback for learning
3. **Offline-friendly** - Audio caching, local history
4. **Cross-platform** - Works on macOS, Linux, Windows
5. **Privacy-conscious** - Local storage, no telemetry

## 12. Project Structure

```
dictcli/
├── cmd/dictcli/           # CLI entry points
├── internal/
│   ├── config/           # Configuration management
│   ├── service/          # Business logic
│   ├── storage/          # Data persistence
│   ├── tui/              # Terminal UI
│   └── types/            # Shared types
└── go.mod
```


## 13. Logging

- Prepare a custom logging infrastructure specifically for UI debugging
- Save logs under the logs directory
- When started with --debug, always record details during screen transitions and key operations


## 14. Implementation Notes

- Start with core dictation flow (generate → play → input → grade)
- Add statistics and history as secondary features
- Keep prompts simple and embedded in code
- Use standard Go error handling patterns
- Prefer explicit over implicit behavior
- follow golang best practice
- The UI (TUI) has higher uncertainty in this case. Prioritize designing the task order to ensure the TUI functions correctly first.