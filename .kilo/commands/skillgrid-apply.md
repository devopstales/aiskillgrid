---
name: /skillgrid-apply
id: skillgrid-apply
category: Workflow
description: Implement from tasks; OpenSpec apply loop; keep PRD and tasks.md in sync
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Task
argument-hint: "[change-id; optional task focus]"
---

<objective>

You are executing **`/skillgrid-apply`** (BUILD phase) for the Skillgrid workflow.

Implement from the OpenSpec change’s **apply** instructions. **Always use hybrid persistence:** keep **PRD**, **`tasks.md`**, and **Engram** aligned — update checkboxes on disk and **`mem_save`** task progress to a stable `topic_key` (e.g. `skillgrid/{change}/tasks`) when you complete or restructure work.

**Status on exit:** Set the PRD’s **`Status:`** to **`inprogress`** once this run advances implementation (e.g. lands code, marks tasks, or checkboxes change). If the run only loads context and does not advance work, you may keep the prior value. Next lifecycle values: **`devdone`** after **`/skillgrid-review`**, **`done`** after **`/skillgrid-finish`** (see **`/skillgrid-init` → PRD / change `Status` lifecycle**).

</objective>

<process>

## Part A — OpenSpec apply loop (on-disk OpenSpec)

**Input**: Optionally the change name (e.g. `/skillgrid-apply add-auth`). If omitted, infer from context, or if only one active change exists, use it. If ambiguous, run `openspec list --json` and ask the user to choose.

1. **Announce the change** — State you are **using change: `<name>`** and how to override (e.g. pass a different name next time).

2. **Status**

   ```bash
   openspec status --change "<name>" --json
   ```

   Note `schemaName` and which artifact holds tasks (often `tasks`).

3. **Apply instructions**

   ```bash
   openspec instructions apply --change "<name>" --json
   ```

   This yields context file paths, progress (total / complete / remaining), task list, and a dynamic note for the current state.

   **Handle states:**

   - **`blocked`** (missing artifacts): show which artifacts are missing; complete them with **`/skillgrid-plan`** or **`/skillgrid-breakdown`**, or create files per `openspec instructions <artifact-id> --change "<name>" --json`. Do not guess past schema requirements.
   - **`all_done`**: congratulate; suggest **`/skillgrid-review`** and **`/skillgrid-finish`**.
   - Otherwise: proceed to implementation.

4. **Read context** — Read every path in `contextFiles` from the apply output (commonly proposal, specs, design, `tasks` for spec-driven flows). Prefer CLI output over hard-coded filenames. Also read **`.skillgrid/tasks/context_<name>.md`** if it exists (rolling handoff; see `docs/workflow.md` — *Filesystem handoff*).

5. **Filesystem handoff (subagents)** — If you use **`Task`** to delegate **research, design critique, or exploration** to a subagent for this change: create or open **`.skillgrid/tasks/context_<name>.md`** and **include its path in the subagent prompt**. After the subagent returns, **read** that file and any **`.skillgrid/tasks/research/<name>/`** files it cites **before** writing product code. The parent session keeps full implementation context; subagents spill long output to disk.

6. **Show progress** — Display schema, **N/M tasks** complete, remaining work, and the CLI’s current instruction.

7. **Implement (loop)** — For each **pending** task line:

   - State which task you are on.
   - Make minimal, focused code changes.
   - **Immediately** set the task checkbox in **`openspec/changes/<name>/tasks.md`**: `- [ ]` → `- [x]`.
   - **Immediately** mirror the same line in the **PRD** **Implementation tasks** section (**`.skillgrid/prd/PRD<NN>_<slug>.md`**). The two must stay identical.
   - Continue until done, blocked, or interrupted.

8. **Pause if** the task is unclear, implementation contradicts the design, an error is hit, or the user stops you—report and wait.

9. **Post-implementation housekeeping** — After a meaningful run (especially when code or the tree changed):

    - **Ask** (or self-check with the user): *Did this change affect the architecture or repo structure?*
    - **Checklist — update `.skillgrid/project/` when needed:** If you **added, renamed, or removed** top-level directories, **new services** or subsystems, **major patterns**, or **runtime topology**, update the right file(s) **in this pass** (do not defer without reason):
      - **`.skillgrid/project/STRUCTURE.md`** — folder and package layout; where code lives.
      - **`.skillgrid/project/ARCHITECTURE.md`** — design decisions, boundaries, new subsystems, pattern or topology shifts. When you add a **new ADR**, **cross-reference** it from here.
      - **`.skillgrid/project/PROJECT.md`** — new **dependencies**, **build** / **CI** config, and **tooling** the team must know about.
    - **Graph** — When the project uses **graphify**, run **`graphify update .`** to refresh the repo map (see also **Part B**). Typical triggers: this housekeeping step after a significant apply; **structural** code edits; **major** file additions/removals (e.g. new packages, deleted trees, large moves).

10. **End of session** — Summarize completed tasks, N/M progress, housekeeping actions (if any), and next action (continue apply, run review, or archive).

### Output shape (optional)

Use clear headings, e.g. `Implementing: <change> (schema: …)`, per-task progress, and a short completion or pause block.

### Apply guardrails

- Read context from the **apply** instruction output every time; do not assume fixed filenames beyond what the schema provides.
- Keep changes scoped; one task at a time.
- Update checkboxes in **`tasks.md` and the PRD** in the same edit pass when possible.
- On ambiguity or schema mismatch, stop and ask—do not invent missing artifacts.

### Fluid workflow

- Apply can be invoked with partial progress; you may need to **update design or specs** if implementation finds issues—suggest that explicitly before coding around the schema.

---

## Part B — Skillgrid-specific (hybrid bookkeeping)

1. **Worktrees** — Prefer **git worktrees** for parallel or isolated work when it helps the team.
2. **Task bookkeeping (mandatory)** — Whenever a task is **done** or **split/deferred/corrected**, update:

   - `openspec/changes/<change-id>/tasks.md` (when it exists)
   - The **PRD** **Implementation tasks** (same line numbers and checkbox state)
   - The **PRD** **`Status:`** line to **`inprogress`** when this session advances work (per **objective**)

   If there is no PRD file, update `tasks.md` (or Engram) only and add a pointer when the team expects a PRD.
3. **Scope** — Smallest vertical slice that satisfies the current tasks; no unrelated refactors.
4. **Contracts** — Preserve agreed APIs and error behavior at public boundaries; document contract changes in design/PRD when you must change them.
5. **TDD** — When tests are in play: red–green–refactor; write the failing test first when the task calls for it.
6. **Frameworks** — Ground behavior in the official docs for the stack the repo uses; cite in commits or code comments when useful.
7. **Quality** — Small, reviewable commits; work in **vertical slices** (implement, verify, commit).
8. **Migrations** — For schema/data changes, follow safe migration practices (one migration per logical change, rollback story).
9. **Graph & project docs** — At end of a substantive apply run, follow **Part A — step 9 (Post-implementation housekeeping)**: **`graphify update .`** when needed, and keep **`.skillgrid/project/ARCHITECTURE.md`**, **STRUCTURE.md**, and **PROJECT.md** aligned with what you changed (see step 9 for which file to touch). **`/skillgrid-finish`** runs a final consistency pass before archive.

## Notes

- Inspect the repo; do not assume stack or layout.
- Checklist format and PRD link live in **`/skillgrid-breakdown`**.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: change name, which **`openspec instructions apply`** steps completed, files changed in the repo, PRD **`Implementation tasks`** sync, and (if any) **`.skillgrid/project/`** or **`graphify`** updates from post-implementation housekeeping.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **`/skillgrid-test`** to run automated checks; if tests are not in scope, **`/skillgrid-review`**.

</process>
