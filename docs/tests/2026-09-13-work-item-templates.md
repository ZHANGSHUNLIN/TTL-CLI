# W-015 测试计划与结果

任务：W-015
状态：测试通过，待交付评审

## 已覆盖

- 模板目录包含七种类型，每种包含阶段清单和六个文档槽位。
- 创建成功会分配下一个 `W-XXX`、写入 Inbox metadata、追加 Markdown 任务并生成六份文档。
- 文档冲突、非法 slug 和缺失旧 metadata 的场景不会修改原清单或已有文档。
- 模板查询和创建 HTTP API 返回预期状态码与 JSON 结构。

## 检查结果

- [x] dashboard `go test ./...`
- [x] dashboard `go vet ./...`
- [x] `node --check workflow/web/app.js`
- [x] 项目 `git diff --check`、`jq empty WORK_ITEMS.json`
- [x] `ensure-project.sh` 运行态健康检查和模板接口查询
- [x] 浏览器端新建任务完整交互和窄屏截图冒烟（owner 于 2026-09-13 手工确认通过）

## 后续

自动检查和 owner 手工验收均已完成，进入 `delivery_review`。
