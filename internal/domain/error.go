package domain

import "errors"

var (
	ErrNotFound   = errors.New("record not found")
	ErrConflict   = errors.New("record already exists")
	ErrInvalidOTP = errors.New("invalid otp")
	ErrInvalidJWT = errors.New("invalid jwt")
)
