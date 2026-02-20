package store

import "order-bot-mgmt-svc/internal/apperr"

var (
	ErrInvalidTx = apperr.Err{
		Code: "ErrInvalidTx",
		Msg:  "transaction not started",
	}
	ErrNotFound = apperr.Err{
		Code: "ErrNotFound",
		Msg:  "record not found",
	}
	ErrUserExists = apperr.Err{
		Code: "ErrUserExists",
		Msg:  "user already exists",
	}
	ErrBotNotFound = apperr.Err{
		Code: "ErrBotNotFound",
		Msg:  "bot not found",
	}
	ErrMenuNotFound = apperr.Err{
		Code: "ErrMenuNotFound",
		Msg:  "menu not found",
	}
	ErrUserBotNotFound = apperr.Err{
		Code: "ErrUserBotNotFound",
		Msg:  "user bot not found",
	}
)
