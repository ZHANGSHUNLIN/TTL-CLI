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

- W-010 建立云端服务独立交付链路
  - 类型：feature
  - 描述：把 `ttl-server` 从“能编译的服务端入口”准备成可独立构建、发布、部署、探活、升级和回滚的云端服务制品。
  - 优先级：P1
  - 当前阶段：design_review
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-004,W-006,W-017
  - 产物：需求=[`docs/requirements/2026-09-13-server-independent-delivery.md`](docs/requirements/2026-09-13-server-independent-delivery.md)；方案=[`docs/tech-designs/2026-09-13-server-independent-delivery.md`](docs/tech-designs/2026-09-13-server-independent-delivery.md)；评审=[`docs/reviews/2026-09-13-server-independent-delivery-design.md`](docs/reviews/2026-09-13-server-independent-delivery-design.md)；WBS=[`docs/task-breakdowns/2026-09-13-server-independent-delivery.md`](docs/task-breakdowns/2026-09-13-server-independent-delivery.md)；测试=[`docs/tests/2026-09-13-server-independent-delivery.md`](docs/tests/2026-09-13-server-independent-delivery.md)；验收=[`docs/acceptance/2026-09-13-server-independent-delivery.md`](docs/acceptance/2026-09-13-server-independent-delivery.md)；部署=[`docs/server-deployment.md`](docs/server-deployment.md)；决策=[`docs/decisions/2026-09-13-server-independent-delivery-baseline.md`](docs/decisions/2026-09-13-server-independent-delivery-baseline.md)（`proposed`）；提交：待完成
  - 阻塞原因：无（等待方案评审门禁）
  - 下一步：完成方案评审，确认 Linux/systemd 基线、TLS 终止边界、健康检查语义和升级回滚策略
  - 目标：让 `ttl-server` 作为独立应用工程发布、部署、升级和回滚，云端运行环境不依赖 `ttl` 客户端或源码目录。
  - 验收：CI 分别生成可独立下载的 `ttl` 和 `ttl-server` 制品；服务端部署包或最小容器只包含明确的运行文件；配置、密钥、数据卷、日志、健康检查、优雅关闭、TLS 边界和回滚方式有明确约定。
  - 检查：需求/方案/方案评审文档人工核对；核对 CI 当前入口与服务端生命周期现状；方案阶段完成后运行双二进制构建、服务端依赖边界、发布制品内容、独立部署冒烟和回滚演练。
  - 决策：[`docs/decisions/2026-09-13-server-independent-delivery-baseline.md`](docs/decisions/2026-09-13-server-independent-delivery-baseline.md)（`proposed`）；沿用 [`docs/decisions/2026-09-12-separate-client-server-layout.md`](docs/decisions/2026-09-12-separate-client-server-layout.md) 的应用边界。
  - 备注：需求和技术方案已按现状补齐，明确 Linux amd64/arm64 服务端 tarball、独立 CI 制品、`/healthz`、优雅关闭、外部 TLS 终止和指针式回滚；当前尚未实现、评审或验收。项目及 Skill 文档统一使用中文，命令、代码标识、路径、协议字段和文件格式保留原文。单仓库和单 Go module 可以保留；能单独编译不作为独立交付完成证据。

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
