package fake

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type UserBotStore struct {
	CreateFn       func(ctx context.Context, tx store.Tx, userBot entities.UserBot) error
	FindByUserIDFn func(ctx context.Context, tx store.Tx, userID string) ([]entities.UserBot, error)
}

func (f *UserBotStore) Create(ctx context.Context, tx store.Tx, userBot entities.UserBot) error {
	return f.CreateFn(ctx, tx, userBot)
}
func (f *UserBotStore) FindByUserID(ctx context.Context, tx store.Tx, userID string) ([]entities.UserBot, error) {
	return f.FindByUserIDFn(ctx, tx, userID)
}
