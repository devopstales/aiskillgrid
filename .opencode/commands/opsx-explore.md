---
name: /opsx-explore
id: opsx-explore
category: Workflow
description: "Enter explore mode - think through ideas, investigate problems, clarify requirements"
---

<objective>

You are executing **`/opsx-explore`**, a thin OpenSpec command wrapper.

Explore an OpenSpec idea or existing change. For full Skillgrid PRD, handoff, and Engram behavior use `/skillgrid-explore`.

</objective>

<process>

## Required Skill

Read and follow `.agents/skills/openspec-explore/SKILL.md`.

## Boundary

This command owns OpenSpec mechanics only. It does not create Skillgrid PRDs, update `.skillgrid/tasks/context_<change-id>.md`, manage Engram, create checkpoints, or perform tracker hygiene unless the referenced OpenSpec skill explicitly requires it.

For full Skillgrid behavior, use the matching `/skillgrid-*` command.

## Completion Report

Summarize the OpenSpec change, artifacts inspected or updated, commands run, blockers, and the recommended next OpenSpec or Skillgrid command.

</process>
