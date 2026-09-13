# Archived Work Items

这里保存已经完成并提交的历史工作项。`WORK_ITEMS.md` 保存当前任务正文，`WORK_ITEMS.json` 保存当前任务状态；本文件只用于查阅历史，不承载进行中的状态。

归档规则：

- 只有完成自动检查、人工审核并创建本地 commit 的工作项才能归档。
- 归档时保留原工作项编号、验收条件、检查结果、决策链接和 commit 备注。
- 新任务和未完成任务继续维护在 [`WORK_ITEMS.md`](WORK_ITEMS.md) 中。

## Done

- [x] W-010 建立云端服务独立交付链路
  - 类型：feature
  - 描述：把 `ttl-server` 从“能编译的服务端入口”准备成可独立构建、发布、部署、探活、升级和回滚的云端服务制品。
  - 优先级：P1
  - 当前阶段：commit
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-004,W-006,W-017
  - 产物：需求=[`docs/requirements/2026-09-13-server-independent-delivery.md`](docs/requirements/2026-09-13-server-independent-delivery.md)；方案=[`docs/tech-designs/2026-09-13-server-independent-delivery.md`](docs/tech-designs/2026-09-13-server-independent-delivery.md)；评审=[`docs/reviews/2026-09-13-server-independent-delivery-design.md`](docs/reviews/2026-09-13-server-independent-delivery-design.md)（`PASS`）；代码评审=[`docs/reviews/2026-09-13-server-independent-delivery-code.md`](docs/reviews/2026-09-13-server-independent-delivery-code.md)（`PASS`）；WBS=[`docs/task-breakdowns/2026-09-13-server-independent-delivery.md`](docs/task-breakdowns/2026-09-13-server-independent-delivery.md)；测试=[`docs/tests/2026-09-13-server-independent-delivery.md`](docs/tests/2026-09-13-server-independent-delivery.md)；验收=[`docs/acceptance/2026-09-13-server-independent-delivery.md`](docs/acceptance/2026-09-13-server-independent-delivery.md)；部署=[`docs/server-deployment.md`](docs/server-deployment.md)；决策=[`docs/decisions/2026-09-13-server-independent-delivery-baseline.md`](docs/decisions/2026-09-13-server-independent-delivery-baseline.md)（`adopted`）；提交：`5cbf3a3`（`feat: establish independent server delivery`）
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：让 `ttl-server` 作为独立应用工程发布、部署、升级和回滚，云端运行环境不依赖 `ttl` 客户端或源码目录。
  - 验收：CI 分别生成可独立下载的 `ttl` 和 `ttl-server` 制品；服务端部署包或最小容器只包含明确的运行文件；配置、密钥、数据卷、日志、健康检查、优雅关闭、TLS 边界和回滚方式有明确约定。
  - 检查：`./scripts/verify.sh`、`scripts/server-delivery-smoke.sh`、`scripts/server-upgrade-rollback-smoke.sh`、双二进制构建、`go test ./internal/server/...`、`git diff --check`、`jq empty WORK_ITEMS.json`、locale JSON 校验、`bash -n scripts/*.sh` 和 CI/release YAML 解析均已通过；远程 Debian/systemd 部署、健康检查、重启后认证 API 和私有地址直连边界已验证，真实 tag 制品清单和外部 TLS 确认待补。
  - 决策：[`docs/decisions/2026-09-13-server-independent-delivery-baseline.md`](docs/decisions/2026-09-13-server-independent-delivery-baseline.md)（`adopted`）；沿用 [`docs/decisions/2026-09-12-separate-client-server-layout.md`](docs/decisions/2026-09-12-separate-client-server-layout.md) 的应用边界。
  - 备注：方案评审、T-01～T-04 实现、自动验证、远程 systemd 部署演练、测试环境验收和 owner 代码评审均已完成；`./scripts/verify.sh` 于 2026-09-13 复跑通过，提交 `5cbf3a3` 后完成收尾并归档。生产外部 TLS 和真实 tag 制品核对属于后续交付项，不阻塞本次测试环境完成结论。项目及 Skill 文档统一使用中文，命令、代码标识、路径、协议字段和文件格式保留原文。单仓库和单 Go module 可以保留；能单独编译不作为独立交付完成证据。

- [x] W-007 实现基础 TUI 本地资源闭环
  - 类型：feature
  - 优先级：P1
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-006
  - 产物：需求=[`docs/requirements/2026-09-12-basic-tui.md`](docs/requirements/2026-09-12-basic-tui.md)；方案=[`docs/tech-designs/2026-09-12-basic-tui.md`](docs/tech-designs/2026-09-12-basic-tui.md)；评审=[`docs/reviews/2026-09-12-basic-tui-design.md`](docs/reviews/2026-09-12-basic-tui-design.md)；WBS=[`docs/task-breakdowns/2026-09-12-tui-scope-split.md`](docs/task-breakdowns/2026-09-12-tui-scope-split.md)；测试=[`docs/tests/2026-09-12-basic-tui.md`](docs/tests/2026-09-12-basic-tui.md)；验收=[`docs/acceptance/2026-09-12-basic-tui.md`](docs/acceptance/2026-09-12-basic-tui.md)；代码评审=[`docs/reviews/2026-09-12-basic-tui-code.md`](docs/reviews/2026-09-12-basic-tui-code.md)（`PASS`）；提交：`f00a9c7`（`feat: complete local resource workflows`）
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：实现复用 core 用例的 `ttl ui`，只覆盖本地资源浏览、搜索、查看、详情页打开、创建/修改、标签维护和安全删除。
  - 验收：TUI 不解析 CLI 文本；CLI/TUI 使用同一存储契约；空态、错误态、未保存编辑、详情页打开成功退出/失败留页、删除确认和终端恢复可观察；网络或同步不可用不阻塞本地闭环。
  - 检查：`go test ./internal/client/app ./internal/client/cli ./internal/client/tui ./integration_test/...`、`go test ./...`、`go vet ./...`、`./scripts/regression.sh /tmp/ttl-w007` 和完整 `./scripts/verify.sh` 已通过；PTY 冒烟已观察空态、alternate screen 进入/退出、光标和 bracketed-paste 恢复；新增长列表分页、平台保存快捷键和详情页打开回归测试；owner 已确认完整创建/搜索/详情/编辑/标签/删除/打开流程以及窄屏验收。
  - 决策：已放行；沿用 [`docs/decisions/2026-09-12-tui-cli-interaction-strategy.md`](docs/decisions/2026-09-12-tui-cli-interaction-strategy.md) 的基础 TUI / W-007 采用范围。
  - 备注：W-006/T-01 已完成并提交；W-007 T-02 已实现并完成自动检查、owner 真实终端完整流程验收、TUI 文案本地化整改和代码评审复审。已提交 `f00a9c7` 并归档；不依赖候选 `ttl pick`，也未加入远端/同步、工作区、复制、外部编辑器或批量操作；详情页 `o` 打开当前资源已纳入本地闭环。

- [x] W-013 实现 `ttl pick` 快速选择模式
  - 类型：feature
  - 优先级：P2
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：W-006
  - 产物：需求=[`docs/requirements/2026-09-12-ttl-pick.md`](docs/requirements/2026-09-12-ttl-pick.md)；方案=[`docs/tech-designs/2026-09-13-ttl-pick.md`](docs/tech-designs/2026-09-13-ttl-pick.md)；评审=[`docs/reviews/2026-09-13-ttl-pick-design.md`](docs/reviews/2026-09-13-ttl-pick-design.md)；WBS=[`docs/task-breakdowns/2026-09-13-ttl-pick.md`](docs/task-breakdowns/2026-09-13-ttl-pick.md)；测试=[`docs/tests/2026-09-13-ttl-pick.md`](docs/tests/2026-09-13-ttl-pick.md)；验收=[`docs/acceptance/2026-09-13-ttl-pick.md`](docs/acceptance/2026-09-13-ttl-pick.md)；代码评审=[`docs/reviews/2026-09-13-ttl-pick-code.md`](docs/reviews/2026-09-13-ttl-pick-code.md)（`PASS`）；提交：`f00a9c7`（`feat: complete local resource workflows`）
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：独立评估短生命周期的资源搜索、选择和输出入口，不与完整 `ttl ui` 绑定。
  - 验收：独立需求明确 query、候选、选择、取消、无结果、输出和退出行为；交互、管道输出和无 TTY 失败路径可测试；不包含 `--copy`/`--open`；`pick` 不记录 history/audit。
  - 检查：`go test ./internal/client/cli`、`./scripts/cli-composability.sh`、`go test ./...`、`go vet ./...`、`./scripts/regression.sh`、`./scripts/verify.sh` 已通过；黑盒、owner TTY 验收和代码评审证据齐全。
  - 决策：无需新增决策记录；v1 契约已冻结在需求、技术方案和方案评审中。
  - 备注：W-013 已完成编码、单测、黑盒回归、owner 真实 TTY 验收和代码评审，已提交 `f00a9c7` 并归档。

- [x] W-014 升级工作流看板任务模型与文档体系
  - 类型：feature
  - 优先级：P1
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：无
  - 目标：让看板区分任务类型、生命周期状态和研发阶段，并展示阶段产物、依赖和阻塞信息，支持 feature、bugfix、docs 等不同工作类型。
  - 验收：旧任务可兼容加载；任务类型、阶段、阶段清单和产物可展示；卡片与详情页显示阶段进度；状态、评审、归档和项目文档功能不回归；流程文档和模板与实际字段一致。
  - 检查：dashboard `gofmt -s -l workflow/*.go`、`go test ./...`、`go vet ./...`、`node --check workflow/web/app.js`、`git diff --check`、`jq empty WORK_ITEMS.json`、运行态 `/health` 和 board API、浏览器手工冒烟、完整 `./scripts/verify.sh` 均通过。
  - 决策：[`docs/decisions/2026-09-13-workflow-board-model.md`](docs/decisions/2026-09-13-workflow-board-model.md)。
  - 文档：需求：[`docs/requirements/2026-09-13-workflow-board-model.md`](docs/requirements/2026-09-13-workflow-board-model.md)；技术方案：[`docs/tech-designs/2026-09-13-workflow-board-model.md`](docs/tech-designs/2026-09-13-workflow-board-model.md)；方案评审：[`docs/reviews/2026-09-13-workflow-board-model-design.md`](docs/reviews/2026-09-13-workflow-board-model-design.md)（`PASS`）；WBS=[`docs/task-breakdowns/2026-09-13-workflow-board-model.md`](docs/task-breakdowns/2026-09-13-workflow-board-model.md)；测试=[`docs/tests/2026-09-13-workflow-board-model.md`](docs/tests/2026-09-13-workflow-board-model.md)；验收=[`docs/acceptance/2026-09-13-workflow-board-model.md`](docs/acceptance/2026-09-13-workflow-board-model.md)；代码评审=[`docs/reviews/2026-09-13-workflow-board-model-code.md`](docs/reviews/2026-09-13-workflow-board-model-code.md)（`PASS`）；提交：`f00a9c7`（`feat: complete local resource workflows`）。
  - 备注：自动检查、owner 看板验收和代码评审已通过。共享 dashboard 实现在 `/Users/v_zhangshun01/.codex/skills/personal-workflow-dashboard`，字段全部可选以兼容其他项目；不增加远程服务、数据库或页面内元数据编辑。已提交 `f00a9c7` 并归档。

- [x] W-016 让 ttl get 默认按 key 和 tag 模糊查询
  - 类型：bugfix
  - 优先级：P1
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：无
  - 产物：需求=[`docs/requirements/2026-09-13-W-016-get-key-only-search.md`](docs/requirements/2026-09-13-W-016-get-key-only-search.md)；方案=[`docs/tech-designs/2026-09-13-W-016-get-key-only-search.md`](docs/tech-designs/2026-09-13-W-016-get-key-only-search.md)；方案评审=[`docs/reviews/2026-09-13-W-016-get-key-only-search-design.md`](docs/reviews/2026-09-13-W-016-get-key-only-search-design.md)；WBS=[`docs/task-breakdowns/2026-09-13-W-016-get-key-only-search.md`](docs/task-breakdowns/2026-09-13-W-016-get-key-only-search.md)；测试=[`docs/tests/2026-09-13-W-016-get-key-only-search.md`](docs/tests/2026-09-13-W-016-get-key-only-search.md)；验收=[`docs/acceptance/2026-09-13-W-016-get-key-only-search.md`](docs/acceptance/2026-09-13-W-016-get-key-only-search.md)；代码评审=[`docs/reviews/2026-09-13-W-016-get-key-only-search-code.md`](docs/reviews/2026-09-13-W-016-get-key-only-search-code.md)（`PASS`）；决策=[`docs/decisions/2026-09-13-get-key-only-search.md`](docs/decisions/2026-09-13-get-key-only-search.md)；提交：`f00a9c7`（`feat: complete local resource workflows`）
  - 阻塞原因：无
  - 下一步：无；已完成并提交
  - 目标：支持通过资源 key 或 tag 获取资源，同时避免 value 命中导致 ttl get 出现非预期候选，并提供显式条件搜索 value
  - 验收：默认查询匹配 key/tag 但不匹配 value；传入 `--value`、`-v` 或 `-val` 时额外匹配 value；既有 TUI 和 pick 搜索行为不回归
  - 检查：`go test ./internal/client/app ./internal/client/cli`、`go test ./...`、`go vet ./...`、`./scripts/regression.sh`、`./scripts/cli-composability.sh`、`./scripts/verify.sh`、`gofmt -s -l`、locale JSON 校验、`git diff --check`；真实 CLI 冒烟已验证默认 key/tag、不命中 value，以及 `--value`、`-v`、`-val` 三种显式 value 搜索入口
  - 决策：[`docs/decisions/2026-09-13-get-key-only-search.md`](docs/decisions/2026-09-13-get-key-only-search.md)
  - 备注：当前修正默认搜索为 key/tag；`--value`、`-v` 和 `-val` 显式扩展到 value。TUI/`pick` 保持原语义。自动检查、tag 命中回归、真实 CLI、owner 验收和代码评审均已完成，已提交 `f00a9c7` 并归档。

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
