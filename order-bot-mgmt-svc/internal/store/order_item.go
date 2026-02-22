package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type OrderItem interface {
	CreateOrderItems(ctx context.Context, tx Tx, items []entities.OrderItem) error
}
