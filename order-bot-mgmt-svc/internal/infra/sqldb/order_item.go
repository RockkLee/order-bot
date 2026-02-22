package sqldb

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type OrderItemRecord struct {
	ID               string `gorm:"column:id;primaryKey"`
	OrderID          string `gorm:"column:order_id"`
	MenuItemID       string `gorm:"column:menu_item_id"`
	Name             string `gorm:"column:name"`
	Quantity         int32  `gorm:"column:quantity"`
	UnitPriceScaled  int64  `gorm:"column:unit_price_scaled"`
	TotalPriceScaled int64  `gorm:"column:total_price_scaled"`
}

func (OrderItemRecord) TableName() string { return "order_item" }

func OrderItemRecordFromModel(item entities.OrderItem) OrderItemRecord {
	return OrderItemRecord{ID: item.ID, OrderID: item.OrderID, MenuItemID: item.MenuItemID, Name: item.Name, Quantity: item.Quantity, UnitPriceScaled: item.UnitPriceScaled, TotalPriceScaled: item.TotalPriceScaled}
}

type OrderItemStore struct{ db *gorm.DB }

func NewOrderItemStore(db *DB) *OrderItemStore {
	if db == nil {
		panic("sqldb.NewOrderItemStore(), the db ptr is nil")
	}
	return &OrderItemStore{db: db.Gorm()}
}

func (s *OrderItemStore) CreateOrderItems(ctx context.Context, tx store.Tx, items []entities.OrderItem) error {
	db, err := resolveDB(s.db, tx)
	if err != nil {
		return fmt.Errorf("sqldb.OrderItemStore.CreateOrderItems: %w", err)
	}
	records := make([]OrderItemRecord, 0, len(items))
	for _, item := range items {
		records = append(records, OrderItemRecordFromModel(item))
	}
	if err := db.WithContext(ctx).Create(&records).Error; err != nil {
		return fmt.Errorf("sqldb.OrderItemStore.CreateOrderItems: %w", err)
	}
	return nil
}
