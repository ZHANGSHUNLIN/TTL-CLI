# Personal Engineering Skills

本项目把 DSH 中适合个人使用的工程能力复刻为本地 Skill。Skill 只提供做事方法，不会自动修改 GitHub 状态，也不会替代人工审核。

## 文档语言规范

项目 README、`docs/`、`WORK_ITEMS*` 和所有项目 Skill 文档统一使用中文。Skill 名称、代码标识符、命令、路径、协议字段、文件格式和外部工具名称保留原文；英文引用必须配中文解释。新增文档不得使用英文正文，已有英文内容在被任务触及时优先翻译。

## 能力分层

```text
WORK_ITEMS.md + WORK_ITEMS.json
  -> personal-work-item
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

docs/decisions/
  -> personal-decision-records

按需使用：
  -> personal-test-reliability
  -> personal-find-simplifications
  -> personal-prose-standard
```

## 一次大任务怎么用

1. 使用 `personal-work-item` 在 `WORK_ITEMS.md` 创建任务正文，并在 `WORK_ITEMS.json` 中将任务状态设为 `Doing`。
2. 使用 `personal-requirement-analysis` 产出需求文档。
3. 使用 `personal-tech-design` 产出技术方案；需要时使用 `personal-decision-records` 记录重要取舍。
4. 使用 `personal-design-review` 审核方案，只有 `PASS` 或已解决的 `CONDITIONAL` 才能继续。
5. 使用 `personal-task-breakdown` 产出带依赖和验收标准的 WBS。
6. 对每个 `T-XX` 使用 `personal-write-code`，再使用 `personal-write-unit-test` 补测试。
7. 使用 `personal-pre-review-checks` 和 `personal-regression-testing` 选择并执行检查。
8. 使用 `personal-code-review` 做人工审核，处理所有 `BLOCK` 问题。
9. 使用 `personal-commit` 创建本地 commit，再在 `WORK_ITEMS.json` 中将任务状态设为 `Done`。

## 与 DSH 能力的对应关系

| DSH 能力 | 本项目个人版 |
| --- | --- |
| `dev-analyze-requirement` | `personal-requirement-analysis` |
| `dev-design-tech-solution` | `personal-tech-design` |
| `dev-review-tech-design` | `personal-design-review` |
| `dev-breakdown-tasks` | `personal-task-breakdown` |
| `dev-write-code` | `personal-write-code` |
| `dev-write-unit-test` | `personal-write-unit-test` |
| `dsh-pre-push-checks` | `personal-pre-review-checks` |
| `dsh-code-review` | `personal-code-review` |
| `dsh-ci-test-reliability` | `personal-test-reliability` |
| DSH layered regression checks | `personal-regression-testing` + `docs/regression-testing.md` |
| `dev-commit-code` | `personal-commit`，只创建本地 commit，不建立 GitHub PR |
| `dsh-find-simplifications` | `personal-find-simplifications` |
| `dsh-prose-standard` | `personal-prose-standard` |
| Agent Notes | `docs/decisions/` + `personal-decision-records` |
| Issue/Project 状态 | `WORK_ITEMS.json` |

Sub-agent、Workflow 和 Agent Teams 不在这次个人 Skill 包中自动启用。它们是运行时协作能力，不是文档规则；个人项目先使用单个 Agent 加人工审核，确实需要并行委派时再单独引入。
