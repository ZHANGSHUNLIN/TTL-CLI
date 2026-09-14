# 方案评审：本地与云端存储模式收敛

任务：W-019
类型：feature
评审日期：2026-09-14
评审人：owner
结论：PASS

## 已确认的设计覆盖

- local/cloud 互斥访问目标、workspace 优先级和默认 local 已定义。
- 远程 profile 可以配置多个，但当前应用或 workspace 只有一个活动 profile；一次普通命令
  只访问当前 workspace 的一个远程。
- 新版本不兼容旧 bbolt/旧 SQLite，不提供迁移命令、自动转换或兼容窗口。
- 服务端每租户独立 SQLite，`users.json` 暂不迁移，租户路径不由请求指定。
- SQLite WAL、busy timeout、连接生命周期、前向 schema 版本、备份、恢复和删除前置检查已
  纳入方案。
- cloud 只访问 HTTP API；同步、版本、CRDT 和冲突处理明确由 W-020 负责。
- 单元、集成、CLI 回归和失败恢复的验证层次已列出。

## 评审门禁

方案已经具备进入 WBS 和编码的条件。本记录的评审人、日期和结论已补齐；`PASS` 只表示
设计可执行，不代表代码、测试或验收已经完成。
