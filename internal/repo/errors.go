package repo

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("requested item not found")
	ErrDuplicateEmail    = errors.New("user with this email already exists")
	ErrDuplicateUsername = errors.New("user with this username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized      = errors.New("unauthorized access")
)