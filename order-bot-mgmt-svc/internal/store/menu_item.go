package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type MenuItem interface {
	FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error)
	DeleteMenuItems(ctx context.Context, tx Tx, menuID string) error
	CreateMenuItems(ctx context.Context, tx Tx, items []entities.MenuItem) error
}
