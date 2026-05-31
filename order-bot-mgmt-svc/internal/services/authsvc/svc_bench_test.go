package authsvc

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/apperr"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/resource"
	"order-bot-mgmt-svc/internal/store"
	"order-bot-mgmt-svc/internal/store/fake"
	"order-bot-mgmt-svc/internal/util"
	"testing"
)

func BenchmarkSvcSignup(b *testing.B) {
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
		b.Run(tt.name, func(b *testing.B) {
			fakeUserStore := &fake.UserStore{}
			fakeUserStore.CreateFn = func(_ context.Context, _ store.Tx, user entities.User) error {
				return tt.errUsrStoreCrt
			}
			fakeUserStore.UpdateTokensFn = func(ctx context.Context, tx store.Tx, id string, accessToken string, refreshToken string) error {
				return nil
			}
			svc := NewSvc(resource.New(nil, nil, resource.GRPCConn{}), ctxFunc, testCfg, fakeUserStore)
			b.ResetTimer() // Reset the benchmark timer so setup work above is excluded from measurement.
			for i := 0; i < b.N; i++ {
				_, _, err := svc.Signup(ctx, nil, tt.email, tt.password)
				if !errors.Is(err, tt.out.err) {
					var apperror apperr.Err
					if !errors.As(err, &apperror) {
						panic("FAILED: errors.As")
					}
					b.Errorf("Signup(%q, %q): got %v, want %v",
						tt.email, tt.password, apperror.Code, tt.out.err.(apperr.Err).Code)
				}
			}
		})
	}
}
