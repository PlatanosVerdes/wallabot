<p align="center">
  <img src="assets/icon.jpg" width="96" alt="">
</p>

<h1 align="center">wallabot</h1>

<p align="center">
  <a href="https://github.com/PlatanosVerdes/wallabot/tags"><img src="https://img.shields.io/github/v/tag/PlatanosVerdes/wallabot?label=release&color=2f7d3a" alt="Latest release"></a>
  <a href="https://t.me/WuallaBot"><img src="https://img.shields.io/badge/telegram-@WuallaBot-26a5e4?logo=telegram&logoColor=white" alt="Telegram @WuallaBot"></a>
</p>

**[@WuallaBot](https://t.me/WuallaBot)** is a public Telegram bot that tells you when something new appears in your Wallapop searches. Write what you are looking for, optionally with a price and a radius, and each new listing arrives with its photo, price and town.

It checks every few minutes without a Wallapop account. Go, no dependencies, a CalVer tag on every push to main. It runs on the Raspberry from [rpi-services](https://github.com/PlatanosVerdes/rpi-services). Reactivating my own listings is a separate service, [wallapop-catalogue](https://github.com/PlatanosVerdes/wallapop-catalogue).

## Running it

```bash
export WALLA_TELEGRAM_TOKEN=<bot token>
export WALLA_TELEGRAM_CHAT=<owner chat id>
wallabot serve                          # the bot, the rounds, /healthz
```

`/healthz` answers on port 8000. In production it is the image `ghcr.io/platanosverdes/wallabot`, pinned by CalVer tag; the compose entry is in [deploy/docker.md](deploy/docker.md).

| Variable | Default | |
| :--- | :--- | :--- |
| `WALLA_TELEGRAM_TOKEN` | – | Bot token |
| `WALLA_TELEGRAM_CHAT` | – | Owner's chat |
| `WALLA_MAX_USERS` | `20` | Chats that can join |
| `WALLA_MAX_SEARCHES` | `4` | Searches per chat |
| `WALLA_DATA_DIR` | `./data` | Data folder |

## Docs

| | |
| :--- | :--- |
| [Using the bot](docs/usage.md) | How to write a search, the commands, what each alert contains |
| [How it works](docs/internals.md) | Rounds, duplicates, price drops, the Wallapop search API, stored data |
| [Operations](docs/operations.md) | The CLI, every setting, health and metrics |
| [Deploy](deploy/docker.md) | The compose entry used on the Pi |
