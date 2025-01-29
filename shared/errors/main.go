package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrCodeMismatch = errors.New("code doesn't match")
	ErrUnexpected   = errors.New("unexpected error")
)
