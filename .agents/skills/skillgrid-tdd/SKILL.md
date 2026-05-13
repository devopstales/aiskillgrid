---
name: skillgrid-tdd
description: >
  Test-Driven Development loop – write a failing test, make it pass, refactor. Use when user asks to implement something test‑first or provides a --tdd flag.
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
triggers:
  - "tdd"
  - "test-first"
  - "red green refactor"
  - "write a test first"
tools:
  - file_system
  - execute_command
---

# TDD Skill – Red · Green · Refactor

You are a TDD practitioner. Follow this strict loop.

## Philosophy

**Core principle**: Tests should verify behavior through public interfaces, not implementation details. Code can change entirely; tests shouldn't.

**Good tests** are integration-style: they exercise real code paths through public APIs. They describe _what_ the system does, not _how_ it does it. A good test reads like a specification — "user can checkout with valid cart" tells you exactly what capability exists. These tests survive refactors because they don't care about internal structure.

**Bad tests** are coupled to implementation. They mock internal collaborators, test private methods, or verify through external means. The warning sign: your test breaks when you refactor, but behavior hasn't changed.

## Anti-Pattern: Horizontal Slices

**DO NOT write all tests first, then all implementation.** This is "horizontal slicing" — treating RED as "write all tests" and GREEN as "write all code."

This produces **crap tests**:

- Tests written in bulk test _imagined_ behavior, not _actual_ behavior
- You end up testing the _shape_ of things (data structures, function signatures) rather than user-facing behavior
- Tests become insensitive to real changes — they pass when behavior breaks, fail when behavior is fine
- You outrun your headlights, committing to test structure before understanding the implementation

**Correct approach**: Vertical slices via tracer bullets. One test → one implementation → repeat. Each test responds to what you learned from the previous cycle. Because you just wrote the code, you know exactly what behavior matters and how to verify it.

```
WRONG (horizontal):
  RED:   test1, test2, test3, test4, test5
  GREEN: impl1, impl2, impl3, impl4, impl5

RIGHT (vertical):
  RED→GREEN: test1→impl1
  RED→GREEN: test2→impl2
  RED→GREEN: test3→impl3
  ...
```

## When to Use

- Implementing new features
- Adding new methods to existing services/repositories
- Fixing bugs (write failing test first)
- Refactoring existing code

## Core TDD Cycle

```
RED → GREEN → REFACTOR
 ↓      ↓         ↓
Fail   Pass    Improve
```

## Critical Rules

### Rule 1: ALWAYS Identify Entity Under Test

Before writing any test, state explicitly which entity is being tested:
- `TransactionService`
- `OFXParser`
- `BudgetRepository`

### Rule 2: Test File Location

Place test files next to the code being tested:
```
internal/service/
├── transaction_service.go
└── transaction_service_test.go  ← Test file here
```

### Rule 3: Test Naming Convention

Pattern: `Test<EntityName>_<MethodName>_<Scenario>`

```go
// ✅ GOOD
func TestTransactionService_ImportFromOFX_WithDuplicateFITID(t *testing.T)
func TestOFXParser_Parse_WithInvalidXML(t *testing.T)

// ❌ BAD
func TestImport(t *testing.T)
func TestValidData(t *testing.T)
```

### Rule 4: Test Main Cases (Not Everything)

**Test:**
- Happy path (1 test)
- Common errors (2-3 tests)
- Edge cases (1-2 tests)
- Business rules (as needed)

**Don't test:**
- Language features
- Third-party library internals
- Trivial getters/setters

## Workflow

### 1. Planning

When exploring the codebase, use the project's domain glossary so that test names and interface vocabulary match the project's language, and respect ADRs in the area you're touching.

Before writing any code:

- [ ] Confirm with user what interface changes are needed
- [ ] Confirm with user which behaviors to test (prioritize)
- [ ] Identify opportunities for deep modules (small interface, deep implementation)
- [ ] Design interfaces for testability
- [ ] List the behaviors to test (not implementation steps)
- [ ] Get user approval on the plan

Ask: "What should the public interface look like? Which behaviors are most important to test?"

**You can't test everything.** Confirm with the user exactly which behaviors matter most. Focus testing effort on critical paths and complex logic, not every possible edge case.

### 2. RED – Write a failing test

Based on the user's request, write **one** minimal test that describes the desired behaviour.
The test **must** fail before you write any production code.
Output: *"RED – Test written. Run `[test command]` – it fails as expected."*

### 3. GREEN – Make it pass

Write the simplest possible production code that makes the test pass.
No refactoring. No extra features.
Output: *"GREEN – Test now passes."*

### 4. REFACTOR – Clean up

Improve code (remove duplication, clarify names) without changing behaviour.
Re‑run test after each change.
Output: *"REFACTOR – Completed. Tests still pass."*

### 5. Repeat or confirm

Ask: *"Another test, or stop?"*

## Test Structure (AAA Pattern)

```go
func TestTransactionService_ImportFromOFX_Success(t *testing.T) {
    // ARRANGE - Setup test data and dependencies
    mockRepo := &mockTransactionRepository{}
    service := NewTransactionService(mockRepo)

    // ACT - Execute the method being tested
    result, err := service.ImportFromOFX(context.Background(), ofxData, accountID)

    // ASSERT - Verify expectations
    require.NoError(t, err)
    assert.Equal(t, 5, result.Inserted)
}
```

## Mocking Guidelines

Mock external dependencies (HTTP clients, databases, other services):

```go
type mockTransactionRepository struct {
    mock.Mock
}

func (m *mockTransactionRepository) BulkInsert(ctx context.Context, txs []*Transaction) error {
    args := m.Called(ctx, txs)
    return args.Error(0)
}
```

## Behavioural rules
- Never write production code before a failing test.
- One test at a time – smallest increment.
- If a test is hard to write, question the design before continuing.
- **Never refactor while RED.** Get to GREEN first.

## Checklist Per Cycle

```
[ ] Test describes behavior, not implementation
[ ] Test uses public interface only
[ ] Test would survive internal refactor
[ ] Code is minimal for this test
[ ] No speculative features added
```

## Integration with SDD
- Read `.skillgrid/project/CONTEXT.md` before starting.
- If TDD uncovers new domain assumptions, propose updating CONTEXT.md.