package config

import (
	"testing"
	"time"
)

// A field declared but not filled in Load reads as zero, and a zero wait spins.
func TestDefaultsAreUsable(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range []struct {
		name  string
		value time.Duration
	}{
		{"RetryEvery", cfg.RetryEvery},
		{"WatchMin", cfg.WatchMin},
		{"WatchMax", cfg.WatchMax},
		{"DeepEvery", cfg.DeepEvery},
	} {
		if d.value <= 0 {
			t.Errorf("%s defaulted to %s", d.name, d.value)
		}
	}
}

func TestSpinningWaitsAreRefused(t *testing.T) {
	for _, key := range []string{"WALLA_RETRY_EVERY", "WALLA_WATCH_MIN"} {
		t.Setenv(key, "0s")
		if _, err := Load(); err == nil {
			t.Errorf("%s=0s was accepted", key)
		}
	}
}

func TestAnUnknownSignSchemeIsRefused(t *testing.T) {
	t.Setenv("WALLA_SIGN_SCHEME", "nope")
	if _, err := Load(); err == nil {
		t.Fatal("an unknown scheme was accepted")
	}
}
