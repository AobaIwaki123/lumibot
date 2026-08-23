# lumibot 開発・協業規約 (AGENTS.md)

本ドキュメントは、AIエージェントおよび開発者が `lumibot` リポジトリで作業する際の共通規約を定義します。

---

## 0. リポジトリ概要

### プロジェクト情報

| 項目 | 内容 |
| :--- | :--- |
| リポジトリ | `github.com/AobaIwaki123/lumibot` |
| 言語 | Go 1.23 |
| ライセンス | MIT |
| コンテナレジストリ | `ghcr.io/aobaiwaki123/lumibot` |

**概要**: TimeTree 公開カレンダーのデイリーダイジェスト通知、スケジュール照会（スラッシュコマンド）、変更速報を配信する軽量 Discord Bot。

### ディレクトリ構成

```
lumibot/
├── cmd/
│   └── lumibot/          # エントリポイント (main.go)
├── pkg/
│   ├── bot/              # Discord セッション / イベントハンドラ / スラッシュコマンド
│   ├── client/           # lumitree API クライアント wrapper
│   ├── config/           # 設定値 (環境変数)
│   ├── cron/             # 定期配信 (デイリーダイジェスト & 差分検知)
│   └── store/            # SQLite 永続化層 (subscriptions, guild_settings)
├── k8s/                  # Kubernetes マニフェスト
├── scripts/              # 開発・検証補助スクリプト (verify-all.sh, worktree.sh)
└── .github/workflows/    # CI/CD ワークフロー
```

---

## 1. テスト駆動開発 (TDD) & 品質基準

- **テストファーストの原則**:
  - 新機能追加やバグ修正時は、まずテストコード（Unit Test）を作成し、失敗（Red）を確認してから実装（Green）に進みます。
  - リファクタリング時は常にテストが PASS し続けることを確認します。
- **CI が通るまで絶対にマージしない (CI-Green 原則)**:
  - すべての PR は、GitHub Actions CI（`golangci-lint`、`go test -race`、`go build`）が 100% PASS することを確認するまでマージしてはなりません。
- **事前ローカル検証の徹底**:
  - コミット・Push 前に必ず `./scripts/verify-all.sh` を実行し、Linter 警告 0 件・テスト通過を確認します。

---

## 2. Git & Pull Request & Worktree 運用ルール

- **`main` ブランチへの直接コミットおよび直接 Push は禁止**します。
- **すべてのタスク作業は Git Worktree（`.worktrees/<branch-name>`）を作成して行います**（`./scripts/worktree.sh create <branch-name>` を活用）。
- **PR の単一責務の原則 (Single Responsibility PR)**: 各 PR は単一の明確な目的に限定し、無関係な変更の混在を禁止します。
- **PR タイトルは具体的かつ明瞭な日本語で記述します**（例: `feat: Discord スラッシュコマンド /today の実装`, `fix: SQLite コネクションリークの修正`）。
- **エージェントによる PR の自律的なマージ・クローズは禁止**します。マージはユーザーのレビュー・判断に委ねます。

---

## 3. ドキュメント & コードスタイル

- README や設計仕様書などの公式ドキュメントでは、**原則として絵文字（emoji）を使用しません**。
- すべてのエクスポートされた関数・構造体・インターフェースには GoDoc コメントを付与し、Linter 警告を 0 に保ちます。
