# 技术方案：本地与云端存储模式收敛

任务：W-019
类型：feature
状态：方案评审通过，编码完成

## 方案摘要

客户端把活动存储源收敛为 `local` 和 `cloud` 两个互斥值：`local` 打开本地 SQLite，
`cloud` 通过现有 HTTP 存储契约访问服务端。服务端继续按租户隔离，但将租户资源文件从
bbolt 切换为 SQLite。新版本不读取或转换旧 bbolt/旧 SQLite 文件；模式切换只改变访问目标，
不触发复制或同步。同步、版本和 CRDT 由 W-020 另行设计。

## 当前执行链路

1. `cmd/ttl` 进入 `internal/client/cli/root.go`，解析 `--storage`、cloud 连接参数和
   `--conf`。
2. `newPreRun` 根据命令行、workspace 和全局配置选择 storage 类型，并调用
   `internal/client/app.OpenStorage`。
3. `OpenStorage` 只支持 local、cloud；旧值直接拒绝并提示新模式，不提供兼容入口。
4. 服务端 `internal/server/api/middleware.go` 从 `StorageManager` 获取当前用户存储，
   handler 通过 `core/storage.Storage` 访问资源、历史、审计和日志。

## 目标模块职责

| 模块 | 责任 |
| --- | --- |
| `internal/client/cli` | 解析存储参数、展示兼容错误、传递单次覆盖，不保存秘密 |
| `internal/config` | 保存全局/workspace 存储类型、数据库路径和远程配置引用；解析优先级 |
| `internal/client/app` | 根据已解析模式创建并初始化 SQLite 或 HTTP storage |
| `internal/storage/sqlite` | 实现 `core/storage.Storage`、SQLite schema 初始化和前向版本检查 |
| `internal/client/remote` | 实现 cloud API 请求、超时、认证和错误映射 |
| `internal/server/tenant` | 按用户 ID 管理 SQLite 句柄、打开/关闭、备份和删除前置检查 |
| `internal/server/api` | 保持资源 API 契约，不感知具体数据库实现 |
| `internal/core/storage` | 保持两种模式共享的资源、标签、历史、审计和日志接口 |

## 存储模式与配置契约

### 模式解析

解析顺序固定为：

```text
命令行 --storage
  > 当前 workspace 的 storage_type
  > 全局 [storage] type
  > 默认 local
```

`--storage` 只对当前命令生效。合法值只有 `local`、`cloud`；`sqlite`、`bbolt`、`sync`
等旧值直接拒绝并提示新模式，不创建对应存储。远程 profile 可以配置多个，但当前应用或
workspace 只激活一个；`cloud` 所需地址、账户和凭据引用优先从当前 workspace 读取，再回退
到全局配置；命令行 `--cloud-url`、`--cloud-key` 仅作为本次调用覆盖，API Key 不写入日志和仓库。

### 客户端打开边界

- `local`：调用 `sqlite.NewSQLiteStorage`，解析 workspace 数据路径，创建目录并初始化
  schema；新默认文件名为 `data.sqlite`，显式 `db_path` 仍可指定完整路径；数据库打开
  失败直接返回本地错误。
- `cloud`：调用 `remote.NewStorage`，先执行健康/认证校验，再使用 HTTP API；网络或认证
  失败直接返回远程错误，不读取本地数据库。
- 普通命令只能注入一个 `core/storage.Storage` 实例；W-019 不在此处加入镜像写入、
  outbox 或同步线程。

## 服务端 SQLite 组织

目标数据目录为：

```text
/var/lib/ttl-server/
├── users.json
└── tenants/
    ├── <user-id>/data.sqlite
    └── <user-id>/data.sqlite
```

`StorageManager` 将句柄表改为 `map[string]*sqlite.SQLiteStorage`，按用户首次请求惰性
打开 `<dataDir>/<userID>/data.sqlite`，进程退出时统一关闭。用户 ID 必须经过现有安全校验
后才能参与路径拼接；禁止通过请求参数指定数据库路径。每租户独立文件能延续当前隔离模型，
也使备份、恢复和删除可以按租户执行。

SQLite 初始化要求：启用 WAL（Windows 使用兼容模式）、`busy_timeout`、受控连接池和
`synchronous=NORMAL`；schema 通过 `schema_version` 或等价迁移表记录版本，初始化和每次
升级按顺序执行幂等迁移。现有 resources、audit、history、logs 表字段和
`core/storage` 语义保持兼容；如需字段变化，必须新增迁移版本和回滚说明。

## 加密、生命周期与并发

- 继续使用现有资源加密能力；新数据库初始化不得把加密值静默转换为明文。发现旧格式或
  无法识别加密状态时直接拒绝并保留源文件。
- `Init` 负责目录、连接、pragma 和新格式 schema 初始化；`Close` 由客户端命令或服务退出统一调用。
- 同一租户的句柄在 `StorageManager` 内串行创建，重复请求复用已打开句柄；关闭和删除前先
  从句柄表移除并等待关闭完成。
- SQLite 写入使用事务，读写错误携带租户和操作上下文；API 层只返回不泄露路径、密钥或
  资源值的错误信息。

## 旧格式处理与租户删除

新版本不提供迁移命令或兼容窗口。打开路径时检测文件签名和 SQLite schema：发现旧 bbolt、
旧 SQLite 或未知格式即返回不兼容错误，不读取、不转换、不覆盖，也不自动删除。用户需要
自行备份旧文件，并为新版本指定新的 `data.sqlite` 路径。

租户删除必须先停止接受该租户请求、关闭句柄并按运维确认生成备份或隔离副本，再执行删除；
不能无条件调用 `os.RemoveAll`。恢复只能将备份还原到独立的新格式目录，校验后再切回配置。

## 备选方案

- **共享 SQLite 文件**：需要所有查询携带 `tenant_id`，任一遗漏都会造成跨租户风险；不采用。
- **继续使用 bbolt**：无法统一客户端和服务端格式，也不能复用现有 SQLite 实现；不采用。
- **PostgreSQL**：对当前单机、少量并发部署引入运维依赖和迁移成本过高；暂不采用，规模
  变化时另立决策记录。

## 测试与回滚计划

- 单元测试：模式解析优先级、旧值错误、SQLite schema/pragma、句柄复用和关闭、租户路径
  隔离、旧格式检测和不兼容错误。
- 集成测试：两个租户通过 API 读写互不影响；服务重启后数据可读；cloud 不可用时不触碰
  local；服务重启后资源、标签、历史、审计和日志保持一致。
- CLI 回归：`--storage local/cloud`、workspace 覆盖、单次覆盖、多个远程 profile 选择、
  缺少 cloud 配置和旧值拒绝；不修改用户真实 home 数据。
- 回滚：实现验证失败时恢复旧 release 和原数据目录；新 SQLite 文件保留诊断，旧文件不被
  自动删除或转换。

## 阶段边界

W-019 不定义同步触发、版本号、操作日志、增量协议、CRDT 或冲突解决。W-020 只能依赖本
任务交付的 `local/cloud` 打开边界和稳定资源语义，具体同步协议另行评审。

## 评审与实现检查

- [x] 所有需求验收项都有模块负责人和可观察验证方式
- [x] SQLite schema 前向版本和加密不兼容策略得到评审确认
- [x] 服务端租户删除、备份和恢复流程得到实现与测试覆盖
- [x] 方案评审通过后进入 WBS 和编码，生产代码已按方案实现
