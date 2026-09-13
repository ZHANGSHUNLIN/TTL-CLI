# 任务模板与文档骨架需求

日期：2026-09-13
任务：W-015
状态：reviewed
基线：[`docs/personal-development-workflow.md`](../personal-development-workflow.md)

## 问题与目标

当前创建任务依赖手工编辑 `WORK_ITEMS.md`，容易只写一行标题，遗漏需求分析、技术设计、方案评审、WBS、测试和验收材料。看板已经能展示任务类型和阶段，但没有把“登记任务”和“建立研发证据”连接起来。

目标是提供可选择的任务模板。用户填写任务标题、类型、优先级、目标等最小信息后，系统一次性创建任务记录和一组文档骨架，让任务从 Inbox 开始就具备可追踪的研发产物槽位。

## 范围

### In Scope

- 提供 `feature`、`bugfix`、`refactor`、`docs`、`chore`、`spike`、`other` 七种模板。
- 模板定义类型标签、默认阶段清单、文档槽位和每份文档的标题/章节提示。
- 看板提供模板列表和创建任务表单；服务端提供模板查询和任务创建 API。
- 创建时分配未被 `WORK_ITEMS.md` 或归档占用的下一个 `W-数字` ID。
- 在项目 `docs/requirements`、`docs/tech-designs`、`docs/reviews`、`docs/task-breakdowns`、`docs/tests`、`docs/acceptance` 下生成六份 Markdown 骨架，并把链接写入任务的“产物”字段。
- 任务正文写入 `WORK_ITEMS.md`，状态元数据写入 `WORK_ITEMS.json`，初始状态为 `Inbox`、完成标记为 `false`。
- 对未知模板、空标题、非法 slug、文档路径冲突和已有任务 ID 冲突返回明确错误；失败时不留下任何新文件或部分任务。
- 旧项目和没有新字段的旧任务继续可读；不自动修改既有任务。

### Out Of Scope

- 自动填写真实需求、设计、代码或测试内容。
- 自动移动生命周期状态、自动开始阶段、自动提交 Git 或发起远程评审。
- 页面内编辑已生成 Markdown；继续使用项目文档浏览器查看，文件由开发者编辑。
- 自定义模板持久化、模板版本迁移和跨项目模板同步。

## 模板契约

所有模板都生成相同的六类文档槽位，具体章节根据类型给出提示。`spike` 和 `other` 可以在文档中标记“不适用”，但不能省略槽位。

| 模板 | 默认阶段清单 | 重点文档提示 |
| --- | --- | --- |
| feature | requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit | 用户价值、业务规则、接口/数据、验收场景 |
| bugfix | requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit | 现象、影响、根因、修复边界、回归证据 |
| refactor | requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit | 行为不变、影响范围、兼容性和性能 |
| docs | requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit | 受众、信息架构、链接/示例校验 |
| chore | requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit | 运维/依赖影响、回滚和验证 |
| spike | requirements,design,design_review,breakdown,tests,delivery_review,commit | 问题、调查方法、结论；实现章节标记不适用 |
| other | requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit | 创建后必须由负责人补充范围和适用阶段说明 |

## 创建输入

```json
{
  "template": "feature",
  "title": "支持资源批量导入",
  "priority": "P1",
  "goal": "让用户可以一次导入多个本地资源",
  "acceptance": "导入成功、部分失败和空输入均有可验证结果",
  "slug": "bulk-import",
  "parent": "",
  "dependencies": []
}
```

`template`、`title`、`goal` 为必填；`priority` 缺省为 `P2`；`slug` 缺省由标题生成，仅允许 ASCII 小写字母、数字和连字符，长度 1-80。客户端展示模板默认值，服务端仍重新校验。

## 原子性与冲突规则

1. 服务端先读取并校验项目文件、模板和所有目标路径，再生成内存中的 Markdown、JSON 和文档内容。
2. 目标文档已存在、标题对应的 slug 冲突或 ID 已占用时，整个请求失败，响应 `409`，不修改任务或文档。
3. 写入使用同目录临时文件和 rename；任务 Markdown、JSON 和文档全部成功后才返回 `201`。中途失败必须删除本次创建的文件并恢复原始两份清单。
4. 创建过程由 dashboard 进程的互斥锁保护；不承诺跨进程并发编辑的强一致性。

## 验收标准

- [ ] `GET /api/projects/{id}/templates` 返回七种模板及阶段/文档定义。
- [ ] `POST /api/projects/{id}/items` 能创建任务、六份文档和初始 JSON 状态，并返回任务与文档路径。
- [ ] 每种模板的类型、阶段清单和文档提示正确；`spike`/`other` 没有伪造实现结论。
- [ ] 空标题、未知模板、非法 slug、路径穿越、重复路径和重复任务均被拒绝且无残留。
- [ ] 写入失败可回滚，原始 `WORK_ITEMS.md`、`WORK_ITEMS.json` 和已有文档内容不变。
- [ ] 旧项目、归档和现有状态/评审接口保持兼容。
- [ ] 页面能完成选择模板、填写最小输入、提交、显示成功文档列表和错误提示。
- [ ] 单元、HTTP、静态脚本和浏览器冒烟证据齐全。
