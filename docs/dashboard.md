# Skillgrid Dashboard

The Skillgrid dashboard is a local, file-backed web UI for PRDs, workflow state, preview artifacts, subagent activity, and graphify output.

Run it from the repository root:

```bash
node .skillgrid/scripts/skillgrid-ui.mjs
```

Then open `http://127.0.0.1:8787`.

## Data Sources

The dashboard reads existing Skillgrid artifacts. It does not add a database, daemon, telemetry layer, platform adapter, or workflow engine.

| Source | Purpose |
|---|---|
| `.skillgrid/prd/PRD*.md` | Kanban cards, title, status, PRD metadata, linked change id. |
| `.skillgrid/tasks/context_<change-id>.md` | Current phase, current status, blockers, AFK-ready work, next action. |
| `.skillgrid/tasks/events/<change-id>.jsonl` | Append-only workflow and subagent event timeline. |
| `.skillgrid/preview/*.html` | Preview links on Board cards and Workflow artifacts. |
| `.skillgrid/tasks/research/<change-id>/` | Research and subagent report artifact links. |
| `graphify-out/graph.html` | Optional Graph tab. |

## Views

### Board

The Board view keeps the PRD Kanban columns:

```text
draft | todo | inprogress | devdone | done
```

Cards show the PRD title, file, change id, current handoff phase, latest event status, blocked badge, preview link, external link, and a Workflow shortcut.

Dragging a card or changing its status dropdown updates only the PRD `status` field.

### Workflow

The Workflow view shows one selected PRD/change at a time:

- current phase and next action from the handoff
- HITL blockers and AFK-ready work
- phase status for brainstorm, design, plan, breakdown, apply, test, security, validate, and finish
- chronological event log
- artifact links for PRD, handoff, previews, research reports, and other event artifacts

### Subagents

The Subagents view is a read-only activity feed over event-log entries that identify delegated work.

An event appears there when it includes one of:

- `agent`
- `subagent`
- `role`
- a `node`, `phase`, or `summary` mentioning agent, reviewer, implementer, critic, auditor, or verifier

Use this for seeing what subagents did, their status, blockers, and output artifacts without reading the full event timeline for every PRD.

### Graph

The Graph tab appears only when `graphify-out/graph.html` exists. It opens the graphify HTML report in a new browser tab.

## PRD Metadata

The dashboard supports YAML frontmatter and Skillgrid title-block metadata. Recommended title-block fields:

```markdown
- **Spec / change:** `openspec/changes/<change-id>/`
- **Session context:** `.skillgrid/tasks/context_<change-id>.md`
- **Status:** `draft`
- **Preview:** `.skillgrid/preview/<change-id>-options.html`
- **External:** local
```

If no change id is found, the dashboard derives one from the PRD filename.

## Preview Discovery

Preview links are found in two ways:

1. Explicit PRD metadata: `Preview`.
2. Convention-based discovery:

```text
.skillgrid/preview/<change-id>*.html
.skillgrid/preview/<prd-slug>*.html
```

Explicit metadata is preferred when the preview filename does not start with the change id or PRD slug.

## Handoff Event Log

Events live at:

```text
.skillgrid/tasks/events/<change-id>.jsonl
```

Each line is one JSON object:

```json
{"time":"2026-04-28T08:00:00Z","changeId":"planning-dashboard","prd":"PRD02_dashboard.md","node":"apply","phase":"apply","status":"started","summary":"Started AFK slice 2","artifacts":[".skillgrid/tasks/context_planning-dashboard.md"]}
```

Supported fields:

- `time`: ISO timestamp.
- `changeId`: OpenSpec/Skillgrid change id.
- `prd`: optional PRD filename.
- `node`: workflow node or command id.
- `phase`: Skillgrid phase.
- `status`: `started`, `progress`, `passed`, `failed`, `blocked`, `completed`, or `skipped`.
- `summary`: one-line event summary.
- `blocker`: optional HITL or technical blocker.
- `artifacts`: optional paths to PRDs, previews, reports, handoff files, tests, or research.
- `external`: optional issue key or URL.
- `agent` or `subagent`: optional delegated agent name.
- `role`: optional role, such as `implementer`, `spec-reviewer`, or `code-quality-reviewer`.
- `task`: optional delegated task label.
- `output`: optional primary output path.

The event log is the timeline. The handoff file remains the current-state summary.

## HTTP API

Read APIs:

```text
GET /api/prds
GET /api/prds/:file
GET /api/dashboard
GET /api/agents
GET /api/changes/:changeId/events
GET /api/changes/:changeId/artifacts
GET /preview/<file>
GET /graphify/graph.html
```

Write API:

```text
PATCH /api/prds/:file
```

No event-writing API exists. Skillgrid commands and agents append JSONL events directly on disk.
