package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// SigningMethodRSA implements the RSA family of signing methods.
// Expects *rsa.PrivateKey for signing and *rsa.PublicKey for validation
type SigningMethodRSA struct {
	name string
	hash crypto.Hash
}

var (
	RS256 *SigningMethodRSA
	RS384 *SigningMethodRSA
	RS512 *SigningMethodRSA
)

func init() {
	RS256 = &SigningMethodRSA{
		name: "RS256",
		hash: crypto.SHA256,
	}
	RS384 = &SigningMethodRSA{
		name: "RS384",
		hash: crypto.SHA384,
	}
	RS512 = &SigningMethodRSA{
		name: "RS512",
		hash: crypto.SHA512,
	}

	RegisterSigningMethod(RS256.Name(), func() SigningMethod { return RS256 })
	RegisterSigningMethod(RS384.Name(), func() SigningMethod { return RS384 })
	RegisterSigningMethod(RS512.Name(), func() SigningMethod { return RS512 })
}

const rsaBits = 2048

func (s *SigningMethodRSA) Name() string {
	return s.name
}

func (s *SigningMethodRSA) GenerateKeyPair() *KeyPair {
	key, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		panic(err)
	}
	return &KeyPair{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}
}

func (s *SigningMethodRSA) Sign(partialToken []byte, key crypto.PrivateKey) ([]byte, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return nil, ErrHashUnavailable
	}
	hasher := s.hash.New()
	hasher.Write(partialToken)
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, s.hash, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}
	return Base64Encode(sig), nil
}

func (s *SigningMethodRSA) Verify(partialToken []byte, signature []byte, key crypto.PublicKey) error {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return ErrHashUnavailable
	}
	hasher := s.hash.New()
	hasher.Write(partialToken)
	sig := Base64Decode(signature)
	return rsa.VerifyPKCS1v15(rsaKey, s.hash, hasher.Sum(nil), sig)
}
