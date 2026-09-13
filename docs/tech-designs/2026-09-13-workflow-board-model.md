# 工作流看板模型升级技术方案

日期：2026-09-13
任务：W-014
需求：`docs/requirements/2026-09-13-workflow-board-model.md`
状态：reviewed

## 当前行为

看板服务位于共享 Skill `/Users/v_zhangshun01/.codex/skills/personal-workflow-dashboard`，项目通过 `.agents/skills/personal-workflow-dashboard` 链接使用它。服务读取项目的 `WORK_ITEMS.md` 和 `WORK_ITEMS.json`，`workflow.Item` 只包含生命周期状态、完成标记、目标、验收、检查、决策、备注和评审记录。

Markdown 解析器按任务标题和固定字段生成 `Item`；JSON 只保存状态和评审记录。HTTP board API 返回解析后的 `Item`，网页只按状态分栏并渲染固定详情字段。因此当前系统没有任务类型、优先级、阶段、阶段清单、依赖或产物的事实来源。

## 方案归属

| 区域 | 负责组件 | 职责 |
| --- | --- | --- |
| 任务内容 | 项目 `WORK_ITEMS.md` | 保存类型、优先级、当前阶段、阶段清单、依赖和产物等可读元数据 |
| 任务状态 | 项目 `WORK_ITEMS.json` | 继续保存生命周期状态、完成标记和评审记录；保持现有写入接口 |
| 解析与校验 | `workflow/board.go` | 解析可选字段、提供兼容默认值、生成阶段进度信息 |
| API | `workflow/server.go` | 继续返回 board JSON，不新增写接口 |
| 页面 | `workflow/web/index.html`、`app.js`、`styles.css` | 显示类型、阶段时间线、产物、依赖和阻塞信息 |
| 流程文档 | 项目 `docs/`、`WORK_ITEMS.md` | 规定字段、阶段 ID、任务类型模板和状态门禁 |

## 数据契约

### Markdown 可选字段

在现有任务字段旁增加以下字段，支持中英文名称：

| 字段 | 解析类型 | 缺失默认值 |
| --- | --- | --- |
| `类型` / `Type` | string | `other` |
| `优先级` / `Priority` | string | 空 |
| `当前阶段` / `Stage` | string | 空 |
| `阶段清单` / `Stage Plan` | comma-separated string | 空数组 |
| `父任务` / `Parent` | string | 空 |
| `依赖` / `Dependencies` | comma-separated IDs | 空数组 |
| `产物` / `Artifacts` | raw Markdown/string | 空 |
| `阻塞原因` / `Blocker` | raw Markdown/string | 空 |
| `下一步` / `Next` | raw Markdown/string | 空 |

字段保存在 `workflow.Item`，由 board API 原样返回。阶段清单只保存有序 ID；阶段完成度由当前阶段在清单中的位置计算，不持久化第二套完成状态，避免与任务状态产生冲突。

### 阶段进度响应

`Item` 增加：

```go
Type         string   `json:"type,omitempty"`
Priority     string   `json:"priority,omitempty"`
Stage        string   `json:"stage,omitempty"`
StagePlan    []string `json:"stagePlan,omitempty"`
Parent       string   `json:"parent,omitempty"`
Dependencies []string `json:"dependencies,omitempty"`
Artifacts    string   `json:"artifacts,omitempty"`
Blocker      string   `json:"blocker,omitempty"`
Next         string   `json:"next,omitempty"`
```

页面根据 `StagePlan` 和 `Stage` 展示 `currentIndex/total`。若阶段清单为空、当前阶段不在清单中或存在重复项，API 仍成功返回，页面显示警告而不猜测修复。

### 兼容策略

- `WORK_ITEMS.json` metadata version 继续为 `1`；新字段只在 Markdown 中增加，不改变 JSON 必需字段。
- 旧任务缺少新字段时使用上述默认值，现有 board/archive/status/review API 响应保持兼容。
- 解析器遇到未知字段、未知类型或未知阶段只保留原任务并返回可展示的原始值，不因单项错误拒绝整个项目。
- 不新增自动迁移、页面写入或远程同步；任务拥有者仍通过项目文件维护元数据。

## 页面行为

- 卡片：ID、标题、类型 badge、优先级 badge、生命周期状态、当前阶段和阶段进度。
- 详情头部：生命周期状态、类型、优先级和阶段进度。
- 详情新增阶段时间线：阶段 ID 使用固定本地化标签；已完成、当前、未开始和跳过由阶段清单位置表达。
- 详情新增产物、父任务、依赖、阻塞原因和下一步；路径仍通过现有项目文档浏览器打开，不新增文件读取接口。
- 未配置元数据时显示“未分类”“阶段未设置”，而不是隐藏字段或推断值。
- 将 Doing 的页面文案从“正在实现”改为“处理中”，与流程文档一致。

## 生命周期与失败行为

- 解析失败只影响对应字段；任务标题、状态和既有五个详情字段仍可显示。
- 阶段元数据不参与状态移动校验；状态接口仍只校验五个生命周期状态。
- Review、评论、打回和 Done 归档流程保持现有行为。
- 产物链接不在服务端解析外部路径；页面只展示 Markdown 文本并沿用文档 API 的项目目录安全边界。

## 安全与兼容性

- 不执行 Markdown 中的命令、链接或外部路径；现有 Markdown 渲染继续转义 HTML。
- 不改变状态文件权限、原子写入、项目注册和归档行为。
- 共享 Skill 的改动会影响所有已注册项目；因此新字段全部可选，旧项目无需立即修改即可使用。

## 未选择的方案

- **把需求、设计、编码拆成多个看板列**：不同任务类型会有不同适用阶段，固定列会制造大量“伪进行中”任务，且破坏现有状态 API。
- **只在 JSON 增加字段**：状态文件不适合承载可读的任务计划和文档链接，且会让人工维护脱离 `WORK_ITEMS.md`。
- **自动根据标题推断类型和阶段**：推断不可审计，容易把候选或 Bug 修复误判为 feature。
- **引入数据库或远程项目管理服务**：超出本地个人工作流边界。

## 测试计划

| 层级 | 覆盖内容 |
| --- | --- |
| workflow 单测 | 新字段解析、默认值、列表分割、未知值、阶段不一致和旧 Markdown 兼容 |
| HTTP 单测 | board API 返回新增字段；旧 metadata 仍可读取；状态/评论接口无回归 |
| 页面冒烟 | 卡片 badge、阶段进度、详情时间线、长产物文本、缺失元数据和项目文档入口 |
| 回归 | dashboard 模块 `go test ./...`；项目侧 `git diff --check`；手工刷新、状态移动、评论和归档 |

## 发布与恢复

- 先编译并启动共享 dashboard 的新版本，在当前项目和一个无新字段的临时项目上读取 board。
- 若页面或解析出现问题，停止新进程并恢复旧 dashboard 进程；项目文件不需要迁移或回滚。
- 新字段逐个加入项目任务，未更新的任务继续显示兼容默认值。
