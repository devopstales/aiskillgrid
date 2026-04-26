---
description: Close the change: optional spec sync, archive, git hygiene, ship checklist, PR
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[change-id or branch name]"
---

<objective>

You are executing **`/skillgrid-finish`** (SHIP phase) for the Skillgrid workflow.

Order of operations: **(1) optional** sync of delta specs into main `openspec/specs/`, **(2)** final alignment of **`.skillgrid/project/`** (ARCHITECTURE, STRUCTURE, PROJECT) with the **merged** codebase, **(3)** archive the change on disk, **(4)** **`mem_save`** ship/closure to Engram (e.g. `topic_key` `skillgrid/<change>/archive`), **(5)** open/update PR, CI, and ship notes. **Always hybrid:** never treat disk and Engram as optional alternatives.

**Status on exit:** Set the relevant PRD **`Status:`** to **`done`** (and INDEX / ticket table if used) when this phase completes successfully — last value in the lifecycle under **`/skillgrid-init` → PRD / change `Status` lifecycle** (`draft` → `todo` → `inprogress` → `devdone` → `done`).

</objective>

<process>

## Flow

```mermaid
flowchart TD
    START([User calls /skillgrid-finish])
    SYNC{Delta specs exist?}
    START --> SYNC
    SYNC -->|Yes| MERGE[Merge delta specs into main openspec/specs/]
    SYNC -->|No| ALIGN[Align .skillgrid/project/ docs]
    MERGE --> ALIGN
    ALIGN --> ARCHIVE[openspec status, warn if incomplete]
    ARCHIVE --> MOVE[mv openspec/changes/{name} archive/YYYY-MM-DD-{name}]
    MOVE --> ENGRAM[mem_save archive record]
    ENGRAM --> PREVIEW[Cleanup .skillgrid/preview/ if user agrees]
    PREVIEW --> PR[Create/update PR, CI checks]
    PR --> DONE[Set PRD status to done]
    DONE --> END([Merge / deploy, suggest new cycle])
```

## 1 — Optional: sync delta specs to main (before or without archive)

**When:** The change has **delta** specs under `openspec/changes/<name>/specs/` and the team wants **`openspec/specs/<capability>/spec.md`** updated before or without archiving the change.

**Input**: Change name, or ask via `openspec list --json` if ambiguous. **Do not** auto-pick a change when multiple are active and the user did not name one.

**Steps**

1. List changes that have a `specs/` subtree under the change. If none, skip this section.

2. For each `openspec/changes/<name>/specs/<capability>/spec.md`:

   - Read the **delta** (sections like `## ADDED Requirements`, `## MODIFIED Requirements`, `## REMOVED Requirements`, `## RENAMED Requirements`).
   - Read or create the **main** file `openspec/specs/<capability>/spec.md`.
   - **Merge intelligently** — apply partial updates (e.g. add a scenario under MODIFIED without pasting the whole file). Preserve main-spec content not touched by the delta.
   - **ADDED** — add or reconcile requirements; **MODIFIED** — apply edits; **REMOVED** — remove blocks; **RENAMED** — apply FROM/TO.
   - If the capability is new, create the directory and a minimal **Purpose** plus the added requirements.

3. Summarize: which capabilities were updated and what class of change (add/modify/remove/rename).

**Idempotency** — Re-reading delta and main and re-applying should be safe; when unsure, ask.

---

## 2 — Final alignment: `.skillgrid/project/` (before archive)

**When:** The repo uses **`.skillgrid/project/`** (from **`/skillgrid-init`** or equivalent). Run this **after** the implementation branch reflects the work you are closing, and **before** moving the change to `openspec/changes/archive/`.

**Goal:** **ARCHITECTURE.md**, **STRUCTURE.md**, and **PROJECT.md** match the **current** tree, dependencies, and how the system runs—so the merged `main` is not ahead of the onboarding docs.

**Steps**

1. **Skim the three files** and compare to the repo: layout (**STRUCTURE**), decisions and subsystems (**ARCHITECTURE**), deps and tooling (**PROJECT**). If **`/skillgrid-apply`** already updated them during the change, still verify nothing drifted in the last commits.
2. **Update** any section that is wrong or missing after the change (new top-level services, package moves, new external deps, build/CI or runtime topology, **ADRs** — link new ADRs from **ARCHITECTURE.md** as appropriate).
3. **Commit** — Prefer a **dedicated** commit for doc-only fixes (e.g. `docs: align Skillgrid project files with <change-id>`) when the diff is non-trivial, so code review can separate product changes from doc catch-up. Small fixes may fold into the final ship commit if the team prefers.

**If the project has no `.skillgrid/project/`** — Skip; optionally note in the finish summary.

---

## 3 — Archive the OpenSpec change (on-disk)

**Input**: Change name, or ask via `openspec list --json` if missing/ambiguous. Prefer explicit user choice when multiple changes exist.

**Steps**

1. **Status**

   ```bash
   openspec status --change "<name>" --json
   ```

   List artifacts not `done`. Warn and get confirmation to continue if something is incomplete.

2. **Tasks** — Open `openspec/changes/<name>/tasks.md` (or the path your schema uses). Count `- [ ]` vs `- [x]`. If incomplete, warn and confirm before archiving.

3. **Sync prompt when delta specs exist** — If `openspec/changes/<name>/specs/` exists, compare to main `openspec/specs/`. If main is behind, offer:

   - **Sync now (recommended)** — run the **sync procedure** in **section 1**, then continue (you may re-check **section 2** if specs implied structural edits).
   - **Archive without syncing** — document that choice in the finish summary.
   - **Cancel** — stop.

4. **Archive move**

   ```bash
   mkdir -p openspec/changes/archive
   mv openspec/changes/<name> openspec/changes/archive/YYYY-MM-DD-<name>
   ```

   Use **today’s date** in `YYYY-MM-DD` (authoritative: session “today”). If `openspec/changes/archive/YYYY-MM-DD-<name>` **already exists**, do not overwrite: report the conflict and list options (rename old archive, remove duplicate, or use another date).

5. **Summary** — Print change id, schema if known, final archive path, whether specs were synced, and any warnings (incomplete artifacts or tasks, sync skipped).

6. **Engram** — `mem_save` a short closure record: archive path, PR link or branch, spec sync outcome. **`topic_key`:** e.g. `skillgrid/<name>/archive`. If you could not archive on disk (e.g. missing `openspec/`), still save to Engram and explain the gap.

---

## 3a — Cleanup preview

Check **`.skillgrid/preview/`** for files related to the change you’re closing.

- If the directory is empty or gitignored and contains only stale files, list them.
- Ask the user: *“Can I delete these preview files? [list]”*  
  On confirmation, remove them.
- If the team keeps previews for reference, skip with a brief note.
   ```

---

## 4 — Pull request, CI, and ship

1. **Pull request** — Open or update a PR: clear description, risks, and links to **`openspec/changes/…`**, the numbered PRD file(s) (**`.skillgrid/prd/PRD<NN>_<slug>.md`**), and tests. Keep commits small and reviewable; follow team branching rules.
2. **CI / gates** — Ensure required checks pass; feature flags and deployment hooks as the project does.
3. **Deprecation** — If old paths are retired, document timelines and follow-up issues.
4. **Deploy** — When shipping to production, use your **launch checklist** (rollout, monitoring, rollback) as the team defines.
5. **Documentation** — ADRs, API docs, and “why” comments if behavior or contracts changed.

6. **PRD `Status`** — Set the relevant PRD’s **`Status:`** to **`done`** (and INDEX / ticket table if used). This is the terminal state in **`/skillgrid-init` → PRD / change `Status` lifecycle**.

## Notes

- Inspect the repo; do not assume stack or layout.
- First-time OpenSpec **guided** walkthrough: see [`/opsx-onboard`](./opsx-onboard.md) for a full cycle narrative; finish here is the **archive + ship** end of that story.

## Anti-patterns

- **Archiving without syncing specs** – Never move a change to archive if delta specs haven’t been merged into `openspec/specs/` (unless the user explicitly chooses not to).
- **Stale preview files** – Don’t leave `.skillgrid/preview/` files that were only used for A/B comparison; ask before archiving.
- **Forgetting the final status** – Never close the change without setting the PRD `Status:` to `done` and updating the INDEX/ticket table.
- **Merge without documentation alignment** – If the change added new dependencies, services, or patterns, `.skillgrid/project/` files must be updated before the PR is considered ready.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: spec sync (if any), **`.skillgrid/project/`** final alignment (if applicable), archive moves, PR/CI prep, and Engram or closure steps completed.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — None for this change cycle—**merge / deploy**; for **new** work, start with **`/skillgrid-plan`** or **`/skillgrid-explore`**.

</process>
