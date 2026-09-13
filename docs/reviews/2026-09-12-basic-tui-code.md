# Basic TUI Code Review

日期：2026-09-12
任务：W-007 / T-02
方案：`docs/tech-designs/2026-09-12-basic-tui.md`
结论：CONDITIONAL

## Findings

- `[MEDIUM]` `internal/client/tui/model.go` user-facing copy：
  - 问题：TUI 主界面、状态、错误提示和帮助文案目前使用英文常量；只有 Cobra 的 `ui` 简介进入了五份 locale。
  - 后果：功能可用且不会破坏数据，但不满足项目“用户可见字符串放入 i18n locale”的完整规范，中文环境会出现混合语言。
  - 修复：在进入 `Done` 前将稳定的 TUI copy 收敛到 locale 或建立明确的 TUI message catalog；按键符号和资源原文不翻译。
  - 验证：i18n locale JSON 校验、语言切换单测和人工检查关键页面。

- `[LOW]` manual platform evidence：
  - 问题：已通过伪终端观察空态、alternate screen 退出、光标和 bracketed-paste 恢复，但尚未由 owner 在真实终端走完创建、value 搜索、详情长文滚动、编辑、标签和删除，也未完成 80 列以下窗口验收。
  - 后果：自动 model/CLI 测试已覆盖状态规则，但真实按键、宽字符和终端体验仍有残余风险。
  - 修复：owner 在本地终端完成 WORK_ITEMS.md 中的人工流程；发现问题则回到 Doing 修正。
  - 验证：记录实际终端、路径和结果。

## Acceptance Coverage

- [x] `ttl ui` 已注册且只允许 SQLite/bbolt，本地/TTY 门禁在 storage open 前执行。
- [x] 浏览、value/key/tag 搜索、详情、创建/修改、标签添加/删除和删除确认复用 application service。
- [x] 空态、无结果、初始读取错误、保存/删除失败、未保存确认、长内容滚动和窄终端状态有 model 测试或实现路径。
- [x] CLI 与 TUI 共享删除编排；SQLite/bbolt 更新时间语义有跨 backend 测试。
- [x] cloud/sync、工作区、`ttl pick`、复制、打开、外部编辑器和批量操作未进入实现。
- [x] `./scripts/verify.sh` 完整通过。
- [ ] TUI 全部稳定文案尚未完成本地化。
- [ ] owner 真实终端完整人工验收尚未完成。

## Review Decision

`CONDITIONAL`。未发现数据丢失、生命周期、CLI 兼容或范围越界 blocker；自动检查和 PTY 冒烟通过。剩余两项是进入 `Done` 前的交付门槛：补齐 TUI 文案本地化，并由 owner 完成真实终端完整流程。任务继续保持 `Doing`，不创建提交。
