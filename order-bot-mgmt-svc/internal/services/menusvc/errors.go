package menusvc

import "order-bot-mgmt-svc/internal/apperr"

var (
	ErrInvalidMenu = apperr.Err{
		Code: "ErrInvalidMenu",
		Msg:  "invalid menu request",
	}
)
