# 测试计划与结果：移除当前工程的后端实现

任务：W-021
状态：passed

## 覆盖目标

- 后端入口和实现目录不存在。
- 客户端不导入已移除的后端包。
- cloud HTTP 适配器继续通过本地 mock server 验证。
- 客户端构建、CLI 黑盒、本地集成、race 和静态检查通过。
- 验证脚本不构建或启动后端。

## 执行结果

2026-09-18 执行 `./scripts/verify.sh`，以下步骤全部通过：

- `git diff --check`、`gofmt -s -l .` 和脚本 `bash -n`；
- `go build -o <temp>/ttl ./cmd/ttl`；
- `go test ./internal/architecture`；
- `./scripts/regression.sh <temp>/ttl`；
- `go test ./...`；
- `go test ./integration_test/...`；
- `go test -race ./...`；
- `go vet ./...`。

`go test ./internal/client/...` 也单独通过，其中 `internal/client/remote` 使用 `httptest` 覆盖认证、资源 CRUD、服务端错误和超时。

## 未执行

没有连接真实独立后端做端到端验证，因为其地址和交付物不在本任务输入中。该风险不影响仓库内后端删除，但协议兼容需要后续跨工程验证。
