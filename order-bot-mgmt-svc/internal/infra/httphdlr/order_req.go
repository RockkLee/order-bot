package httphdlr

type orderItemRes struct {
	ID               string `json:"id"`
	OrderID          string `json:"order_id"`
	MenuItemID       string `json:"menu_item_id"`
	Name             string `json:"name"`
	Quantity         int    `json:"quantity"`
	UnitPriceScaled  int    `json:"unit_price_scaled"`
	TotalPriceScaled int    `json:"total_price_scaled"`
}

type orderRes struct {
	ID          string         `json:"id"`
	BotID       string         `json:"bot_id"`
	CartID      string         `json:"cart_id"`
	SessionID   string         `json:"session_id"`
	TotalScaled int            `json:"total_scaled"`
	Items       []orderItemRes `json:"items"`
}
