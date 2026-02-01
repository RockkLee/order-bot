package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type UserBot interface {
	Create(ctx context.Context, tx Tx, userBot entities.UserBot) error
	FindByUserID(ctx context.Context, tx Tx, userID string) ([]entities.UserBot, error)
}
