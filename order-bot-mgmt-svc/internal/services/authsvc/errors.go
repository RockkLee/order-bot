package authsvc

import "errors"

var (
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrLoggedOut           = errors.New("logged out")
)
