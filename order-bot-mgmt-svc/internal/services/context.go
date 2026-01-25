package services

import (
	"context"
	"time"
)

func QueryContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
