---
name: /skillgrid-init
id: skillgrid-init
category: Workflow
description: Bootstrap workflow: structure, graphify, CocoIndex, OpenCode, baseline skills
---

You are executing **`/skillgrid-init`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-init`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Establish or verify project folder structure per team conventions.
2. Initialize graphify when this repo uses it.
3. Initialize CocoIndex Code: run `ccc init` from project root if needed, then `ccc index`.
4. Initialize OpenCode layout if the project uses OpenCode.
5. Ensure OpenSpec CLI is available when using OpenSpec; run onboarding only if appropriate.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/sdd-init/SKILL.md`
- `.agents/skills/openspec-onboard/SKILL.md`
- `.agents/skills/search-first/SKILL.md`
- `.agents/skills/karpathy-guidelines/SKILL.md`
- `.agents/skills/skill-creator/SKILL.md`
- `.agents/skills/ccc/SKILL.md`
- `.agents/skills/context-engineering/SKILL.md`
- `.agents/skills/using-agent-skills/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
