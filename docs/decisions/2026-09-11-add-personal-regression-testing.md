# Add Personal Regression Testing

日期：2026-09-11
状态：adopted
任务：W-002

## 背景

项目已有单元测试、集成测试和一个包含 CLI 流程的 `scripts/verify.sh`，但各层职责没有单独说明，CLI 黑盒回归也不能独立运行。DeepSeek Harness 的分层测试经验适合借鉴，但它的 Session Snapshot、真实模型、浏览器和 GitHub CI 对个人 Go CLI 练习项目过重。

## 决定

项目采用四层个人回归体系：聚焦单元测试、集成测试、构建后二进制 CLI 黑盒回归和完整验证。新增 `scripts/regression.sh` 作为黑盒回归入口，`scripts/verify.sh` 作为构建、黑盒、单测、集成测试和 `go vet` 的总入口。黑盒测试必须使用临时 `HOME`、配置、数据库和加密密钥。

## 备选方案

- 继续只保留一个大验证脚本：可以运行，但不能按变更范围选择检查，失败定位也更粗。
- 原样移植 DSH 的 Session Snapshot 和真实模型回归：能覆盖更多 Agent 场景，但当前 CLI 没有对应的模型会话，维护成本与收益不匹配。
- 只运行 Go 包测试：反馈快，但可能遗漏真实命令参数、配置加载、输出和构建产物问题。

## 影响

个人开发可以按风险选择最小检查，也可以用一个入口完成提交前验证。黑盒测试会增加少量脚本维护成本，但能证明真实 CLI 入口和持久化结果。项目暂不提供远程 CI、浏览器快照、真实外部 API 和跨平台矩阵，这些能力在需要相应运行环境时再单独增加。

## 验证

`docs/regression-testing.md`、`personal-regression-testing` Skill、`scripts/regression.sh` 和 `scripts/verify.sh` 共同定义并实现该流程；执行 `git diff --check`、Go 格式检查、静态检查和可用的测试命令确认落地结果。
