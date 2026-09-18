# 交付验收：移除当前工程的后端实现

任务：W-021
状态：ready-for-owner-review

## 验收结果

- [x] `cmd/ttl-server` 和 `internal/server` 已删除。
- [x] 服务端专属集成测试、冒烟与升级回滚脚本已删除。
- [x] CI 和 release 只构建、打包客户端。
- [x] `scripts/regression.sh` 与 `scripts/verify.sh` 不再构建或启动后端。
- [x] cloud 配置、`internal/client/remote` 和 `ttl sync` 保留。
- [x] 架构测试锁定后端目录不得重新出现。
- [x] README、AGENTS、项目概览、数据流和工程边界已改为“客户端 + 外部后端”。
- [x] `./scripts/verify.sh` 全部通过。
- [x] owner 已确认 diff，并要求创建 commit 后推送。

## 遗留风险

真实独立后端尚未纳入本仓验证，当前只能证明客户端 mock HTTP 契约。后续接入独立后端时需要补充跨工程契约或端到端检查。

## 回滚

变更只删除 Git 跟踪的仓库内实现，不触碰已部署服务或数据。提交前可直接修改当前 diff；提交后可通过新的回滚提交恢复历史实现。
