# Slice Specification with Review Checklist

**Slice:** `<slice-slug>`
**Change:** `<change-id>`
**Spec version:** 1.0
**Status:** draft | approved | in-review | merged

---

## Goal

One-sentence description of user-visible outcome this slice delivers.

**Example:** "User can register with email/password and receive verification email."

## Scope

### In Scope

- ✅ Specific behaviors this slice implements
- ✅ APIs created or modified
- ✅ UI components (if frontend)
- ✅ Data schema changes
- ✅ Tests to be written

### Out of Scope

- ❌ Explicitly excluded related features (to prevent scope creep)
- ❌ Items deferred to later slices
- ❌ Non-functional requirements unless critical

## Acceptance Criteria

Each criterion must be testable and traceable to a task.

- [ ] AC-1: [specific, measurable condition]
- [ ] AC-2: [specific, measurable condition]
- [ ] AC-3: [specific, measurable condition]

## Technical Approach

Brief design decisions relevant to this slice (details in `design.md`).

- Pattern chosen: [controller/service/repository]
- Key dependencies: [libraries, external services]
- Data model: [entity relationships]
- Error handling strategy: [how errors surface]

## Implementation Tasks Reference

From `tasks.md`:

| Task ID | Description | TDD Phase | Files |
|---------|-------------|-----------|-------|
| 1.1 | Write failing test | RED | `tests/...` |
| 1.2 | Implement minimal code | GREEN | `src/...` |
| 1.3 | Refactor | REFACTOR | `src/...` |

---

## Review Checklist

For use during `sdd-verify` and `sdd-review`:

### Spec Compliance (`sdd-verify`)

- [ ] Every acceptance criterion has ≥1 task implementing it
- [ ] All tasks completed per tasks.md
- [ ] No out-of-scope features added without HITL approval
- [ ] Tests cover all acceptance criteria
- [ ] Edge cases from spec addressed

### Code Quality (`sdd-review`)

- [ ] Functions ≤ 50 lines
- [ ] Names descriptive
- [ ] No duplication
- [ ] Error handling explicit and logged
- [ ] Security: no secrets, input validated
- [ ] Performance: no O(n²) patterns
- [ ] Test quality: realistic assertions, not over-mocked

### Pre-Merge (`pre-merge-verification`)

- [ ] Tests fully green
- [ ] Lint/typecheck clean
- [ ] No debug leftovers
- [ ] Commit messages follow convention
- [ ] Branch mergeable with main

---

## Evidence

**Spec compliance matrix:**

| Requirement | Evidence | Status |
|-------------|----------|--------|
| AC-1 | `src/auth/register.ts:15-42` | ✅ |
| AC-2 | `tests/auth/register.test.ts:8-23` | ✅ |
| AC-3 | — | ❌ missing |

**Code quality findings:**

| Severity | Count | Examples |
|----------|-------|----------|
| CRITICAL | 0 | — |
| IMPORTANT | 1 | DRY violation in error handling (lines 45-50, 67-71) |
| MINOR | 2 | Variable naming inconsistencies |

---

## Sign-off

**Spec Compliance:** PASS / FAIL — [reviewer name] — [date]
**Code Quality Review:** APPROVED / CHANGES_REQUESTED — [reviewer] — [date]
**Verified By:** [SHA of commit that finalized slice]

---

**Template:** `.skillgrid/templates/template-slice-spec-with-review.md`
