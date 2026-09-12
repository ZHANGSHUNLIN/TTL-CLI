# Client Storage Lifecycle Code Review

日期：2026-09-12
任务：W-006 / T-01
结论：`PASS`

## Findings

无 `BLOCK`、`MEDIUM` 或 `LOW` 发现。

评审中发现并已修复一处历史记录错误传播差异：`Service.RecordCommandHistory` 现在与旧门面一致，将存储写入错误返回给客户端入口。资源历史清理的 debug 输出也留在 CLI 层，`app.Service` 不产生终端副作用。

## Acceptance Review

- `command`、`internal/client/cli` 和同步实现不再 import `ttl-cli/db` 或读取 `db.Stor`。
- `internal/client/cli` 创建、注入并关闭一个显式 `app.Service`；命令执行失败时同样执行关闭。
- `ttl` 不再暴露 `server` 命令，也不再依赖 `internal/server`；根目录旧入口已删除。
- 资源、标签、审计、历史、日志、加密、导入导出、workspace 和显式同步均使用新的服务或 `core/storage.Storage` 契约。
- 没有修改数据库、配置、JSON 或 `/api/v1` 格式。

## Evidence

- `gofmt -s -l .`
- `git diff --check`
- `go test ./...`
- `go test ./integration_test/...`
- `go test -race ./...`
- `go vet ./...`
- `./scripts/regression.sh`
- `go build -o ttl ./cmd/ttl`
- `go build -o ttl-server ./cmd/ttl-server`
- `go test ./internal/architecture`
- `go list -deps ./cmd/ttl` 不包含 `ttl-cli/db` 或 `ttl-cli/internal/server/*`

上述检查均通过。

## Deferred Cleanup

旧 `db` 包仍被部分历史测试直接使用。它已经退出客户端生产依赖图，物理删除与测试迁移单独处理，不阻塞 W-006 的显式客户端生命周期验收。
