# 交付验收：修复 TUI 打开 Markdown 链接失败

任务：W-018
日期：2026-09-13
状态：accepted

## 缺陷关闭条件

- [x] Markdown 链接值交给平台打开器前已提取目标 URL。
- [x] 纯 URL 行为保持兼容。
- [x] 无效值错误清晰，TUI 失败时保留详情页。
- [x] TUI 成功打开后仍自动退出。
- [x] 代码评审、测试和 CLI 回归证据齐全。

## 人工验收

已通过单元测试验证 `[www.example.com](http://www.example.com)` 提取为 `http://www.example.com`，并通过 macOS 临时 `open` 可执行文件捕获参数，确认 TUI 集成传递目标 URL；现有 TUI 状态测试确认打开成功后退出、失败时保留详情页。未启动真实浏览器。

## 遗留风险

不支持复杂 Markdown 语法（引用标题等）；此类值会明确失败，不会修改原资源。当前实现对 URL 内的普通括号保持原样传递。
