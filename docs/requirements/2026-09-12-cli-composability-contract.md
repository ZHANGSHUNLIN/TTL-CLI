# CLI 可组合性契约

日期：2026-09-12
任务：W-012
状态：reviewed

## Problem And Goal

当前 `ttl` 的核心命令主要面向人在终端中使用：正常结果、提示和调试信息可能写入同一输出流；`get` 在匹配多个资源时会读取终端选择；所有执行错误通常都返回退出码 `1`。这些行为便于人工操作，但脚本和 CI 无法稳定地区分结果、诊断和错误类型。

本需求为现有核心资源命令增加一个显式的机器调用模式。调用者能够获得稳定 JSON、明确的 stdin/stdout/stderr 行为和可分类退出码，同时默认的人类可读模式保持现有语义。

## Scope

### In Scope

- 为现有 `add`、`get`、`update`、`del`、`tag`、`dtag` 命令提供 `--json` 机器输出。
- 提供全局 `--non-interactive`，禁止数字选择、确认提示及其他隐式终端读取。
- `add <key> -` 和 `update <key> -` 显式从 stdin 读取完整 value，支持多行内容。
- 约定正常结果只写 stdout，诊断和错误只写 stderr。
- 约定稳定的退出码和机器可读错误编码。
- 为无 TTY、管道输入、JSON 成功和 JSON 错误增加黑盒验收。

### Out Of Scope

- 新增 `list`、`search`、`pick`、`edit`、`copy` 等命令。
- 改变默认文本模式的核心命令语义或文案。
- 为 `sync`、导入导出、工作空间、审计、历史、日志、加密和服务端命令增加 JSON 契约。
- 定义服务端 `/api/v1` 契约或让 TUI 解析 CLI JSON。
- 增加删除确认、批量操作或 shell 自动补全能力。

## Use Cases

1. 脚本运行 `ttl get api-token --json`，从 stdout 解析资源，而不需要解析本地化文本。
2. CI 运行不存在的资源查询，根据退出码和 stderr 中的错误编码判断为 `not_found`，而不是把它当成存储故障。
3. 用户通过管道运行 `ttl add note - --json`，保存多行内容并获得创建结果。
4. 无 TTY 环境运行可能产生多个匹配项的查询时，命令立即返回 `ambiguous`，不等待输入。
5. 人在终端中不传新参数时，继续获得当前文本输出和交互选择。

## Rules

- `--json` 表示机器调用模式，并隐含 `--non-interactive`。
- JSON 成功结果写 stdout；JSON 错误对象写 stderr；错误时 stdout 必须为空。
- JSON 字段名、错误编码和退出码属于稳定契约；本地化 `message` 只供人阅读，不作为机器判断依据。
- `--non-interactive` 可以单独用于文本输出。它禁止隐式读取 stdin；只有 value 参数明确为 `-` 时才允许读取 stdin 到 EOF。
- stdin 内容按原样读取，不自动删除末尾换行；空 stdin 是合法的空 value。
- `get` 无参数保持“列出资源”的现有含义；带查询参数且存在多个匹配时，机器调用模式返回歧义错误，不自动选中。
- 默认文本模式不承诺稳定排版；只有本需求定义的 JSON schema、输出流和退出码属于机器契约。
- 第一版契约版本为 `1`。未来不兼容变更必须引入新版本或经过单独兼容性决策，不能静默修改。

## Edge And Failure Cases

- 参数数量或 flag 无效：返回参数错误，不初始化或修改存储。
- 查询没有结果：返回 `not_found`。
- 查询存在多个结果且禁止交互：返回 `ambiguous`，错误详情可以携带候选 key。
- 添加重复 key：返回 `conflict`，原资源不变。
- stdin 读取失败：返回 `system_error`，不写入资源。
- 存储初始化、读取或写入失败：返回 `system_error`，不得输出伪成功 JSON。
- JSON 编码失败：返回 `system_error`；已经发生的存储写入不得被错误地描述为回滚成功。
- `--debug --json` 下调试信息只能写 stderr，不能污染 stdout JSON。

## Acceptance Criteria

- [ ] 六个核心命令在 `--json` 下输出符合版本 1 schema 的单个 JSON 文档；其他命令使用 `--json` 时明确返回参数错误，不输出混合格式。
- [ ] JSON 成功时退出码为 `0`，stderr 为空；JSON 失败时 stdout 为空，stderr 是单个机器可读错误对象。
- [ ] 参数错误、资源不存在、歧义/冲突和系统错误具有不同且稳定的退出码。
- [ ] `--json` 和 `--non-interactive` 在无 TTY 环境中均不会等待数字选择或确认输入。
- [ ] `add/update <key> -` 能从管道保存多行和空 value，非 `-` 参数不会隐式读取 stdin。
- [ ] 不传新参数时，现有核心命令黑盒回归继续通过。
- [ ] TUI 和服务端不依赖 CLI JSON schema。

## Open Questions And Assumptions

- 第一版只覆盖 `add/get/update/del/tag/dtag`，其余命令按后续实际自动化场景扩展；不支持的命令不能静默忽略机器模式。
- 假设 `get <query>` 暂时保留当前模糊匹配能力；机器模式以 `ambiguous` 代替交互选择。
- 假设项目仍处开发阶段，因此可以建立第一版稳定契约；一旦实现并写入帮助文档，后续变更按兼容接口管理。
