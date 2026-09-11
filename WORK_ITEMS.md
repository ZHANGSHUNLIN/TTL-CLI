# Work Items

这是本项目的个人开发清单，也是任务状态的唯一中转处。不依赖 GitHub
Issue、Project、Pull Request 或远程审批。

## Inbox

<!-- 新任务先放这里。每项至少写清目标和验收条件。 -->

## Doing

<!-- 当前正在处理的任务。通常只保留一项。 -->

## Review

<!-- 代码已完成，等待人工检查和验证。 -->

## Blocked

<!-- 被外部依赖、信息缺失或技术问题阻塞的任务。写明阻塞原因和下一步。 -->

## Done

<!-- 人工审核通过并完成提交的任务。 -->

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
