# 客户端与外部后端边界

状态：已采用
更新：2026-09-18
决策：[`后端服务移出 ttl-cli 仓库`](decisions/2026-09-18-externalize-backend-service.md)

## 当前边界

本仓库只包含 `ttl` 客户端。后端是独立工程、独立制品和独立运行环境，不在本仓库构建、测试、发布或部署。

```text
cmd/ttl
  -> internal/client/cli
  -> internal/client/app
  -> local: internal/storage/sqlite
  -> cloud: internal/client/remote -> HTTP -> 独立后端服务
```

## 依赖规则

1. `cmd/ttl` 是唯一产品可执行入口。
2. `internal/client/remote` 是远端访问唯一入口，只依赖 HTTP、认证头和 JSON DTO。
3. 本仓库不得出现 `internal/server`、后端 handler、用户管理、租户存储或服务端发布实现。
4. `internal/core` 是客户端内部契约，不作为与后端共享源码的承诺。
5. 远端协议变化必须更新客户端 mock HTTP 测试，并由跨工程协作确认后端兼容性。

## 保留能力

- local/cloud 两种 storage mode。
- 远程 profile、API 地址和凭据环境变量。
- Bearer 认证与当前 `/api/v1` 客户端 DTO。
- `ttl sync` 的 diff、pull、push 和 dry-run。

## 不属于本仓库

- 后端 API 实现、认证策略和租户模型。
- 服务端数据库、数据迁移、账号与 API Key 管理。
- 后端二进制、容器、CI、release、部署、健康检查和回滚。
- 真实服务端兼容性与生产运维。

## 验证

```bash
go test ./internal/architecture
go test ./internal/client/remote ./internal/client/sync
go test ./...
go test ./integration_test/...
go test -race ./...
go vet ./...
./scripts/regression.sh
```

架构测试检查后端入口和实现不存在；远端 mock 测试检查客户端协议。独立后端准备好后，应在独立工程或跨工程流水线增加真实契约/端到端验证。
