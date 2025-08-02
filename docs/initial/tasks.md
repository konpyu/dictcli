# task list

## NOTE

- Confirm that the Test and Linter have completed successfully at the end of each step
- Check the checkbox at the end of each step

## Step 1: Project Setup & Core Infrastructure
- [x] Initialize npm project with TypeScript configuration
- [x] Set up directory structure as defined in tech-design.md
- [x] Install core dependencies (ink v4, zustand v4, commander.js, openai sdk)
- [x] Configure TypeScript with strict mode
- [x] Set up ESLint and Prettier configuration
- [x] Create basic CLI entry point (cli.tsx)
- [x] Implement environment variable loading (OPENAI_API_KEY)
- [x] Set up vitest configuration for testing
- [x] Run `npm run lint` and ensure no errors
- [x] Run `npm run test` and ensure basic test setup works

## Step 2: Storage Layer Implementation
- [x] Implement settings.ts for user preferences persistence
- [x] Implement history.ts for learning history (JSONL format)
- [x] Create default settings configuration
- [x] Add unit tests for storage layer
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 3: OpenAI Service Layer
- [x] Implement OpenAI client wrapper (client.ts)
- [x] Create problem generator service (generator.ts) with CEFR levels
- [x] Implement TTS service (tts.ts) with voice mapping
- [x] Create scorer service (scorer.ts) with Japanese feedback
- [x] Implement audio caching mechanism
- [x] Add unit tests for OpenAI services with mocked responses
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Human Quality Check 1: OpenAI API Verification
- [x] Create test script for Text API (test-openai-text.ts)
  - Generate sample sentences at different CEFR levels
  - Test scoring functionality with Japanese feedback
  - Verify response format and error handling
- [x] Create test script for Audio API (test-openai-audio.ts)
  - Generate TTS for all voice options
  - Test different playback speeds
  - Verify audio file generation and caching
- [x] Execute both test scripts and confirm API connectivity
- [x] Document any API-related issues or limitations

## Step 4: Audio Player Service
- [x] Implement audio player service using play-sound
- [x] Add macOS afplay support
- [x] Implement playback speed control
- [x] Add error handling for audio playback
- [x] Create unit tests for audio service
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Phase 1: 基本UIフレームワークと状態管理

## Step 5: State Management with Zustand (Minimal)
- [x] Implement minimal useStore.ts (only viewState management)
- [x] Implement Ink v4 workaround (manual subscription)
- [x] Create simple state change test
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 6: Core Static UI Components
- [x] Create Header.tsx (static display only)
- [x] Create StatusBar.tsx (static display only)
- [x] Basic styling with Box/Text
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 7: Minimal App.tsx with View Switching
- [x] Create App.tsx with screen switching only
- [x] Implement keyboard navigation (Q to quit, S for settings)
- [x] Check for memory leaks and infinite loops
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## 🔍 Human Quality Check 1: Basic Operation Verification
- [x] Run `npm run dev` and verify startup
- [x] Test screen transitions (Learning ↔ Settings)
- [x] Check CPU/memory usage is normal
- [x] Verify Q key exits properly
- [x] Confirm no infinite rendering occurs

## Phase 2: 各画面の実装（モックデータ使用）

## Step 8: LearningView with Mock Data
- [x] Create LearningView.tsx (no audio, text input only)
- [x] Implement slash command detection
- [x] Test with mock data
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 9: SlashCommandMenu Implementation
- [x] Create SlashCommandMenu.tsx
- [x] Implement dropdown display and arrow key navigation
- [x] Add real-time filtering
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 10: ResultView with Mock Data
- [x] Create ResultView.tsx (mock scoring results)
- [x] Implement keyboard shortcuts (Enter, N, R, S, Q)
- [x] Verify Japanese display
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 11: SettingsModal Implementation
- [x] Create SettingsModal.tsx
- [x] Implement arrow key navigation for settings
- [x] Connect to actual storage layer
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## 🔍 Human Quality Check 2: Full UI Operation Verification
- [x] Verify all screens display correctly (with mock data)
- [x] Test slash command menu operation
- [x] Test keyboard operations on each screen
- [x] Verify settings save and load
- [x] Confirm no memory leaks

## Phase 3: API統合と完全動作

## Step 12: State Management Full Implementation
- [x] Add all states to useStore.ts
- [x] Implement Round management, audio playback state
- [x] Create DebugPanel.tsx (when DEBUG=true)
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 13: Orchestrator Service
- [x] Implement orchestrator.ts
- [x] Integrate with OpenAI API
- [x] Integrate audio generation and playback
- [x] Implement pre-generation strategy
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 14: Full Integration
- [x] Complete learning flow with actual API
- [x] Implement error handling
- [x] Implement caching functionality
- [x] Run `npm run lint` and fix any issues
- [x] Run `npm run test` and ensure all tests pass

## Step 15: CLI & Performance
- [x] Implement CLI options
- [x] Performance optimization
- [x] Verify binary build
- [x] Run `npm run build`

## 🔍 Human Quality Check 3: Complete Operation Verification
- [x] Test actual learning flow with OpenAI API
- [x] Verify audio playback works
- [x] Test scoring and Japanese feedback
- [x] Test CLI options
- [x] Verify no memory leaks after extended use
- [x] Test global installation with `npm run build && npm link`

## Step 16: Documentation & Release Preparation
- [x] Create comprehensive README.md
- [x] Document API setup instructions
- [x] Add usage examples
- [x] Create CHANGELOG.md
- [x] Prepare package.json for npm publish
- [x] Test installation from npm (dry run)
- [x] Final code review and cleanup