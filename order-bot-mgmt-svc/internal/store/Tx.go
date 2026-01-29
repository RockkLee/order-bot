package store

import "context"

type Tx interface {
	Commit() error
	Rollback() error
}

type TxStore interface {
	BeginTx(ctx context.Context) (Tx, error)
}
