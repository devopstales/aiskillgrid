# Diagnostic Report — <issue-id>

**Issue:** [short description]
**Change:** [change-id if applicable]
**Diagnostician:** [agent/human]
**Timestamp:** [ISO datetime]

---

## Phase 1 — Root Cause Investigation

### Reproduction

**Steps to reproduce:**
1.
2.
3.

**Reproduction rate:** always / sometimes (<X% flake) / intermittent

**Expected behavior:** [what should happen]
**Actual behavior:** [what actually happens]

**Evidence captured:**
- Logs: [path]
- Stack trace: [paste]
- Screenshot/recording: [link if applicable]

### Recent Changes Check

```bash
$ git log --oneline -10
$ git diff origin/main...HEAD --name-only
```

**Relevant recent changes:**
- [file] — [commit SHA] — [why relevant]

### Multi-Component Instrumentation (if applicable)

For each boundary, logged state:

| Layer | Entry state | Exit state | Pass/Fail |
|-------|-------------|------------|-----------|
| API gateway | headers: {auth: Bearer…} | — | ✅ |
| Auth service | token decoded: {sub: 123} | user: null | ❌ |

**Finding:** Failure occurs at [specific boundary].

### Data Flow Trace (for deep errors)

Traced backward from symptom:
```
symptom: userService.get() returns null
  ← called by orderService.createOrder() line 120
    ← user param passed from req.user set by auth middleware line 45
      ← auth middleware: token expired → req.user unset (correct behavior)
Root cause: Client using expired token, code correctly returns null. Not a bug.
```

---

## Phase 2 — Pattern Analysis

### Working Examples Found

| Example | Location | Why it works |
|---------|----------|--------------|
| Email validation | src/validators/email.ts | Uses regex with boundaries |
| Password hashing | src/auth/password.ts | bcrypt with 12 rounds |

### Differences Identified

| Aspect | Working | Broken | Gap |
|--------|---------|--------|-----|
| Input sanitization | `email.trim()` | missing | Add trim |
| Error type thrown | `ValidationError` | `Error` | Use specific type |

### Dependencies & Assumptions

What this component needs:
- Config: `JWT_SECRET` env var
- External: Redis for session store
- Precondition: user must be authenticated

**Missing/incorrect assumption:** [if any]

---

## Phase 3 — Hypothesis & Testing

### Ranked Hypotheses

1. **H1: Token expiration not checked** (confidence: HIGH)
   - If true: adding expiry check → null prevented
   - Test: mock expired token, verify 401 response

2. **H2: Redis session store evicts early** (confidence: MEDIUM)
   - If true: check Redis TTL configuration
   - Test: inspect session TTL directly

3. **H3: Race condition in concurrent requests** (confidence: LOW)
   - If true: reproducing with sequential requests → bug disappears
   - Test: run rep Nougat 100× sequential → if no failures, not race

**User input:** [user re-ranking or ruled-out hypotheses]

### Instrumentation Results

| Hypothesis | Test performed | Result | Succeeds/Fails prediction |
|------------|----------------|--------|--------------------------|
| H1 | Mock expired token, call endpoint | Returns 401 | ✅ Confirmed |
| H2 | Check Redis TTL | TTL = 86400 (correct) | ❌ Ruled out |
| H3 | Reproduce 100× sequentially | 0/100 failures | ❌ Ruled out |

---

## Phase 4 — Implementation

### Regression Test Created

```typescript
// tests/auth/middleware.test.ts
test('rejects request with expired token', async () => {
  const expiredToken = generateToken({ sub: 123, exp: Math.floor(Date.now() / 1000) - 60 });
  const response = await request(app).get('/api/profile').set('Authorization', `Bearer ${expiredToken}`);
  expect(response.status).toBe(401);
});
```

**Test status:** RED → GREEN confirmed.

### Fix Applied

```diff
--- a/src/auth/middleware.ts
+++ b/src/auth/middleware.ts
@@
- const decoded = jwt.verify(token, secret);
+ const decoded = jwt.verify(token, secret);
+ if (decoded.exp < Math.floor(Date.now() / 1000)) {
+   throw new UnauthorizedError('Token expired');
+ }
```

**Commit:** `abc123` — "fix: reject expired tokens (refs #456)"

### Verification

**Original repro:** No longer reproduces.
**Regression test:** PASS.
**Full test suite:** 127/127 PASS.

---

## Post-Mortem

**What would have prevented this bug?**
- Token expiry check was missing from middleware spec → update PRD acceptance criteria
- No test covered expired token case → add to standard auth test matrix

**Architecture friction:**
- None — test seam available and effective

**Recommendations:**
1. Update CONTEXT.md: "Access tokens must be validated for expiry"
2. Add "expired token rejection" to auth middleware spec
3. Add to `sdd-tasks` template: always include negative-path auth tests

---

## Escalation (if any)

**Escalated to:** `sdd-architecture-review` on 2025-05-14 — recommended test-seam audit for auth flows.

**HITL required:** [Yes/No] — [reason]

---

**Status:** COMPLETED — root cause fixed and verified
**Follow-up:** [next action if any]
