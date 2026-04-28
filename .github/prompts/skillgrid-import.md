---
description: Import existing PRDs and OpenSpec changes into Skillgrid PRD/index artifacts
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: path, change-id, --dry-run]"
---

<objective>

You are executing **`/skillgrid-import`** for the Skillgrid workflow.

Import existing PRDs and/or OpenSpec changes into the canonical `.skillgrid/prd/` workflow without deleting legacy sources or silently merging ambiguous artifacts.

**Status on exit:** inferred per imported PRD by `skillgrid-import-artifacts` and `skillgrid-prd-artifacts`.

</objective>

<process>

## Required Skills

Load these first for this command:

- `.agents/skills/skillgrid-import-artifacts/SKILL.md`
- `.agents/skills/skillgrid-prd-artifacts/SKILL.md`
- `.agents/skills/skillgrid-spec-artifacts/SKILL.md`
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md`
- `.agents/skills/skillgrid-hybrid-persistence/SKILL.md`
- `.agents/skills/skillgrid-questioning/SKILL.md`
- `.agents/skills/ccc/SKILL.md` — optional semantic **`ccc search`** when locating scattered PRDs or legacy docs in huge trees
- `.agents/skills/references/indexing-and-memory.md` — optional **graphify** / **ccc** ordering when mapping the repo before import

## Steps

1. Detect import mode from arguments:
   - no args: scan known PRD and OpenSpec locations.
   - path: import only that PRD file or folder.
   - change id: import/backfill only that OpenSpec change.
   - `--dry-run`: report planned changes without writing.
2. Discover canonical PRDs, legacy PRDs, active OpenSpec changes, and archived OpenSpec changes.
3. Build an import map: matched PRD + OpenSpec change, PRD without OpenSpec, OpenSpec without PRD, skipped items, and ambiguous candidates.
4. Ask only for ambiguous choices. Otherwise proceed automatically.
5. Create or update `.skillgrid/prd/PRD<NN>_<slug>.md`, preserving source intent and adding Skillgrid title-block fields when missing.
6. Update `.skillgrid/prd/INDEX.md` in execution order.
7. Create or update `.skillgrid/tasks/context_<change-id>.md` when a matched OpenSpec change exists.
8. Save a concise Engram summary when hybrid persistence is active.

## Boundaries

- Do not delete legacy PRDs.
- Do not overwrite canonical PRDs without reading them first.
- Do not create remote issues during import.
- Do not implement product code.

## Completion Report

Report imported files, updated index or handoff files, skipped files, ambiguous items, inferred statuses, memory saves, and the recommended next command.

</process>
