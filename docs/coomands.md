
# Commands

Skills (appear in autocomplete):

- `/sdd-init` → initialize SDD context; detects stack, bootstraps persistence
- `/sdd-explore <topic>` → investigate an idea; reads codebase, compares approaches; no files created
- `/sdd-apply [change]` → implement tasks in batches; checks off items as it goes
- `/sdd-verify [change]` → validate implementation against specs; reports CRITICAL / WARNING / SUGGESTION
- `/sdd-archive [change]` → close a change and persist final state in the active artifact store
- `/sdd-ui-design <description>` → develop UI and UX direction, previews, design constraints, and design artifacts`

Meta-commands (type directly — orchestrator handles them, won't appear in autocomplete):

- `/sdd-brainstorm <change>` → start a new change by delegating exploration + proposal to sub-agents
