# Work Items

这是本项目的个人开发清单，也是任务状态的唯一中转处。不依赖 GitHub
Issue、Project、Pull Request 或远程审批。

已完成并提交的历史工作项归档在 [`WORK_ITEMS_ARCHIVE.md`](WORK_ITEMS_ARCHIVE.md)；当前任务状态仍以本文件为准。

## Inbox

<!-- 新任务先放这里。每项至少写清目标和验收条件。 -->

- [ ] W-007 实现 TUI 并收敛兼容入口
  - 目标：实现复用 core 用例的 `ttl ui`，并根据版本兼容承诺处理根入口和 `ttl server` 代理。
  - 验收：TUI 不解析 CLI 文本；CLI/TUI 使用同一存储契约；兼容入口的保留或移除有明确版本策略；安装和发布说明同步更新。
  - 检查：TUI 聚焦测试、CLI 黑盒、双二进制构建、`./scripts/verify.sh` 和人工交互验收。
  - 决策：实现前更新客户端交互与兼容策略决策记录。
  - 备注：依赖 W-006 的显式存储生命周期。

## Doing

<!-- 当前正在处理的任务。通常只保留一项。 -->

## Review

<!-- 代码已完成，等待人工检查和验证。 -->

- [ ] W-012 定义并实现 CLI 可组合性契约
  - 目标：为 `add/get/update/del/tag/dtag` 定义 JSON、非交互、stdin/stdout/stderr、稳定退出码和机器可读错误，服务脚本和自动化场景。
  - 验收：六个核心命令提供版本 1 JSON；既有默认文本语义保持兼容；无 TTY、stdin、输出流和退出码有黑盒证据；未覆盖命令明确拒绝机器模式；不阻塞 W-007。
  - 检查：`go test ./internal/client/app ./internal/client/cli`、`./scripts/cli-composability.sh`、`./scripts/regression.sh`、`go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`bash -n scripts/cli-composability.sh`、`gofmt -s -l .`、`git diff --check` 和完整 `./scripts/verify.sh` 均通过。
  - 决策：[`docs/decisions/2026-09-12-adopt-cli-composability-contract.md`](docs/decisions/2026-09-12-adopt-cli-composability-contract.md)（`adopted`）。
  - 文档：[`docs/requirements/2026-09-12-cli-composability-contract.md`](docs/requirements/2026-09-12-cli-composability-contract.md)、[`docs/tech-designs/2026-09-12-cli-composability-contract.md`](docs/tech-designs/2026-09-12-cli-composability-contract.md)、[`docs/reviews/2026-09-12-cli-composability-contract-design.md`](docs/reviews/2026-09-12-cli-composability-contract-design.md)（`PASS`）、[`docs/reviews/2026-09-12-cli-composability-contract-code.md`](docs/reviews/2026-09-12-cli-composability-contract-code.md)（`PASS`）。
  - 备注：实现、测试和代码评审已完成，等待 owner 人工确认和本地 commit；不扩大到 search、pick、sync、工作空间命令或 W-006 旧 handler 清理。

## Blocked

<!-- 被外部依赖、信息缺失或技术问题阻塞的任务。写明阻塞原因和下一步。 -->

## Done

<!-- 人工审核通过并完成提交的任务；历史记录见 [WORK_ITEMS_ARCHIVE.md](WORK_ITEMS_ARCHIVE.md)。 -->

- [x] W-006 消除客户端对全局存储的依赖
  - 目标：让客户端命令和同步流程通过显式构造参数使用存储，不再读取 `db.Stor`。
  - 验收：`command`、`internal/client/cli` 和 `internal/client/sync` 不直接读取 `db.Stor`；客户端生产链路不再依赖旧存储门面；当前 CLI 行为测试通过。
  - 检查：`gofmt -s -l .`、`git diff --check`、`go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`./scripts/regression.sh`、`go build -o ttl ./cmd/ttl`、`go build -o ttl-server ./cmd/ttl-server` 和 `go test ./internal/architecture` 均已通过。
  - 决策：已决策；沿用 [`docs/decisions/2026-09-12-separate-client-server-layout.md`](docs/decisions/2026-09-12-separate-client-server-layout.md)，由 `internal/client/cli` 负责创建并注入显式客户端服务，命令通过 context 使用，命令执行结束后由客户端入口关闭。实现结果已回写该记录，不新增重复 ADR。
  - 文档：技术设计：[`docs/tech-designs/2026-09-12-client-storage-lifecycle.md`](docs/tech-designs/2026-09-12-client-storage-lifecycle.md)；设计评审：[`docs/reviews/2026-09-12-client-storage-lifecycle-design.md`](docs/reviews/2026-09-12-client-storage-lifecycle-design.md)；代码评审：[`docs/reviews/2026-09-12-client-storage-lifecycle-code.md`](docs/reviews/2026-09-12-client-storage-lifecycle-code.md)（`PASS`）。
  - 提交：`c58c644`（`refactor: inject client storage service`）。
  - 备注：客户端生产代码已脱离全局 `db.Stor`，删除了根目录入口和 `ttl server` 兼容代理；旧 `db` 包仍被历史测试使用，物理删除单独处理，避免扩大本任务范围。


## Task Format

```md
- [ ] W-001 简短任务名称
  - 目标：要解决的问题或要交付的结果
  - 验收：可以客观判断完成与否的条件
  - 检查：需要运行的命令或需要人工确认的内容
  - 决策：需要的 docs/decisions 记录；没有则写“无需记录，原因：...”
  - 备注：依赖、风险或后续工作
```

状态迁移：

```text
Inbox → Doing → Review → Done
             ↘ Blocked
```
