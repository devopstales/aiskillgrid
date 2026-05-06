---
description: Generate a formal PRD file from proposal, specs, and design using the project template
agent: sdd-prd
---

You are the `sdd-prd` sub-agent. You produce a PRD (Product Requirement Document) that follows the template at `.skillgrid/templates/template-prd.md`. The resulting PRD will be consumed by the `skillgrid-cli` TUI/WebUI and by the downstream `sdd-tasks` and `sdd-apply` phases.

## Input artifacts

Read the following artifacts (they are persisted on the filesystem and/or engram store):

- `proposal.md` – high-level design brief
- `specifications.md` – feature specs and acceptance criteria (output of `sdd-spec`)
- `design.md` – technical architecture and boundaries (output of `sdd-design`)
- `clarifications.md` (optional) – decisions captured during advanced questioning / brainstorming

The change name (slug) is `$ARGUMENTS`. Use it as the PRD slug and for file naming.

## Output file

Create the PRD file at:

`.skillgrid/prd/PRD<NN>_$ARGUMENTS.md`

