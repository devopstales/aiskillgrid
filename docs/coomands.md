
# Commands

Skills (appear in autocomplete):

- `/sdd-init` → initialize SDD context; detects stack, bootstraps persistence
- `/sdd-explore <topic>` → investigate an idea; reads codebase, compares approaches; no files created
- `/sdd-apply [change]` → implement tasks in batches; checks off items as it goes
- `/sdd-loop [change]` → run controlled AFK-safe build loop iterations with reassessment after each slice
- `/sdd-verify [change]` → validate implementation against specs; reports CRITICAL / WARNING / SUGGESTION
- `/sdd-board <decision>` → run specialist persona board and record accepted/rejected options with conflict gating
- `/sdd-archive [change]` → close a change and persist final state in the active artifact store
- `/sdd-design-ui <description>` → develop UI and UX direction, previews, design constraints, and design artifacts

Meta-commands (type directly — orchestrator handles them, won't appear in autocomplete):

- `/sdd-brainstorm <change>` → start a new change by delegating exploration + proposal to sub-agents
