# 需求分析：移除当前工程的后端实现

任务：W-021
类型：refactor
状态：ready

## 问题与目标

当前仓库同时维护 `ttl` 客户端和 `ttl-server` 后端。后续后端将由专门的服务工程负责，继续保留 `cmd/ttl-server`、`internal/server`、服务端发布任务和部署脚本会形成两份实现与两条发布链路。

本任务把当前仓库收敛为客户端工程：保留本地存储、远程 HTTP 适配器、同步客户端和远程协议测试，移除后端启动、HTTP handler、认证、租户存储、用户管理及其交付链路。

## 范围

### 包含

- 删除 `cmd/ttl-server` 和 `internal/server`。
- 删除只服务于内置后端的单元测试、集成测试、冒烟脚本和部署说明。
- CI 和 release 只构建、测试、发布 `ttl` 客户端。
- `scripts/regression.sh` 与 `scripts/verify.sh` 只验证客户端。
- 保留 `internal/client/remote`、cloud 配置和 `ttl sync`，其服务端由独立工程提供。
- 用 `httptest` 验证客户端现有 `/api/v1` 请求/响应契约，不从新后端工程复制实现。
- 更新工程边界、README、项目导览和架构检查。

### 不包含

- 创建、迁移或部署新的后端工程。
- 修改远程 API 路径、认证头、DTO 或同步语义。
- 删除 cloud 模式、远程 profile 或 `ttl sync`。
- 删除仅因历史原因存在的已完成任务证据；历史文档保留其当时事实，但当前入口文档不得继续宣称本仓库提供后端。
- 迁移现有服务端数据、账号或 API Key。

## 兼容与失败规则

- `ttl` 的本地模式、cloud 模式、命令参数和配置格式保持不变。
- 本仓库不再生成 `ttl-server` 二进制或发布物；依赖该制品的使用者必须转向独立后端服务。
- 客户端访问不可用或不兼容的远端时继续返回现有网络或 API 错误，不回退到内置后端。
- 删除仓库后端源码不会删除任何已部署服务、用户目录或租户数据；本任务不操作仓库外数据。

## 验收标准

- [ ] `cmd/ttl-server`、`internal/server` 和服务端专属脚本不存在。
- [ ] 源码、当前入口文档和 workflow 不再依赖或构建仓库内后端。
- [ ] `go build -o <temp>/ttl ./cmd/ttl` 通过。
- [ ] `go test ./...`、`go test ./integration_test/...`、`go test -race ./...` 和 `go vet ./...` 通过。
- [ ] `./scripts/regression.sh` 和 `./scripts/verify.sh` 通过，且不构建或启动服务端。
- [ ] `internal/client/remote` 的 mock HTTP 测试继续覆盖认证、CRUD、错误和超时。
- [ ] README 和 `PROJECT_OVERVIEW.md` 明确当前仓库只提供客户端，远端服务由独立工程提供。

## 假设与开放项

- 假设独立后端继续兼容客户端当前 `/api/v1` 契约；协议演进由后续跨工程工作单独处理。
- 独立后端仓库地址和部署文档尚未提供，因此本仓库只说明外部依赖，不添加不可验证链接。
