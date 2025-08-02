# Repository Guidelines

## Project Structure & Module Organization
- `src/cli.tsx`: CLI entry point rendered with Ink.
- `src/components/*`: React components (e.g., `App`, `LearningView`, `ResultView`, `SettingsModal`).
- `src/services/*`: App services (`openai/*` generator/scorer/tts, `audio/player`, `orchestrator`, `sceneLoader`).
- `src/storage/*`: Persistent data (`settings`, `history`). Settings saved to `~/.dictcli/settings.json`.
- `src/store/*`: Zustand store. `src/types/*`: shared types. `src/utils/*`: helpers.
- `tests/{unit,integration}`: Vitest test suites. `dist/`: compiled output.

## Build, Test, and Development Commands
- `npm run dev`: Start the CLI from sources.
- `npm run dev:tech` / `dev:advanced` / `dev:beginner`: Preset demo runs.
- `npm run build`: TypeScript compile to `dist/`.
- `npm test`: Run Vitest once. `npm run test:watch`: watch mode.
- `npm run lint`: ESLint over `src`. `npm run format`: Prettier format `src/**/*.{ts,tsx}`.
- After build: `node dist/cli.js` or, when linked, `dictcli`.
- Requirements: Node.js >= 20, macOS for audio (uses `afplay`).

## Coding Style & Naming Conventions
- Language: TypeScript (ESM). Strict compiler options enabled.
- Formatting: Prettier (2 spaces, single source of style). Linting: ESLint with `@typescript-eslint`, `react`, `prettier`.
- Rules: no unused vars (prefix ignored args with `_`), no `any`, JSX without importing React is allowed.
- Naming: PascalCase for components/types, camelCase for functions/variables, directories/files match existing patterns (e.g., `ResultView.tsx`, `services/openai/*`). Keep `.js` extensions in imports from TS output where present.

## Testing Guidelines
- Framework: Vitest (node environment). Coverage via V8: run `vitest --run --coverage` if needed.
- File locations: place tests under `tests/unit` or `tests/integration`; name files `*.test.ts[x]`.
- Expectations: keep tests fast, deterministic; mock OpenAI/audio where applicable.

## Commit & Pull Request Guidelines
- Commit style (observed): short, imperative messages (e.g., "fix lint error", "adjust speed setting"). Prefer one topic per commit.
- Include: brief rationale, scope if helpful, and related file paths in body when non-trivial.
- PRs should include: summary, before/after or CLI screenshots when UI changes, test plan (commands), and linked issues.
- Ensure `npm run lint`, `npm test`, and `npm run build` pass locally before requesting review.

## Security & Configuration Tips
- Set `OPENAI_API_KEY` in your shell. Use `DICTCLI_DEBUG=true` for verbose logs.
- Never commit keys or `~/.dictcli/*`. Audio features target macOS (`afplay`).
