# 代码评审：修复 TUI 打开 Markdown 链接失败

任务：W-018
日期：2026-09-13
结论：PASS

## Findings

未发现阻塞或条件性问题。

- `internal/client/opener.Target` 只做输入归一化，不修改资源数据；Markdown 链接目标和纯值边界有 focused tests。
- TUI 和 CLI 的平台打开器均复用同一 helper，避免两个入口对同一资源值行为不一致。
- macOS/Windows 的外部进程参数保持单参数传递，Linux 和其他系统的既有限制不变。
- TUI 打开成功仍返回 `tea.Quit`，打开失败仍停留在详情页并展示错误。
- 变更未触及存储、API、配置、加密、同步或资源生命周期。

## Evidence

- `go test ./internal/client/opener ./internal/client/tui ./internal/client/cli`
- `go test ./...`
- `go test -race ./internal/client/opener ./internal/client/tui ./internal/client/cli`
- `go vet ./...`
- `go build -o /tmp/ttl-w018 ./cmd/ttl`
- `go build -o /tmp/ttl-server-w018 ./cmd/ttl-server`
- `./scripts/regression.sh`
- `git diff --check`

## Decision

PASS。满足 W-018 验收条件，可以创建本地提交。
