package sqldbexecutor

import (
	"context"
	"database/sql"
	"fmt"
	"order-bot-mgmt-svc/internal/store"
)

type SqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func Executor(db *sql.DB, tx store.Tx) (SqlExecutor, error) {
	if tx == nil {
		return db, nil
	}
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, fmt.Errorf("sqldb.Executor: expected *sql.Tx, got %T", tx)
	}
	return sqlTx, nil
}
