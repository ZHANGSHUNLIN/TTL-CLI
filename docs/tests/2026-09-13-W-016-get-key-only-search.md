# 测试计划与结果：让 ttl get 默认按 key 和 tag 模糊查询

任务：W-016
类型：bugfix
状态：通过，待代码评审

## 测试结果

- [x] `go test ./internal/client/app ./internal/client/cli`
- [x] service 测试证明默认 key/tag 不命中 value-only、能命中 tag-only，开启 `IncludeValue` 后命中 value。
- [x] CLI 测试证明 `ttl get` 默认不命中 value、可以命中 tag，`ttl get --value` 可以命中 value。
- [x] `go test ./...`
- [x] `go vet ./...`
- [x] `gofmt -s -l`（变更 Go 文件无输出）
- [x] 五种 locale JSON 可解析。
- [x] `git diff --check`
- [x] `./scripts/regression.sh`
- [x] `./scripts/cli-composability.sh`（覆盖 tag 查询下的歧义用例）。
- [x] 临时 bbolt 配置真实 CLI 冒烟：默认 `get secret` 不命中 value，`get --value secret`、`get -v secret` 和 `get -val secret` 均命中。
- [x] `scripts/regression.sh` 覆盖默认 key/tag 语义及三个 value 搜索入口。

## 验收状态

自动检查、真实 CLI 冒烟和 owner 验收均已完成，当前待代码评审。

## 验收与证据

实现文件：`internal/client/app/resources.go`、`internal/client/cli/resource_commands.go`、`internal/client/app/resources_test.go`、`internal/client/cli/get_test.go`。

## 备注

这是由任务模板生成的文档骨架，不代表本阶段已经完成。
