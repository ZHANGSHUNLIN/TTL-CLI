---
name: personal-test-reliability
description: Use when adding or debugging ttl-cli tests involving databases, files, HTTP servers, subprocesses, goroutines, shared process state, clocks, or integration resources.
---

# Personal Test Reliability

Design tests for the repository's real execution environment, not only a quiet single-test run.

## Isolation

- Use `t.TempDir()` for filesystem and database state.
- Use `httptest` for HTTP servers and close every server or response body.
- Avoid fixed ports, fixed paths, shared global configuration, and user-home data.
- Do not depend on test order or a previously created database.
- Use explicit cleanup for goroutines, subprocesses, channels, and background sync work.

## Concurrency

Before using `t.Parallel()`, confirm that the test does not share:

- environment variables;
- process-global configuration;
- database files or ports;
- mutable package-level state;
- a common user or workspace directory.

If a resource is shared, document its owner and use an atomic allocation or a serialized test instead of relying on timing.

## Failure diagnosis

When a test is flaky, record:

- the shared resource;
- the expected readiness signal;
- the cleanup owner;
- the quiescent signal;
- whether the failure reproduces under `go test -count=20` or `go test -race`.

Do not make a test reliable by adding arbitrary sleeps. Replace timing guesses with readiness, completion, or cancellation signals.

## Evidence

Run the narrowest reproducer first. For lifecycle or cross-package changes, include the relevant integration test package and race detector when practical. Keep fixtures deterministic and change the fixture or ownership model rather than weakening assertions.
