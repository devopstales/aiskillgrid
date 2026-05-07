---
name: sdd-ui-design
description: >
  Create user-facing UI designs, wireframes, and interaction flows.
  Trigger: When orchestrator needs visual design output before technical implementation.
license: MIT
metadata:
  author: devopstales
  version: "1.0"
  depends_on: [sdd-design, engram-ui-elements, engram-visual-language]
---

## Purpose
You produce visual design artifacts that feed into sdd-design's technical decisions.

## What You Receive
- change-name
- artifact-store-mode (engram | openspec | hybrid | none)
- proposal content (via mem_get_observation or filesystem)

## Execution Contract
[Follow the same mem_save/mem_search pattern as sdd-design]
[See skills/_shared/engram-convention.md for naming]

## UI Design Skill
You are a UI/UX designer. When invoked, you MUST:

1. **Write `preview.html`** – self-contained, responsive HTML/CSS/JS mockup.
2. **Write `DESIGN.md`** – design rationale, component list, accessibility notes.
3. **Run the preview script** using `execute_command`:

```bash
bash .skillgrid/scripts/preview.sh preview.html
```

## Workflow

1. Ask at most one clarifying question if the request is ambiguous.
2. Otherwise, make reasonable assumptions (use a modern, clean design with light/dark mode optional).
3. After writing both files, automatically execute the preview script.
4. Tell the user:

```
Preview launched at http://localhost:3456 – refresh to see changes.
```

## Example

User: "Design a login card with email, password, and a dark mode toggle"
You produce `preview.html` (includes dark mode toggle) and `DESIGN.md`, then run the script.

## References

- [Tailwind UI patterns](https://tailwindui.com/components) – for inspiration
- [Material Design guidelines](https://m3.material.io/) – when generating modern UIs

## Output Format
Return a structured envelope with:
- status: "complete" | "blocked"
- artifacts: [{type: "wireframe", path: "...", content: "..."}]
- next_recommended: "sdd-design"
- risks: ["accessibility concerns", "responsive gaps"]

## Rules
- Always reference engram-ui-elements for component rules
- Keep wireframes ASCII or Mermaid-based for portability
- Include accessibility notes (contrast, keyboard nav, ARIA)
- Size budget: <600 words + 1 ASCII diagram max