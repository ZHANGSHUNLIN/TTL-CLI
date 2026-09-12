# Work Items

这是本项目的个人开发清单，也是任务状态的唯一中转处。不依赖 GitHub
Issue、Project、Pull Request 或远程审批。

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

## Blocked

<!-- 被外部依赖、信息缺失或技术问题阻塞的任务。写明阻塞原因和下一步。 -->

## Done

<!-- 人工审核通过并完成提交的任务。 -->

- [x] W-006 消除客户端对全局存储的依赖
  - 目标：让客户端命令和同步流程通过显式构造参数使用存储，不再读取 `db.Stor`。
  - 验收：`command`、`internal/client/cli` 和 `internal/client/sync` 不直接读取 `db.Stor`；客户端生产链路不再依赖旧存储门面；当前 CLI 行为测试通过。
  - 检查：`gofmt -s -l .`、`git diff --check`、`go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`./scripts/regression.sh`、`go build -o ttl ./cmd/ttl`、`go build -o ttl-server ./cmd/ttl-server` 和 `go test ./internal/architecture` 均已通过。
  - 决策：已决策；沿用 [`docs/decisions/2026-09-12-separate-client-server-layout.md`](docs/decisions/2026-09-12-separate-client-server-layout.md)，由 `internal/client/cli` 负责创建并注入显式客户端服务，命令通过 context 使用，命令执行结束后由客户端入口关闭。实现结果已回写该记录，不新增重复 ADR。
  - 文档：技术设计：[`docs/tech-designs/2026-09-12-client-storage-lifecycle.md`](docs/tech-designs/2026-09-12-client-storage-lifecycle.md)；设计评审：[`docs/reviews/2026-09-12-client-storage-lifecycle-design.md`](docs/reviews/2026-09-12-client-storage-lifecycle-design.md)；代码评审：[`docs/reviews/2026-09-12-client-storage-lifecycle-code.md`](docs/reviews/2026-09-12-client-storage-lifecycle-code.md)（`PASS`）。
  - 提交：`c58c644`（`refactor: inject client storage service`）。
  - 备注：客户端生产代码已脱离全局 `db.Stor`，删除了根目录入口和 `ttl server` 兼容代理；旧 `db` 包仍被历史测试使用，物理删除单独处理，避免扩大本任务范围。

- [x] W-005 补齐回归与验收自动化
  - 目标：把 race、架构依赖、旧数据兼容和测试环境隔离变成可重复执行的自动化证据。
  - 验收：完整验证包含 race；依赖边界由测试断言；SQLite/bbolt 旧格式数据有兼容 fixture；测试不读取真实 `~/.ttl`；文档与实际入口一致。
  - 检查：`./scripts/regression.sh`、直接运行的 `go test -race ./...`、`./scripts/verify.sh`、`git diff --check` 和 `gofmt -s -l .` 已通过。
  - 决策：无需记录，原因：补齐既有回归设计的实现证据，不改变产品架构或兼容策略。
  - 备注：人工审核结果为 `PASS`；TUI、`db.Stor` 清理和移除兼容入口已分别登记为 W-006、W-007。

- [x] W-004 拆分客户端与后端服务工程边界
  - 目标：让 CLI/TUI 客户端、远端服务和共享核心的源码与构建入口一眼可辨。
  - 验收：形成可执行的目标目录、依赖规则、迁移步骤和验收命令；明确现有文件到新目录的映射；迁移期间保持 `ttl` CLI 命令和 HTTP `/api/v1` 契约兼容。
  - 检查：`./scripts/verify.sh`、`git diff --check` 与 `gofmt -s -l .` 已通过。
  - 决策：`docs/decisions/2026-09-12-separate-client-server-layout.md`（adopted）。
  - 备注：已提交为 `b6493cd`；显式存储依赖与 `db.Stor` 清理作为后续任务。
- [x] W-003 梳理工程目录与项目导览
  - 目标：明确根目录、核心包、启动链路和常见修改入口，减少目录认知成本。
  - 验收：根目录有与当前代码一致的 `PROJECT_OVERVIEW.md`；README 的目录树不再引用不存在的文件；本地构建产物不会污染 Git 状态。
  - 检查：完整验证及文档路径人工核对已通过。
  - 决策：无需记录，原因：只整理工程文档和忽略规则。
  - 备注：已提交为 `b6493cd`。
- [x] W-001 搭建个人工程 Skill
  - 目标：建立覆盖任务状态、检查、审核、测试、简化、文档和决策记录的个人工程 Skill。
  - 验收：`.agents/skills/`、`AGENTS.md`、流程文档和决策记录均已落地。
  - 检查：`./scripts/verify.sh` 已通过。
  - 决策：`docs/decisions/2026-09-11-add-personal-engineering-skills.md`。
  - 备注：已提交为 `b6493cd`。
- [x] W-002 搭建个人版分层回归测试
  - 目标：建立单元、集成、CLI 黑盒和完整验证流程。
  - 验收：回归脚本、总验证脚本、文档、Skill 和决策记录均已落地。
  - 检查：`./scripts/verify.sh` 已通过，并使用临时 `HOME` 隔离真实用户数据。
  - 决策：`docs/decisions/2026-09-11-add-personal-regression-testing.md`。
  - 备注：已提交为 `b6493cd`。

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
