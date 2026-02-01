package util

import (
	"context"
	"time"
)

type CtxFunc func() (context.Context, context.CancelFunc)

func NewCtxFunc(timeout time.Duration) CtxFunc {
	return func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), timeout)
	}
}

func CallCtxFunc(ctx context.Context, ctxFunc CtxFunc) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return ctxFunc()
	}
	return ctx, func() {}
}
