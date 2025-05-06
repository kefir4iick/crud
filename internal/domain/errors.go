package domain

import "errors"

var (
	ErrCarNotFound      = errors.New("car not found")
	ErrDuplicateCarID   = errors.New("car with this ID already exists")
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidLimit     = errors.New("invalid limit value")
	ErrInvalidOffset    = errors.New("invalid offset value")
)
