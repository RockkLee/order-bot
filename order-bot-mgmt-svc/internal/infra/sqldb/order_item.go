package sqldb

import (
	"context"
	"fmt"
	"order-bot-mgmt-svc/internal/models/entities"
	"order-bot-mgmt-svc/internal/store"

	"gorm.io/gorm"
)

type OrderItemRecord struct {
	Base             BaseRecord `gorm:"embedded"`
	ID               string     `gorm:"column:id;primaryKey"`
	OrderID          string     `gorm:"column:order_id"`
	MenuItemID       string     `gorm:"column:menu_item_id"`
	Name             string     `gorm:"column:name"`
	Quantity         int        `gorm:"column:quantity"`
	UnitPriceScaled  int        `gorm:"column:unit_price_scaled"`
	TotalPriceScaled int        `gorm:"column:total_price_scaled"`
}

func (OrderItemRecord) TableName() string { return "order_item" }

func (r OrderItemRecord) ToModel() entities.OrderItem {
	return entities.OrderItem{
		ID:               r.ID,
		OrderID:          r.OrderID,
		MenuItemID:       r.MenuItemID,
		Name:             r.Name,
		Quantity:         r.Quantity,
		UnitPriceScaled:  r.UnitPriceScaled,
		TotalPriceScaled: r.TotalPriceScaled,
	}
}

type OrderItemStore struct{ db *gorm.DB }

func NewOrderItemStore(db *DB) *OrderItemStore {
	if db == nil {
		panic("sqldb.NewOrderItemStore(), the db ptr is nil")
	}
	return &OrderItemStore{db: db.Gorm()}
}

func (s *OrderItemStore) FindByOrderIDs(ctx context.Context, orderIDs []string) ([]entities.OrderItem, error) {
	if len(orderIDs) == 0 {
		return []entities.OrderItem{}, nil
	}
	var records []OrderItemRecord
	if err := s.db.WithContext(ctx).Where("order_id IN ?", orderIDs).Order("order_id desc").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("sqldb.OrderItemStore.FindByOrderIDs: %w", err)
	}
	items := make([]entities.OrderItem, 0, len(records))
	for _, rec := range records {
		items = append(items, rec.ToModel())
	}
	return items, nil
}

func (s *OrderItemStore) InsertMany(ctx context.Context, tx store.Tx, items []entities.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	db, errDb := resolveDB(s.db, tx)
	if errDb != nil {
		return fmt.Errorf("sqldb.OrderItemStore.InsertMany: %w", errDb)
	}
	records := make([]OrderItemRecord, 0, len(items))
	for _, item := range items {
		records = append(records, OrderItemRecord{
			ID:               item.ID,
			OrderID:          item.OrderID,
			MenuItemID:       item.MenuItemID,
			Name:             item.Name,
			Quantity:         item.Quantity,
			UnitPriceScaled:  item.UnitPriceScaled,
			TotalPriceScaled: item.TotalPriceScaled,
		})
	}
	if err := db.WithContext(ctx).Create(&records).Error; err != nil {
		return fmt.Errorf("sqldb.OrderItemStore.InsertMany: %w", err)
	}
	return nil
}
