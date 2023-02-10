package v4

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
)

type HMACPrivateKey []byte
type HMACPublicKey []byte

func (k HMACPrivateKey) Public() crypto.PublicKey {
	return k
}

func (k HMACPrivateKey) Equal(x crypto.PrivateKey) bool {
	return bytes.Equal(k, x.([]byte))
}

func (k HMACPublicKey) Equal(x crypto.PublicKey) bool {
	return bytes.Equal(k, x.([]byte))
}

// SigningMethodHMAC implements the HMAC-SHA family of signing methods.
// Expects key type of []byte for both signing and validation
type SigningMethodHMAC struct {
	name string
	hash crypto.Hash
}

var (
	HS256 *SigningMethodHMAC
	HS384 *SigningMethodHMAC
	HS512 *SigningMethodHMAC
)

func init() {
	HS256 = &SigningMethodHMAC{
		name: "HS256",
		hash: crypto.SHA256,
	}
	HS384 = &SigningMethodHMAC{
		name: "HS384",
		hash: crypto.SHA384,
	}
	HS512 = &SigningMethodHMAC{
		name: "HS512",
		hash: crypto.SHA512,
	}

	RegisterSigningMethod(HS256.Name(), func() SigningMethod { return HS256 })
	RegisterSigningMethod(HS384.Name(), func() SigningMethod { return HS384 })
	RegisterSigningMethod(HS512.Name(), func() SigningMethod { return HS512 })
}

func (s *SigningMethodHMAC) Name() string {
	return s.name
}

func (s *SigningMethodHMAC) GenerateKeyPair() *KeyPair {
	buf := make([]byte, s.hash.Size())
	_, err := rand.Read(buf)
	if err != nil {
		return nil
	}
	base64.RawURLEncoding.Encode(buf, buf)
	return &KeyPair{
		PrivateKey: HMACPrivateKey(buf),
		PublicKey:  HMACPublicKey(buf),
	}
}

func (s *SigningMethodHMAC) Sign(partialToken []byte, key crypto.PrivateKey) ([]byte, error) {
	k, ok := key.([]byte)
	if !ok {
		return nil, ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return nil, ErrHashUnavailable
	}
	hasher := hmac.New(s.hash.New, k)
	hasher.Write(partialToken)
	return Base64Encode(hasher.Sum(nil)), nil
}

func (s *SigningMethodHMAC) Verify(partialToken []byte, signature []byte, key crypto.PublicKey) error {
	k, ok := key.([]byte)
	if !ok {
		return ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return ErrHashUnavailable
	}
	sig := Base64Decode(signature)
	// The HMAC signing method is a symmetric one. We will validate the
	// signature by reproducing the signature from the partial token,
	// then compare it against the provided signature.
	hasher := hmac.New(s.hash.New, k)
	hasher.Write(partialToken)
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}
	return nil
}
