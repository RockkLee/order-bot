package fake

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"
)

type BotStore struct {
	CreateFn   func(ctx context.Context, tx store.Tx, bot entities.Bot) error
	FindByIDFn func(ctx context.Context, tx store.Tx, id string) (entities.Bot, error)
}

func (f *BotStore) Create(ctx context.Context, tx store.Tx, bot entities.Bot) error {
	return f.CreateFn(ctx, tx, bot)
}
func (f *BotStore) FindByID(ctx context.Context, tx store.Tx, id string) (entities.Bot, error) {
	return f.FindByIDFn(ctx, tx, id)
}
