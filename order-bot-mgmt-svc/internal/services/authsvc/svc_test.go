package authsvc

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/apperr"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/util"
	"testing"
	"time"
)

type fakeUserStore struct {
	CreateFn       func(ctx context.Context, tx store.Tx, user entities.User) error
	FindByEmailFn  func(ctx context.Context, tx store.Tx, email string) (entities.User, error)
	FindByIDFn     func(ctx context.Context, tx store.Tx, id string) (entities.User, error)
	UpdateTokensFn func(ctx context.Context, tx store.Tx, id string, accessToken string, refreshToken string) error
}

func (f *fakeUserStore) Create(ctx context.Context, tx store.Tx, user entities.User) error {
	return f.CreateFn(ctx, tx, user)
}

func (f *fakeUserStore) FindByEmail(ctx context.Context, tx store.Tx, email string) (entities.User, error) {
	return f.FindByEmailFn(ctx, tx, email)
}

func (f *fakeUserStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.User, error) {
	return f.FindByIDFn(ctx, tx, id)
}

func (f *fakeUserStore) UpdateTokens(ctx context.Context, tx store.Tx, id string, accessToken string, refreshToken string) error {
	return f.UpdateTokensFn(ctx, tx, id, accessToken, refreshToken)
}

var testCfg = config.Config{
	Auth: config.Auth{
		AccessSecret:    "access",
		RefreshSecret:   "refresh",
		AccessTokenTTL:  time.Minute,
		RefreshTokenTTL: time.Minute,
	},
	Others: config.Others{QryCtxTimeout: time.Second},
}

func TestSvcSignup(t *testing.T) {
	type output struct {
		err error
	}
	tests := []struct {
		name           string
		email          string
		password       string
		errUsrStoreCrt error
		out            output
	}{
		{name: "happy path", email: "123", password: "123", errUsrStoreCrt: nil, out: output{nil}},
		{name: "unexpected store error", email: "123", password: "123", errUsrStoreCrt: store.ErrInvalidTx, out: output{store.ErrInvalidTx}},
		{name: "empty email", email: "", password: "123", errUsrStoreCrt: nil, out: output{ErrInvalidCredentials}},
		{name: "empty password", email: "123", password: "", errUsrStoreCrt: nil, out: output{ErrInvalidCredentials}},
		{name: "user already exists", email: "123", password: "456", errUsrStoreCrt: fmt.Errorf("sqldb.UserStore.Create: %w", store.ErrUserExists), out: output{ErrUserExists}},
	}

	ctxFunc := util.NewCtxFunc(testCfg.Others.QryCtxTimeout)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeUserStore := &fakeUserStore{}
			fakeUserStore.CreateFn = func(_ context.Context, _ store.Tx, user entities.User) error {
				return tt.errUsrStoreCrt
			}
			fakeUserStore.UpdateTokensFn = func(ctx context.Context, tx store.Tx, id string, accessToken string, refreshToken string) error {
				return nil
			}
			svc := NewSvc(nil, ctxFunc, testCfg, fakeUserStore)
			_, _, err := svc.Signup(ctx, nil, tt.email, tt.password)
			if !errors.Is(err, tt.out.err) {
				var apperror apperr.Err
				if !errors.As(err, &apperror) {
					panic("FAILED: errors.As")
				}
				t.Errorf("Signup(%q, %q): got %v, want %v",
					tt.email, tt.password, apperror.Code, tt.out.err.(apperr.Err).Code)
			}
		})
	}
}
