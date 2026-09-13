# Add Complete Personal Development Flow

日期：2026-09-12
状态：adopted
任务：W-009

## 背景

项目已有任务清单、回归测试、人工代码审核和本地提交规则，但需求分析、技术设计、方案评审、WBS 拆分、编码和单测没有独立的 Skill 或固定产物。大功能容易直接进入编码，导致边界、依赖和验收条件在实现过程中反复变化。

## 决定

项目采用完整的八阶段个人研发流程：需求分析、技术方案、方案评审、WBS 拆分、逐项编码、单元测试、代码评审和本地提交。`WORK_ITEMS.md` 保存任务正文，`WORK_ITEMS.json` 管理阶段状态；需求、方案、评审和任务拆分分别归档到 `docs/requirements/`、`docs/tech-designs/`、`docs/reviews/` 和 `docs/task-breakdowns/`；通用格式归档到 `docs/templates/`。小功能可以合并阶段，但不能省略与变更风险匹配的验收、检查和人工审核。

## 备选方案

- 继续只使用一份流程说明：维护成本低，但阶段输入、输出和放行条件不清晰。
- 原样移植 DSH 的企业流程和远程审批：与个人项目的本地任务清单目标冲突，且引入不必要的平台依赖。
- 只增加编码和测试 Skill：能规范实现，但无法提前发现需求边界、接口归属和任务依赖问题。

## 影响

大功能拥有可追踪的阶段产物和放行条件，Agent 可以按单个 WBS 任务执行，人工审核也有明确依据。代价是多维护几类短文档；通过“小功能合并阶段”和模板化控制额外负担。流程不建立 GitHub Issue、Pull Request、审批分数、远程 CI 或自动合并。

## 验证

新增七个个人版阶段 Skill，补充五类文档模板，更新 `AGENTS.md`、工程 Skill 索引和个人工作流，并在 `WORK_ITEMS.md`/`WORK_ITEMS.json` 登记 W-009。通过 `git diff --check`、Skill frontmatter 检查和人工检查链接、文件职责与流程一致性验证。
