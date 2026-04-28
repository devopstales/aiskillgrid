---
name: skillgrid-import-artifacts
description: >
  Imports existing PRDs and OpenSpec changes into Skillgrid's canonical `.skillgrid/prd/` and INDEX workflow.
  Trigger: Importing legacy PRDs, backfilling PRDs for existing OpenSpec changes, or normalizing brownfield planning artifacts.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a repository already has PRDs, OpenSpec changes, or archived OpenSpec work that should be represented in Skillgrid's canonical `.skillgrid/prd/` index.

This skill is used directly by `/skillgrid-import` and automatically by `/skillgrid-explore` when exploration finds existing artifacts without canonical PRD coverage.

## Critical Patterns

### Conservative Import

- Do not delete source PRDs from legacy paths.
- Do not overwrite canonical `.skillgrid/prd/` files without reading them first.
- Do not silently merge ambiguous PRD/OpenSpec matches.
- Do not create remote issues during import.
- Preserve source intent; normalize structure only enough to fit Skillgrid.

### Discovery Sources

Scan these locations:

```text
.skillgrid/prd/PRD*.md
prd/PRD*.md
prd/**/*.md
docs/PRD/**/*.md
docs/prd/**/*.md
**/PRD*.md
openspec/changes/*/
openspec/changes/archive/*/
```

Ignore generated mirrors and dependency folders.

### Match Strategy

Match PRDs to OpenSpec changes in this order:

1. Explicit `openspec/changes/<change-id>/` link in the PRD.
2. `Spec / change:` field in the PRD title block.
3. Shared slug in filename and change directory.
4. Title similarity after normalizing case, spaces, hyphens, and underscores.
5. If still unclear, mark ambiguous and ask the user.

### Status Inference

| Source state | Skillgrid `Status:` |
|---|---|
| OpenSpec change under `openspec/changes/archive/` | `done` |
| `tasks.md` has one or more checked tasks | `inprogress` |
| `tasks.md` exists but no checked tasks | `todo` |
| proposal/design/spec exists but no tasks | `draft` |
| imported PRD with no OpenSpec change | `draft` |

If the source PRD already has a Skillgrid status, keep it unless the OpenSpec archive state clearly proves `done`.

### Numbering

Use `skillgrid-prd-artifacts` rules:

- List existing `.skillgrid/prd/PRD*.md`.
- Preserve existing numbers.
- Assign imported PRDs the next free contiguous number unless the source has an obvious execution order.
- Update `.skillgrid/prd/INDEX.md` after every import batch.

## Import Map Template

Build this before writing files:

```markdown
# Skillgrid Import Map

## Matched PRDs and OpenSpec changes

| Source PRD | OpenSpec change | Target PRD | Status | Confidence |
|---|---|---|---|---|
| `<path>` | `openspec/changes/<id>/` | `.skillgrid/prd/PRD<NN>_<slug>.md` | `draft` | explicit-link |

## PRDs without OpenSpec changes

| Source PRD | Target PRD | Status | Note |
|---|---|---|---|
| `<path>` | `.skillgrid/prd/PRD<NN>_<slug>.md` | `draft` | no change found |

## OpenSpec changes without PRDs

| OpenSpec change | Target PRD | Status | Source artifacts |
|---|---|---|---|
| `openspec/changes/<id>/` | `.skillgrid/prd/PRD<NN>_<slug>.md` | `draft` | proposal, design, specs, tasks |

## Ambiguous items

| Item | Candidates | Question |
|---|---|---|
| `<path-or-change>` | `<candidate list>` | <what user must choose> |
```

## Canonical PRD Creation

When importing a legacy PRD:

1. Keep original headings and content where possible.
2. Add or normalize the Skillgrid title block.
3. Add missing `Spec / change`, `Session context`, and `Status` fields.
4. Add an `Imported from:` note when the source path differs from the target.
5. Do not summarize away detailed requirements.

When backfilling from OpenSpec only:

1. Use `proposal.md` for problem, why, goals, and scope.
2. Use `design.md` for codebase touchpoints, technical constraints, and risks.
3. Use delta specs for success criteria.
4. Use `tasks.md` for implementation tasks.
5. Mark missing information as open questions instead of inventing it.

## Import Report Template

```markdown
## Skillgrid Import Report

### Imported

- `<source>` -> `.skillgrid/prd/PRD<NN>_<slug>.md` (`Status: <status>`)

### Updated

- `.skillgrid/prd/INDEX.md`
- `.skillgrid/tasks/context_<change-id>.md`

### Skipped

- `<path>` — <reason>

### Ambiguous

- `<path-or-change>` — needs user choice: <question>

### Next command

<Usually `/skillgrid-explore`, `/skillgrid-plan`, or `/skillgrid-breakdown`.>
```

## Dry Run

When `--dry-run` is requested, produce the import map and report only. Do not create, rename, or update files.

## Resources

- PRD rules: `skillgrid-prd-artifacts`
- OpenSpec alignment: `skillgrid-spec-artifacts`
- Handoff files: `skillgrid-filesystem-handoff`
- Persistence: `skillgrid-hybrid-persistence`
