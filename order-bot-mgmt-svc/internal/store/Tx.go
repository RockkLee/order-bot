package store

import "context"

type Tx interface {
	Commit() error
	Rollback() error
}

type TxBeginner interface {
	BeginTx(ctx context.Context) (Tx, error)
}
