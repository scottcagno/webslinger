package v4

import (
	"crypto"
	"sync"
)

var methods sync.Map

func GetSigningMethod(name string) SigningMethod {
	fn, found := methods.Load(name)
	if !found {
		return nil
	}
	methodFn, ok := fn.(func() SigningMethod)
	if !ok {
		panic("could not get signing method: method func error")
	}

	var method SigningMethod
	method = methodFn()
	return method
}

func RegisterSigningMethod(name string, fn func() SigningMethod) {
	methods.Store(name, fn)
}

type KeyPair struct {
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
}

// SigningMethod is an interface for implementing singing and verifying methods
type SigningMethod interface {

	// Name should return the name of the signing method.
	Name() string

	// GenerateKeyPair should generate and return a key pair complient
	// with the implementing signing method.
	GenerateKeyPair() *KeyPair

	// Sign should take a base64 encoded header and payload and return a
	// valid signature.
	Sign(partialToken []byte, key crypto.PrivateKey) (signature []byte, err error)

	// Verify should take a token and signature and verify the token using the
	// provided signature.
	Verify(partialToken []byte, signature []byte, key crypto.PublicKey) error
}
