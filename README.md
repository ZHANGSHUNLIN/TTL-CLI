<div align="center">

# ttl

### Your Personal Knowledge Archive

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[简体中文](README.zh-CN.md) | [日本語](README.ja.md) | [Español](README.es.md) | [Français](README.fr.md) | [Português](README.pt.md)

---

*A lightweight CLI tool for personal data management. Store anything as key-value pairs and search instantly.*

</div>

---

## 📖 The Story

Every developer has been there:

> "Wait, what was that Docker command I used last month?"
> "Where did my colleague share that config file?"
> "I know I read an article about this... but can't find it anywhere."

We store knowledge everywhere — browser tabs, Slack messages, emails, bookmark folders, notes apps. When we actually need it, we waste time digging through endless tabs and scrolling through chat history.

**This was my pain too.**

So I built **ttl**.

The name comes from "Time to Live" — but with a different meaning. Instead of expiring, it's about giving your knowledge a **time to live** forever.

- Store everything in one place as key-value pairs
- Tag it for easy organization
- Search instantly by keyword

No more searching through old emails or scrolling through Slack history. Just `ttl get <keyword>` and you have it.

**ttl is your personal knowledge archive — everything you need, when you need it.**

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🗄️ **Local KV Storage** | Fast, zero-config embedded database (bbolt) |
| 🏷️ **Tag System** | Organize resources with flexible, searchable tags |
| 🔍 **Fuzzy Search** | Find what you need instantly across keys and tags |
| 📝 **Work Log** | Track and filter daily work |
| ☁️ **Cloud Sync** | Self-host multi-tenant server with per-user data isolation |
| 🚀 **Smart Open** | Open URLs and files with system default programs |
| 📤 **Export** | Export data as JSON or CSV |

---

## 🚀 Quick Start

### Installation

#### Linux / macOS

```bash
# Install from GitHub releases
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"

# Or build from source
go build -o ttl ./cmd/ttl
sudo mv ttl /usr/local/bin/
```

#### Windows

```powershell
# Install from GitHub releases
irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

#### Custom Download URL

For internal networks or custom mirrors:

```bash
# Linux/macOS
TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-linux-amd64" /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"
```

```powershell
# Windows
$env:TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-windows-amd64.zip"; irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```
```

### Basic Usage

```bash
# Add a resource
ttl add my-link https://example.com

# Add with tags
ttl add docker-cmd "docker run -d -p 8080:80 nginx"
ttl tag docker-cmd dev ops

# Search resources
ttl get docker

# Open in browser
ttl open my-link

# Delete
ttl del old-key
```

---

## 📝 Work Log

Track and filter your daily work.

```bash
# Write a log entry
ttl log write "Finished user module refactoring" --tags "projectA,dev"

# View logs
ttl log list                    # Today's logs
ttl log list --range week       # This week
ttl log list --range month      # This month

```

---

## ☁️ Cloud Server & Sync

### Start Your Server

```bash
# Create a user
ttl server user add alice

# Start multi-tenant server
ttl server start --port 8080
```

### Sync Your Data

```bash
# Configure remote server
ttl config
# Edit the server section with your endpoint and API key

# Sync local <-> remote
ttl sync
```

**Architecture:**
- Multi-tenant design with per-user isolated databases
- API Key authentication
- REST API for programmatic access

---

## ⚙️ Configuration

Config file: `~/.ttl/ttl.ini`

```ini
[default]
db_path = ~/.ttl/data.db

[server]
endpoint  = https://your-server.com
api_key   = your-user-api-key
```

```bash
# View current config
ttl config

```

---

## 📤 Export Your Data

```bash
# Export as JSON
ttl export --format json

# Export as CSV
ttl export --format csv

# Export to specific file
ttl export --format json --output backup.json
```

---

## 🏗️ Project Structure

```
ttl-cli/
├── main.go                 # Compatible root build for the ttl client
├── cmd/ttl/                # Canonical local client entry point
├── cmd/ttl-server/         # Standalone backend server entry point
├── internal/client/cli/    # Client command tree, sync, and migration
├── internal/server/        # Backend API, tenant data, and server commands
├── command/                # User-facing CLI commands
├── db/                     # Storage interface and bbolt/SQLite/cloud backends
├── sync/                   # Local/remote diff and push/pull execution
├── conf/                   # INI configuration and workspace management
├── crypto/                 # Data encryption and key lifecycle
├── i18n/                   # Localization loader and locale resources
├── models/                 # Shared persisted and API-facing types
├── util/                   # Small shared helpers
├── integration_test/       # Cross-package and server/sync scenarios
├── scripts/                # CLI regression and full verification scripts
└── docs/                   # Engineering workflow and decision records
```

For startup flow, package responsibilities, data flow, and common change
locations, see [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md).

The local client and backend service can now be built independently:

```bash
go build -o ttl .
go build -o ttl ./cmd/ttl
go build -o ttl-server ./cmd/ttl-server
```

`ttl server ...` remains available during migration. New deployments can use
`ttl-server serve ...` and `ttl-server user ...` directly.

---

## 🔧 Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | [Go 1.23](https://golang.org) |
| CLI Framework | [cobra](https://github.com/spf13/cobra) |
| Storage | [bbolt](https://github.com/etcd-io/bbolt) |
| Configuration | [ini.v1](https://gopkg.in/ini.v1) |

---

## 🌐 Translations

- [简体中文](README.zh-CN.md)
- [日本語](README.ja.md)
- [Español](README.es.md)
- [Français](README.fr.md)
- [Português](README.pt.md)

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

For major changes, please open an issue first to discuss what you'd like to change.

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [cobra](https://github.com/spf13/cobra) for the excellent CLI framework
- [bbolt](https://github.com/etcd-io/bbolt) for the reliable embedded key-value store
- The open-source community

---

<div align="center">

**Made with ❤️ by developers who hate searching for lost knowledge**

</div>
