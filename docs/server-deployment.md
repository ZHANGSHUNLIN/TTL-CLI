# ttl-server 参考部署

本文是 W-010 面向 Linux 主机、systemd 或等价进程管理器的参考流程，描述部署边界，
不负责创建云账号、DNS、防火墙、反向代理或证书。

## 文件布局

将可执行 release 与可变服务数据分开：

```text
/opt/ttl-server/releases/<version>/ttl-server
/opt/ttl-server/current -> /opt/ttl-server/releases/<version>
/var/lib/ttl-server/users.json
/var/lib/ttl-server/tenants/<user-id>/data.db  # 当前 W-010 实现
# W-019 完成迁移后：
/var/lib/ttl-server/tenants/<user-id>/data.sqlite
```

服务以专用非 root 用户运行。数据目录持久化且只允许 owner 访问；`users.json` 包含
API Key，必须保持 `0600`。发布归档不能包含用户数据、客户端配置或秘密。

W-019 当前仍处于文档和方案阶段，线上已有租户继续使用 `data.db`；不得手动将该文件按
SQLite 打开。完成 W-019 的显式迁移、备份和校验后，目标文件名为 `data.sqlite`，迁移失败
时保留源文件和备份，详见[服务端按租户使用独立 SQLite 文件](decisions/2026-09-13-server-tenant-sqlite.md)。

## 启动契约

进程管理器使用显式参数执行当前 release：

```text
/opt/ttl-server/current/ttl-server serve \
  --listen 127.0.0.1 \
  --port 8080 \
  --data-dir /var/lib/ttl-server
```

如果目标机器限制业务端口范围，应在 `8000~9000` 内选择空闲端口。本次测试机实际绑定
`10.99.48.2:8900`，并同步修改 systemd 的 `ExecStart`；直接绑定私有地址时必须同时配置
防火墙和生产环境 TLS/反向代理边界。

默认监听用于避免意外公网暴露。使用私有服务网络时，防火墙必须只允许可信反向代理或
负载均衡器访问该端口。

进程将运维日志写入 stdout/stderr，由 supervisor 收集；不得记录 API Key、Authorization
头或资源值。`GET /healthz` 是无需 API Key 的就绪探针。

## TLS 边界

公网 HTTPS 在受管反向代理或负载均衡器处终止，代理只转发到私有监听地址并探测
`/healthz`。不支持直接从公网访问明文 HTTP 端口；证书签发和轮换由部署环境负责。

## 安装与验证

1. 下载 Linux 服务端归档及 `.sha256` 文件。
2. 校验 SHA-256 后再解压。
3. 检查归档清单，只允许 `ttl-server` 和发布元数据，不得有 `.git`、源码、客户端或数据。
4. 解压到 `/opt/ttl-server/releases/` 下新的不可变版本目录。
5. 确保运行用户可以执行二进制，只能读写所需的数据目录。
6. 将 `current` 指向新目录，重启 supervisor，先轮询 `/healthz`。
7. 完成一次认证 `/api/v1` 请求后再将版本标记为就绪。

## 升级与回滚

每次升级前停止或静默服务，并对 `/var/lib/ttl-server` 创建一致备份或文件系统快照。新
版本解压到旁边的 release 目录，校验后切换 `current`、重启并等待健康和 API 探针。

若启动或健康检查失败，停止不健康服务，恢复旧的 `current` 目标并重启。失败版本目录
保留供诊断，不自动恢复或删除数据目录。未来不兼容的 schema 迁移必须另立迁移和恢复
决策。

## 运维边界

- `ttl-server user add/list/enable/disable/reset-key/delete` 与服务使用同一个 `--data-dir`。
- 服务端主机不需要客户端二进制或客户端工作区。
- 服务不读取客户端配置，`/healthz` 不返回用户/API Key 数据。
- 日志、备份、防火墙、TLS、证书轮换和 supervisor 策略由部署环境负责。
