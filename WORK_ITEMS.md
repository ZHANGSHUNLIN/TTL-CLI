# Work Items

这是本项目的个人开发清单。任务正文、目标、验收、检查和文档链接放在这里，任务状态由同目录的 [`WORK_ITEMS.json`](WORK_ITEMS.json) 维护。这样 Markdown 适合阅读和编辑，JSON 适合管理台和 Agent 更新状态。

不依赖 GitHub Issue、Project、Pull Request 或远程审批。已完成并提交的历史工作项归档在 [`WORK_ITEMS_ARCHIVE.md`](WORK_ITEMS_ARCHIVE.md)；当前任务正文仍以本文件为准，当前状态以 `WORK_ITEMS.json` 为准。

## Tasks

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

- W-017 重新梳理并收敛客户端与服务端代码结构
  - 类型：refactor
  - 优先级：P1
  - 当前阶段：commit
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-004, W-006
  - 产物：需求=[`docs/requirements/2026-09-13-W-017-code-structure-convergence.md`](docs/requirements/2026-09-13-W-017-code-structure-convergence.md)；方案=[`docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md`](docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md)；方案评审=[`docs/reviews/2026-09-13-W-017-code-structure-convergence-design.md`](docs/reviews/2026-09-13-W-017-code-structure-convergence-design.md)；代码评审=[`docs/reviews/2026-09-13-W-017-code-structure-convergence-code.md`](docs/reviews/2026-09-13-W-017-code-structure-convergence-code.md)（`PASS`）；WBS=[`docs/task-breakdowns/2026-09-13-W-017-code-structure-convergence.md`](docs/task-breakdowns/2026-09-13-W-017-code-structure-convergence.md)；测试=[`docs/tests/2026-09-13-W-017-code-structure-convergence.md`](docs/tests/2026-09-13-W-017-code-structure-convergence.md)；验收=[`docs/acceptance/2026-09-13-W-017-code-structure-convergence.md`](docs/acceptance/2026-09-13-W-017-code-structure-convergence.md)；决策=[`docs/decisions/2026-09-13-code-structure-convergence.md`](docs/decisions/2026-09-13-code-structure-convergence.md)
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：以当前代码为事实重新定义并一次性统一客户端、服务端、核心模型、存储和旧兼容包的最终代码边界，消除新旧结构并存造成的维护混乱。
  - 验收：一次性完成客户端、服务端、核心模型、存储、同步和旧包结构收敛；所有生产/测试消费者迁移，旧包删除，行为和数据格式保持不变，全量回归通过，文档与任务证据同步。
  - 检查：`gofmt -s -l .`、`go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`go test ./internal/architecture`、`./scripts/regression.sh`、`./scripts/cli-composability.sh`、双二进制构建、`./scripts/verify.sh`、`git diff --check` 均通过。
  - 决策：[`docs/decisions/2026-09-13-code-structure-convergence.md`](docs/decisions/2026-09-13-code-structure-convergence.md)（`adopted`，一次性切换与最终归属）；方案评审已更新为 `PASS`。
  - 备注：W-003/W-004/W-006 的历史完成证据不回退；W-017 以当前代码为事实完成一次性结构切换。T-01～T-07 仅表示内部依赖顺序，不产生中间交付物。2026-09-13 owner code review 结论为 `PASS`，实现、测试、交付验收和本地 commit 均已完成；提交：`3ca3dcb`（`refactor: converge client and server package boundaries`）。

- W-018 修复 TUI 打开 Markdown 链接失败
  - 类型：bugfix
  - 优先级：P1
  - 当前阶段：commit
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：无
  - 产物：需求=[`docs/requirements/2026-09-13-W-018-tui-open-markdown-url.md`](docs/requirements/2026-09-13-W-018-tui-open-markdown-url.md)；方案=[`docs/tech-designs/2026-09-13-W-018-tui-open-markdown-url.md`](docs/tech-designs/2026-09-13-W-018-tui-open-markdown-url.md)；方案评审=[`docs/reviews/2026-09-13-W-018-tui-open-markdown-url-design.md`](docs/reviews/2026-09-13-W-018-tui-open-markdown-url-design.md)（`PASS`）；代码评审=[`docs/reviews/2026-09-13-W-018-tui-open-markdown-url-code.md`](docs/reviews/2026-09-13-W-018-tui-open-markdown-url-code.md)（`PASS`）；WBS=[`docs/task-breakdowns/2026-09-13-W-018-tui-open-markdown-url.md`](docs/task-breakdowns/2026-09-13-W-018-tui-open-markdown-url.md)；测试=[`docs/tests/2026-09-13-W-018-tui-open-markdown-url.md`](docs/tests/2026-09-13-W-018-tui-open-markdown-url.md)；验收=[`docs/acceptance/2026-09-13-W-018-tui-open-markdown-url.md`](docs/acceptance/2026-09-13-W-018-tui-open-markdown-url.md)
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：TUI 详情页按 o 打开资源时，支持 Markdown 链接值并在 macOS 正确启动目标 URL。
  - 验收：Markdown 链接值可提取并成功交给平台打开器；纯 URL 行为保持兼容；无法识别的值返回清晰错误并保留详情页；回归测试覆盖 macOS 打开参数和失败状态。
  - 检查：`go test ./internal/client/opener ./internal/client/tui ./internal/client/cli`、`go test ./...`、`go test -race ./internal/client/opener ./internal/client/tui ./internal/client/cli`、`./scripts/regression.sh`、`go build -o /tmp/ttl-w018 ./cmd/ttl`、`go build -o /tmp/ttl-server-w018 ./cmd/ttl-server`、`go vet ./...`、`gofmt -s -l .`、`git diff --check` 均通过；macOS 临时 `open` 命令参数捕获测试通过
  - 决策：无需记录，原因：仅修复客户端打开值的解析，不改变架构、存储、协议或配置格式
  - 备注：根因是平台打开器收到完整 Markdown 字符串而非目标 URL；修复范围限定为客户端 TUI 和 `ttl open` 的输入归一化，平台分支和退出行为保持不变。owner code review 结论为 `PASS`；提交：`8e5057b`（`fix: open markdown links from tui`）。

- W-019 收敛本地与云端存储模式
  - 类型：feature
  - 优先级：P1
  - 当前阶段：Done
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-010
  - 产物：需求=[`docs/requirements/2026-09-13-W-019-storage-model-redesign.md`](docs/requirements/2026-09-13-W-019-storage-model-redesign.md)；方案=[`docs/tech-designs/2026-09-13-W-019-storage-model-redesign.md`](docs/tech-designs/2026-09-13-W-019-storage-model-redesign.md)；评审=[`docs/reviews/2026-09-13-W-019-storage-model-redesign-design.md`](docs/reviews/2026-09-13-W-019-storage-model-redesign-design.md)；WBS=[`docs/task-breakdowns/2026-09-13-W-019-storage-model-redesign.md`](docs/task-breakdowns/2026-09-13-W-019-storage-model-redesign.md)；测试=[`docs/tests/2026-09-13-W-019-storage-model-redesign.md`](docs/tests/2026-09-13-W-019-storage-model-redesign.md)；验收=[`docs/acceptance/2026-09-13-W-019-storage-model-redesign.md`](docs/acceptance/2026-09-13-W-019-storage-model-redesign.md)
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：将客户端当前数据源固定为 local 或 cloud 二选一；普通命令只访问当前数据源；本地存储统一为 SQLite，远程连接可一次配置，切换模式不隐式复制数据。
  - 验收：客户端只暴露 local/cloud 两种存储模式；local 只访问本地 SQLite，cloud 只访问远程服务；服务端每租户使用独立 data.sqlite；远程 profile 可配置多个但当前应用或 workspace 只有一个活动项；sync 不再作为存储类型；切换数据源不自动复制或覆盖另一端；旧 bbolt/SQLite 文件明确拒绝且不读取、不转换、不覆盖、不删除；W-019 不引入 W-020 的同步、版本或 CRDT 协议。
  - 检查：方案评审后执行 `gofmt -s -l .`、`go test ./...`、`go test -race ./...`、`go test ./integration_test/...`、`./scripts/regression.sh`、`./scripts/cli-composability.sh`、`go vet ./...`、`./scripts/verify.sh` 和 `git diff --cached --check` 均通过；远端 `10.99.48.2:8900` 部署、健康检查、认证资源读写、租户 SQLite 权限和删除隔离完成；本地模式隔离和旧格式拒绝演练确认源文件未被修改。
  - 决策：[`服务端按租户使用独立 SQLite 文件`](docs/decisions/2026-09-13-server-tenant-sqlite.md)（`adopted`）；CRDT 同步决策已转交 [`W-020 暂定 CRDT 决策`](docs/decisions/2026-09-13-versioned-crdt-sync.md)（`proposed`）
  - 备注：本任务仅改造 local/cloud 存储模式、客户端配置、服务端租户 SQLite、旧格式拒绝和生命周期；不提供旧数据迁移或兼容窗口。版本化同步、CRDT 和冲突处理拆分到 W-020，不在本任务中讨论或实现。owner code review、方案评审和交付验收结论均为 `PASS`；代码提交：`f9d9af7`（`feat: converge local and cloud storage modes`）。

- W-020 设计本地与云端数据同步及冲突处理
  - 类型：feature
  - 优先级：P1
  - 当前阶段：requirements
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-019
  - 产物：需求=[`docs/requirements/2026-09-13-W-020-local-cloud-data-sync.md`](docs/requirements/2026-09-13-W-020-local-cloud-data-sync.md)；方案=[`docs/tech-designs/2026-09-13-W-020-local-cloud-data-sync.md`](docs/tech-designs/2026-09-13-W-020-local-cloud-data-sync.md)；评审=[`docs/reviews/2026-09-13-W-020-local-cloud-data-sync-design.md`](docs/reviews/2026-09-13-W-020-local-cloud-data-sync-design.md)；WBS=[`docs/task-breakdowns/2026-09-13-W-020-local-cloud-data-sync.md`](docs/task-breakdowns/2026-09-13-W-020-local-cloud-data-sync.md)；测试=[`docs/tests/2026-09-13-W-020-local-cloud-data-sync.md`](docs/tests/2026-09-13-W-020-local-cloud-data-sync.md)；验收=[`docs/acceptance/2026-09-13-W-020-local-cloud-data-sync.md`](docs/acceptance/2026-09-13-W-020-local-cloud-data-sync.md)
  - 阻塞原因：无
  - 下一步：补充需求分析文档并完成 requirements 阶段
  - 目标：在 local 与 cloud 两个独立存储源之间提供版本化双向同步，单次只处理当前 workspace，采用分类型 CRDT 处理可合并变更，并为无法自动判断的冲突提供可观察、可恢复的人工解决流程。
  - 验收：sync 不再作为存储类型；同步可独立访问 local/cloud 且一次只作用于一个 workspace；版本、幂等、增量拉取、OR-Set、Multi-Value Register、tombstone 和冲突解决行为有明确需求、方案、测试和验收证据；未评审通过前不实现具体协议。
  - 检查：待补充
  - 决策：[`docs/decisions/2026-09-13-versioned-crdt-sync.md`](docs/decisions/2026-09-13-versioned-crdt-sync.md)（`proposed`；具体协议和实现待后续方案评审）
  - 备注：本任务承接从 W-019 拆出的版本化同步、CRDT 和冲突处理，当前只保留需求骨架和暂定决策，不在 W-019 实施期间展开。

- W-021 移除当前工程的后端实现
  - 类型：refactor
  - 优先级：P1
  - 当前阶段：commit
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：无
  - 产物：需求=[`docs/requirements/2026-09-18-W-021-remove-in-repo-backend.md`](docs/requirements/2026-09-18-W-021-remove-in-repo-backend.md)；方案=[`docs/tech-designs/2026-09-18-W-021-remove-in-repo-backend.md`](docs/tech-designs/2026-09-18-W-021-remove-in-repo-backend.md)；方案评审=[`docs/reviews/2026-09-18-W-021-remove-in-repo-backend-design.md`](docs/reviews/2026-09-18-W-021-remove-in-repo-backend-design.md)（`PASS`）；代码评审=[`docs/reviews/2026-09-18-W-021-remove-in-repo-backend-code.md`](docs/reviews/2026-09-18-W-021-remove-in-repo-backend-code.md)（`PASS`）；WBS=[`docs/task-breakdowns/2026-09-18-W-021-remove-in-repo-backend.md`](docs/task-breakdowns/2026-09-18-W-021-remove-in-repo-backend.md)；测试=[`docs/tests/2026-09-18-W-021-remove-in-repo-backend.md`](docs/tests/2026-09-18-W-021-remove-in-repo-backend.md)；验收=[`docs/acceptance/2026-09-18-W-021-remove-in-repo-backend.md`](docs/acceptance/2026-09-18-W-021-remove-in-repo-backend.md)；决策=[`docs/decisions/2026-09-18-externalize-backend-service.md`](docs/decisions/2026-09-18-externalize-backend-service.md)
  - 阻塞原因：无
  - 下一步：创建本地 commit 并推送到远端
  - 目标：删除当前仓库中的后端入口、实现、交付链路和专属测试，使 ttl-cli 只保留客户端及其对外部后端服务的远程访问能力
  - 验收：cmd/ttl-server 与 internal/server 不再存在；CI、release、验证脚本不再构建或发布后端；客户端本地和 cloud 模式可构建并通过测试；文档明确后端由独立工程提供
  - 检查：`go test ./internal/client/...` 和完整 `./scripts/verify.sh` 均通过；后者包含格式、脚本语法、客户端构建、架构、CLI 黑盒、`go test ./...`、集成测试、race 和 `go vet ./...`。
  - 决策：[`docs/decisions/2026-09-18-externalize-backend-service.md`](docs/decisions/2026-09-18-externalize-backend-service.md)（`adopted`）
  - 备注：owner 已确认提交并推送。代码评审结论为 `PASS`。真实独立后端未提供，因此未执行跨工程端到端验证；客户端 mock HTTP 契约测试保留。
## Task Format

```md
- W-001 简短任务名称
  - 类型：feature | bugfix | refactor | docs | chore | spike | other
  - 描述：任务背景和要解决的问题
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
