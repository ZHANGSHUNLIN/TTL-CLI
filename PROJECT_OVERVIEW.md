# ttl-cli 项目概览

## 1. 项目定位

`ttl-cli` 是 Go 编写的本地优先个人知识归档客户端。用户以 key-value 形式保存资源，通过标签、搜索、工作日志、历史记录、工作空间和导入导出管理数据。

本仓库只维护 `ttl` 客户端。远端后端由独立服务工程维护和部署；这里保留 cloud 配置、HTTP 适配器和同步客户端，但不包含服务端入口、handler、认证、租户存储、用户管理或服务端发布流程。

## 2. 技术栈与常用命令

| 领域 | 实现 |
| --- | --- |
| 语言与模块 | Go 1.25，模块名 `ttl-cli` |
| CLI / TUI | Cobra、Bubble Tea |
| 默认本地存储 | SQLite（`modernc.org/sqlite`） |
| 配置 | INI（`gopkg.in/ini.v1`） |
| 本地化 | go-i18n，语言文件位于 `internal/i18n/locales/` |
| 远端访问 | Go 标准库 `net/http` |

```bash
go build -o ttl ./cmd/ttl
go run ./cmd/ttl <command>
go test ./...
go test ./integration_test/...
go vet ./...
./scripts/regression.sh
./scripts/verify.sh
```

手动验证应使用临时 `HOME` 和配置文件，避免修改真实 `~/.ttl` 数据。

## 3. 目录结构

```text
.
├── cmd/ttl/                       # 唯一产品可执行入口
├── internal/client/cli/           # Cobra 命令树与执行生命周期
├── internal/client/cli/commands/  # 命令适配
├── internal/client/app/           # 客户端用例与存储生命周期
├── internal/client/tui/           # TUI 适配
├── internal/client/remote/        # 外部后端 HTTP 适配
├── internal/client/sync/          # 同步 diff、push/pull 与镜像存储
├── internal/core/                 # 客户端内部模型与存储契约
├── internal/storage/              # SQLite 与旧格式兼容检查
├── internal/config/               # 配置与工作空间
├── internal/crypto/               # 加密与密钥生命周期
├── internal/i18n/                 # 本地化加载器与语言资源
├── integration_test/              # 跨包客户端场景
├── scripts/                       # CLI 回归与完整验证
└── docs/                          # 需求、设计、决策和研发流程
```

## 4. 启动与数据流

```text
cmd/ttl
  -> internal/client/cli
  -> internal/client/app.Service
  -> local: internal/storage/sqlite
  -> cloud: internal/client/remote -> 外部后端 HTTP API
```

客户端命令执行前按 `--storage` 和配置创建 storage，注入 `app.Service`，执行完成后关闭资源。本地模式只访问 SQLite；cloud 模式只通过 HTTP 访问已配置的外部服务。

同步命令会同时打开本地 SQLite 与远端 HTTP 适配器，使用 `internal/client/sync` 计算差异并执行 pull 或 push。独立后端应兼容客户端当前 `/api/v1`、Bearer 认证和 JSON DTO；本仓库以 `httptest` mock 验证客户端契约，不提供服务端实现。

## 5. 包职责

- `internal/client/cli/commands`：解析参数、输出用户文案、调用当前 `app.Service`。
- `internal/client/app`：资源、审计、历史和日志用例，以及 storage 创建与关闭。
- `internal/client/remote`：请求外部后端并转换 `/api/v1` DTO。
- `internal/client/sync`：比较两份资源并按方向更新目标 storage。
- `internal/storage/sqlite`：当前本地持久化实现。
- `internal/storage/bbolt`：仅保留旧格式识别/兼容性测试所需实现，不是当前 storage mode。
- `internal/core/resource`：资源、标签、审计、历史和日志类型。
- `internal/core/storage`：客户端内部统一 storage 契约。
- `internal/config`、`internal/crypto`、`internal/i18n`：配置、加密和本地化基础能力。

## 6. 配置和安全边界

| 内容 | 默认位置或入口 | 注意事项 |
| --- | --- | --- |
| 配置 | `~/.ttl/ttl.ini` | 测试使用 `--conf` 和临时目录 |
| 默认数据库 | `~/.ttl/data.sqlite` | 工作空间可覆盖路径 |
| 加密密钥 | `~/.ttl/.key` | 不提交、不写入文档内容 |
| 远端凭据 | 环境变量或命令参数 | 配置只保存环境变量名，不保存明文凭据 |

外部后端的账号、API Key 颁发、租户隔离、数据目录、部署和备份不属于本仓库。客户端改动远程协议时，必须同步更新 mock 契约测试，并与独立后端协调兼容性。

## 7. 测试分层

| 层级 | 位置或命令 | 覆盖内容 |
| --- | --- | --- |
| 包级单元测试 | 各包旁 `*_test.go` | 局部行为、远端 mock 契约 |
| 跨包集成测试 | `integration_test/` | 本地存储、加密、导入导出、生命周期 |
| CLI 黑盒 | `scripts/regression.sh` | 构建后客户端的参数、输出和持久化 |
| 完整验证 | `scripts/verify.sh` | 格式、构建、架构、测试、race 和 vet |

涉及文件、HTTP、数据库或进程环境的测试应使用 `t.TempDir`、`httptest` 或临时 `HOME`。真实外部后端兼容性应由后续跨工程契约或端到端流水线验证。

## 8. 常见修改入口

| 需求 | 优先查看 | 同时检查 |
| --- | --- | --- |
| 新增 CLI 命令 | `internal/client/cli/commands/`、`internal/client/cli/root.go` | i18n、命令测试、CLI 回归 |
| 修改资源行为 | `internal/client/app/`、`internal/client/cli/commands/` | local/cloud 适配器测试 |
| 修改工作空间 | `internal/config/`、workspace command | 配置测试、CLI 回归 |
| 修改本地存储 | `internal/storage/sqlite/` | storage 契约、兼容性测试 |
| 修改远端协议 | `internal/client/remote/` | mock HTTP 契约、独立后端协调 |
| 修改同步策略 | `internal/client/sync/`、sync command | 单元测试、远端适配器测试 |
| 修改发布安装 | `install.sh`、`install.ps1`、workflows | 客户端跨平台制品 |

## 9. 维护规则

- 不在本仓库新增后端入口或实现；相关需求进入独立后端工程。
- cloud 客户端只依赖 HTTP 契约，不导入或复制服务端内部代码。
- 修改启动、配置、存储、远端协议或同步链路时更新本文和相关决策。
- 不记录账号、密码、API Key、加密密钥或真实用户数据。
- 设计理由写入 `docs/decisions/`，本文件只描述当前有效结构。
