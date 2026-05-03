# Tasks: <change-id>

Cross-slice ordering and integration checklist. Per-slice detail belongs in `specs/<vertical-slice-slug>/spec.md` (see `template-openspec-slice-spec.md`).

## Slice map

| Slice slug | Spec | Tag | Blocked by |
|------------|------|-----|------------|
| `<slice-a>` | [spec](specs/<slice-a>/spec.md) | `[AFK]` | — |
| `<slice-b>` | [spec](specs/<slice-b>/spec.md) | `[AFK]` | `<slice-a>` |

## Dependency notes

- Record `blockedBy` / `unblocks` between slices or external decisions.

## Checklist (example)

### Wave 1

- [ ] `[AFK]` <task> — verify: `<command>`

### Wave 2

- [ ] `[HITL]` <decision or approval>
- [ ] `[AFK]` <task> — verify: `<command>`
