# SDD Phase — Common Protocol

Boilerplate identical across SDD phase skills. Sub-agents should load this alongside their phase-specific `SKILL.md`.

Executor boundary: every SDD phase agent is an EXECUTOR, not an orchestrator. Do the phase work directly. Do not launch sub-agents unless the phase skill explicitly permits it.

## A. Skill Loading

1. Check if orchestrator injected a `## Project Standards (auto-resolved)` block. If present, follow it and do not read additional `SKILL.md` files.
2. If not injected, follow explicit `SKILL: Load` instructions.
3. If neither exists, fallback to skill registry discovery:
   - `mem_search(query: "skill-registry", project: "{project}")`
   - `mem_get_observation(id)` for full registry content
   - fallback file `.atl/skill-registry.md` if it exists
4. If no registry exists, proceed with the phase skill only.

## B. Shared Contracts

Load these shared contracts when they apply:

- `skills/_shared/sdd-engram-retrieval.md`
- `skills/_shared/sdd-artifact-store-modes.md`
- `skills/_shared/sdd-return-envelope.md`
- `skills/_shared/sdd-label-gate-contract.md`
- `skills/_shared/sdd-enforcement-contract.md`

## C. Return Envelope

Use the canonical envelope defined in `skills/_shared/sdd-return-envelope.md`.

## D. Compatibility Note

Legacy filename `skills/_shared/sdd-pahse-common.md` is deprecated (typo). Use this file as canonical.
