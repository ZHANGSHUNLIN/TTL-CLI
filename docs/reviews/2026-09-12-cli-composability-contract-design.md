# CLI 可组合性契约设计评审

日期：2026-09-12
任务：W-012
方案：[`docs/tech-designs/2026-09-12-cli-composability-contract.md`](../tech-designs/2026-09-12-cli-composability-contract.md)
结论：PASS

## Findings

- `[MEDIUM]` 需求范围与全局 flag：
  - 问题：第一版只覆盖六个命令，但原方案没有说明其他命令收到全局机器模式 flag 后的行为。
  - 后果：例如 `ttl version --json` 可能成功输出普通文本，脚本会把它误认为 JSON。
  - 修复：设计已明确未覆盖命令必须在存储初始化前返回 `invalid_argument`。
  - 验证：增加未覆盖命令的 CLI 黑盒失败用例，断言 stdout 为空且退出码为 `2`。
- `[MEDIUM]` 默认文本兼容与退出码：
  - 问题：需求要求默认文本模式保持现有语义，原方案却让默认文本模式采用新的分类退出码；当前删除不存在资源还会成功返回。
  - 后果：没有启用机器模式的既有脚本可能发生兼容性变化。
  - 修复：分类退出码只在 `--json` 或 `--non-interactive` 生效；默认文本模式继续使用现有退出行为。
  - 验证：既有黑盒回归通过，并新增机器模式退出码断言。
- `[LOW]` 敏感值与 debug 输出：
  - 问题：资源 JSON 本身可能包含敏感 value，debug 文本又可能破坏单文档约束。
  - 后果：非显式调用可能泄露内容，额外输出会破坏解析。
  - 修复：只有显式 `--json` 输出 value；JSON 模式不额外输出 debug 文本，错误 details 不包含 value、API Key 或完整敏感路径。
  - 验证：stdout/stderr 单文档测试和错误 details 检查。

## Acceptance Coverage

- [x] 六个核心命令、JSON schema、stdin 和非交互范围明确。
- [x] stdout、stderr、错误 code 和退出码可由黑盒测试观察。
- [x] 默认文本兼容和未覆盖命令行为明确。
- [x] W-006 客户端 Service 是业务与 CLI 表示层之间的依赖边界。
- [x] 无存储、配置、HTTP 或持久化格式变更。

## Review Dimensions

- [x] 需求覆盖
- [x] 架构和代码归属
- [x] 接口、数据和配置
- [x] 错误、并发、生命周期和安全
- [x] 迁移、契约同步和恢复
- [x] 测试和回归
- [x] 可拆分性

## Decision

`PASS`。两项兼容性问题已在需求和技术设计中修正；实现可以在 W-006 显式 Service 边界上进行，不需要引入新依赖或改变持久化格式。
