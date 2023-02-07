package v2

import (
	"crypto"
	"errors"
)

var (
	// ErrMissing indicates that a given token to `Verify` is empty.
	ErrMissing = errors.New("jwt: token is empty")
	// ErrTokenForm indicates that the extracted token has not the expected form .
	ErrTokenForm = errors.New("jwt: invalid token form")
	// ErrTokenAlg indicates that the given algorithm does not match the extracted one.
	ErrTokenAlg = errors.New("jwt: unexpected token algorithm")
)

// PrivateKey is a semi-generic type responsible for signing the token.
type PrivateKey = crypto.PrivateKey

// PublicKey is a semi-generic type responsible for verifying the token.
type PublicKey interface {
	Equal(x crypto.PublicKey) bool
}
