package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Bot interface {
	Create(ctx context.Context, bot entities.Bot) error
	FindByID(ctx context.Context, id string) (entities.Bot, error)
}
