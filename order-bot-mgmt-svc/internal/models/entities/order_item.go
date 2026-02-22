package entities

type OrderItem struct {
	ID               string
	OrderID          string
	MenuItemID       string
	Name             string
	Quantity         int32
	UnitPriceScaled  int64
	TotalPriceScaled int64
}
