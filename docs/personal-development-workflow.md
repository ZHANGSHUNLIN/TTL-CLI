# Personal Development Workflow

本项目采用面向个人练习的本地开发流程。`WORK_ITEMS.md` 是当前项目任务状态记录；已完成任务的历史保存在 [`WORK_ITEMS_ARCHIVE.md`](../WORK_ITEMS_ARCHIVE.md)。不使用 GitHub Issue、Project、Pull Request、reviewer 请求或审批积分。

工程 Skill 的索引见 [`docs/engineering-skills.md`](engineering-skills.md)，具体规则位于 [`.agents/skills/`](../.agents/skills/)。回归测试的分层说明见 [`docs/regression-testing.md`](regression-testing.md)。Skill 是可按场景复用的操作规范，不会替代任务清单、人工审核或本地 commit。

## 日常流程

1. 在 `WORK_ITEMS.md` 的 `Inbox` 创建任务，写明目标、验收条件和检查命令。
2. 开始处理时把任务移动到 `Doing`，并保留一个稳定的任务编号。
3. 多步骤任务在当前对话中使用会话级 todo；会话级 todo 只记录本次执行，不替代 `WORK_ITEMS.md`。
4. 如果实现过程中出现重要设计取舍，新增或更新 `docs/decisions/` 下的决策记录。
5. 代码、测试和文档完成后，把任务移动到 `Review`。
6. 执行任务对应的自动检查，并进行人工审核。
7. 审核通过后创建本地 Git commit，再把任务移动到 `Done`。
8. 遇到无法继续的问题，把任务移动到 `Blocked`，记录原因、等待对象和下一步。

完成任务并创建 commit 后，可以把完整条目从 `WORK_ITEMS.md` 的 `Done` 区域迁移到 [`WORK_ITEMS_ARCHIVE.md`](../WORK_ITEMS_ARCHIVE.md)。归档只减少当前任务板的历史内容，不改变任务编号、完成证据或状态规则。

没有必要为每个练习创建分支。需要保留多个实验版本时，使用本地分支；分支提交前仍然以 `WORK_ITEMS.md` 中的任务为准。

## 状态定义

| 状态 | 含义 |
| --- | --- |
| `Inbox` | 已记录但还没有开始处理 |
| `Doing` | 当前正在实现、调查或验证 |
| `Review` | 实现已经完成，等待检查和人工确认 |
| `Blocked` | 暂时不能继续，原因已经记录 |
| `Done` | 自动检查和人工审核均通过，并已提交 |

## 自动检查

根据改动范围选择最小的检查集合：

| 改动范围 | 检查 |
| --- | --- |
| Go 源码 | `gofmt -s -w`、相关 `go test` |
| 普通功能改动 | `go test ./...`、`go vet ./...` |
| API、同步、数据库或跨包改动 | `go test ./...`、`go test ./integration_test/...`、`go vet ./...` |
| CLI 命令、输出或用户可见数据流程 | 相关单测 + `./scripts/regression.sh` |
| 提交前完整验证 | `./scripts/verify.sh` |
| 依赖变化 | `go mod tidy`，然后运行相关测试 |

自动检查通过不等于任务完成。检查命令应写回任务的“检查”字段，必要时在提交信息或工作日志中保留结果。

## 人工审核

个人使用时，人工审核由任务拥有者本人完成，重点检查变更是否符合目标，而不是重复执行所有测试。

### 代码审核

- 阅读 `git diff`，确认每个改动都能对应任务目标。
- 检查是否误改无关文件、配置、数据文件或本地凭据。
- 检查错误处理、边界条件、资源关闭和并发行为。
- 检查用户可见的 CLI 输出是否清晰，新增文本是否经过 i18n。
- 检查数据库、同步和 API 改动是否保持兼容，必要时补测试。
- 检查测试是否验证真实行为，而不是只验证函数被调用。
- 检查是否需要新增或更新 `docs/decisions/` 中的决策记录。

### 结果审核

- 对照任务中的验收条件逐项确认。
- 检查测试输出和实际文件、数据库或 CLI 行为。
- 对破坏性变更、迁移和安全相关改动进行一次手动演练。
- 如果仍有疑问，不要移动到 `Done`，保留在 `Review` 或移到 `Blocked`。

## 设计决策记录

`docs/decisions/` 记录重要决定的背景、取舍和影响。它不替代任务清单，也不替代用户文档：

- `WORK_ITEMS.md` 记录任务状态。
- `docs/decisions/` 记录为什么采用某个方案。
- README 或命令帮助记录用户如何使用功能。
- 测试记录行为是否按预期工作。

以下情况需要写决策记录：

- 选择或更换核心技术方案。
- 修改数据库、配置、文件格式或同步协议。
- 改变 CLI 用户可见行为、默认值或兼容性策略。
- 引入安全、迁移、隐私或数据丢失风险。
- 放弃一个看起来合理、以后可能被重新提出的方案。
- 修改本地开发流程或人工审核规则。

普通拼写修复、局部 Bug 修复和不改变行为的简单重构通常不需要写。任务进入 `Review` 时，如果没有写决策记录，应在任务的“决策”字段写明无需记录的原因。

新增记录时复制 [`docs/decisions/TEMPLATE.md`](decisions/TEMPLATE.md)。已经采用的决定使用 `adopted` 状态；还在考虑的使用 `proposed`；明确放弃的使用 `rejected`。

## 提交规则

提交前确认：

```sh
git diff --check
git status --short
```

提交信息使用简短的 Conventional Commit 风格，例如：

```text
feat: add workspace export
fix: handle duplicate resource keys
refactor: simplify storage routing
docs: describe personal development workflow
```

提交后把任务移动到 `Done`，并在任务下保留实际运行过的关键检查。

## 与 TTL 工作日志的关系

`ttl log write` 适合记录每天做了什么和生成周期总结；`WORK_ITEMS.md` 负责记录任务当前处于哪个状态。两者可以同时使用，但不能用工作日志代替任务状态。
