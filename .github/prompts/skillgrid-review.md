---
name: /skillgrid-review
id: skillgrid-review
category: Workflow
description: Verify change artifacts, code review, performance, data, docs (pre-security gate)
allowed-tools: Read, Glob, Grep, Bash, Task
argument-hint: "[change-id, PR, or scope]"
---

<objective>

You are executing **`/skillgrid-review`** (REVIEW — quality and alignment) for the Skillgrid workflow.

Run **`/skillgrid-security`** after this for a **dedicated security pass**, or **`/skillgrid-validate`** to run both.

This command combines **OpenSpec-style verification** on disk with **Engram** retrieval and reporting. **Always use hybrid persistence:** read artifacts from **`openspec/changes/<name>/`**, cross-check **Engram** with `mem_search` → `mem_get_observation` for full text, and **persist** the verification report to **`openspec/changes/<name>/verify-report.md`** **and** **`mem_save`** with `topic_key` **`skillgrid/<change-name>/verify-report`** (or your team’s prefix).

**Status on exit:** Set the change’s PRD **`Status:`** to **`devdone`** (and INDEX / ticket table if used) when this review run completes. Final step is **`done`** in **`/skillgrid-finish`** (see **`/skillgrid-init` → PRD / change `Status` lifecycle**).

</objective>

<process>

## A — On-disk OpenSpec verify

**Input**: Optionally a change name. If omitted, use context or `openspec list --json` and **ask the user to choose**; do not guess when multiple changes exist.

1. **Status**

   ```bash
   openspec status --change "<name>" --json
   ```

   Note `schemaName` and which artifacts exist.

2. **Apply context (artifact paths)**

   ```bash
   openspec instructions apply --change "<name>" --json
   ```

   Read all `contextFiles`.

3. **Report dimensions** — Track **completeness**, **correctness**, and **coherence**; tag issues as CRITICAL, WARNING, or SUGGESTION.

4. **Completeness**

   - **Tasks:** If `tasks.md` is in context, count `- [ ]` vs `- [x]`. Incomplete core tasks → **CRITICAL** with a concrete “complete or mark N” recommendation.
   - **Delta specs:** If `openspec/changes/<name>/specs/` exists, extract **### Requirement:** blocks and check the codebase for likely implementation; unimplemented requirements → **CRITICAL** when clear.

5. **Correctness**

   - For each requirement, search for implementation evidence; note file paths. Divergence from intent → **WARNING** with file:line and requirement id.
   - For **#### Scenario:** blocks, check handling and tests; missing coverage → **WARNING** or **SUGGESTION** by severity.

6. **Coherence**

   - If `design.md` is present, check major decisions and “Approach/Architecture” style sections against the code; contradictions → **WARNING**.
   - Check new code against repo patterns (naming, layout); big deviations → **SUGGESTION** unless they break contracts.

7. **Scorecard and verdict**

   - Summarize in a short table: completeness / correctness / coherence.
   - List CRITICAL, then WARNING, then SUGGESTION with **actionable** follow-ups and `file:line` where possible.
   - Verdict: ready for merge vs must-fix vs ship with known warnings.
   - **Graceful degradation:** If only `tasks.md` exists, focus on task completion; if no delta specs, skip requirement-level correctness from specs.

8. **Heuristics** — Prefer SUGGESTION over WARNING, WARNING over CRITICAL, when evidence is weak. Be specific, not vague.

9. **Report file (required)** — Write **`openspec/changes/<name>/verify-report.md`** with the same content as the user-facing summary. **Also** `mem_save` that report to Engram (section B).

10. **PRD `Status`** — Set the linked PRD’s **`Status:`** to **`devdone`** (and INDEX / ticket table if used). Align with **`/skillgrid-init` → PRD / change `Status` lifecycle**.

### Verification: plain-language reference

- **Completeness** — Objectives: checklists and requirements lists.
- **Correctness** — Keyword search, structure, and reasonable inference; not legal proof.
- **Coherence** — Obvious design/code mismatches, not style nitpicks unless the repo demands it.

---

## B — Engram (required with on-disk)

**Previews are truncated** — for anything you rely on from Engram, `mem_get_observation` the full text.

1. **Search** — `mem_search` for proposal, spec, design, and tasks for this change (scope **`project`** as your team defines).

2. **Retrieve** — `mem_get_observation(id: …)` for **full** content for each id used in the review. **Do not** reason from `mem_search` snippets alone.

3. **Reconcile** — Compare Engram copy with **`openspec/changes/<name>/`** files; note any drift. Run the same **completeness / correctness / coherence** checks against **both** sources and the repository.

4. **Persist** — `mem_save` the final markdown report; **`topic_key`:** `skillgrid/<change-name>/verify-report` (align with **`/skillgrid-init`**). The body must match **`openspec/changes/<name>/verify-report.md`**.

5. **Execution evidence** (recommended) — Run tests and build when the repo supports it; attach results in both the file and the `mem_save` body (see `package.json`, `Makefile`, or `openspec/config.yaml`).

**If Engram is unavailable in-session** — complete section A, write `verify-report.md` on disk, and state that Engram `mem_save` is pending.

---

## C — General code and product review

1. **Code review** — Readability, coupling, error handling, and change size; use severity labels (blocker / major / minor) if helpful.
2. **Assumptions** — Re-read PRD or proposal **success criteria** and **scope**; flag drift.
3. **Performance** — If relevant: measure or profile before large optimizations; check budgets and hot paths the change touched.
4. **Data** — For SQL/storage: review schema, queries, and safety (indexes, RLS, migrations) as appropriate.
5. **Documentation** — If behavior or public APIs changed, update ADRs, README, or user-facing docs in the same effort or list it as a follow-up.

## Optional: IDE personas

- **`skillgrid-spec-verifier`** ([`.cursor/agents/skillgrid-spec-verifier.md`](../../.cursor/agents/skillgrid-spec-verifier.md)) — traceability: specs ↔ `tasks.md` ↔ implementation.
- **`skillgrid-code-reviewer`** ([`.cursor/agents/skillgrid-code-reviewer.md`](../../.cursor/agents/skillgrid-code-reviewer.md)) — deeper multi-axis code review in a subagent.

See [`.cursor/agents/README.md`](../../.cursor/agents/README.md) for running personas in parallel.

## Notes

- This phase is **not** a full security sign-off; use **`/skillgrid-security`** when it matters.
- Inspect the repo with tools; do not assume stack or layout.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: change id, `verify-report` or review outcome, spec vs code gaps, and Engram save if any.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **`/skillgrid-security`** for threat/scanner pass; or **`/skillgrid-validate`** to combine gates if your workflow uses a single step.

</process>
