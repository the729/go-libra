package serialization

import "errors"

var (
	ErrBufferOverflow = errors.New("dest buffer too short")
	ErrWrongSize      = errors.New("input buffer size is invalid")
)
