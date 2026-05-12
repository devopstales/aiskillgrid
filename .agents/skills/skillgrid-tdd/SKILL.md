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

### 1. RED – Write a failing test
- Based on the user's request, write **one** minimal test that describes the desired behaviour.
- The test **must** fail before you write any production code.
- Output: *"RED – Test written. Run `[test command]` – it fails as expected."*

### 2. GREEN – Make it pass
- Write the simplest possible production code that makes the test pass.
- No refactoring. No extra features.
- Output: *"GREEN – Test now passes."*

### 3. REFACTOR – Clean up
- Improve code (remove duplication, clarify names) without changing behaviour.
- Re‑run test after each change.
- Output: *"REFACTOR – Completed. Tests still pass."*

### 4. Repeat or confirm
- Ask: *"Another test, or stop?"*

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

## Integration with SDD
- Read `.skillgrid/project/CONTEXT.md` before starting.
- If TDD uncovers new domain assumptions, propose updating CONTEXT.md.