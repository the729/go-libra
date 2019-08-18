package types

import "errors"

var (
	// ErrNilInput is error when nil input is not expected.
	ErrNilInput = errors.New("input is nil")
)
