# TTL 客户端与后端服务拆分方案

状态：已确认；阶段 A、B 已实施，阶段 C 已完成存储包迁移并保留兼容门面

## 1. 目标与术语

本方案把用户所说的“GI/CI”统一理解为本地交互客户端，包括当前 CLI 和产品方向中计划增加的 TUI；后端指为远端同步提供多租户 HTTP API 的服务。

拆分目标不是简单增加两个文件夹，而是让以下三类职责拥有明确边界：

- `client`：CLI/TUI 入口、本地配置、本地工作空间、交互输出和远端 API 客户端。
- `server`：HTTP 路由、鉴权、用户管理、租户隔离和服务端启动。
- `core`：客户端与服务端共同使用的领域模型、存储接口和资源业务规则。

必须保持现有 `ttl` 命令、数据文件和 `/api/v1` HTTP 契约兼容；目录拆分本身不改变用户行为。

## 2. 当前混乱的根因

当前目录不只是命名问题，而是依赖边界交叉：

| 现状 | 问题 |
| --- | --- |
| `main.go` 同时定义根 CLI、服务端用户管理和同步命令 | 客户端与服务端运维入口混在一个 456 行文件中 |
| `command/` 只包含部分 Cobra 命令 | 看到该目录无法确认所有命令都在哪里 |
| `api/` 是服务端 HTTP 层，但 `db/` 同时包含本地、远端客户端和租户存储 | “后端代码”分散在 `api/`、`db/` 和根目录 |
| `models/` 同时承担持久化模型和跨端协议模型 | 共享契约与存储细节没有区分 |
| `sync/` 是客户端同步流程，但名字像通用后端能力 | 容易误判它属于哪一端 |

因此，仅把 `command/` 改名为 `frontend/`、把 `api/` 改名为 `backend/` 并不能真正分开。

## 3. 目标目录

采用一个 Go module、两个可执行入口、共享内部包的布局：

```text
ttl-cli/
├── cmd/
│   ├── ttl/                         # 本地客户端可执行程序
│   │   └── main.go
│   └── ttl-server/                  # 远端服务可执行程序
│       └── main.go
├── internal/
│   ├── client/
│   │   ├── cli/                     # Cobra 根命令和客户端子命令
│   │   ├── tui/                     # 后续 TUI；未实现前不创建空包
│   │   ├── remote/                  # /api/v1 HTTP 客户端
│   │   └── sync/                    # diff、push、pull 及交互确认
│   ├── server/
│   │   ├── api/                     # 路由、handler、DTO、响应封装
│   │   ├── auth/                    # API Key 认证中间件
│   │   ├── tenant/                  # 用户与租户存储生命周期
│   │   └── app/                     # server 组装和启动
│   ├── core/
│   │   ├── resource/                # Resource、Tag、Audit、History 领域类型
│   │   └── storage/                 # Storage 接口和共享存储契约
│   ├── storage/
│   │   ├── sqlite/                  # 客户端默认本地实现
│   │   └── bbolt/                   # bbolt 与服务端租户数据库实现
│   ├── config/                      # 客户端配置和工作空间
│   ├── cryptox/                     # 数据加密与密钥生命周期
│   └── i18n/                        # 本地化加载与语言资源
├── integration_test/
│   ├── client/                      # CLI、本地数据和工作空间场景
│   └── server/                      # API、鉴权、租户和同步场景
├── docs/
├── scripts/
├── install.sh
├── install.ps1
└── go.mod
```

这里保留一个 module，而不是立即拆为两个仓库或多个 module。客户端和服务端仍共享数据结构与存储语义，单 module 能先获得清晰目录，同时避免版本联动和发布流程复杂化。

## 4. 依赖规则

拆分后用依赖方向判断代码应放在哪里：

```text
cmd/ttl
  -> internal/client/*
  -> internal/core/*
  -> internal/storage/*

cmd/ttl-server
  -> internal/server/*
  -> internal/core/*
  -> internal/storage/bbolt

禁止：
internal/client/* -> internal/server/*
internal/server/* -> internal/client/*
internal/core/*   -> internal/client/* 或 internal/server/*
```

具体规则：

1. 客户端访问远端只能通过 `internal/client/remote` 的 HTTP 契约，不能直接调用 server handler；迁移期仅 `internal/client/cli` 可依赖 `internal/server/cli` 以保留同进程的 `ttl server` 兼容入口。
2. 服务端不依赖 Cobra、终端输出、客户端配置或工作空间。
3. `core` 不依赖具体数据库、HTTP、Cobra 或全局变量。
4. SQLite 与 bbolt 实现依赖 `core/storage`，反向依赖禁止。
5. API DTO 只属于 server API；若客户端需要同一 JSON 契约，提取为小型共享 protocol 包，不直接共享 handler 内部类型。

## 5. 现有文件迁移映射

| 当前文件/目录 | 目标位置 | 说明 |
| --- | --- | --- |
| `main.go` | 拆到 `cmd/ttl/main.go`、`internal/client/cli/`、`internal/server/app/` | 根入口只保留依赖组装和退出码 |
| `migrate.go` | `internal/client/cli/migrate.go` | 这是客户端数据运维命令 |
| `command/*.go` | `internal/client/cli/` | CLI/TUI 共享逻辑要先从 Cobra handler 中提取 |
| `sync/` | `internal/client/sync/` | 属于本地客户端对远端的同步用例 |
| `api/` | `internal/server/api/` | handler 与 HTTP DTO 保持在服务端 |
| `db/tenant_storage.go` | `internal/server/tenant/storage.go` | 只服务于多租户后端 |
| `db/user_store.go` | `internal/server/tenant/users.go` | 后端用户和 API Key 管理 |
| `db/context.go` | `internal/server/api/storage_context.go` | 请求级存储上下文 |
| `db/sqlite.go` | `internal/storage/sqlite/` | 默认本地后端 |
| `db/db.go` 中 `LocalStorage` | `internal/storage/bbolt/` | bbolt 具体实现 |
| `db/db.go` 中 `CloudStorage` | `internal/client/remote/` | 它是 HTTP 客户端，不是数据库 |
| `db/db.go` 中 `SyncStorage` | 评估后移入 `internal/client/sync/` 或删除 | 与显式 push/pull 语义重叠，先确认生产消费者 |
| `db/storage.go` | 拆到 `internal/core/storage/` 与客户端组装层 | 去掉全局 `db.Stor` 后再完成 |
| `models/` | `internal/core/resource/`，用户模型归 `internal/server/tenant/` | 避免服务端用户模型泄漏给客户端 |
| `conf/` | `internal/config/` | 当前主要是客户端本地配置和工作空间 |
| `crypto/` | `internal/cryptox/` | 避免与标准库语义混淆 |
| `i18n/` | `internal/i18n/` | 主要服务本地交互；服务端错误后续单独规范 |
| `util/` | 按真实消费者下沉 | 不保留泛化的杂物包 |

## 6. 可执行程序和兼容性

最终产出两个二进制：

```bash
go build -o ttl ./cmd/ttl
go build -o ttl-server ./cmd/ttl-server
```

建议的用户入口：

```bash
ttl add ...
ttl get ...
ttl sync ...
ttl ui

ttl-server serve --port 8080 --data-dir /var/lib/ttl
ttl-server user add --id alice --name Alice
ttl-server user list
```

迁移期保留原来的 `ttl server ...` 作为兼容代理，并在后续版本明确弃用周期。第一阶段只新增 `ttl-server`，不立即删除旧入口。`/api/v1` 路径、请求/响应 JSON 和 API Key 头保持不变。

## 7. 分阶段实施

### 阶段 A：先拆可执行入口，不搬核心包

- [x] 把 server 与 user 命令从 `main.go` 拆到 `internal/server/cli`。
- [x] 新增 `cmd/ttl-server`，提供 `serve` 和 `user` 命令。
- [x] 保留根 `ttl` 构建和 `ttl server` 兼容入口，两者复用同一命令实现。
- [x] 建立 `cmd/ttl` 和可测试的客户端 `NewRootCommand`，根目录只保留兼容包装。
- 验收：两个二进制可构建；原 CLI 黑盒回归通过；server 单元与集成测试通过。

### 阶段 B：拆 HTTP 客户端与服务端存储

- [x] 将 `CloudStorage` 从 `db` 移到客户端 remote 包。
- [x] 将 HTTP API、`UserStore` 和租户存储管理移到 `internal/server`。
- [x] 将兼容的 `SyncStorage` 移到客户端同步边界。
- [x] 用接口构造 server handler，不再回退到全局 `db.Stor`。
- 验收：客户端包不 import server 包；server 包不 import Cobra/客户端包；API 契约测试通过。

### 阶段 C：建立共享 core 并迁移本地存储

- [x] 抽出最小 `Storage` 接口与资源领域类型；`models` 与 `db.Storage` 暂以类型别名保持源码和数据兼容。
- [x] SQLite、bbolt 实现迁入 `internal/storage`；`db` 仅保留类型别名、构造器和全局门面以兼容现有命令。
- [x] server handler 通过构造函数注入存储，不再读取全局 `db.Stor`。
- 消除客户端命令对 `db.Stor` 的依赖；server 已通过请求上下文注入，不再依赖该全局变量。
- 验收：`core` 无具体存储和传输依赖；包级、集成和 race 检查通过。

### 阶段 D：接入 TUI 并清理兼容层

- 在 `internal/client/tui` 实现 `ttl ui`，复用 core 用例而非解析 CLI 文本。
- 根据已发布版本的兼容承诺，决定何时删除根目录旧入口与 `ttl server` 代理。
- 更新安装包，明确 `ttl` 与 `ttl-server` 是否分别发布。

## 8. 测试与验收

每个阶段至少运行：

```bash
gofmt -s -l .
go test ./...
go test ./integration_test/...
go vet ./...
./scripts/regression.sh
```

拆分完成后增加结构检查：

```bash
go list -deps ./cmd/ttl
go list -deps ./cmd/ttl-server
```

人工验收：

- `ttl` 不包含 server 用户管理的实现依赖。
- `ttl-server` 不读取客户端工作空间，也不依赖终端交互输出。
- CLI 和未来 TUI 操作同一套 core 用例与本地数据。
- 客户端与服务端只通过 `/api/v1` 契约通信。
- 现有 SQLite/bbolt 数据无需迁移即可继续读取。
- 多租户认证、用户禁用、API Key 重置和租户隔离行为不变。

## 9. 风险与控制

| 风险 | 控制措施 |
| --- | --- |
| 一次移动大量文件导致 diff 无法审核 | 四阶段迁移，每阶段单独任务和验证 |
| `internal/` 使外部 Go 使用者无法 import | 当前产品是可执行程序；若以后承诺 SDK，再建立稳定 `pkg/` |
| 两个二进制增加发布成本 | 第一阶段保留 `ttl server`，先验证独立 server 构建再调整发布脚本 |
| 去除全局存储改变初始化顺序 | 先增加构造函数和测试，再迁移调用者 |
| 模型拆分破坏 JSON/数据库兼容 | 保持字段名和序列化标签，增加兼容数据测试 |
| 目录看似清楚但业务仍在 handler 中重复 | CLI/TUI 共享行为必须进入 core 用例，入口层只处理输入输出 |

## 10. 明确不采用的方案

- 不使用 `frontend/` / `backend/` 两个大包：Go 中容易形成新的杂物目录，且共享模型、存储接口无处安放。
- 不立即拆成两个仓库：当前共享代码多，跨仓版本和发布维护成本高。
- 不立即拆成多个 Go module：会引入 replace、版本联动和测试矩阵，不能直接解决职责混杂。
- 不把 TUI 当作独立后端消费者：本地 TUI 应直接复用 core；只有远端同步经过 HTTP。

## 11. 完成定义

当仓库达到以下状态时，才算真正完成前后端拆分：

1. 从 `cmd/ttl` 和 `cmd/ttl-server` 能分别看见两个产品入口。
2. 从 `internal/client`、`internal/server` 和 `internal/core` 能判断代码归属。
3. 依赖规则被测试或静态检查验证，而不是只写在文档里。
4. 现有 CLI、HTTP、数据库和同步兼容性测试全部通过。
5. README、PROJECT_OVERVIEW、安装和发布说明与新目录一致。
