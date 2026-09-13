# 技术方案：让 ttl get 默认按 key 和 tag 模糊查询

任务：W-016
类型：bugfix
状态：reviewed

## 当前行为

`internal/client/app.Service.FindResourcesWithOptions` 为 TUI、`pick` 和 CLI `get` 提供可配置的 key/value/tag 搜索。CLI `get` 当前默认只启用 key，导致 tag 命中缺失。

## 方案

- 在 application service 增加 `SearchOptions` 和 `FindResourcesWithOptions`；key 始终参与匹配，value/tag 由选项控制。
- 保留 `FindResources` 的全字段默认语义，避免改变 TUI 和 `pick`。
- CLI `get` 增加 `--value` bool flag；默认调用 key+tag 选项，传入 flag 时调用 key+value+tag 选项。
- 更新五种 locale 的 `get` 长帮助和 flag 文案。

## 待补充内容

## 兼容性与风险

只改变 `ttl get` 的默认搜索范围；service 的既有 `FindResources`、TUI 和 `pick` 不变。回滚只需恢复 CLI flag 选择和 service 选项调用，不涉及数据迁移。

## 验收与证据

- [x] service 和 CLI 归属明确。
- [x] 默认、显式 value、既有入口和本地化帮助均有测试/检查计划。

## 备注

这是由任务模板生成的文档骨架，不代表本阶段已经完成。
