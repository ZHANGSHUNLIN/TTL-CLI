# 客户端与服务端代码结构收敛测试计划

日期：2026-09-13  
任务：W-017  
需求：[`docs/requirements/2026-09-13-W-017-code-structure-convergence.md`](../requirements/2026-09-13-W-017-code-structure-convergence.md)  
方案：[`docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md`](../tech-designs/2026-09-13-W-017-code-structure-convergence.md)  
状态：completed

## Test Objective

证明代码结构迁移完成后具备唯一、可理解、可自动检查的包边界，同时既有 CLI、TUI、服务端 API、存储格式、加密格式和同步语义不回归。

本任务不以“目录存在”作为证据；一次性切换完成后已同时验证 import 图、真实行为和文档状态。T-01 至 T-06 只是切换前后的内部执行顺序。

## Baseline Checks

- [x] 保存迁移前 `go list -json ./...` 的包和依赖快照（本轮盘点记录在 W-017 方案中）。
- [x] 记录生产代码仍引用旧顶层包的清单，并在切换后确认清零。
- [x] 记录现有 CLI 黑盒、集成、架构和 race 检查基线结果。
- [x] 确认 `go build -o ttl ./cmd/ttl` 和 `go build -o ttl-server ./cmd/ttl-server` 基线通过。

## Unit Tests

### Core and model migration

- [x] `internal/core/resource` 类型字段、JSON tag、排序和比较行为保持不变。
- [x] `internal/core/storage` 接口的所有实现仍覆盖资源、标签、审计、历史和日志操作。
- [x] 删除或停用 model alias 后，核心包测试不依赖旧路径。

### Client command adapters

- [x] 命令从 context 获取显式 service，缺少 service 时返回稳定错误而不 panic。
- [x] 没有包级全局 writer 交叉污染；命令运行保持既有 stdout/stderr 行为。
- [x] 参数校验、机器 JSON、非交互、退出码、历史快捷参数和本地化 key 与迁移前一致。

### Config, crypto, i18n, text

- [x] 默认配置和 workspace 路径解析保持不变。
- [x] 加密/解密、密钥生成、导入导出、权限和错误行为保持不变。
- [x] 所有 locale key 可加载，语言选择和 fallback 行为保持不变。
- [x] 文本去重、大小写匹配和转义规则保持不变；无调用 helper 不再编译进产品。

### Sync

- [x] `local_only`、`remote_only`、`conflict`、`in_sync` 计算保持不变。
- [x] pull、push、auto、dry-run、镜像写入顺序和 close 行为保持不变。

## Architecture Tests

- [x] client 不依赖 server。
- [x] server 不依赖 client、TUI、客户端命令适配器或根目录 `sync`。
- [x] core 不依赖 client、server、storage adapter、Cobra、HTTP 或本地配置。
- [x] `cmd/ttl` 不依赖 `ttl-cli/command`、`ttl-cli/models`、`ttl-cli/db` 或根 `ttl-cli/sync`。
- [x] 迁移完成后仓库无 `ttl-cli/db`、`db.Stor`、`ttl-cli/models` 和根 `ttl-cli/sync` 的生产/测试引用。
- [x] 直接依赖和传递依赖均通过 `go list -json` 检查，而不是只检查源文件文本。

命令：

```bash
go test ./internal/architecture
```

## Integration Tests

- [x] 资源生命周期：add/get/update/tag/dtag/del/rename。
- [x] 工作区、配置路径和临时 HOME 隔离。
- [x] 加密启用、迁移、解密、密钥导入导出和错误恢复。
- [x] 导入/导出、审计、历史和日志。
- [x] HTTP API、API Key、用户状态和租户隔离。
- [x] local/cloud/sync 存储以及 push/pull/冲突处理。

命令：

```bash
go test ./integration_test/...
```

## CLI And Build Regression

- [x] 使用临时 `HOME` 和配置文件运行 CLI 黑盒回归，不触碰真实 `~/.ttl`。
- [x] 验证文本输出、JSON 输出、stdin/stdout/stderr、无 TTY、稳定退出码和命令帮助。
- [x] 验证 `ttl ui` 的资源浏览、编辑、标签、删除、打开和终端恢复行为不变。
- [x] 验证客户端二进制不链接 server；服务端二进制不链接 client、旧 command/db/sync。

命令：

```bash
./scripts/regression.sh
./scripts/cli-composability.sh
go build -o ttl ./cmd/ttl
go build -o ttl-server ./cmd/ttl-server
```

## Full Verification

```bash
gofmt -s -l .
go test ./...
go test -race ./...
go vet ./...
./scripts/verify.sh
git diff --check
```

## Evidence And Failure Handling

- 在一次性切换的任务记录中保存实际运行过的命令和结果，不用“计划执行”代替证据。
- 任一行为回归、架构越界、数据格式变化或测试仍依赖旧包时，停止切换并整体回滚 W-017。
- 测试迁移必须使用临时目录、`httptest` 或显式 service/storage 生命周期，不能为了通过测试重新引入全局状态。
- 全量验证通过后，先进行 owner code review，再更新 acceptance；没有整体迁移的本地 commit 不进入 `Done`。
