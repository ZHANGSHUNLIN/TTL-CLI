# 方案评审：让 ttl get 默认按 key 和 tag 模糊查询

任务：W-016
类型：bugfix
状态：reviewed

## 评审结论

`PASS`。方案将 CLI 默认行为和共享 service 搜索能力解耦：`get` 默认保留 key/tag 查询并排除 value，显式 `--value` 扩展到 value，不破坏 TUI/`pick` 既有 key/value/tag 能力。

## 待补充内容

- 风险评估和方案评审结论
- 负责人：待填写
- 评审人：待填写

## 验收与证据

- [x] 需求覆盖默认 key/tag 与显式 value 条件。
- [x] 兼容性、帮助文案和回滚边界明确。
- [x] 测试可覆盖 service、CLI 和现有入口。

## 备注

这是由任务模板生成的文档骨架，不代表本阶段已经完成。
