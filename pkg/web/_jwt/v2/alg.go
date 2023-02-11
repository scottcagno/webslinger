package v2

import (
	"crypto"
	"crypto/rsa"
	"errors"
)

var (
	// ErrTokenSignature indicates that the verification failed.
	ErrTokenSignature = errors.New("jwt: invalid token signature")
	// ErrInvalidKey indicates that an algorithm required secret key is not a valid type.
	ErrInvalidKey = errors.New("jwt: invalid key")
)

// Alg represents a signing and verifying algorithm
type Alg interface {

	// Name should return the string representation of the
	// JWT "alg" field
	Name() string

	// Sign should take the private key and a base64 encoded header
	// and payload data and return the signature.
	Sign(k crypto.PrivateKey, payload []byte) ([]byte, error)

	// Verify should take the public key, a base64 encoded header,
	// payload, and signature and verify the signature against
	// the header and payload.
	Verify(k crypto.PublicKey, payload, signature []byte) error
}

// Available signing algorithms
var (
	// HMAC shared secrets, are optimized for speed.
	HS256 Alg = &algHMAC{"HS256", crypto.SHA256}
	HS384 Alg = &algHMAC{"HS384", crypto.SHA384}
	HS512 Alg = &algHMAC{"HS512", crypto.SHA512}
)

type algHMAC struct {
	name   string
	secret crypto.SignerOpts
}

func (a *algHMAC) Name() string {
	// TODO implement me
	panic("implement me")
}

func (a *algHMAC) Sign(k crypto.PrivateKey, payload []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (a *algHMAC) Verify(k crypto.PublicKey, payload, signature []byte) error {
	// TODO implement me
	panic("implement me")
}

var (
	// RSA signing algorithms.
	// Sign   key: *rsa.PrivateKey
	// Verify key: *rsa.PublicKey (or *rsa.PrivateKey with its PublicKey filled)
	//
	// $ openssl genpkey -algorithm rsa -out private_key.pem -pkeyopt rsa_keygen_bits:2048
	// Derive the public key from the private key:
	// $ openssl rsa -pubout -in private_key.pem -out public_key.pem
	RS256 Alg = &algRSA{"RS256", crypto.SHA256}
	RS384 Alg = &algRSA{"RS384", crypto.SHA384}
	RS512 Alg = &algRSA{"RS512", crypto.SHA512}
)

type algRSA struct {
	name   string
	secret crypto.SignerOpts
}

func (a *algRSA) Name() string {
	// TODO implement me
	panic("implement me")
}

func (a *algRSA) Sign(k crypto.PrivateKey, payload []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (a *algRSA) Verify(k crypto.PublicKey, payload, signature []byte) error {
	// TODO implement me
	panic("implement me")
}

var (
	// RSASSA-PSS is another signature scheme with appendix based on RSA.
	PS256 Alg = &algRSAPSS{"PS256", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA256}}
	PS384 Alg = &algRSAPSS{"PS384", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA384}}
	PS512 Alg = &algRSAPSS{"PS512", &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: crypto.SHA512}}
)

type algRSAPSS struct {
	name   string
	secret crypto.SignerOpts
}

func (a *algRSAPSS) Name() string {
	// TODO implement me
	panic("implement me")
}

func (a *algRSAPSS) Sign(k crypto.PrivateKey, payload []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (a *algRSAPSS) Verify(k crypto.PublicKey, payload, signature []byte) error {
	// TODO implement me
	panic("implement me")
}

var (
	// ECDSA signing algorithms.
	// Sign   key: *ecdsa.PrivateKey
	// Verify key: *ecdsa.PublicKey (or *ecdsa.PrivateKey with its PublicKey filled)
	ES256 Alg = &algECDSA{"ES256", crypto.SHA256, 32, 256}
	ES384 Alg = &algECDSA{"ES384", crypto.SHA384, 48, 384}
	ES512 Alg = &algECDSA{"ES512", crypto.SHA512, 66, 521}
)

type algECDSA struct {
	name   string
	secret crypto.SignerOpts
	X, Y   int
}

func (a *algECDSA) Name() string {
	// TODO implement me
	panic("implement me")
}

func (a *algECDSA) Sign(k crypto.PrivateKey, payload []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (a *algECDSA) Verify(k crypto.PublicKey, payload, signature []byte) error {
	// TODO implement me
	panic("implement me")
}

var (
	// Ed25519 Edwards-curve Digital Signature Algorithm.
	// The algorithm's name is: "EdDSA".
	// Sign   key: ed25519.PrivateKey
	// Verify key: ed25519.PublicKey
	EdDSA Alg = &algEdDSA{"EdDSA"}
)

type algEdDSA struct {
	name string
}

func (a *algEdDSA) Name() string {
	// TODO implement me
	panic("implement me")
}

func (a *algEdDSA) Sign(k crypto.PrivateKey, payload []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (a *algEdDSA) Verify(k crypto.PublicKey, payload, signature []byte) error {
	// TODO implement me
	panic("implement me")
}
