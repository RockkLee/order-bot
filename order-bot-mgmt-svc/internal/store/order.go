package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Order interface {
	FindByBotID(ctx context.Context, tx Tx, botId string) ([]entities.Order, error)
}
