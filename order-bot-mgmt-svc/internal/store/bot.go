package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Bot interface {
	Create(ctx context.Context, tx Tx, bot entities.Bot) error
	FindByID(ctx context.Context, tx Tx, id string) (entities.Bot, error)
}
