# Regression Testing

本项目采用适合个人练习的分层回归测试。它借鉴 DeepSeek Harness 的测试思想，但不引入 GitHub CI、模型响应录制、浏览器快照或跨平台矩阵。

## 四层测试

### 1. 聚焦单元测试

用于验证一个包或一个局部规则：

```sh
go test ./command -run TestName
go test ./db -run TestName
```

修改小范围逻辑时先运行这一层，反馈最快。

### 2. 集成测试

用于验证数据库、加密、API、同步、文件系统、旧数据兼容和跨包生命周期：

```sh
go test ./...
go test ./integration_test/...
go test -race ./...
go vet ./...
```

集成测试应使用临时目录、临时数据库和 `httptest`，不要依赖个人的 `~/.ttl` 或固定端口。涉及默认配置和密钥路径的测试包会在 `TestMain` 中设置独立的临时用户目录；旧数据兼容测试直接生成拆分前的 SQLite 表和 bbolt bucket/JSON 格式，再由当前存储实现读取。

`scripts/verify.sh` 还会为整个验证进程设置临时 `HOME`，同时保留当前 `GOPATH`。这是对所有命令的第二层隔离；测试自身不能把脚本隔离当作读取个人配置的理由。

### 3. CLI 黑盒回归

用于验证用户真正调用的构建后 `ttl` 二进制：

```sh
./scripts/regression.sh
```

脚本会创建临时 `HOME`、配置、数据库和加密密钥，覆盖资源生命周期、标签、导入导出、日志、历史、审计、加密和工作空间，并检查命令输出、实际文件位置和后续数据状态。它不会修改真实用户数据。

如果已经有构建产物，也可以把它作为参数传入：

```sh
./scripts/regression.sh ./ttl
```

### 4. 完整验证

用于进入本地提交前的宽范围检查：

```sh
./scripts/verify.sh
```

完整验证依次执行补丁/Go/shell 格式检查、双二进制构建、架构依赖检查、CLI 黑盒回归、全部 Go 测试、集成测试、race 检查和 `go vet`。架构检查通过 `go list -json` 分析直接与传递依赖；服务端二进制一旦依赖客户端、旧 `command`、旧 `db` 或旧 `sync` 包就会失败。

## 自动验收证据

| 验收目标 | 自动化证据 |
| --- | --- |
| 补丁无空白错误，Go 与验证脚本格式有效 | `git diff --check`、`gofmt -s -l .`、`bash -n` |
| 客户端和服务端均可构建 | `scripts/verify.sh` 的双二进制构建 |
| 用户入口和持久化生命周期可用 | `scripts/regression.sh` |
| 自定义配置、数据库、密钥和工作空间不写入真实用户目录 | 黑盒文件路径断言和测试包临时 `HOME` |
| client、server、core 依赖方向受控 | `internal/architecture/dependencies_test.go` |
| 拆分前 SQLite/bbolt 数据可继续读取 | `TestSQLiteStorage_LegacyDataCompatibility`、`TestBboltStorage_LegacyDataCompatibility` |
| 并发访问没有已知数据竞争 | `go test -race ./...` |

人工验收仍负责判断产品目标和迁移阶段是否完成，例如是否可以移除 `ttl server` 兼容入口、是否已完成 TUI，以及当前改动是否符合任务范围。自动检查不能替代这些产品决策。

## 按变更选择检查

| 变更范围 | 建议检查 |
| --- | --- |
| 一个 helper 或局部规则 | 聚焦单元测试 |
| 命令行为、CLI 输出或数据生命周期 | 相关单元测试 + `./scripts/regression.sh` |
| 数据库、加密、API、同步、迁移或跨包逻辑 | `go test ./...`、`go test ./integration_test/...`、`go test -race ./...`、`go vet ./...` |
| 依赖变化、宽范围重构或提交前总检查 | `./scripts/verify.sh` |
| 只改文档、Skill 或流程 | `git diff --check` + 人工检查链接和内容 |

不需要每次都运行完整验证。选择检查时，目标是让证据覆盖变更的真实入口和风险，而不是累积命令数量。

## 与人工审核的关系

自动测试只回答“已执行的行为是否符合断言”。进入 `Review` 后仍要人工检查：

- 黑盒测试是否覆盖了任务验收条件；
- 命令是否真的读取了临时文件、数据库或配置；
- 是否误改了真实用户数据、凭据或无关文件；
- 失败路径、迁移、删除、加密和同步风险是否被检查；
- 用户可见文案和文档是否与实现一致。

测试证据写回 `WORK_ITEMS.md`，任务只有在自动检查、人工审核和本地 commit 都完成后才能进入 `Done`。
