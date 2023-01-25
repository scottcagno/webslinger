package jwt

import (
	"crypto"
)

// PublicKey is a generic public key type that is
// also compatible with all the crypto.PublicKey
// types and implementations in the standard library
// and is also extensible.
type PublicKey interface {
	Equal(x crypto.PublicKey) bool
}

// PrivateKey is a generic private key type that is
// also compatible with all the crypto.PrivateKey
// types and implementations in the standard library
// and is also extensible.
type PrivateKey interface {
	Public() crypto.PublicKey
	Equal(x crypto.PrivateKey) bool
}

// Signer is an interface for signing a JWT token, and
// return the JWT token as a signed encoded string.
type Signer interface {
	// Sign should take a Token and a private key and generate
	// the signature and return a fully encoded JWT token string.
	Sign(token *Token, key PrivateKey) (string, error)
}

// Verifier is an interface for validating a JWT token
// string reporting if it is valid or invalid.
type Verifier interface {
	// Verify should take a signed JWT token string and verify
	// that the signature is valid.
	Verify(token string, key PublicKey) error
}

// Parser is an interface for validating a JWT token string
// and returning a decoding it back into a *Token type.
type Parser interface {
	// Parse should take a JWT token string, and a public
	// key and use it to verify
	Parse(token string, publicKey string, cust func()) (*Token, error)
}
