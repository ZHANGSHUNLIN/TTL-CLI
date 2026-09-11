---
name: personal-find-simplifications
description: Use when looking for dead code, duplicated logic, unnecessary abstractions, speculative configuration, hand-rolled helpers, or documentation that makes ttl-cli harder to maintain.
---

# Personal Simplification Review

Find a few evidence-backed simplifications instead of producing a broad list of preferences. Read the task, current code, tests, and `docs/decisions/` before proposing removal.

## Strong candidates

- A public function, config field, command, or helper has no production consumer.
- Two code paths implement the same behavior with different edge cases.
- A package or dependency exists only for a narrow test or abandoned feature.
- A configuration option is not needed to support a real deployment choice.
- A comment or document repeats code, records review history, or describes an obsolete design.
- Hand-written code duplicates a maintained dependency or a standard-library capability.
- A test hook protects behavior that no user or production path relies on.

## Guardrails

- Do not remove behavior merely because it is small or unfamiliar.
- Check CLI compatibility, data formats, migrations, localization, and integration callers.
- Preserve a useful decision record when a rejected alternative could be proposed again.
- Separate independent cleanup from a behavior change.
- Prefer one small deletion or extraction with tests over a speculative refactor.

## Output

For each candidate, state:

- what can be simplified;
- evidence that it is unused, duplicated, or overbuilt;
- behavior and compatibility risk;
- the smallest safe change;
- the check that would prove the simplification.

This Skill proposes simplifications. Make edits only when the task explicitly asks for them.
