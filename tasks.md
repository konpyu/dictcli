# task list

## NOTE

- Confirm that the Test and Linter have completed successfully at the end of each step
- Check the checkbox at the end of each step

## Step 1: Project Setup & Core Infrastructure
- [] Initialize npm project with TypeScript configuration
- [] Set up directory structure as defined in tech-design.md
- [] Install core dependencies (ink v4, zustand v4, commander.js, openai sdk)
- [] Configure TypeScript with strict mode
- [] Set up ESLint and Prettier configuration
- [] Create basic CLI entry point (cli.tsx)
- [] Implement environment variable loading (OPENAI_API_KEY)
- [] Set up vitest configuration for testing
- [] Run `npm run lint` and ensure no errors
- [] Run `npm run test` and ensure basic test setup works

## Step 2: Storage Layer Implementation
- [] Implement settings.ts for user preferences persistence
- [] Implement history.ts for learning history (JSONL format)
- [] Create default settings configuration
- [] Add unit tests for storage layer
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 3: OpenAI Service Layer
- [] Implement OpenAI client wrapper (client.ts)
- [] Create problem generator service (generator.ts) with CEFR levels
- [] Implement TTS service (tts.ts) with voice mapping
- [] Create scorer service (scorer.ts) with Japanese feedback
- [] Implement audio caching mechanism
- [] Add unit tests for OpenAI services with mocked responses
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Human Quality Check 1: OpenAI API Verification
- [] Create test script for Text API (test-openai-text.ts)
  - Generate sample sentences at different CEFR levels
  - Test scoring functionality with Japanese feedback
  - Verify response format and error handling
- [] Create test script for Audio API (test-openai-audio.ts)
  - Generate TTS for all voice options
  - Test different playback speeds
  - Verify audio file generation and caching
- [] Execute both test scripts and confirm API connectivity
- [] Document any API-related issues or limitations

## Step 4: Audio Player Service
- [] Implement audio player service using play-sound
- [] Add macOS afplay support
- [] Implement playback speed control
- [] Add error handling for audio playback
- [] Create unit tests for audio service
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 5: State Management with Zustand
- [] Implement useStore.ts with Ink v4 workaround
- [] Define all application states
- [] Implement manual subscription pattern
- [] Add state management utilities
- [] Create unit tests for store functionality
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 6: Core UI Components
- [] Create Header.tsx component
- [] Create StatusBar.tsx component
- [] Create DebugPanel.tsx component (for debug mode)
- [] Style components with Ink's Box and Text
- [] Add unit tests using ink-testing-library
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 7: Learning View Implementation
- [] Create LearningView.tsx component
- [] Implement audio playback integration
- [] Add text input handling
- [] Implement slash command detection
- [] Create SlashCommandMenu.tsx with dropdown functionality
- [] Add real-time filtering and arrow key navigation
- [] Write comprehensive tests for slash command functionality
- [] Test all keyboard interactions with ink-testing-library
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 8: Result View Implementation
- [] Create ResultView.tsx component
- [] Display scoring results with Japanese explanations
- [] Implement hybrid keyboard shortcuts (Enter, N, R, S, Q)
- [] Add WER calculation display
- [] Show error highlighting and alternatives
- [] Write tests for all keyboard shortcuts
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 9: Settings Modal Implementation
- [] Create SettingsModal.tsx component
- [] Implement arrow key navigation for settings
- [] Add voice, level, topic, and word count selection
- [] Integrate with storage layer
- [] Test all navigation and persistence
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 10: Orchestrator Service
- [] Implement orchestrator.ts for main flow control
- [] Add round management logic
- [] Implement pre-generation strategy
- [] Handle state transitions
- [] Create integration tests for complete flow
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 11: Main App Component
- [] Create App.tsx as main container
- [] Implement view transitions
- [] Add global keyboard shortcut handling
- [] Integrate all components
- [] Add debug mode support
- [] Test complete application flow
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Human Quality Check 2: TUI Mock Testing
- [] Create TUI mock application with all screens
  - Learning View with slash command menu
  - Result View with scoring display
  - Settings Modal with all options
  - Debug panel (when DEBUG=true)
- [] Test all screen transitions
  - Startup → Learning View
  - Learning View → Result View (on answer submission)
  - Any View → Settings Modal (via /settings or S key)
  - Result View → Learning View (on Enter/N key)
- [] Verify all keyboard interactions
  - Slash command menu navigation and selection
  - Result view shortcuts (Enter, N, R, S, Q)
  - Settings modal arrow key navigation
  - Escape key handling in modals
- [] Test edge cases
  - Empty input handling
  - Long text overflow
  - Rapid key press handling
  - Window resize behavior
- [] Document any UI/UX issues found

## Step 12: CLI Interface & Binary
- [] Implement command-line argument parsing
- [] Add --topic, --level, --words, --voice options
- [] Create npm binary configuration
- [] Test global installation
- [] Add help documentation
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 13: Performance Optimization
- [] Implement TTS caching with TTL
- [] Add cache size management
- [] Implement pre-generation for next round
- [] Add performance metrics tracking
- [] Test cache effectiveness
- [] Run `npm run lint` and fix any issues
- [] Run `npm run test` and ensure all tests pass

## Step 14: Final Integration Testing
- [] Run full E2E tests with ink-testing-library
- [] Test complete learning flow multiple times
- [] Verify data persistence across sessions
- [] Test error recovery scenarios
- [] Check memory usage and performance
- [] Run `npm run lint` - must pass with no errors
- [] Run `npm run test` - all tests must pass
- [] Run `npm run build` - must complete successfully

## Step 15: Documentation & Release Preparation
- [] Create comprehensive README.md
- [] Document API setup instructions
- [] Add usage examples
- [] Create CHANGELOG.md
- [] Prepare package.json for npm publish
- [] Test installation from npm (dry run)
- [] Final code review and cleanup