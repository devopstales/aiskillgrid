---
description: Security review, hardening, SAST, dependency/config audit, attack-surface and deprecation hygiene
allowed-tools: Read, Glob, Grep, Bash, Task
argument-hint: "[scope: path, surface, or change-id]"
---

<objective>

You are executing **`/skillgrid-security`** (REVIEW — security) for the Skillgrid workflow.

Prefer running **`/skillgrid-review`** first so behavior matches specs; this phase focuses on **threats, misuse, and supply-chain** surfaces.

</objective>

<process>

## Steps

1. **Code security review** — AuthN/Z, secrets, injection, SSRF, path traversal, deserialization, and trust boundaries; minimal blast radius.
2. **Static analysis** — Run **Semgrep** and other project-configured SAST.
3. **Artifacts and dependencies** — Run **Trivy** (and related scanners) on containers, IaC, and lockfiles as configured.
4. **Risk framing** — Threat model and prioritize when exposure is meaningful; document assumptions.
5. **Agent and IDE config** — When relevant, audit `.claude/`, MCP servers, hooks, and agent definitions for unsafe defaults.
6. **Deprecation / attack surface** — When removing or replacing endpoints, auth paths, or dependencies, retire old entry points and document timelines.

## Practices (inline)

- State explicit **threat assumptions** (who is trusted, what data is sensitive).
- Prefer measurable checks (scanner commands, policy-as-code) over generic advice.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If scanners are not installed, say what is missing and what commands would run once configured.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: scanners or checklists run, findings severity summary, and files written (e.g. reports) if any.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **`/skillgrid-validate`** for combined review+security sign-off, or **`/skillgrid-finish`** if you already ran review separately and are ready to ship.

</process>
