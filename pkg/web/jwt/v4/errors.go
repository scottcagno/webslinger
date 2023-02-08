package v4

import (
	"errors"
)

var (
	ErrInvalidKeyType   = errors.New("alg: invalid key type")
	ErrHashUnavailable  = errors.New("alg: hash unavailable")
	ErrSignatureInvalid = errors.New("alg: signature is invalid")
)
