package notifier

import "errors"

var (
	ErrNotifierClosed = errors.New("notifier was closed")
)
