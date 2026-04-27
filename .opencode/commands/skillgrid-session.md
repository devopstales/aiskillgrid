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

## Flow

```mermaid
flowchart TD
    START([User calls /skillgrid-session])
    CHARTER[Capture session goal & constraints]
    AGENTS[Read AGENTS.md and openspec/config.yaml]
    OPS[openspec list --json]
    PRDS[Glob .skillgrid/prd/PRD*.md, check statuses]
    CONTEXT[Check .skillgrid/tasks/context_*.md for active changes]
    ENGRAM[mem_search for skillgrid-init and change keys]
    DIRTY[git status --short]
    BUDGET[Decide what to load now vs later]
    TOOLS[Enable needed MCPs/CLIs]
    CHECKPOINT[Note resume marker]
    START --> CHARTER
    CHARTER --> AGENTS --> OPS --> PRDS --> CONTEXT --> ENGRAM --> DIRTY
    DIRTY --> BUDGET --> TOOLS --> CHECKPOINT --> DONE([Handoff to substantive phase])
```

## Steps

1. **Charter** — Capture a short session note: goal, constraints, success criteria, and stack or repo assumptions. Align with how your project records session context (e.g. project playbook or `NOTE.md`).
2. **Restore project memory (hybrid)**  
   - Read **`AGENTS.md`** and **`openspec/config.yaml`** (if present) for stack, conventions, and rules.  
   - Read **`.skillgrid/config.json`** (if present) and note **`ticketing.provider`** and **`artifactStore.mode`** (`hybrid` \| `openspec` \| `engram`, per **`/skillgrid-init`** / SDD **`sdd-init`** contract) — do not assume GitHub when the provider is not `github`; do not assume `openspec/` exists when mode is **`engram`**. If the file is missing or **`artifactStore`** is absent, mention **`/skillgrid-init`** to set defaults (**`hybrid`** for artifact store).  
   - Run `openspec list --json` to identify active changes.  
   - Glob **`.skillgrid/prd/PRD*.md`** and check `Status:` lines to find work in flight.  
   - For each active change, check for **`.skillgrid/tasks/context_<change-id>.md`** — if one exists, display its `state`, `current goal`, and `last checkpoint` so the user can resume immediately.  
   - If Engram is available, `mem_search` for `skillgrid-init/{project-name}` and any change-scoped keys (`skillgrid/<change>/plan`, `skillgrid/<change>/verify-report`). Offer to load the full text via `mem_get_observation`.

3. **Detect uncommitted and parallel work**  
   - Run `git status --short` to surface uncommitted work from a previous session.
   - **Optional (only if your team uses extra checkouts):** if `.worktree/` exists or `git worktree list` shows more than the main worktree, note it briefly. Skillgrid’s default is a **single working tree** with per-change handoff (see `docs/workflow.md` — *Filesystem handoff*); do not require worktrees to resume a change.
4. **Context budget** — Decide what must be loaded now versus later: rules, skills, MCP servers, and large files. Prefer minimal viable context until a phase needs depth.
5. **Tooling** — Enable only the MCPs and CLIs needed for this session; defer the rest to avoid noise and token burn.
6. **Checkpoint** — Note where to resume if interrupted (open change, branch, last task id, or spec section).
7. **Research vs build** — If discovery is still open, prefer a quick literature pass (web, docs MCPs, repo search) before locking design; if building, defer broad research unless a risk appears.

## Practices (inline)

- **Context** — Pack rules and files in layers; load deep docs only when a phase needs them.
- **Assumptions** — State what you believe about stack and scope; keep edits small and falsifiable.
- **Packaged instructions** — Optional: some IDEs auto-discover extra instruction files in the repo; this command does not link to any of them by path.

## Notes

- Phase 0 is **orthogonal** to `/skillgrid-init` (repo bootstrap). Run session when the **agent context** resets; run init when the **repository** still needs structure or tooling setup.
- Inspect the repo with tools; do not assume stack or layout.

## Anti-patterns

- **Relying on chat history** – Never assume context from a previous chat; always run `/skillgrid-session` first to restore active changes, per-change handoff files, and Engram state.
- **Loading everything** – Don’t enable all MCPs and files at once; use the context budget to load only what the imminent phase needs.
- **Skipping the checkpoint** – Never start heavy work without noting a resume marker (change id, task number, branch, etc.).
- **Ignoring dirty state** – Don’t forget uncommitted changes and `context_*.md` handoff before modifying files; optional worktrees are secondary.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: Engram or memory steps run, `topic_key` or summaries saved, context loaded, **`ticketing.provider`** and **`artifactStore.mode`** from **`.skillgrid/config.json`** (if read), and any repo paths refreshed.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — The user’s substantive task: usually **`/skillgrid-plan`**, **`/skillgrid-init`**, or **`/skillgrid-apply`** depending on their goal; state which you recommend and why in one line.

</process>
