---
name: /opsx-onboard
id: opsx-onboard
category: Workflow
description: Guided onboarding - walk through a complete OpenSpec workflow cycle with narration
---

<objective>

You are executing **`/opsx-onboard`**, a thin OpenSpec command wrapper.

Run a guided OpenSpec onboarding flow. For full Skillgrid bootstrap use `/skillgrid-init`.

</objective>

<process>

## Required Skill

Read and follow `.agents/skills/openspec-onboard/SKILL.md`.

## Boundary

This command owns OpenSpec mechanics only. It does not create Skillgrid PRDs, update `.skillgrid/tasks/context_<change-id>.md`, manage Engram, create checkpoints, or perform tracker hygiene unless the referenced OpenSpec skill explicitly requires it.

For full Skillgrid behavior, use the matching `/skillgrid-*` command.

## Completion Report

Summarize the OpenSpec change, artifacts inspected or updated, commands run, blockers, and the recommended next OpenSpec or Skillgrid command.

</process>
