package authsvc

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/sqldbold/pqsqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"testing"
	"time"
)

type fakeUserStore struct {
	users map[string]entities.User
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
	return entities.User{}, fmt.Errorf("fakeUserStore.FindByID: %w", store.ErrNotFound)
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

func TestSvcSignupAndLogin(t *testing.T) {
	cfg := config.Config{
		Auth: config.Auth{
			AccessSecret:    "access",
			RefreshSecret:   "refresh",
			AccessTokenTTL:  time.Minute,
			RefreshTokenTTL: time.Minute,
		},
		Others: config.Others{QryCtxTimeout: time.Second},
	}
	userStore := &fakeUserStore{users: make(map[string]entities.User)}
	ctxFunc := util.NewCtxFunc(cfg.Others.QryCtxTimeout)
	svc := NewSvc(&pqsqldb.DB{}, ctxFunc, cfg, userStore)

	ctx := context.Background()
	tokenPair, userID, err := svc.Signup(ctx, nil, "test@example.com", "secret")
	if err != nil {
		t.Fatalf("expected signup to succeed, got error: %v", err)
	}
	if tokenPair.AccessToken == "" || tokenPair.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens after signup")
	}

	user, exists := userStore.users["test@example.com"]
	if !exists {
		t.Fatalf("expected user to be stored after signup")
	}
	if user.ID != userID {
		t.Fatalf("expected stored user ID %q to match returned %q", user.ID, userID)
	}
	if user.AccessToken != tokenPair.AccessToken || user.RefreshToken != tokenPair.RefreshToken {
		t.Fatalf("expected stored tokens to match signup tokens")
	}

	loginTokens, err := svc.Login(ctx, "test@example.com", "secret")
	if err != nil {
		t.Fatalf("expected login to succeed, got error: %v", err)
	}
	if loginTokens.AccessToken == "" || loginTokens.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens after login")
	}

	updatedUser := userStore.users["test@example.com"]
	if updatedUser.AccessToken != loginTokens.AccessToken || updatedUser.RefreshToken != loginTokens.RefreshToken {
		t.Fatalf("expected stored tokens to match login tokens")
	}
}
