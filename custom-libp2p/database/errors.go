package database

import "errors"

var (
	ErrDatabaseRunning = errors.New("database already running")
	ErrDatabaseStopped = errors.New("database not running")
)
