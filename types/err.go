package types

import "errors"

var (
	ErrWrongSize   = errors.New("text or data has wrong size")
	ErrInvalidText = errors.New("text has invalid char")
)
