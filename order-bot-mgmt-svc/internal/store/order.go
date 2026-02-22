package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Order interface {
	FindOrders(ctx context.Context) ([]entities.Order, error)
}
