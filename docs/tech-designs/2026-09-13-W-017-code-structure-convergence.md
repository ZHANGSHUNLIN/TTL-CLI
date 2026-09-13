# 客户端与服务端代码结构一次性收敛技术方案

日期：2026-09-13  
任务：W-017  
需求：[`docs/requirements/2026-09-13-W-017-code-structure-convergence.md`](../requirements/2026-09-13-W-017-code-structure-convergence.md)  
状态：draft

## Current Behavior

仓库已经有两个正式可执行入口：`cmd/ttl` 调用 `internal/client/cli`，`cmd/ttl-server` 调用 `internal/server/cli`。客户端服务生命周期由 `internal/client/app.Service` 管理，SQLite/bbolt 位于 `internal/storage`，共享 `Storage` 接口位于 `internal/core/storage`。

迁移没有完成的部分集中在旧包并存和依赖方向：

- `internal/client/cli/root.go` 仍注册顶层 `command/` 的全局 Cobra 命令，并直接依赖 `models`、`conf`、`i18n` 和根目录 `sync/`。
- `models/` 只是 `internal/core/resource` 的别名，但生产包、服务端和测试仍普遍引用它。
- 根目录 `sync/` 同时承载 diff/push/pull，`internal/client/sync/` 承载镜像存储，边界不统一。
- `db/` 已不在客户端主生产链路中，但集成测试和服务端测试仍使用其全局门面和旧构造器。
- `conf/`、`crypto/`、`i18n/`、`util/` 仍是顶层包，跨客户端、服务端和存储适配器的归属没有收敛。
- 当前架构测试只禁止 client/server 互相依赖和部分旧包依赖，没有覆盖所有目标淘汰包。

本方案只做一次结构切换和依赖统一，不改变 CLI/TUI/API、数据文件或同步协议；切换前后不把新旧结构作为两个长期版本维护。

## Cutover Result

一次性切换已完成：客户端命令位于 `internal/client/cli/commands`，配置和加密基础位于共享的 `internal/config`、`internal/crypto`，资源模型和存储契约由 `internal/core` 唯一提供，同步实现统一到 `internal/client/sync`。旧顶层 `command`、`models`、`db`、`sync`、`conf`、`crypto`、`i18n` 和 `util` 包已删除，测试改为使用具体适配器或显式 service/storage 生命周期。`ttl` 与 `ttl-server` 的传递依赖和全量行为检查已通过。

## Proposed Ownership

| Area | Owner | Responsibility |
| --- | --- | --- |
| `cmd/ttl` | 客户端入口 | 创建客户端 root command 并执行一次客户端生命周期 |
| `internal/client/cli` | 客户端组合根 | 持久化 flags、服务注入、命令注册、退出和错误映射 |
| `internal/client/cli/commands` | 客户端命令适配 | Cobra 参数解析、调用 app service、输出和本地化；不拥有存储生命周期 |
| `internal/client/app` | 客户端用例 | 资源、审计、历史、日志、迁移和存储选择；不依赖 Cobra |
| `internal/client/tui` | TUI 适配 | 终端状态、交互和 app service 调用；不解析 CLI 输出 |
| `internal/client/remote` | HTTP 客户端 | `/api/v1` 请求、响应和远端存储适配 |
| `internal/client/sync` | 客户端同步 | diff、push、pull、镜像存储和同步交互；不放入共享 core |
| `internal/config` | 共享配置基础 | INI、工作区和本地路径解析；客户端与 bbolt 适配器共用，不依赖入口 |
| `internal/crypto` | 共享加密基础 | 本地密钥和资源加解密；客户端与 bbolt 适配器共用，不进入 server API |
| `internal/i18n` | 共享本地化基础 | CLI、TUI 和 `ttl-server` 运维命令共用的本地化加载；语言文件随包移动 |
| `internal/core/resource` | 共享领域模型 | 资源、标签、审计、历史、日志和排序类型的唯一来源 |
| `internal/core/storage` | 共享存储契约 | `Storage` 接口及其稳定行为约束 |
| `internal/core/text` | 共享纯文本规则 | 去重、忽略大小写匹配、转义等无状态规则；不依赖入口或存储 |
| `internal/storage/sqlite` | SQLite 适配器 | SQLite 连接、表结构和资源持久化 |
| `internal/storage/bbolt` | bbolt 适配器 | bbolt 连接、加密读写和资源持久化 |
| `internal/server/api` | 服务端 HTTP | 路由、鉴权上下文、DTO 和 handler |
| `internal/server/tenant` | 服务端租户 | 用户、API Key、租户存储生命周期 |
| `internal/server/cli` | 服务端入口 | `serve` 和用户运维命令；不依赖 client |

### Canonical dependency direction

```text
cmd/ttl -> internal/client/cli -> internal/client/cli/commands
                         -> internal/client/app -> internal/core/*
                         -> internal/config|crypto|remote|sync

cmd/ttl-server -> internal/server/cli -> internal/server/api|tenant
                                     -> internal/core/*

internal/storage/* -> internal/core/*
internal/config|crypto -> no client/server dependency
internal/i18n          -> no client/server/core dependency
internal/core/*   -> standard library only
```

顶层 `command/`、`models/`、`db/`、根目录 `sync/`、`conf/`、`crypto/` 和 `i18n/` 在本次切换完成后不再作为生产代码入口保留。`util/` 不整体搬迁，按函数消费者拆到 `internal/core/text` 或所属适配器；无调用的 debug helper 直接删除。

## Interfaces

### Stable contracts

- `internal/core/resource` 的类型字段、JSON tag 和 `internal/core/storage.Storage` 方法签名保持现状。
- `internal/client/app.Service` 继续通过 context 注入命令，TUI 继续依赖窄化的 `ResourceService`。
- `internal/client/remote` 保持现有 `/api/v1` 路径、认证头、响应包装和错误转换。
- `ttl`、`ttl-server` 的命令、参数、输出和退出码保持兼容。

### Command registration

`internal/client/cli` 只保留 root command、持久化生命周期和少量客户端专属命令。原 `command/*.go` 分批迁入 `internal/client/cli/commands`，每个注册函数返回独立 `*cobra.Command` 或命令组；不再暴露可变的包级 `AddCmd`、`OutputWriter` 或共享 flag 状态。命令输出统一写入 `cmd.OutOrStdout()`/`cmd.ErrOrStderr()`。

### Sync ownership

将根目录 `sync` 的 `DiffResult`、比较和执行函数迁入 `internal/client/sync`，与现有 `MirroredStorage` 一起形成客户端同步包。`internal/client/cli` 只负责参数和交互选择，不能直接实现 diff 规则；`internal/core` 只保留资源类型和存储接口。

### Compatibility policy

本次采用一次性切换，不保留跨交付周期的兼容 alias。实现时可以按依赖顺序编辑文件，但最终变更必须同时完成所有仓库内消费者迁移和旧包删除；禁止新增引用旧路径。

## Data And Configuration

- 不修改 SQLite 表、bbolt bucket、JSON 字段、INI 字段、加密格式或 HTTP DTO。
- `conf` 迁移到 `internal/config` 时保持默认路径、工作区解析和配置键不变。
- `crypto` 迁移到 `internal/crypto` 时保持密钥路径、AES-GCM 格式和错误语义不变。
- `i18n` 迁移到 `internal/i18n` 时保持 locale key、语言选择和嵌入资源内容不变。
- 远端服务独立交付、部署配置、TLS、数据卷和回滚仍由 W-010 负责。
- 如任何迁移发现必须改变数据或配置格式，立即停止该步骤，另建迁移需求和决策记录。

## Lifecycle And Failure Behavior

- **Creation**：`cmd/ttl` 创建 root；`internal/client/cli` 在 pre-run 中打开一个 app service；server 入口只创建自己的 API/tenant 资源。
- **Mutation**：命令和 TUI 通过 app/service 或明确的 sync/remote 接口执行；适配层不直接打开数据库。
- **Cleanup**：客户端 root 执行结束关闭 service；迁移测试中的临时存储由测试负责关闭；整体切换结束前完成旧包消费者清零和删除。
- **Failure and recovery**：整体切换在合入前必须完成编译、架构、行为和文档检查；import cycle、行为回归或数据格式变化都回滚整个切换，不保留半迁移状态。
- **Concurrency**：不引入新的全局可变状态；现有单进程 service 生命周期和存储锁语义保持不变。

## Security And Compatibility

- 不改变 API Key、租户隔离、密钥文件权限或本地加密行为。
- `internal/server` 不得依赖客户端组合根、TUI、命令输出或客户端专属工作区状态；共享 `internal/config`、`internal/crypto` 基础包可被存储适配器使用。
- 迁移后 server binary 的传递依赖不得包含 `internal/client`、`command`、`db` 或根目录 `sync`。
- 兼容 alias 不作为新公共接口；删除前必须迁移所有仓库内测试和调用方。

## Alternatives

- **继续保留新旧两套结构**：短期改动小，但新功能会继续产生错误归属，拒绝。
- **无设计、无依赖顺序的一次性重写全部目录**：风险集中，难以定位行为回归和 import cycle，拒绝；本次仍采用一次性切换，只在同一变更窗口内按依赖顺序编辑和验证。
- **拆成客户端/服务端两个仓库或 Go module**：超过本任务目标，增加版本和发布联动，保留单仓库/单 module。
- **把所有代码塞进 `internal/client`**：仍会混合入口、用例和适配器，拒绝；采用 client/commands、app、adapter 分层。

## Test And Regression Plan

| Layer | Evidence |
| --- | --- |
| Unit | 切换过程中保留原包行为测试；新增 core、commands、config、sync 和 adapter 的边界测试 |
| Architecture | `go test ./internal/architecture` 检查 client/server/core/legacy import 方向和传递依赖 |
| Integration | `go test ./integration_test/...` 覆盖资源、加密、导入导出、日志、标签和同步 |
| CLI | `./scripts/regression.sh`、`./scripts/cli-composability.sh` 验证命令、输出、退出码和临时 HOME |
| Build | `go build -o ttl ./cmd/ttl`、`go build -o ttl-server ./cmd/ttl-server`，检查二进制边界 |
| Full | `go test -race ./...`、`go vet ./...`、`./scripts/verify.sh`、`git diff --check` |

## Rollout And Recovery

按 WBS 由契约到适配器的依赖顺序在同一次结构切换中实施。T-01 至 T-06 是内部执行步骤，不是独立交付物，也不在主线保留半迁移状态；T-07 统一完成全量验证。若出现行为、数据或 API 回归，回滚整个 W-017 切换，不回滚已经完成的 W-003/W-004/W-006 历史提交。
