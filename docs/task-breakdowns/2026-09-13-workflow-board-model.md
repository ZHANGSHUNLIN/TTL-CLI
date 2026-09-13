# 工作流看板模型升级任务分解

日期：2026-09-13
任务：W-014
需求：`docs/requirements/2026-09-13-workflow-board-model.md`
方案：`docs/tech-designs/2026-09-13-workflow-board-model.md`
方案评审：`docs/reviews/2026-09-13-workflow-board-model-design.md`

## Dependency Graph

```text
T-01 解析模型与兼容默认值
  -> T-02 board API 与页面阶段展示
      -> T-03 测试、流程文档和示例任务
```

## WBS

### T-01 扩展任务元数据解析模型

- 目标：让 board API 能读取任务类型、优先级、阶段、阶段清单、父子关系、依赖、产物、阻塞原因和下一步。
- 输入：现有 `workflow/board.go`、需求契约和旧项目样例。
- 输出：可选字段解析、兼容默认值和稳定 JSON 字段；旧任务继续正常加载。
- 依赖：无。
- 负责区域：共享 dashboard `workflow/board.go` 及其单测。
- 验收：
  - [ ] 中英文字段可解析，列表字段保留顺序。
  - [ ] 缺失、未知、重复和不一致值不会阻塞整板。
  - [ ] 旧 Markdown/JSON 的 board/archive/status/review 行为不变。
- 检查：`go test ./workflow -run 'Test(Parse|Load|Build)'` 或等价的 dashboard 单测。
- 风险或后续：阶段进度只由当前阶段和阶段清单计算，不持久化第二套完成状态。

### T-02 展示任务类型和阶段进度

- 目标：卡片和详情页展示类型、优先级、阶段进度、阶段时间线、产物、依赖、阻塞原因和下一步。
- 输入：T-01 的 board API 字段、现有网页渲染和项目文档 API。
- 输出：兼容旧任务的 UI；缺失元数据显示明确默认状态；Doing 文案与流程一致。
- 依赖：T-01。
- 负责区域：共享 dashboard `workflow/web/index.html`、`app.js`、`styles.css`。
- 验收：
  - [ ] 卡片显示类型、优先级、当前阶段和 `current/total` 进度。
  - [ ] 详情时间线正确处理跳过阶段、未知阶段和空阶段清单。
  - [ ] 产物和依赖信息可读，项目文档入口和状态操作无回归。
- 检查：dashboard HTTP 测试、浏览器手工冒烟和响应 JSON 人工核对。
- 风险或后续：共享 Skill 改动会影响其他项目，必须用无新字段临时项目验证。

### T-03 测试、流程文档和项目示例

- 目标：补齐解析/API 测试，更新流程文档、模板、索引和 W-014 示例，使字段来源和生命周期门禁一致。
- 输入：T-01/T-02 实现和需求验收标准。
- 输出：覆盖正常、边界、旧数据兼容和页面关键状态的验证证据；项目文档可指导新任务创建。
- 依赖：T-02。
- 负责区域：dashboard `workflow/*_test.go`、项目 `docs/`、`WORK_ITEMS.md`、`WORK_ITEMS.json`。
- 验收：
  - [ ] 单测和 HTTP 测试覆盖新增字段及异常场景。
  - [ ] 模板说明任务类型、阶段清单和轻重流程选择。
  - [ ] W-014 链接需求、方案、评审、WBS，项目任务状态仍符合生命周期规则。
- 检查：dashboard `go test ./...`、项目 `git diff --check`、页面手工验收。
- 风险或后续：不在本任务中增加页面编辑元数据或自动推断任务类型。

## Delivery Gate

- [x] 每个任务都有输入、输出、依赖和验收标准。
- [x] 每个任务可以独立实现或有明确前置依赖。
- [x] 没有按文件名机械拆分；每项对应一个可观察闭环。
- [x] 测试和文档任务已纳入。
