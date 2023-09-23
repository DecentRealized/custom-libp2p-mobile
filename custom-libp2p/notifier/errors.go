package notifier

import "errors"

var (
	ErrClosed = errors.New("notifier was closed")
)
