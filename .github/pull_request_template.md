## Summary
- What and why of this change (1–2 sentences).

## Changes
- Bullet the key changes (files, modules, behaviors).
- Call out user-facing or CLI UX updates.

## Screenshots / Recordings (optional)
- CLI output, GIF, or before/after snippets if UI/UX changed.

## Test Plan
- Local commands run and results:
  - `npm run lint` → 
  - `npm test` → 
  - `npm run build` → 
  - Manual run: `npm run dev` or `node dist/cli.js`
- Add/modify tests under `tests/{unit,integration}` and note coverage if relevant.

## Checklist
- [ ] Lint passes (`npm run lint`)
- [ ] Tests pass (`npm test`); added/updated tests where appropriate
- [ ] Build succeeds (`npm run build`)
- [ ] Docs updated (README/AGENTS.md) if behavior or commands changed
- [ ] Backward-compatible (no breaking changes) or noted below
- [ ] No secrets committed (e.g., `OPENAI_API_KEY`); config remains in `~/.dictcli/settings.json`
- [ ] Audio behavior verified or safely stubbed (macOS `afplay`)

## Linked Issues
- Closes #

## Breaking Changes (if any)
- Impact, migration steps, and communication notes.

## Notes for Reviewers
- Areas needing extra attention, trade-offs, or follow-ups.

## Release Notes (optional)
- 1–2 lines for CHANGELOG.

