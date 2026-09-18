<div align="center">

# ttl

### あなたの個人ナレッジアーカイブ


[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md) | [Español](README.es.md) | [Français](README.fr.md) | [Português](README.pt.md)

---


</div>

---

## 📖 ストーリー

すべての開発者が経験したことがあります：

> 「あの先月使った Docker コマンド、なんだったっけ？」
> 「同僚が共有してくれた設定ファイル、どこにあったっけ？」
> 「関連記事を読んだはずだけど…見つからない」

私たちは知識を至る所に保存しています —— ブラウザのタブ、Slack メッセージ、メール、ブックマーク、メモアプリ。本当に必要なとき、無限のタブやチャット履歴をスクロールして時間を浪費しています。

**これは私の悩みでもありました。**

だから **ttl** を作りました。

名前は "Time to Live" に由来しますが、少し意味が違います。期限切れのことではなく、知識を**永遠に保存**することです。

- すべてを一箇所にキーバリュー形式で保存
- タグを付けて整理
- キーワードで即座に検索

もう古いメールを探したり、チャット履歴をスクロールする必要はありません。`ttl get <キーワード>` だけですぐに見つかります。

**ttl はあなたの個人のナレッジアーカイブ —— 必要なものを、必要なときに。**

---

## ✨ 機能

| 機能 | 説明 |
|------|------|
| 🗄️ **ローカル KV ストレージ** | 高速、設定不要の組み込みデータベース (bbolt) |
| 🏷️ **タグシステム** | 柔軟で検索可能なタグでリソースを整理 |
| 🔍 **ファジー検索** | キーとタグをまたいで即座に検索 |
| 📝 **作業ログ** | 日次作業を記録・絞り込み |
| ☁️ **クラウド同期** | 別プロジェクトで運用される TTL バックエンドに接続 |
| 🚀 **スマートオープン** | システムデフォルトプログラムで URL とファイルを開く |
| 📤 **エクスポート** | JSON または CSV 形式でエクスポート |

---

## 🚀 クイックスタート

### インストール

#### Linux / macOS

```bash
# GitHub releases からインストール
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"

# またはソースからビルド
go build -o ttl ./cmd/ttl
sudo mv ttl /usr/local/bin/
```

#### Windows

```powershell
# GitHub releases からインストール
irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

#### カスタムダウンロード URL

社内ネットワークやカスタムミラー用：

```bash
# Linux/macOS
TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-linux-amd64" /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.sh)"
```

```powershell
# Windows
$env:TTL_DOWNLOAD_URL="https://your-mirror.com/ttl-cli-v1.0.0-windows-amd64.zip"; irm https://raw.githubusercontent.com/ZHANGSHUNLIN/TTL-CLI/main/install.ps1 | iex
```

### 基本的な使用方法

```bash
# リソースを追加
ttl add my-link https://example.com

# タグを付けて追加
ttl add docker-cmd "docker run -d -p 8080:80 nginx"
ttl tag docker-cmd dev ops

# リソースを検索
ttl get docker

# ブラウザで開く
ttl open my-link

# 削除
ttl del old-key
```

---


## 📝 作業ログ

日次作業を記録・絞り込みます。

```bash
# ログを書く
ttl log write "ユーザーモジュールのリファクタリング完了" --tags "プロジェクトA,開発"

# ログを表示
ttl log list                    # 今日のログ
ttl log list --range week       # 今週
ttl log list --range month      # 今月

```

---


## ☁️ クラウドサービスと同期

バックエンドは別のサービスプロジェクトで保守・運用されます。このリポジトリには `ttl` クライアントと HTTP アダプターだけが含まれ、サーバー実行ファイルはビルドも配布もしません。

### データを同期

```bash
# リモートサーバーを設定
ttl config
# server セクションを編集してエンドポイントと API キーを入力

# ローカルとリモートを同期
ttl sync
```

**アーキテクチャ：**
- ユーザーごとの分離データベースを持つマルチテナント設計
- API Key 認証
- プログラム的アクセス用の REST API

---

## ⚙️ 設定

設定ファイル：`~/.ttl/ttl.ini`

```ini
[default]
db_path = ~/.ttl/data.db

[server]
endpoint  = https://your-server.com
api_key   = your-user-api-key
```

```bash
# 現在の設定を表示
ttl config

```

---

## 📤 データのエクスポート

```bash
# JSON でエクスポート
ttl export --format json

# CSV でエクスポート
ttl export --format csv

# 指定ファイルにエクスポート
ttl export --format json --output backup.json
```

---

## 🏗️ プロジェクト構造

```
ttl-cli/
├── cmd/ttl/                    # クライアントエントリーポイント
├── internal/client/            # CLI、TUI、リモートアクセス、同期
├── internal/core/              # クライアント内部モデルと契約
├── internal/storage/           # ローカルストレージアダプター
├── internal/config/            # 設定とワークスペース
├── internal/crypto/            # 暗号化とキー管理
├── integration_test/           # クライアント統合テスト
└── scripts/                    # 回帰テストと完全検証
```

---

## 🔧 技術スタック

| コンポーネント | 技術 |
|------|------|
| 言語 | [Go 1.23](https://golang.org) |
| CLI フレームワーク | [cobra](https://github.com/spf13/cobra) |
| ストレージ | [bbolt](https://github.com/etcd-io/bbolt) |
| 設定 | [ini.v1](https://gopkg.in/ini.v1) |

---

## 🌐 翻訳

- [English](README.md)
- [简体中文](README.zh-CN.md)
- [Español](README.es.md)
- [Français](README.fr.md)
- [Português](README.pt.md)

---

## 🤝 貢献

貢献を歓迎します！以下の方法で協力できます：

1. リポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. Pull Request を開く

大きな変更については、まず何を変更したいか議論するために Issue を開いてください。

---

## 📄 ライセンス

このプロジェクトは Apache License 2.0 の下でライセンスされています — 詳細は [LICENSE](LICENSE) ファイルをご覧ください。

---

## 🙏 謝辞

- 素晴らしい CLI フレームワーク [cobra](https://github.com/spf13/cobra)
- 信頼性の高い組み込みキーバリューストレージ [bbolt](https://github.com/etcd-io/bbolt)
- オープンソースコミュニティ

---

<div align="center">

**失われた知識を探すのが嫌いな開発者が ❤️ で作成**

</div>
