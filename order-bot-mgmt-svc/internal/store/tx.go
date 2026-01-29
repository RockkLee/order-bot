package store

import "context"

type TxBeginner[T any] interface {
	BeginTx(ctx context.Context) (T, error)
}
