package fake

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type UserStore struct {
	CreateFn       func(ctx context.Context, tx store.Tx, user entities.User) error
	FindByEmailFn  func(ctx context.Context, tx store.Tx, email string) (entities.User, error)
	FindByIDFn     func(ctx context.Context, tx store.Tx, id string) (entities.User, error)
	UpdateTokensFn func(ctx context.Context, tx store.Tx, id string, accessToken, refreshToken string) error
}

func (f *UserStore) Create(ctx context.Context, tx store.Tx, user entities.User) error {
	return f.CreateFn(ctx, tx, user)
}
func (f *UserStore) FindByEmail(ctx context.Context, tx store.Tx, email string) (entities.User, error) {
	return f.FindByEmailFn(ctx, tx, email)
}
func (f *UserStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.User, error) {
	return f.FindByIDFn(ctx, tx, id)
}
func (f *UserStore) UpdateTokens(ctx context.Context, tx store.Tx, id string, accessToken, refreshToken string) error {
	return f.UpdateTokensFn(ctx, tx, id, accessToken, refreshToken)
}
