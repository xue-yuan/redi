package constants

import "errors"

var (
	ErrUnsupportedType      = errors.New("unsupported type")
	ErrCreateShortURLFailed = errors.New("create failed")
	ErrDuplicateShortURL    = errors.New("duplicate short url")
)
