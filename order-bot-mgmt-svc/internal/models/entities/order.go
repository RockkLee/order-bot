package entities

type Order struct {
	ID          string
	CartID      string
	SessionID   string
	TotalScaled int
}
