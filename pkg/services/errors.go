package services

import "errors"

var (
	ErrorRowNotFound     = errors.New("desired row not found")
	ErrorInvalidArgument = errors.New("invalid argument")
)
