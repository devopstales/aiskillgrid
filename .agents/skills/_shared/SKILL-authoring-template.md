# SKILL.md authoring template

Copy into `.agents/skills/<skill-name>/SKILL.md` (or your host’s skills root). Replace angle-bracket placeholders. Keep one **primary job** per file.

---

```yaml
---
name: <skill-id-matches-folder>
description: >
  <One paragraph: what the skill does. Include triggers as natural language —
  e.g. "Use when …" / "Trigger: …". List 2–4 example user phrases.>
license: Apache-2.0   # or project default
metadata:
  author: <org-or-you>
  version: "1.0"
---
```

## Triggers

- When …
- Example user prompts: "…", "…"

## Out of scope (do not do here)

- …
- Defer to skill `<other-skill>` for …

## Workflow

1. …
2. …

## Tools and evidence

- Prefer: …
- Avoid: …

## Stop conditions

- Stop when …
- Escalate or hand off when …

## Example prompt (user message)

```text
<Short realistic user request that should load this skill>
```

## Optional sections

- **Related skills** — bullets with backtick names.
- **Commands** — only if a fixed shell entrypoint exists.

---

**Compactness:** Agents load full `SKILL.md`; keep under ~200 lines unless the domain forces more. Put long examples in linked docs under `docs/` or `templates/`.
