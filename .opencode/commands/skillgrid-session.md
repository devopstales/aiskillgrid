---
name: /skillgrid-session
id: skillgrid-session
category: Workflow
description: Phase 0 — session charter, context budget, MCP/tool selection, checkpoints
---

You are executing **`/skillgrid-session`** (Phase 0) for the Skillgrid workflow.

Use this at the **start of a new agent session**, after context compaction, or when switching tracks so work stays bounded and reproducible.

## Steps

1. **Charter** — Capture a short session note: goal, constraints, success criteria, and stack or repo assumptions. Align with how your project records session context (e.g. project playbook or `NOTE.md`).
2. **Context budget** — Decide what must be loaded now versus later: rules, skills, MCP servers, and large files. Prefer minimal viable context until a phase needs depth.
3. **Tooling** — Enable only the MCPs and CLIs needed for this session; defer the rest to avoid noise and token burn.
4. **Checkpoint** — Note where to resume if interrupted (open change, branch, last task id, or spec section).
5. **Research vs build** — If discovery is still open, prefer `search-first` / `deep-research` before locking design; if building, defer broad research unless a risk appears.

## Skills to read and follow

- `.agents/skills/context-engineering/SKILL.md` — rules, context packing, MCP usage.
- `.agents/skills/using-agent-skills/SKILL.md` — meta: how to use the agent-skills pack.
- `.agents/skills/karpathy-guidelines/SKILL.md` — assumptions, simplicity, verifiable steps.
- `.agents/skills/search-first/SKILL.md` — when to research before building.
- `.agents/skills/continuous-learning/SKILL.md` — optional: refresh patterns or docs when the task needs it.

## Notes

- Phase 0 is **orthogonal** to `/skillgrid-init` (repo bootstrap). Run session when the **agent context** resets; run init when the **repository** still needs structure or tooling setup.
- Inspect the repo with tools; do not assume stack or layout.
