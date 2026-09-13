# Personal Development Workflow

本项目采用面向个人练习的本地开发流程。`WORK_ITEMS.md` 保存当前任务正文和交付证据，`WORK_ITEMS.json` 保存任务状态元数据；已完成任务的历史保存在 [`WORK_ITEMS_ARCHIVE.md`](../WORK_ITEMS_ARCHIVE.md)。不使用 GitHub Issue、Project、Pull Request、reviewer 请求或审批积分。

工程 Skill 的索引见 [`docs/engineering-skills.md`](engineering-skills.md)，具体规则位于 [`.agents/skills/`](../.agents/skills/)。回归测试的分层说明见 [`docs/regression-testing.md`](regression-testing.md)。Skill 是可按场景复用的操作规范，不会替代任务清单、人工审核或本地 commit。

## 完整研发流程

大功能按下面的八个阶段推进。`WORK_ITEMS.md` 记录任务正文，`WORK_ITEMS.json` 记录当前状态，阶段文档记录过程产物，Git commit 记录最终交付。

```text
 W-XXX 父任务
  -> personal-requirement-analysis  需求 / PRD
  -> personal-tech-design           技术方案
  -> personal-design-review         方案评审
  -> personal-task-breakdown       WBS 拆分
  -> personal-write-code            逐个任务编码
  -> personal-write-unit-test       补充单元测试
  -> personal-code-review           人工代码审核
  -> personal-commit                本地提交
  -> Done
```

看板状态是任务生命周期，不是一一对应研发阶段的八个状态。一个大功能从开始需求分析起进入 `Doing`，需求、技术方案、方案评审、WBS、编码、单测和验证阶段都保持 `Doing`；只有实现、测试和必需文档完成后才进入 `Review`。因此，阶段进度应通过对应的 `docs/` 产物和任务备注表达，不要把尚未实现的任务提前移动到 `Review`。

| 阶段 | 输入 | 输出 | 放行条件 |
| --- | --- | --- | --- |
| 需求分析 | 沟通记录、当前行为 | `docs/requirements/` 需求文档 | 范围、规则、异常和验收标准明确 |
| 技术设计 | 需求文档、代码和架构约束 | `docs/tech-designs/` 技术方案 | 接口、归属、数据和风险完整 |
| 方案评审 | 需求文档、技术方案 | `docs/reviews/` 方案评审 | `PASS`，或已解决 `CONDITIONAL` |
| WBS 拆分 | 已通过的技术方案 | `docs/task-breakdowns/` 任务拆分 | 每项有输入、输出、依赖和验收 |
| 编码 | 一个 `T-XX` 任务 | 代码 diff、变更说明 | 只实现当前任务，验收条件可验证 |
| 单测 | 代码 diff、任务验收 | 包级行为测试 | 正常、边界和失败路径有证据 |
| 代码评审 | 代码、测试、阶段文档 | `docs/reviews/` 代码评审 | `PASS`，或已解决所有阻塞项 |
| 本地提交 | 已审核 diff、检查结果 | Git commit、任务完成记录 | commit 成功后才进入 `Done` |

阶段模板位于 [`docs/templates/`](templates/)。文档文件名使用 `YYYY-MM-DD-<slug>.md`，并在 `WORK_ITEMS.md` 中互相链接。

## 从模板创建任务

运行中的 personal workflow dashboard 提供“新建任务”入口。选择任务类型并填写标题、目标等最小信息后，服务端会为任务分配 `W-XXX` 编号，创建 `Inbox` 任务，并同时生成以下六个可编辑文档骨架：

```text
docs/requirements/YYYY-MM-DD-W-XXX-<slug>.md
docs/tech-designs/YYYY-MM-DD-W-XXX-<slug>.md
docs/reviews/YYYY-MM-DD-W-XXX-<slug>-design.md
docs/task-breakdowns/YYYY-MM-DD-W-XXX-<slug>.md
docs/tests/YYYY-MM-DD-W-XXX-<slug>.md
docs/acceptance/YYYY-MM-DD-W-XXX-<slug>.md
```

模板只决定默认类型、阶段清单和提示语，不代表任何阶段已经完成。`feature`、`bugfix`、`refactor`、`docs`、`chore`、`spike` 和 `other` 都会生成完整槽位；轻量任务须在对应文档中明确“不适用”及理由。创建失败会整体拒绝，不能留下半个任务或孤立文档。

小功能可以合并阶段，但必须满足同样的门禁。例如一个单文件的局部 Bug 修复，可以把需求、方案和拆分压缩到一个 `WORK_ITEMS.md` 条目中；如果涉及多个包、公共接口、存储、同步、迁移或用户流程，不应跳过需求、方案和评审。

## 日常流程

1. 在 `WORK_ITEMS.md` 的 `Tasks` 区域创建任务，写明目标、验收条件和检查命令；在 `WORK_ITEMS.json` 增加同一个任务 ID 和初始状态 `Inbox`。
2. 开始处理时把 JSON 中的任务状态更新为 `Doing`，并保留一个稳定的任务编号。
3. 多步骤任务在当前对话中使用会话级 todo；会话级 todo 只记录本次执行，不替代 `WORK_ITEMS.md` 和 `WORK_ITEMS.json`。
4. 大功能按“需求 → 方案 → 方案评审 → WBS”推进，未通过方案评审不得开始实现。
5. 每次只实现一个已拆分的 `T-XX` 任务；必要时把独立子任务另列为 `W-XXX`。
6. 如果实现过程中出现重要设计取舍，新增或更新 `docs/decisions/` 下的决策记录。
7. 代码、测试和文档完成后，把 JSON 中的任务状态更新为 `Review`。
8. 执行任务对应的自动检查，并进行人工审核。
9. 审核通过后创建本地 Git commit，再把 JSON 中的任务状态更新为 `Done` 并把完成标记设为 `true`。
10. 遇到无法继续的问题，把 JSON 中的任务状态更新为 `Blocked`，并在 Markdown 备注中记录原因、等待对象和下一步。

完成任务并创建 commit 后，可以把完整条目从 `WORK_ITEMS.md` 的 `Tasks` 区域迁移到 [`WORK_ITEMS_ARCHIVE.md`](../WORK_ITEMS_ARCHIVE.md)，并从 `WORK_ITEMS.json` 移除对应状态。归档只减少当前任务板的历史内容，不改变任务编号、完成证据或状态规则。

没有必要为每个练习创建分支。需要保留多个实验版本时，使用本地分支；分支提交前仍然以 `WORK_ITEMS.md` 中的任务和 `WORK_ITEMS.json` 中的状态为准。

## 状态定义

| 状态 | 含义 |
| --- | --- |
| `Inbox` | 已记录但还没有开始处理 |
| `Doing` | 已开始处理，正在进行需求分析、方案设计、评审、拆分、实现、测试或验证 |
| `Review` | 实现已经完成，等待检查和人工确认 |
| `Blocked` | 暂时不能继续，原因已经记录 |
| `Done` | 自动检查和人工审核均通过，并已提交 |

`Review` 是“等待人工审核”的状态，不代表审核已经通过。`Done` 必须同时具备自动检查证据、人工审核结论和 commit id；`WORK_ITEMS.json` 中的 `completed` 应同步为 `true`。

任务的 `reviews` 数组可以保存过程中的普通评论，但普通评论不会自动改变任务状态，也不等同于代码评审。正式代码评审应在任务进入 `Review` 后进行；只有 `Review` 状态的任务才允许使用“打回修改”动作，打回后回到 `Doing`。

任务类型和阶段进度属于任务内容元数据，写在 `WORK_ITEMS.md` 的可选字段中，由看板解析和展示。推荐字段为 `类型`、`优先级`、`当前阶段`、`阶段清单`、`父任务`、`依赖`、`产物`、`阻塞原因` 和 `下一步`。生命周期状态仍只写入 `WORK_ITEMS.json`；不要用阶段字段替代 `Inbox`、`Doing`、`Review`、`Blocked`、`Done`。

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

自动检查通过不等于任务完成。检查命令应写回 `WORK_ITEMS.md` 任务的“检查”字段，必要时在提交信息或工作日志中保留结果。

## 人工审核

个人使用时，人工审核由任务拥有者本人完成，重点检查变更是否符合目标，而不是重复执行所有测试。

### 代码审核

- 阅读 `git diff`，确认每个改动都能对应任务目标。
- 检查是否误改无关文件、配置、数据文件或本地凭据。
- 检查错误处理、边界条件、资源关闭和并发行为。
- 检查用户可见的 CLI 输出是否清晰，新增文本是否经过 i18n。
- 检查数据库、同步和 API 改动是否同步更新全部调用方、当前数据和测试。
- 检查测试是否验证真实行为，而不是只验证函数被调用。
- 检查是否需要新增或更新 `docs/decisions/` 中的决策记录。

### 结果审核

- 对照任务中的验收条件逐项确认。
- 检查测试输出和实际文件、数据库或 CLI 行为。
- 对破坏性变更、迁移和安全相关改动进行一次手动演练。
- 如果仍有疑问，不要移动到 `Done`，保留在 `Review` 或移到 `Blocked`。

## 设计决策记录

`docs/decisions/` 记录重要决定的背景、取舍和影响。它不替代任务清单，也不替代用户文档：

- `WORK_ITEMS.md` 记录任务正文和交付证据。
- `WORK_ITEMS.json` 记录任务状态和完成标记。
- `docs/decisions/` 记录为什么采用某个方案。
- README 或命令帮助记录用户如何使用功能。
- 测试记录行为是否按预期工作。

以下情况需要写决策记录：

- 选择或更换核心技术方案。
- 修改数据库、配置、文件格式或同步协议。
- 改变 CLI 用户可见行为、默认值或契约变更策略。
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

提交后把 `WORK_ITEMS.json` 中的任务更新为 `Done`，并在 `WORK_ITEMS.md` 的任务下保留实际运行过的关键检查。

## 与 TTL 工作日志的关系

`ttl log write` 适合记录每天做了什么和生成周期总结；`WORK_ITEMS.md` 负责记录任务正文，`WORK_ITEMS.json` 负责记录任务当前处于哪个状态。两者可以同时使用，但不能用工作日志代替任务状态。
