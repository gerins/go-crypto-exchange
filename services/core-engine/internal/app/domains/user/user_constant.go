package user

import "errors"

var (
	ErrInvalidPassword = errors.New("Invalid username or password")
	ErrUserBlocked     = errors.New("user blocked")
)
