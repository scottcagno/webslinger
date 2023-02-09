package v4

import (
	"crypto"
	"sync"
)

// SigningMethod is an interface for implementing singing and verifying methods
type SigningMethod interface {

	// Name should return the name of the signing method.
	Name() string

	// Sign should take a base64 encoded header and payload and return a
	// valid signature.
	Sign(partialToken []byte, key crypto.PrivateKey) (signature []byte, err error)

	// Verify should take a token and signature and verify the token using the
	// provided signature.
	Verify(partialToken []byte, signature []byte, key crypto.PublicKey) error
}

var methods sync.Map

func GetSigningMethod(name string) SigningMethod {
	method, found := methods.Load(name)
	if !found {
		return nil
	}
	return method.(SigningMethod)
}
