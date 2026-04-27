---
description: Archive a completed change in the experimental workflow
---

<objective>

You are executing **`/opsx-archive`**, a thin OpenSpec command wrapper.

Archive an OpenSpec change. For PRD status, tracker hygiene, and checkpoint cleanup use `/skillgrid-finish`.

</objective>

<process>

## Required Skill

Read and follow `.agents/skills/openspec-archive-change/SKILL.md`.

## Boundary

This command owns OpenSpec mechanics only. It does not create Skillgrid PRDs, update `.skillgrid/tasks/context_<change-id>.md`, manage Engram, create checkpoints, or perform tracker hygiene unless the referenced OpenSpec skill explicitly requires it.

For full Skillgrid behavior, use the matching `/skillgrid-*` command.

## Completion Report

Summarize the OpenSpec change, artifacts inspected or updated, commands run, blockers, and the recommended next OpenSpec or Skillgrid command.

</process>
