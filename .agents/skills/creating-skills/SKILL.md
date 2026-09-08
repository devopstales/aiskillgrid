---
name: creating-skills
description: Create or improve an agent skill (SKILL.md plus optional scripts/references) from real expertise. Use when the user asks to make, write, scaffold, refine, or audit a skill; escape skill hell; apply the writing-great-skills checklist (trigger, structure, steering, pruning); or validate skill structure.
license: MIT
metadata:
  author: devopstales
  version: "1.1"
  part-of: skillgrid
---

# Creating skills

A skill is a directory with a `SKILL.md` (YAML frontmatter + Markdown), optionally plus `scripts/`, `references/`, and `assets/`. Ground every skill in real expertise — steps that worked, corrections needed, project facts. General knowledge alone produces vague, worthless skills.

Rubric for great skills (trigger → structure → steering → pruning). Use it when writing **or** auditing. Details and examples: [`references/checklist.md`](references/checklist.md). Frontmatter rules: [`references/anatomy.md`](references/anatomy.md). Content patterns: [`references/patterns.md`](references/patterns.md).

## Workflow

Progress:
- [ ] Step 1: Gather expertise
- [ ] Step 2: Decide the trigger (user-invoked vs model-invoked)
- [ ] Step 3: Name, scope, and location
- [ ] Step 4: Structure as steps + reference; keep SKILL.md minimal
- [ ] Step 5: Steer — leading words and enough leg work per step
- [ ] Step 6: Write frontmatter and body
- [ ] Step 7: Prune — DRY, sediment, no-ops, deletion tests
- [ ] Step 8: Validate — `scripts/validate_skill.sh <skill-dir>`
- [ ] Step 9: Test on a real task; fold every correction into Gotchas

### Step 1: Gather expertise

Extract from a recent or repeating task:

- Steps that worked (successful sequence)
- Corrections ("use X, not Y", "always check Z first")
- Concrete formats in and out
- Project-specific facts, conventions, edge cases

Otherwise synthesize from runbooks, config, review comments, fix commits, failure reports. Ask which material to draw from — do not invent.

### Step 2: Decide the trigger

| Mode | Mechanism | Cost | Prefer when |
|------|-----------|------|-------------|
| **Model-invoked** | `description` stays in agent context as a context pointer | Context load every request + unpredictability (agent may skip) | Agent (or another skill) must discover/fire it |
| **User-invoked** | `disable-model-invocation: true` — description is human-facing only | Cognitive load on the user (must remember to invoke) | You want deterministic control; avoid evals for "did it fire?" |

Model-invoked is not "better" — it is more flexible and more expensive. Prefer user-invoked when predictability matters more than autonomous discovery. See checklist for harness notes.

### Step 3: Name, scope, and location

One coherent unit of work — like a function. Query + format results: one skill. Query + database admin: two. Prefer a purpose word in the name: `creating-skills`, `deploy-staging`.

Default: project's `.agents/skills/` (or `~/.agents/skills/` if it should travel). Match naming/structure of sibling skills first.

### Step 4: Structure — steps + reference; minimal SKILL.md

Compose from two units:

1. **Steps** — the procedure the agent walks
2. **Reference** — supporting material those steps need (definitions, templates)

Skills may be steps-only, reference-only, or both. Keep **SKILL.md as small as possible** (maintainability + tokens).

**Branches:** if reference is only needed on one path, hide it behind a **context pointer** ("If updating the glossary, read the matching file under `references/`"). Always-needed reference for a single-branch skill can stay in SKILL.md. Bundled code → `scripts/` (run it); templates/data → `assets/`. References one level deep from SKILL.md.

### Step 5: Steer — leading words and leg work

**Leading words** pack meaning into a short phrase the agent will echo in thinking and output (e.g. `vertical slice`, not a long "don't code layer by layer…" essay). Repeat the phrase consistently through the skill. If the agent ignores you, strengthen or replace the leading word — watch reasoning traces for adoption.

**Leg work:** when a step under-invests because a later goal is visible (classic: skim clarifying questions, rush the plan), split into a focused skill so the agent only sees the current phase. Not always required — use when you need extra depth on one step.

### Step 6: Write frontmatter and body

Minimal frontmatter (full rules in anatomy):

```markdown
---
name: <dir-name>
description: <what it does> + "Use when <triggers and keywords>."
# user-invoked only:
# disable-model-invocation: true
---
```

Body rules of thumb (patterns for templates/checklists):

- Add what the agent lacks; omit what it knows
- Default + escape hatch, not a menu of equal options
- Procedures over one-off answers; calibrate control per section
- Gotchas = concrete corrections; working example for non-obvious formats
- Target a tight SKILL.md — move branch-only detail out

**Description = discovery surface (SDO).** For model-invoked skills, the `description` is the *only* thing the agent sees when deciding whether to load the skill. Write it for triggering, not summarizing:

- Start with "Use when …" — describe the *triggering conditions*, not the workflow
- Include the keywords an agent would use when facing the problem ("bug", "test failure", "refactor")
- Do NOT summarize the workflow in the description — the agent may follow the summary instead of reading the skill
- Do NOT use first person ("I will…") or vague abstractions ("helps with development")
- Technology-specific skills: name the technology explicitly ("Use when writing Playwright E2E tests")
- Token budget: aim for <200 chars. Every character is in context on every request.

```
❌ "A skill for writing great tests using TDD principles"
✅ "Use when implementing any feature or bugfix, before writing implementation code"
```

### Step 7: Prune

Before shipping, run a pruning pass:

- **DRY / single source of truth** — no duplicated reference or restated steps
- **Sediment** — stale or drive-by additions; delete or move to the right branch
- **No-ops** — instructions the agent would follow anyway; deletion-test each paragraph
- **Massive body** — usually a symptom of the above, not a goal

### Step 8: Validate

Run `scripts/validate_skill.sh <path-to-skill-dir>`. Fix every failure; re-run until clean.

### Step 9: Test and iterate (TDD for skills)

**Writing skills IS TDD applied to process documentation.** If you didn't watch an agent fail without the skill, you don't know if the skill teaches the right thing.

**RED — baseline (before writing the skill):**
1. Give a fresh subagent the task the skill will govern, *without* the skill loaded.
2. Read the actual execution and reasoning traces — not just the final output.
3. Document the exact rationalizations, skips, and mistakes the agent made. These are the failure modes the skill must close.

**GREEN — write the minimal skill:**
4. Write the smallest SKILL.md that would have prevented each documented failure.
5. Give a fresh subagent the same task *with* the skill loaded.
6. Read the traces: did the agent comply? If not, the skill's wording is wrong — fix it and re-test.

**REFACTOR — close loopholes:**
7. Micro-test wording: change one sentence, re-run the subagent, see if compliance holds. A rule that survives paraphrase is a rule that works.
8. Cut instructions the agent followed without help (deletion-test each paragraph).
9. Add every remaining mistake to the red-flags table.

**Match the form to the failure:**

| Skill type | Test focus |
|---|---|
| **Discipline** (rules, requirements) | Does the agent still rationalize past the rule? Close every loophole explicitly. |
| **Technique** (how-to guide) | Does the agent produce the correct output format? |
| **Pattern** (mental model) | Does the agent apply the pattern in a *new* context, not just the example? |
| **Reference** (docs, API) | Does the agent find and use the right section without being pointed? |

**Bulletproofing discipline skills against rationalization:**
- Close every loophole explicitly ("no exceptions," "not just this once")
- Address "spirit vs letter" arguments ("violating the letter IS violating the spirit")
- Build a red-flags table — every rationalization the agent actually used becomes a row
- Update the red-flags table for new violation symptoms found in testing

One execute→revise loop helps a lot; hard domains need more.

## Red flags

| Thought | Reality |
|---|---|
| "I'll write the skill, then test it on a real task" | That's GREEN without RED. Run the baseline first — watch the agent fail *without* the skill, then write the minimal skill that fixes the failure. |
| "The description can summarize the workflow" | The agent may follow the summary instead of reading the skill. Description = triggering conditions, not workflow. |
| "I'll name it `My-Skill` with a capital" | `name` must match the directory exactly. Uppercase is invalid. |
| "I'll add more detail to the description so it triggers better" | Every character is in context on every request. <200 chars. More keywords, less prose. |
| "The agent already handles this well, but a skill would be nice" | Context cost with no gain. Don't create a skill the agent already handles. |
| "I'll leave the branch-only template in SKILL.md for convenience" | It bloats every invocation. Context-pointer it out to `references/`. |
| "I'll skip the red-flags table, the rules are clear" | Clear to you. The agent will find the loophole. Every rationalization it actually uses becomes a row. |
| "The script will work, the agent knows the setup" | The agent won't know your local setup. Scripts must be self-contained — document every dependency. |
| "I'll add a 'similar to the other skill' reference" | Cross-reference the other skill by name. Don't restate its content. |
| "I'll commit the example with the real API key, it's just a sample" | Never commit secrets into examples or references. |
