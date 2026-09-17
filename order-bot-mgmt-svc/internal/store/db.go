package store

import (
	"context"
	"database/sql"
)

// DB provides the database operations shared by application services and transports.
type DB interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() (map[string]string, error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// Conn returns the underlying SQL connection.
	Conn() *sql.DB

	// WithTx runs the given function within a transaction.
	WithTx(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error

	// GetWithTx runs the given function within a transaction.
	GetWithTx(ctx context.Context, fn func(ctx context.Context, tx Tx) (any, error)) (any, error)
}
