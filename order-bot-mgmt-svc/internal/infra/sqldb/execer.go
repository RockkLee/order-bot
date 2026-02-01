package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/store"
)

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func executorForTx(db *sql.DB, tx store.Tx) (sqlExecutor, error) {
	if tx == nil {
		return db, nil
	}
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("sqldb.executorForTx: expected *sql.Tx, got %T", tx)
	}
	return sqlTx, nil
}
