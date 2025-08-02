# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Core Development Workflow
- `npm run dev` - Start development server with tsx
- `npm run build` - Build TypeScript to dist/ folder
- `npm test` - Run all tests with Vitest
- `npm run test:watch` - Run tests in watch mode
- `npm run lint` - Lint TypeScript files with ESLint
- `npm run format` - Format code with Prettier

### Running the Application
- Development: `npm run dev` or with options: `npm run dev -- --topic Technology --level CEFR_B1 --words 15`
- Production: Build first (`npm run build`), then run `./dist/cli.js`
- CLI options: `--topic <topic> --level <level> --words <number> --voice <voice>`

### Testing
- Run single test: `npm test -- tests/unit/components/App.test.tsx`
- Run with coverage: `npm test -- --coverage`
- Watch mode: `npm run test:watch`

### Environment Setup
**Required**: `OPENAI_API_KEY` environment variable must be set
**Optional**: `DICTCLI_DEBUG=true` to enable debug mode in development

## Architecture Overview

This is a **LLM-First Dictation TUI App** for Japanese English learners built with:
- **TUI Framework**: Ink v4 + React for terminal UI
- **State Management**: Zustand v4 with custom Ink compatibility layer
- **LLM Integration**: OpenAI API for problem generation, TTS, and scoring
- **Audio**: macOS-only audio playback via play-sound/afplay

### Three-Layer Architecture
1. **TUI Layer** (`src/components/`) - Ink React components for terminal UI
2. **Service Layer** (`src/services/`) - OpenAI API integration and audio handling
3. **Storage Layer** (`src/storage/`) - Local file persistence (JSONL history, JSON settings)

### Key Components
- `App.tsx` - Main container with view state management and global keyboard handling
- `LearningView.tsx` - Audio playback and text input with slash command menu
- `ResultView.tsx` - Score display with hybrid keyboard shortcuts (Enter/single keys)
- `SettingsModal.tsx` - Configuration modal with arrow key navigation
- `SlashCommandMenu.tsx` - Dropdown for `/replay`, `/settings`, `/quit`, `/giveup` commands
- `orchestrator.ts` - Main flow controller coordinating all services

### State Management Pattern (Important!)
Due to Ink v4 incompatibility with React 18's `useSyncExternalStore`, this project uses a **custom Zustand wrapper** in `src/store/useStore.ts`:
- Manual subscription pattern with forced re-renders
- Use `useStore(selector)` hook for components
- Direct access via `store.getState()` for services
- The `useStore` hook implements Ink v4 compatibility workarounds

### OpenAI Integration Details
- **Problem Generation**: Uses gpt-4o-mini with CEFR levels (A1-C2) and topics
- **TTS**: Uses OpenAI TTS with voice mapping (ALEX→echo, SARA→shimmer, etc.) - see `VOICE_MAPPING` in types
- **Scoring**: Uses gpt-4o-mini for Japanese-language error explanations and WER calculation
- **Caching**: TTS audio files cached with 15-minute TTL in temp directory

## Testing Strategy

### Unit Tests
- Components: `tests/unit/components/` using ink-testing-library
- Services: `tests/unit/services/` with mocked OpenAI API calls
- Store: `tests/unit/store/` testing Zustand state management

### Test Configuration
- Framework: Vitest with Node environment
- Coverage: v8 provider with text/json/html reporters
- Mock strategy: Mock OpenAI API responses, use real audio for integration tests

## Development Guidelines

### Code Style
- TypeScript strict mode enabled
- ESLint + Prettier for consistent formatting
- React functional components with hooks
- Async/await for promises (avoid .then/.catch)

### LLM Model Usage
- **Text Generation/Scoring**: gpt-4o-mini only
- **TTS**: OpenAI TTS API (NOT deprecated tts-1 or tts-1-hd models)
- All API keys from environment variables only

### Audio Handling
- **Platform**: macOS only (uses afplay command)
- **Speed Control**: 0.8x to 1.2x supported
- **File Management**: Temporary files auto-cleaned after 15 minutes

### Debug Mode
When `DICTCLI_DEBUG=true`:
- DebugPanel component shows API calls, timing, and state changes at top of UI
- All debug logs available in store (`debugState.logs`)
- Essential for development and troubleshooting

### Performance Optimizations
- **Pre-generation**: Next problem generated during current round to reduce latency
- **TTS Caching**: Audio files cached by hash of (text + voice + speed)
- **State Slicing**: Use specific selectors in `useStore()` to minimize re-renders

## Important Implementation Notes

### Voice Configuration
Voice display names (ALEX, SARA) map to OpenAI voice IDs internally via `VOICE_MAPPING` in types. Always use display names in UI.

### View State Management
Three main views: 'learning' | 'result' | 'settings'
- Learning: Slash command menu for `/replay`, `/settings`, etc.
- Result: Hybrid shortcuts (Enter for next, R for replay, S for settings, Q for quit)
- Settings: Arrow keys for navigation, Enter to save, Esc to cancel

### Error Handling
- All API calls wrapped in try/catch with debug logging
- Audio errors logged but don't crash the app
- User-facing errors shown in Japanese when appropriate

### File Structure Standards
- Components: PascalCase, default exports
- Services: camelCase classes with singleton instances
- Types: Centralized in `src/types/index.ts`
- Tests: Mirror source structure in `tests/` directory