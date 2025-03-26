package redis

import "errors"

var (
	ErrCacheMiss   = errors.New("cache: key not found")
	ErrCacheFailed = errors.New("cache: failed to get key")
	ErrUnmarshal   = errors.New("cache: failed to unmarshal data")
	ErrMarshal     = errors.New("cache: failed to marshal data")
)
