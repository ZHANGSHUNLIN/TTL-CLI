# W-010 测试计划：云端服务独立交付

日期：2026-09-13
任务：W-010
状态：自动检查和测试环境验收通过，待 commit

## 确定性检查

### 构建与架构

- [x] `go build -o <temp>/ttl ./cmd/ttl`
- [x] `go build -o <temp>/ttl-server ./cmd/ttl-server`
- [x] `go test ./internal/architecture`
- [x] `gofmt -s -l .`
- [x] `go vet ./...`

CI 和 release workflow 已检查为从明确的 `cmd/` 路径构建，而不是从已删除的根入口构建。
真实 GitHub tag 运行记录仍待人工发布演练。

### 运行时单元与集成

- [x] 校验空监听地址、端口范围、数据目录和关闭超时；失败时不得打开文件或 socket。
- [x] 使用临时数据目录和 loopback 临时端口启动应用，确认就绪后 `/healthz` 返回 200，且
      不含凭据。
- [x] 使用临时用户存储执行一次认证 `/api/v1` 请求。
- [x] 取消 context 或发送 SIGTERM，确认租户存储关闭，进程在关闭超时内退出。
- [x] 损坏 `users.json`、数据目录路径不可用和端口占用时，确认服务在提供请求前失败，且不删除
      已有数据；`internal/server/app/server_test.go` 已覆盖这些进程级失败边界。

### 打包与干净目录冒烟

- [x] 构建服务端二进制，只将它复制到干净临时 release 目录，并使用另一个临时数据目录运行。
- [ ] 检查两个 Linux 架构的真实发布归档清单和 SHA-256 校验值。当前已完成本地等价的归档规则
      和服务端独立交付冒烟，真实 tag 制品仍待发布后核对。
- [x] 确认本地服务端交付冒烟不带客户端二进制、源码、`.git`、`users.json`、租户数据库或客户端配置。
- [x] 探测 `/healthz`，通过服务端 CLI 创建用户，发起一次认证 API 请求，停止进程，并确认
  数据仍在数据目录中。

### 升级与回滚

- [x] 在临时部署目录创建 `releases/v1`、`releases/v2` 和 `current` 符号链接。
- [x] 启动 v1，备份数据，切换 v2，重启并验证健康和 API 探针。
- [x] 让 v2 探活失败，恢复 v1 指针，重启并确认 v1 仍能读取升级前数据。
- [x] 确认失败 release 保留供诊断，不会自动执行破坏性恢复。

## 既有回归

本轮已执行并通过：

```text
go test ./internal/server/...
go build -o /tmp/ttl-server-w010 ./cmd/ttl-server
scripts/server-delivery-smoke.sh /tmp/ttl-server-w010
scripts/server-upgrade-rollback-smoke.sh /tmp/ttl-server-w010
./scripts/verify.sh
git diff --check
jq empty WORK_ITEMS.json
jq empty internal/i18n/locales/*.json
bash -n scripts/*.sh
Ruby YAML 解析 .github/workflows/ci.yml
Ruby YAML 解析 .github/workflows/release.yml
```

`./scripts/verify.sh` 已覆盖客户端/服务端构建、交付和回滚冒烟、架构检查、CLI 黑盒回归、
全量 Go 测试、集成测试、race 和 vet。

## 远程部署验证

2026-09-13 在 Debian 12 x86_64 测试机 `10.99.48.2:8022` 上完成一次实际部署：

- 本地交叉编译 `GOOS=linux GOARCH=amd64 CGO_ENABLED=0` 服务端二进制并上传。
- 创建 `ttl-server` 系统用户，发布目录为 `/opt/ttl-server/releases/2026-09-13-w010`，
  当前指针为 `/opt/ttl-server/current`。
- 创建 `/etc/systemd/system/ttl-server.service`，服务以非 root 用户运行，数据目录为
  `/var/lib/ttl-server`，监听 `10.99.48.2:8900`（该机器要求业务端口位于 `8000~9000`）。
- `systemctl enable --now ttl-server.service` 成功，服务状态为 `active`，重启后仍能就绪。
- 远端本机 `GET /healthz` 返回 `{"status":"ok"}`；未认证 API 返回 `401`。
- 创建临时用户后重启服务，认证 `GET /api/v1/resources` 返回成功响应；随后删除临时用户并
  再次重启，最终数据目录只保留空用户文件。
- 本机直接访问 `http://10.99.48.2:8900/healthz` 返回 `{"status":"ok"}`，无需 SSH 隧道。
- 使用临时 HOME 的 `ttl --storage cloud --cloud-url http://10.99.48.2:8900` 完成远程
  `add`、`get`、`update`、`del`；删除后再次 `get` 正确返回资源不存在。测试用户随后已清理。
- 服务只绑定该私有地址，没有配置公网 TLS；生产环境仍需防火墙和反向代理/TLS 边界。

## 人工审核

- [x] 负责人确认本次测试环境无需安装客户端或准备源码 checkout 即可完成参考部署。
- [x] 负责人确认本次测试环境的私有监听、数据权限和日志收集符合部署演练；公网 TLS 边界留待生产环境。
- [x] 升级失败和回滚在临时目录中可观察、可重复。
- [x] 测试不依赖真实 home 目录、生产数据或公网服务。
