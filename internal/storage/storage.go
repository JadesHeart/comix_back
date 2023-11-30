package storage

import "errors"

var (
	ErrURLNotFound   = errors.New("URL NOT FOUND")
	ErrURLExists     = errors.New("URL EXISTS")
	ComixTagIsExists = errors.New("TAG EXISTS")
)
