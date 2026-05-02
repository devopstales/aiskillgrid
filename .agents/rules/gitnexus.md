---
description: GitNexus code graph context
alwaysApply: true
---

This project uses GitNexus for local code graph context.

- Before answering architecture or codebase questions, prefer GitNexus MCP/resources for repository context, clusters, processes, impact analysis, and symbol context.
- If `.gitnexus/` exists, treat it as local generated index state and avoid broad raw file reads before checking graph context.
- After substantive code or layout changes, run `npx -y gitnexus@1.3.11 analyze` to keep the graph current. Use `npx -y gitnexus@1.3.11 analyze --force` only when a full rebuild is needed.
