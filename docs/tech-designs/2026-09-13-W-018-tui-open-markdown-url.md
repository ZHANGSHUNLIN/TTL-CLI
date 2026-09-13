# 技术方案：修复 TUI 打开 Markdown 链接失败

任务：W-018
类型：bugfix
状态：accepted
日期：2026-09-13

## 根因

TUI 和 CLI 的平台打开分支直接把资源值传给 `open`/`explorer`。资源值允许保存 Markdown 链接，完整字符串不是平台打开器可识别的目标。

## 修复边界

- 新增 `internal/client/opener`，提供无平台副作用的 `Target` 解析函数。
- `Target` 对完整 Markdown 链接提取目标，对非 Markdown 的非空值保持原样，对空值和不完整链接返回错误。
- TUI 的 `openExternalResource` 和 CLI `open` 命令在执行平台命令前调用该 helper。
- 不改变平台分支、错误退出码、TUI 成功退出和失败留在详情页的状态机。

## 兼容性与回滚

纯 URL 和原有非空值保持兼容；没有数据迁移或协议变化。若修复异常，可回滚本次客户端代码和测试提交，不影响已有资源数据。

## 风险

Markdown 目标包含空格或嵌套括号时不做复杂渲染，返回错误而不是把原始标记交给系统；这避免误打开错误目标。
