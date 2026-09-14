# 交付验收：本地与云端存储模式收敛

任务：W-019
类型：feature
状态：自动验收、人工确认与 owner code review 通过

## 验收前提

方案评审结论为 `PASS`，实现、自动化测试、人工确认和 owner code review 均已完成。

## 自动验收

- [x] CLI 帮助和配置只暴露 `local`、`cloud`；旧 storage 值明确拒绝并提示新模式。
- [x] local 只访问本地 SQLite，cloud 只访问远程 API；失败时不自动切换另一端。
- [x] 模式优先级为命令行 > workspace > 全局 > 默认 local；单次覆盖不写回配置。
- [x] 服务端按租户使用独立 `data.sqlite`，租户间资源、标签、历史、审计和日志互不可见。
- [x] SQLite 启用约定的并发参数并执行 schema 迁移；服务重启后数据可读。
- [x] 旧 bbolt/旧 SQLite 被明确拒绝，不读取、不转换、不覆盖、不删除，`users.json` 不受影响。
- [x] 可配置多个远程 profile，但当前应用或 workspace 只有一个活动 profile，一次普通命令
  只访问该远程。
- [x] 租户删除先关闭句柄并备份或隔离，未调用无条件 `os.RemoveAll`。
- [x] 项目单测、race、集成、CLI 回归、`go vet` 和 `git diff --check` 均有记录。

## 人工确认

- [x] 检查部署目录权限：服务用户可读写租户目录，`users.json` 保持 `0600`；远程节点租户目录为
  `0700`，租户 SQLite 文件为 `0600`。
- [x] 使用临时租户演练模式切换，确认另一端数据没有复制、删除或覆盖。
- [x] 使用旧格式文件进行启动演练，确认收到不兼容错误且源文件未被修改。
- [x] 检查部署日志和测试输出未记录 API Key 或 Authorization 头；资源值仅出现在受控的 API
  响应验证中，未写入交付日志。
- [x] 确认 W-019 代码和文档没有引入 W-020 的同步、版本、outbox、CRDT 或冲突协议。

## 遗留风险

- 远程 profile 的配置字段已落地，API Key 仅通过命令参数或环境变量读取。
- 若未来并发量超出单机 SQLite 适用范围，应另立 PostgreSQL 或共享数据库决策，不在本
  任务中临时改变存储模型。

## 证据

## 自动验收证据

- 2026-09-13：`./scripts/verify.sh` 全部 10 个阶段通过。
- 2026-09-13：`./scripts/regression.sh`、`./scripts/cli-composability.sh` 通过。
- 2026-09-13：`go test ./...`、`go test -race ./...`、`go test ./integration_test/...`、
  `go vet ./...` 和 `git diff --check` 通过。
- 2026-09-14：构建 `linux/amd64` 服务端并部署到 `10.99.48.2:8900`，当前 release 为
  `/opt/ttl-server/releases/2026-09-14-w019-r2`；`/healthz` 返回 `{"status":"ok"}`。
- 2026-09-14：远端临时租户闭环通过：未认证请求返回 `401`，认证创建与读取资源成功，租户
  `data.sqlite` 权限为 `600`，删除后目录进入 `tenants/.deleted` 隔离副本；临时用户已删除，
  服务已重启并确认无残留活动用户。
- 2026-09-14：本地临时配置演练确认 `local` 使用 SQLite、`cloud` 使用活动远程 profile，
  切换不复制或覆盖另一端；`sqlite`、`bbolt`、`sync` 及旧 SQLite schema 均被拒绝，源文件
  保持不变。

提交号在本地 commit 完成后回写到任务记录。
