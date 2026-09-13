---
name: personal-tech-design
description: Use after requirements are clear to design a ttl-cli implementation, ownership, interfaces, data changes, lifecycle, risks, and verification plan.
---

# Personal Technical Design

Use this stage for a feature that crosses packages, changes storage or API behavior, adds a command, or needs a deliberate compatibility decision. Read the requirement document and the relevant current code before designing.

## Input

- A reviewed or sufficiently complete document under `docs/requirements/`.
- `AGENTS.md`, relevant package code, tests, configuration, and existing decision records.

## Output

Create `docs/tech-designs/YYYY-MM-DD-<slug>.md` and link it from the work item. Include:

- current behavior and affected execution path;
- proposed ownership and package/file responsibilities;
- public or package interfaces, command arguments, request/response fields, and errors;
- data, configuration, file-format, migration, and compatibility changes;
- lifecycle, cleanup, concurrency, security, and failure behavior;
- alternatives and why they were not selected;
- test and regression coverage by layer;
- rollout, rollback, or recovery plan when applicable.

## Rules

- Reuse existing project patterns before adding abstractions or dependencies.
- Keep command handlers, reusable use cases, storage, API, and shared models in their existing ownership areas.
- State which component owns creation, mutation, cleanup, and error reporting.
- Separate stable contracts from implementation details.
- A new dependency, data format, migration, security rule, or user-visible default requires a decision record under `docs/decisions/`.
- Do not write production code in this stage.

## Gate

The design is ready for review only when every requirement has an implementation owner, every changed interface has inputs/outputs/errors, and the test plan covers normal, failure, and compatibility behavior.

## Not this stage

Do not silently redefine the requirement or begin implementation. Send missing product questions back to `personal-requirement-analysis`; send implementation work to `personal-task-breakdown` after design review.
