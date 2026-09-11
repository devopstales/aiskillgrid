# Task Classification Convention (shared across all SDD skills)

Classify the change before planning or coding. Keep trivial work light — no ceremony for local edits.

## `trivial`

Use for: typo fixes, rename-only changes, tiny local copy/logging edits, single-line adjustments with obvious blast radius.

Required: keep explanation short, no big plan, lightest meaningful check (verification ladder L1).

Failure mode to avoid: turning a tiny change into a ceremony.

Classification in `change.md` is **optional** for trivial — one line (`Classification: trivial`) suffices, verification floor may be omitted.

## `standard`

Use for: ordinary bug fixes, bounded feature work, localized test updates, straightforward config changes.

Required: inspect the exact files involved, restate the task concretely, choose the narrowest viable change, run targeted verification (ladder L2).

Failure mode to avoid: "probably correct" coding without execution.

## `risky`

Use for: shared infrastructure or auth paths, migrations, concurrency/state-heavy code, multi-file refactors, destructive or high-impact CLI workflows.

Required: surface assumptions and alternate interpretations early, minimize change surface, choose a stronger verification level (ladder L3 or L4), mention residual risk and rollback concerns.

Failure mode to avoid: confident action on under-specified requirements.

## `review`

Use for: PR review, self-review, diff audit, logic or regression check.

Required: findings first, cite file and concrete impact, prioritize correctness/regressions/missing tests over style, mention residual gaps if nothing is proven.

Failure mode to avoid: style commentary before behavioral risk.

## Tie-break rule

If between two classes: choose the lighter class for truly local edits; choose the stricter class if blast radius or uncertainty is material.
