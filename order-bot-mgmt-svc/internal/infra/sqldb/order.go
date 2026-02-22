package sqldb

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type OrderRecord struct {
	ID          string `gorm:"column:id;primaryKey"`
	CartID      string `gorm:"column:cart_id"`
	SessionID   string `gorm:"column:session_id"`
	TotalScaled int64  `gorm:"column:total_scaled"`
}

func (OrderRecord) TableName() string { return "orders" }

func OrderRecordFromModel(order entities.Order) OrderRecord {
	return OrderRecord{ID: order.ID, CartID: order.CartID, SessionID: order.SessionID, TotalScaled: order.TotalScaled}
}

type OrderStore struct{ db *gorm.DB }

func NewOrderStore(db *DB) *OrderStore {
	if db == nil {
		panic("sqldb.NewOrderStore(), the db ptr is nil")
	}
	return &OrderStore{db: db.Gorm()}
}

func (s *OrderStore) CreateOrder(ctx context.Context, tx store.Tx, order entities.Order) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.OrderStore.CreateOrder: %w", err)
	}
	record := OrderRecordFromModel(order)
	if err := db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("sqldb.OrderStore.CreateOrder: %w", err)
	}
	return nil
}
