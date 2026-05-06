# Skillgrid file templates

Canonical blanks for Skillgrid and OpenSpec artifacts. **Naming:** every template is `template-<kebab-case>.md` (lowercase, hyphenated). Copy the file and rename to the target path (for example `template-prd.md` → `.skillgrid/prd/PRD01_feature.md`).

Planning logic and how pieces connect: **`docs/skillgrid-templates-and-logic.md`**.

| File | Typical destination |
|------|---------------------|
| `template-adr.md` | `.skillgrid/adr/NNNN-short-title.md` (MADR only; keep `adr/` free of README or other non-ADR files) |
| `template-prd.md` | `.skillgrid/prd/PRD<NN>_<slug>.md` |
| `template-index.md` | `.skillgrid/prd/INDEX.md` |
| `template-openspec-tasks.md` | `openspec/changes/<id>/tasks.md` |
| `template-openspec-slice-spec.md` | `openspec/changes/<id>/specs/<slice>/spec.md` |
| `template-handoff-context.md` | `.skillgrid/tasks/context_<change-id>.md` |
| `template-project.md` | `.skillgrid/project/PROJECT.md` |
| `template-architecture.md` | `.skillgrid/project/ARCHITECTURE.md` |
| `template-structure.md` | `.skillgrid/project/STRUCTURE.md` |
| `template-design.md` | Repo root `DESIGN.md` |
