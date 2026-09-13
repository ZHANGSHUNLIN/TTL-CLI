---
name: personal-work-item
description: Use when starting, updating, reviewing, blocking, or completing a task in ttl-cli. Keeps WORK_ITEMS.md for task content and WORK_ITEMS.json for task state, connecting implementation, checks, manual review, decisions, and commits.
---

# Personal Work Item

Use `WORK_ITEMS.md` for current project-level task content and `WORK_ITEMS.json` for current task state. Completed items may be moved, with their full evidence intact, to the repository's `WORK_ITEMS_ARCHIVE.md` history file. Do not introduce GitHub Issue, Project, Pull Request, reviewer assignment, or approval-score steps for this local workflow.

## Start a task

For a new feature, bug fix, refactor, docs change, maintenance task, spike, or other bounded work, use the `personal-workflow-dashboard` template creation flow first. It creates one item under the `Tasks` section in `WORK_ITEMS.md`, matching state metadata under `items` in `WORK_ITEMS.json` with `status: "Inbox"`, and six linked Markdown skeletons for requirements, technical design, design review, WBS, tests, and acceptance. If the dashboard is unavailable, reproduce that same complete set manually from `docs/templates/`; do not create an incomplete one-line item.

Each task needs:

- a stable `W-XXX` id;
- a short objective;
- observable acceptance criteria;
- focused checks;
- a decision-record expectation;
- risks or dependencies.

The generated task starts at `当前阶段：requirements`. Treat every generated document as a placeholder waiting for real content. The presence of a file or link never satisfies a stage gate by itself.

The initial bundle does not pre-create a code-review conclusion: `personal-code-review` writes that evidence after the implementation diff and tests exist. Commit evidence is added only after the review gate passes.

For a large feature, link the requirement, technical design, design review, and task breakdown documents from the item. Use `T-XX` ids inside the breakdown and create separate `W-XXX` items only when a child task needs its own project-level status.

When work starts, update the item's JSON metadata to `status: "Doing"` without changing its id.

## During implementation

Keep `WORK_ITEMS.md` concise and keep state changes in `WORK_ITEMS.json`. Use the session todo list for temporary substeps, but do not use it as a replacement for the project task item.

If work cannot continue, update the item to `status: "Blocked"` and record in its Markdown notes:

- the concrete blocker;
- what is waiting;
- the next action that will unblock it.

Do not use `Blocked` for ordinary uncertainty or unfinished work. Unresolved product or design questions stay in the upstream document until answered.

## Review and completion

Update the item to `status: "Review"` only after implementation, tests, and required docs are present. Then:

1. Run `personal-pre-review-checks`.
2. Perform `personal-code-review`.
3. Confirm the acceptance criteria manually.
4. Confirm the decision-record requirement.
5. Run `git diff --check` and inspect `git status --short`.
6. Create a local commit.
7. Update the item to `status: "Done"` and `completed: true`, then record the important checks.

If the repository is not committed yet, keep the JSON item in `Review`; do not claim `Done`.
