package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type UserBot interface {
	Create(ctx context.Context, userBot entities.UserBot) error
	FindByUserID(ctx context.Context, userID string) ([]entities.UserBot, error)
}
