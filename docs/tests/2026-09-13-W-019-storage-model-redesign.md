# 测试计划与结果：本地与云端存储模式收敛

任务：W-019
类型：feature
状态：自动化测试已执行，人工确认待完成

## 测试范围

覆盖客户端模式解析、SQLite 生命周期、服务端租户隔离、旧格式拒绝和失败恢复。同步、
版本、CRDT 和冲突处理属于 W-020，不在本测试计划中验证。

## 用例矩阵

| 编号 | 场景 | 预期结果 | 层级 |
| --- | --- | --- | --- |
| S-01 | 无配置启动普通命令 | 使用 local，创建/打开 SQLite | 配置单测、CLI |
| S-02 | workspace 与全局模式冲突 | workspace 优先；`--storage` 再覆盖且不写回 | 配置单测 |
| S-03 | `--storage cloud` 缺地址/凭据 | 明确配置错误，不访问 local | CLI、集成 |
| S-04 | cloud API 不可用或认证失败 | 明确远程错误，不降级到 local | httptest、CLI |
| S-05 | `sqlite`、`bbolt`、`sync` 作为 storage | 明确拒绝并提示新模式或独立同步能力，不打开旧存储 | CLI 回归 |
| S-06 | 两个租户并发读写 | 数据互不可见，句柄复用且关闭干净 | 服务端集成、race |
| S-07 | 服务重启后读取租户数据 | 资源、标签、历史、审计和日志保持一致 | 服务端集成 |
| S-08 | 配置多个远程 profile | 当前 workspace 只有一个活动 profile，普通命令只访问该远程 | 配置单测、CLI |
| S-09 | bbolt/旧 SQLite、损坏文件或未知加密状态 | 明确不兼容，文件不读取、不转换、不覆盖、不删除 | 单测、故障注入 |
| S-10 | 租户删除 | 先关闭句柄并备份/隔离，不无条件删除 | 服务端单测、手测 |

## 通过标准

- 所有 S-01 至 S-10 均有自动化或可复现的手动证据。
- 本地和云端对 `core/storage` 的资源、标签、历史、审计和日志基本语义一致。
- 测试日志和报告不包含 API Key、Authorization 头、资源值或完整本地路径中的敏感信息。
- 旧格式被拒绝后原文件保持不变；新数据库初始化和恢复只针对新格式 SQLite。

## 执行命令

实现后按变更范围执行：

```text
go test ./internal/config ./internal/client/... ./internal/storage/... ./internal/server/...
go test -race ./internal/config ./internal/client/... ./internal/storage/... ./internal/server/...
go test ./integration_test/...
./scripts/regression.sh
go vet ./...
git diff --check
```

## 执行结果

以下检查于 2026-09-13 在仓库根目录执行并通过：

```text
go test ./...
go test -race ./...
go test ./integration_test/...
./scripts/regression.sh
./scripts/cli-composability.sh
go vet ./...
git diff --check
./scripts/verify.sh
```

覆盖证据包括：local/cloud 模式解析和远程 profile、SQLite 生命周期和旧格式拒绝、租户隔离
与删除隔离、副作用边界、CLI 黑盒行为、组合式 JSON 契约、竞态检查和静态检查。人工部署权限、
模式切换和日志敏感信息确认仍保留在验收清单中。
