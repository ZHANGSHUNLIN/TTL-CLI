# CLI 可组合性契约技术设计

日期：2026-09-12
任务：W-012
需求：[`docs/requirements/2026-09-12-cli-composability-contract.md`](../requirements/2026-09-12-cli-composability-contract.md)
状态：reviewed

## Current Behavior

`cmd/ttl` 调用 `internal/client/cli.Run`，由 Cobra 解析参数并组装顶层 `command` 包中的全局命令。命令通过 `command.OutputWriter` 输出文本，但部分入口和错误仍直接使用 `fmt.Print*`；`Run` 将所有执行错误打印到 stdout，并统一返回退出码 `1`。

`get` 通过 `PossiblyRun` 做 key/tag 模糊匹配。多匹配时直接使用 `fmt.Scan` 读取选择，因此无 TTY 调用可能阻塞或得到难以分类的输入错误。`add` 和 `update` 要求 value 必须来自位置参数，尚不支持显式 stdin。业务错误目前主要是本地化字符串，调用方无法稳定分类。

本设计依赖 W-006/T-01 提供结构化资源服务和可判断错误；W-012 不应继续在 Cobra handler 中复制资源查找和写入规则。

## Proposed Ownership

| Area | Owner | Responsibility |
| --- | --- | --- |
| 资源查询和修改 | W-006 建立的 `internal/client` 服务层 | 返回结构化资源结果和 typed error；不感知 JSON、终端或 Cobra |
| CLI 参数与模式 | `internal/client/cli` | 定义 `--json`、`--non-interactive`，把参数和 stdin 转成服务调用 |
| CLI 输出适配 | `internal/client/cli` 下的小型 formatter/error mapper | 生成版本化 JSON、文本结果、stderr 错误和退出码 |
| 进程入口 | `internal/client/cli.Run` / `cmd/ttl` | 注入 stdin/stdout/stderr，执行命令，关闭存储并返回最终退出码 |
| 黑盒契约 | `scripts/regression.sh` 及独立 CLI 测试 | 从构建后的 `ttl` 验证输出流、JSON 和退出码 |

不新增第三方依赖。JSON 使用标准库 `encoding/json`。

## Interfaces

### Flags And Input

根命令新增两个 persistent flags：

```text
--json              输出版本 1 JSON，并隐含非交互模式
--non-interactive   禁止隐式终端输入，保留文本输出
```

`add` 和 `update` 的位置参数数量保持不变。当 value 等于单独的 `-` 时，handler 从 `cmd.InOrStdin()` 读取到 EOF；其他值继续按现有规则处理。这样不需要检测 TTY，也不会把管道存在误判成用户要求读取。

只有 `add/get/update/del/tag/dtag` 接受机器模式。其他命令出现 `--json` 或 `--non-interactive` 时返回 `invalid_argument`，避免调用者把文本误当 JSON。该检查在存储初始化前完成。

### JSON Schema

成功对象统一使用 envelope：

```json
{
  "schema_version": 1,
  "ok": true,
  "data": {}
}
```

资源使用 CLI DTO，不直接暴露数据库模型的 `val`、`tag`、`originKey` 字段名：

```json
{
  "key": "api-token",
  "value": "secret",
  "tags": ["work"],
  "created_at": 1789142400,
  "updated_at": 1789142400
}
```

- `get` 无参数：`data.resources` 为数组，按当前文本列表采用的确定性顺序输出。
- `get <query>`、`add`、`update`、`tag`、`dtag`：`data.resource` 为资源对象。
- `del`：`data` 为 `{"key":"...","deleted":true}`。
- 空 tags 编码为 `[]`，不编码为 `null`。时间字段保持整数 Unix 秒。

错误对象写 stderr：

```json
{
  "schema_version": 1,
  "ok": false,
  "error": {
    "code": "not_found",
    "message": "resource not found",
    "details": {}
  }
}
```

`code` 稳定且不本地化；`message` 可以本地化；`details` 只放非敏感结构化上下文。`ambiguous` 可以返回排序后的候选 key，但不得包含资源 value。

### Exit Codes

| Exit code | Stable error code | Meaning |
| --- | --- | --- |
| `0` | 无 | 成功 |
| `1` | `system_error` | 存储、配置、网络、I/O、编码或未分类内部错误 |
| `2` | `invalid_argument` | Cobra 参数、flag 或输入格式错误 |
| `3` | `not_found` | 目标资源不存在 |
| `4` | `conflict` / `ambiguous` / `interaction_required` | 重复、多个候选或当前模式不允许交互 |

退出码只表达粗粒度类别；脚本需要精确原因时读取 JSON `error.code`。分类退出码只在 `--json` 或 `--non-interactive` 模式生效；默认文本模式继续沿用成功为 `0`、执行错误为 `1` 的现有行为，避免本任务扩大兼容性变化。

### Internal Result And Error Boundary

W-006 服务层应至少能够表达以下稳定类别，而不是让 CLI 解析错误文本：

```go
type ErrorKind string

const (
    ErrorNotFound ErrorKind = "not_found"
    ErrorConflict ErrorKind = "conflict"
    ErrorAmbiguous ErrorKind = "ambiguous"
)
```

具体 Go 类型和文件名可在 W-006 设计中确定，但 W-012 的 formatter 只能通过 `errors.Is/errors.As` 或等价 typed result 分类。Cobra 参数错误在 CLI 层映射为 `invalid_argument`，未识别错误统一安全地映射为 `system_error`。

`Run` 调整为接收或内部构造明确的 stdin/stdout/stderr，并且只在一个位置输出最终错误。命令 handler 不再自行打印同一错误，避免重复文本或两个 JSON 文档。

## Data And Configuration

- 不改变 SQLite、bbolt、远端 API、配置文件或资源持久化格式。
- CLI JSON DTO 是新的外部接口，不复用持久化 JSON tag，避免数据库模型重构破坏脚本。
- 不增加环境变量或配置开关；机器模式必须由每次调用的 flag 显式启用。
- schema 版本固定为数字 `1`。第一版不提供用户选择其他版本的 flag。

## Lifecycle And Failure Behavior

- Creation：根命令创建 invocation-scoped options 和输入输出流；存储仍由 W-006 设计的客户端生命周期负责。
- Mutation：handler 先完整读取和校验输入，再调用服务；成功后只编码一次结果。
- Cleanup：无论命令成功还是失败都关闭已创建的存储；关闭失败在原命令成功时转为 `system_error`，在已有错误时作为 stderr 诊断保留原始退出类别。
- Failure and recovery：服务调用失败时不输出成功 envelope；JSON 编码失败返回 `system_error`。若写操作已成功但结果编码失败，不尝试反向修改存储，错误信息明确结果输出失败，调用者可重新读取确认。
- Concurrency：本任务不改变当前数据库并发模型。

## Security And Compatibility

- JSON 可能包含资源 value，只有用户显式传 `--json` 才输出；错误详情不得包含 value、API Key、配置密钥或完整敏感路径。
- `--debug --json` 的诊断只写 stderr。JSON 错误模式下 stderr 仍必须保持单个错误对象，因此 debug 详情放入受控的 `error.details`，不得额外打印文本行。
- 不带新 flag 的文本、交互和本地化行为保持现有黑盒兼容。
- 默认文本模式下删除不存在资源仍保持当前成功返回和提示；机器模式将同一情况表达为 `not_found`。这是为保留旧脚本行为而设置的显式兼容分支。
- `--json` 一旦实现并文档化，schema、字段含义、错误 code 和退出码按对外兼容接口管理。
- 本契约需要新增 CLI 兼容性决策记录；实现前把 W-012 从候选范围转为阶段承诺。

## Alternatives

- 只提供 `--quiet`：只能减少文本，不能表达结构化结果和错误，因此不采用。
- 每条命令自行定义 JSON：实现简单但字段和错误会漂移，因此采用统一 envelope 和 CLI DTO。
- 复用服务端 API envelope：CLI 与服务端生命周期和错误类别不同，会造成不必要耦合，因此只借鉴 envelope 思路，不共享 DTO。
- 自动检测 stdout 是否为 TTY：脚本重定向不等于用户一定需要 JSON，行为也难预测，因此采用显式 flag。
- JSON 错误写 stdout：会污染只消费成功结果的管道，因此错误对象写 stderr，stdout 保持为空。

## Test And Regression Plan

| Layer | Evidence |
| --- | --- |
| Unit | formatter 对每种成功 DTO、空 tags、错误 kind 和编码失败的测试；stdin 多行/空输入测试；非交互歧义测试 |
| Client service | W-006 的精确匹配、无结果、多个结果、重复 key 和存储错误测试 |
| CLI command | 注入 bytes.Buffer 验证 stdout/stderr 分离、单个 JSON 文档和参数错误映射 |
| CLI black-box | 构建 `./cmd/ttl`，在临时 HOME/conf 下验证六个命令、管道输入、无 TTY、退出码 0/1/2/3/4 和 JSON schema |
| Compatibility | 运行既有 `./scripts/regression.sh`，确认不带 flag 的核心命令行为不变 |
| Full | `go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...`、`git diff --check` |

黑盒测试解析 JSON 字段，不用字符串包含代替 schema 验证；错误测试分别捕获 stdout、stderr 和退出码。所有测试继续使用临时目录，不能读写真实用户数据目录。

## Rollout And Recovery

- 先完成 W-006 的结构化服务和 typed error，再实现 formatter、全局 flags，最后逐条接入六个核心命令。
- 接入过程中默认文本路径始终由现有回归保护；某条命令尚未接入时不得接受 `--json` 后输出文本伪装成功。
- 若实现阶段发现 schema 需要不兼容调整，在首次发布前更新本设计和决策记录；发布后通过新 schema version 演进。
- 回滚 W-012 只删除新增机器模式，不回滚或复制 W-006 服务层。
