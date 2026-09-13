---
name: personal-task-breakdown
description: Use after technical design review to split a large ttl-cli feature into independently executable and verifiable WBS tasks.
---

# Personal Task Breakdown

Split the approved design by feature closure, not by file list. A task may touch several files when those files form one independently verifiable behavior.

## Input

- Requirement document.
- `PASS` or resolved `CONDITIONAL` design review.
- Technical design and repository constraints.

## Output

Create `docs/task-breakdowns/YYYY-MM-DD-<slug>.md` and link it from the parent work item. Give each task a stable `T-XX` id with:

- objective;
- input;
- output;
- dependencies;
- affected ownership area;
- acceptance criteria;
- focused checks;
- unresolved risk or follow-up.

Use `WORK_ITEMS.md` for project-level state. Link child work items to the breakdown when a task is large enough to need its own status.

## Good task boundaries

- One coherent behavior or lifecycle.
- One clear owner and completion condition.
- Independently testable through a package, integration, CLI, or persisted result.
- Small enough for one coding and review cycle.
- Dependencies are explicit and ordered.

Typical feature-closure order is:

```text
contract / model
  -> core behavior
  -> command or API entry
  -> persistence / migration
  -> tests and regression
  -> documentation
```

Do not force this order when the design has a better dependency graph.

## Gate

Every task must have input, output, dependency, and acceptance criteria. If a task cannot be tested without all other tasks, split it or explain the unavoidable dependency. If a task is only “modify file X”, it is probably not a useful WBS unit.

## Not this stage

Do not redesign the feature or write code. Return to `personal-tech-design` when the design cannot be split without inventing missing interfaces.
