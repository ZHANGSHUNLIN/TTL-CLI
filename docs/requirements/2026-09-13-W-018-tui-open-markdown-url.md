# 需求分析：修复 TUI 打开 Markdown 链接失败

任务：W-018
类型：bugfix
状态：accepted
日期：2026-09-13

## 现象

TUI 详情页按 `o` 打开资源时，如果资源值是 Markdown 链接，例如 `[www.example.com](http://www.example.com)`，macOS `open` 收到的是完整 Markdown 文本，返回 `exit status 1`，TUI 停留在详情页并显示打开失败。

## 影响范围

- 受影响：TUI 详情页的外部打开动作；同样的输入也会影响 `ttl open`。
- 不受影响：资源存储格式、详情页展示、纯 URL 值、Linux 不支持打开的既有行为、打开成功后的 TUI 退出行为。

## 期望行为

1. 完整 Markdown 链接值提取括号中的目标 URL，再交给当前操作系统的打开器。
2. 纯 URL 值保持原样传递。
3. 空值或不完整 Markdown 链接返回清晰错误；TUI 保留详情页和错误状态，不退出。

## 验收标准

- [x] `[label](https://example.com)` 被解析为 `https://example.com`。
- [x] 纯 URL 解析结果不变。
- [x] 空值和不完整 Markdown 链接被拒绝并返回错误。
- [x] TUI 打开成功仍请求退出，打开失败仍保留详情页。
- [x] `ttl open` 复用相同解析逻辑。

## 备注

不新增 Markdown 渲染器，不改变资源值本身；只在调用平台打开器前做输入归一化。
