package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"order-bot-mgmt-svc/internal/config"
	"order-bot-mgmt-svc/internal/store"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service interface {
	Health() (map[string]string, error)
	Close() error
	Conn() *sql.DB
	WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error
	GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error)
}

type DB struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func New(cfg config.Db) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s search_path=%s sslmode=disable",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.Port,
		cfg.Schema,
	)
	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, fmt.Errorf("sqldb.New: %w", err)
	}
	sdb, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("sqldb.New: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := sdb.PingContext(ctx); err != nil {
		_ = sdb.Close()
		return nil, fmt.Errorf("sqldb.New: %w", err)
	}
	return &DB{gormDB: gdb, sqlDB: sdb}, nil
}

func (s *DB) Gorm() *gorm.DB { return s.gormDB }

func (s *DB) Health() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stats := make(map[string]string)
	if err := s.sqlDB.PingContext(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats, fmt.Errorf("sqldb.DB.Health: %w", err)
	}
	dbStats := s.sqlDB.Stats()
	stats["status"] = "up"
	stats["message"] = "It's healthy"
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)
	return stats, nil
}

func (s *DB) Close() error {
	log.Printf("Disconnected from database")
	if err := s.sqlDB.Close(); err != nil {
		return fmt.Errorf("sqldb.DB.Close: %w", err)
	}
	return nil
}

func (s *DB) Conn() *sql.DB { return s.sqlDB }

func (s *DB) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) error) error {
	if fn == nil {
		return fmt.Errorf("sqldb.DB.WithTx: fn is nil")
	}
	return s.gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if tx == nil {
			return fmt.Errorf("sqldb.DB.WithTx: nil gorm tx")
		}
		return fn(ctx, tx)
	})
}

func (s *DB) GetWithTx(ctx context.Context, fn func(ctx context.Context, tx store.Tx) (any, error)) (any, error) {
	if fn == nil {
		return nil, fmt.Errorf("sqldb.DB.GetWithTx: fn is nil")
	}
	var out any
	err := s.gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result, err := fn(ctx, tx)
		if err != nil {
			return err
		}
		out = result
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("sqldb.DB.GetWithTx: %w", err)
	}
	return out, nil
}

// Allow services and handle can open a transaction with infra separation
func resolveDB(base *gorm.DB, tx store.Tx) (*gorm.DB, error) {
	if tx == nil {
		return base, nil
	}
	gtx, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("sqldb.resolveDB: expected *gorm.DB, got %T", tx)
	}
	return gtx, nil
}
