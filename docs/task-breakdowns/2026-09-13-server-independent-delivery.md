# W-010 WBS：云端服务独立交付

日期：2026-09-13
任务：W-010
需求：`docs/requirements/2026-09-13-server-independent-delivery.md`
技术方案：`docs/tech-designs/2026-09-13-server-independent-delivery.md`
方案评审：`docs/reviews/2026-09-13-server-independent-delivery-design.md`
状态：实现完成，待交付评审

## WBS

### T-01 构建与发布制品

- 输入：通过评审的构建/发布契约和现有 GitHub Actions workflow
- 输出：客户端/服务端独立 CI job、发布矩阵、确定性归档名、校验值和清单检查
- 依赖：方案评审通过
- 验收：`ttl` 从 `./cmd/ttl` 构建，`ttl-server` 从 `./cmd/ttl-server` 构建；Linux
  `amd64`/`arm64` 制品存在；制品不包含源码、客户端或数据
- 证据：workflow 运行记录或本地等价输出、归档清单和校验结果

### T-02 服务端运行时生命周期

- 输入：通过评审的运行时接口，以及现有 `internal/server/api`、`internal/server/tenant`
- 输出：显式监听地址、配置校验、`/healthz`、就绪切换、信号/context 关闭、超时排空和
  租户存储清理
- 依赖：方案评审通过；本地测试不依赖 T-01
- 验收：启动失败不绑定端口；健康端点免认证且不泄露敏感信息；SIGTERM/SIGINT 或
  context 取消在限定时间内关闭资源；现有 `/api/v1` 行为不变
- 证据：单元测试、服务端集成测试和 loopback 进程冒烟

### T-03 参考部署与恢复

- 输入：T-01 制品和 T-02 运行时契约
- 输出：`docs/server-deployment.md`、systemd/等价服务示例、权限、反向代理/TLS 边界、
  备份、升级、探活和回滚流程
- 依赖：T-01、T-02
- 验收：非 root 的临时部署可以从干净 release 目录启动，数据独立持久化，通过健康/API
  探针，切换版本，并在探活失败后恢复旧版本
- 证据：命令记录、临时目录布局和回滚结果，不要求真实生产主机

### T-04 交付回归与兼容性

- 输入：T-01 至 T-03 的实现
- 输出：确定性的 `scripts/server-delivery-smoke.sh` 以及更新后的测试和验收记录
- 依赖：T-02；T-01 完成后打包可以与 T-03 并行
- 验收：现有客户端、API、用户管理、租户隔离和集成测试通过；交付冒烟覆盖干净目录
  启动、损坏状态、端口占用、关闭和回滚
- 证据：`./scripts/verify.sh`、交付冒烟输出和人工清单

## 执行顺序

```text
方案评审 -> T-01 与 T-02 -> T-03 -> T-04 -> 交付评审
```

T-01 和 T-02 在方案通过后可独立测试；T-03 必须使用前两项确定的真实制品布局和运行
参数；T-04 是交付门禁，不能用一次成功的 `go build` 代替。

## 交付门禁

- [x] 需求已评审，高影响假设已确认
- [x] 技术方案已评审并得到 `PASS` 或解决后的 `CONDITIONAL`
- [x] T-01 至 T-04 已附本地自动化验收证据；真实 tag 归档证据仍待补充
- [x] 参考部署和回滚流程已在临时目录自动演练；负责人现场部署确认待完成
- [x] 本地服务端交付冒烟确认制品没有泄露客户端、源码或数据文件；真实 tag 制品待核对

## 当前证据

- `./scripts/verify.sh` 已通过，覆盖双二进制构建、交付冒烟、升级回滚、架构检查、CLI 黑盒
  回归、全量测试、集成测试、race 和 vet。
- `scripts/server-delivery-smoke.sh` 和 `scripts/server-upgrade-rollback-smoke.sh` 已通过。
- 测试环境验收已通过，任务进入 commit 门禁；不得在本地 commit 创建并补充提交证据前标记 Done。
