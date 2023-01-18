package tree

import (
	"errors"
)

var (
	ErrNoop            = errors.New("noop")
	ErrInvalidCommand  = errors.New("invalid command")
	ErrMissingCommand  = errors.New("missing command")
	ErrMissingArgument = errors.New("missing argument")
)
