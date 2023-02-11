package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

// SigningMethodECDSA implements the ECDSA family of signing methods.
// Expects *ecdsa.PrivateKey for signing and *ecdsa.PublicKey for verification
type SigningMethodECDSA struct {
	name      string
	hash      crypto.Hash
	KeySize   int
	CurveBits int
	curve     elliptic.Curve
}

var (
	ES256 *SigningMethodECDSA
	ES384 *SigningMethodECDSA
	ES512 *SigningMethodECDSA
)

func init() {
	ES256 = &SigningMethodECDSA{
		name:  "ES256",
		hash:  crypto.SHA256,
		curve: elliptic.P256(),
	}
	ES384 = &SigningMethodECDSA{
		name:  "ES384",
		hash:  crypto.SHA384,
		curve: elliptic.P384(),
	}
	ES512 = &SigningMethodECDSA{
		name:  "ES512",
		hash:  crypto.SHA512,
		curve: elliptic.P521(),
	}

	RegisterSigningMethod(ES256.Name(), func() SigningMethod { return ES256 })
	RegisterSigningMethod(ES384.Name(), func() SigningMethod { return ES384 })
	RegisterSigningMethod(ES512.Name(), func() SigningMethod { return ES512 })
}

func (s *SigningMethodECDSA) Name() string {
	return s.name
}

func (s *SigningMethodECDSA) GenerateKeyPair() *KeyPair {
	key, err := ecdsa.GenerateKey(s.curve, rand.Reader)
	if err != nil {
		return nil
	}
	return &KeyPair{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}
}

func (s *SigningMethodECDSA) Sign(partialToken []byte, key crypto.PrivateKey) ([]byte, error) {
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return nil, ErrHashUnavailable
	}
	hasher := s.hash.New()
	hasher.Write(partialToken)
	sig, err := ecdsa.SignASN1(rand.Reader, ecdsaKey, hasher.Sum(nil))
	if err != nil {
		return nil, err
	}
	return Base64Encode(sig), nil
}

func (s *SigningMethodECDSA) Verify(partialToken []byte, signature []byte, key crypto.PublicKey) error {
	ecdsaKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return ErrHashUnavailable
	}
	hasher := s.hash.New()
	hasher.Write(partialToken)
	sig := Base64Decode(signature)
	valid := ecdsa.VerifyASN1(ecdsaKey, hasher.Sum(nil), sig)
	if !valid {
		return ErrSignatureInvalid
	}
	return nil
}
