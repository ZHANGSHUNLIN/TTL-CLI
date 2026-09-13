# 需求分析：让 ttl get 默认按 key 和 tag 模糊查询

任务：W-016
类型：bugfix
状态：reviewed

## 问题与目标

`ttl get <query>` 当前只按资源 key 搜索，带标签的资源无法通过标签获取；value 命中也不应在默认查询中产生候选。

目标：`ttl get` 默认对 key 和 tag 做不区分大小写的子串匹配；需要按内容查找时提供显式 `--value` 条件。

## 范围

- `ttl get <query>` 默认匹配 ORIGIN 资源 key 或 tags。
- `ttl get --value <query>` 匹配 key、value 和 tags。
- TUI 和 `ttl pick` 继续沿用 key/value/tag 的既有搜索语义。
- 不改变精确删除、更新、标签命令或存储格式。

## 待补充内容

### 验收场景

1. 资源 key 为 `deployment-note`、value 为 `contains-secret` 时，`ttl get secret` 不应命中。
2. 资源 tag 为 `production` 时，`ttl get production` 应命中并输出 value。
3. 同一资源执行 `ttl get --value secret` 应命中并输出 value。
4. key/tag 命中和 `--value` 的 key/value/tag 命中都保持原有单条输出、多条交互和非交互歧义语义。
5. TUI/`pick` 仍能按 value 和 tag 搜索。

## 验收与证据

- [x] 默认 key/tag 与显式 `--value` 行为可通过 CLI 和 service 测试观察。
- [x] TUI/`pick` 搜索行为保持兼容。

## 备注

这是由任务模板生成的文档骨架，不代表本阶段已经完成。
