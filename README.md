# lumibot

Lightweight Discord Bot for TimeTree public calendar briefings and notifications.

[![CI - Go](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-go.yml/badge.svg)](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-go.yml)
[![CI - Shell Scripts](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-scripts.yml/badge.svg)](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-scripts.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 概要

`lumibot` は、[`lumitree`](https://github.com/AobaIwaki123/lumitree) API と連携し、TimeTree 公開カレンダーのスケジュールを Discord サーバーへ自動配信・照会する軽量 Bot です。

- **朝の定期ダイジェスト配信**: 毎朝決まった時間に登録グループの「今日の予定」を Embed 形式で自動投稿。
- **スラッシュコマンド照会**: `/today`, `/tomorrow`, `/list`, `/add`, `/remove` による即時確認。
- **超軽量・省リソース**: Pure-Go (CGo-free SQLite) によるシングルバイナリ設計（数十 MB 以下で常駐）。
- **ゼロ・クレデンシャル**: ユーザーのアカウント情報（ID・パスワード）不要で安全に運用可能。

---

## 開発とローカル検証 (TDD)

```bash
# 全テスト & リント & ビルド検証 (Push 前に必ず実行)
./scripts/verify-all.sh

# Git Worktree を作成してタスク作業を開始
./scripts/worktree.sh create feat/slash-commands
```

---

## 環境変数設定

| 環境変数 | 必須 | デフォルト値 | 説明 |
| :--- | :---: | :--- | :--- |
| `DISCORD_TOKEN` | ◯ | - | Discord Bot Token |
| `DISCORD_APP_ID` | - | - | Discord Application ID (グローバルコマンド登録用) |
| `LUMITREE_API_URL` | - | `https://lumitree.aooba.net` | 接続先 lumitree API サーバー URL |
| `SQLITE_DB_PATH` | - | `lumibot.db` | SQLite データベースファイルのパス |
| `LOG_LEVEL` | - | `info` | ログレベル (`debug`, `info`, `warn`, `error`) |

---

## ライセンス

[MIT License](LICENSE) © 2026 Aoba Iwaki
