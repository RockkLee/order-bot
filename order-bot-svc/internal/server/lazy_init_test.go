package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-bot-svc/internal/database"
)

type fakeDB struct {
	health map[string]string
}

func (f *fakeDB) Health() map[string]string {
	return f.health
}

func (f *fakeDB) Close() error {
	return nil
}

func TestServerLazyDBInit(t *testing.T) {
	called := 0
	db := &fakeDB{health: map[string]string{"status": "ok"}}
	server := NewServer(0, func() database.Service {
		called++
		return db
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if called != 0 {
		t.Fatalf("expected db factory to be unused for hello handler, got %d", called)
	}

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	rec = httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if called != 1 {
		t.Fatalf("expected db factory to be called once, got %d", called)
	}

	var payload map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}
	if payload["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", payload["status"])
	}
}
