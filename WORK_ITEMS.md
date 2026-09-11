# Work Items

这是本项目的个人开发清单，也是任务状态的唯一中转处。不依赖 GitHub
Issue、Project、Pull Request 或远程审批。

## Inbox

<!-- 新任务先放这里。每项至少写清目标和验收条件。 -->

## Doing

<!-- 当前正在处理的任务。通常只保留一项。 -->

## Review

<!-- 代码已完成，等待人工检查和验证。 -->

- [ ] W-004 拆分客户端与后端服务工程边界
  - 目标：让 CLI/TUI 客户端、远端服务和共享核心的源码与构建入口一眼可辨。
  - 验收：形成可执行的目标目录、依赖规则、迁移步骤和验收命令；明确现有文件到新目录的映射；迁移期间保持 `ttl` CLI 命令和 HTTP `/api/v1` 契约兼容。
  - 检查：阶段 A、B 已完成；阶段 C 已建立 core 契约并把 SQLite/bbolt 实现迁入 `internal/storage`。`./scripts/verify.sh` 已通过（双二进制构建、CLI 黑盒、单元、集成、vet）；`git diff --check` 与 `gofmt -s -l .` 已通过。
  - 决策：`docs/decisions/2026-09-12-separate-client-server-layout.md`（adopted）。
  - 备注：客户端 remote/sync、服务端 API/tenant、共享 core 与具体 storage 已有明确目录；`models`、`db` 暂作为兼容门面。人工检查未发现阻断项；显式存储依赖与 `db.Stor` 清理作为后续任务。

- [ ] W-003 梳理工程目录与项目导览
  - 目标：明确根目录、核心包、启动链路和常见修改入口，减少目录认知成本。
  - 验收：根目录有与当前代码一致的 `PROJECT_OVERVIEW.md`；README 的目录树不再引用不存在的文件；本地构建产物不会污染 Git 状态。
  - 检查：`git diff --check` 已通过；文档中的关键路径与命令已人工核对。
  - 决策：无需记录，原因：本轮只整理工程文档和忽略规则，不调整包边界、运行行为或兼容性。
  - 备注：保留多语言 README 和根目录安装脚本的稳定入口；后续代码拆分另立任务。人工审核并创建本地 commit 后再移动到 `Done`。

- [ ] W-001 搭建个人工程 Skill
  - 目标：为个人练手项目建立可复用的工程 Skill，覆盖任务状态、进入 Review 前检查、代码审核、测试可靠性、简化审查、工程文档和决策记录。
  - 验收：`.agents/skills/` 下有 Skill 索引和 7 个可执行规则；`AGENTS.md` 与流程文档能找到入口；`docs/decisions/` 有对应决策记录。
  - 检查：`git diff --check`、`gofmt -s -l .`、`go test ./...`、`go vet ./...`；已验证 `git diff --check`、`gofmt -s -l .`、`go vet ./...`，`go test ./...` 目前被现有 `db/TestGetDBPath` 的 `.bbolt`/`.db` 断言失败阻塞。
  - 决策：`docs/decisions/2026-09-11-add-personal-engineering-skills.md`
  - 备注：人工审核通过并创建本地 commit 后再移动到 `Done`。
- [ ] W-002 搭建个人版分层回归测试
  - 目标：把 DSH 的分层回归思路移植为个人 Go CLI 可执行的单元、集成、CLI 黑盒和完整验证流程。
  - 验收：有独立的 `scripts/regression.sh`；`scripts/verify.sh` 汇总构建、黑盒、单测、集成测试和 `go vet`；有回归测试文档、Skill 和决策记录；黑盒回归不写入真实 `~/.ttl`。
  - 检查：`git diff --check`、`gofmt -s -l .`、`go vet ./...`、`go test ./integration_test/...`、`./scripts/regression.sh`、`./scripts/verify.sh`；总验证脚本内部使用临时 `HOME` 和现有 `GOPATH` 运行 `go test ./...`。
  - 决策：`docs/decisions/2026-09-11-add-personal-regression-testing.md`
  - 备注：上述检查已通过。直接在当前用户配置下运行 `go test ./...` 仍会触发既有 `db/TestGetDBPath` 的 `.bbolt`/`.db` 断言问题；总验证已通过临时 `HOME` 隔离该环境影响。人工审核通过并创建本地 commit 后再移动到 `Done`。

## Blocked

<!-- 被外部依赖、信息缺失或技术问题阻塞的任务。写明阻塞原因和下一步。 -->

## Done

<!-- 人工审核通过并完成提交的任务。 -->

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
