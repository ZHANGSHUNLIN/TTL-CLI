# 技术方案：移除当前工程的后端实现

任务：W-021
类型：refactor
状态：ready

## 当前执行路径

`cmd/ttl-server` 通过 `internal/server/cli` 组装用户管理和 `serve` 命令，随后进入 `internal/server/app`、`internal/server/api` 和 `internal/server/tenant`。CI、release、完整验证和 CLI 回归均会构建或运行该入口；`integration_test/server_sync_test.go` 直接导入服务端 handler。

客户端远程能力位于 `internal/client/remote`，只通过 HTTP 与 `/api/v1` 通信，并已有独立 `httptest` 测试。它不需要仓库内服务端实现。

## 方案与归属

| 区域 | 处理 | 最终归属 |
| --- | --- | --- |
| `cmd/ttl-server` | 删除 | 独立后端工程 |
| `internal/server` | 删除 | 独立后端工程 |
| `internal/client/remote` | 保留 | 本仓库客户端 HTTP 适配 |
| `internal/client/sync` | 保留 | 本仓库客户端同步编排 |
| `internal/core` | 保留 | 客户端内部模型与存储契约，不再宣称与仓库内服务端共享 |
| `internal/storage` | 保留 | 客户端本地 SQLite 与兼容性检查 |
| `integration_test/server_sync_test.go` | 删除 | 跨工程契约测试由后续工作建立 |
| 服务端 smoke/deployment | 删除 | 独立后端工程 |
| CI/release | 移除 server job | 只交付 `ttl` |

## 协议与兼容性

不修改客户端请求契约：Bearer 认证、`/api/v1/resources`、`/api/v1/audit/stats` 和 `/api/v1/history` 保持现状。`internal/client/remote/storage_test.go` 继续用本地 mock server 验证请求与响应，不导入后端实现。

这是制品和源码边界的破坏性变化：本仓库不再支持构建 `ttl-server`。客户端 CLI、INI 字段、本地数据库格式和远程协议不变，不需要数据迁移。

## 删除与恢复

删除目标仅限 Git 跟踪的源码、脚本和文档，不访问或删除任何部署目录。若需要回滚，可从本次变更前的 Git 版本恢复；独立后端工程的迁移不属于本任务。

## 测试与门禁

- 架构测试列举 `internal/client`、`internal/core`、`internal/storage` 和 `cmd/ttl`，并显式检查后端目录/入口不存在。
- CLI 黑盒回归删除 `ttl-server --help` 检查。
- 完整验证保留格式、脚本语法、客户端构建、架构、CLI、单元、集成、race 和 vet。
- CI 只保留 lint、test、integration-test 和 client-build。
- release 只打包客户端，发布说明不再宣称包含服务端。

## 文档策略

当前入口文档改写为“客户端 + 外部后端服务”边界。W-010/W-017/W-019 等历史任务文档继续作为当时交付证据；相关 adopted 决策增加被 W-021 替代说明，避免把历史结论误读为当前结构。

## 风险

- 独立后端尚未提供时，cloud 模式无法完成真实端到端验证；以客户端 mock HTTP 契约测试保留本仓库可验证边界。
- 新后端可能与当前 `/api/v1` 漂移；后续应在独立工程或跨工程流水线建立契约测试。
- 历史文档会出现旧路径；通过“历史证据保留、当前文档更新”区分，而不重写已完成任务的事实。
