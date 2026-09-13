---
name: personal-requirement-analysis
description: Use when turning a large or unclear ttl-cli feature request into a reviewable requirement document with scope, rules, edge cases, and acceptance criteria.
---

# Personal Requirement Analysis

Use this stage before technical design for a feature that has multiple behaviors, packages, user flows, or compatibility risks. Small bug fixes may collapse this stage into the work item when the scope and acceptance criteria are already clear.

## Input

- User conversation, notes, or issue description.
- Current `AGENTS.md`, relevant README or docs, and existing behavior.
- Existing `WORK_ITEMS.md` item when one already exists.

## Output

Create `docs/requirements/YYYY-MM-DD-<slug>.md` and link it from the work item. The document must contain:

- problem and user goal;
- in-scope and out-of-scope behavior;
- user flows or use cases;
- business and compatibility rules;
- normal, boundary, and failure cases;
- observable acceptance criteria;
- open questions and assumptions.

## Rules

- Describe what the product must do, not how the code will do it.
- Use the repository's current behavior and terminology; do not invent domain rules.
- Make negative cases and data-loss risks explicit.
- Acceptance criteria must be testable through a package, integration, CLI, file, database, or API observation.
- Record unknowns as questions or assumptions instead of silently deciding them.

## Gate

Do not start technical design when the document has no clear scope, acceptance criteria, or unresolved high-impact questions. A requirement document is ready when another person can determine whether a proposed implementation satisfies it without reading the implementation.

## Not this stage

Do not choose packages, database tables, APIs, dependencies, or code structure here. Those belong to `personal-tech-design`.
