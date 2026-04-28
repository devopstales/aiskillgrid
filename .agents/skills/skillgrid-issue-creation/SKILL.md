---
name: skillgrid-issue-creation
description: >
  Creates and links external tracker issues for Skillgrid only when configured by .skillgrid/config.json.
  Trigger: Creating GitHub, GitLab, or Jira issues from Skillgrid PRDs, specs, or vertical slices.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a Skillgrid phase may create, link, or close external tracker work for a PRD, OpenSpec change, or vertical slice.

Do not use this skill for generic issue creation outside Skillgrid.

## Critical Patterns

### Provider Preflight

Always read `.skillgrid/config.json` before creating or linking remote issues.

Supported `ticketing.provider` values:

| Provider | Behavior |
|---|---|
| `local` or missing | Do not create remote issues. Use PRD `Status:` and `.skillgrid/prd/INDEX.md`. |
| `github` | Use GitHub issues and links when the user wants tracker-backed work. |
| `gitlab` | Use GitLab issues and links when configured. |
| `jira` | Use Jira issue keys and transition hints when configured. |

Never assume GitHub just because the repository is hosted on GitHub.

### Local-First Rule

Skillgrid works without a remote tracker. In `local` mode:

- PRD files are the ticket source.
- `.skillgrid/prd/INDEX.md` is the work board.
- `Status:` is the lifecycle field.
- External issue columns are optional and should not be created unless useful.

### Issue Scope

Create issues for reviewable slices, not for vague epics.

Good issue units:

- one vertical slice from `skillgrid-vertical-slices`
- one `[HITL]` decision needing a human answer
- one `[AFK]` implementation slice ready for autonomous apply
- one security or validation blocker

Avoid issue spam for every tiny checkbox when the PRD and `tasks.md` already track the work well.

### QA Session Intake

When the user is reporting bugs conversationally:

1. Listen first and ask at most 2-3 short clarifying questions.
2. Capture expected behavior, actual behavior, reproduction steps, and whether the issue is consistent or intermittent.
3. Explore the relevant codebase area only to learn domain language and behavior boundaries; do not turn the issue body into an implementation diagnosis.
4. Decide whether the report is one issue or several independently fixable issues.
5. If using a remote tracker, create issues in dependency order so blockers can reference real issue numbers.

Remote issues from QA should be durable and user-facing. Avoid file paths, line numbers, private function names, or guesses about root cause unless the user explicitly asks for an implementation ticket.

### Issue Body

Issue bodies should include:

- PRD link: `.skillgrid/prd/PRD<NN>_<slug>.md`
- OpenSpec change: `openspec/changes/<change-id>/`
- slice or task range
- acceptance criteria
- verification command or manual check
- blocker state, if `[HITL]`

Use this template:

```markdown
## Summary

<One sentence describing the slice or blocker.>

## Source artifacts

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`
- OpenSpec change: `openspec/changes/<change-id>/`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`

## Scope

### In

- <Included work>

### Out

- <Explicitly excluded work>

## Acceptance criteria

- [ ] <Observable result>
- [ ] <Quality or edge-case criterion>

## Verification

- [ ] `<command>`
- [ ] <manual check, if needed>

## Skillgrid task mapping

- Slice: `<slice name or number>`
- Tag: `[AFK]` or `[HITL]`
- Blocks: `<dependent issue/task or none>`

## Notes

<Provider-specific labels, Jira component, GitHub milestone, or human decision needed.>
```

For user-reported bugs, prefer this behavior-oriented template:

```markdown
## What happened

<Actual behavior in the user's words, using project domain language.>

## What I expected

<Expected behavior.>

## Steps to reproduce

1. <Concrete step>
2. <Concrete step>

## Scope and blockers

- Slice: <single issue or part of a breakdown>
- Blocked by: <issue/task or None>

## Additional context

<Concise observations that help reproduce or understand impact. No stale file paths or line numbers.>
```

### INDEX External Column Template

```markdown
| Order | PRD | Status | Spec / change | External |
|---|---|---|---|---|
| 01 | [`PRD01_<slug>.md`](PRD01_<slug>.md) | `todo` | `openspec/changes/<change-id>/` | `<issue-key-or-url>` |
```

### Backlinks

When an issue is created or linked:

- update `.skillgrid/prd/INDEX.md` with an optional `External` value if the table uses one
- add a link in the PRD or task section when it helps navigation
- include the issue link/key in `.skillgrid/tasks/context_<change-id>.md` if the handoff exists
- save a short Engram summary when hybrid persistence is active and the issue changes workflow state

## Commands

```bash
gh issue create --title "<title>" --body "<body>"
glab issue create --title "<title>" --description "<body>"
```

Use Jira through the project’s configured tool or MCP if present; do not invent credentials or project keys.

## Resources

- Provider configuration: `.skillgrid/config.json`
- PRD lifecycle: `skillgrid-prd-artifacts`
- Slice rules: `skillgrid-vertical-slices`
- Handoff rules: `skillgrid-filesystem-handoff`
- Provider-specific inspiration: `engram-issue-creation`
