package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Order interface {
	CreateOrder(ctx context.Context, tx Tx, order entities.Order) error
}
