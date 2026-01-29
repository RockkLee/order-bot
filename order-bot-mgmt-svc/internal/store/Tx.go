package store

import "context"

type Tx interface {
	Commit() error
	Rollback() error
}

type DBStore interface {
	BeginTx(ctx context.Context) (Tx, error)
}
