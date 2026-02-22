package entities

type OrderItem struct {
	ID               string
	OrderID          string
	MenuItemID       string
	Name             string
	Quantity         int
	UnitPriceScaled  int
	TotalPriceScaled int
}
