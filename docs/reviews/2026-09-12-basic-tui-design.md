# Basic TUI Design Review

日期：2026-09-12
任务：W-007 / T-02
方案：`docs/tech-designs/2026-09-12-basic-tui.md`
结论：PASS

## Findings

- `[BLOCK]` Command contract / `internal/client/cli/root.go:newPreRun`：
  - 问题：方案要求 `ui` 对 `cloud`/`sync` 在不发网络请求的前提下报错，但现有 `PersistentPreRunE` 会在 `RunE` 之前调用 `app.OpenStorage`。如果门禁只放在 `newUICommand`，远端连接已经发生。非 TTY 检查同样会在数据库打开之后才执行。
  - 后果：基础 TUI 会违反“本地闭环不依赖网络”的范围边界，并在本应立即拒绝的调用中触发连接等待或远端副作用。
  - 修复：已在方案中定义 pre-run 级别的 UI 门禁及可测试 seam；先解析实际 storage type、验证仅为 SQLite/bbolt 且 stdin/stdout 为 TTY，再调用 `OpenStorage`。门禁失败时不得构造 service。
  - 验证：CLI 单测用 fake opener/TTY detector 断言远端和非 TTY 路径的 open 调用次数为 0；黑盒验证无 TTY 快速失败。
  - 状态：已整改。

- `[BLOCK]` Application contract / delete use case：
  - 问题：方案让 TUI 直接调用 `DeleteResourceByKey`，而现有 CLI 删除在 handler 中先调用 `CleanupResourceHistory` 再删资源。两种入口会得到不同的审计/历史结果，删除编排仍留在 Cobra 表示层。
  - 后果：违反 W-007“CLI/TUI 共享业务规则”和“TUI 不复制 CLI 业务规则”的核心验收；从 TUI 删除会留下与 CLI 不同的辅助数据。
  - 修复：方案已把删除编排收敛为 `app.Service.DeleteResourceWithCleanup`，规定清理为 best-effort、资源删除为 hard-fail，并由结构化 `DeleteResult` 返回清理警告；CLI 与 TUI 同时使用该入口。
  - 验证：service 单测覆盖清理成功、清理失败但资源删除成功、资源删除失败；CLI 既有 debug 输出和默认行为保持；TUI 只消费结构化结果。
  - 状态：已整改。

- `[MEDIUM]` Application contract / search scope：
  - 问题：当前 `FindResources` 遍历全部 storage entries；方案只写了补 value 匹配，没有明确排除 `TAG` alias entry。
  - 后果：bbolt 或兼容数据中可能把标签索引项作为可浏览资源返回，列表/详情和删除目标不再只对应 origin resource。
  - 修复：方案已冻结搜索仅返回 `Key.Type == ORIGIN` 的资源，key/value/tag 匹配都在 origin value 上执行；空 query 与 `ListResources` 语义一致。
  - 验证：service 单测加入 origin 与 tag alias 同时存在的 fixture，断言结果不含 alias。
  - 状态：已整改。

- `[MEDIUM]` Mutation semantics / metadata preservation：
  - 问题：方案要求刷新后展示持久化时间，却没有处理当前 `UpdateResourceValue` 与 bbolt `UpdateResource` 之间的元数据差异：service 构造的新 value 未携带 `CreatedAt`/`UpdatedAt`，bbolt update 会写入零值，而 SQLite save 会自行保留创建时间。
  - 后果：TUI 编辑同一资源后，不同本地 backend 展示不同且可能丢失的时间元数据，也削弱 CLI/TUI 交叉观察的可信度。
  - 修复：方案已规定 application 传递原 `CreatedAt`，storage adapter 保留创建时间并生成新的 `UpdatedAt`；SQLite/bbolt 必须行为一致。
  - 验证：跨 backend 测试断言修改 value 后创建时间不变、更新时间不早于原值、tags 不变。
  - 状态：已整改。

- `[LOW]` TUI runner seam：
  - 问题：`RunOptions.Width/Height` 与 Bubble Tea 的 `WindowSizeMsg` 责任重叠，TTY 判断也只描述为 adapter 行为但没有依赖注入点。
  - 后果：实现时容易把真实终端探测带进 model 单测，或出现初始尺寸与 resize 消息冲突。
  - 修复：方案已将尺寸定义为可选测试初值，运行时以 `WindowSizeMsg` 为准；CLI options 注入 opener、TTY detector 和 TUI runner。
  - 验证：model 尺寸测试不访问真实终端；CLI unit test 使用 fake detector/runner。
  - 状态：已整改。

## Acceptance Coverage

- [x] 浏览、详情、创建/修改、标签、删除确认有状态设计。
- [x] 空数据库、搜索无结果、内容过长、未保存编辑、保存/删除错误和窄终端有可观察行为。
- [x] 远端、同步、工作区、复制、打开、外部编辑器和批量操作被排除。
- [x] 本地存储和 TTY 门禁明确发生在资源初始化前，并有 fake opener 断言。
- [x] CLI/TUI 删除共享一个完整 application use case 和清理警告结果。
- [x] origin-only 搜索和本地 backend 元数据一致性已冻结。

## Review Dimensions

- [x] 需求覆盖
- [x] 架构和代码归属
- [x] 接口、数据和配置
- [x] 错误、并发、生命周期和安全
- [x] 迁移、契约同步和恢复
- [x] 测试和回归
- [x] 可拆分性

## Decision

`PASS`。初审的两个 BLOCK 和三个补充问题均已在技术方案中整改：UI 门禁前置到 storage open 之前，删除编排收敛到共享 application use case，搜索限定 origin resource，并统一 SQLite/bbolt 的更新时间责任。W-007 可以按修订后的 T-02 拆分进入编码，但不得扩大到远端/同步、工作区、`ttl pick` 或其他候选能力。
