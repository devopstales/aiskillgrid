# Granular Task Template

Use this structure for each task in the tasks file. Every task must be atomic (2–5 min) and complete.

---

### Task X.Y: [Action-oriented title, ≤8 words]

**Files:**
- Create: `path/to/file.ext`
- Modify: `path/to/existing.ext:line-start-line-end`  (or full file if small)
- Delete: `path/to/remove.ext` (if applicable)

**TDD Phase:** RED | GREEN | REFACTOR

**Step:**
- [ ] One atomic action (write test / implement / refactor)

```language
[Complete code block exactly as it should appear]
[No placeholders, TODOs, or "fill in later"]
```

**Command:**
```bash
[Exact command to run — include flags, paths]
```

**Expected Output:**
```
[Exact expected output — copy from successful run]
```

**Atomic duration:** 2-5 minutes estimate

**Rollback:** (optional — how to undo if something breaks)
```bash
git revert HEAD  # or specific instruction
```

**Validation risk:** low | medium | high  (high = affects shared resource, run with care)

---

### Task X.Y+1: [Next atomic step]

[Same structure]

---

## Quality Gate Checklist (for sdd-tasks validator)

- [ ] Every task has exact file paths
- [ ] No "TBD", "TODO", "later", "appropriate", "proper", etc.
- [ ] All code blocks complete and executable
- [ ] Commands have expected output
- [ ] TDD phase labeled if project uses TDD
- [ ] Duration per task ≤5 minutes estimated
- [ ] Task touches ≤3 files
- [ ] No "and" in title (split if so)
- [ ] All tasks under vertical slice headings, not flat list

---

## Example Correct Task

### Task 1.1: Write failing test for auth middleware

**Files:**
- Create: `tests/auth/middleware.test.ts`

**TDD Phase:** RED

**Step:**
- [ ] Write failing test

```typescript
import { request } from '@iyio/test-utils';
import { app } from '../src/index';

describe('AuthMiddleware', () => {
  test('rejects request without Authorization header', async () => {
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
  × rejects request without Authorization header
    Expected status 401 but received 404 (handler not defined)
```

**Atomic duration:** 3 minutes

---

## Antipatterns to Avoid

❌ **BAD — No file paths:**
```
### Task 1: Add authentication
- Write code for login
```

❌ **BAD — Placeholder code:**
```typescript
function validateToken(token: string) {
  // TODO: implement validation
}
```

❌ **BAD — "and" combining multiple things:**
```
Create middleware and error handlers and tests
```

❌ **BAD — Vague verification:**
```
Make sure it works
```

❌ **BAD — No expected output:**
```
Command: npm test
```

✅ **GOOD — Specific, file-located, complete code, expected output, single action**
