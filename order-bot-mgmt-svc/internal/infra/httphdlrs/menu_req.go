package httphdlrs

type menuItemRequest struct {
	Name string `json:"name" validate:"required"`
}

type menuRequest struct {
	BotID string            `json:"bot_id" validate:"required"`
	Items []menuItemRequest `json:"items"`
}

type menuItemResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type menuResponse struct {
	ID    string             `json:"id"`
	BotID string             `json:"bot_id"`
	Items []menuItemResponse `json:"items"`
}
