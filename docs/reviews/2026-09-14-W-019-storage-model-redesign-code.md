# 代码评审：本地与云端存储模式收敛

任务：W-019
评审日期：2026-09-14
评审人：owner
结论：PASS

## 评审范围

本次评审覆盖客户端 `local`/`cloud` 存储选择、SQLite 本地和租户实现、远程 profile
解析、旧格式拒绝、租户生命周期、CLI 回归脚本、集成测试和相关文档。工作区中与 W-020
同步协议、项目流程语言规范及其他任务相关的独立改动不纳入本次提交。

## 关键检查

- `OpenStorage` 只创建本地 SQLite 或远程 HTTP storage，`sqlite`、`bbolt`、`sync` 等旧
  storage 值明确拒绝；`sync` 仍作为独立命令访问 local 与 cloud。
- 配置优先级为命令行 > workspace > 全局 > 默认 `local`；远程 profile 只保存地址、账户
  和凭据环境变量名，不记录 API Key。
- SQLite 初始化检查文件签名、schema 版本和旧同名文件；不兼容输入不会读取、转换、覆盖
  或删除；错误路径关闭已打开的数据库句柄。
- 服务端按租户使用独立 `data.sqlite`，租户目录和文件权限收紧；删除前关闭句柄，关闭失败
  时停止操作，成功后将目录移动到 `.deleted` 隔离副本。
- workspace 展示只读取 `local` SQLite；cloud workspace 不会误打开本地数据库。
- 代码、测试、脚本和文档未写入远程 API Key 或其他凭据。

## 验证证据

以下命令均通过：

```text
go test ./...
go test -race ./...
go test ./integration_test/...
./scripts/regression.sh
./scripts/cli-composability.sh
./scripts/verify.sh
go vet ./...
git diff --check
```

远程节点 `10.99.48.2:8900` 已完成部署闭环：健康检查返回 `{"status":"ok"}`，未认证请求
返回 `401`，认证资源创建/读取成功，租户 SQLite 文件权限为 `0600`，删除后进入
`tenants/.deleted` 隔离副本。临时用户已删除，服务已重启，当前 release 为
`/opt/ttl-server/releases/2026-09-14-w019-r2`。

## 遗留事项

- W-020 的版本、同步、CRDT 和冲突协议不在本次实现范围内。
- 远程用户管理目前由独立 CLI 修改 `users.json`，运行中的服务需要重启才能加载用户变更；
  这属于现有服务生命周期约束，未在 W-019 扩大处理范围。
