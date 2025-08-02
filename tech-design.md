# DictCLI Technical Design v2

> **Date:** 2025-08-02
> **Node Version:** v20 or higher  
> **TUI Framework:** Ink v4  
> **Architecture:** Three-layer Architecture  
> **TypeScript:** v5.x with strict mode  
> **Type Definitions:** Always use @types/react and ink's built-in types  

## 1. Architecture Overview

```
┌─────────────────────────────┐
│   TUI Layer (Ink)           │ - State management, UI rendering
├─────────────────────────────┤
│      Service Layer          │ - OpenAI API, Audio, Orchestration  
├─────────────────────────────┤
│      Storage Layer          │ - History, Setting(Config)
└─────────────────────────────┘
```

### 1.1 Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **TUI** | Ink v4, React | Terminal UI rendering |
| **State** | Zustand v4 | Global state management (with Ink v4 workaround) |
| **API** | OpenAI SDK | Text generation, TTS, scoring |
| **Audio** | play-sound | macOS audio playback |
| **Storage** | Node fs + JSONL | Local data persistence |
| **CLI** | Commander.js | Command-line argument parsing |
| **Testing** | Vitest + ink-testing-library | Unit & integration tests |
| **Linting** | ESLint + Prettier | Code quality & formatting |

### 1.2 LLM Model

- TTS: gpt-4o-mini-tts
  - NOTE: using tts-1-hd or tts-1 is prohibited
- Text: gpt-4o-mini

## 2. Directory Structure

```
dictcli/
├── src/
│   ├── components/          # Ink UI components
│   │   ├── App.tsx         # Main app component
│   │   ├── LearningView.tsx
│   │   ├── ResultView.tsx
│   │   ├── SettingsModal.tsx
│   │   ├── SlashCommandMenu.tsx  # NEW: Slash command dropdown
│   │   └── common/
│   │       ├── Header.tsx
│   │       ├── StatusBar.tsx
│   │       └── DebugPanel.tsx   # NEW: Debug mode display
│   ├── services/           # Business logic
│   │   ├── openai/
│   │   │   ├── client.ts   # OpenAI SDK wrapper
│   │   │   ├── generator.ts # Problem generation
│   │   │   ├── tts.ts      # Text-to-speech
│   │   │   └── scorer.ts   # Answer scoring
│   │   ├── audio/
│   │   │   └── player.ts   # Audio playback
│   │   └── orchestrator.ts # Main flow control
│   ├── storage/            # Data persistence
│   │   ├── history.ts      # Learning history
│   │   └── settings.ts     # User preferences
│   ├── store/              # State management
│   │   └── useStore.ts     # Zustand store
│   ├── types/              # TypeScript definitions
│   │   └── index.ts
│   ├── utils/              # Utilities
│   │   ├── cache.ts        # TTS caching
│   │   └── metrics.ts      # WER calculation
│   ├── cli.tsx             # CLI entry point
│   └── index.tsx           # Main app entry
├── tests/
│   ├── unit/
│   └── integration/
├── package.json
├── .gitignore
├── tsconfig.json
├── .eslintrc.js
├── .prettierrc
└── README.md
```

## 3. Component Architecture

### 3.1 App Component
Main container managing view transitions and global keyboard shortcuts.

### 3.2 LearningView Component
Handles audio playback and text input during learning phase.
- Implements slash command detection and menu display
- Auto-complete functionality when "/" is typed

### 3.3 ResultView Component
Displays scoring results with Japanese explanations.
- Implements hybrid keyboard shortcuts (Enter for next, single keys for actions)
- No slash commands needed for better flow

### 3.4 SettingsModal Component
Modal for adjusting learning parameters with arrow key navigation.

### 3.5 SlashCommandMenu Component (NEW)
Dropdown menu for slash commands with:
- Real-time filtering based on input
- Arrow key navigation
- Command descriptions
- Visual highlighting of selected item

### 3.6 State Management with Zustand
Due to Ink v4's custom renderer not supporting `useSyncExternalStore`, we need to implement manual subscription:
- Use `store.subscribe()` for state change detection
- Force re-renders with local state updates
- Subscribe only to specific slices for performance

## 4. Service Layer Design

### 4.1 OpenAI Integration
- Problem generation with CEFR levels (A1 default)
- Text-to-speech with voice selection
  - Display names (ALEX, SARA, etc.) are mapped to OpenAI voice IDs internally
  - Default voice: ALEX (male) → alloy
- Answer scoring with Japanese feedback
- WER calculation and error highlighting

### 4.2 Audio Player Service
- macOS audio playback via afplay
- Speed control support (0.8x-1.2x)
- Audio file caching

### 4.3 Orchestrator Service
- Main flow control
- Round management
- Pre-generation strategy
- State transitions

## 5. Data Models

### 5.1 Core Types

```typescript
// Voice Display Names
type VoiceDisplayName = 'ALEX' | 'SARA' | 'EVAN' | 'NOVA' | 'NICK' | 'FAYE';

// OpenAI Voice IDs
type OpenAIVoice = 'alloy' | 'shimmer' | 'echo' | 'nova' | 'onyx' | 'fable';

// Voice Mapping
const VOICE_MAPPING: Record<VoiceDisplayName, OpenAIVoice> = {
  'ALEX': 'alloy',    // Default male voice
  'SARA': 'shimmer',  // Default female voice
  'EVAN': 'echo',
  'NOVA': 'nova',
  'NICK': 'onyx',
  'FAYE': 'fable'
};

// CEFR Levels
type Level = 'CEFR_A1' | 'CEFR_A2' | 'CEFR_B1' | 'CEFR_B2' | 'CEFR_C1' | 'CEFR_C2';

// Topics
type Topic = 'Business' | 'Tech' | 'Travel' | 'Daily' | 'Technology' | 'Health';

// Settings
interface Settings {
  voice: VoiceDisplayName;
  level: Level;
  topic: Topic;
  wordCount: number; // 5-30
}

// Slash Commands
interface SlashCommand {
  command: '/replay' | '/settings' | '/quit' | '/giveup';
  description: string;
  action: () => void;
}

// Round Data
interface Round {
  id: string;
  sentence: string;
  userInput: string;
  score: number;
  wer: number;
  errors: Array<{
    expected: string;
    actual: string;
    explanation: string; // Japanese
  }>;
  alternatives: string[];
  timestamp: Date;
}
```

### 5.2 Storage Schema
- Settings persistence in JSON
- History tracking in JSONL format
- Cache for TTS audio files

## 6. UI Implementation Details

### 6.1 Slash Command Menu
```typescript
// When user types "/"
- Show dropdown menu below input
- Filter commands as user types
- Highlight selected item
- Handle arrow key navigation
- Execute on Enter
- Dismiss on Escape
```

### 6.2 Result View Keyboard Handling
```typescript
// Hybrid approach
- Enter: Next round (default action)
- N: Next round
- R: Replay audio
- S: Open settings
- Q: Quit application
```

### 6.3 Debug Mode UI
```typescript
// When DEBUG=true
- Show debug panel at top
- Display API calls and responses
- Show timing information
- Log state changes
```

## 7. Performance Optimization

### 7.1 Caching Strategy
- TTS audio caching with 15-minute TTL
- Cache key: hash of (text + voice + speed)
- Maximum cache size: 100MB

### 7.2 Pre-generation
- Generate next problem during current round
- Pre-cache TTS audio
- Reduce perceived latency

## 8. Testing Strategy

### 8.1 Unit Tests
- Service layer testing
- Component testing (including SlashCommandMenu)
- Utility function testing

### 8.2 Integration Tests
- Test OpenAI API integration with mocked responses
- Test audio playback on macOS
- Test file system operations
- Test slash command functionality

### 8.3 E2E Tests with ink-testing-library
- Complete learning flow simulation
- Keyboard navigation testing (both slash commands and shortcuts)
- Settings persistence verification
- Audio playback control testing
- Error state handling

## 9. Configuration

### 9.1 Environment Variables
- OPENAI_API_KEY (Required)
- DICTCLI_DATA_DIR (Optional, default: ~/.dictcli)
- DICTCLI_CACHE_SIZE (Optional, default: 100MB)
- DICTCLI_DEBUG (Optional, enables debug mode)

### 9.2 Build Configuration
- TypeScript compilation with strict mode
- Binary generation for npm distribution
- Development workflow with hot reload

## 10. Deployment

### 10.1 NPM Package
- Global installation support: `npm install -g dictcli`
- CLI binary distribution
- Automatic dependency resolution

### 10.2 Platform Support
- macOS only (native audio support via afplay)
- Node.js v20+

## 11. Security Considerations

- API key stored in environment variable only
- No user data sent to external services except OpenAI
- Local storage only, no cloud sync
- Input sanitization for file paths
- No telemetry or analytics

## 12. Technical Limitations & Workarounds

### 12.1 Ink v4 + Zustand v4 Compatibility
- Ink v4's custom renderer doesn't support React 18's `useSyncExternalStore`
- Zustand's automatic re-rendering with selectors won't work
- Workaround: Manual subscription pattern with forced re-renders
- Future: Wait for Ink v5 which may support `useSyncExternalStore`

### 12.2 Audio Playback
- Limited to macOS afplay command
- No cross-platform audio support in MVP

## 13. Implementation Priorities

1. **Core Learning Flow**
   - Problem generation → Audio playback → Input → Scoring
   
2. **Slash Command UI**
   - Detection and menu display
   - Keyboard navigation
   - Command execution

3. **Result View Shortcuts**
   - Hybrid keyboard handling
   - Smooth transitions

4. **Debug Mode**
   - API call logging
   - Performance metrics

5. **Storage & Settings**
   - Persistence
   - Default values

## 14. Future Considerations (Post-MVP)

- Voice input support using Whisper API
- Spaced repetition for missed words
- Custom problem sets
- Progress analytics dashboard
- Multi-language support (beyond Japanese)
- Cross-platform audio support
- Migration to Ink v5 when available (remove Zustand workaround)
- Theme customization
- GEMINI API support