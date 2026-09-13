# TUI 与 CLI 范围拆分任务分解

日期：2026-09-12
任务：W-007、W-012、W-013
需求：[`docs/requirements/2026-09-12-capability-boundary-baseline.md`](../requirements/2026-09-12-capability-boundary-baseline.md)
方案：[`docs/tui-master-plan.md`](../tui-master-plan.md)（v0.2，提案）
方案评审：[`docs/reviews/2026-09-12-tui-client-boundary-design.md`](../reviews/2026-09-12-tui-client-boundary-design.md)（`PASS`，仅基础 TUI / W-007 范围）

## 拆分原则

本拆分以功能闭环为单位，不按文件名拆任务。基础 TUI、通用客户端前置、CLI 可组合性和快速选择模式拥有不同的用户目标、验收方式和延期条件，因此分别管理。

基线约束：基础 `ttl ui` 本地资源闭环和已独立批准的 W-012/W-013 属于各自工作项的阶段承诺；其他路线图能力仍是候选。W-013 不阻塞基础 TUI。

## Dependency Graph

```text
W-006 / T-01 通用客户端前置
              |
              v
       W-007 / T-02 基础 ttl ui

W-012 / T-03 CLI 可组合性（已独立交付）
W-013 / T-04 ttl pick（已批准，实施中）
              ^
              |
       T-01 可复用的客户端服务契约
```

T-03 和 T-04 可以在 T-01 完成后分别立项；它们不依赖 T-02，也不能反向阻塞 T-02。是否从候选转为阶段承诺，必须先更新能力基线并完成独立需求和方案评审。

## WBS

### T-01 通用客户端服务与显式存储生命周期

状态：已完成，由 W-006 跟踪并提交（`c58c644`）。

- 目标：为 CLI、TUI 和后续交互入口提供可复用的资源用例和显式存储生命周期，消除入口之间复制业务规则的必要。
- 输入：W-006、基线中的客户端工程边界、现有 CLI 行为和 `Storage` 契约。
- 输出：可由多个客户端入口复用的资源/搜索服务边界；客户端命令不直接依赖全局存储门面；既有 CLI 行为保持兼容。
- 依赖：无；W-006 是项目级跟踪项。
- 负责区域：`internal/client`、`internal/core`、客户端命令组装和存储生命周期。
- 验收：
  - [x] `command`、`internal/client/cli` 和 `internal/client/sync` 不直接读取全局 `db.Stor`。
  - [x] 资源创建、查看、修改、标签、删除和 key/value/tag 搜索由可复用服务承载。
  - [x] 既有 `add/get/update/del/tag` 核心 CLI 语义不变。
  - [x] 本地存储失败、资源不存在和重复 key 等错误仍能被入口识别。
- 检查：相关单元测试、`go test ./integration_test/...`、`./scripts/regression.sh`、`go test -race ./...` 和依赖方向检查。
- 证据：[`docs/reviews/2026-09-12-client-storage-lifecycle-code.md`](../reviews/2026-09-12-client-storage-lifecycle-code.md)（`PASS`）；旧 `db` 包仍被历史测试使用，但已退出客户端生产依赖图。
- 风险或后续：不在本任务中设计新的 `ttl pick` 契约；该候选能力仍由 T-04 独立跟踪。

### T-02 基础 `ttl ui` 本地资源闭环

状态：实现和自动检查已完成，等待交付验收；W-007 保持 `Doing`。

- 目标：实现基线阶段承诺的 TUI 首版，只覆盖本地资源浏览、按 key/value/tag 搜索、详情、创建/修改、标签维护和安全删除。
- 输入：T-01 输出、基线的阶段承诺、现有资源字段和 TUI 交互约束。
- 输出：可启动的 `ttl ui`；用户可以在本地完成资源闭环；CLI 与 TUI 操作同一份数据。
- 依赖：T-01；不依赖 T-03 或 T-04。
- 负责区域：`cmd/ttl`、`internal/client/tui`、共享 application/core 契约和 TUI 聚焦测试。
- 验收：
  - [x] 空数据库、搜索无结果、内容过长、未保存编辑、删除确认和数据库错误均有可观察处理。
  - [x] TUI 不解析 CLI 文本输出，也不复制 CLI 业务规则。
  - [x] 创建、查看、修改、标签维护和删除后的结果可以通过 CLI 或重新进入 TUI 观察到。
  - [x] 本地模式不依赖网络、API Key 或同步服务可用。
  - [x] 删除和保存失败不会静默造成数据丢失，终端退出状态可恢复。
- 检查：TUI model/update 测试、相关 service/storage 测试、CLI 黑盒回归、窄终端和异常退出人工验收。
- 证据：[`docs/tests/2026-09-12-basic-tui.md`](../tests/2026-09-12-basic-tui.md)；[`docs/acceptance/2026-09-12-basic-tui.md`](../acceptance/2026-09-12-basic-tui.md)；[`docs/reviews/2026-09-12-basic-tui-code.md`](../reviews/2026-09-12-basic-tui-code.md)（`CONDITIONAL`）。长列表分页、最后一页页码和平台保存快捷键已有回归测试。
- 风险或后续：工作区管理、同步页面、数据管理、审计、日志、批量操作、复制和外部编辑器不属于本任务；详情页打开当前资源已纳入本任务，Linux 平台仍不支持。

### T-03 CLI 可组合性契约

状态：已完成，由 W-012 跟踪并提交（`8b0717d`）。

- 目标：在基础 TUI 不受影响的前提下，为脚本和自动化场景定义并实现 JSON、非交互、stdin/stdout、稳定退出码和机器可读错误。
- 输入：T-01 的结构化服务结果、基线候选能力状态、现有 CLI 黑盒行为。
- 输出：独立的 CLI 契约需求、兼容性说明、实现和黑盒回归；不改变既有核心命令默认语义。
- 依赖：T-01；不依赖 T-02。开始实现前需要基线范围变更确认和独立方案评审。
- 负责区域：`internal/client/cli`、CLI 输出/错误适配、黑盒回归和帮助文档。
- 验收：
  - [x] 需求明确列出 JSON schema、退出码、stdin/stdout 和非交互行为。
  - [x] stdout 只输出机器结果，诊断信息与调试信息不混入结果流。
  - [x] 无 TTY 和交互终端的行为分别有黑盒覆盖。
  - [x] 既有 CLI 核心命令在未启用新模式时保持兼容。
- 检查：独立 CLI 需求/设计评审、`./scripts/regression.sh`、无 TTY 黑盒、错误路径和兼容性回归。
- 证据：[`docs/reviews/2026-09-12-cli-composability-contract-code.md`](../reviews/2026-09-12-cli-composability-contract-code.md)（`PASS`）；提交 `8b0717d`。
- 风险或后续：不得把该任务的契约反向写入 TUI 首版依赖；如果契约设计影响现有命令，必须单独更新基线和兼容决策。

### T-04 `ttl pick` 快速选择模式

状态：已完成编码和自动验证；W-013 进入交付评审前准备。

- 目标：评估并实现短生命周期的资源搜索/选择入口，服务单条资源快速取用，不把它与完整 `ttl ui` 绑定。
- 输入：T-01 的资源搜索服务、W-013 独立需求和通过的技术方案评审。
- 输出：`ttl pick [query]` 实现、单元测试和黑盒回归；v1 不提供 `--print`。
- 依赖：T-01；不依赖 T-02 或 T-03。
- 负责区域：客户端入口、快速选择 TUI、结构化资源输出和交互回归。
- 验收：
  - [x] 需求明确用户输入、选择、取消、无结果和退出行为。
  - [x] 快速选择模式不打开完整设置、工作区或同步页面。
  - [x] `--copy`、`--open` 不因本任务自动进入范围。
  - [x] 该模式与 `ttl ui` 共享服务契约，不复制资源业务规则。
- 检查：`go test ./internal/client/cli`、`./scripts/cli-composability.sh`、`./scripts/regression.sh`、`go test ./...` 和 `go vet ./...`。
- 风险或后续：v1 不提供 `--print`；需要结构化输出或复制/打开能力时另立需求。

## Delivery Gate

- [x] 每个任务都有独立目标、输入、输出、依赖、负责区域和验收标准。
- [x] T-02 不依赖候选 T-03/T-04，候选延期不会阻塞基础 TUI。
- [x] T-04 在实现前完成需求、技术方案和方案评审，并更新 W-013 的阶段承诺状态；T-03 已由 W-012 完成该门禁。
- [x] 没有把工作区、同步、批量、复制、外部编辑器、Web UI 或多设备一致性混入 T-02；详情页 `o` 打开属于本地资源闭环增量。
- [ ] T-02 尚缺 owner 真实终端验收；T-04 自动交付证据已闭合，仍需 owner 交付评审后单独完成。
