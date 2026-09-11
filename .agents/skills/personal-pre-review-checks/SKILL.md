---
name: personal-pre-review-checks
description: Use before moving a ttl-cli task to Review or claiming that local validation passes. Selects the smallest Go checks that cover the current diff and reports exactly what ran.
---

# Personal Pre-Review Checks

Select checks from the changed surface. Do not run the full verification script by default for a narrow change, and do not claim checks passed when they were not run.

## Inspect the change

Run:

```sh
git status --short
git diff --stat
git diff --check
```

Read the complete diff and identify whether it touches:

- Go source or tests;
- API, sync, database, storage, migration, or integration code;
- dependencies or generated files;
- documentation, workflow rules, or localization.

## Check matrix

| Changed surface | Minimum evidence |
| --- | --- |
| Go source or tests | `gofmt -s -l .`, focused `go test` |
| Normal behavior change | `go test ./...`, `go vet ./...` |
| API, sync, database, storage, migration, or cross-package change | `go test ./...`, `go test ./integration_test/...`, `go vet ./...` |
| Dependency change | `go mod tidy`, then the relevant tests |
| Broad or release-sensitive change | `./scripts/verify.sh` |
| Documentation-only change | `git diff --check` and manual link/content review |

If `gofmt -s -l .` prints files, format the changed Go files and inspect the resulting diff. Do not format unrelated files without checking why they changed.

## Report

Record in the task item:

- commands actually run;
- whether they passed;
- any known skipped check and why;
- the next action if a check failed.

Tests are evidence about behavior, not a substitute for reading the diff or checking acceptance criteria.
