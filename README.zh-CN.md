<div align="center">

# ttl

### 你的个人知识归档

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[English](README.md) | [日本語](README.ja.md) | [Español](README.es.md) | [Français](README.fr.md) | [Português](README.pt.md)

---

*一个轻量级的个人数据管理 CLI 工具。将任何内容以键值对形式存储，即时搜索。*

</div>

---

## 📖 项目故事

每个开发者都经历过这样的时刻：

> "那个上个月用的 Docker 命令是什么来着？"
> "同事分享的配置文件放哪了？"
> "我看过一篇相关的文章……但找不到了。"

我们把知识存得到处都是 —— 浏览器标签、Slack 消息、邮件、书签夹、笔记应用。真正需要的时候，却要在无尽的标签和历史记录中浪费时间翻找。

**这也是我的痛点。**

所以我开发了 **ttl**。

名字来自 "Time to Live" —— 但有不一样的含义。不是过期，而是让知识**永久保存**。

- 把所有内容存到一个地方，用键值对的方式
- 加上标签，方便整理
- 按关键词即时搜索

不再翻找旧邮件、滚动聊天记录。只需要 `ttl get <关键词>`，就能立刻找到。

**ttl 就像你的个人知识归档库 —— 你需要什么，什么时候都有。**

---

## ✨ 功能特性

| 功能 | 描述 |
|------|------|
| 🗄️ **本地 KV 存储** | 快速、零配置的嵌入式数据库 (bbolt) |
| 🏷️ **标签系统** | 用灵活、可搜索的标签组织资源 |
| 🔍 **模糊搜索** | 跨键名和标签即时查找所需内容 |
| 📝 **工作日志** | 记录并按条件筛选每日工作 |
| ☁️ **云端同步** | 连接由独立工程维护的 TTL 后端服务 |
| 🚀 **智能打开** | 用系统默认程序打开 URL 和文件 |
| 📤 **数据导出** | 导出为 JSON 或 CSV 格式 |

---

## 🚀 快速开始

### 安装

#### Linux / macOS

```bash
# 从 GitHub releases 安装
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"

# 或从源码构建
go build -o ttl ./cmd/ttl
sudo mv ttl /usr/local/bin/
```

#### Windows

```powershell
# 从 GitHub releases 安装
irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

#### 自定义下载地址

适用于内网或自定义镜像：

```bash
# Linux/macOS
TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-linux-amd64" /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"
```

```powershell
# Windows
$env:TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-windows-amd64.zip"; irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

### 基础用法

```bash
# 添加资源
ttl add my-link https://example.com

# 添加并打标签
ttl add docker-cmd "docker run -d -p 8080:80 nginx"
ttl tag docker-cmd dev ops

# 搜索资源
ttl get docker

# 在浏览器中打开
ttl open my-link

# 删除
ttl del old-key
```

---

## 📝 工作日志

记录并按条件筛选每日工作。

```bash
# 写日志
ttl log write "完成用户模块重构" --tags "项目A,开发"

# 查看日志
ttl log list                    # 今天的日志
ttl log list --range week       # 本周
ttl log list --range month      # 本月

```

---

## ☁️ 云端服务与同步

后端由独立服务工程维护和部署。本仓库只包含 `ttl` 客户端及其 HTTP 适配器，不构建或发布服务端可执行文件。

### 同步数据

```bash
# 配置远程服务器
ttl config
# 编辑 server 部分，填入端点和 API 密钥

# 同步本地与远程数据
ttl sync
```

**架构特点：**
- 多租户设计，每用户独立数据库
- API Key 认证
- REST API 用于程序化访问

---

## ⚙️ 配置

配置文件：`~/.ttl/ttl.ini`

```ini
[default]
db_path = ~/.ttl/data.db

[server]
endpoint  = https://your-server.com
api_key   = your-user-api-key
```

```bash
# 查看当前配置
ttl config

```

---

## 📤 导出数据

```bash
# 导出为 JSON
ttl export --format json

# 导出为 CSV
ttl export --format csv

# 导出到指定文件
ttl export --format json --output backup.json
```

---

## 🏗️ 项目结构

```
ttl-cli/
├── cmd/ttl/                    # 客户端入口
├── internal/client/            # CLI、TUI、远端访问与同步
├── internal/core/              # 客户端内部模型与存储契约
├── internal/storage/           # 本地存储适配
├── internal/config/            # 配置与工作空间
├── internal/crypto/            # 加密与密钥生命周期
├── integration_test/           # 跨包客户端测试
└── scripts/                    # 回归与完整验证
```

---

## 🔧 技术栈

| 组件 | 技术 |
|------|------|
| 语言 | [Go 1.23](https://golang.org) |
| CLI 框架 | [cobra](https://github.com/spf13/cobra) |
| 存储 | [bbolt](https://github.com/etcd-io/bbolt) |
| 配置 | [ini.v1](https://gopkg.in/ini.v1) |

---

## 🌐 翻译

- [English](README.md)
- [日本語](README.ja.md)
- [Español](README.es.md)
- [Français](README.fr.md)
- [Português](README.pt.md)

---

## 🤝 贡献

欢迎贡献！你可以这样帮忙：

1. Fork 仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

对于重大更改，请先开 Issue 讨论你想做什么。

---

## 📄 许可证

本项目采用 Apache License 2.0 许可证 — 详见 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

- [cobra](https://github.com/spf13/cobra) 提供优秀的 CLI 框架
- [bbolt](https://github.com/etcd-io/bbolt) 提供可靠的嵌入式键值存储
- 开源社区

---

<div align="center">

**由讨厌寻找丢失知识的开发者用 ❤️ 打造**

</div>
