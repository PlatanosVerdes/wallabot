package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PlatanosVerdes/wallabot/internal/watch"
)

func get(t *testing.T, h *Health) (int, map[string]any) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("the body is not JSON: %v (%s)", err, rec.Body.String())
	}
	return rec.Code, body
}

func TestHealthAnswersBeforeTheFirstRound(t *testing.T) {
	code, body := get(t, &Health{Version: "test", DataDir: t.TempDir()})
	if code != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("answered %d with %v", code, body)
	}
	if _, ok := body["last_watch"]; ok {
		t.Errorf("a round was reported before any ran: %v", body)
	}
}

func TestHealthReportsTheRounds(t *testing.T) {
	dir := t.TempDir()
	if err := watch.SaveResult(dir, watch.Result{Users: 3, Watched: 7}); err != nil {
		t.Fatal(err)
	}
	next := time.Date(2026, 10, 18, 1, 0, 0, 0, time.UTC)

	_, body := get(t, &Health{DataDir: dir, NextWatch: func() time.Time { return next }})
	last, ok := body["last_watch"].(map[string]any)
	if !ok || last["users"] != 3.0 || last["watched"] != 7.0 {
		t.Errorf("last_watch came back as %v", body["last_watch"])
	}
	if body["next_watch"] != "2026-10-18T01:00:00Z" {
		t.Errorf("next_watch is %v", body["next_watch"])
	}
}
