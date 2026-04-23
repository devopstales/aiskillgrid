---
description: Security review, hardening, static analysis, dependency and config audit
---

You are executing **`/skillgrid-security`** (REVIEW — security) for the Skillgrid workflow.

Prefer running **`/skillgrid-review`** first so behavior matches specs; this phase focuses on **threats, misuse, and supply-chain** surfaces.

## Steps

1. **Code security review** — Apply `security-review` and `security-and-hardening` (auth, secrets, boundaries, OWASP-oriented patterns).
2. **Static analysis** — Run **Semgrep** and other project-configured SAST; follow `semgrep-security` for conventions.
3. **Artifacts and dependencies** — Run **Trivy** (and related scanners) on containers, IaC, and lockfiles per `trivy-security`.
4. **Risk framing** — Use `vulnerability-scanner` for threat modeling and prioritization when the change has meaningful exposure.
5. **Agent and IDE config** — When relevant, `security-scan` for `.claude/`, MCP, hooks, and similar automation surfaces.

## Skills to read and follow

- `.agents/skills/security-review/SKILL.md` — code-focused security review.
- `.agents/skills/security-and-hardening/SKILL.md` — OWASP-oriented patterns, auth, secrets, boundaries.
- `.agents/skills/semgrep-security/SKILL.md` — Semgrep static analysis.
- `.agents/skills/trivy-security/SKILL.md` — Trivy for containers/deps and related surfaces.
- `.agents/skills/vulnerability-scanner/SKILL.md` — threat modeling and prioritization.
- `.agents/skills/security-scan/SKILL.md` — audit agent/IDE config.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If scanners are not installed, say what is missing and what commands would run once configured.
