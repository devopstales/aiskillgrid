---
description: Start a new SDD change — runs exploration then creates a proposal
agent: sdd-orchestrator
---

Follow the SDD orchestrator workflow for starting a new change named "$ARGUMENTS".

WORKFLOW:
1. Launch sdd-explore sub-agent to investigate the codebase for this change
2. Present the exploration summary to the user
3. Launch sdd-clarify sub-agent to refine understanding:
   - It will ask questions to resolve ambiguity about goals, scope, and constraints.
   - It captures the mutual understanding in a clarifications.md artifact.
4. Present the clarifications summary and ask the user to confirm or iterate.
5. Once mutual understanding is confirmed, launch sdd-propose sub-agent
6. Present the proposal summary and ask the user if they want to continue with specs and design

Run these sub-agents in sequence after clarification:
1. sdd-propose — create the proposal (using the clarifications)
2. sdd-spec — write specifications
3. sdd-design — create technical design
4. sdd-prd — consolidate into PRD (.skillgrid/prd/PRD<NN>_$ARGUMENTS.md)
5. sdd-tasks — break down into vertically sliced implementation tasks

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Change name: $ARGUMENTS
- Artifact store mode: hybrid

ENGRAM NOTE:
Sub-agents handle persistence automatically. Each phase saves its artifact to engram with topic_key "sdd/$ARGUMENTS/{type}".
FILESYSTEM PERSISTENCE:
  Reade .agents/skills/_shared/skillgrid-handoff.md for filesystem persistence instructions.

Read the orchestrator instructions to coordinate this workflow. Do NOT execute phase work inline — delegate to sub-agents.