package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Menu interface {
	FindByID(ctx context.Context, menuID string) (entities.Menu, error)
	FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error)
	CreateMenu(ctx context.Context, tx Tx, menu entities.Menu) error
	UpdateMenu(ctx context.Context, tx Tx, menu entities.Menu) error
	DeleteMenu(ctx context.Context, tx Tx, menuID string) error
	DeleteMenuItems(ctx context.Context, tx Tx, menuID string) error
	CreateMenuItems(ctx context.Context, tx Tx, items []entities.MenuItem) error
}
