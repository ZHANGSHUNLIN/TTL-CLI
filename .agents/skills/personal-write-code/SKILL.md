---
name: personal-write-code
description: Use to implement one approved ttl-cli WBS task with the smallest repository-consistent code change and explicit validation evidence.
---

# Personal Write Code

Implement one `T-XX` task at a time. Read the parent requirement, technical design, task breakdown, `AGENTS.md`, affected code, and nearby tests before editing.

## Preconditions

- The task has observable acceptance criteria and known dependencies.
- The design review is `PASS` or its `CONDITIONAL` items are resolved.
- The task is marked `Doing` in `WORK_ITEMS.md` when project-level tracking is needed.

If these are missing, stop and return to the relevant upstream stage. Do not guess.

## Implementation rules

- Inspect the current owner and an existing analogous implementation first.
- Preserve package boundaries, CLI compatibility, storage formats, localization, and error conventions.
- Prefer existing helpers and standard dependencies over new abstractions.
- Keep the diff scoped to the task; do not perform unrelated cleanup.
- For changes touching three or more files, public APIs, storage, migration, synchronization, authentication, or other core flows, state a short implementation plan before editing.
- Add comments only for non-obvious invariants or ownership rules.

## Output

Return:

- changed files and behavior;
- acceptance criteria covered;
- focused checks actually run;
- remaining risks or deferred work.

Move the task to `Review` only after code, required tests, and docs are present. Use `personal-write-unit-test` and `personal-regression-testing` when the task needs those layers.

## Not this stage

Do not mark the task `Done`, create a commit, or claim a passing check that was not run.
