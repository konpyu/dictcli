# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-08-02

### 🎉 Initial Release

#### Features
- **LLM-Powered Dictation**: Dynamic sentence generation using OpenAI API
- **Multi-Level Support**: CEFR levels from A1 to C2
- **Text-to-Speech**: 6 voice options (3 male, 3 female)
- **Japanese Feedback**: Detailed error explanations in Japanese
- **Beautiful TUI**: Built with Ink v4 for smooth terminal experience
- **Slash Commands**: Quick access to replay, settings, hints, and quit
- **Customizable Settings**: Topic, level, word count, and voice selection
- **Progress Tracking**: Local history storage in JSONL format
- **Audio Speed Control**: Playback speed adjustment (0.8x-1.2x)
- **Pre-generation**: Next problem generated in background for faster rounds

#### Technical Details
- Built with TypeScript and strict mode enabled
- Zustand v4 for state management with Ink compatibility workaround
- Comprehensive test suite with Vitest
- ESLint and Prettier for code quality
- macOS audio support via afplay

#### Known Limitations
- Audio playback only works on macOS
- Requires OpenAI API key
- No offline mode

### Contributors
- Initial implementation by [konpyu]

---

For more details, see the [README](README.md).
