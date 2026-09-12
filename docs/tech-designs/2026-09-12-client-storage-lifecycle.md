# Client Storage Lifecycle Technical Design

日期：2026-09-12
任务：W-006 / T-01
需求：`docs/requirements/2026-09-12-capability-boundary-baseline.md`
状态：reviewed

## Current Behavior

`db/storage.go` 通过全局 `db.Stor` 负责存储构造、初始化、关闭，以及资源、审计、历史和日志操作。`command`、`internal/client/cli` 和顶层 `sync` 直接调用这些门面。客户端 root command 在 PersistentPreRun 中初始化全局存储，进程结束后再关闭它。

## Proposed Ownership

| Area | Owner | Responsibility |
| --- | --- | --- |
| `internal/client/app` | Client service | Holds one explicit `core/storage.Storage`, exposes storage operations and migration construction |
| `internal/client/cli` | Client composition root | Reads flags/config, opens the selected storage, injects the service into command context, closes it after execution |
| `command` | Command adapters | Reads the service from command context, maps arguments/errors/output, never owns storage lifecycle |
| `internal/client/sync` and `sync` | Sync use cases | Accept `core/storage.Storage` interfaces as arguments; no global storage access |
| `internal/storage/*` | Storage adapters | Own database and HTTP resources after `Init`; release them in `Close` |

## Interfaces

`internal/client/app.Service` wraps `core/storage.Storage` and delegates resource, audit, history, and log operations. `app.OpenStorage` constructs and initializes `sqlite`, `local`/`bbolt`, `cloud`, or `sync` storage from the existing CLI options. `app.WithService` and `app.ServiceFromContext` provide the per-command dependency seam.

The command surface and `/api/v1` request/response fields do not change. Missing service context returns an explicit `storage not initialized` error.

## Data And Configuration

No database, config, JSON, or API format changes. Existing storage constructors and configuration resolution remain the source of paths and timeouts. The `db` package is removed from client production imports; storage adapters continue to use the existing model aliases until the separate core-model cleanup task.

## Lifecycle And Failure Behavior

- Creation: the client composition root creates and initializes one service before resource commands run.
- Mutation: command handlers call service methods; storage errors retain the existing localized command context.
- Cleanup: `Run` closes the service after command execution. Workspace switch closes the active service before returning.
- Failure and recovery: if storage construction or initialization fails, the partially created storage is closed before returning the initialization error. A command without injected service fails explicitly.
- Concurrency: one service is scoped to one client process; no mutable package-level storage state remains in client code.

## Security And Compatibility

The change preserves existing local storage encryption, cloud API key handling, and storage formats. Removing `ttl server` from the client prevents the client binary from linking server command code. This is a development-stage breaking cleanup already adopted by the layout decision.

## Alternatives

- Keep `db.Stor` and add more wrappers: rejected because it preserves hidden lifecycle and prevents CLI/TUI reuse.
- Put a new global in `command`: rejected for the same lifecycle and test-isolation problem.
- Pass every dependency through Cobra constructors immediately: deferred; context injection keeps the current command registration surface stable while making ownership explicit.

## Test And Regression Plan

| Layer | Evidence |
| --- | --- |
| Unit | service delegation/context tests; storage initialization and close tests; command tests for missing context |
| Integration | existing resource, encryption, export/import, log, tag, and sync scenarios using explicit service/storage setup |
| CLI / regression | `./scripts/regression.sh`, both binary builds, architecture dependency test, `go test -race ./...`, `go vet ./...` |

## Rollout And Recovery

This is a development-stage internal refactor with no migration. If validation fails, keep W-006 in `Review` and revert only the task commit; existing database files remain untouched.
