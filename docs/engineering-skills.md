# Personal Engineering Skills

本项目把 DSH 中适合个人使用的工程能力简化为本地 Skill。Skill 只提供做事方法，不会自动修改 GitHub 状态，也不会替代人工审核。

## 能力分层

```text
WORK_ITEMS.md
  -> personal-work-item
  -> personal-pre-review-checks
  -> personal-regression-testing
  -> personal-code-review
  -> local commit

docs/decisions/
  -> personal-decision-records

按需使用：
  -> personal-test-reliability
  -> personal-find-simplifications
  -> personal-prose-standard
```

## 一次任务怎么用

1. 使用 `personal-work-item` 在 `WORK_ITEMS.md` 创建任务并移动到 `Doing`。
2. Agent 实现代码，必要时使用会话 todo 拆分临时步骤。
3. 重要设计取舍使用 `personal-decision-records` 创建记录。
4. 实现完成后移动到 `Review`。
5. 使用 `personal-pre-review-checks` 选择最小检查集合。
6. 使用 `personal-code-review` 做人工审核，处理所有 `BLOCK` 问题。
7. 审核通过后创建本地 commit，再移动到 `Done`。

## 与 DSH 能力的对应关系

| DSH 能力 | 本项目个人版 |
| --- | --- |
| `dsh-pre-push-checks` | `personal-pre-review-checks` |
| `dsh-code-review` | `personal-code-review` |
| `dsh-ci-test-reliability` | `personal-test-reliability` |
| DSH layered regression checks | `personal-regression-testing` + `docs/regression-testing.md` |
| `dsh-find-simplifications` | `personal-find-simplifications` |
| `dsh-prose-standard` | `personal-prose-standard` |
| Agent Notes | `docs/decisions/` + `personal-decision-records` |
| Issue/Project 状态 | `WORK_ITEMS.md` |

Sub-agent、Workflow 和 Agent Teams 不在这次个人 Skill 包中自动启用。它们是运行时协作能力，不是文档规则；个人项目先使用单个 Agent 加人工审核，确实需要并行委派时再单独引入。
