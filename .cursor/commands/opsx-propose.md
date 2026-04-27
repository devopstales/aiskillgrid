---
name: /opsx-propose
id: opsx-propose
category: Workflow
description: Propose a new change - create it and generate all artifacts in one step
---

<objective>

You are executing **`/opsx-propose`**, a thin OpenSpec command wrapper.

Create OpenSpec proposal/design/spec/tasks artifacts. For Skillgrid PRD-first planning use `/skillgrid-plan`.

</objective>

<process>

## Required Skill

Read and follow `.agents/skills/openspec-propose/SKILL.md`.

## Boundary

This command owns OpenSpec mechanics only. It does not create Skillgrid PRDs, update `.skillgrid/tasks/context_<change-id>.md`, manage Engram, create checkpoints, or perform tracker hygiene unless the referenced OpenSpec skill explicitly requires it.

For full Skillgrid behavior, use the matching `/skillgrid-*` command.

## Completion Report

Summarize the OpenSpec change, artifacts inspected or updated, commands run, blockers, and the recommended next OpenSpec or Skillgrid command.

</process>
