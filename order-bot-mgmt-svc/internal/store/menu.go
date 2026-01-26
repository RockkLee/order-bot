package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Menu interface {
	BeginTx(ctx context.Context) (MenuTx, error)
	FindByID(ctx context.Context, menuID string) (entities.Menu, error)
	FindItems(ctx context.Context, menuID string) ([]entities.MenuItem, error)
}

type MenuTx interface {
	Commit() error
	Rollback() error
	CreateMenu(ctx context.Context, menu entities.Menu) error
	UpdateMenu(ctx context.Context, menu entities.Menu) error
	DeleteMenu(ctx context.Context, menuID string) error
	DeleteMenuItems(ctx context.Context, menuID string) error
	CreateMenuItems(ctx context.Context, items []entities.MenuItem) error
}
