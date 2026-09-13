# Archived Work Items

这里保存已经完成并提交的历史工作项。`WORK_ITEMS.md` 保存当前任务正文，`WORK_ITEMS.json` 保存当前任务状态；本文件只用于查阅历史，不承载进行中的状态。

归档规则：

- 只有完成自动检查、人工审核并创建本地 commit 的工作项才能归档。
- 归档时保留原工作项编号、验收条件、检查结果、决策链接和 commit 备注。
- 新任务和未完成任务继续维护在 [`WORK_ITEMS.md`](WORK_ITEMS.md) 中。

## Done

- [x] W-005 补齐回归与验收自动化
  - 目标：把 race、架构依赖、存储格式回归和测试环境隔离变成可重复执行的自动化证据。
  - 验收：完整验证包含 race；依赖边界由测试断言；SQLite/bbolt 当前格式有回归测试；测试不读取真实 `~/.ttl`；文档与实际入口一致。
  - 检查：`./scripts/regression.sh`、直接运行的 `go test -race ./...`、`./scripts/verify.sh`、`git diff --check` 和 `gofmt -s -l .` 已通过。
  - 决策：无需记录，原因：补齐既有回归设计的实现证据，不改变产品架构。
  - 备注：人工审核结果为 `PASS`；TUI、`db.Stor` 和旧入口清理已分别登记为 W-006、W-007。

- [x] W-004 拆分客户端与后端服务工程边界
  - 目标：让 CLI/TUI 客户端、远端服务和共享核心的源码与构建入口一眼可辨。
  - 验收：形成可执行的目标目录、依赖规则、迁移步骤和验收命令；明确现有文件到新目录的映射；客户端与服务端调用方、测试和文档随契约变更同步调整。
  - 检查：`./scripts/verify.sh`、`git diff --check` 与 `gofmt -s -l .` 已通过。
  - 决策：`docs/decisions/2026-09-12-separate-client-server-layout.md`（adopted）。
  - 备注：已提交为 `b6493cd`，完成的是源码入口和包边界；显式存储依赖、`db.Stor` 清理和云端独立交付分别由后续任务跟踪，不能据此宣称部署拆分已经完成。

- [x] W-006 消除客户端对全局存储的依赖
  - 目标：让客户端命令和同步流程通过显式构造参数使用存储，不再读取 `db.Stor`。
  - 验收：`command`、`internal/client/cli` 和 `internal/client/sync` 不直接读取 `db.Stor`；客户端生产链路不再依赖旧存储门面；当前 CLI 行为测试通过。
  - 检查：`gofmt -s -l .`、`git diff --check`、`go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`./scripts/regression.sh`、`go build -o ttl ./cmd/ttl`、`go build -o ttl-server ./cmd/ttl-server` 和 `go test ./internal/architecture` 均已通过。
  - 决策：已决策；沿用 [`docs/decisions/2026-09-12-separate-client-server-layout.md`](docs/decisions/2026-09-12-separate-client-server-layout.md)，由 `internal/client/cli` 负责创建并注入显式客户端服务，命令通过 context 使用，命令执行结束后由客户端入口关闭。实现结果已回写该记录，不新增重复 ADR。
  - 文档：技术设计：[`docs/tech-designs/2026-09-12-client-storage-lifecycle.md`](docs/tech-designs/2026-09-12-client-storage-lifecycle.md)；设计评审：[`docs/reviews/2026-09-12-client-storage-lifecycle-design.md`](docs/reviews/2026-09-12-client-storage-lifecycle-design.md)；代码评审：[`docs/reviews/2026-09-12-client-storage-lifecycle-code.md`](docs/reviews/2026-09-12-client-storage-lifecycle-code.md)（`PASS`）。
  - 提交：`c58c644`（`refactor: inject client storage service`）。
  - 备注：客户端生产代码已脱离全局 `db.Stor`，删除了根目录入口和 `ttl server` 兼容代理；旧 `db` 包仍被历史测试使用，物理删除单独处理，避免扩大本任务范围。

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
