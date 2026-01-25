package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type MenuItem interface {
	Create(ctx context.Context, item entities.MenuItem) error
	FindByMenuID(ctx context.Context, menuID string) ([]entities.MenuItem, error)
	DeleteByMenuID(ctx context.Context, menuID string) error
}
