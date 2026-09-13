# W-010 代码评审：云端服务独立交付

日期：2026-09-13
任务：W-010
评审范围：服务端运行时、健康检查、CLI 参数、CI/release workflow、交付冒烟脚本和回归测试
评审人：owner
结论：PASS

## 评审依据

- 需求：[`docs/requirements/2026-09-13-server-independent-delivery.md`](../requirements/2026-09-13-server-independent-delivery.md)
- 技术方案：[`docs/tech-designs/2026-09-13-server-independent-delivery.md`](../tech-designs/2026-09-13-server-independent-delivery.md)
- 方案评审：[`docs/reviews/2026-09-13-server-independent-delivery-design.md`](2026-09-13-server-independent-delivery-design.md)
- WBS：[`docs/task-breakdowns/2026-09-13-server-independent-delivery.md`](../task-breakdowns/2026-09-13-server-independent-delivery.md)
- 决策：[`docs/decisions/2026-09-13-server-independent-delivery-baseline.md`](../decisions/2026-09-13-server-independent-delivery-baseline.md)

## 检查结果

### 生命周期和资源

- `internal/server/app` 在绑定端口前完成配置、数据目录和用户文件检查。
- 默认监听地址为 loopback；监听器由应用层创建，关闭通过 `http.Server.Shutdown` 执行。
- SIGINT、SIGTERM 和 context 取消都会进入有限时长的关闭流程，并始终调用租户存储
  `CloseAll`。
- `/healthz` 位于认证中间件之前，只返回就绪状态，不读取或返回用户、API Key 和资源数据。

### 兼容性和范围

- `NewHandler(storage)` 保持 API 测试入口；现有 `/api/v1` 路由和认证中间件未改变。
- `ttl-server user` 命令继续使用同一 `--data-dir`，客户端入口和数据格式没有迁移。
- W-017 的服务端边界没有被重新复制；新增生命周期归属集中在 `internal/server/app`。

### 制品和脚本

- CI/release 从 `./cmd/ttl` 和 `./cmd/ttl-server` 构建，服务端只发布 Linux `amd64`/`arm64`。
- 服务端归档清单在 workflow 中显式检查；交付脚本覆盖干净目录、健康/API 探针、关闭、
  版本切换和失败回滚。
- 脚本使用临时目录和 loopback 端口，不读取真实 home 目录或生产数据。

### 测试证据

以下检查已通过：

```text
./scripts/verify.sh
go test ./internal/server/...
git diff --check
jq empty WORK_ITEMS.json
jq empty internal/i18n/locales/*.json
bash -n scripts/*.sh
```

`./scripts/verify.sh` 已包含双二进制构建、交付/升级回滚冒烟、架构依赖、CLI 黑盒回归、
全量 Go 测试、集成测试、race 和 vet。新增测试覆盖损坏 `users.json`、数据目录路径不可用
和端口占用时在监听前失败。

## 发现项

没有发现需要阻塞合并的正确性、安全性、数据安全、资源生命周期或兼容性问题。

以下事项属于交付验收，不是本次代码评审阻塞项：

- 需要负责人对真实 GitHub tag 归档清单和 SHA-256 文件做一次核对。
- 需要负责人确认实际部署环境中的非 root 用户、数据卷、日志收集和外部 TLS 边界。

## 放行结论

代码变更通过 owner review，测试环境验收也已完成，可以进入本地 commit 门禁。真实 tag
制品核对和生产环境外部 TLS 配置属于后续交付核对，不影响本次测试环境验收结论。
