package store

import "errors"

var (
	ErrNotFound         = errors.New("user not found")
	ErrUserExists       = errors.New("user already exists")
	ErrBotNotFound      = errors.New("bot not found")
	ErrMenuNotFound     = errors.New("menu not found")
	ErrMenuItemNotFound = errors.New("menu item not found")
	ErrUserBotNotFound  = errors.New("user bot not found")
)
