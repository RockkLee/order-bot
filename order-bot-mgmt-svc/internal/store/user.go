package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type User interface {
	Create(ctx context.Context, user entities.User) error
	FindByEmail(ctx context.Context, email string) (entities.User, error)
	FindByID(ctx context.Context, id string) (entities.User, error)
	UpdateTokens(ctx context.Context, id string, accessToken string, refreshToken string) error
}
