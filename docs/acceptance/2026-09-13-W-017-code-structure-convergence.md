# 客户端与服务端代码结构收敛交付验收

日期：2026-09-13  
任务：W-017  
需求：[`docs/requirements/2026-09-13-W-017-code-structure-convergence.md`](../requirements/2026-09-13-W-017-code-structure-convergence.md)  
方案：[`docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md`](../tech-designs/2026-09-13-W-017-code-structure-convergence.md)  
测试计划：[`docs/tests/2026-09-13-W-017-code-structure-convergence.md`](../tests/2026-09-13-W-017-code-structure-convergence.md)  
状态：accepted

## Acceptance Scope

本验收只覆盖结构收敛和行为不变，不覆盖 W-010 的云端独立部署交付，也不新增 CLI/TUI/API 功能。

## Delivery Criteria

- [x] `cmd/ttl` 和 `cmd/ttl-server` 是唯一正式可执行入口，构建和运行边界清晰。
- [x] 客户端 root、命令适配、app 用例、TUI、remote 和 sync 的职责与技术方案一致。
- [x] `internal/core/resource` 是资源、标签、审计、历史和日志模型的唯一来源。
- [x] `internal/core/storage` 是共享存储契约的唯一来源，core 不依赖 adapter、入口或 HTTP。
- [x] SQLite、bbolt、remote 和 sync 适配器分别位于规定的 owner 包，不存在重复实现或反向依赖。
- [x] 客户端命令不再从顶层 `command/` 注册全局命令，不使用全局输出 writer 或全局存储状态。
- [x] 根目录 `sync/` 已合并或删除，客户端同步职责集中在 `internal/client/sync`。
- [x] `db/`、`models/` 等旧包的生产和测试消费者已迁移；无消费者后才物理删除。
- [x] `conf`、`crypto`、`i18n`、`util` 的最终归属已在代码和文档中明确，且无跨层反向依赖。
- [x] 架构测试能阻止旧依赖回流，并检查直接和传递依赖。

## Behavior Preservation

- [x] `ttl add/get/update/del/tag/dtag/rename` 的参数、输出、错误和退出码不变。
- [x] `ttl ui` 的列表、搜索、查看、创建、编辑、标签、删除、打开和终端恢复行为不变。
- [x] `ttl-server serve/user` 的启动、用户管理、API Key 和租户隔离行为不变。
- [x] SQLite/bbolt 数据库、配置文件、加密密钥和 HTTP `/api/v1` 数据格式不变。
- [x] `local/cloud/sync` 的读取、写入、diff、pull、push、auto、dry-run 和冲突语义不变。

## Required Evidence

- [x] `go test ./...`
- [x] `go test ./integration_test/...`
- [x] `go test -race ./...`
- [x] `go vet ./...`
- [x] `go test ./internal/architecture`
- [x] `./scripts/regression.sh`
- [x] `./scripts/cli-composability.sh`
- [x] `go build -o ttl ./cmd/ttl`
- [x] `go build -o ttl-server ./cmd/ttl-server`
- [x] `./scripts/verify.sh`
- [x] `git diff --check`
- [x] owner 完成 diff、文档、临时 HOME 和关键 CLI/TUI 流程人工验收。

## Rollback Gate

- [x] 任一内部迁移步骤失败时可以整体回滚 W-017，不回退 W-003/W-004/W-006 历史提交。
- [x] 发现数据格式、API 契约或用户行为变化时，停止验收并另立迁移/行为变更任务。
- [x] 未完成整体迁移的 owner code review 和本地 commit 前，不将 W-017 标记为 `Done`；两项门禁现已完成。

## Sign-off

- 负责人：已确认
- 评审日期：2026-09-13
- 结论：通过，已完成本地提交 `f27315a`，W-017 可归档
