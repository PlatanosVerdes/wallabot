# How it works

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

