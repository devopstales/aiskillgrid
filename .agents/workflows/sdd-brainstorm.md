---
description: Start a new SDD change — runs exploration then creates a proposal
agent: sdd-orchestrator
---

Follow the SDD orchestrator workflow for starting a new change named "$ARGUMENTS".

WORKFLOW:
1. Launch sdd-explore sub-agent to investigate the codebase for this change
2. Present the exploration summary to the user
3. Launch sdd-propose sub-agent to create a proposal based on the exploration
4. Present the proposal summary and ask the user if they want to continue with specs and design

Run these sub-agents in sequence:
1. sdd-propose — create the proposal
2. sdd-spec — write specifications
3. sdd-design — create technical design
4. sdd-tasks — break down into implementation tasks

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Change name: $ARGUMENTS
- Artifact store mode: hybrid

ENGRAM NOTE:
Sub-agents handle persistence automatically. Each phase saves its artifact to engram with topic_key "sdd/$ARGUMENTS/{type}".

Read the orchestrator instructions to coordinate this workflow. Do NOT execute phase work inline — delegate to sub-agents.