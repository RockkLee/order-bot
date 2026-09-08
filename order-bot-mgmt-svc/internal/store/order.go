package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Order interface {
	FindByBotID(ctx context.Context, tx Tx, botId string) ([]entities.Order, error)
	Insert(ctx context.Context, tx Tx, order entities.Order) error
}
