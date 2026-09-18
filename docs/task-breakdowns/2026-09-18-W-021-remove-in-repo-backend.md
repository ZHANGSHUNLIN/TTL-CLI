# WBS 拆分：移除当前工程的后端实现

任务：W-021
状态：approved

## T-01 删除后端实现与专属测试

- 目标：移除本仓库可执行后端和所有仅验证该实现的测试。
- 输入：通过评审的 W-021 方案。
- 输出：删除 `cmd/ttl-server`、`internal/server`、`integration_test/server_sync_test.go`。
- 依赖：无。
- 负责区域：Go 入口、服务端包、集成测试。
- 验收：`rg` 无生产或测试 import 指向 `internal/server`；客户端包测试可编译。
- 检查：`go test ./internal/client/... ./internal/architecture`。

## T-02 收敛构建、发布和验证链路

- 目标：确保本仓库只构建和发布客户端。
- 输入：T-01。
- 输出：删除服务端 smoke 脚本；更新 CI、release、regression、verify 和架构检查。
- 依赖：T-01。
- 负责区域：`.github/workflows`、`scripts`、`internal/architecture`。
- 验收：所有验证入口不引用 `cmd/ttl-server` 或服务端脚本。
- 检查：`bash -n scripts/*.sh`、`go test ./internal/architecture`、workflow 人工核对。

## T-03 更新当前工程文档与决策

- 目标：让用户和维护者不会在本仓库寻找或构建后端。
- 输入：T-01、T-02。
- 输出：更新 AGENTS、README、项目导览、数据流/边界文档和相关决策；删除当前服务端部署说明。
- 依赖：T-01、T-02。
- 负责区域：根目录文档、`docs/`、W-021 证据。
- 验收：当前入口文档明确外部后端边界，历史证据仍可识别为历史。
- 检查：链接和路径人工核对、`git diff --check`。

## T-04 全量验证与 owner 评审

- 目标：证明删除后客户端行为和工程门禁完整。
- 输入：T-01～T-03。
- 输出：测试记录、交付验收和代码评审结论。
- 依赖：T-01～T-03。
- 负责区域：全仓库。
- 验收：W-021 需求中的自动化检查全部通过；任务进入 `Review`，等待用户决定是否提交。
- 检查：`./scripts/verify.sh`、`git diff --check`、`git status --short`。
