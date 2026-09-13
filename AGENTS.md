# Repository Guidelines

## Project Structure & Module Organization

`ttl-cli` contains a local client and a separately buildable backend service. `cmd/ttl/` and `cmd/ttl-server/` are the only target executable entries for the product; root `main.go`, the client-side `ttl server` command, `db/` facades, and `models/` aliases are temporary code to remove rather than supported entry points. Put client command assembly in `internal/client/cli/`, reusable client handlers in `command/`, backend HTTP/API/tenant code in `internal/server/`, concrete storage implementations in `internal/storage/`, and shared contracts in `internal/core/`. Supporting packages include `conf/` (INI and workspaces), `crypto/`, `i18n/` (with `i18n/locales/`), `sync/`, and `util/`. Unit tests live beside implementation; end-to-end server and sync scenarios are in `integration_test/`. The shared `personal-workflow-dashboard` Skill is outside this repository and is not part of the product CLI. Installers are at the root. `scripts/regression.sh` provides the built-CLI black-box regression layer, and `scripts/verify.sh` provides the full local verification pass.

## Build, Test, and Development Commands

Run these from the repository root:

```bash
go mod download                 # fetch dependencies
go build -o ttl ./cmd/ttl       # build the client
go build -o ttl-server ./cmd/ttl-server # build the backend
go run ./cmd/ttl <command>      # run the client without installing
go test ./...                   # run all package tests
go test -race -coverprofile=coverage.out ./...  # CI-style unit run
go test ./integration_test/...  # run integration tests
go vet ./...                    # static checks
gofmt -s -w .                   # format Go files
/Users/v_zhangshun01/.codex/skills/personal-workflow-dashboard/scripts/ensure-project.sh "$PWD"  # initialize or connect the shared workflow dashboard
./scripts/regression.sh         # built CLI black-box regression
./scripts/verify.sh             # build, regression, unit, integration, and vet checks
```

Use temporary directories and configuration files for manual CLI checks so local `~/.ttl` data is not changed. Run `go mod tidy` when dependencies change.

## Coding Style & Naming Conventions

Follow idiomatic Go and let `gofmt -s` define layout (tabs in source). Use mixed-case exported identifiers with Go doc comments where appropriate; use short lower-case names for locals and packages. Keep user-facing strings in i18n locale files. Preserve package boundaries and return contextual errors from command and storage layers.

## Testing Guidelines

Name tests `Test<Type>_<Scenario>` (for example, `TestUserStore_AddUser`), and isolate them with `t.TempDir` or `httptest`. Add focused package tests for behavior changes, then run `go test ./...`; run `go test ./integration_test/...` for API, cloud-storage, sync, or cross-package changes. CI exercises the race detector and publishes coverage; no numeric threshold is configured.

## Commit Guidelines

Use concise, imperative Conventional Commit-style subjects such as `feat: add workspace support`, `fix: handle duplicate keys`, or `refactor: simplify storage routing`. Keep unrelated changes separate. This personal workflow uses local commits and does not require a GitHub Pull Request. The commit or linked work item should explain behavior, affected packages, validation commands, and configuration, migration, security, or CLI output changes. Include sample output when it clarifies a user-facing change.

## Personal Practice Workflow

For local practice, use [`WORK_ITEMS.md`](WORK_ITEMS.md) for current task content, [`WORK_ITEMS.json`](WORK_ITEMS.json) for current task state, and [`WORK_ITEMS_ARCHIVE.md`](WORK_ITEMS_ARCHIVE.md) for completed task history. The local workflow has five states: `Inbox`, `Doing`, `Review`, `Blocked`, and `Done`; the dashboard keeps `Done` tasks out of the active flow, shows the completed count, and exposes history through a board popup. It also provides a read-only project Markdown browser for requirements, designs, reviews, and decisions. It does not require GitHub Issues, Projects, Pull Requests, reviewer requests, or approval scores. Follow [`docs/personal-development-workflow.md`](docs/personal-development-workflow.md) for task transitions, focused checks, manual review, and lightweight decision records. Important design choices go in [`docs/decisions/`](docs/decisions/); routine fixes can state in the task that no decision record is needed. A task reaches `Done` only after the relevant checks pass, the owner reviews the diff, required decision records are present, and a local Git commit is created; update both the status metadata and completion evidence before moving its full entry to the archive.

Project-specific engineering Skills live under [`.agents/skills/`](.agents/skills/). The `personal-workflow-dashboard` entry is a symlink to the shared global Skill at `/Users/v_zhangshun01/.codex/skills/personal-workflow-dashboard`; use its `ensure-project.sh` entrypoint to initialize or connect this project. Use [`docs/engineering-skills.md`](docs/engineering-skills.md) for the overview and choose the smallest applicable Skill:

任务类型和研发阶段进度写在 `WORK_ITEMS.md` 的可选元数据字段中；`WORK_ITEMS.json` 只维护生命周期状态、完成标记和评审记录。看板中的 `Inbox`、`Doing`、`Review`、`Blocked`、`Done` 是生命周期状态，不是需求、设计、编码等阶段列。

- `personal-requirement-analysis` turns a large feature request into a reviewable requirement document.
- `personal-tech-design` defines implementation ownership, interfaces, data changes, lifecycle, risks, and tests.
- `personal-design-review` reviews the technical design before implementation.
- `personal-task-breakdown` converts an approved design into independently executable WBS tasks.
- `personal-write-code` implements one WBS task with a scoped diff.
- `personal-write-unit-test` adds behavior-focused unit tests for that diff.
- `personal-work-item` manages `WORK_ITEMS.md` state and completion evidence.
- `personal-pre-review-checks` selects checks for the current diff.
- `personal-regression-testing` selects and runs the layered unit, integration, CLI black-box, and full checks.
- `personal-code-review` performs the owner-led manual review.
- `personal-commit` creates the final local commit after review.
- `personal-test-reliability` covers isolation and lifecycle risks in tests.
- `personal-find-simplifications` proposes evidence-backed reductions in complexity.
- `personal-prose-standard` reviews code comments, docs, help, and localized copy.
- `personal-decision-records` manages meaningful decisions in `docs/decisions/`.
