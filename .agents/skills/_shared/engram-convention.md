# Engram Artifact Convention (reference documentation)

NOTE: Critical engram calls (`mem_search`, `mem_save`, `mem_get_observation`) are inlined directly in each skill's SKILL.md. This document is supplementary reference — sub-agents do NOT need to read it to function.

## Naming Rules

ALL SDD artifacts persisted to Engram MUST follow this deterministic naming:

```
title:     sdd/{change-name}/{artifact-type}
topic_key: sdd/{change-name}/{artifact-type}
type:      architecture
project:   {detected or current project name}
scope:     project
```

### Artifact Types

| Artifact Type | Produced By | Description |
|---------------|-------------|-------------|
| `explore` | sdd-explore | Exploration analysis |
| `proposal` | sdd-propose | Change proposal |
| `spec` | sdd-spec | Delta specifications (all domains concatenated) |
| `design` | sdd-design | Technical design |
| `tasks` | sdd-tasks | Task breakdown |
| `apply-progress` | sdd-apply | Implementation progress (one per batch) |
| `verify-report` | sdd-verify | Verification report |
| `archive-report` | sdd-archive | Archive closure with lineage |
| `state` | orchestrator | DAG state for recovery after compaction |

Exception: `sdd-init` uses `sdd-init/{project-name}` as both title and topic_key.

### State Artifact

```
mem_save(
  title: "sdd/{change-name}/state",
  topic_key: "sdd/{change-name}/state",
  type: "architecture",
  project: "{project}",
  content: "change: {change-name}\nphase: {last-phase}\nartifact_store: engram\nartifacts:\n  proposal: true\n  specs: true\n  design: false\n  tasks: false\ntasks_progress:\n  completed: []\n  pending: []\nversion:\n  id: {monotonic-int}\n  previous_id: {prior-monotonic-int-or-null}\ntimestamps:\n  created_at: {ISO-8601-UTC}\n  updated_at: {ISO-8601-UTC}"
)
```

Recovery: `mem_search("sdd/{change-name}/state")` → `mem_get_observation(id)` → parse YAML → restore state.

## Versioning and Timestamps (required)

Every persisted Engram artifact must carry explicit version metadata in `content`.

Required fields:

```yaml
version:
  id: <monotonic integer, starts at 1>
  previous_id: <integer or null>
timestamps:
  created_at: <ISO-8601 UTC timestamp, immutable>
  updated_at: <ISO-8601 UTC timestamp, changes on every write>
```

Rules:

- `created_at` is set on first write and never rewritten.
- `updated_at` is set to write time on each upsert/update.
- `version.id` increments by exactly 1 per successful write for the same `topic_key`.
- `version.previous_id` points to the last committed `version.id` (or `null` on first write).
- Writers must read latest artifact state before modifying versioned content.

## Conflict Resolution (required)

When multiple writes race on the same `topic_key`, resolve deterministically:

1. **Version wins first**: Higher `version.id` always wins.
2. **Timestamp tie-breaker**: If `version.id` matches, newer `timestamps.updated_at` wins.
3. **Stable final tie-breaker**: If both match (clock skew/replay), keep the existing record and reject the incoming write as duplicate.
4. **Non-overlapping merge**: If both branches edited different top-level sections, merge sections and emit a new version (`id = max + 1`, `previous_id = max`).
5. **Overlapping edit conflict**: If both changed the same section semantically, preserve winner and append a `conflict_note` block documenting discarded candidate metadata (writer, version, timestamp).

## Recovery Protocol (2 steps)

```
Step 1: mem_search(query: "sdd/{change-name}/{artifact-type}", project: "{project}") → truncated preview + ID
Step 2: mem_get_observation(id: {observation-id}) → complete content
```

When retrieving multiple artifacts, group all searches first, then all retrievals:

```
STEP A — SEARCH (get IDs only):
  mem_search(query: "sdd/{change-name}/proposal", ...) → save ID
  mem_search(query: "sdd/{change-name}/spec", ...) → save ID
  mem_search(query: "sdd/{change-name}/design", ...) → save ID

STEP B — RETRIEVE FULL CONTENT (mandatory):
  mem_get_observation(id: {proposal_id})
  mem_get_observation(id: {spec_id})
  mem_get_observation(id: {design_id})
```

Loading project context:
```
mem_search(query: "sdd-init/{project}", project: "{project}") → get ID
mem_get_observation(id) → full project context
```

## Writing Artifacts

Standard write:
```
mem_save(
  title: "sdd/{change-name}/{artifact-type}",
  topic_key: "sdd/{change-name}/{artifact-type}",
  type: "architecture",
  project: "{project}",
  content: "{full markdown content}"
)
```

Concrete example — saving a proposal for `add-dark-mode`:
```
mem_save(
  title: "sdd/add-dark-mode/proposal",
  topic_key: "sdd/add-dark-mode/proposal",
  type: "architecture",
  project: "my-app",
  content: "version:\n  id: 1\n  previous_id: null\ntimestamps:\n  created_at: 2026-05-08T00:00:00Z\n  updated_at: 2026-05-08T00:00:00Z\n\n## Proposal\n\nAdd dark mode toggle..."
)
```

Update existing artifact (when you have the observation ID):
```
mem_update(id: {observation-id}, content: "{updated full content}")
```

Use `mem_update` when you have the exact ID. Use `mem_save` with same `topic_key` for upserts.

### Browsing All Artifacts for a Change

```
mem_search(query: "sdd/{change-name}/", project: "{project}")
→ Returns all artifacts for that change
```

## Why This Convention

- Deterministic titles → recovery works by exact match
- `topic_key` → enables upserts without duplicates
- `sdd/` prefix → namespaces all SDD artifacts
- Two-step recovery → search previews are always truncated; `mem_get_observation` is the only way to get full content
- Lineage → archive-report includes all observation IDs for complete traceability
