// Package server answers /healthz: alive, and how the rounds of searches are going.
package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/PlatanosVerdes/wallabot/internal/watch"
)

type Health struct {
	Version   string
	DataDir   string
	NextWatch func() time.Time
}

type payload struct {
	Status    string        `json:"status"`
	Version   string        `json:"version"`
	NextWatch string        `json:"next_watch,omitempty"`
	LastWatch *watch.Result `json:"last_watch,omitempty"`
}

func (h *Health) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// A failed search retries in minutes, so it does not decide the status.
		body := payload{Status: "ok", Version: h.Version}
		if res, ok := watch.LoadResult(h.DataDir); ok {
			body.LastWatch = &res
		}
		if h.NextWatch != nil {
			if next := h.NextWatch(); !next.IsZero() {
				body.NextWatch = next.Format(time.RFC3339)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(body)
	})
	return mux
}
