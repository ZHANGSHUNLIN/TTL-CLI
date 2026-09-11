# Use Local Work Items For Personal Practice

日期：2026-09-11
状态：adopted
任务：W-000

## 背景

这个项目用于个人练习。完整的 GitHub Issue、Project、Pull Request、reviewer 请求和审批积分适合多人协作，但对个人练习来说维护成本偏高，也会把注意力从代码、测试和复盘转移到平台配置上。

## 决定

使用仓库内的 `WORK_ITEMS.md` 作为项目任务状态记录。任务只保留 `Inbox`、`Doing`、`Review`、`Blocked` 和 `Done` 五个状态。人工审核由任务拥有者本人完成，审核重点是 diff、验证证据和是否需要决策记录。

重要设计取舍记录在 `docs/decisions/`。任务状态和设计决策分开维护：`WORK_ITEMS.md` 记录“现在做到哪一步”，决策记录说明“为什么这样做”。

## 备选方案

- 继续使用 GitHub Issue 和 Project：适合多人协作，但个人练习需要维护远程状态、标签和自动化规则。
- 只使用 commit message：提交信息能说明改了什么，但不能稳定记录备选方案和设计代价。
- 把决策写在任务备注里：短期方便，但任务完成后难以检索，也容易把长期设计理由埋在任务流转细节中。

## 影响

这个流程更轻，离线可用，也更适合反复练习。代价是没有远程看板、自动审批和多人 reviewer 保护；需要任务拥有者在 `Review` 阶段主动检查 diff、运行相关命令，并判断是否补充决策记录。

## 验证

`WORK_ITEMS.md` 定义任务状态和任务模板。`docs/personal-development-workflow.md` 定义日常流程、自动检查、人工审核和决策记录检查。`docs/decisions/README.md` 定义何时写决策记录以及审核标准。

