# 任务模板与文档骨架技术方案

日期：2026-09-13
任务：W-015
需求：[`docs/requirements/2026-09-13-work-item-templates.md`](../requirements/2026-09-13-work-item-templates.md)
状态：reviewed

## 设计摘要

在共享 dashboard `workflow` 包中增加只读模板目录和 `CreateWorkItem` 服务函数。HTTP 层增加模板查询与任务创建路由，网页在当前项目工具区增加创建模态框。创建函数在进程互斥锁内完成校验、临时写入、提交和失败清理。

## 组件归属

| 组件 | 变更 |
| --- | --- |
| `workflow/templates.go` | 模板定义、文档槽位、slug 校验和骨架渲染 |
| `workflow/create.go` | 读取清单、分配 ID、事务式文件提交和回滚 |
| `workflow/server.go` | `GET templates` 与 `POST items` 路由、请求校验和响应 |
| `workflow/web/index.html`/`app.js`/`styles.css` | 创建按钮、表单、成功/错误状态 |
| `workflow/*_test.go` | 模板契约、创建、冲突、回滚和 HTTP 覆盖 |

## 数据与文件契约

任务 Markdown 继续是可读事实来源，新增任务沿用现有字段格式。JSON metadata 只新增 `items[ID] = {status: "Inbox", completed: false}`，版本保持 1。每次创建的文档基于日期、ID 和 slug 命名：

```text
docs/requirements/YYYY-MM-DD-W-XXX-<slug>.md
docs/tech-designs/YYYY-MM-DD-W-XXX-<slug>.md
docs/reviews/YYYY-MM-DD-W-XXX-<slug>-design.md
docs/task-breakdowns/YYYY-MM-DD-W-XXX-<slug>.md
docs/tests/YYYY-MM-DD-W-XXX-<slug>.md
docs/acceptance/YYYY-MM-DD-W-XXX-<slug>.md
```

请求可显式传入 slug，但不能传目录或扩展名。所有目标路径用 `filepath.Join(projectRoot, relative)` 后以 `isWithin` 校验；文档目录按需创建。

## 创建事务

1. 加锁并读取 `WORK_ITEMS.md`、`WORK_ITEMS.json`，解析所有任务和归档 ID。
2. 校验模板、字段、slug、目标路径不存在，并计算最大数字 ID 加一。
3. 生成任务 Markdown 块、六份文档内容和新 metadata；保留原文件权限。
4. 将清单和文档写入临时文件。先 rename 文档，最后 rename `WORK_ITEMS.md` 与 `WORK_ITEMS.json`；记录本次已替换/创建的路径。
5. 任意步骤失败时删除新文档、恢复被替换的清单备份，返回错误；成功返回创建结果。

由于普通文件 rename 不能跨文件提供真正的系统事务，回滚通过同目录备份和受控清理实现；每一步都在进程锁内执行，避免 dashboard 自身并发请求交错。

## API

- `GET /api/projects/{project}/templates` -> `{templates:[{id,label,description,stagePlan,documents}]}`
- `POST /api/projects/{project}/items` -> `201 {item, documents}`
- 错误：`400` 输入/模板/path，`404` 项目，`409` 冲突，`500` 文件系统失败。

## 页面

工具栏增加“新建任务”按钮。表单包含模板、标题、优先级、slug、目标、验收、父任务和依赖；提交期间禁用按钮，成功后刷新 board/archive 并列出六个文档路径，失败显示服务端错误。无需在页面编辑文档正文。

## 风险与验证

- 文件系统失败和权限问题通过注入/临时目录 HTTP 测试覆盖。
- 旧项目缺少 `docs/tests` 或 `docs/acceptance` 时由创建流程建立目录，不影响旧读取。
- 浏览器冒烟验证窄屏表单、长标题、错误提示和成功刷新。
