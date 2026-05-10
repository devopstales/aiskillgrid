---
name: parallel-delegate
description: >
  Parent workflow: split independent work across parallel sub-agents (Cursor Task or equivalent), define handoffs, merge results safely.
  Trigger: Multiple independent files or modules; parallel research or doc passes; parallel review slices; user asks to run subagents in parallel or to delegate lanes.
  Not a replacement for SDD phase skills (sdd-propose, sdd-tasks, sdd-apply); use those for full change lifecycle.
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
---

## Role

You are the **coordinator**. You decompose, launch parallel children, then **integrate** their outputs into one coherent outcome for the user. You do not duplicate SDD artifact production inside each child unless the task is explicitly SDD-shaped.

## When to use

- Several **non-overlapping** scopes (different paths, packages, or research questions).
- **Read-heavy** passes (explore, impact, docs) that do not all edit the same files.
- **Review or audit** slices on disjoint areas.

## When not to use

- Single file or tightly coupled edits (serial is simpler and safer).
- Children would **race on the same mutable paths** — merge design first or serialize writes.
- User asked for a **full SDD / OpenSpec phase** — use the matching `sdd-*` / `openspec-*` skill instead.

## Decomposition

1. List work packages; for each, note **inputs**, **outputs**, and **files touched**.
2. Mark **read-only** vs **write** lanes; prefer parallel only when write sets are disjoint or children are read-only.
3. If two packages share a file, assign **one** owner for that file or split read vs write phases.

## Child prompt template (paste per sub-agent)

Use a fixed shape so merges stay predictable:

```text
## Role
<one sentence: reviewer | researcher | implementer for scope X only>

## Project standards
<optional: paste compact rules from .skillgrid/project/SKILL_REGISTRY.md for this scope>

## In scope
- Paths / questions: …

## Out of scope
- …

## Inputs
- …

## Deliverable
- Format: bullets | patch list | table | file paths
- Stop when: …

## Constraints
- Tools / branches / readonly: …
```

## Launch

- Use the host’s **parallel agent / Task** capability: one Task per package, **distinct** `description` strings, prompts that include the template above.
- Prefer **readonly** children for exploration; use **at most one** writer per repo path per wave.

## Merge

1. Collect each child’s deliverable; drop duplicates.
2. **Conflict scan** — if two children touched the same path or contradicted facts, reconcile in the coordinator turn (ask user if product intent is unclear).
3. Produce **one** user-facing summary: decisions, file list, open risks, next step.
4. If children wrote files, run quick **sanity** (build/tests) once after merge, not N times in parallel writers.

## Stop

Stop when the user’s original question is answered or the delegated wave is integrated. Escalate to serial work if merge conflicts or ambiguous ownership appear.

## Related

- **`skillgrid-skill-registry`** — inject compact rules into child prompts.
- **`skillgrid-vertical-slices`** — split work into testable slices before delegating implementation.
