---
name: personal-regression-testing
description: Use when choosing, extending, or running ttl-cli regression checks. Maps a change to focused unit, integration, CLI black-box, or full verification evidence without requiring GitHub CI.
---

# Personal Regression Testing

Use the smallest test layer that can prove the changed behavior. The repository keeps four layers:

1. **Focused unit tests** for a package-level rule or helper.
2. **Integration tests** for storage, API, sync, encryption, filesystem, or multi-package behavior.
3. **CLI black-box regression** for behavior visible through the built `ttl` command.
4. **Full verification** for broad, release-sensitive, or cross-cutting changes.

## Select the layer

| Changed surface | Minimum evidence |
| --- | --- |
| One helper or isolated package rule | Focused `go test ./path -run TestName` |
| Command behavior or user-visible output | Focused tests plus `./scripts/regression.sh` |
| Database, encryption, API, sync, migration, or cross-package behavior | `go test ./...` and `go test ./integration_test/...` |
| Broad refactor, dependency change, or pre-commit confidence pass | `./scripts/verify.sh` |
| Documentation or workflow-only change | `git diff --check` and manual content review |

Run `go vet ./...` for normal behavior changes when the selected checks do not already include it.

## CLI black-box rules

`scripts/regression.sh` must exercise the built binary, not only package functions. It must:

- use a temporary `HOME`, config file, database, and encryption key;
- verify observable output and persisted behavior after commands finish;
- cover the smallest representative lifecycle for changed commands;
- fail on an unexpected success, missing output, wrong exit status, or lost data;
- avoid real user files, fixed ports, remote services, and developer-specific paths.

When a command changes, extend the black-box flow only if the behavior is part of the normal user entry point. Do not add a second test for an internal implementation detail that a focused package test already proves.

## DSH mapping

This is the personal equivalent of DSH's layered regression model:

- unit and integration tests cover package and capability contracts;
- the built CLI regression covers expected user-visible output and persisted results;
- `go build` and `go vet` cover artifact and static gates;
- real external API, browser snapshot, recorded model session, platform matrix, and GitHub workflow gates are intentionally out of scope for this local CLI practice project.

## Report evidence

Record the commands actually run in the `WORK_ITEMS.md` task. If a check fails, keep the task in `Review` or move it to `Blocked`, record the concrete failure, and do not describe the regression layer as passing.
