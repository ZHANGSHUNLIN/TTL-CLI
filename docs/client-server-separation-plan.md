# TTL 客户端与后端服务拆分方案

状态：已确认；W-017 一次性结构切换已完成；W-010 云端独立交付尚未完成

## 1. 目标与术语

本方案把本地 CLI 和产品方向中计划增加的 TUI 归入同一个客户端应用；后端是为远端同步提供多租户 HTTP API 的云端服务。

本文所说的“工程”是可独立构建、测试、发布、部署和回滚的应用单元，不等同于 Git 仓库或 Go module。项目必须包含两个独立应用工程：

- `ttl` 客户端工程：运行在用户设备，包含 CLI、未来的 TUI、本地配置、本地数据和远端 API 客户端。
- `ttl-server` 云服务工程：运行在服务端，包含 HTTP API、认证、用户管理、租户隔离和服务端数据。

两个应用工程当前可以共用一个源码仓库和一个 Go module，但不得共用运行制品、部署配置或运行时权限。TUI 是客户端的另一种交互入口，不是第三个独立应用工程。

拆分目标不是简单增加两个文件夹，而是让以下三类职责拥有明确边界：

- `client`：CLI/TUI 入口、本地配置、本地工作空间、交互输出和远端 API 客户端。
- `server`：HTTP 路由、鉴权、用户管理、租户隔离和服务端启动。
- `core`：客户端与服务端共同使用的领域模型、存储接口和资源业务规则。

项目当前处于开发阶段，直接删除旧命令入口、旧包门面和旧数据格式支持。拆分以目标边界为准，调用方、测试和文档随实现同步调整。

## 2. 原结构混乱的根因

拆分前的目录不只是命名问题，而是依赖边界交叉。下表用于解释迁移来源，不代表所有问题仍存在：

| 现状 | 问题 |
| --- | --- |
| `main.go` 同时定义根 CLI、服务端用户管理和同步命令 | 客户端与服务端运维入口混在一个 456 行文件中 |
| 客户端命令、存储和同步曾分散在顶层包 | 维护者无法确认修改入口 |
| 服务端 HTTP 层与租户存储曾与客户端存储混杂 | “后端代码”职责边界不清 |
| 共享模型与存储契约曾通过别名转发 | 共享契约与实现细节没有区分 |
| 同步 diff 与镜像写入曾位于不同目录 | 容易误判同步能力归属 |

因此，仅改目录名称并不能真正分开；W-017 已按依赖方向完成一次性归属收敛。

## 3. 目标目录

采用一个源码仓库和 Go module、两个独立应用工程、两个可执行入口以及受控共享内部包的布局：

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
│   │   │   └── commands/            # Cobra 命令适配
│   │   ├── tui/                     # TUI 适配
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
│   ├── config/                      # 共享配置和工作空间
│   ├── crypto/                      # 共享数据加密与密钥生命周期
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

这里保留一个 module，而不是立即拆为两个仓库或多个 module。客户端和服务端仍共享数据结构与存储语义，单 module 可以避免过早引入共享模块版本联动，但不降低应用和部署边界：云端最终只运行 `ttl-server` 制品，客户端最终只分发 `ttl` 制品。

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

1. 客户端访问远端只能通过 `internal/client/remote` 的 HTTP 契约，不能直接调用 server handler；任何客户端包都不得依赖 `internal/server`。
2. `internal/server/api` 和 `internal/server/tenant` 不依赖客户端 Cobra 命令、终端交互或客户端专属状态；共享 `internal/config`、`internal/crypto` 基础包可被 bbolt 适配器使用；`internal/server/cli` 可以使用 Cobra 提供服务端自身的 `serve` 和用户管理命令。
3. `core` 不依赖具体数据库、HTTP、Cobra 或全局变量。
4. SQLite 与 bbolt 实现依赖 `core/storage`，反向依赖禁止。
5. API DTO 只属于 server API；若客户端需要同一 JSON 契约，提取为小型共享 protocol 包，不直接共享 handler 内部类型。

## 5. 现有文件迁移映射

| 当前文件/目录 | 目标位置 | 说明 |
| --- | --- | --- |
| `main.go` | 删除；入口拆到 `cmd/ttl/main.go`、`cmd/ttl-server/main.go` | 不保留根目录构建入口 |
| `migrate.go` | `internal/client/cli/migrate.go` | 这是客户端数据运维命令 |
| `command/*.go` | `internal/client/cli/commands/` | 已完成一次性迁移为客户端 Cobra 适配层 |
| `sync/` | `internal/client/sync/` | 已完成一次性迁移，属于本地客户端对远端的同步用例 |
| `api/` | `internal/server/api/` | handler 与 HTTP DTO 保持在服务端 |
| `db/tenant_storage.go` | `internal/server/tenant/storage.go` | 只服务于多租户后端 |
| `db/user_store.go` | `internal/server/tenant/users.go` | 后端用户和 API Key 管理 |
| `db/context.go` | `internal/server/api/storage_context.go` | 请求级存储上下文 |
| `db/sqlite.go` | `internal/storage/sqlite/` | 默认本地后端 |
| `db/db.go` 中 `LocalStorage` | `internal/storage/bbolt/` | bbolt 具体实现 |
| `db/db.go` 中 `CloudStorage` | `internal/client/remote/` | 它是 HTTP 客户端，不是数据库 |
| `db/db.go` 中 `SyncStorage` | `internal/client/sync/` | 已完成一次性迁移，与显式 push/pull 统一归属 |
| `db/storage.go` | `internal/core/storage/` 与 `internal/client/app/` | 已去掉全局 `db.Stor` 并删除旧门面 |
| `models/` | `internal/core/resource/`，用户模型归 `internal/server/tenant/` | 避免服务端用户模型泄漏给客户端 |
| `conf/` | `internal/config/` | 当前主要是客户端本地配置和工作空间 |
| `crypto/` | `internal/crypto/` | 共享基础包，避免服务端通过 bbolt 反向依赖客户端 |
| `i18n/` | `internal/i18n/` | 主要服务本地交互；服务端错误后续单独规范 |
| `util/` | 按真实消费者下沉 | 不保留泛化的杂物包 |

## 6. 可执行程序

两个应用工程必须分别产出二进制：

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

服务端能力只通过 `ttl-server` 暴露，`ttl` 不提供 `server` 或服务端用户管理命令。客户端与服务端需要同时调整协议时，直接更新 `/api/v1` 实现、调用方和测试。

### 6.1 构建与部署边界

- 客户端发布物只包含 `ttl`，面向用户设备安装；CLI 和 TUI 共享该客户端制品。
- 云端发布物只包含 `ttl-server` 及其运行所需文件，不包含 `ttl`、客户端配置、工作空间数据或客户端安装脚本。
- 客户端和服务端使用独立的 CI 构建任务、制品名称、校验和、发布步骤与回滚目标；两者可以由同一个 Git tag 触发，但任一方必须能够单独重建和部署。
- 构建环境可以检出整个 monorepo；运行环境只能获得对应制品。不得把源码仓库整体复制到云端作为运行目录。
- 服务端部署必须明确监听地址、数据卷、密钥注入、日志、健康检查、优雅关闭、TLS 终止位置和最小运行权限。

## 7. 一次性结构切换与后续交付

W-017 已在一个变更窗口内完成入口、模型、配置、加密、同步、存储和测试的统一迁移；内部依赖顺序只用于编辑和验证，不代表多个中间版本或分阶段交付。当前目录和 import 图以本文件第 3、4 节为准。

### 后续：建立云端独立交付

- 为 `ttl-server` 建立独立的 CI 构建和发布制品，不能只在本地验证可编译。
- 提供只包含服务端运行内容的部署包或最小容器镜像。
- 定义服务端配置、密钥、数据卷、健康检查、优雅关闭、日志和 TLS 边界。
- 增加部署冒烟、制品内容检查以及独立升级和回滚验证。

W-010 当前的交付基线、技术方案和参考部署见：
[`docs/requirements/2026-09-13-server-independent-delivery.md`](requirements/2026-09-13-server-independent-delivery.md)、
[`docs/tech-designs/2026-09-13-server-independent-delivery.md`](tech-designs/2026-09-13-server-independent-delivery.md)
和 [`docs/server-deployment.md`](server-deployment.md)。

### 已完成：TUI 与旧入口收敛

- `internal/client/tui` 已实现 `ttl ui`，复用 client app/core 用例而非解析 CLI 文本。
- 已删除根目录旧入口与 `ttl server` 代理。
- 客户端安装包不依赖 `ttl-server`。

## 8. 测试与验收

一次性切换完成后统一运行：

```bash
gofmt -s -l .
go test ./...
go test ./integration_test/...
go test -race ./...
go vet ./...
./scripts/regression.sh
```

拆分完成后增加结构检查：

```bash
go test ./internal/architecture
```

独立交付阶段还必须验证：发布系统同时生成可分别下载的 `ttl` 和 `ttl-server` 制品；服务端部署只安装 `ttl-server` 制品；服务端冒烟测试不依赖客户端配置、客户端二进制或源码目录。

结构测试内部使用 `go list -json` 检查直接与传递依赖。`ttl` 不得依赖任何服务端包；`ttl-server` 不得依赖客户端或旧 `command`、`db`、`sync` 包。

人工验收：

- 云端运行目录或容器只包含 `ttl-server` 及明确声明的运行文件。
- `ttl-server` 可以脱离客户端制品、客户端配置和源码目录独立启动、升级和回滚。
- `ttl` 不包含 server 用户管理的实现依赖。
- `ttl-server serve` 不读取客户端工作空间，也不要求交互式终端输入。
- CLI 和未来 TUI 操作同一套 core 用例与本地数据。
- 客户端与服务端只通过 `/api/v1` 契约通信。
- 多租户认证、用户禁用、API Key 重置和租户隔离行为不变。

## 9. 风险与控制

| 风险 | 控制措施 |
| --- | --- |
| 一次移动大量文件导致 diff 无法审核 | 先完成设计和依赖盘点，再在同一变更窗口按内部顺序编辑，最后统一评审和回归 |
| `internal/` 使外部 Go 使用者无法 import | 当前产品是可执行程序；若以后承诺 SDK，再建立稳定 `pkg/` |
| 两个应用增加发布成本 | 使用同一仓库和可复用流水线模板，但为客户端和服务端生成独立制品与部署步骤 |
| 去除全局存储改变初始化顺序 | 先增加构造函数和测试，再迁移调用者 |
| 模型或协议调整造成两端不一致 | 同一变更同步更新客户端、服务端和契约测试 |
| 目录看似清楚但业务仍在 handler 中重复 | CLI/TUI 共享行为必须进入 core 用例，入口层只处理输入输出 |

## 10. 明确不采用的方案

- 不使用 `frontend/` / `backend/` 两个大包：Go 中容易形成新的杂物目录，且共享模型、存储接口无处安放。
- 不立即拆成两个仓库：两个应用工程的独立性由构建、依赖、制品和部署边界保证；当前共享代码多，跨仓版本维护成本高。
- 不立即拆成多个 Go module：会引入 replace、版本联动和测试矩阵，不能直接解决职责混杂。
- 不把 TUI 当作独立后端消费者：本地 TUI 应直接复用 core；只有远端同步经过 HTTP。

## 11. 完成定义

当仓库达到以下状态时，才算真正完成前后端拆分：

1. 从 `cmd/ttl` 和 `cmd/ttl-server` 能分别看见客户端与云端服务两个独立应用工程入口。
2. 从 `internal/client`、`internal/server` 和 `internal/core` 能判断代码归属。
3. 依赖规则被测试或静态检查验证，而不是只写在文档里。
4. 当前 CLI、HTTP、数据库和同步行为测试全部通过。
5. 发布系统生成两个独立制品，云端环境不需要 `ttl` 客户端或源码仓库即可部署、升级和回滚 `ttl-server`。
6. README、PROJECT_OVERVIEW、安装、发布和服务端部署说明与实际边界一致。
