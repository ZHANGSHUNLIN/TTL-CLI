# 客户端与服务端代码结构收敛代码评审

日期：2026-09-13  
任务：W-017  
需求：[`docs/requirements/2026-09-13-W-017-code-structure-convergence.md`](../requirements/2026-09-13-W-017-code-structure-convergence.md)  
方案：[`docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md`](../tech-designs/2026-09-13-W-017-code-structure-convergence.md)  
方案评审：[`docs/reviews/2026-09-13-W-017-code-structure-convergence-design.md`](2026-09-13-W-017-code-structure-convergence-design.md)  
结论：PASS

## Review Scope

本次评审覆盖 W-017 一次性结构切换的完整 diff、迁移后的客户端/服务端依赖边界、共享资源模型和存储契约、SQLite/bbolt/remote/sync 适配器、测试迁移、旧包删除以及同步后的项目文档。工作区中已有的 W-010、AGENTS 和工作流文档改动不属于本任务，不纳入提交。

## Findings

未发现阻塞或条件性问题。

- 客户端与服务端正式入口分别位于 `cmd/ttl` 和 `cmd/ttl-server`，服务端不再传递依赖 `internal/client`。
- 生产代码和测试均已迁移到最终 owner 包；`command`、`models`、`db`、`sync`、`conf`、`crypto`、`i18n`、`util` 旧包不存在于 `go list ./...`。
- `internal/core/resource` 和 `internal/core/storage` 成为共享模型与契约唯一来源，适配器没有反向依赖 core 之外的入口层。
- 本次变更没有修改命令、API、配置文件、数据库格式、加密密钥或同步语义；集成、黑盒和竞态检查均通过。
- 迁移后的命令适配仍保留既有命令树和用户可见输出；包级命令变量属于客户端命令 owner 包，不构成旧顶层命令入口或全局存储状态。

## Evidence

- `gofmt -s -l .`
- `go test ./...`
- `go test ./integration_test/...`
- `go test -race ./...`
- `go vet ./...`
- `go test ./internal/architecture`
- `go build -o ttl ./cmd/ttl`
- `go build -o ttl-server ./cmd/ttl-server`
- `./scripts/regression.sh`
- `./scripts/cli-composability.sh`
- `./scripts/verify.sh`
- `git diff --check`
- owner 已人工核对完整 diff、临时 HOME 下的关键 CLI/TUI 流程和 W-017 验收条件。

## Decision

PASS。W-017 的验收条件已满足，可以创建本地 Conventional Commit；提交成功前保持任务状态为 `Review`，提交后再更新为 `Done`。
