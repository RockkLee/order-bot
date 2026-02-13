package httpserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/services/menusvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"strings"
	"testing"
	"time"

	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/services"
)

type fakeRepository struct {
	health map[string]string
	calls  int
}

func (f *fakeRepository) Health() (map[string]string, error) {
	f.calls++
	return f.health, nil
}

func (f *fakeRepository) Close() error {
	return nil
}

func (f *fakeRepository) Conn() *sql.DB {
	return nil
}

func (f *fakeRepository) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	if fn == nil {
		return fmt.Errorf("fakeRepository.WithTx: fn is nil")
	}
	return fn(ctx, nil)
}

type fakeUserStore struct {
	users map[string]entities.User
}

type fakeBotStore struct{}

func (f *fakeBotStore) Create(_ context.Context, _ store.Tx, _ entities.Bot) error {
	return nil
}

func (f *fakeBotStore) FindByID(_ context.Context, _ store.Tx, _ string) (entities.Bot, error) {
	return entities.Bot{}, nil
}

type fakeUserBotStore struct{}

func (f *fakeUserBotStore) Create(_ context.Context, _ store.Tx, _ entities.UserBot) error {
	return nil
}

func (f *fakeUserBotStore) FindByUserID(_ context.Context, _ store.Tx, _ string) ([]entities.UserBot, error) {
	return nil, nil
}

func (f *fakeUserStore) Create(_ context.Context, _ store.Tx, user entities.User) error {
	if _, exists := f.users[user.Email]; exists {
		return fmt.Errorf("fakeUserStore.Create: %w", store.ErrUserExists)
	}
	f.users[user.Email] = user
	return nil
}

func (f *fakeUserStore) FindByEmail(_ context.Context, _ store.Tx, email string) (entities.User, error) {
	user, exists := f.users[email]
	if !exists {
		return entities.User{}, fmt.Errorf("fakeUserStore.FindByEmail: %w", store.ErrNotFound)
	}
	return user, nil
}

func (f *fakeUserStore) FindByID(_ context.Context, _ store.Tx, id string) (entities.User, error) {
	for _, user := range f.users {
		if user.ID == id {
			return user, nil
		}
	}
	return entities.User{}, fmt.Errorf("fakeUserStore.FindByBotID: %w", store.ErrNotFound)
}

func (f *fakeUserStore) UpdateTokens(_ context.Context, _ store.Tx, id string, accessToken string, refreshToken string) error {
	for email, user := range f.users {
		if user.ID == id {
			user.AccessToken = accessToken
			user.RefreshToken = refreshToken
			f.users[email] = user
			return nil
		}
	}
	return fmt.Errorf("fakeUserStore.UpdateTokens: %w", store.ErrNotFound)
}

func TestServerDependencies(t *testing.T) {
	db := &fakeRepository{health: map[string]string{"status": "ok"}}
	authCfg := config.Auth{
		AccessSecret:    "access",
		RefreshSecret:   "refresh",
		AccessTokenTTL:  time.Minute,
		RefreshTokenTTL: time.Minute,
	}
	cfg := config.Config{
		Auth:   authCfg,
		Others: config.Others{QryCtxTimeout: time.Second},
	}
	authInitCalls := 0
	menuInitCalls := 0
	botInitCalls := 0
	serviceContainer := services.NewServices(
		func() *authsvc.Svc {
			authInitCalls++
			ctxFunc := util.NewCtxFunc(cfg.Others.QryCtxTimeout)
			return authsvc.NewSvc(nil, ctxFunc, cfg, &fakeUserStore{users: make(map[string]entities.User)})
		},
		func() *menusvc.Svc {
			menuInitCalls++
			return nil
		},
		func() *botsvc.Svc {
			botInitCalls++
			return botsvc.NewSvc(&sqldb.DB{}, nil, cfg, &fakeBotStore{}, &fakeUserBotStore{})
		},
	)
	server := NewServer(0, db, serviceContainer)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"email":"test@example.com","password":"secret","bot_name":"test-bot"}`))
	rec := httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if db.calls != 0 {
		t.Fatalf("expected db health to be unused for auth handler, got %d", db.calls)
	}
	if authInitCalls != 1 {
		t.Fatalf("expected auth init to be called once, got %d", authInitCalls)
	}
	if botInitCalls != 1 {
		t.Fatalf("expected bot init to be called once for signup handler, got %d", botInitCalls)
	}
	if menuInitCalls != 0 {
		t.Fatalf("expected menu init to be unused for auth handler, got %d", menuInitCalls)
	}

	var payload models.TokenPair
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode token response: %v", err)
	}
	if payload.AccessToken == "" || payload.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got access=%q refresh=%q", payload.AccessToken, payload.RefreshToken)
	}

	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Authentication", "Bearer "+payload.RefreshToken)
	rec = httptest.NewRecorder()
	server.Handler.ServeHTTP(rec, req)

	if db.calls != 1 {
		t.Fatalf("expected db health to be called once after health check, got %d", db.calls)
	}
	if menuInitCalls != 0 {
		t.Fatalf("expected menu init to be unused for health handler, got %d", menuInitCalls)
	}
}
