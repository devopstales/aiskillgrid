# SDD Engram Retrieval Contract

Use this whenever `artifactStore.mode` is `engram` or `hybrid`.

## Mandatory Rules

- `mem_search` returns previews only; never use preview text as source material.
- Always call `mem_get_observation(id)` for every artifact that influences decisions.
- Run all related `mem_search` calls in parallel.
- Run all related `mem_get_observation` calls in parallel.

## Canonical Sequence

1. Search IDs:
   - `mem_search(query: "sdd/{change-name}/{artifact}", project: "{project}")`
2. Retrieve full artifacts:
   - `mem_get_observation(id: {artifact_id})`
3. Work only from full retrieved content.

## Persistence Reminder

When writing artifacts, use stable `topic_key` so repeated saves upsert instead of duplicating.
