# ttl-cli 项目脉络梳理

> 面向第一次接触项目的快速导览。先看清入口、包边界和数据流，再按具体需求深入文件。

## 1. 项目定位

`ttl-cli` 是一个 Go 编写的个人知识归档 CLI。用户以 key-value 形式保存资源，通过标签、模糊搜索、工作日志和历史记录管理本地数据；同一套存储抽象还支持工作空间、加密、远端 API 和双向同步。

可以把项目理解为三层：Cobra 命令和 HTTP API 是入口，`internal/client/app.Service` 与 `internal/core/storage.Storage` 是客户端业务访问边界，SQLite、bbolt、云端和组合式同步存储是具体实现。需要区分这里的 `MirroredStorage` 与 `sync/`：前者把写操作同时转发到本地和云端，后者负责计算两端差异并执行显式 push/pull。

从工程边界看，仓库包含两个独立应用：`ttl` 是运行在用户设备的客户端，CLI 和未来的 TUI 都属于该应用；`ttl-server` 是独立发布和部署的云端服务。两者暂时共享一个 Git 仓库和 Go module，但必须拥有独立构建制品、配置、发布和运行环境。仓库数量不等于应用工程数量。

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
go build -o ttl ./cmd/ttl
go run ./cmd/ttl <command>
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
├── cmd/ttl/                # 正式的本地客户端可执行入口
├── cmd/ttl-server/         # 独立云端服务可执行入口
├── internal/client/        # 客户端命令树、远端访问和同步适配
├── internal/server/        # 后端 API、租户数据与运维命令
├── internal/core/          # 客户端与服务端共享的领域类型和存储契约
├── internal/storage/       # SQLite 与 bbolt 具体存储实现
├── command/                # 普通 CLI 子命令
├── db/                     # 待迁移调用方和删除的旧存储门面
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
├── WORK_ITEMS.md           # 本地当前任务正文和交付证据
├── WORK_ITEMS.json         # 本地当前任务状态元数据
├── WORK_ITEMS_ARCHIVE.md   # 已完成任务历史
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
       -> 按配置选择存储后端并创建 internal/client/app.Service
       -> 展开历史快捷参数并记录命令历史
  -> rootCmd.Execute()
       -> command/*、server、sync 或 migrate
  -> 关闭 client/app.Service
```

客户端命令树集中在 `internal/client/cli/`，`cmd/ttl` 是唯一目标入口。普通资源命令仍由 `command/` 提供，客户端层负责组装同步、迁移、存储初始化和生命周期。后端 HTTP API、租户数据和运维命令集中在 `internal/server/`，只由 `cmd/ttl-server` 组装。

独立构建与运行：

```bash
go build -o ttl ./cmd/ttl
go build -o ttl-server ./cmd/ttl-server
ttl-server serve --port 8080
ttl-server user list
```

独立构建只是应用拆分的第一步。正式交付时，客户端发布物只包含 `ttl`；云端发布物只包含 `ttl-server` 及明确声明的运行文件，不依赖客户端二进制、客户端配置或源码目录。当前仓库尚未完成 `ttl-server` 的独立发布和部署流程。

## 5. 目录职责

### `command/`：CLI 交互层

- `commands.go`：资源增删改查、标签、配置、版本、审计和历史。
- `workspace.go`：工作空间创建、切换、查看和删除。
- `import.go` / `export.go`：数据导入导出。
- `encrypt.go`：数据加解密以及密钥导入、导出和校验。
- `log.go`：工作日志。
- `init.go`：Shell completion 初始化。
- `tags.go`：标签统计和标签资源列表。

命令层负责参数、输出和用户可见错误；持久化操作应通过当前命令 context 中的 `client/app.Service` 完成。支持本地化的用户文案应放到 `i18n/locales/`。

### `db/`：持久化边界

- `db.go`：历史兼容类型别名与构造器；客户端生产代码不再依赖该包，后续测试迁移完成后删除。
- `sqlite.go`：待调用方迁移后删除的 SQLite 构造器。
- `storage.go`：全局存储门面、初始化、迁移以及审计/历史/日志代理。

具体 SQLite 和 bbolt 实现已经迁入 `internal/storage/`，远端 HTTP 存储位于 `internal/client/remote/`，租户存储位于 `internal/server/tenant/`。`db/` 仅因调用方迁移尚未完成而存在，迁移完成后直接删除，不能再增加新的服务端或客户端实现。

### `internal/server/`：后端服务

- `api/`：组装 HTTP 路由，完成 API Key 校验、租户存储注入，并处理资源、标签、审计和历史接口。
- `tenant/`：维护 `users.json`、API Key 以及每个用户的隔离数据库。
- `cli/`：构造 `ttl-server serve/user` 命令。

### 其他核心包

- `sync/` 只负责比较两份资源和执行同步方向，不负责读取 CLI 参数。
- `conf/` 负责 `~/.ttl/ttl.ini`、自定义配置文件和工作空间路径。
- `crypto/` 负责加密格式与密钥生命周期；密钥不应进入仓库。
- `models/` 是待迁入 `internal/core/resource` 的旧跨包类型目录。
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

客户端本地存储、远端连接、后端租户数据和三种远端交互模式的开发阶段说明见 `docs/client-server-data-flow.md`。该文档同时记录当前同步边界与已知缺口；调整 API、连接参数或同步语义时应一并更新。

## 7. 配置、数据与安全边界

| 内容 | 默认位置或入口 | 注意事项 |
| --- | --- | --- |
| 配置 | `~/.ttl/ttl.ini` | 测试使用 `--conf` 和临时目录 |
| 默认数据库 | `~/.ttl/data.db` | 工作空间可覆盖路径和存储类型 |
| 加密密钥 | `~/.ttl/.key` | 不提交、不写入文档内容 |
| 服务端用户 | `<data-dir>/users.json` | 包含 API Key，不能作为样例提交 |
| 租户数据 | `<data-dir>/tenants/` | 删除用户与删除租户数据是不同动作 |

`Storage` 接口、资源 JSON 字段、INI 结构和 HTTP DTO 是跨组件契约。开发阶段允许直接调整，但必须在同一变更中更新全部调用方、当前数据、测试和文档，不保留旧格式读取或过渡代理。

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
| 新增普通 CLI 命令 | `command/`、`internal/client/cli/root.go` 的命令注册 | `i18n/locales/`、命令测试、黑盒回归 |
| 修改资源增删改查 | `command/commands.go` | `db/`、`internal/server/api/handlers.go`、集成测试 |
| 修改工作空间 | `command/workspace.go`、`conf/ini.go` | `conf/workspace_test.go`、CLI 回归 |
| 修改存储后端 | `internal/client/app/` 与 `internal/storage/` | `Storage` 接口、迁移、API、同步 |
| 修改 HTTP API | `internal/server/api/server.go`、`internal/server/api/handlers.go` | 中间件、DTO、handler 测试、集成测试 |
| 修改同步策略 | `sync/sync.go`、`internal/client/cli/root.go` 的 sync command | 单元测试、服务端同步集成测试 |
| 修改加密 | `crypto/`、`command/encrypt.go` | 当前数据重建或迁移、加密集成测试 |
| 修改配置格式 | `conf/ini.go`、`models.TtlIni` | 当前配置更新、工作空间、决策记录 |
| 修改用户文案 | `i18n/locales/` 和对应命令 | 各语言 key 一致性、README 示例 |
| 修改客户端发布安装 | `install.sh`、`install.ps1` | 根目录稳定 URL、跨平台行为 |
| 修改云端发布部署 | `cmd/ttl-server/`、发布流水线和部署说明 | 服务端制品内容、配置、密钥、数据卷、健康检查与回滚 |

## 10. 当前结构的已知整理点

客户端（CLI/TUI）与远端服务的目标边界和分阶段迁移方案见 `docs/client-server-separation-plan.md`。当前已经建立两个应用入口并完成主要服务端和存储包迁移，但旧门面、客户端应用服务和独立交付仍在迁移中。

阶段 A、B 已建立 `cmd/ttl` 与 `cmd/ttl-server` 两个入口，并把远端客户端、服务端 API 和租户存储迁入对应边界。W-006 已将客户端命令改为显式服务注入并删除根目录入口和 `ttl server` 代理。后续整理按以下顺序单独立项：

1. 为 `ttl-server` 建立独立制品、部署配置、升级和回滚流程；能单独编译不能替代独立交付验收。
2. 在客户端工程中实现 `ttl ui`，CLI 与 TUI 复用 application/core 契约。

每项都应独立验证并单独审阅，不能一次性搬完所有目录。

## 11. 建议阅读顺序

1. `README.md`：用户能力和命令用法。
2. `internal/client/cli/root.go`：客户端启动生命周期与命令入口。
3. `models/models.go`：核心数据结构。
4. `internal/client/app/` 和 `internal/storage/`：客户端服务、存储接口与具体实现。
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

一句话总结：单仓库中包含本地客户端和云端服务两个独立应用工程，入口层把 CLI/TUI 或 HTTP 请求转换为共享 core 契约上的操作，两个应用分别构建和交付。
