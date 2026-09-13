# Add Personal Engineering Skills

日期：2026-09-11
状态：adopted
任务：W-001

## 背景

项目已经使用本地 `WORK_ITEMS.md`、人工审核和 `docs/decisions/`，但缺少可以重复使用的工程操作规范。直接移植 DeepSeek Harness 的完整 GitHub、PR、stack 和审批流程会给个人练手项目带来不必要的复杂度。

## 决定

项目在 `.agents/skills/` 下维护轻量个人工程 Skill，覆盖任务状态、进入 Review 前检查、代码审核、测试可靠性、简化审查、工程文档和决策记录。Skill 只描述工作方法，不自动管理 GitHub 状态，也不替代任务拥有者的人工审核。

## 备选方案

- 直接复制 DeepSeek Harness 的全部 Skill：包含团队协作和 GitHub 流程，超出个人项目需要。
- 只保留一份开发流程文档：容易形成一份过长的清单，难以按场景复用。
- 不增加 Skill：继续依赖临时对话约定，容易遗漏检查步骤。

## 影响

项目获得了一套与 Go CLI 工程匹配的可复用操作规范。维护成本是新增 Skill 文档需要随命令、目录结构和人工审核规则一起更新；运行时仍由个人决定何时调用 Skill，不能把 Agent 的结论当成自动批准。

## 验证

通过 Markdown 结构检查、`git diff --check`、`go test ./...` 和 `go vet ./...` 验证。Skill 的命令和状态规则与 `AGENTS.md`、`WORK_ITEMS.md`、`WORK_ITEMS.json`、`docs/personal-development-workflow.md` 保持一致。
