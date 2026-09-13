# W-010 技术方案：云端服务独立交付

日期：2026-09-13
任务：W-010
需求：`docs/requirements/2026-09-13-server-independent-delivery.md`
状态：方案评审通过（PASS）

## 方案摘要

保留一个仓库和一个 Go module，但在所有构建和发布路径中明确两个应用入口。
新增小型服务端运行时组合层，负责监听配置、就绪状态、信号处理和资源清理。
服务端发布 Linux `amd64`/`arm64` tarball，参考部署采用带版本号的目录和
`current` 符号链接。公网 TLS 由部署代理终止，Go 服务端只在私有 HTTP 跳上监听。

本方案不引入容器运行时、新存储格式或新 API 协议。

## 当前执行链路

`cmd/ttl-server/main.go` 调用 `internal/server/cli.Run`。CLI 解析 `--listen`、
`--port`、`--data-dir` 和关闭超时，再调用 `internal/server/app.RunWithSignals`。
应用层加载 `tenant.UserStore`，创建 `tenant.StorageManager`，通过
`internal/server/api.NewServerHandler` 组装健康路由、`/api/v1` 路由和认证中间件，
最后运行 `http.Server` 并在收到信号或 context 取消时关闭租户存储。

`.github/workflows/ci.yml` 和 `release.yml` 仍执行 `go build .`，并使用单一的
`ttl-cli` 制品名；仓库根目录没有受支持的非测试 Go 入口。

## 归属与包变更

| 区域 | W-010 中的职责 |
| --- | --- |
| `cmd/ttl-server` | 保持轻量可执行入口，不承载部署策略 |
| `internal/server/cli` | 解析 `--listen`、`--port`、`--data-dir`、可选关闭超时，并调用运行时 |
| `internal/server/app` | 提议的组合/生命周期归属：加载用户存储、创建租户管理器、构造 `http.Server`、管理就绪、信号、关闭和清理 |
| `internal/server/api` | 构造 API 路由和中间件，提供接收就绪状态的健康处理器；保留 `NewHandler` 供测试使用 |
| `internal/server/tenant` | 负责用户存储和每用户存储创建、`CloseAll`；不处理发布和进程信号 |
| `internal/storage/bbolt` | 现有租户数据库实现；W-010 不做格式迁移 |
| `.github/workflows/ci.yml` | 使用 `cmd/` 入口构建并分别上传验证制品 |
| `.github/workflows/release.yml` | 分开构建客户端/服务端矩阵、归档、校验并发布 |
| `scripts/server-delivery-smoke.sh` | 新增干净目录启动、探活、API、关闭和发布布局冒烟脚本 |
| `docs/server-deployment.md` | 新增 systemd/tarball 部署、权限、升级、备份和回滚说明 |

`internal/server/app` 的理由是把当前位于 HTTP 包中的进程生命周期和部署配置
分离出来，使 handler 可独立测试。实现前必须由 W-017 确认它是唯一服务端组合归属；
若 W-017 采用其他包名，只迁移相同职责，不得新增第二个入口。

## 运行时接口

### 服务端配置

定义由服务端运行时拥有的配置（名称可在实现阶段固定）：

```go
type Config struct {
    ListenAddress   string
    Port            int
    DataDir         string
    ShutdownTimeout time.Duration
}
```

CLI 保留 `--port`、`--data-dir`，新增 `--listen`。v1 默认值为 `127.0.0.1`；
私有服务网络需显式设置。配置校验非空数据目录、`1..65535` 端口、非空监听地址
和正关闭超时；无效配置在打开文件或绑定端口前失败。

### 生命周期 API

应用层应提供可测试的函数，等价于：

```go
func Run(ctx context.Context, cfg Config) error
```

`Run` 同步执行启动，创建 `http.Server`，仅在用户存储加载且路由安装完成后标记
就绪，并持续运行到监听器退出或上下文取消。CLI 使用显式配置路径；测试入口
`NewHandler(storage)` 保持兼容，生产启动统一经过 `internal/server/app`。

HTTP server 设置请求头和空闲连接超时。关闭顺序为：

1. 先取消就绪并关闭监听器，拒绝新请求；
2. 用配置的上下文调用 `http.Server.Shutdown`；
3. 即使 HTTP 排空返回错误，也调用 `StorageManager.CloseAll`；
4. 将监听或强制关闭错误返回 CLI，形成非零退出。

SIGINT 和 SIGTERM 由应用层转换为上下文取消；测试可以直接取消 context，不依赖
操作系统信号。

### 健康端点

`internal/server/api` 在认证中间件之前增加 `GET /healthz`：

- 启动就绪后返回 `200 application/json` 和 `{"status":"ok"}`；
- 就绪前返回 `503 application/json` 和 `{"status":"starting"}`；
- 不返回用户、租户、API Key 或资源数据；
- 非 GET 方法返回 `405`。

健康响应只表示进程就绪，不作为 API 版本协商机制。生产路由构造器接收就绪探针
和租户依赖，避免健康端点绕过启动校验；`NewHandler(storage)` 继续作为 API 测试辅助。

## 数据与文件系统契约

不修改数据库 schema 或资源 payload。服务端继续使用：

```text
<data-dir>/users.json
<data-dir>/tenants/<user-id>/data.db
```

启动时在平台允许的范围内以 owner-only 权限创建数据目录和租户父目录；
`users.json` 保持 `0600`，不得打印内容或 API Key。现有数据原地打开，不新增自动
迁移或双写。

发布目录与数据目录分离：

```text
/opt/ttl-server/releases/<version>/ttl-server
/opt/ttl-server/current -> /opt/ttl-server/releases/<version>
/var/lib/ttl-server/users.json
/var/lib/ttl-server/tenants/<user-id>/data.db
```

运行用户拥有 `/var/lib/ttl-server`，只执行发布二进制，不能修改已经激活的旧版本
目录。部署文档说明最小读写权限，密钥不进入发布归档。

## 构建与发布方案

### CI

- 客户端矩阵改为 `go build ./cmd/ttl`。
- 增加 Linux `amd64`/`arm64` 服务端矩阵，使用 `go build ./cmd/ttl-server` 和
  `CGO_ENABLED=0`。
- 单测、集成、race、vet 和架构检查独立于归档打包。
- 验证制品使用不同名称，例如 `ttl-linux-amd64` 和 `ttl-server-linux-amd64`，
  服务端 job 不复用客户端路径。

### Release

tag 为 `vX.Y.Z` 时采用：

```text
ttl-vX.Y.Z-linux-amd64.tar.gz
ttl-vX.Y.Z-darwin-arm64.tar.gz
ttl-server-vX.Y.Z-linux-amd64.tar.gz
ttl-server-vX.Y.Z-linux-arm64.tar.gz
<archive>.sha256
```

客户端保留现有系统矩阵，服务端只构建承诺的 Linux 矩阵。每个归档由显式 staging
目录生成，只包含匹配的可执行文件和可选的版本元数据；不得包含 `.git`、源码、
`users.json`、租户数据库、客户端配置或安装脚本。若健康日志需要版本信息，应使用
专用 version 包或 linker symbol，不得依赖已删除根入口的 `main.version`。

## 部署与恢复方案

参考部署由 systemd（或等价管理器）执行：

```text
/opt/ttl-server/current/ttl-server serve \
  --listen 127.0.0.1 \
  --port 8080 \
  --data-dir /var/lib/ttl-server
```

服务以非 root 用户运行，数据独立持久化，stdout/stderr 交给 supervisor 收集。
反向代理/负载均衡器终止 TLS，只向私有监听地址转发，并用 `/healthz` 探活；防火墙
禁止公网直接访问明文 HTTP 端口。

升级流程：校验归档；停止或静默服务并备份数据；解压到新版本目录；切换
`current`；重启；在限定时间内检查 `/healthz` 和认证 API 请求。

回滚流程：停止不健康版本，保留失败目录；恢复旧 `current` 指针；重启并检查健康和
API；只有单独的恢复决策要求时才恢复数据备份，W-010 不自动执行破坏性恢复。

## 安全与失败行为

- 不记录 API Key、Authorization 头、`users.json` 内容或资源值。
- 启动错误说明失败操作和路径，但不回显秘密。
- 用户文件损坏、权限错误、参数无效或端口占用时，监听前非零退出，不删除数据或
  改变活动指针。
- 健康端点免认证但不包含敏感状态。
- 关闭始终尝试关闭租户存储，即使 HTTP 排空超时。
- 公网 TLS 和证书轮换由代理/部署环境负责。

## 备选方案

| 方案 | 结论 |
| --- | --- |
| 从根目录构建两个二进制 | 拒绝：根目录不是受支持入口，无法表达独立制品 |
| 只发布一个服务端二进制 | 拒绝：权限、健康、数据和恢复行为会继续隐含 |
| 现在给 Go 服务端增加直连 TLS | 后置：扩大密钥生命周期，重复部署代理职责 |
| 优先发布 OCI 镜像 | 后置：仓库没有镜像、registry 和运行时契约；tarball 更容易用现有工具验证 |
| 升级时自动迁移数据库 | 拒绝：迁移兼容性和回滚需要独立版本化设计 |
| 新增服务端配置文件格式 | 拒绝：保留显式 CLI 参数，避免无版本配置迁移 |

## 测试与验证计划

| 层级 | 覆盖内容 |
| --- | --- |
| 单元 | 配置校验、健康状态/方法、就绪切换、关闭清理、归档命名辅助函数 |
| 服务端集成 | 临时 loopback 端口启动，创建用户、探活、认证 CRUD、取消 context、确认存储关闭 |
| CLI 黑盒 | `ttl-server serve --listen ... --port ... --data-dir ...`，无效参数、端口占用、损坏用户文件、缺少客户端文件 |
| 交付冒烟 | 将服务端二进制复制到无源码临时目录，启动、探活、调用 API、SIGTERM、确认数据持久化 |
| 打包 | 检查归档清单和校验值，确认无客户端/源码/数据文件及 Linux 两种架构命名 |
| 升级回滚 | 临时目录中的 v1/v2 和 `current` 指针，成功切换、探活失败、指针恢复和 API 验证 |
| 既有回归 | `go test ./...`、集成测试、`go test -race ./...`、`go vet ./...`、`./scripts/regression.sh`、架构检查 |

## 放行门禁

需求和本方案评审通过、TLS/部署边界确认、WBS 任务具备独立证据后才能编码。本
文档本身不修改生产代码。

## 相关记录

- [`docs/requirements/2026-09-13-server-independent-delivery.md`](../requirements/2026-09-13-server-independent-delivery.md)
- [`docs/decisions/2026-09-12-separate-client-server-layout.md`](../decisions/2026-09-12-separate-client-server-layout.md)
- [`docs/client-server-separation-plan.md`](../client-server-separation-plan.md)
- [`docs/client-server-data-flow.md`](../client-server-data-flow.md)
