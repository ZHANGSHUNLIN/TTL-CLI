# Server Independent Delivery

日期：2026-09-13
任务：W-010
状态：draft

## Problem And Goal

让 `ttl-server` 能独立构建、发布、部署、升级和回滚，运行环境不依赖客户端二进制、客户端配置或源码目录。

## Scope

### In Scope

- 独立 CI 制品、服务端运行包/镜像、配置和密钥注入、数据卷、日志、健康检查、优雅关闭、TLS 和回滚。

### Out Of Scope

- 改变客户端 CLI/TUI 功能、同步协议或仓库拆分为多个 Go module。

## Acceptance Criteria

- [ ] `ttl` 与 `ttl-server` 可单独构建并下载。
- [ ] 服务端运行包不包含客户端运行依赖或源码目录。
- [ ] 部署冒烟、升级和回滚均有可重复证据。
