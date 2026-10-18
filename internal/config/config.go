// Package config reads everything from the environment, so compose is the only config.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/PlatanosVerdes/wallabot/internal/places"
	"github.com/PlatanosVerdes/wallabot/internal/wallapop"
)

type Config struct {
	DataDir    string
	BaseURL    string
	WebURL     string
	Scheme     wallapop.SignScheme
	AppVersion string

	// RetryEvery is the shortest wait after a failed round.
	RetryEvery time.Duration
	Port       int

	// Pushgateway is empty outside the Pi.
	Pushgateway string

	TelegramToken string
	// TelegramChat is the owner's: always let in, and the default for "searches add".
	TelegramChat string
	// The bot is public and every search is requests from the Pi, hence the caps.
	MaxSearches    int
	MaxUsers       int
	WatchMin       time.Duration
	WatchMax       time.Duration
	WatchMaxAge    time.Duration
	WatchMaxAlerts int
	WatchPhotos    int
	// SearchPages is read on a deep round. One page is 40 listings.
	SearchPages int
	// DeepEvery keeps the prices further down a search watched.
	DeepEvery time.Duration
	// WatchDrop is the share of its lowest price a listing must shed to be announced; 0 is off.
	WatchDrop     float64
	WatchMinPause time.Duration
	WatchMaxPause time.Duration
	SeenTTL       time.Duration
	// PlacesURL is the geocoder for "en <town>"; empty turns it off.
	PlacesURL string

	LogJSON bool
	// StatsAddr exposes /bot-stats; empty turns it off. It lists names: never publish it.
	StatsAddr string
	BotName   string
}

func Load() (Config, error) {
	cfg := Config{
		DataDir:     env("WALLA_DATA_DIR", "./data"),
		BaseURL:     env("WALLA_BASE_URL", wallapop.DefaultBaseURL),
		WebURL:      env("WALLA_WEB_URL", wallapop.DefaultWebURL),
		Scheme:      wallapop.SignScheme(env("WALLA_SIGN_SCHEME", string(wallapop.SchemeNone))),
		AppVersion:  env("WALLA_APP_VERSION", wallapop.DefaultAppVersion),
		RetryEvery:  duration("WALLA_RETRY_EVERY", 15*time.Minute),
		Port:        number("WALLA_PORT", 8000),
		Pushgateway: env("WALLA_PUSHGATEWAY", ""),

		TelegramToken:  env("WALLA_TELEGRAM_TOKEN", ""),
		TelegramChat:   env("WALLA_TELEGRAM_CHAT", ""),
		MaxSearches:    number("WALLA_MAX_SEARCHES", 4),
		MaxUsers:       number("WALLA_MAX_USERS", 20),
		WatchMin:       duration("WALLA_WATCH_MIN", 5*time.Minute),
		WatchMax:       duration("WALLA_WATCH_MAX", 15*time.Minute),
		WatchMaxAge:    duration("WALLA_WATCH_MAX_AGE", 24*time.Hour),
		WatchMaxAlerts: number("WALLA_WATCH_MAX_ALERTS", 10),
		WatchPhotos:    number("WALLA_WATCH_PHOTOS", 2),
		SearchPages:    number("WALLA_SEARCH_PAGES", 3),
		DeepEvery:      duration("WALLA_DEEP_EVERY", time.Hour),
		WatchDrop:      percent("WALLA_WATCH_DROP", 5) / 100,
		WatchMinPause:  duration("WALLA_WATCH_MIN_PAUSE", 3*time.Second),
		WatchMaxPause:  duration("WALLA_WATCH_MAX_PAUSE", 15*time.Second),
		SeenTTL:        duration("WALLA_SEEN_TTL", 30*24*time.Hour),
		PlacesURL:      env("WALLA_PLACES_URL", places.DefaultURL),

		LogJSON:   env("WALLA_LOG_JSON", "") == "1",
		StatsAddr: env("WALLA_STATS_ADDR", ""),
		BotName:   env("WALLA_BOT_NAME", ""),
	}

	// A zero wait turns the ticker into a spin loop.
	if cfg.RetryEvery < time.Minute {
		return cfg, fmt.Errorf("WALLA_RETRY_EVERY (%s) is under a minute, which would spin", cfg.RetryEvery)
	}
	if cfg.WatchMin < time.Minute {
		return cfg, fmt.Errorf("WALLA_WATCH_MIN (%s) is under a minute, which would spin", cfg.WatchMin)
	}
	if cfg.WatchMax < cfg.WatchMin {
		return cfg, fmt.Errorf("WALLA_WATCH_MAX (%s) is below WALLA_WATCH_MIN (%s)", cfg.WatchMax, cfg.WatchMin)
	}
	if cfg.WatchMaxPause < cfg.WatchMinPause {
		return cfg, fmt.Errorf("WALLA_WATCH_MAX_PAUSE (%s) is below WALLA_WATCH_MIN_PAUSE (%s)", cfg.WatchMaxPause, cfg.WatchMinPause)
	}
	if !wallapop.ValidScheme(cfg.Scheme) {
		return cfg, fmt.Errorf("WALLA_SIGN_SCHEME %q is not one of none, pipe, legacy", cfg.Scheme)
	}

	// A change on their side should be a redeploy, not a rebuild.
	if v := os.Getenv("WALLA_PATH_SEARCH"); v != "" {
		wallapop.PathSearch = v
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func number(key string, fallback int) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}

func percent(key string, fallback float64) float64 {
	v, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil || v < 0 {
		return fallback
	}
	return v
}

func duration(key string, fallback time.Duration) time.Duration {
	v, err := time.ParseDuration(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return v
}
