package sqldb

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"

	"gorm.io/gorm"
)

type OrderRecord struct {
	ID          string `gorm:"column:id;primaryKey"`
	BotID       string `gorm:"column:bot_id"`
	CartID      string `gorm:"column:cart_id"`
	SessionID   string `gorm:"column:session_id"`
	TotalScaled int    `gorm:"column:total_scaled"`
}

func (OrderRecord) TableName() string { return "orders" }

func (r OrderRecord) ToModel() entities.Order {
	return entities.Order{
		ID:          r.ID,
		BotID:       r.BotID,
		CartID:      r.CartID,
		SessionID:   r.SessionID,
		TotalScaled: r.TotalScaled,
	}
}

type OrderStore struct{ db *gorm.DB }

func NewOrderStore(db *DB) *OrderStore {
	if db == nil {
		panic("sqldb.NewOrderStore(), the db ptr is nil")
	}
	return &OrderStore{db: db.Gorm()}
}

func (s *OrderStore) FindOrders(ctx context.Context) ([]entities.Order, error) {
	var records []OrderRecord
	if err := s.db.WithContext(ctx).Order("id desc").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("sqldb.OrderStore.FindOrders: %w", err)
	}
	orders := make([]entities.Order, 0, len(records))
	for _, rec := range records {
		orders = append(orders, rec.ToModel())
	}
	return orders, nil
}
