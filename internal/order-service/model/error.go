package model

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderSubmit   = errors.New("unable to submit order")
)
