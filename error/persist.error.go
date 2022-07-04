package error

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrNoQueries = errors.New("you need to provide a query")
)