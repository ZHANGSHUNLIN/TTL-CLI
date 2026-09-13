---
name: personal-commit
description: Use after local review to create a scoped conventional Git commit for a ttl-cli task and finish its WORK_ITEMS.md content and WORK_ITEMS.json state without GitHub PR or approval workflow.
---

# Personal Commit

This is the final local delivery stage. It creates a Git commit only after the owner has reviewed the diff and all required evidence is available.

## Preconditions

- The task is in `Review`.
- `personal-pre-review-checks`, required unit/integration/regression checks, and manual review are complete.
- All `BLOCK` findings are resolved.
- Required requirement, design, breakdown, review, documentation, and decision records are present when the selected workflow needs them; a small change may keep its rationale in `WORK_ITEMS.md`.
- The worktree contains no unrelated generated files, credentials, or local data.

## Procedure

1. Inspect `git status --short`, the complete diff, and `git diff --check`.
2. Reconfirm the task acceptance criteria and record the commands actually run in the task entry in `WORK_ITEMS.md`.
3. Stage only files belonging to the task.
4. Create a concise imperative Conventional Commit message such as `feat: add workspace export`.
5. Inspect the created commit with `git show --stat --oneline HEAD`.
6. Only after the commit succeeds, update the task to `Done` and `completed: true` in `WORK_ITEMS.json`, then retain the commit id in `WORK_ITEMS.md`.

## Rules

- This project uses local commits, not GitHub Issues, Pull Requests, approval scores, or automatic merge.
- Never use destructive reset, checkout, or raw force operations as part of this stage.
- Do not include unrelated user changes. If the worktree is mixed, stage explicit paths or pause for manual separation.
- A failed commit or failed check leaves the task in `Review`; record the concrete blocker.

## Output

Return the commit id, message, files included, checks run, and final work-item state. A commit is delivery evidence, not a substitute for review.
