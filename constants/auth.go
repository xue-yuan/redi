package constants

import "errors"

type AuthType int8

const (
	HardAuth AuthType = iota
	SoftAuth
)

var (
	ErrJWTMissing = errors.New("missing JWT")
)
