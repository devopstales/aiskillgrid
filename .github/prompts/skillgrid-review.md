---
description: Verify specs, code review, performance, data, docs (pre-security gate)
---

You are executing **`/skillgrid-review`** (REVIEW — quality and alignment) for the Skillgrid workflow.

Run **`/skillgrid-security`** after this when you want a **dedicated security pass**, or use **`/skillgrid-validate`** to run both in one invocation.

## Steps

1. **Spec / SDD verify** — Run `openspec-verify-change` and/or `sdd-verify` so implementation matches the change’s proposal, design, tasks, and delta specs.
2. **Code review** — Use `code-review-and-quality` (and `clean-code`, `code-simplification` only when simplification is safe and behavior-preserving).
3. **Assumptions** — Cross-check success criteria and scope with `karpathy-guidelines`.
4. **Performance** — When applicable, `performance-optimization`: measure first, then Web Vitals, profiling, bundle checks.
5. **Data** — For SQL/storage, `database-reviewer` (schema, RLS, performance) as appropriate.
6. **Documentation** — Update ADRs, API docs, and inline *why* via `documentation-and-adrs` when behavior or contracts changed.

## Skills to read and follow

- `.agents/skills/openspec-verify-change/SKILL.md` — implementation matches change artifacts.
- `.agents/skills/sdd-verify/SKILL.md` — SDD verification against specs, design, and tasks.
- `.agents/skills/karpathy-guidelines/SKILL.md` — success criteria and scope discipline.
- `.agents/skills/clean-code/SKILL.md` — review for clarity and coupling.
- `.agents/skills/code-review-and-quality/SKILL.md` — multi-axis review, change sizing, severity labels.
- `.agents/skills/code-simplification/SKILL.md` — simplify without changing behavior (Chesterton’s fence, rule of 500).
- `.agents/skills/performance-optimization/SKILL.md` — measure first, Web Vitals, profiling, bundle checks.
- `.agents/skills/database-reviewer/SKILL.md` — PostgreSQL-oriented review when applicable.
- `.agents/skills/documentation-and-adrs/SKILL.md` — ADRs, API docs, document the why.

## Optional: IDE personas

- **`skillgrid-spec-verifier`** ([`.cursor/agents/skillgrid-spec-verifier.md`](../../.cursor/agents/skillgrid-spec-verifier.md)) — traceability only: specs ↔ `tasks.md` ↔ implementation; use when you want that split from general code review.
- **`skillgrid-code-reviewer`** ([`.cursor/agents/skillgrid-code-reviewer.md`](../../.cursor/agents/skillgrid-code-reviewer.md)) — five-axis code review in a dedicated subagent context.

With a **parallel** harness, you may run independent personas on the same change and merge reports (see [`.cursor/agents/README.md`](../../.cursor/agents/README.md)).

## Notes

- Do **not** treat this phase as sufficient for security sign-off; follow with `/skillgrid-security` before merge when security matters.
- Inspect the repo with tools; do not assume stack or layout.
