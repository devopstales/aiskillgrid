# Web UI

The AISkillGrid web UI is a local dashboard for project workflow state. It makes agent work visible without requiring a hosted platform.

Start it from the project root:

```bash
node .skillgrid/scripts/skillgrid-ui.mjs
```

Then open the local address printed by the server. The common default is:

```text
127.0.0.1:8787
```

## Why The Web UI Exists

AI work can become invisible when it lives only in chat. The web UI gives users and stakeholders a place to see:

- What PRDs exist.
- Which phase a change is in.
- What is blocked.
- What subagents did.
- Which previews exist.
- Whether graph output is available.
- Whether Engram shared-memory export metadata exists.
- Whether the project skill registry exists.
- What event history has been recorded.

This turns AISkillGrid from a prompt library into an observable workflow.

## Main Views

```mermaid
flowchart TD
  WebUI[Local Web UI] --> Board[Board View]
  WebUI --> Workflow[Workflow View]
  WebUI --> Subagents[Subagents View]
  WebUI --> Previews[Preview Links]
  WebUI --> Graph[Graph View]
  Board --> PRD[PRD Files]
  Workflow --> Events[Event Logs]
  Subagents --> Reports[Research And Reports]
  Graph --> Graphify[graphify Output]
  WebUI --> Memory[Engram And Registry Status]
```

## Board View

The Board view shows PRDs in workflow columns.

It helps answer:

- What work exists?
- What status is each PRD in?
- Which items are blocked?
- Which items have previews?
- Which item should be opened in the Workflow view?

When the workflow status changes, the dashboard updates the PRD status field. That keeps the board tied to files instead of a hidden database.

## Workflow View

The Workflow view focuses on one selected PRD or change.

It can show:

- Current phase.
- Current state.
- Next recommended action.
- HITL blockers.
- AFK-ready work.
- Event timeline.
- Artifacts such as PRD, handoff, previews, research, tests, and review reports.

This is the view to use when asking, “Can the agent keep going, or does a human need to decide something?”

## Subagents View

The Subagents view collects delegated work activity.

It is useful for:

- Seeing what reviewers, researchers, critics, auditors, and verifiers did.
- Finding their output files.
- Spotting blockers.
- Checking whether independent reports agree.

This makes multiagent work easier to trust because the activity is visible.

## Preview Links

When a workflow produces HTML previews, the dashboard can surface them from:

```text
.skillgrid/preview/
```

This is especially useful for UI design, prototypes, visual comparisons, or generated documentation pages.

## Graph View

When graphify output exists, the dashboard can expose the graph report.

Typical source:

```text
graphify-out/
```

This gives users a quick way to jump from workflow state into codebase structure.

## Data Sources

The dashboard reads files that already belong to the Skillgrid workflow:

| Source | What It Powers |
|---|---|
| `.skillgrid/prd/` | Board cards and product intent |
| `.skillgrid/tasks/context_<change-id>.md` | Current state and next action |
| `.skillgrid/tasks/events/<change-id>.jsonl` | Timeline and subagent activity |
| `.skillgrid/tasks/research/<change-id>/` | Reports and research artifacts |
| `.skillgrid/preview/` | Preview links |
| `graphify-out/` | Graph view |
| `.engram/manifest.json` | Engram export counts, when team memory sync is used |
| `.skillgrid/project/SKILL_REGISTRY.md` | Skill registry availability and skill count |

No separate database is required for the core local dashboard model.

The dashboard reads `.engram/manifest.json` directly when it exists. It does not call the Engram CLI or expose memory contents.

## Practical Advantage

The web UI gives AISkillGrid a visible operating surface. New users can understand what is happening without reading every artifact by hand. Leads and reviewers can inspect state without asking the agent to summarize itself.

That visibility is a major difference between a full workflow solution and a collection of isolated prompts.
