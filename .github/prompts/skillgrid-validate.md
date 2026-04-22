---
description: Verify against specs, review code, security, performance, docs
---

You are executing **`/skillgrid-validate`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-validate`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Run `openspec-verify-change` and/or `sdd-verify` when the project uses those artifacts.
2. Perform structured review (`code-review-and-quality`); simplify only with `code-simplification` when appropriate.
3. Apply security skills; run Trivy/Semgrep projects when configured.
4. Update ADRs and API docs where decisions changed (`documentation-and-adrs`).

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/openspec-verify-change/SKILL.md`
- `.agents/skills/sdd-verify/SKILL.md`
- `.agents/skills/karpathy-guidelines/SKILL.md`
- `.agents/skills/clean-code/SKILL.md`
- `.agents/skills/code-review-and-quality/SKILL.md`
- `.agents/skills/code-simplification/SKILL.md`
- `.agents/skills/security-review/SKILL.md`
- `.agents/skills/security-and-hardening/SKILL.md`
- `.agents/skills/performance-optimization/SKILL.md`
- `.agents/skills/semgrep-security/SKILL.md`
- `.agents/skills/trivy-security/SKILL.md`
- `.agents/skills/vulnerability-scanner/SKILL.md`
- `.agents/skills/security-scan/SKILL.md`
- `.agents/skills/database-reviewer/SKILL.md`
- `.agents/skills/documentation-and-adrs/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
