# Manifesto and principles

This document states how AISkillGrid thinks about **AI-assisted development**: not as a path to a lights-out “dark factory,” but as **harnessed** agentic work where **you** own the pipeline, the gates, and validation. **Spec-driven development** is the spine that keeps models aligned with human intent.

---

## Manifesto

**The goal is not autonomy for its own sake.** Some visions of agentic coding end in a fully automated factory: humans fade out, models ship on their own. We do not treat that as the end state. It compresses risk, hides mistakes, and confuses motion with progress.

**The goal is a controllable agentic pipeline.** The future we build for is **workflows you design**: phases, checkpoints, explicit **human-in-the-loop (HITL)** gates, and **user validation** at the points that matter. Agents accelerate execution; **you** retain authority over scope, quality bar, and when something is “done enough” to merge or release.

**Specs guide the AI.** Requirements and design artifacts are not paperwork—they are the steering signal. **Spec-driven development** (OpenSpec-style changes, PRDs, delta specs, tasks) gives the model something stable to implement against and something reviewers can verify against. Chat alone drifts; **files and specs anchor** the work.

**Harness engineering.** A model in isolation is not a process. A **harness** is commands, skills, rules, MCP tools, memory, indexing, handoff files, and dashboards that bound what agents do, when they stop, and what evidence they must produce. AISkillGrid is that harness: structured phases (`/skillgrid-*`), reusable skills, subagent discipline, and durable artifacts so work survives compaction and handoff.

**Humans stay in the loop.** Automation without validation is debt. Validation is not a failure mode—it is the mechanism that makes agentic pipelines trustworthy in real teams.

---

## Core concepts of AI-assisted development

These ideas show up throughout the hub documentation and commands.

| Concept | Meaning here |
|--------|----------------|
| **Human-in-the-loop (HITL)** | You approve boundaries, risky moves, and “done.” Commands like `/skillgrid-loop` and validate/finish steps assume you can stop or steer the process. |
| **User validation** | Explicit checks: spec alignment, review gates, test evidence, security passes—not silent merge-by-default. |
| **Spec-driven development** | Intent lives in specs and change artifacts before (and during) code. AI implements *toward* those specs; verification traces back to them. |
| **Agentic pipeline** | A **sequence of phases** (explore, plan, apply, test, validate, finish) with tools and skills attached—not one long autonomous chat. |
| **Harness** | The configured layer around the model: rules, skills, MCP, memory, indexing, handoff paths, UI—so behavior is repeatable and auditable. |
| **Artifacts over transcripts** | PRDs, OpenSpec changes, handoff markdown, logs, and checkpoints are the system of record; chat is ephemeral. |
| **Specialists, not a monoculture** | Subagent personas (review, security, test, research) provide **bounded** reports; the parent session synthesizes and decides. |
| **Local-first and portable** | State lives in the repo and local services where possible; no single vendor runtime is required to resume work. |

---

## How this maps to AISkillGrid

- **Commands** encode phases and gates; they are the skeleton of your pipeline.
- **Skills** encode *how* work in a phase should be done (TDD, review, GitNexus-first exploration, and so on).
- **`.configs/AGENTS.md`** and **rules** encode non-negotiables for agents working in a project.
- **Memory and indexing** reduce re-discovery; they do not replace your judgment on product tradeoffs.
- **Web UI and ticketing** make state visible so validation is informed, not blind.

If a single idea ties the table together: **AI assists; you validate; specs and harnesses keep both honest.**

---

## Further reading

- [00-start-here.md](00-start-here.md) — solution overview and reading order.
- [09-workflow-usage.md](09-workflow-usage.md) — operating the Skillgrid workflow day to day.
- [docs/workflow.md](workflow.md) — phases, config, and handoff (repository root).
- [03-skills.md](03-skills.md) — how skills package operating procedures for agents.
