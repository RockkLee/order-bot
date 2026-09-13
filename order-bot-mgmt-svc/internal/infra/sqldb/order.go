package sqldb

import (
	"context"
	"errors"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type OrderRecord struct {
	Base        BaseRecord `gorm:"embedded"`
	ID          string     `gorm:"column:id;primaryKey"`
	BotID       string     `gorm:"column:bot_id"`
	CartID      string     `gorm:"column:cart_id"`
	SessionID   string     `gorm:"column:session_id"`
	TotalScaled int        `gorm:"column:total_scaled"`
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

func (s *OrderStore) FindByBotID(ctx context.Context, tx store.Tx, botId string) ([]entities.Order, error) {
	db, errDb := resolveDB(s.db, tx)
	if errDb != nil {
		return nil, fmt.Errorf("sqldb.OrderStore.FindByBotID: %w", errDb)
	}
	var records []OrderRecord
	if err := db.WithContext(ctx).Where("bot_id = ?", botId).Find(&records).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("sqldb.OrderStore.FindByBotID: %w", store.ErrBotNotFound)
		}
		return nil, fmt.Errorf("sqldb.OrderStore.FindByBotID: %w", err)
	}
	orders := make([]entities.Order, 0, len(records))
	for _, rec := range records {
		orders = append(orders, rec.ToModel())
	}
	return orders, nil
}

func (s *OrderStore) Insert(ctx context.Context, tx store.Tx, order entities.Order) error {
	db, errDb := resolveDB(s.db, tx)
	if errDb != nil {
		return fmt.Errorf("sqldb.OrderStore.Insert: %w", errDb)
	}
	record := OrderRecord{
		ID:          order.ID,
		BotID:       order.BotID,
		CartID:      order.CartID,
		SessionID:   order.SessionID,
		TotalScaled: order.TotalScaled,
	}
	if err := db.WithContext(ctx).Create(&record).Error; err != nil {
		return fmt.Errorf("sqldb.OrderStore.Insert: %w", err)
	}
	return nil
}
