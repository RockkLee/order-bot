package store

import "errors"

var (
	ErrInvalidTx       = errors.New("transaction not started")
	ErrNotFound        = errors.New("record not found")
	ErrUserExists      = errors.New("user already exists")
	ErrBotNotFound     = errors.New("bot not found")
	ErrMenuNotFound    = errors.New("menu not found")
	ErrUserBotNotFound = errors.New("user bot not found")
)
