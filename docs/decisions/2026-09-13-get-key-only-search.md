# ttl get 默认按 key 搜索

日期：2026-09-13
状态：adopted
任务：W-016

## Context

`ttl get` 直接复用 key/value/tag 全字段搜索，导致资源内容中的普通词也会改变候选结果。CLI get 的常见意图是按资源 key 取值。

## Decision

`ttl get` 默认模糊匹配 key 和 tag，但不匹配 value；需要把 value 纳入搜索时显式使用 `--value`。TUI 和 `pick` 保持原有全字段搜索，不因 CLI 默认变化而改变共享 service 的既有语义。

## Alternative

继续让所有入口默认搜索 key/value/tag，优点是命令少一个 flag，但会保留用户反馈的非预期候选，且无法区分“找资源”与“找内容”两种意图。

## Verification

service 和 CLI 测试分别覆盖 key/tag 默认搜索与显式 value；`go test ./...`、`go vet ./...` 和 locale JSON 校验通过。
