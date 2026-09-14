# Work Item Template

复制到 `WORK_ITEMS.md` 的 `Tasks` 区域，并为任务选择合适的类型和阶段清单。

```md
- W-XXX 任务标题
  - 类型：feature | bugfix | refactor | docs | chore | spike | other
  - 优先级：P0 | P1 | P2 | P3
  - 当前阶段：requirements
  - 阶段清单：requirements,design,design_review,breakdown,implementation,tests,delivery_review,commit
  - 父任务：无
  - 依赖：无
  - 产物：需求=；方案=；评审=；WBS=；测试=；提交=
  - 阻塞原因：无
  - 下一步：
  - 目标：
  - 验收：
  - 检查：
  - 决策：无需记录，原因：
  - 备注：
```

`WORK_ITEMS.json` 只维护该任务的生命周期状态和完成标记。开始处理任务时将对应状态设为 `Doing`；实现和测试完成后才设为 `Review`。看板不保存评审评论或打回记录。
