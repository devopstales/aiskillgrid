# Skillgrid PRD Index

## Execution snapshot

- **Current phase:** <plan | breakdown | apply | validate | finish>
- **Active change:** `openspec/changes/<change-id>/`
- **Active slice:** `<slice-slug>` (if any) — `specs/<slice-slug>/spec.md`
- **In progress:** <bullets>
- **Recently completed:** <bullets>
- **Next up:** <bullets>
- **Discovered during work:** <bullets>
- **Open questions / blockers:** <bullets>
- **Session notes:** <2–5 lines>

## PRD index (dependency order)

Local Kanban dashboard:

Run `node .skillgrid/scripts/skillgrid-ui.mjs`, then open `http://127.0.0.1:8787`.

| Order | PRD | Status | Spec / change | Depends on | External |
|---|---|---|---|---|---|
| 01 | [`PRD01_<slug>.md`](PRD01_<slug>.md) | `draft` | [`openspec/changes/<change-id>/`](../../openspec/changes/<change-id>/) | — | local |
