# WBS 拆分：让 ttl get 默认按 key 和 tag 模糊查询

任务：W-016
类型：bugfix
状态：reviewed

## WBS

### T-01 搜索选项与 service 行为

- 输入：现有 `FindResources` key/value/tag 逻辑。
- 输出：可控制 value/tag 是否参与的搜索选项，既有全字段入口保持兼容。
- 依赖：无。
- 验收：默认 key/tag 不命中 value-only 资源且能命中 tag-only 查询；开启 value 后可命中 value。

### T-02 CLI get 条件与帮助

- 输入：T-01、现有 Cobra 命令和 locale。
- 输出：`--value` flag，默认 key/tag，帮助文案同步。
- 依赖：T-01。
- 验收：CLI 默认和 `--value` 行为正确，TUI/`pick` 不回归。

### T-03 测试与交付证据

- 输入：T-01/T-02。
- 输出：service/CLI 测试、检查结果和验收记录。
- 依赖：T-02。
- 验收：自动检查通过，任务进入 Review。

## 待补充内容

- 修复步骤、回归范围和依赖
- 负责人：待填写
- 评审人：待填写

## 验收与证据

- [x] 每项包含输入、输出、依赖和验收。

## 备注

这是由任务模板生成的文档骨架，不代表本阶段已经完成。
