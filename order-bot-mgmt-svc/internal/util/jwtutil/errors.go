package jwtutil

import "order-bot-mgmt-svc/internal/apperr"

var (
	ErrInvalidToken = apperr.Err{
		Code: "ErrInvalidToken",
		Msg:  "invalid token",
	}
	ErrExpiredToken = apperr.Err{
		Code: "ErrExpiredToken",
		Msg:  "expired token",
	}
)
