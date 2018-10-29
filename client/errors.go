package client

import "errors"

var (
	ErrDereferencing = errors.New("dereferencing slice of pointers to slice of elements failed")
)
