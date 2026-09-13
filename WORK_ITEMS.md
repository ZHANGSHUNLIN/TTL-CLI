# Work Items

这是本项目的个人开发清单。任务正文、目标、验收、检查和文档链接放在这里，任务状态由同目录的 [`WORK_ITEMS.json`](WORK_ITEMS.json) 维护。这样 Markdown 适合阅读和编辑，JSON 适合管理台和 Agent 更新状态。

不依赖 GitHub Issue、Project、Pull Request 或远程审批。已完成并提交的历史工作项归档在 [`WORK_ITEMS_ARCHIVE.md`](WORK_ITEMS_ARCHIVE.md)；当前任务正文仍以本文件为准，当前状态以 `WORK_ITEMS.json` 为准。

## Tasks

- W-007 实现基础 TUI 本地资源闭环
  - 目标：实现复用 core 用例的 `ttl ui`，只覆盖本地资源浏览、搜索、查看、创建/修改、标签维护和安全删除。
  - 验收：TUI 不解析 CLI 文本；CLI/TUI 使用同一存储契约；空态、错误态、未保存编辑、删除确认和终端恢复可观察；网络或同步不可用不阻塞本地闭环。
  - 检查：`go test ./internal/client/app ./internal/client/cli ./internal/client/tui ./integration_test/...`、`./scripts/regression.sh /tmp/ttl-w007` 和完整 `./scripts/verify.sh` 已通过；PTY 冒烟已观察空态、alternate screen 进入/退出、光标和 bracketed-paste 恢复。仍需 owner 人工完成完整创建/搜索/详情/编辑/标签/删除流程以及窄终端验收。
  - 决策：已放行；沿用 [`docs/decisions/2026-09-12-tui-cli-interaction-strategy.md`](docs/decisions/2026-09-12-tui-cli-interaction-strategy.md) 的基础 TUI / W-007 采用范围。
  - 文档：技术方案：[`docs/tech-designs/2026-09-12-basic-tui.md`](docs/tech-designs/2026-09-12-basic-tui.md)（`reviewed`）；方案评审：[`docs/reviews/2026-09-12-basic-tui-design.md`](docs/reviews/2026-09-12-basic-tui-design.md)（`PASS`）；代码评审：[`docs/reviews/2026-09-12-basic-tui-code.md`](docs/reviews/2026-09-12-basic-tui-code.md)（`CONDITIONAL`）；任务拆分：[`docs/task-breakdowns/2026-09-12-tui-scope-split.md`](docs/task-breakdowns/2026-09-12-tui-scope-split.md)。
  - 备注：W-006/T-01 已完成并提交；W-007 T-02 已实现并完成自动检查与首轮代码评审。任务保持 `Doing`，待补齐 TUI 稳定文案本地化并完成 owner 真实终端完整流程后复审；不依赖候选 `ttl pick`，也未加入远端/同步、工作区、复制、打开、外部编辑器或批量操作。

- W-013 评估并实现 `ttl pick` 快速选择模式（候选）
  - 目标：独立评估短生命周期的资源搜索、选择和输出入口，不与完整 `ttl ui` 绑定。
  - 验收：独立需求明确 query、候选、选择、取消、无结果、输出和退出行为；交互、管道输出和无 TTY 失败路径可测试；不包含 `--copy`/`--open`；延期不阻塞 W-007。
  - 检查：独立需求/设计评审、交互选择测试、管道输出回归和无 TTY 失败路径测试。
  - 决策：范围确认后补充快速选择交互决策记录。
  - 文档：[`docs/requirements/2026-09-12-ttl-pick.md`](docs/requirements/2026-09-12-ttl-pick.md)、[`docs/task-breakdowns/2026-09-12-tui-scope-split.md`](docs/task-breakdowns/2026-09-12-tui-scope-split.md)。
  - 备注：当前处于需求分析阶段，需求文档已补齐，尚未进入技术设计；任务保持 `Doing`，直到需求、方案、评审、WBS、编码和测试全部完成。候选能力依赖 W-006；开始实现前仍需单独确认是否纳入首期以及 `--print` 是否承诺。

- W-014 升级工作流看板任务模型与文档体系
  - 类型：feature
  - 优先级：P1
  - 当前阶段：tests
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：无
  - 目标：让看板区分任务类型、生命周期状态和研发阶段，并展示阶段产物、依赖和阻塞信息，支持 feature、bugfix、docs 等不同工作类型。
  - 验收：旧任务可兼容加载；任务类型、阶段、阶段清单和产物可展示；卡片与详情页显示阶段进度；状态、评审、归档和项目文档功能不回归；流程文档和模板与实际字段一致。
  - 检查：dashboard `go test ./...`、`git diff --check`、无新字段临时项目兼容验证、浏览器手工冒烟。
  - 决策：[`docs/decisions/2026-09-13-workflow-board-model.md`](docs/decisions/2026-09-13-workflow-board-model.md)。
  - 文档：需求：[`docs/requirements/2026-09-13-workflow-board-model.md`](docs/requirements/2026-09-13-workflow-board-model.md)；技术方案：[`docs/tech-designs/2026-09-13-workflow-board-model.md`](docs/tech-designs/2026-09-13-workflow-board-model.md)；方案评审：[`docs/reviews/2026-09-13-workflow-board-model-design.md`](docs/reviews/2026-09-13-workflow-board-model-design.md)（`PASS`）；WBS：[`docs/task-breakdowns/2026-09-13-workflow-board-model.md`](docs/task-breakdowns/2026-09-13-workflow-board-model.md)。
  - 备注：共享 dashboard 实现在 `/Users/v_zhangshun01/.codex/skills/personal-workflow-dashboard`，字段全部可选以兼容其他项目；不增加远程服务、数据库或页面内元数据编辑。

- W-015 为任务类型提供模板并原子创建研发文档
  - 类型：feature
  - 优先级：P1
  - 当前阶段：commit
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-014
  - 产物：需求=[`docs/requirements/2026-09-13-work-item-templates.md`](docs/requirements/2026-09-13-work-item-templates.md)；方案=[`docs/tech-designs/2026-09-13-work-item-templates.md`](docs/tech-designs/2026-09-13-work-item-templates.md)；评审=[`docs/reviews/2026-09-13-work-item-templates-design.md`](docs/reviews/2026-09-13-work-item-templates-design.md)；代码评审=[`docs/reviews/2026-09-13-work-item-templates-code.md`](docs/reviews/2026-09-13-work-item-templates-code.md)（`PASS`）；WBS=[`docs/task-breakdowns/2026-09-13-work-item-templates.md`](docs/task-breakdowns/2026-09-13-work-item-templates.md)；测试=[`docs/tests/2026-09-13-work-item-templates.md`](docs/tests/2026-09-13-work-item-templates.md)；验收=[`docs/acceptance/2026-09-13-work-item-templates.md`](docs/acceptance/2026-09-13-work-item-templates.md)；提交：`39c3e63`（`feat: add workflow task templates`）
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：按 feature、bugfix、docs、chore、spike、refactor、other 模板创建任务时，同时生成需求、技术方案、方案评审、WBS、测试计划和交付验收文档骨架，保证新任务不会只有一句笼统描述。
  - 验收：模板列表和字段契约可观察；创建一次只产生一个任务 ID 和对应文档；任一文档或任务写入失败都不留下半成品；重复标题/路径、非法 slug 和未知模板被拒绝；旧项目仍可读取；页面可选择模板、填写标题和目标并看到生成结果。
  - 检查：dashboard `go test ./...`、HTTP 创建/冲突/回滚测试、`node --check workflow/web/app.js`、项目 `git diff --check`、浏览器创建任务冒烟。
  - 决策：实现前补充 [`docs/decisions/2026-09-13-work-item-templates.md`](docs/decisions/2026-09-13-work-item-templates.md)，记录原子创建和文档槽位策略。
  - 备注：任务创建只写本地项目文件，不自动提交 Git；文档生成内容是可编辑骨架，不能替代人工需求、方案和评审。2026-09-13 owner 已人工确认浏览器新建任务流程和窄屏验收通过；本轮 `gofmt -s -l .`、`go test ./...`、`go vet ./...`、`jq empty WORK_ITEMS.json` 和 `git diff --check` 均通过。代码评审结论为 `PASS`，已提交 `39c3e63`，W-015 完成。

- W-010 建立云端服务独立交付链路
  - 目标：让 `ttl-server` 作为独立应用工程发布、部署、升级和回滚，云端运行环境不依赖 `ttl` 客户端或源码目录。
  - 验收：CI 分别生成可独立下载的 `ttl` 和 `ttl-server` 制品；服务端部署包或最小容器只包含明确的运行文件；配置、密钥、数据卷、日志、健康检查、优雅关闭、TLS 边界和回滚方式有明确约定。
  - 检查：双二进制跨平台构建、服务端依赖边界测试、发布制品内容检查、独立部署冒烟和回滚演练。
  - 决策：实现时更新 `docs/decisions/2026-09-12-separate-client-server-layout.md`。
  - 备注：单仓库和单 Go module 可以保留；能单独编译不作为独立交付完成证据。

- W-012 定义并实现 CLI 可组合性契约
  - 目标：为 `add/get/update/del/tag/dtag` 定义 JSON、非交互、stdin/stdout/stderr、稳定退出码和机器可读错误，服务脚本和自动化场景。
  - 验收：六个核心命令提供版本 1 JSON；既有默认文本语义保持兼容；无 TTY、stdin、输出流和退出码有黑盒证据；未覆盖命令明确拒绝机器模式；不阻塞 W-007。
  - 检查：`go test ./internal/client/app ./internal/client/cli`、`./scripts/cli-composability.sh`、`./scripts/regression.sh`、`go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`bash -n scripts/cli-composability.sh`、`gofmt -s -l .`、`git diff --check` 和完整 `./scripts/verify.sh` 均通过。
  - 决策：[`docs/decisions/2026-09-12-adopt-cli-composability-contract.md`](docs/decisions/2026-09-12-adopt-cli-composability-contract.md)（`adopted`）。
  - 文档：[`docs/requirements/2026-09-12-cli-composability-contract.md`](docs/requirements/2026-09-12-cli-composability-contract.md)、[`docs/tech-designs/2026-09-12-cli-composability-contract.md`](docs/tech-designs/2026-09-12-cli-composability-contract.md)、[`docs/reviews/2026-09-12-cli-composability-contract-design.md`](docs/reviews/2026-09-12-cli-composability-contract-design.md)（`PASS`）、[`docs/reviews/2026-09-12-cli-composability-contract-code.md`](docs/reviews/2026-09-12-cli-composability-contract-code.md)（`PASS`）。
  - 提交：`8b0717d`（`feat: add CLI composability contract`）。
  - 备注：owner 已人工确认，实施、测试、代码评审和本地 commit 均已完成；不扩大到 search、pick、sync、工作空间命令或 W-006 旧 handler 清理。

- W-009 复刻并完善个人版完整研发流程
  - 目标：补齐需求分析、技术方案、方案评审、WBS 拆分、编码、单元测试、代码评审和本地提交阶段，使个人工程具备完整的八阶段研发流程。
  - 验收：`.agents/skills/` 有八阶段个人版 Skill；`docs/templates/` 有需求、方案、方案评审、WBS 和代码评审模板；流程文档、AGENTS、Skill 索引和决策记录保持一致；明确小功能合并阶段的规则。
  - 检查：`git diff --check`、新建 Skill frontmatter 检查、文档链接和内容人工核对。
  - 决策：`docs/decisions/2026-09-12-add-complete-personal-development-flow.md`
  - 备注：上述检查已通过。保持本地 `WORK_ITEMS.md`、`WORK_ITEMS.json`、人工审核和本地 commit；不引入 GitHub PR、远程审批或自动合并。人工审核通过并创建本地 commit 后再把状态更新为 `Done`。

- W-011 建立工程能力边界基线并展开首期需求
  - 目标：将产品范围、客户端/服务端工程边界、能力状态和范围变更门禁固化为全工程基线，并基于基线展开首期本地资源管理需求。
  - 验收：`docs/ttl-product-direction.md` 升级为带版本和生效条件的范围基线；新增 `docs/requirements/2026-09-12-capability-boundary-baseline.md` 和 `docs/ttl-future-roadmap.md`；当前需求只保留首期承诺，候选能力集中到后续规划；文档互相链接。
  - 检查：`git status --short`、`git diff --stat`、`git diff --check` 已通过；人工核对文档链接、章节引用、范围状态和验收标准；未运行 Go 测试，原因：本任务仅修改文档和任务记录。
  - 文档：[`docs/ttl-product-direction.md`](docs/ttl-product-direction.md)、[`docs/ttl-future-roadmap.md`](docs/ttl-future-roadmap.md)、[`docs/requirements/2026-09-12-capability-boundary-baseline.md`](docs/requirements/2026-09-12-capability-boundary-baseline.md)。
  - 决策：无需单独记录，原因：本任务建立范围基线和需求门禁，不改变产品架构、接口、数据格式或实现行为。
  - 备注：基础 TUI / W-007 范围已完成评审确认；CLI 增强、同步增强和其他候选能力不得借本任务直接进入实现。

- W-008 记录开发阶段的客户端与后端数据交互
  - 目标：把当前客户端存储、后端实现、连接方式和同步时机整理成可维护的开发文档。
  - 验收：文档以代码现状为准，覆盖本地、cloud、sync 和显式 push/pull；说明认证、租户隔离、数据位置、失败语义、客户端竞态与已知限制；明确开发阶段行为可直接调整。
  - 检查：人工核对实现路径与命令示例，检查文档链接，运行 `git diff --check`。
  - 决策：无需记录，原因：只记录当前实现与已知缺口，不改变架构、协议或用户行为。
  - 备注：已明确 CLI/TUI 是同一客户端、CI 只是 CLI 的运行环境，并记录单账户不分享资源的边界；竞态作为同一账户多进程或多设备的防御性技术风险保留，不代表多人协作范围。关联设计评审见 `docs/reviews/2026-09-12-tui-client-boundary-design.md`。后续实现变化时继续更新本文档，README 中与现状不一致的示例另行处理。

## Task Format

```md
- W-001 简短任务名称
  - 类型：feature | bugfix | refactor | docs | chore | spike | other
  - 优先级：P0 | P1 | P2 | P3
  - 当前阶段：requirements | design | design_review | breakdown | implementation | tests | delivery_review | commit
  - 阶段清单：按任务实际适用阶段填写，逗号分隔
  - 父任务：W-XXX 或“无”
  - 依赖：W-XXX, W-YYY 或“无”
  - 产物：需求、方案、评审、WBS、测试和提交证据链接
  - 阻塞原因：无法继续时填写，否则写“无”
  - 下一步：下一项可执行动作
  - 目标：要解决的问题或要交付的结果
  - 验收：可以客观判断完成与否的条件
  - 检查：需要运行的命令或需要人工确认的内容
  - 决策：需要的 docs/decisions 记录；没有则写“无需记录，原因：...”
  - 备注：依赖、风险或后续工作
```

状态值和任务 ID 在 [`WORK_ITEMS.json`](WORK_ITEMS.json) 中维护：

```json
{
  "version": 1,
  "items": {
    "W-001": {
      "status": "Inbox"
    }
  }
}
```
