# 任务模板与文档骨架 WBS

需求：[`docs/requirements/2026-09-13-work-item-templates.md`](../requirements/2026-09-13-work-item-templates.md)
方案：[`docs/tech-designs/2026-09-13-work-item-templates.md`](../tech-designs/2026-09-13-work-item-templates.md)
方案评审：[`docs/reviews/2026-09-13-work-item-templates-design.md`](../reviews/2026-09-13-work-item-templates-design.md)

## Dependency Graph

```text
T-01 模板契约与骨架渲染
  -> T-02 创建事务与 API
      -> T-03 页面入口
          -> T-04 测试、文档和验收
```

## WBS

### T-01 模板目录与文档渲染

- 输入：W-015 需求、现有字段和目录边界。
- 输出：七种模板、阶段清单、六类文档槽位、slug/path 校验和骨架内容。
- 验收：模板数据稳定；所有类型生成六份 Markdown；非法 slug 被拒绝。
- 检查：模板单测和 `git diff --check`。

### T-02 任务创建事务与 HTTP API

- 输入：T-01、现有 board/parser/metadata API。
- 输出：ID 分配、清单追加、metadata 初始化、文档写入、冲突检测、回滚、模板查询和创建接口。
- 验收：成功返回 201；任一步失败原文件不变且无残留；旧接口不回归。
- 检查：workflow 单测和 HTTP 测试。

### T-03 看板创建入口

- 输入：T-02 API 契约和现有页面样式。
- 输出：模板选择表单、输入校验、提交状态、成功文档列表和错误提示。
- 验收：桌面/窄屏可用；成功刷新任务板；错误不会清空已有任务。
- 检查：`node --check workflow/web/app.js` 和浏览器冒烟。

### T-04 测试、流程文档与交付验收

- 输入：T-01 至 T-03。
- 输出：边界/兼容/回滚测试、流程文档和模板索引更新、W-015 验收证据。
- 验收：`go test ./...`、静态检查和手工创建流程通过；任务可进入 Review。
- 检查：dashboard `go test ./...`、`go vet ./...`、项目 `git diff --check`。

## Delivery Gate

- [x] 每项有输入、输出、依赖和客观验收。
- [x] 文件事务、页面交互和测试均有独立闭环。
- [x] 没有把“生成文档内容”误认为完成真实需求/设计工作。
