# Separate Client And Server Layout

日期：2026-09-12
状态：adopted
任务：W-004

## 背景

当前仓库把本地 CLI、远端 HTTP 服务、云端客户端、多租户存储和共享模型放在根目录、`command/`、`api/` 与 `db/` 中。`main.go` 同时定义客户端和服务端运维命令，`db/` 同时包含 SQLite、bbolt、HTTP 客户端和租户能力。目录名称无法准确表达代码属于本地客户端、远端服务还是共享核心。

项目方向还计划增加 TUI。如果继续让 TUI 直接依赖 Cobra 全局状态和 `db.Stor`，客户端入口之间难以共享可测试的业务逻辑。

## 决定

提议将仓库整理为一个 Go module 下的两个可执行入口和三个明确边界：

- `cmd/ttl` 构建本地 CLI/TUI 客户端。
- `cmd/ttl-server` 构建远端多租户服务。
- `internal/client`、`internal/server`、`internal/core` 分别承载客户端用例、服务端用例和共享领域契约。

SQLite、bbolt 等具体存储实现放在 `internal/storage`。客户端与服务端不得直接互相 import，只通过共享 core 或 `/api/v1` 契约交互。实施分阶段进行，先增加独立入口和兼容层，再迁移实现并去除全局存储。

## 备选方案

- 使用顶层 `frontend/` 和 `backend/`：名称直观，但无法合理安置共享模型和存储接口，容易形成两个新的大杂物包。
- 拆为客户端与服务端两个仓库：隔离最强，但当前共享数据模型和存储语义较多，会立刻引入版本联动、发布与兼容维护。
- 保持当前包，只补充文档：改动最小，但不能解决 `main.go`、`db/` 和全局存储造成的实际依赖交叉。
- 立即拆为多个 Go module：能强化依赖边界，但当前阶段会增加 replace、版本和测试矩阵成本。

## 影响

收益是客户端、服务端和共享核心能从路径直接识别，CLI 与 TUI 可以复用业务用例，服务端也能独立构建与部署。代价是迁移会影响大量 import、测试和构建脚本，并增加第二个二进制的发布维护。

迁移期间必须保留现有 `ttl` 命令、数据文件格式和 `/api/v1` 契约。`ttl server` 先作为兼容入口保留，是否弃用需要在正式发布前单独决定。阶段 A 已建立 `cmd/ttl-server` 和 `internal/server/cli`，因此该决定已进入采用状态。

## 验证

按 `docs/client-server-separation-plan.md` 分阶段验证两个二进制构建、包依赖方向、CLI 黑盒回归、API/同步集成测试、旧数据兼容以及 `go vet`。第一阶段完成且目录开始遵循该方案后，将状态更新为 `adopted`。
