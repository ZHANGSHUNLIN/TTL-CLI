# Personal Engineering Skills

本目录保存本项目的个人工程 Skill。它们是给 Agent 和项目维护者使用的工作规则，不是 GitHub Issue、Pull Request 或 CI 平台配置。

`personal-workflow-dashboard` 是全局共享 Skill 的项目接入入口，实际实现位于 `/Users/v_zhangshun01/.codex/skills/personal-workflow-dashboard/`；当前项目只保留这条软链接，不保存管理台源码。

## Skills

| Skill | 使用场景 |
| --- | --- |
| `personal-requirement-analysis` | 把大功能需求整理成范围、规则、异常和验收标准 |
| `personal-tech-design` | 根据需求设计代码归属、接口、数据、生命周期和验证方案 |
| `personal-design-review` | 在编码前评审技术方案是否可落地 |
| `personal-task-breakdown` | 把已评审方案拆成可独立执行的 WBS |
| `personal-write-code` | 按单个 WBS 任务实现最小代码变更 |
| `personal-write-unit-test` | 为代码变更补充行为、边界和失败路径单测 |
| `personal-work-item` | 创建任务、推进 `WORK_ITEMS.json` 状态、维护 `WORK_ITEMS.md` 正文 |
| `personal-pre-review-checks` | 进入 `Review` 前选择并执行最小检查 |
| `personal-regression-testing` | 按变更范围选择单元、集成、CLI 黑盒或完整回归 |
| `personal-code-review` | 以人工审核为中心检查 diff、行为、风险和测试 |
| `personal-commit` | 审核通过后创建本地 commit 并完成任务状态 |
| `personal-test-reliability` | 设计或排查并发、临时目录、数据库和集成测试问题 |
| `personal-find-simplifications` | 找到可以删除、合并或降低复杂度的实现 |
| `personal-prose-standard` | 审核 README、注释、帮助文本和本地化文案 |
| `personal-decision-records` | 新增、更新和检查 `docs/decisions/` 决策记录 |
| `personal-workflow-dashboard` | 初始化当前项目并接入唯一的全局流程管理台，查看已完成任务和项目文档 |

## 推荐顺序

```text
personal-work-item
  -> personal-requirement-analysis
  -> personal-tech-design
  -> personal-design-review
  -> personal-task-breakdown
  -> personal-write-code
  -> personal-write-unit-test
  -> personal-pre-review-checks
  -> personal-regression-testing
  -> personal-code-review
  -> personal-commit
```

其他 Skill 按改动类型按需使用：

- 测试出现隔离、时序或资源问题时使用 `personal-test-reliability`。
- 发现重复实现或过度设计时使用 `personal-find-simplifications`。
- 修改文档、注释、CLI 帮助或本地化文本时使用 `personal-prose-standard`。

## 与 DSH 的关系

这些 Skill 借鉴 DSH 的工程实践，但保留个人项目的简化边界：

- 使用 `WORK_ITEMS.md` 保存当前任务正文，使用 `WORK_ITEMS.json` 代替 Issue、Project 和 PR 状态，使用 `WORK_ITEMS_ARCHIVE.md` 保存已完成任务历史。
- 使用 `docs/decisions/` 代替大型 Agent Notes 体系。
- 人工审核由任务拥有者本人完成。
- 不要求 stacked PR、加权审批、远程机器人或自动合并。
