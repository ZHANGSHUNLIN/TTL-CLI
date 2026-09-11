# Repository Guidelines

## Project Structure & Module Organization

`ttl-cli` is a Go CLI application. `main.go` wires Cobra commands and startup. Put command handlers in `command/`, storage in `db/`, HTTP server code in `api/`, and shared types in `models/`. Supporting packages include `conf/` (INI and workspaces), `crypto/`, `i18n/` (with `i18n/locales/`), `sync/`, and `util/`. Unit tests live beside implementation; end-to-end server and sync scenarios are in `integration_test/`. Installers are at the root, and `scripts/verify.sh` provides a full verification pass.

## Build, Test, and Development Commands

Run these from the repository root:

```bash
go mod download                 # fetch dependencies
go build -o ttl .               # build the CLI
go run . <command>              # run without installing
go test ./...                   # run all package tests
go test -race -coverprofile=coverage.out ./...  # CI-style unit run
go test ./integration_test/...  # run integration tests
go vet ./...                    # static checks
gofmt -s -w .                   # format Go files
./scripts/verify.sh             # build, regression, unit, and integration checks
```

Use temporary directories and configuration files for manual CLI checks so local `~/.ttl` data is not changed. Run `go mod tidy` when dependencies change.

## Coding Style & Naming Conventions

Follow idiomatic Go and let `gofmt -s` define layout (tabs in source). Use mixed-case exported identifiers with Go doc comments where appropriate; use short lower-case names for locals and packages. Keep user-facing strings in i18n locale files. Preserve package boundaries and return contextual errors from command and storage layers.

## Testing Guidelines

Name tests `Test<Type>_<Scenario>` (for example, `TestUserStore_AddUser`), and isolate them with `t.TempDir` or `httptest`. Add focused package tests for behavior changes, then run `go test ./...`; run `go test ./integration_test/...` for API, cloud-storage, sync, or cross-package changes. CI exercises the race detector and publishes coverage; no numeric threshold is configured.

## Commit & Pull Request Guidelines

Use concise, imperative Conventional Commit-style subjects such as `feat: add workspace support`, `fix: handle duplicate keys`, or `refactor: simplify storage routing`. Keep unrelated changes separate. Pull requests should explain behavior and affected packages, link an issue when available, list validation commands, and call out configuration, migration, security, or CLI output changes. Include sample output when it clarifies a user-facing change.

## Personal Practice Workflow

For local practice, use [`WORK_ITEMS.md`](WORK_ITEMS.md) as the project task board. The local workflow has five states: `Inbox`, `Doing`, `Review`, `Blocked`, and `Done`; it does not require GitHub Issues, Projects, Pull Requests, reviewer requests, or approval scores. Follow [`docs/personal-development-workflow.md`](docs/personal-development-workflow.md) for task transitions, focused checks, manual review, and lightweight decision records. Important design choices go in [`docs/decisions/`](docs/decisions/); routine fixes can state in the task that no decision record is needed. A task reaches `Done` only after the relevant checks pass, the owner reviews the diff, required decision records are present, and a local Git commit is created.

Project-specific engineering Skills live under [`.agents/skills/`](.agents/skills/). Use [`docs/engineering-skills.md`](docs/engineering-skills.md) for the overview and choose the smallest applicable Skill:

- `personal-work-item` manages `WORK_ITEMS.md` state and completion evidence.
- `personal-pre-review-checks` selects checks for the current diff.
- `personal-code-review` performs the owner-led manual review.
- `personal-test-reliability` covers isolation and lifecycle risks in tests.
- `personal-find-simplifications` proposes evidence-backed reductions in complexity.
- `personal-prose-standard` reviews code comments, docs, help, and localized copy.
- `personal-decision-records` manages meaningful decisions in `docs/decisions/`.
