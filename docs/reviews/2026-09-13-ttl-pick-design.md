# TTL Pick Design Review

日期：2026-09-13
任务：W-013
方案：`docs/tech-designs/2026-09-13-ttl-pick.md`
状态：PASS

## Findings

- 范围已确认纳入当前阶段；v1 不提供 `--print`。
- 空数据库和无匹配统一为 `not_found`；query 首尾空格按原样匹配。
- stdout 只承载选中 value，候选、提示和诊断写 stderr；无 TTY 在读取 stdin 前返回 `interaction_required`。
- 复用现有 ORIGIN-only 搜索和稳定排序，不改变 `get`，且 `pick` 不记录 history/audit。
- 测试边界覆盖参数、候选顺序、输出流、取消、非法选择、EOF、无 TTY、只读和黑盒回归。

## Review Gate

- [x] 范围和首期承诺已确认
- [x] 接口、stdout/stderr、退出码和无 TTY 行为已设计
- [x] 测试和回归边界已评审

## Decision

PASS。方案允许进入编码阶段。
