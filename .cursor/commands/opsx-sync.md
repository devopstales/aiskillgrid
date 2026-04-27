---
name: /opsx-sync
id: opsx-sync
category: Workflow
description: Sync delta specs from a change to main specs
---

<objective>

You are executing **`/opsx-sync`**, a thin OpenSpec command wrapper.

Sync delta specs into main specs without archiving. For ship-phase closure use `/skillgrid-finish`.

</objective>

<process>

## Required Skill

Read and follow `.agents/skills/openspec-sync-specs/SKILL.md`.

## Boundary

This command owns OpenSpec mechanics only. It does not create Skillgrid PRDs, update `.skillgrid/tasks/context_<change-id>.md`, manage Engram, create checkpoints, or perform tracker hygiene unless the referenced OpenSpec skill explicitly requires it.

For full Skillgrid behavior, use the matching `/skillgrid-*` command.

## Completion Report

Summarize the OpenSpec change, artifacts inspected or updated, commands run, blockers, and the recommended next OpenSpec or Skillgrid command.

</process>
