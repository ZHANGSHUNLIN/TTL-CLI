# 交付验收：让 ttl get 默认按 key 和 tag 模糊查询

任务：W-016
类型：bugfix
状态：验收和代码评审完成，待本地提交

## 验收结果

- [x] 默认 `ttl get <query>` 模糊匹配 key 和 tag，不匹配 value。
- [x] `ttl get --value <query>` 显式扩展到 value 匹配，同时保留 key/tag 匹配。
- [x] `ttl get -v <query>` 和 `ttl get -val <query>` 可作为 value 搜索简写。
- [x] TUI 和 `pick` 继续使用既有 key/value/tag 搜索。
- [x] 自动测试、静态检查和 locale 校验通过。
- [x] owner 已完成真实 CLI 使用确认（默认 key/tag 搜索、value 搜索显式开关及 `-v`/`-val` 简写，2026-09-13 确认）。
- [x] 代码评审通过；本地 commit 待创建。

## 验收结论

自动检查、owner 真实 CLI 验收和代码评审均通过；尚未创建本地 commit。

## 验收与证据

- 实现：`internal/client/app/resources.go`、`internal/client/cli/resource_commands.go`
- 测试：`internal/client/app/resources_test.go`、`internal/client/cli/get_test.go`
- 检查：`go test ./...`、`go vet ./...`、`./scripts/regression.sh`、`./scripts/cli-composability.sh`、`git diff --check`

## 备注

这是由任务模板生成的文档骨架，不代表本阶段已经完成。
