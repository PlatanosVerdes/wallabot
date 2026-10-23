# Operations

## CLI

```bash
wallabot serve                          # the bot, the rounds, /healthz
wallabot watch --dry-run                # one round, print instead of send
wallabot searches                       # every chat and its searches
wallabot searches add --chat <id> '<address or text>'
```

## Configuration

| Variable | Default | |
| :--- | :--- | :--- |
| `WALLA_TELEGRAM_TOKEN` | – | Bot token |
| `WALLA_TELEGRAM_CHAT` | – | Owner's chat |
| `WALLA_MAX_USERS` | `20` | Chats that can join |
| `WALLA_MAX_SEARCHES` | `4` | Searches per chat |
| `WALLA_WATCH_MIN` / `MAX` | `5m` / `15m` | Wait between rounds |
| `WALLA_WATCH_MIN_PAUSE` / `MAX_PAUSE` | `3s` / `15s` | Pause between requests |
| `WALLA_SEARCH_PAGES` | `3` | Pages on a deep round |
| `WALLA_DEEP_EVERY` | `1h` | How often a round is deep |
| `WALLA_WATCH_MAX_AGE` | `24h` | Older listings are not announced |
| `WALLA_WATCH_MAX_ALERTS` | `10` | Messages per chat per round |
| `WALLA_WATCH_PHOTOS` | `2` | Photos hashed per listing |
| `WALLA_WATCH_DROP` | `5` | % drop to announce a cheaper listing (`0` = off) |
| `WALLA_SEEN_TTL` | `720h` | How long a listing is remembered |
| `WALLA_PLACES_URL` | Nominatim | Geocoder for `en <town>` (empty = off) |
| `WALLA_RETRY_EVERY` | `15m` | Shortest wait after a failed round |
| `WALLA_DATA_DIR` | `./data` | Data folder |
| `WALLA_PORT` | `8000` | `/healthz` port |
| `WALLA_PUSHGATEWAY` | – | Where metrics go |
| `WALLA_LOG_JSON` | – | `1` for JSON logs |
| `WALLA_STATS_ADDR` | – | An address such as `:9180` that exposes the bot's stats at `/bot-stats`: who uses it, its traffic and what the gate turned away (`gonana/botstats`, kept in `bot-stats.json` in the data folder). It lists names: never publish it |
| `WALLA_BOT_NAME` | – | The bot's username, named in its stats |

Endpoints can be overridden without a rebuild: `WALLA_BASE_URL`, `WALLA_WEB_URL`, `WALLA_PATH_SEARCH`, `WALLA_SIGN_SCHEME`, `WALLA_APP_VERSION`.

## Health

`/healthz` always answers 200 while the process runs, with the last round and the next one. A failed search retries in minutes, so it does not change the `status`.

Metrics go to the Pushgateway as `wallapop_watch_*` gauges, under the job `wallapop-watch`.
