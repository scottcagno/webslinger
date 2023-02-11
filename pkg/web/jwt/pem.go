package jwt

import (
	"crypto"
	"crypto/dsa"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	ErrParsingPEMPublicKey  = errors.New("failed to decode PEM block containing public key")
	ErrParsingPEMPrivateKey = errors.New("failed to decode PEM block containing private key")
)

// ParsePrivateKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 public key
func ParsePrivateKeyFromPEM(key []byte) (crypto.PrivateKey, error) {
	// Parse PEM block
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, ErrParsingPEMPrivateKey
	}
	// Parse the key
	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, ErrParsingPEMPublicKey
	}
	return pri, nil

}

// ParsePublicKeyFromPEM parses a PEM encoded PKCS1 or PKCS8 public key
func ParsePublicKeyFromPEM(key []byte) (crypto.PublicKey, error) {
	// Parse PEM block
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, ErrParsingPEMPublicKey
	}
	// Parse the key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrParsingPEMPublicKey
	}
	return pub, nil
}

func PrivateKeyType(key crypto.PrivateKey) {
	switch pri := key.(type) {
	case *rsa.PrivateKey:
		fmt.Println("Private key type RSA:", pri)
	case *dsa.PrivateKey:
		fmt.Println("Private key type DSA:", pri)
	case *ecdsa.PrivateKey:
		fmt.Println("Private key type ECDSA:", pri)
	case ed25519.PrivateKey:
		fmt.Println("Private key type Ed25519:", pri)
	case *ecdh.PrivateKey:
		fmt.Println("Private key type ECDH:", pri)
	default:
		fmt.Println("Private key type unknown")
	}
}

func PublicKeyType(key crypto.PublicKey) {
	switch pub := key.(type) {
	case *rsa.PublicKey:
		fmt.Println("Public key type RSA:", pub)
	case *dsa.PublicKey:
		fmt.Println("Public key type DSA:", pub)
	case *ecdsa.PublicKey:
		fmt.Println("Public key type ECDSA:", pub)
	case ed25519.PublicKey:
		fmt.Println("Public key type Ed25519:", pub)
	case *ecdh.PublicKey:
		fmt.Println("Public key type ECDH:", pub)
	default:
		fmt.Println("Public key type unknown")
	}
}
