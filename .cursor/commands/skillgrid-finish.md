---
name: /skillgrid-finish
id: skillgrid-finish
category: Workflow
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
    DONE --> CHECKPOINTS[Clean checkpoints for closed change]
    CHECKPOINTS --> END([Merge / deploy, suggest new cycle])
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

---

## 3b — Remote tracker hygiene (read `.skillgrid/config.json`)

After the PR/merge path is clear, read **`.skillgrid/config.json`** and **`ticketing.provider`** (default mentally to **`local`** if missing). Suggest the matching closure habit:

| Provider | Suggestion |
|----------|------------|
| **`local`** | Set PRD **`Status:`** to **`done`** and refresh **`.skillgrid/prd/INDEX.md`**. If the team uses the PRD Kanban script, note that a browser refresh may be needed. |
| **`github`** | In the PR or merge commit body, use **`Fixes #123`** / **`Closes #123`** to close linked issues when applicable. |
| **`gitlab`** | Use **`Closes …`** / issue references per [GitLab cross‑linking](https://docs.gitlab.com/ee/user/project/issues/managing_issues.html#closing-issues-automatically) and team practice. |
| **`jira`** | Transition or close the linked Jira issue per workflow; paste the Jira key in the ship notes if this change fulfilled it. |

Do **not** require API access: if the user works in the web UI, list the **exact** follow‑up (issue keys, states) in the finish summary.

---

## 3c — Cleanup checkpoints

When the change is **finished** (merged locally, PR created and ready, or intentionally discarded), clean up change-scoped checkpoint entries from **`.skillgrid/tasks/checkpoints.log`**.

**Default behavior**

- Keep checkpoints for **unrelated** changes.
- Remove entries whose checkpoint name or context clearly belongs to the finishing change, for example:
  - `before-apply-<name>`
  - `before-apply-<name>-...`
  - entries whose `contexts=` field includes `.skillgrid/tasks/context_<name>.md`
- If the log becomes empty, leave an empty file or delete it—follow the repo’s existing style; either is acceptable.
- Record the cleanup in the finish summary.

**Do not clean checkpoints** when the user chooses **Leave as-is**; those checkpoints are still useful for resume and comparison.

**Ambiguity rule:** If a checkpoint might belong to multiple changes, keep it and mention that you left it untouched.

---

## 4 — Finish options (choose one)

After archiving, present the user with exactly these four choices and act accordingly:

1. **Merge locally and push**
   - Merge the change branch into `main` (or the default branch) locally: `git checkout main && git merge <branch>`
   - Push directly (if the team allows): `git push origin main`
   - Run any required CI checks that can be triggered locally; ensure all gates pass before pushing.
   - Clean change-scoped checkpoints per **section 3c**.
   - Proceed to **post‑merge steps** (deprecation, deploy, documentation, PRD status → `done`).

2. **Create a pull request**
   - Ensure the branch is pushed.
   - Open a PR with a clear description, links to the `openspec/changes/…` directory, the PRD file(s), and a summary of testing.
   - Trigger CI; wait for checks to pass (or note expected failures).
   - Once the PR is approved and merged by the team, update the PRD status to `done` (this may happen outside the session).
   - If you can’t merge now, leave the change in `devdone` status and note the PR URL.
   - Clean change-scoped checkpoints per **section 3c** once the PR is ready or merged according to team practice.

3. **Leave as‑is (do not merge yet)**
   - Keep the feature branch intact (and any optional local checkout the team uses); no further action.
   - The PRD status stays as `devdone` (or whatever it was).
   - Suggest the user return later with `/skillgrid-finish` to finalize.

4. **Discard the change**
   - Confirm with the user that they want to permanently discard this change.
   - **If** the team uses an **optional** extra git worktree (e.g. `.worktree/<slug>/` — not required by Skillgrid): remove it first: `git worktree remove .worktree/<slug>/ --force`
   - Delete the branch: `git branch -D <branch-name>`
   - Archive the `openspec/changes/<name>/` directory (already done in step 3) but set the PRD status to `archived` (or `done` with a note of discard). Optionally move the PRD to an archived section of the index.
   - Clean up any remaining preview files (step 3a already asked, but confirm again).
   - Clean change-scoped checkpoints per **section 3c**.

### Post‑merge actions (only after options 1 or 2 complete)

- **Deprecation** — If old paths are retired, document timelines and follow‑up issues.
- **Deploy** — Use the team’s launch checklist (rollout, monitoring, rollback).
- **Documentation** — Update ADRs, API docs, and “why” comments if behavior or contracts changed.
- **PRD status** — Set the relevant PRD’s **`Status:`** to **`done`** (and INDEX / ticket table). This is the terminal state.

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

1. **What I did** — Bullets: spec sync (if any), **`.skillgrid/project/`** final alignment (if applicable), archive moves, PR/CI prep, checkpoint cleanup, and Engram or closure steps completed.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — None for this change cycle—**merge / deploy**; for **new** work, start with **`/skillgrid-plan`** or **`/skillgrid-explore`**.

</process>
