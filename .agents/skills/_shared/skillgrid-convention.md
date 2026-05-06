Based on what you detected, create the skillgrid files by populating the templates:

- `.skillgrid/project/ARCHITECTURE.md`: Copy `.skillgrid/templates/template-architecture.md` and replace placeholders with detected architecture patterns and tech stack.

- `.skillgrid/project/PROJECT.md`: Copy `.skillgrid/templates/template-project.md` and replace placeholders with detected project conventions.

- `.skillgrid/project/STRUCTURE.md`: Copy `.skillgrid/templates/template-structure.md` and replace placeholders with detected project structure.

If the templates do not exist, print that they are missing and skip file creation, but still return the detected context in the summary.

Refresh `.skillgrid/project/*` and `DESIGN.md` only for stable discovered facts.