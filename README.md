# wallabot

**[@WuallaBot](https://t.me/WuallaBot)**: a public Telegram bot that tells you when something new appears in your Wallapop searches.

Runs on the Raspberry from [rpi-services](https://github.com/PlatanosVerdes/rpi-services). Go, no dependencies, CalVer tag on every push to main. Reactivating my own listings is a separate service, [wallapop-catalogue](https://github.com/PlatanosVerdes/wallapop-catalogue).

## The bot

1. `/start` to join.
2. Write what you are looking for and it is saved (a greeting gets a joke instead). A price and a radius are picked out of the text:

   | Written | Searches |
   | :--- | :--- |
   | `iphone 13` | "iphone 13", all of Spain |
   | `kallax hasta 40` · `máx 40` · `menos de 40` | up to 40 € |
   | `moto desde 1.500` · `más de 1500` | from 1500 € |
   | `bici 100-300` · `entre 100 y 300` | 100 to 300 € |
   | `sofá en Sant Cugat` | 30 km around Sant Cugat |
   | `sofá en Sant Cugat a 10 km` | 10 km around Sant Cugat |

   The town after the last `en` is looked up on OpenStreetMap; when it is no town (`funda en piel`), it stays part of the text.
3. Optional: send a location (📎 → Location, any point on the map) to move the latest search there, 30 km unless the text said otherwise.

A listing shared from the app (`wallapop.com/item/...`) becomes a search for things like it: the first three words of its title, in its category, up to 20% dearer.

For filters the text cannot say (category, condition, brand and model), make the search on es.wallapop.com in a phone browser (typed in the address bar, so it does not jump to the app) and paste its address. Whatever is written next to it is the name, e.g. `coches top https://es.wallapop.com/search?...`

| Command | What it does |
| :--- | :--- |
| `/nueva <search>` | Same as writing it, or pasting the address |
| `/busquedas` | Your searches: 🔔 silence, ✏️ rename, 🗑 delete, and when the last round ended |
| `/ahora` | Check one search, or all, right now |
| `/baja` | Leave and delete your searches |
| `/ayuda` | Help |

Each listing arrives with its photo, price, town, the details Wallapop shows for it (year and km for a motorbike, condition, colour and material for furniture, whether it ships), and two buttons: open it, or silence that search. Details already in the title are left out, and a listing whose details cannot be read is sent without them.

Each chat only sees its own searches. Up to `WALLA_MAX_USERS` chats, `WALLA_MAX_SEARCHES` searches each.

## How it works

- A round runs every 5 to 15 minutes, at random.
- Each search reads **one page** (40 listings, newest first). Every hour, and on a new search's first round, it reads 3.
- Two chats with the **same search** share one request.
- The **first round** of a new search is silent: those listings were already there.
- **Duplicates** (the same vehicle posted by many dealer accounts) are detected by photo: a 64-bit image hash plus a similar price.
- A known listing that gets **5% cheaper** than its lowest price is announced again.
- Searches need no Wallapop account: every request is anonymous.

## Search API notes

The bot uses Wallapop's private `GET /api/v3/search`, which needs no login.

| Parameter | Note |
| :--- | :--- |
| `source` | Required, any value (400 without it) |
| `distance_in_km` | The radius. Only works with `latitude` and `longitude`; without them, all of Spain |
| `distance` | Ignored by the API |
| `order_by` | Always forced to `newest` |

## Data

In `WALLA_DATA_DIR`:

| File | Holds |
| :--- | :--- |
| `searches.json` | Every chat and its searches |
| `users/<chat>/seen.json` | What that chat has already been told |
| `last_watch.json` | The last round |

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
| `WALLA_STATS_ADDR` | – | An address such as `:9180` where [muninn](https://github.com/PlatanosVerdes/muninn) reads who uses the bot, its traffic and what the gate turned away (`/bot-stats`, kept in `bot-stats.json` in the data folder). It lists names: never publish this port |
| `WALLA_BOT_NAME` | – | The bot's username, as muninn shows it |

Endpoints can be overridden without a rebuild: `WALLA_BASE_URL`, `WALLA_WEB_URL`, `WALLA_PATH_SEARCH`, `WALLA_SIGN_SCHEME`, `WALLA_APP_VERSION`.

## Health

`/healthz` always answers 200 while the process runs, with the last round and the next one. A failed search retries in minutes, so it does not change the `status`.

Metrics go to the Pushgateway as `wallapop_watch_*` gauges, under the job `wallapop-watch`.
