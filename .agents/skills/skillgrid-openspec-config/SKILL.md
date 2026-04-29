---
name: skillgrid-openspec-config
description: >
  Keeps openspec/config.yaml aligned with Skillgrid ticketing and artifact-store configuration without clobbering local context.
  Trigger: Initializing OpenSpec, refreshing config.yaml, or reconciling Skillgrid config with OpenSpec instructions.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when Skillgrid commands inspect, create, or update `openspec/config.yaml` in a repository that also uses `.skillgrid/config.json`.

## Critical Patterns

### Read Both Configs

Before changing `openspec/config.yaml`, read:

- `.skillgrid/config.json`
- `openspec/config.yaml` if it exists

Preserve existing project context unless it contradicts the Skillgrid configuration.

### Required Context Lines

When OpenSpec is active, `openspec/config.yaml` context should mention:

- ticketing provider from `.skillgrid/config.json`
- artifact store mode from `.skillgrid/config.json`
- PRD workflow source, status ids, and phase-to-status mapping from `.skillgrid/config.json`
- canonical PRD path: `.skillgrid/prd/`
- handoff path: `.skillgrid/tasks/context_<change-id>.md`

### Merge, Do Not Clobber

- Prefer small edits over rewriting the whole file.
- Keep existing useful project context.
- Do not erase team-specific instructions.
- If the config is missing or incoherent, explain the merge and show the intended replacement.

### Context Size

Keep OpenSpec context concise. Large research, architecture maps, and full PRDs belong in their own files and should be linked or summarized.

### config.yaml Context Template

Adapt to the local OpenSpec schema, preserving existing keys:

```yaml
project: <project-name>
schema: spec-driven
context: |
  This repository uses Skillgrid with OpenSpec.

  Ticketing:
  - Provider: <local|github|gitlab|jira>
  - Local mode uses `.skillgrid/prd/INDEX.md` and PRD `Status:` as the work tracker.
  - Remote issues are created only when the provider is configured.

  PRD workflow:
  - Source: <preset|custom|provider-import>
  - Statuses: <ordered status ids, e.g. draft, todo, inprogress, devdone, done>
  - Phase status map: plan=<status>, breakdown=<status>, apply=<status>, validate=<status>, finish=<status>

  Artifact store:
  - Mode: <hybrid|openspec|engram>
  - PRDs: `.skillgrid/prd/`
  - OpenSpec changes: `openspec/changes/`
  - Handoff: `.skillgrid/tasks/context_<change-id>.md`
  - Research spill: `.skillgrid/tasks/research/<change-id>/`

  Agent rules:
  - Keep PRD intent, OpenSpec artifacts, and tasks aligned.
  - Use `[HITL]` for human-blocked work and `[AFK]` for autonomous slices.
  - Do not paste raw instruction/context blocks into generated artifacts.
```

### Merge Report Template

When changing config, summarize:

```markdown
## OpenSpec Config Alignment

- **Read:** `.skillgrid/config.json`
- **Read:** `openspec/config.yaml`
- **Ticketing provider:** `<provider>`
- **Artifact store:** `<mode>`
- **PRD workflow:** `<source>; <ordered-statuses>`

## Changes made

- <Line or section updated>

## Preserved

- <Existing project context kept intact>

## Follow-up

- <None, or command to run next>
```

## Commands

```bash
openspec instructions
openspec list --json
openspec status --change "<change-id>" --json
```

## Resources

- Skillgrid config: `.skillgrid/config.json`
- OpenSpec artifacts: `skillgrid-spec-artifacts`
- Persistence: `skillgrid-hybrid-persistence`
