# W-016 代码评审

任务：W-016
结论：PASS
评审日期：2026-09-13

## 评审范围

- `ttl get` 默认 key/tag 模糊匹配和显式 value 搜索条件。
- `--value`、`-v`、`-val` 参数兼容性，以及 TUI/`pick` 既有搜索语义。
- 未命中、歧义、非交互和真实 bbolt CLI 行为。

## 评审结果

- application service 通过 `SearchOptions` 控制字段范围，保留 `FindResources` 的既有 key/value/tag 语义供 TUI/`pick` 使用。
- CLI 默认不匹配 value，显式 value 参数只扩展 value 条件，不改变 key/tag 匹配。
- `-val` 在根命令解析前归一化为 `--value`，避免 pflag 将多字符短写误解析为组合短参数。

## 证据

- `go test ./internal/client/app ./internal/client/cli`
- `go test ./...`
- `go vet ./...`
- `./scripts/regression.sh`
- `./scripts/cli-composability.sh`
- 临时 bbolt 配置真实 CLI 冒烟覆盖默认 key/tag、`--value`、`-v` 和 `-val`。
- `git diff --check`

## 结论

验收条件、兼容性、错误语义和测试证据均已核对，没有发现阻塞性问题。W-016 可以进入本地提交阶段。
