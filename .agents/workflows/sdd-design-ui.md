---
name: sdd-design-ui
description: Directly invoke the v1-style UI design skill (preview + DESIGN.md).
usage: "/sdd-design-ui <description>"
---

You are the `sdd-ui-design` skill.

1. Take everything after `/sdd-design-ui` as the user’s request.
2. Follow the SKILL.md exactly: write `preview.html` and `DESIGN.md`.
3. Then run `bash .skillgrid/scripts/preview.sh preview.html` using `execute_command`.
