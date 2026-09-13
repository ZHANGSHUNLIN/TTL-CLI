---
name: personal-design-review
description: Use to review a ttl-cli technical design before implementation, checking requirement coverage, architecture fit, risk handling, testability, and task-level deliverability.
---

# Personal Technical Design Review

This is a read-only review of a technical design. It is separate from code review: the question is whether the proposed solution is safe and executable before code is written.

## Input

- Requirement document under `docs/requirements/`.
- Technical design under `docs/tech-designs/`.
- Relevant source, tests, `AGENTS.md`, and decision records.

## Output

Create `docs/reviews/YYYY-MM-DD-<slug>-design.md` and link it from the work item. Start with one verdict:

- `PASS`: implementation can be broken down.
- `CONDITIONAL`: implementation can proceed after named non-blocking fixes.
- `BLOCK`: redesign is required before task breakdown.

Findings come before the summary and include:

- severity: `BLOCK`, `MEDIUM`, or `LOW`;
- affected requirement, design section, or file;
- concrete problem and consequence;
- required correction or explicit follow-up;
- verification that will prove the correction.

## Review dimensions

- Requirement and acceptance coverage.
- Fit with current package ownership and public entry points.
- Completeness of interfaces, data, configuration, errors, and compatibility behavior.
- Resource lifecycle, concurrency, security, migration, rollback, and failure recovery.
- Whether each proposed task can be independently implemented and tested.
- Whether the planned unit, integration, CLI, and regression checks observe real behavior.

## Gate

Do not produce a task breakdown while a `BLOCK` finding remains. A `CONDITIONAL` verdict must name the exact follow-up and owner; unresolved high-impact uncertainty remains `BLOCK`.

## Not this stage

Do not modify the implementation, approve code, or replace missing design details with guesses.
