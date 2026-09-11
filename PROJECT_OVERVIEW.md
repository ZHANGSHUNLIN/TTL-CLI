# ttl-cli 项目脉络梳理

> 面向第一次接触项目的快速导览。先看清入口、包边界和数据流，再按具体需求深入文件。

## 1. 项目定位

`ttl-cli` 是一个 Go 编写的个人知识归档 CLI。用户以 key-value 形式保存资源，通过标签、模糊搜索、工作日志和历史记录管理本地数据；同一套存储抽象还支持工作空间、加密、远端 API 和双向同步。

可以把项目理解为三层：Cobra 命令和 HTTP API 是入口，`db.Storage` 是业务访问边界，SQLite、bbolt、云端和组合式同步存储是具体实现。需要区分这里的 `SyncStorage` 与 `sync/`：前者把写操作同时转发到本地和云端，后者负责计算两端差异并执行显式 push/pull。

## 2. 技术栈与常用命令

| 领域 | 实现 |
| --- | --- |
| 语言与模块 | Go 1.25，模块名 `ttl-cli` |
| CLI | Cobra |
| 默认本地存储 | SQLite（`modernc.org/sqlite`） |
| 可选本地存储 | bbolt |
| 配置 | INI（`gopkg.in/ini.v1`） |
| 本地化 | go-i18n，语言文件位于 `i18n/locales/` |
| HTTP 服务 | Go 标准库 `net/http` |

```bash
go build -o ttl .
go run . <command>
go test ./...
go test ./integration_test/...
go vet ./...
./scripts/regression.sh
./scripts/verify.sh
```

手动验证 CLI 时应使用临时 `HOME` 和配置文件，避免修改真实的 `~/.ttl`。

## 3. 根目录结构

```text
.
├── main.go                 # 兼容根目录构建的 ttl 客户端入口
├── cmd/ttl/                # 正式的本地客户端可执行入口
├── cmd/ttl-server/         # 独立后端服务可执行入口
├── internal/client/cli/    # 客户端命令树、同步和迁移命令
├── internal/server/        # 后端 API、租户数据与运维命令
├── command/                # 普通 CLI 子命令
├── db/                     # 存储接口、后端实现和存储门面
├── sync/                   # 同步差异计算与 push/pull
├── conf/                   # 配置文件与工作空间
├── crypto/                 # 数据加解密与密钥管理
├── i18n/                   # 本地化加载器和语言资源
├── models/                 # 跨包共享的数据结构
├── util/                   # 无状态通用辅助函数
├── integration_test/       # 跨包、API 和同步场景测试
├── scripts/                # 黑盒回归与完整验证
├── docs/                   # 开发流程、测试说明和决策记录
├── .agents/skills/         # 本项目的个人工程 Skill
├── WORK_ITEMS.md           # 本地任务状态板
└── README*.md              # 对外说明及多语言版本
```

多语言 README 需要保留在根目录，方便 GitHub 页面相互跳转。`install.sh` 和 `install.ps1` 也应继续保留稳定的根目录下载地址。它们看起来分散，但属于发布入口，不适合仅为视觉整齐而移动。

## 4. 应用启动链路

```text
main()
  -> 初始化 i18n
  -> 刷新 Cobra 命令描述
  -> PersistentPreRunE
       -> 注入 debug / confFile 上下文
       -> 按配置选择存储后端并调用 db.InitDB
       -> 展开历史快捷参数并记录命令历史
  -> rootCmd.Execute()
       -> command/*、server、sync 或 migrate
  -> db.CloseDB()
```

客户端命令树集中在 `internal/client/cli/`，`cmd/ttl` 是正式入口；根目录 `main.go` 只是保留 `go build .` 的兼容包装。普通资源命令仍由 `command/` 提供，客户端层负责组装同步、迁移、存储初始化和生命周期。后端 HTTP API、租户数据和运维命令集中在 `internal/server/`，既供兼容入口 `ttl server` 使用，也供独立入口 `cmd/ttl-server` 使用。

独立构建与运行：

```bash
go build -o ttl .
go build -o ttl ./cmd/ttl
go build -o ttl-server ./cmd/ttl-server
ttl-server serve --port 8080
ttl-server user list
```

## 5. 目录职责

### `command/`：CLI 交互层

- `commands.go`：资源增删改查、标签、配置、版本、审计和历史。
- `workspace.go`：工作空间创建、切换、查看和删除。
- `import.go` / `export.go`：数据导入导出。
- `encrypt.go`：数据加解密以及密钥导入、导出和校验。
- `log.go`：工作日志。
- `init.go`：Shell completion 初始化。
- `tags.go`：标签统计和标签资源列表。

命令层负责参数、输出和用户可见错误；持久化操作应通过 `db` 层完成。支持本地化的用户文案应放到 `i18n/locales/`。

### `db/`：持久化边界

- `db.go`：`Storage` 接口，以及 bbolt、云端和组合式同步存储实现。
- `sqlite.go`：默认 SQLite 实现。
- `storage.go`：全局存储门面、初始化、迁移以及审计/历史/日志代理。
- `tenant_storage.go`：服务端按用户隔离存储。
- `user_store.go`：服务端用户及 API Key 文件。
- `context.go`：在请求上下文中传递存储实例。

这里的文件名有一处容易误解：`db.go` 不只是初始化，而包含 bbolt 实现；真正的全局初始化入口 `InitDB` 位于 `storage.go`。修改存储时不要只凭文件名判断职责。

### `internal/server/`：后端服务

- `api/`：组装 HTTP 路由，完成 API Key 校验、租户存储注入，并处理资源、标签、审计和历史接口。
- `tenant/`：维护 `users.json`、API Key 以及每个用户的隔离数据库。
- `cli/`：构造 `ttl-server serve/user` 和兼容的 `ttl server` 命令。

### 其他核心包

- `sync/` 只负责比较两份资源和执行同步方向，不负责读取 CLI 参数。
- `conf/` 负责 `~/.ttl/ttl.ini`、自定义配置文件和工作空间路径。
- `crypto/` 负责加密格式与密钥生命周期；密钥不应进入仓库。
- `models/` 是持久化结构和跨包契约，字段变化需要检查兼容性与迁移。
- `util/` 只放无状态、跨入口复用的小工具。

## 6. 核心业务流程

### 本地 CLI 读写

```text
Cobra command
  -> db 全局门面
  -> Storage 接口
  -> SQLite 或 bbolt
  -> workspace 对应的数据文件
```

默认配置目录为 `~/.ttl`。若命令传入 `--conf`，配置和工作空间路径以该文件为起点；测试和脚本利用这一点隔离真实用户数据。

### 多租户 HTTP API

```text
HTTP request
  -> MultiTenantAuthMiddleware
  -> UserStore 校验 API Key
  -> TenantStorageManager 选择用户存储
  -> request context
  -> handler
  -> Storage
```

服务端数据目录中，`users.json` 保存用户信息，各租户数据库位于 `tenants/`。涉及认证或租户隔离的改动必须覆盖 `internal/server/api/`、`internal/server/tenant/` 和集成测试。

### 同步

```text
本地 Storage + CloudStorage
  -> sync.ComputeDiff
  -> local_only / remote_only / conflict
  -> ExecutePull 或 ExecutePush
  -> 目标 Storage
```

同步以资源 key 为比较单位，冲突处理会覆盖目标侧。修改比较规则或方向语义时，应同步检查 `sync/sync_test.go` 和 `integration_test/server_sync_test.go`。

## 7. 配置、数据与安全边界

| 内容 | 默认位置或入口 | 注意事项 |
| --- | --- | --- |
| 配置 | `~/.ttl/ttl.ini` | 测试使用 `--conf` 和临时目录 |
| 默认数据库 | `~/.ttl/data.db` | 工作空间可覆盖路径和存储类型 |
| 加密密钥 | `~/.ttl/.key` | 不提交、不写入文档内容 |
| 服务端用户 | `<data-dir>/users.json` | 包含 API Key，不能作为样例提交 |
| 租户数据 | `<data-dir>/tenants/` | 删除用户与删除租户数据是不同动作 |

`Storage` 接口、`models` 中的 JSON 字段、INI 结构和 HTTP DTO 都是兼容性边界。调整这些内容前应评估数据迁移、旧客户端和同步行为。

## 8. 测试分层

| 层级 | 位置或命令 | 适用场景 |
| --- | --- | --- |
| 包级单元测试 | 各包旁的 `*_test.go` | 局部规则和边界条件 |
| 跨包集成测试 | `integration_test/` | API、同步、加密和资源生命周期 |
| CLI 黑盒回归 | `scripts/regression.sh` | 真实二进制参数、输出和持久化 |
| 完整验证 | `scripts/verify.sh` | 宽范围或提交前验证 |

涉及文件、HTTP、数据库、全局状态或进程环境的测试应使用 `t.TempDir`、`httptest` 或临时 `HOME`，并确保资源被关闭。

## 9. 常见修改入口速查

| 需求 | 优先查看 | 同时检查 |
| --- | --- | --- |
| 新增普通 CLI 命令 | `command/`、`main.go` 的命令注册 | `i18n/locales/`、命令测试、黑盒回归 |
| 修改资源增删改查 | `command/commands.go` | `db/`、`internal/server/api/handlers.go`、集成测试 |
| 修改工作空间 | `command/workspace.go`、`conf/ini.go` | `conf/workspace_test.go`、CLI 回归 |
| 修改存储后端 | `db/storage.go` 与对应实现 | `Storage` 接口、迁移、API、同步 |
| 修改 HTTP API | `internal/server/api/server.go`、`internal/server/api/handlers.go` | 中间件、DTO、handler 测试、集成测试 |
| 修改同步策略 | `sync/sync.go`、`main.go` 的 `syncCmd` | 单元测试、服务端同步集成测试 |
| 修改加密 | `crypto/`、`command/encrypt.go` | 数据兼容、迁移、加密集成测试 |
| 修改配置格式 | `conf/ini.go`、`models.TtlIni` | 旧配置兼容、工作空间、决策记录 |
| 修改用户文案 | `i18n/locales/` 和对应命令 | 各语言 key 一致性、README 示例 |
| 修改发布安装 | `install.sh`、`install.ps1` | 根目录稳定 URL、跨平台行为 |

## 10. 当前结构的已知整理点

客户端（CLI/TUI）与远端服务的目标边界和分阶段迁移方案见 `docs/client-server-separation-plan.md`。当前代码仍处于迁移前结构，不要仅凭目标目录寻找实现。

阶段 A 已建立 `cmd/ttl` 与 `cmd/ttl-server` 两个入口，客户端命令树位于 `internal/client/cli/`，服务端运维命令位于 `internal/server/cli/`。现有 `go build .` 和 `ttl server` 作为兼容入口继续使用同一实现。后续整理按以下顺序单独立项：

1. 把 `CloudStorage` HTTP 客户端从 `db/` 迁入 client 边界。
2. 抽出共享 core 和具体存储实现，并逐步消除全局 `db.Stor`。
3. 发布流程全部改用 `cmd/ttl` 后，再评估移除根目录兼容入口。

每项都应独立验证并单独审阅，不能一次性搬完所有目录。

## 11. 建议阅读顺序

1. `README.md`：用户能力和命令用法。
2. `internal/client/cli/root.go`：客户端启动生命周期与命令入口。
3. `models/models.go`：核心数据结构。
4. `db/storage.go` 和 `db/db.go`：存储门面、接口与 bbolt 实现。
5. `command/commands.go`：主要 CLI 行为。
6. `internal/server/api/`、`internal/server/tenant/`：服务端链路。
7. `sync/sync.go`：同步语义。
8. `integration_test/` 和 `scripts/regression.sh`：系统实际承诺的行为。

## 12. 文档维护规则

- 新增顶层包、重要命令或存储后端时，更新目录结构和修改入口。
- 修改启动、认证、租户或同步链路时，更新对应流程。
- 配置路径、默认值或验证命令变化时，同时更新 README 和本导览。
- 不记录账号、密码、API Key、加密密钥或真实用户数据。
- 设计理由写入 `docs/decisions/`，本文件只描述当前有效结构。

一句话总结：入口层把 CLI 与 HTTP 请求转换为存储操作，`db.Storage` 隔离具体后端，配置、加密、同步和多租户能力围绕这条主线展开。
