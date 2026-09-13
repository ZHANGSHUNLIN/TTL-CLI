# 客户端与服务端代码结构收敛方案评审

日期：2026-09-13  
任务：W-017  
方案：[`docs/tech-designs/2026-09-13-W-017-code-structure-convergence.md`](../tech-designs/2026-09-13-W-017-code-structure-convergence.md)  
结论：PASS（负责人已确认目标归属和一次性切换规则）

## Findings

- `[RESOLVED]` 评审结论：负责人已确认目标包归属
  - 位置：技术方案“Proposed Ownership”、需求“Open Questions And Assumptions”。
  - 问题：方案给出了 `internal/config`、`internal/crypto`、`internal/i18n` 和 `internal/core/text` 的建议归属，但这些移动会影响客户端、服务端、存储适配器和嵌入资源。
  - 后果：若负责人不同意其中任一归属，一次性切换的目标架构就不明确，不能开始编码。
  - 修复：已确认归属表；实现中因 bbolt 的共享依赖将配置和加密基础包收敛到 `internal/config`、`internal/crypto`，并同步方案和 WBS。
  - 验证：方案中的 ownership 表、内部执行顺序、WBS T-01/T-03/T-04 和架构测试规则保持一致。

- `[RESOLVED]` 旧测试兼容层的删除条件必须逐项核对
  - 位置：技术方案“Compatibility aliases”和“Rollout And Recovery”。
  - 问题：当前集成测试和服务端测试仍直接使用 `db.Stor`，不能只迁移生产 import 就删除 `db/`。
  - 后果：过早删除会导致测试覆盖断裂，或者迫使测试重新引入全局状态。
  - 修复：测试迁移已纳入本次整体切换的 T-05；当前仓库已无 `ttl-cli/db` consumer，旧包已删除。
  - 验证：`rg 'ttl-cli/db|db\.Stor|db\.InitDB' --glob '*.go'` 无结果，`go test ./...` 和集成测试通过。

- `[RESOLVED]` 文档状态需要在方案评审通过后统一升级
  - 位置：`PROJECT_OVERVIEW.md`、`docs/client-server-separation-plan.md` 和 W-017 文档状态。
  - 问题：现有文档同时描述阶段 A/B 已完成、阶段 C 部分未完成和旧包待删除；若只更新代码不更新状态，维护者仍会把目标结构误认为现状。
  - 后果：后续任务会重复盘点或误用旧入口。
  - 修复：本次切换已统一更新当前结构、迁移状态和实现前证据；最终完成证据待全量回归和代码评审后补齐。
  - 验证：文档中的目录、依赖和阶段勾选与 `go list ./...`、架构测试和 W-017 状态一致。

## Acceptance Coverage

- [x] 需求覆盖当前新旧结构、目标边界、行为不变范围和删除门禁。
- [x] 技术方案为客户端、服务端、core、存储、同步和辅助包指定了责任边界。
- [x] 接口、数据格式、配置格式和错误兼容规则已写明。
- [x] 生命周期、迁移失败、import cycle、测试兼容和回滚场景已覆盖。
- [x] 测试计划包含 unit、architecture、integration、CLI、build 和 full verification。
- [x] WBS 按可验证的功能闭环拆分，而不是简单按文件罗列。
- [x] 负责人确认 ownership 表和“行为不变、无隐式数据迁移”门禁。
- [x] 负责人确认后，结论更新为 `PASS` 并进入一次性实现切换。

## Review Dimensions

- [x] 需求和验收覆盖
- [x] 架构和代码归属
- [x] 接口、数据和配置
- [x] 错误、并发、生命周期和安全
- [x] 迁移、契约同步和恢复
- [x] 测试和回归
- [x] 可拆分性

## Decision

当前方案采用一次性结构切换，WBS 仅表示内部依赖顺序，不表示多个中间交付物。负责人已确认目标包归属、整体切换和整体回滚规则，以及不改变用户行为、数据格式和 API 的门禁，结论为 `PASS`，允许执行一次性编码切换。
