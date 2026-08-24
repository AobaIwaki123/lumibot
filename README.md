# lumibot

Lightweight Discord Bot for TimeTree public calendar briefings and notifications.

[![CI - Go](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-go.yml/badge.svg)](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-go.yml)
[![CI - Shell Scripts](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-scripts.yml/badge.svg)](https://github.com/AobaIwaki123/lumibot/actions/workflows/ci-scripts.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## Overview

`lumibot` is a lightweight Discord Bot that integrates with the [`lumitree`](https://github.com/AobaIwaki123/lumitree) API to deliver public TimeTree calendar schedules directly to your Discord servers.

- **Daily Morning Briefings**: Automatically posts today's schedule to registered channels every morning at 08:00 JST.
- **Slash Commands**: Interactive commands like `/add`, `/list`, `/remove`, and `/today` for on-demand schedule checks and subscription management.
- **Lightweight**: Pure-Go architecture utilizing `modernc.org/sqlite` (CGo-free) for zero-dependency local storage.
- **Secure**: Operates exclusively with public calendars, requiring no user account credentials.

---

## Setup Instructions

### 1. Create a Discord Application

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications).
2. Click "New Application" and give it a name.
3. Navigate to the "Bot" tab and click "Add Bot".
4. Copy the "Token" and save it securely. This will be your `DISCORD_TOKEN`.
5. Note your "Application ID" from the "General Information" tab.

### 2. Generate an Invitation Link

To invite the bot to your server, generate an OAuth2 URL:

1. In the Developer Portal, go to "OAuth2" -> "URL Generator".
2. Select the `bot` and `applications.commands` scopes.
3. Select the required bot permissions:
   - `Send Messages`
   - `Embed Links`
4. Copy the generated URL and open it in your browser to invite the bot to your server.

### 3. Run the Bot Locally

You can run the bot directly using Go:

```bash
export DISCORD_TOKEN="your-bot-token"
export LUMITREE_API_URL="http://localhost:8080" # URL of your running lumitree instance

go run ./cmd/lumibot
```

Or build the binary:

```bash
go build -o lumibot ./cmd/lumibot
./lumibot
```

---

## Environment Variables

| Variable | Required | Default Value | Description |
| :--- | :---: | :--- | :--- |
| `DISCORD_TOKEN` | Yes | - | Discord Bot Token required for authentication. |
| `LUMITREE_API_URL` | No | `https://api.lumitree.example.com` | URL of the upstream `lumitree` API server. |
| `LUMIBOT_DB_PATH` | No | `lumibot.db` | File path for the SQLite database. |

---

## Development & Local Testing (TDD)

All tasks should be performed in a dedicated Git Worktree. Do not commit directly to `main`.

```bash
# Create a new worktree for your task
./scripts/worktree.sh create feat/your-feature-name

# Run all verification checks (format, lint, test, build)
./scripts/verify-all.sh
```

---

## License

[MIT License](LICENSE) (c) 2026 Aoba Iwaki
