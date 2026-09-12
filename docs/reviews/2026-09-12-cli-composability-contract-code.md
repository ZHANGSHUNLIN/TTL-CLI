# CLI 可组合性契约代码评审

日期：2026-09-12
任务：W-012
结论：PASS

## Findings

未发现需要阻塞进入人工 Review 的问题。评审中发现并已修正以下问题：

- JSON 成功结果此前会在存储关闭前写入 stdout；关闭失败可能同时出现成功结果和失败退出码。现改为缓冲 JSON，命令执行和存储关闭均成功后才提交 stdout。
- 新资源 Service 的英文错误曾直接进入默认文本模式。现由 CLI 适配回既有 i18n 文案，同时保留 typed error，供 `--non-interactive` 使用分类退出码。
- history/audit debug 信息曾可能污染 JSON stdout。现 JSON 模式抑制这些非契约诊断，并由黑盒测试校验成功和错误各自只有一个 JSON 文档。
- `UpdateResourceValue` 曾错误调用 `SaveResource`。现恢复使用 `UpdateResource`，保持已有存储及远端 PUT 语义，并由单元测试锁定。
- 默认文本错误和 debug 输出流曾被统一改为 stderr。现默认模式继续保持历史 stdout 语义；只有机器模式执行错误写 stderr。
- 机器 flag 的预检测此前只识别裸 flag。现支持 `--json=true/false`，并将非法布尔值稳定映射为 `invalid_argument`。

## Acceptance Coverage

- [x] `add/get/update/del/tag/dtag` 输出版本 1 JSON envelope。
- [x] JSON 成功只写 stdout；失败只写单个 stderr JSON，stdout 为空。
- [x] 退出码 1/2/3/4 分别覆盖系统、参数、不存在和冲突/歧义。
- [x] `--json` 隐含非交互，`--non-interactive` 不等待数字选择。
- [x] `add/update <key> -` 保留多行、末尾换行及空输入。
- [x] 未覆盖命令拒绝机器模式。
- [x] 默认文本模式的核心黑盒回归及 i18n 错误语义保持兼容。
- [x] TUI、服务端和持久化格式不依赖 CLI DTO。

## Evidence

- `go test ./internal/client/app ./internal/client/cli`
- `./scripts/cli-composability.sh`
- `./scripts/regression.sh`
- `go test ./...`
- `go test ./integration_test/...`
- `go test -race ./...`
- `go vet ./...`
- `bash -n scripts/cli-composability.sh`
- `gofmt -s -l .`（无输出）
- `git diff --check`
- `./scripts/verify.sh`

## Scope And Remaining Risk

旧 `command` 包的全局核心 handler 暂时保留；客户端 root 已使用 `internal/client/cli` 的新 handler。本任务不扩大到 W-006 的旧入口清理，也不改服务端、TUI、配置或存储格式。主要剩余风险是 JSON schema 发布后的长期兼容维护，已由 ADR 记录。
