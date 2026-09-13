# 客户端与服务端代码结构收敛任务拆分

日期：2026-09-13  
任务：W-017  
需求：[`docs/requirements/2026-09-13-W-017-code-structure-convergence.md`](../requirements/2026-09-13-W-017-code-structure-convergence.md)  
方案：[`docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md`](../tech-designs/2026-09-13-W-017-code-structure-convergence.md)  
方案评审：[`docs/reviews/2026-09-13-W-017-code-structure-convergence-design.md`](../reviews/2026-09-13-W-017-code-structure-convergence-design.md)

状态：completed；一次性 implementation cutover 已完成

## Dependency Graph

```text
T-01 资源模型与存储契约唯一化
  +-> T-02 客户端命令适配层收敛
  +-> T-03 配置、加密、本地化与文本工具归属
  +-> T-04 客户端同步边界收敛
  -> T-05 旧 db 门面与测试迁移/删除
  -> T-06 架构门禁与项目文档同步
  -> T-07 全量回归与交付验收
```

T-01 至 T-06 是同一次切换内的内部编辑和验证顺序，不是独立交付物、独立主线状态或独立提交。T-02、T-03 和 T-04 依赖 T-01 的 canonical resource/storage 类型，可在同一变更窗口内并行编辑；T-05 在所有消费者已迁移后完成旧包删除；T-06/T-07 统一收口本次切换。

## WBS

### T-01 资源模型与存储契约唯一化

- 目标：让 `internal/core/resource` 和 `internal/core/storage` 成为共享模型与存储契约的唯一来源。
- 输入：当前 `models/` 别名、core 类型、所有生产和测试 import。
- 输出：调用方迁移到 core 包；本次切换结束时删除 `models` 别名，禁止新增旧路径引用。
- 依赖：无。
- 负责区域：`internal/core/`、`models/`、客户端/服务端/存储适配器和相关测试。
- 验收：
  - [x] 所有资源、审计、历史、日志和排序类型的生产引用使用 `internal/core/resource`。
  - [x] `Storage` 方法签名和 JSON tag 行为不变。
  - [x] `go test ./...`、集成测试和两个二进制构建通过。
- 检查：`rg 'ttl-cli/models' --glob '*.go'` 无结果、`go test ./...`、`go vet ./...`。
- 风险或后续：若发现外部兼容承诺，停止整体切换并另建兼容任务，不在主线保留别名。

### T-02 客户端命令适配层收敛

- 目标：消除顶层 `command/` 与 `internal/client/cli` 的双重命令入口，让 CLI 命令由客户端组合根统一装配。
- 输入：`command/*.go`、`internal/client/cli/root.go`、机器模式和 TUI 共用的 app service。
- 输出：`internal/client/cli` 负责 root/lifecycle，`internal/client/cli/commands` 负责 Cobra 命令适配；在本次切换中删除可变包级命令、全局输出 writer 和顶层 `command/`。
- 依赖：T-01。
- 负责区域：`internal/client/cli/`、顶层 `command/`、CLI 测试和回归脚本。
- 验收：
  - [x] `cmd/ttl` 的命令注册只经过 `internal/client/cli`，不再导入顶层 `command`。
  - [x] 命令适配不再持有包级可变 output writer，既有 stdout/stderr、JSON、非交互和退出码语义不变。
  - [x] `go test ./internal/client/cli`、CLI composability 和黑盒回归通过。
- 检查：`rg 'ttl-cli/command' internal cmd --glob '*.go'`、`./scripts/cli-composability.sh`、`./scripts/regression.sh`。
- 风险或后续：命令全局 flag 和历史替换逻辑需要逐条确认，不能通过复制实现留下双写。

### T-03 配置、加密、本地化与文本工具归属

- 目标：按实际消费者收敛顶层 `conf`、`crypto`、`i18n`、`util`，并保持配置、密钥、语言和转义行为不变。
- 输入：顶层包调用关系、嵌入 locale 文件、存储适配器和 server CLI 消费者。
- 输出：`conf` -> `internal/config`；`crypto` -> `internal/crypto`；`i18n` -> `internal/i18n`；无状态文本规则 -> `internal/core/text` 或所属适配器；无调用 helper 删除。
- 依赖：T-01；可与 T-02 并行设计，编码顺序由 import 图决定。
- 负责区域：配置、加密、本地化、文本工具及其测试。
- 验收：
  - [x] 默认配置路径、工作区、加密格式、密钥权限和 locale key 不变。
  - [x] server 不依赖客户端组合根、TUI 或客户端专属状态；共享 `internal/config`、`internal/crypto` 可被存储适配器使用；core 不依赖本地化或适配器。
  - [x] 所有顶层包消费者迁移，相关单测和集成测试通过。
- 检查：locale JSON 校验、`go test ./internal/config ./internal/crypto ./internal/i18n ./internal/...`、`go vet ./...`。
- 风险或后续：若 server CLI 需要本地化，保留共享 `internal/i18n`，不要让 server 反向依赖 client。

### T-04 客户端同步边界收敛

- 目标：把同步 diff、push/pull、交互所需的镜像存储统一到 `internal/client/sync`。
- 输入：根目录 `sync/`、`internal/client/sync/storage.go`、CLI sync command 和同步测试。
- 输出：一个客户端同步包；本次切换结束时删除根目录 `sync/`，CLI 只保留参数和交互选择，core 不包含同步用例。
- 依赖：T-01。
- 负责区域：`internal/client/sync/`、根目录 `sync/`、`internal/client/cli/root.go`、同步测试。
- 验收：
  - [x] `cmd/ttl` 不再导入根目录 `ttl-cli/sync`。
  - [x] `local_only`、`remote_only`、`conflict`、pull/push/auto/dry-run 语义不变。
  - [x] 同步单测、服务端同步集成测试和 CLI 回归通过。
- 检查：`rg 'ttl-cli/sync' --glob '*.go'` 无结果、`go test ./internal/client/sync ./integration_test`。
- 风险或后续：镜像存储的写入顺序和 close 语义不能因移动改变。

### T-05 旧 db 门面与测试迁移/删除

- 目标：在本次切换中迁移仍使用 `db.InitDB`/`db.Stor` 的测试和辅助代码，并删除旧 `db` 包。
- 输入：`db/`、集成测试、server handler 测试、架构测试和当前存储构造器。
- 输出：测试直接使用 `internal/storage`、`internal/server/tenant` 或显式 service/storage；旧全局门面和重复构造器删除。
- 依赖：T-01、T-02、T-03、T-04；仍属于同一次切换。
- 负责区域：`db/`、`integration_test/`、`internal/server/**_test.go`、测试辅助代码。
- 验收：
  - [x] 仓库 Go 代码中无 `ttl-cli/db`、`db.Stor`、`db.InitDB` 引用。
  - [x] 删除 `db/` 后所有测试仍使用临时目录和显式测试 storage/service 生命周期。
  - [x] 资源、加密、导入导出、日志、标签、API 和同步测试通过。
- 检查：`rg 'ttl-cli/db|db\\.Stor|db\\.InitDB' --glob '*.go'`、`go test ./...`、`go test -race ./...`。
- 风险或后续：如存在外部使用承诺，停止物理删除并另建兼容迁移任务。

### T-06 架构门禁与项目文档同步

- 目标：让架构测试阻止旧依赖回流，并让项目导览、拆分方案和任务证据反映真实状态。
- 输入：T-01 至 T-05 的最终 import 图、`PROJECT_OVERVIEW.md`、`docs/client-server-separation-plan.md`。
- 输出：覆盖直接/传递依赖的架构测试；在切换完成后更新目录职责、最终结构、修改入口和 W-017 证据。
- 依赖：T-05。
- 负责区域：`internal/architecture/`、`PROJECT_OVERVIEW.md`、`docs/`、`WORK_ITEMS.md`。
- 验收：
  - [x] client/server/core/legacy 的禁止依赖有自动断言。
  - [x] 文档不再把已迁移或已删除的包描述为当前入口。
  - [x] W-017 文档链路和阶段状态与实际代码一致。
- 检查：`go test ./internal/architecture`、路径链接检查、`git diff --check`。
- 风险或后续：切换完成前不要在文档中把 W-017 的目标结构写成已完成。

### T-07 全量回归与交付验收

- 目标：证明结构迁移没有改变产品行为，并收集进入代码评审/提交阶段的证据。
- 输入：T-01 至 T-06 的代码、测试和文档。
- 输出：完整测试结果、人工验收记录、代码评审输入和最终 acceptance 证据。
- 依赖：T-06。
- 负责区域：全仓库、脚本和 W-017 测试/验收文档。
- 验收：
  - [x] `go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...` 通过。
  - [x] 两个二进制构建、CLI 黑盒回归和架构检查通过。
  - [x] 人工确认 CLI/TUI/API/数据格式/同步语义无未经批准的变化。
- 检查：`./scripts/verify.sh`、`git diff --check`、临时 HOME CLI 冒烟。
- 风险或后续：本任务不包含 W-010 的独立部署验收；一次性切换完成后仍需代码评审和一个本地提交。

## Delivery Gate

- [x] 每个任务都有输入、输出、依赖和验收标准。
- [x] 任务按结构闭环拆分，而不是按单个文件机械拆分。
- [x] 测试、架构门禁、文档和最终验收已纳入。
- [x] 方案评审结论为 `PASS`。
- [x] 负责人确认目标 ownership 和行为不变规则。
- [x] 方案评审通过后执行一次性编码切换。
