# SDD Artifact Store Modes Contract

Defines required behavior by `artifactStore.mode`.

## `engram`

- Read artifacts via Engram (`mem_search` + `mem_get_observation`).
- Persist produced artifacts via `mem_save`.
- Update mutable artifacts with `mem_update` when applicable.

## `openspec`

- Read/write artifacts on filesystem under `openspec/changes/{change-name}/`.
- Do not rely on Engram for required phase outputs.

## `hybrid`

- Do both:
  - filesystem artifacts
  - Engram persistence (`mem_save` / `mem_update`)

## `none`

- Return results inline only.
- Do not write filesystem artifacts.
- Do not call Engram persistence APIs.

## Cross-Mode Rule

Skipping required persistence breaks downstream phases and is treated as a workflow defect.
