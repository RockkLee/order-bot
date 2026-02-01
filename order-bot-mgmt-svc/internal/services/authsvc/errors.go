package authsvc

import "errors"

var (
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrLoggedOut           = errors.New("logged out")
)
