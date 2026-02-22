package grpcsync

type SubmitOrderRequest struct {
	OrderID     string            `json:"order_id"`
	CartID      string            `json:"cart_id"`
	SessionID   string            `json:"session_id"`
	TotalScaled int64             `json:"total_scaled"`
	Items       []SubmitOrderItem `json:"items"`
}

type SubmitOrderItem struct {
	ID               string `json:"id"`
	MenuItemID       string `json:"menu_item_id"`
	Name             string `json:"name"`
	Quantity         int32  `json:"quantity"`
	UnitPriceScaled  int64  `json:"unit_price_scaled"`
	TotalPriceScaled int64  `json:"total_price_scaled"`
}

type SubmitOrderResponse struct {
	Accepted bool `json:"accepted"`
}

type UpdateOrderStatusRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type UpdateOrderStatusResponse struct {
	Updated bool `json:"updated"`
}
