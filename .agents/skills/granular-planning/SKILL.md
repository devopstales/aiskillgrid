---
name: granular-planning
description: >
  Creates highly-detailed implementation plans with 2–5 minute atomic tasks.
  Each task includes exact file paths, complete code blocks, verification commands,
  and no placeholders. Requires fresh-agent context clarity.
  Trigger: After sdd-brainstorm/sdd-design, before sdd-apply.
license: Apache-2.0
metadata:
  author: aiskillgrid-integration
  version: "1.0"
  dependencies:
    - "enforced-tdd-protocol"
  triggers:
    - "planning_required"
    - "sdd-apply_prepare"
---

# Granular Planning

## Overview

Break implementation scope into atomic tasks where each task represents 2–5 minutes of work for a fresh agent. Tasks are so explicit that an engineer with zero project context can execute them without additional research.

**Core principle:** If a task needs more than 5 minutes, it's too big. Split it.

## When to Use

Invoked automatically:
- After `sdd-brainstorm` produces approved spec
- After `sdd-design` completes design artifacts
- When `sdd-apply` detects missing or insufficient tasks
- Before any implementation begins

## Output Location

Plans saved to (in order of preference):
1. User-specified location (passed as parameter)
2. `.skillgrid/plans/YYYY-MM-DD-<change-title>-plan.md`
3. `openspec/changes/<change-id>/tasks.md`

State entry (for tracking):
```
engram: sdd/<change-id>/granular-plan
file:    openspec/changes/<change-id>/tasks.md
```

## Plan Structure

```markdown
# [Feature Name] — Granular Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use sequential-agent-executor (preferred) or batch-executor to implement this plan task-by-task.

**Goal:** [One-sentence description of delivered capability]

**Architecture:** [2-3 sentence approach summary]

**Tech Stack:** [Key technologies, frameworks, libraries used]

**Change ID:** <change-id>
**Plan ID:** granular-plan-<change-id>-<YYYY-MM-DD>
**Total estimated tasks:** <N>

---

## File Structure Map

Before tasks, document which files exist and where:

```
CREATE:
  src/services/user-authentication.service.ts
  tests/unit/services/user-authentication.service.test.ts
  docs/authentication-flow.md

MODIFY:
  src/index.ts:1-15              (add import)
  package.json:12-20             (add devDependency)
  .env.example:8-12              (add JWT_SECRET)
```

## Task List

Requirements per task:

| Field | Content |
|-------|---------|
| ID | Numeric sequence (1., 1.1, 1.2, 2.1… ) |
| Title | Actionable verb phrase (≤ 8 words) |
| Files | Exact paths with line ranges if modifying |
| Code | Complete code blocks (no TODOs) |
| Command | Exact shell command to run |
| Expected | Expected command output |
| TDD | RED / GREEN / REFACTOR phase indicator |
| Duration | 2-5 min estimate |

### Task Template

```markdown
### Task 1.1: [Action-oriented title]

**Files:**
- Create: `src/auth/middleware.ts`
- Modify: `src/index.ts:1-10` (add import)

**TDD Phase:** RED

**Step:**
- [ ] Write the failing test

```typescript
// tests/auth/middleware.test.ts
describe('AuthMiddleware', () => {
  test('rejects requests without Authorization header', () => {
    const response = await request(app).get('/api/profile');
    expect(response.status).toBe(401);
    expect(response.body.error).toBe('Missing credentials');
  });
});
```

**Command:**
```bash
npm test -- tests/auth/middleware.test.ts
```

**Expected Output:**
```
FAIL tests/auth/middleware.test.ts
  × rejects requests without Authorization header
    Expected status 401 but received 404 (function not defined)
```

**Atomic duration:** 3 minutes

---

### Task 1.2: [Implement minimal pass]

**Files:**
- Create: `src/auth/middleware.ts` (first implementation)

**TDD Phase:** GREEN

**Step:**
- [ ] Write minimal implementation

```typescript
// src/auth/middleware.ts
import { Request, Response, NextFunction } from 'express';

export function authMiddleware(req: Request, res: Response, next: NextFunction) {
  const auth = req.headers['authorization'];
  if (!auth) {
    return res.status(401).json({ error: 'Missing credentials' });
  }
  next();
}
```

**Command:**
```bash
npm test -- tests/auth/middleware.test.ts
```

**Expected Output:**
```
PASS tests/auth/middleware.test.ts
Test Suites: 1 passed, 1 total
```

**Atomic duration:** 4 minutes

---

### Task 1.3: [Refactor cleanly]

**Files to modify:**
- `src/auth/middleware.ts` (extract token parsing)

**TDD Phase:** REFACTOR

**Step:**
- [ ] Extract token parsing helper

```typescript
// src/auth/helpers.ts
export function extractToken(authHeader: string): string | null {
  const match = authHeader.match(/^Bearer (.+)$/);
  return match ? match[1] : null;
}

// src/auth/middleware.ts
import { extractToken } from './helpers';
// ... use extractToken
```

**Verification:**
```bash
npm test -- tests/auth/middleware.test.ts && npm test -- tests/auth/helpers.test.ts
```

**Expected:** All tests pass. No behavioral changes.

**Atomic duration:** 3 minutes
```

## Plan Quality Checklist

Before finalizing plan, self-review:

**Coverage:** Every spec requirement maps to ≥1 task
**Completeness:** No "TODO", "TBD", "implement later" placeholders
**Precision:** Every file path exact, line range provided for modifications
**Executability:** Engineer can follow steps without external research
**TDD compliance:** Every code task precedes or accompanies test task
**Atomicity:** No task exceeds 5 minutes or touches >3 files
**Traceability:** Task ID references slice ID and spec requirement ID

## Granularity Guidelines

**Task = 2–5 minutes** means:
- Write one test function (≤ 20 lines)
- Implement one function (≤ 15 lines)
- Add one config key + read it
- Create one file (≤ 50 lines)
- Run one command and verify result

**Split if:**
- Task mentions "and" (e.g., "middleware and error handler") → split
- Task spans multiple directories → split
- Task requires "designing" → split into decisions first, then implementation
- Task duration estimate >5 min → split

**Good task examples:**
- ✅ "Define User interface in `src/types/index.ts`"
- ✅ "Add `validateEmail()` helper function"
- ✅ "Write test: POST /login rejects invalid password"
- ✅ "Register authMiddleware in app.ts"

**Bad task examples:**
- ❌ "Implement authentication" (too vague)
- ❌ "Add auth throughout project" (scope creep)
- ❌ "Set up database schema" (could be many files)
- ❌ "Create login page" (UI typically needs multiple files)

## No-Placeholders Rule

**You must write actual code in every task that produces code.** Never:
- "Implement the actual logic here"
- "Add appropriate error handling"
- "Fill in the rest based on your design"
- "See implementation"
- "Similar to Task X" (code must be repeated)

Every code block must be complete and correct as written.

## AFK/HITL Tagging

Each task inherits slice-level label, refined for the task:
```
[Slice: <slice-name>] [Label: AFK|HITL] [Reason: <why>] [Budget: safe|RISK]
```

AFK tasks: fully deterministic, no human decisions required
HITL tasks: first task in slice if slice requires human decision point

## Vertical Slice Prerequisites

If the spec has multiple slices, plan them separately:
- Plan Slice 1 completely (all tasks, granular)
- Then Slice 2 (which may depend on Slice 1 infrastructure)
- Document dependencies between slices

**Never mix slices in a single flat task list.**

## Rollback and Validation Risk

For each task, consider:
- What if this breaks? How to revert?
- Does it leave inconsistent state if incomplete?
- Is there a verification command that proves correctness?

Document high-risk tasks with rollback steps:
```markdown
**Rollback:** git revert HEAD if tests fail after completion
**Validation risk:** Medium — modifies shared config file
```

## Completion Criteria

Plan is complete when:
- Every acceptance criterion from slice spec has ≥1 corresponding task
- Every task is executable without research
- Total task count matches estimated slice size
- No placeholders remain
- All file paths are validated to exist or will be created
- TDD sequence explicit for every implementation task

## Invocation

Orchestrator calls:
```
/skill: granular-planning
Change: <change-id>
Spec: openspec/changes/<change-id>/specs/<slice>/spec.md
Design: openspec/changes/<change-id>/design.md
Tasks: openspec/changes/<change-id>/tasks.md
```

The skill reads inputs, produces `tasks.md`, performs self-review, then returns.

## See Also

- Source inspiration: `superpowers:writing-plans` (adapted)
- Existing: `sdd-tasks` skill (uses this planning method)
- Templates: `.skillgrid/templates/template-granular-task.md`
