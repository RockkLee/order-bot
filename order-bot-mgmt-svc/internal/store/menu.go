package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type Menu interface {
	FindByBotID(ctx context.Context, menuID string) (entities.Menu, error)
	CreateMenu(ctx context.Context, tx Tx, menu entities.Menu) error
	UpdateMenu(ctx context.Context, tx Tx, menu entities.Menu) error
	DeleteMenu(ctx context.Context, tx Tx, menuID string) error
}
