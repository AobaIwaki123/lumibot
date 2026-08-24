# lumibot

Lightweight Discord Bot for TimeTree public calendar briefings and notifications.

[![CI - Go](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-go.yml/badge.svg)](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-go.yml)
[![CI - Shell Scripts](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-scripts.yml/badge.svg)](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-scripts.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 概要

`lumibot` は、[`lumitree`](https://github.com/AobaIwaki123/lumitree) API と連携し、TimeTree 公開カレンダーのスケジュールを Discord サーバーへ自動配信・照会する軽量 Bot です。

- **朝の定期ダイジェスト配信**: 毎朝決まった時間に登録グループの「今日の予定」を Embed 形式で自動投稿。
- **スラッシュコマンド照会**: `/today`, `/list`, `/add`, `/remove` による即時確認と購読管理。
- **超軽量・省リソース**: Pure-Go (CGo-free SQLite) によるシングルバイナリ設計（数十 MB 以下で常駐）。
- **ゼロ・クレデンシャル**: ユーザーのアカウント情報（ID・パスワード）不要で安全に運用可能。

---

## セットアップ手順

### 1. Discord アプリケーションの作成

1. [Discord Developer Portal](https://discord.com/developers/applications) にアクセスします。
2. 「New Application」をクリックし、Bot の名前を入力します。
3. 左側のメニューから「Bot」タブを開き、「Add Bot」をクリックします。
4. 「Token」をコピーして安全な場所に保存してください（これが `DISCORD_TOKEN` になります）。
5. 「General Information」タブから「Application ID」をメモしておきます。

### 2. Bot の招待リンク作成

サーバーへ Bot を招待するため、以下の手順で URL を生成します。

1. Developer Portal の「OAuth2」->「URL Generator」を開きます。
2. Scopes で `bot` と `applications.commands` を選択します。
3. Bot Permissions で以下の権限を選択します。
   - `Send Messages`
   - `Embed Links`
4. 下部に生成された URL をブラウザで開き、自身のサーバーへ Bot を招待します。

### 3. ローカルでの起動方法

環境変数を設定し、以下のコマンドで起動します。

```bash
export DISCORD_TOKEN="your-bot-token"
export LUMITREE_API_URL="http://localhost:8080" # ローカル起動している lumitree のURL

go run ./cmd/lumibot
```

---

## 環境変数設定

| 環境変数 | 必須 | デフォルト値 | 説明 |
| :--- | :---: | :--- | :--- |
| `DISCORD_TOKEN` | ◯ | - | Discord Bot Token |
| `LUMITREE_API_URL` | - | `https://api.lumitree.example.com` | 接続先 lumitree API サーバー URL |
| `LUMIBOT_DB_PATH` | - | `lumibot.db` | SQLite データベースファイルのパス |

---

## 開発とローカル検証 (TDD)

`main` への直接コミットは禁止されています。必ず Git Worktree を作成して作業してください。

```bash
# タスク用の Git Worktree を作成
./scripts/worktree.sh create feat/your-feature-name

# 全テスト & リント & ビルド検証 (Push 前に必ず実行)
./scripts/verify-all.sh
```

---

## ライセンス

[MIT License](LICENSE) (c) 2026 Aoba Iwaki
