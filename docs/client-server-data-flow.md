# ttl 客户端本地与远端数据流

状态：当前实现
更新：2026-09-18

## 1. 边界

`ttl` 是本仓库唯一应用，CLI 与 TUI 共享客户端用例。后端由独立服务工程提供；本仓库只维护 HTTP 客户端契约，不维护或部署后端实现。

```text
ttl command / TUI
  -> internal/client/app.Service
     -> local -> SQLite
     -> cloud -> internal/client/remote -> 外部后端

ttl sync
  -> local SQLite + remote HTTP Storage
  -> ComputeDiff
  -> ExecutePull / ExecutePush
```

## 2. local 模式

`--storage local` 通过 `internal/storage/sqlite` 访问本地数据库。默认配置和工作空间决定数据库路径。普通资源、标签、审计、历史和日志操作均通过 `internal/client/app.Service` 执行。

local 模式不要求网络或外部服务。测试与手动检查必须使用临时 `HOME`、`--conf` 和临时数据库。

## 3. cloud 模式

`--storage cloud` 通过 `internal/client/remote.Storage` 直接访问已配置外部后端，不维护本地业务副本。客户端支持：

- 服务地址与超时；
- `Authorization: Bearer <credential>`；
- `GET/POST /api/v1/resources`；
- `PUT/DELETE /api/v1/resources/{key}`；
- `GET /api/v1/audit/stats`；
- `GET /api/v1/history`。

远端返回统一 JSON envelope：`code`、`message` 和可选 `data`。`code != 0` 时客户端返回 API 错误；连接、超时、序列化和解析错误保留上下文。当前 mock 契约见 `internal/client/remote/storage_test.go`。

外部后端的用户、API Key 颁发、租户隔离、数据库布局和部署方式不属于本仓库。

## 4. 远程配置

远程 profile 保存 URL、超时和凭据环境变量名，不保存明文凭据。客户端解析当前应用或 workspace 的活动 profile；也可通过 `--cloud-url`、`--cloud-key` 和 `--cloud-timeout` 显式覆盖。

缺少 URL 或凭据时 cloud 与 sync 会明确失败，不自动回退到 local。

## 5. 同步

`ttl sync` 始终以本地 SQLite 和远端 HTTP Storage 为两个独立数据源。当前同步按 key 比较资源：

- `local_only`：仅本地存在；
- `remote_only`：仅远端存在；
- `conflict`：同 key 的 value 或 tags 不同；
- 相同内容视为同步。

pull 让本地匹配远端，push 让远端匹配本地；`--dry-run` 只展示差异。现有行为没有 revision、幂等键或 CRDT，相关演进由 W-020 另行设计。

## 6. 已知边界

- 本仓库只用 `httptest` 验证客户端协议，不能证明独立后端的真实兼容性。
- 当前协议没有版本协商；跨工程变更必须协调客户端和独立后端。
- 创建请求包含 tags，更新请求当前只发送 value；协议完善需要独立需求。
- cloud 模式的部分审计、历史和日志写入仍是空实现或有限实现。
- 同步覆盖目标侧，执行前应先使用 dry-run 检查差异。

## 7. 代码索引

| 职责 | 路径 |
| --- | --- |
| 客户端组合与生命周期 | `internal/client/cli/root.go` |
| 客户端用例 | `internal/client/app/` |
| 远端 HTTP 适配 | `internal/client/remote/storage.go` |
| 远端 mock 契约测试 | `internal/client/remote/storage_test.go` |
| 同步 diff 与执行 | `internal/client/sync/` |
| 本地 SQLite | `internal/storage/sqlite/` |
| 配置与远程 profile | `internal/config/` |
