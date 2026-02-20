package authsvc

import "order-bot-mgmt-svc/internal/apperr"

var (
	ErrUserExists = apperr.Err{
		Code: "ErrUserExists",
		Msg:  "user already exists",
	}
	ErrInvalidCredentials = apperr.Err{
		Code: "ErrInvalidCredentials",
		Msg:  "invalid credentials",
	}
	ErrInvalidRefreshToken = apperr.Err{
		Code: "ErrInvalidRefreshToken",
		Msg:  "invalid refresh token",
	}
	ErrLoggedOut = apperr.Err{
		Code: "ErrLoggedOut",
		Msg:  "logged out",
	}
)
