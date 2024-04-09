package services

import "errors"

var (
	RowNotFoundError = errors.New("desired row not found")
	InvalidArgument  = errors.New("invalid argument")
)
