package services

import (
	"context"
	"time"
)

type ContextFactory func() (context.Context, context.CancelFunc)

func NewContextFactory(timeout time.Duration) ContextFactory {
	return func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), timeout)
	}
}
