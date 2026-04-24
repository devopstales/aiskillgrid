---
name: /skillgrid-session
id: skillgrid-session
category: Workflow
description: Phase 0 — session charter, context budget, MCP/tool selection, checkpoints
allowed-tools: Read, Glob, Grep, Write, Task
argument-hint: "[optional: session goal or resume context]"
---

<objective>

You are executing **`/skillgrid-session`** (Phase 0) for the Skillgrid workflow.

Use this at the **start of a new agent session**, after context compaction, or when switching tracks so work stays bounded and reproducible.

</objective>

<process>

## Steps

1. **Charter** — Capture a short session note: goal, constraints, success criteria, and stack or repo assumptions. Align with how your project records session context (e.g. project playbook or `NOTE.md`).
2. **Context budget** — Decide what must be loaded now versus later: rules, skills, MCP servers, and large files. Prefer minimal viable context until a phase needs depth.
3. **Tooling** — Enable only the MCPs and CLIs needed for this session; defer the rest to avoid noise and token burn.
4. **Checkpoint** — Note where to resume if interrupted (open change, branch, last task id, or spec section).
5. **Research vs build** — If discovery is still open, prefer a quick literature pass (web, docs MCPs, repo search) before locking design; if building, defer broad research unless a risk appears.

## Practices (inline)

- **Context** — Pack rules and files in layers; load deep docs only when a phase needs them.
- **Assumptions** — State what you believe about stack and scope; keep edits small and falsifiable.
- **Packaged instructions** — Optional: some IDEs auto-discover extra instruction files in the repo; this command does not link to any of them by path.

## Notes

- Phase 0 is **orthogonal** to `/skillgrid-init` (repo bootstrap). Run session when the **agent context** resets; run init when the **repository** still needs structure or tooling setup.
- Inspect the repo with tools; do not assume stack or layout.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: Engram or memory steps run, `topic_key` or summaries saved, context loaded, and any repo paths refreshed.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — The user’s substantive task: usually **`/skillgrid-plan`**, **`/skillgrid-init`**, or **`/skillgrid-apply`** depending on their goal; state which you recommend and why in one line.

</process>
