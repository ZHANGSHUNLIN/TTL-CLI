# Store Work Item State In JSON

日期：2026-09-12
状态：adopted
任务：WORK_ITEMS.md 维护

## 背景

把任务正文、状态章节和复选框全部放在 `WORK_ITEMS.md` 中，适合阅读，但流程工具或 Agent 更新状态时需要移动整段 Markdown，容易造成格式冲突和无关 diff。

## 决定

使用同目录的 `WORK_ITEMS.json` 保存任务状态元数据。`WORK_ITEMS.md` 只保存任务标题、目标、验收、检查、决策和备注等正文；`WORK_ITEMS.json` 保存 `version`、每个任务的 `status` 和 `completed`。流程工具读取两个文件并只更新 JSON。

现有按状态分节、带复选框的 Markdown 在 JSON 不存在时可以由流程工具首次使用时迁移。迁移不会覆盖或重写原 Markdown；后续状态变化以 JSON 为准。

## 备选方案

- 继续把状态放在 Markdown 的章节和复选框中：阅读直观，但状态更新需要重排任务块，容易产生文档冲突。
- 引入 SQLite：查询和更新能力更强，但个人练习不需要数据库生命周期，也会产生第二份任务数据。
- 把所有任务正文复制到 JSON：机器处理方便，但会产生两套正文来源，编辑时容易失步。

## 影响

状态更新的 diff 更小，流程工具和 Agent 可以直接编辑结构化文件；任务正文仍然能在 Markdown 中阅读、评审和提交。代价是任务需要同时维护两个文件，新增或归档任务时必须同步更新任务 ID；使用方应校验两边的任务集合，缺失或多余元数据应直接报错。

`WORK_ITEMS.json` 使用 `version: 1`，后续修改格式时应增加版本并保留已提交的旧格式。流程工具默认使用 `WORK_ITEMS.md` 和同目录的 `WORK_ITEMS.json`。

## 验证

通过 JSON 语法检查、任务正文与状态 ID 的人工核对，以及根项目相关检查验证。
