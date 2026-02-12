package pqsqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/store"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() (map[string]string, error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// Conn returns the underlying SQL connection.
	Conn() *sql.DB

	// WithTx runs the given function within a transaction.
	WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error

	// GetWithTx runs the given function within a transaction.
	GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error)
}

type DB struct {
	db *sql.DB
}

func New(cfg config.Db) (*DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Schema,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("pqsqldb.New: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("failed to close database after ping failure: %v", closeErr)
		}
		return nil, fmt.Errorf("pqsqldb.New: %w", err)
	}
	return &DB{
		db: db,
	}, nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *DB) Health() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats, fmt.Errorf("pqsqldb.DB.Health: %w", err)
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats, nil
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *DB) Close() error {
	log.Printf("Disconnected from database")
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("pqsqldb.DB.Close: %w", err)
	}
	return nil
}

func (s *DB) Conn() *sql.DB {
	return s.db
}

func (s *DB) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	if fn == nil {
		return fmt.Errorf("pqsqldb.DB.WithTx: fn is nil")
	}
	tx, errTx := s.db.BeginTx(ctx, nil)
	if errTx != nil {
		return fmt.Errorf("pqsqldb.DB.WithTx() BeginTx: %w", errTx)
	}
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("pqsqldb.DB.WithTx: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("pqsqldb.DB.WithTx: %w", err)
	}
	return nil
}

func (s *DB) GetWithTx(
	ctx context.Context,
	fn func(ctx context.Context, tx store.Tx) (any, error),
) (any, error) {
	if fn == nil {
		return nil, fmt.Errorf("pqsqldb.DB.WithTx: fn is nil")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("pqsqldb.DB.WithTx: BeginTx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	out, err := fn(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("pqsqldb.DB.WithTx: fn: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("pqsqldb.DB.WithTx: Commit: %w", err)
	}
	committed = true

	return out, nil
}
