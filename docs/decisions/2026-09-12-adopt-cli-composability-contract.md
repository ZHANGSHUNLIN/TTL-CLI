# Adopt CLI Composability Contract

日期：2026-09-12
状态：adopted
任务：W-012

## 背景

现有 `ttl` 输出和交互主要服务人工终端操作。脚本需要解析本地化文本，无法区分参数、资源不存在和系统错误；多匹配查询还会读取终端选择。W-012 需要建立一套机器可调用契约，同时避免破坏已有文本使用方式。

## 决定

`ttl` 为 `add/get/update/del/tag/dtag` 提供显式 `--json` 和 `--non-interactive` 模式。JSON 使用带 `schema_version: 1` 的统一 envelope 和独立 CLI DTO；成功结果写 stdout，错误对象写 stderr。机器模式使用稳定的分类退出码和非本地化错误 code。

`add/update <key> -` 显式读取 stdin 到 EOF。机器模式不会依赖 TTY 自动检测，也不会进行数字选择。默认文本模式继续保持现有输出、交互和退出行为。

## 备选方案

- 只增加 `--quiet`：不能提供结构化结果和错误分类，不采用。
- 自动根据 TTY 切换 JSON：重定向输出不等于调用者需要 JSON，行为不够可预测，不采用。
- 直接复用数据库模型或服务端响应：会把持久化字段和 HTTP 生命周期泄漏为 CLI 外部契约，不采用。
- 所有命令一次性支持 JSON：范围过大且缺少真实自动化场景，首版只覆盖核心资源命令。

## 影响

脚本能够稳定消费资源和错误，代价是 JSON schema、错误 code 和退出码一旦发布就必须兼容维护。资源 value 可能敏感，因此只有显式 `--json` 才输出完整结构；错误详情不得携带 value 或凭据。默认文本模式保留兼容分支，后续若要统一退出码需要单独决策。

## 验证

通过 formatter 和错误映射单元测试、六个核心命令的构建后二进制黑盒测试、stdin 多行输入、无 TTY 歧义、stdout/stderr 分离、退出码 0/1/2/3/4，以及既有 `scripts/regression.sh` 验证。实现和验证完成后将状态改为 `adopted`。
