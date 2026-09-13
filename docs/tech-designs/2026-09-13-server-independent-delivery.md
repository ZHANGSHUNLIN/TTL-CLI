# Server Independent Delivery Technical Design

日期：2026-09-13
任务：W-010
需求：`docs/requirements/2026-09-13-server-independent-delivery.md`
状态：draft

## Current Behavior

源码边界已按 [`docs/client-server-separation-plan.md`](../client-server-separation-plan.md) 整理，但云端独立制品、部署和回滚尚未落地。

## Proposed Ownership

- `cmd/ttl-server`：服务端构建入口和启动参数。
- CI：服务端制品、校验和、发布及回滚步骤。
- 部署配置：运行用户、数据卷、密钥、日志、健康检查和 TLS 边界。

## Interfaces And Lifecycle

待评审确定制品目录、配置优先级、健康检查协议、优雅关闭超时和版本回滚策略。

## Gate

未完成方案评审前不得宣称独立交付完成。
