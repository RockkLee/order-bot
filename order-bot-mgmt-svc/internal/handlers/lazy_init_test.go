package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/postgres"
	postgresuser "order-bot-mgmt-svc/internal/postgres/user"
	"order-bot-mgmt-svc/internal/services"
)

type fakeRepository struct {
	health map[string]string
}

func (f *fakeRepository) Health() map[string]string {
	return f.health
}

func (f *fakeRepository) Close() error {
	return nil
}

func (f *fakeRepository) Conn() *sql.DB {
	return nil
}

type fakeUserStore struct {
	users map[string]models.User
}

func (f *fakeUserStore) Create(_ context.Context, user models.User) error {
	if _, exists := f.users[user.Email]; exists {
		return postgresuser.ErrUserExists
	}
	f.users[user.Email] = user
	return nil
}

func (f *fakeUserStore) FindByEmail(_ context.Context, email string) (models.User, error) {
	user, exists := f.users[email]
	if !exists {
		return models.User{}, postgresuser.ErrNotFound
	}
	return user, nil
}

func TestServerLazyServicesInit(t *testing.T) {
	dbCalled := 0
	authCalled := 0
	db := &fakeRepository{health: map[string]string{"status": "ok"}}
	server := NewServer(
		0,
		func() postgres.Service {
			dbCalled++
			return db
		},
		func() *services.Service {
			authCalled++
			return services.NewService(&fakeUserStore{users: make(map[string]models.User)})
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"email":"test@example.com","password":"secret"}`))
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if authCalled != 1 {
		t.Fatalf("expected auth factory to be called once, got %d", authCalled)
	}
	if dbCalled != 0 {
		t.Fatalf("expected db factory to be unused for auth handler, got %d", dbCalled)
	}

	var payload models.TokenPair
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode token response: %v", err)
	}
	if payload.AccessToken == "" || payload.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got access=%q refresh=%q", payload.AccessToken, payload.RefreshToken)
	}

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	rec = httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if dbCalled != 1 {
		t.Fatalf("expected db factory to be called once after health check, got %d", dbCalled)
	}
}
