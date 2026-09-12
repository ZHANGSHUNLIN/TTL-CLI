---
name: personal-work-item
description: Use when starting, updating, reviewing, blocking, or completing a task in ttl-cli. Keeps WORK_ITEMS.md as the local source of task state and connects implementation, checks, manual review, decisions, and commits.
---

# Personal Work Item

Use `WORK_ITEMS.md` as the current project-level task state. Completed items may be moved, with their full evidence intact, to the repository's `WORK_ITEMS_ARCHIVE.md` history file. Do not introduce GitHub Issue, Project, Pull Request, reviewer assignment, or approval-score steps for this local workflow.

## Start a task

Create one item under `Inbox` with:

- a stable `W-XXX` id;
- a short objective;
- observable acceptance criteria;
- focused checks;
- a decision-record expectation;
- risks or dependencies.

When work starts, move the item to `Doing` without changing its id.

## During implementation

Keep `WORK_ITEMS.md` concise. Use the session todo list for temporary substeps, but do not use it as a replacement for the project task item.

If work cannot continue, move the item to `Blocked` and record:

- the concrete blocker;
- what is waiting;
- the next action that will unblock it.

Do not use `Blocked` for ordinary uncertainty or unfinished work.

## Review and completion

Move the item to `Review` only after implementation, tests, and required docs are present. Then:

1. Run `personal-pre-review-checks`.
2. Perform `personal-code-review`.
3. Confirm the acceptance criteria manually.
4. Confirm the decision-record requirement.
5. Run `git diff --check` and inspect `git status --short`.
6. Create a local commit.
7. Move the item to `Done` and record the important checks.

If the repository is not committed yet, keep the item in `Review`; do not claim `Done`.
