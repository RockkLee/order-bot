package httphdlr

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/infra/sqldb"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/services/authsvc"
	"order-bot-mgmt-svc/internal/services/botsvc"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/store/fake"
	"order-bot-mgmt-svc/internal/util"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type fakeAuthServer struct {
	authSvc  *authsvc.Svc
	botSvc   *botsvc.Svc
	WithTxFn func(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error
}

func (f *fakeAuthServer) AuthService() *authsvc.Svc { return f.authSvc }
func (f *fakeAuthServer) BotService() *botsvc.Svc   { return f.botSvc }
func (f *fakeAuthServer) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	return f.WithTxFn(ctx, fn)
}

var cfgFakeAuthHdlr = config.Config{
	Auth: config.Auth{
		AccessSecret:    "access",
		RefreshSecret:   "refresh",
		AccessTokenTTL:  time.Minute,
		RefreshTokenTTL: time.Minute,
	},
	Others: config.Others{QryCtxTimeout: time.Second},
}

func TestSignupHdlrFunc(t *testing.T) {
	type output struct {
		status int
	}
	tests := []struct {
		name   string
		body   string
		server func() AuthServer
		out    output
	}{
		{
			name: "happy path",
			body: `{"email":"test@example.com","password":"secret","bot_name":"my-bot"}`,
			server: func() AuthServer {
				ctxFunc := util.NewCtxFunc(cfgFakeAuthHdlr.Others.QryCtxTimeout)
				userStore := &fake.UserStore{
					CreateFn:       func(_ context.Context, _ store.Tx, _ entities.User) error { return nil },
					UpdateTokensFn: func(_ context.Context, _ store.Tx, _ string, _, _ string) error { return nil },
				}
				botStore := &fake.BotStore{
					CreateFn: func(_ context.Context, _ store.Tx, _ entities.Bot) error { return nil },
				}
				userBotStore := &fake.UserBotStore{
					CreateFn: func(_ context.Context, _ store.Tx, _ entities.UserBot) error { return nil },
				}
				return &fakeAuthServer{
					authSvc: authsvc.NewSvc(nil, ctxFunc, cfgFakeAuthHdlr, userStore),
					botSvc:  botsvc.NewSvc(&sqldb.DB{}, ctxFunc, cfgFakeAuthHdlr, botStore, userBotStore),
					WithTxFn: func(ctx context.Context, fn func(context.Context, store.Tx) error) error {
						return fn(ctx, nil)
					},
				}
			},
			out: output{http.StatusCreated},
		},
		{
			name: "invalid JSON body",
			body: `not-json`,
			server: func() AuthServer {
				// handler returns before reaching WithTx — services are never called
				return &fakeAuthServer{}
			},
			out: output{http.StatusBadRequest},
		},
		{
			name: "missing required field",
			body: `{"email":"test@example.com","password":"secret"}`,
			server: func() AuthServer {
				return &fakeAuthServer{}
			},
			out: output{http.StatusBadRequest},
		},
		{
			name: "user already exists",
			body: `{"email":"test@example.com","password":"secret","bot_name":"my-bot"}`,
			server: func() AuthServer {
				ctxFunc := util.NewCtxFunc(cfgFakeAuthHdlr.Others.QryCtxTimeout)
				userStore := &fake.UserStore{
					CreateFn: func(_ context.Context, _ store.Tx, _ entities.User) error {
						return fmt.Errorf("sqldb.UserStore.Create: %w", store.ErrUserExists)
					},
				}
				botStore := &fake.BotStore{
					CreateFn: func(_ context.Context, _ store.Tx, _ entities.Bot) error { return nil },
				}
				userBotStore := &fake.UserBotStore{
					CreateFn: func(_ context.Context, _ store.Tx, _ entities.UserBot) error { return nil },
				}
				return &fakeAuthServer{
					authSvc: authsvc.NewSvc(nil, ctxFunc, cfgFakeAuthHdlr, userStore),
					botSvc:  botsvc.NewSvc(&sqldb.DB{}, ctxFunc, cfgFakeAuthHdlr, botStore, userBotStore),
					WithTxFn: func(ctx context.Context, fn func(context.Context, store.Tx) error) error {
						return fn(ctx, nil)
					},
				}
			},
			out: output{http.StatusConflict},
		},
		{
			name: "bot creation fails",
			body: `{"email":"test@example.com","password":"secret","bot_name":"my-bot"}`,
			server: func() AuthServer {
				ctxFunc := util.NewCtxFunc(cfgFakeAuthHdlr.Others.QryCtxTimeout)
				userStore := &fake.UserStore{
					CreateFn:       func(_ context.Context, _ store.Tx, _ entities.User) error { return nil },
					UpdateTokensFn: func(_ context.Context, _ store.Tx, _ string, _, _ string) error { return nil },
				}
				botStore := &fake.BotStore{
					CreateFn: func(_ context.Context, _ store.Tx, _ entities.Bot) error {
						return store.ErrInvalidTx
					},
				}
				return &fakeAuthServer{
					authSvc: authsvc.NewSvc(nil, ctxFunc, cfgFakeAuthHdlr, userStore),
					botSvc:  botsvc.NewSvc(&sqldb.DB{}, ctxFunc, cfgFakeAuthHdlr, botStore, &fake.UserBotStore{}),
					WithTxFn: func(ctx context.Context, fn func(context.Context, store.Tx) error) error {
						return fn(ctx, nil)
					},
				}
			},
			out: output{http.StatusInternalServerError},
		},
		{
			name: "invalid credentials",
			body: `{"email":"test@example.com","password":"secret","bot_name":"my-bot"}`,
			server: func() AuthServer {
				return &fakeAuthServer{
					WithTxFn: func(_ context.Context, _ func(context.Context, store.Tx) error) error {
						return authsvc.ErrInvalidCredentials
					},
				}
			},
			out: output{http.StatusBadRequest},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/signup", signupHdlrFunc(tt.server()))

			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.out.status {
				t.Errorf("expected status %d, got %d — body: %s", tt.out.status, rec.Code, rec.Body.String())
			}
		})
	}
}
