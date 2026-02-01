package store

import (
	"context"
	"order-bot-mgmt-svc/internal/models/entities"
)

type User interface {
	Create(ctx context.Context, tx Tx, user entities.User) error
	FindByEmail(ctx context.Context, tx Tx, email string) (entities.User, error)
	FindByID(ctx context.Context, tx Tx, id string) (entities.User, error)
	UpdateTokens(ctx context.Context, tx Tx, id string, accessToken string, refreshToken string) error
}
