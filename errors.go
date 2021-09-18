package disgolf

import "errors"

var (
	// ErrCommandNotExists means that the requested command does not exist.
	ErrCommandNotExists = errors.New("command not exists")
)
