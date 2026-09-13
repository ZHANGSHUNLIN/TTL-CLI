# 测试计划与结果：修复 TUI 打开 Markdown 链接失败

任务：W-018
日期：2026-09-13
状态：passed

## 测试范围

- `internal/client/opener`：Markdown 链接目标提取、纯值兼容、空值和不完整标记错误。
- `internal/client/tui`：打开成功请求退出，打开失败保留详情页。
- `internal/client/cli`：`open` 路径编译并复用目标解析。
- CLI 黑盒：既有资源生命周期和打开命令回归不受影响。

## 计划命令

```text
go test ./internal/client/opener ./internal/client/tui ./internal/client/cli
go test ./...
go vet ./...
gofmt -s -l .
./scripts/regression.sh
git diff --check
```

## 结果

- [x] `gofmt -w` 后 `git diff --check` 通过。
- [x] `go test ./internal/client/opener ./internal/client/tui ./internal/client/cli` 通过。
- [x] `go test ./...` 通过。
- [x] `go test -race ./internal/client/opener ./internal/client/tui ./internal/client/cli` 通过。
- [x] `go vet ./...` 通过。
- [x] `go build -o /tmp/ttl-w018 ./cmd/ttl` 和 `go build -o /tmp/ttl-server-w018 ./cmd/ttl-server` 通过。
- [x] `./scripts/regression.sh` 通过。
- [x] macOS 临时 `open` 可执行文件捕获到的参数为 `https://example.com/path`，未启动真实浏览器。
