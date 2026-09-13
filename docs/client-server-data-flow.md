# 客户端、后端与同步实现说明

> 本文记录当前开发阶段的实际实现，帮助后续开发定位代码和继续调整方案。它不是稳定协议；实现变化时应同步更新。
>
> 设计评语和整改复核见 [TUI、CLI 与客户端边界设计评审与整改复核](reviews/2026-09-12-tui-client-boundary-design.md)。

## 1. 范围与术语

当前项目没有浏览器前端。“客户端”指 `ttl`，包括已经实现的 CLI 和计划增加的 TUI；两者是同一个客户端的不同交互形式。客户端目标入口是 `cmd/ttl/`；“后端”指可独立构建和部署的 `ttl-server`，入口是 `cmd/ttl-server/`。根目录入口和客户端中的 `ttl server` 已删除，不属于目标接口。

CLI 和 TUI 使用相同的账户、API Key、业务服务、存储契约和远端 API。相对于服务端，两者没有身份或数据所有权上的区别；服务端也不需要识别请求来自哪一种界面。客户端、服务端和具体数据库实现共享 `internal/core/storage.Storage` 契约。客户端访问远端时不会直接调用服务端 Go 包，而是通过 HTTP `/api/v1` 接口交互。

当前产品是单用户资源管理工具。每个服务端账户拥有独立租户数据，资源不能分享给其他账户，也没有协作者、资源级权限或多人实时编辑能力。同一账户在多个进程或设备上使用属于技术上可能发生的访问方式，不代表产品提供多人协作。

## 2. 当前整体结构

```text
用户
  |
  v
ttl client（当前 CLI；TUI 计划中）
  |
  +-- sqlite（默认）------------> 本地 SQLite
  |
  +-- local / bbolt -----------> 本地 bbolt
  |
  +-- cloud -- HTTP + API Key -> ttl-server -> 用户独立 bbolt
  |
  +-- sync ---+----------------> 本地 bbolt
              +-- HTTP + API Key -> ttl-server -> 用户独立 bbolt
```

`internal/client/app.OpenStorage` 根据 `--storage` 和配置选择实际存储。CLI 的资源命令通过当前命令 context 中的 `app.Service` 访问存储，因此同一条 `add/get/update/del` 命令可以落到本地数据库、远端 API 或组合存储。

## 3. 客户端如何存储

### 3.1 存储模式

| 模式 | 读取位置 | 写入位置 | 是否自动访问后端 |
| --- | --- | --- | --- |
| `sqlite` | 本地 SQLite | 本地 SQLite | 否 |
| `local` / `bbolt` | 本地 bbolt | 本地 bbolt | 否 |
| `cloud` | 后端 HTTP API | 后端 HTTP API | 是，每次操作都访问后端 |
| `sync` | 本地 bbolt | 先本地 bbolt，后远端 HTTP API | 是，每次修改都访问后端 |

默认模式是 `sqlite`。配置默认位于 `~/.ttl/ttl.ini`，默认 SQLite 数据库位于 `~/.ttl/data.db`。工作空间可以在 INI 中指定独立的数据库路径和存储类型；未指定时沿用全局设置。

`sync` 模式当前固定组合本地 bbolt 和远端 HTTP 存储，不会沿用默认 SQLite 作为本地副本。它读取本地，资源增、删、改时先写本地，成功后再写远端。

### 3.2 客户端写入链路

以 `ttl add` 为例：

```text
Cobra 命令
  -> internal/client/cli/commands.AddCmd
  -> internal/client/app.Service.SaveResource
  -> 当前 core/storage.Storage
       -> SQLiteStorage
       -> bbolt.LocalStorage
       -> remote.Storage
       -> client/sync.MirroredStorage
```

客户端由 `internal/client/cli` 创建一个显式 `internal/client/app.Service`，通过命令 context 注入；命令执行结束后由客户端入口关闭服务。命令和同步流程不读取全局存储状态。

## 4. 后端如何实现

### 4.1 启动与路由

`ttl-server serve` 通过 `internal/server/cli` 调用 `internal/server/app.RunWithSignals`，
由应用层组装 `internal/server/api.NewServerHandler`：

1. 从 `<data-dir>/users.json` 加载用户。
2. 创建以 `<data-dir>/tenants` 为根目录的租户存储管理器。
3. 注册 `/api/v1` HTTP 路由。
4. 用多租户 API Key 中间件包裹路由。
5. 使用显式监听地址创建 Go 标准库 `http.Server`，并在 SIGTERM/SIGINT 或 context
   取消时执行优雅关闭和租户存储清理。

当前主要接口及实现进度如下。“服务端已实现”只表示 HTTP handler 和存储调用存在，不代表远端客户端已经完整接入。

| 请求 | 行为 | 当前进度 |
| --- | --- | --- |
| `GET /api/v1/resources` | 获取当前用户的全部资源和标签 | 服务端、远端客户端均已接入 |
| `POST /api/v1/resources` | 新增资源 | value 已实现；请求 DTO 不接收 tags |
| `PUT /api/v1/resources/{key}` | 更新资源值 | value 已实现；请求 DTO 不接收 tags |
| `DELETE /api/v1/resources/{key}` | 删除资源 | 服务端、远端客户端均已接入 |
| `POST /api/v1/resources/{key}/tags` | 添加标签 | 服务端已实现；`remote.Storage` 未接入 |
| `DELETE /api/v1/resources/{key}/tags/{tag}` | 删除标签 | 服务端已实现；`remote.Storage` 未接入 |
| `POST /api/v1/resources/{key}/rename` | 重命名资源 | 服务端已实现；远端客户端未调用该专用路由，命令层仍以读/删/写组合实现 |
| `GET /api/v1/audit/stats` | 获取审计统计 | 服务端、远端客户端均已接入 |
| `GET /api/v1/history` | 获取历史记录 | 服务端、远端客户端均已接入 |

响应使用统一 JSON 包装：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

客户端主要依据响应体中的 `code` 判断业务是否成功。`code = 0` 表示成功，非零值会转换为客户端错误。

### 4.2 认证与租户隔离

客户端在请求头中发送：

```http
Authorization: Bearer <API_KEY>
```

多租户中间件按以下顺序处理请求：

```text
读取 Bearer API Key
  -> UserStore 查找用户
  -> 检查用户 Active 状态
  -> StorageManager 按 user.ID 获取存储
  -> 将用户 ID 和 Storage 注入请求 context
  -> API handler 执行业务操作
```

服务端默认数据布局为：

```text
<data-dir>/
├── users.json
└── tenants/
    ├── alice/
    │   └── data.db
    └── bob/
        └── data.db
```

`users.json` 包含用户信息、API Key 和启用状态，保存权限为 `0600`。每个用户的数据存放在独立的 bbolt 文件中；服务端首次访问用户存储时打开数据库，并在进程内缓存连接。

## 5. 客户端如何连接后端

### 5.1 启动服务

构建两个可执行文件：

```bash
go build -o ttl ./cmd/ttl
go build -o ttl-server ./cmd/ttl-server
```

创建用户并保存命令输出的 API Key：

```bash
./ttl-server --data-dir /var/lib/ttl user add --id alice --name Alice
```

启动服务：

```bash
./ttl-server --data-dir /var/lib/ttl --port 8080 serve
```

开发和正式云端部署都必须使用独立的 `ttl-server` 制品。当前仓库只能验证该入口可独立构建，发布、部署、升级和回滚链路尚未完成，不能据此宣称云端独立交付已经完成。目标边界见 [`client-server-separation-plan.md`](client-server-separation-plan.md)。

### 5.2 客户端参数

远端连接目前由命令行参数传入：

| 参数 | 含义 | 默认值 |
| --- | --- | --- |
| `--cloud-url` | 服务端根地址，例如 `http://127.0.0.1:8080` | 空 |
| `--cloud-key` | 用户 API Key | 空 |
| `--cloud-timeout` | HTTP 超时秒数 | `30` |

例如，检查同步差异但不修改数据：

```bash
./ttl \
  --cloud-url http://127.0.0.1:8080 \
  --cloud-key '<API_KEY>' \
  sync --direction push --dry-run
```

执行从本地到远端的同步：

```bash
./ttl \
  --cloud-url http://127.0.0.1:8080 \
  --cloud-key '<API_KEY>' \
  sync --direction push
```

当前客户端代码不会从 README 示例中的 `[server] endpoint/api_key` 配置读取连接信息。配置落盘方式尚未收敛前，应以上述参数为准。API Key 不应写入仓库、示例配置或长期文档中的真实值。

`remote.Storage.Init` 只创建带超时的 HTTP Client，不会主动进行健康检查；地址、网络和认证问题通常在第一次 API 请求时返回。

## 6. 什么时候同步

### 6.1 默认本地模式

`sqlite`、`local` 和 `bbolt` 模式不会在启动、退出或后台定时同步。普通资源命令只改变本地数据库。只有用户显式执行 `ttl sync` 才会读取远端并比较两端数据。

显式同步的步骤是：

```text
读取本地全部资源
  -> GET /api/v1/resources 读取远端全部资源
  -> 按资源 key 比较 value 和 tags
  -> 得到 local_only / remote_only / conflict
  -> 显示差异
  -> 按用户指定方向修改目标端
```

方向语义如下：

| 方向 | 权威端 | 结果 |
| --- | --- | --- |
| `pull` | 远端 | 本地最终与远端一致 |
| `push` | 本地 | 远端最终与本地一致 |
| `auto` | 用户现场选择 | 显示差异后选择 pull、push 或 skip |

`--dry-run` 只显示差异，不执行增、删、改。

显式同步不是合并算法。`pull` 会删除本地独有资源，`push` 会删除远端独有资源；冲突由权威端直接覆盖目标端。

### 6.2 `sync` 镜像模式

`--storage sync` 的修改操作会立即尝试写远端：

```text
写本地 bbolt
  -> 本地成功
  -> 调用远端 API
```

本地失败时不会请求远端。远端失败时本地修改已经提交，目前没有事务回滚、离线队列、自动重试或补偿操作，因此可能出现两端不一致。出现这种情况时只能检查差异并再次执行显式同步。

### 6.3 `cloud` 模式

`--storage cloud` 不维护本地业务副本。资源命令直接通过 HTTP 操作后端，因此这里不存在“稍后同步”的阶段；每次请求成功即表示后端已处理该操作。

## 7. 当前同步边界与已知限制

以下内容是当前代码行为，后续开发需要重点核对：

- 显式同步每次读取两端的全部资源，没有分页、变更游标或增量同步。
- 差异只处理 `ORIGIN` 资源，审计、历史和日志不属于显式同步集合。
- 远端存储的部分审计、历史和日志写入方法仍是空实现，不能视为完整数据同步。
- 远端标签读取已实现；服务端也提供独立的标签添加和删除路由，但 `remote.Storage` 尚未接入这些路由。
- 客户端创建资源时会发送 `tags`，但服务端创建请求 DTO 当前没有接收标签；远端更新的客户端和服务端 DTO 都只处理 `value`。因此创建、更新和同步含标签的资源尚未端到端完成，标签可能丢失或无法按权威端覆盖。
- `sync` 镜像模式先写本地再写远端，远端失败不会回滚本地。
- HTTP API 没有协议版本协商；开发阶段变更由客户端和服务端同步发布。
- 服务端使用明文 HTTP 监听，没有内建 TLS、健康检查端点或优雅关闭流程。远程部署时应在受控网络或 HTTPS 反向代理后使用，并由部署层配置防火墙和访问控制。
- 服务端监听 `:port`，即所有网络接口，而不只监听本机回环地址。
- W-010 已完成独立交付：服务端默认监听 loopback，提供 `/healthz`，支持优雅关闭和版本化回滚；生产环境仍需由部署层配置外部 TLS 终止。
- API Key 以明文形式保存在服务端 `users.json` 中。虽然文件权限为 `0600`，仍需要保护数据目录和备份。
- 资源 key 直接拼接到 URL path，当前客户端没有执行路径转义；包含 `/`、空格或特殊字符的 key 需要额外验证。

这些限制不应在文档中被描述成稳定设计。修复或改变其中任一行为时，应补充相应测试，并更新本节。

## 8. 多客户端竞态场景

本节是对同一账户被两个进程或设备同时使用时的防御性分析，不是分享或多人协作设计。当前产品没有跨账户资源分享、协作者、资源级权限或多人实时编辑；建议的正常使用方式仍是一个账户只有一个主动写端。CLI 与未来 TUI 若作为两个进程同时写入，也属于这里的同账户竞态。

当前实现允许多个客户端使用同一 API Key 访问同一租户，但没有资源版本号、ETag、`If-Match`、幂等键、服务端同步锁或冲突响应。因此数据库可以保证单次写操作不会把文件写坏，却不能保证跨请求的业务操作符合用户预期。当前整体语义接近“服务端最后完成的写入生效”。

### 8.1 场景与当前结果

| 场景 | 可能交错 | 当前结果 | 风险 |
| --- | --- | --- | --- |
| 两个客户端同时更新同一 key | A、B 都读取 v1；A 写 v2；B 写 v3 | 最后完成的写入覆盖前一次写入 | A 的修改静默丢失 |
| 一个客户端更新，另一个删除 | A 读取资源；B 删除；A 随后更新 | bbolt 的 `UpdateResource` 会直接 `Put`，SQLite 的更新路径也可重新插入 | 被删除资源可能复活；执行顺序相反时更新可能丢失 |
| 两个客户端同时加标签 | A、B 都读取旧标签；各自添加不同标签并保存整份资源 | 后写入者覆盖先写入者的标签集合 | 某一方新增标签丢失 |
| 一个客户端删标签，另一个更新正文 | 两端都基于旧快照保存整份资源 | 后写请求携带的旧标签可能恢复已删除标签，或正文更新被覆盖 | 字段间相互覆盖 |
| 两个客户端同时创建同一 key | 两端都通过“是否存在”检查，然后分别保存 | 检查和保存不是一个原子操作，两个请求都可能报告成功 | 后写值覆盖先写值，客户端却都以为创建成功 |
| 重命名与任意并发修改 | rename 先删旧 key，再建新 key；另一客户端同时更新或创建旧/新 key | 删除和创建不是同一事务，也没有检查目标 key 是否已存在 | 旧 key 复活、目标 key 被覆盖，或 rename 中途只完成删除 |
| 客户端修改期间执行显式 push/pull | 同步先获取两端快照并计算 diff，之后逐项执行；期间其他客户端继续写 | 同步按旧快照覆盖或删除后到达的修改 | 新修改被回滚、误删，且同步结束后也未必一致 |
| 两个客户端同时执行同步 | 两边分别基于不同快照执行 push/pull | 操作逐项交错，没有全局事务或租户级同步锁 | 结果取决于时序，可能部分成功、相互覆盖 |
| 镜像写入超时后用户重试 | 远端可能已提交，但客户端没有收到响应 | POST 重试可能得到“已存在”；更新和删除则可能被重复执行 | 客户端无法判断第一次操作是否成功，接口缺少幂等语义 |
| 用户管理与在线服务并发 | 独立 CLI 和服务进程分别加载自己的 `UserStore` | 服务进程不会热加载文件变化；多个 CLI 实例也没有跨进程文件锁 | 新增、禁用、重置 Key 不会立即作用于运行中的服务；并发管理可能覆盖 `users.json` |
| 删除用户与该用户的在途请求并发 | 管理命令使用新的 `StorageManager` 删除目录，服务进程持有自己的缓存数据库连接 | 运行中的服务仍可能继续使用旧用户快照和已打开连接 | 禁用或删除不能可靠阻止在途及后续请求，租户数据状态不确定 |

### 8.2 已有保护能解决什么

- bbolt 事务和 SQLite 事务能力可以串行化单个数据库写操作，主要防止存储文件损坏。
- `StorageManager` 的互斥锁只保护同一服务进程内“为用户创建或取得存储实例”的缓存操作，不保护资源业务读改写。
- `UserStore` 的互斥锁只保护同一个 `UserStore` 实例；服务进程和每次管理命令创建的是不同实例，因此没有跨实例、跨进程保护。
- 本地 bbolt 会阻止多个进程同时以当前方式打开同一数据库文件；这会返回锁超时错误，而不是协调两个本地客户端共同修改。

这些机制不能检测客户端使用了旧数据，也不能避免丢失更新、删除后复活或过期同步覆盖。

### 8.3 开发阶段建议的处理优先级

在明确支持多个客户端前，当前使用约束应是：同一租户尽量只有一个主动写客户端；同步前停止其他客户端写入；先使用 `--dry-run` 检查差异；不要让两个客户端同时执行 push/pull；用户和 API Key 管理后重启服务以重新加载 `users.json`。

只有未来明确把同账户多设备作为受支持场景时，才按以下顺序单独立项和评审；这些建议不是当前 TUI 的隐含范围：

1. 为资源引入服务端版本号或单调 revision，并在 GET 响应中返回。
2. 更新、删除和重命名要求提交预期 revision；版本不匹配时返回 `409 Conflict`，不能静默覆盖。
3. 将服务端的 create-if-absent、rename、标签变更做成单个存储事务，避免 handler 层“先读后写”。
4. 同步执行前重新校验目标 revision；只有全量比较不能满足已确认需求时，再评估基于变更日志和游标的同步协议。旧快照不得直接删除或覆盖新版本。
5. 为可重试的写请求增加幂等键和操作结果查询，区分“未执行”和“已执行但响应丢失”。
6. 明确用户管理的运行时一致性：统一到服务端管理 API，或者加入跨进程文件锁和配置热加载；删除用户前先阻断认证并关闭服务进程持有的租户存储。
7. 增加并发集成测试，至少覆盖同 key 双更新、更新与删除、标签并发、rename、双同步和超时重试。

在这些约束落地前，多客户端冲突属于已知的数据一致性风险。当前开发阶段以限制同账户并发写入来控制风险；如果产品以后承诺同账户多设备，再采用 revision 和显式冲突提示等简单机制，不能直接推导出 CRDT 或多人协作。

## 9. 代码入口索引

| 关注点 | 代码位置 |
| --- | --- |
| 客户端命令树与存储初始化 | `internal/client/cli/root.go` |
| 普通资源命令 | `internal/client/cli/commands/commands.go` |
| 客户端存储选择与生命周期 | `internal/client/app/` |
| 共享存储接口 | `internal/core/storage/storage.go` |
| SQLite 实现 | `internal/storage/sqlite/storage.go` |
| bbolt 实现 | `internal/storage/bbolt/storage.go` |
| HTTP 远端存储 | `internal/client/remote/storage.go` |
| 镜像读写存储 | `internal/client/sync/storage.go` |
| 差异计算和 push/pull | `internal/client/sync/sync.go` |
| 服务端命令 | `internal/server/cli/command.go` |
| HTTP 路由和启动 | `internal/server/api/server.go` |
| 认证中间件 | `internal/server/api/middleware.go` |
| HTTP handler 和 DTO | `internal/server/api/handlers.go`、`types.go` |
| 用户和 API Key | `internal/server/tenant/user_store.go` |
| 租户数据库 | `internal/server/tenant/storage.go` |
| 跨端集成场景 | `integration_test/server_sync_test.go` |

## 10. 后续维护规则

- 调整存储模式、默认值或配置来源时，更新第 3、5 节。
- 修改 API 路径、请求 DTO、认证方式或响应结构时，更新第 4 节。
- 修改同步触发时机、冲突策略、删除语义或失败补偿时，更新第 6、7 节。
- 修改资源并发控制、版本校验、幂等策略或用户管理一致性时，更新第 8 节。
- 修复已知限制后删除对应条目，并链接覆盖该行为的测试。
- 长期架构取舍记录在 `docs/decisions/`；本文只描述当前可观察行为。
- 不在本文记录真实地址、用户信息、API Key、加密密钥或生产数据路径。
