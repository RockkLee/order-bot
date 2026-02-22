package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type OrderItem interface {
	FindByOrderIDs(ctx context.Context, orderIDs []string) ([]entities.OrderItem, error)
}
