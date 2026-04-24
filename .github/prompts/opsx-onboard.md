---
description: Guided onboarding - walk through a complete OpenSpec workflow cycle with narration
---

## Alias of Skillgrid (canonical procedures live in `skillgrid-*`)

This command is a **short entrypoint**. Execute the same work using the **canonical** markdown under **`/skillgrid-init`** through **`/skillgrid-finish`**.

**Suggested first cycle**

1. **`/skillgrid-init`** — Bootstrap `.skillgrid/`, persistence mode (`hybrid` recommended), `openspec/` when applicable, Engram `mem_save` when applicable.
2. **`/skillgrid-plan`** — PRD + `openspec new change "<name>"` + `openspec status` / `openspec instructions` artifact loop.
3. **`/skillgrid-breakdown`** — Sync PRD checklists with `openspec/changes/<id>/tasks.md`.
4. **`/skillgrid-apply`** — `openspec instructions apply` loop; keep PRD and `tasks.md` aligned.
5. **`/skillgrid-review`** — Verify artifacts and code; write `verify-report.md` (+ Engram when used).
6. **`/skillgrid-finish`** — Optional delta → main spec sync, archive under `openspec/changes/archive/`, PR/ship.

**Preflight:** `openspec --version` (install OpenSpec CLI per product docs). For narrative context see `docs/wokflow.md`.

## Completion report (required)

After the steps you ran, give a **Session wrap-up**: **(1)** which **`/skillgrid-*`** phases you completed and artifacts touched, **(2)** token usage if the product shows it (else **`Token usage: not shown in this environment`**), **(3)** next command in the numbered list you have not done yet (e.g. after init → **`/skillgrid-plan`**) or the linked **Completion report** for the last command you executed in full.
