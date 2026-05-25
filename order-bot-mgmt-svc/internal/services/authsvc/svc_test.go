package authsvc

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/apperr"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/models"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/store/fake"
	"order-bot-mgmt-svc/internal/util"
	"order-bot-mgmt-svc/internal/util/jwtutil"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
			fakeUserStore := &fake.UserStore{}
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

func TestSvcLogin(t *testing.T) {
	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	mustHashPassword := func(t *testing.T, password string) string {
		t.Helper()

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatalf("bcrypt.GenerateFromPassword() error = %v", err)
		}

		return string(hash)
	}

	tests := []struct {
		name    string
		args    args
		user    entities.User
		findErr error
		want    models.Claims
		wantErr bool
	}{
		{
			name: "Happy path",
			args: args{
				ctx:      context.Background(),
				email:    "user@example.com",
				password: "password123",
			},
			user: entities.User{
				ID:           "user-1",
				Email:        "user@example.com",
				PasswordHash: mustHashPassword(t, "password123"),
			},
			want: models.Claims{
				Sub:   "user-1",
				Email: "user@example.com",
				Typ:   "access",
			},
		},
		{
			name: "The user is not found",
			args: args{
				ctx:      context.Background(),
				email:    "missing@example.com",
				password: "password123",
			},
			findErr: store.ErrNotFound,
			wantErr: true,
		},
		{
			name: "Comparison between the password and the passwordHash is incorrect",
			args: args{
				ctx:      context.Background(),
				email:    "user@example.com",
				password: "wrong-password",
			},
			user: entities.User{
				ID:           "user-3",
				Email:        "user@example.com",
				PasswordHash: mustHashPassword(t, "correct-password"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeUserStore := &fake.UserStore{}
			fakeUserStore.FindByEmailFn = func(_ context.Context, _ store.Tx, email string) (entities.User, error) {
				if email != tt.args.email {
					t.Fatalf("FindByEmail() email = %q, want %q", email, tt.args.email)
				}

				return tt.user, tt.findErr
			}
			fakeUserStore.UpdateTokensFn = func(_ context.Context, _ store.Tx, id string, accessToken string, refreshToken string) error {
				if tt.wantErr {
					t.Fatalf("UpdateTokens() should not be called when error is expected")
				}
				if id != tt.user.ID {
					t.Fatalf("UpdateTokens() id = %q, want %q", id, tt.user.ID)
				}
				if accessToken == "" || refreshToken == "" {
					t.Fatalf("UpdateTokens() got empty tokens")
				}

				return nil
			}

			svc := NewSvc(nil, util.NewCtxFunc(testCfg.Others.QryCtxTimeout), testCfg, fakeUserStore)
			got, err := svc.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.name == "The user is not found" && !errors.Is(err, ErrInvalidCredentials) {
					t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
				}
				if tt.name == "Comparison between the password and the passwordHash is incorrect" && got != (models.TokenPair{}) {
					t.Errorf("Login() got = %v, want empty token pair", got)
				}
				return
			}
			if got.AccessToken == "" || got.RefreshToken == "" {
				t.Errorf("Login() got = %v, want non-empty token pair", got)
			}
			claims, err := jwtutil.ParseJWT([]byte(testCfg.Auth.AccessSecret), got.AccessToken, time.Now())
			if err != nil {
				t.Fatalf("jwtutil.ParseJWT() error = %v", err)
			}
			if claims.Sub != tt.want.Sub {
				t.Errorf("Login() access token sub = %q, want %q", claims.Sub, tt.want.Sub)
			}
			if claims.Email != tt.want.Email {
				t.Errorf("Login() access token email = %q, want %q", claims.Email, tt.want.Email)
			}
			if claims.Typ != tt.want.Typ {
				t.Errorf("Login() access token typ = %q, want %q", claims.Typ, tt.want.Typ)
			}
		})
	}
}
