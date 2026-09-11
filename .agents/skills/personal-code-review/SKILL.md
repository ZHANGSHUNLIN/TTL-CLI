---
name: personal-code-review
description: Use when reviewing a ttl-cli change before local commit. Prioritizes bugs, data loss, compatibility, lifecycle, user-visible behavior, missing tests, and scope mistakes over style nits.
---

# Personal Code Review

Review the current diff as the task owner. Findings come before summary. A short review with one real blocker is better than a long list of cosmetic comments.

## Review order

1. Read `AGENTS.md`, the task in `WORK_ITEMS.md`, and any linked decision record.
2. Read `git diff --stat`, then the complete `git diff`.
3. Read enough surrounding code to understand ownership and error paths.
4. Compare every changed behavior with the task's acceptance criteria.
5. Check the focused evidence selected by `personal-pre-review-checks`.

## Priority checks

- **Data safety:** deletion, overwrite, encryption, import/export, migration, and sync behavior cannot silently lose or corrupt data.
- **Error behavior:** errors retain useful context, user-visible commands fail clearly, and partial operations do not report success.
- **Resource lifecycle:** files, databases, HTTP servers, goroutines, subprocesses, and temporary directories are closed or cleaned up.
- **Compatibility:** existing CLI commands, config files, file formats, localized strings, and public package behavior remain compatible unless the task says otherwise.
- **Tests:** tests prove observable behavior and cover failure or boundary cases introduced by the change.
- **Scope:** no unrelated local files, generated artifacts, credentials, or binary outputs enter the change.
- **Documentation:** README, command help, comments, and localization match the implemented behavior.

## Review outcome

Use one of:

- `PASS`: acceptance criteria, checks, and manual review are complete.
- `CONDITIONAL`: usable after named non-blocking follow-up work.
- `BLOCK`: a correctness, safety, compatibility, or required-evidence issue must be fixed first.

Write findings with file paths and line numbers when possible. Keep `Review` tasks in `WORK_ITEMS.md` until the owner has resolved all `BLOCK` findings and manually confirmed the result.
