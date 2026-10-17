# Deploy

Deployed by [rpi-services](https://github.com/PlatanosVerdes/rpi-services) from the image on GHCR,
pinned to the CalVer tag that `release.yml` creates on every push to main.

Entry for its `docker-compose.yml`:

```yaml
  # --- Wallabot: the public bot that watches Wallapop searches ---
  wallabot:
    image: ghcr.io/platanosverdes/wallabot:${WALLABOT_VERSION}
    container_name: wallabot
    labels:
      prometheus.probe: "http://wallabot:8000/healthz"
    restart: unless-stopped
    profiles: [wallapop, all]
    # Hashing a thumbnail decodes it first.
    mem_limit: 128m
    memswap_limit: 128m
    dns:
      - 8.8.8.8
      - 1.1.1.1
    networks:
      - media-network
    environment:
      TZ: ${TZ:-Europe/Madrid}
      WALLA_PUSHGATEWAY: http://pushgateway:9091
      # A bot of its own: this service is the only reader of its updates.
      WALLA_TELEGRAM_TOKEN: ${WALLA_TELEGRAM_TOKEN}
      # The owner's chat, always a user of the bot.
      WALLA_TELEGRAM_CHAT: ${WALLA_TELEGRAM_CHAT}
    volumes:
      # The users, their searches and what has been seen for each of them.
      - ${APP_CONFIG_PATH}/wallabot:/data
```

Who uses the bot and what they look for, and a round without waiting for the clock:

```bash
docker exec wallabot wallabot searches
docker exec wallabot wallabot searches add --chat <id> '<address of a web search>'
docker exec wallabot wallabot watch --dry-run
```

Searches live in `searches.json` and are added, silenced and deleted from the bot. Nothing in the
watcher reads or writes a Wallapop account.
