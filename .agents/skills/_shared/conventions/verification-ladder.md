# Verification Ladder Convention (shared across all SDD skills)

Pick the lowest level that still gives defensible evidence. Record the floor in `change.md` (`Verification floor: L1–L4`); trivial changes may omit it.

## L0: Inspection only

Use only when the task is analysis/review-only or the environment cannot execute checks. State explicitly the conclusion is inspection-based.

## L1: Structural check

Syntax validation, parse/build check, schema validation, focused static analysis. For trivial edits and structure-only changes.

## L2: Targeted execution

Focused unit test, single command path, minimal repro step. For ordinary bug fixes, localized features, small refactors.

## L3: Behavioral proof

Before/after comparison on the affected path, targeted integration test, focused manual or scripted smoke test. For regressions, shared-path fixes, user-visible behavior changes.

## L4: High-risk validation

Multiple checks across affected layers, rollback-aware migration verification, destructive-command safety verification. For risky refactors, migrations, auth/payment/prod-sensitive paths.

## Selection rule

- `trivial` → L1 (or omit floor)
- `standard` → L2
- `risky` → L3 or L4
- `review` → state whether any claim is below execution proof
