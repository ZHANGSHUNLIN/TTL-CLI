---
name: personal-write-unit-test
description: Use after a ttl-cli code change to add focused production-quality unit tests for its behavior, boundaries, failures, and local invariants.
---

# Personal Write Unit Test

Use the current WBS task, requirement, design, and code diff to add or extend tests. Prefer the existing test file and helpers for the owning package.

## Test design

- Assert observable behavior, returned values, persisted state, errors, and important side effects.
- Cover the normal path plus boundary or failure cases introduced by the change.
- Use `t.TempDir()` for files and databases, `httptest` for local HTTP, and explicit cleanup for servers, goroutines, subprocesses, and responses.
- Avoid fixed paths, ports, user-home data, test order, arbitrary sleeps, and shared mutable global state.
- Mock only an expensive or non-deterministic boundary; keep downstream code real.
- Do not add a test that only proves a function was called when the behavior can be observed directly.

## Output

Return:

- tests added or extended;
- behavior each test proves;
- focused test commands and results;
- coverage that belongs to integration or CLI black-box regression instead of unit tests.

When a command, output, database lifecycle, encryption, synchronization, or built-entry behavior changes, pair the unit tests with `personal-regression-testing` or the relevant integration test.

## Gate

Do not create empty tests, broad snapshots of unstable output, or assertions that merely mirror implementation details. If the test cannot isolate the behavior, report the missing seam and route the design back for correction.
