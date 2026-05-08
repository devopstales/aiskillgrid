# Hooks And Automation

This document defines hook locations, shared scripts, and automation policy across supported agent surfaces.

## Canonical Shared Hook Script

- Shared refresh hook: `.agents/hooks/refresh-indexes.sh`
- Purpose: trigger index refresh policy after merge/bootstrap shell flows
- Current behavior:
  - refresh `ccc` index
  - run GitNexus stale-check using `.gitnexus/meta.json` and `git rev-parse HEAD`
  - preserve embeddings mode when re-analyzing
  - fail-open with logs

## Surface Wiring

- Cursor:
  - `.cursor/hooks.json`
  - `.cursor/hooks/refresh-indexes.sh` (wrapper delegating to shared script)
- GitHub agents:
  - `.github/hooks/refresh-indexes.json`
- Kilo:
  - `.kilo/hook/hooks.md`
- OpenCode:
  - `.opencode/hook/hooks.md`

## Trigger Policy

Refresh automation is intended for shell commands matching:

- `git merge`
- `gh pr merge`
- `git pull`
- `ccc init`

## Operational Policy

1. Keep hook logic deduplicated under `.agents/hooks/`.
2. Surface-specific files should be thin wrappers/config only.
3. Hook failures should not block user workflow unless explicitly set to fail-closed.
4. Any hook behavior change must be reflected in this file and in `100-ide-configs.md`.

## Related Documents

- `01-installation.md`
- `02-workflow-usage.md`
- `11-memory-and-indexing.md`
- `100-ide-configs.md`
