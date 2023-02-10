package v4

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// SigningMethodRSA implements the RSAPSS family of signing methods.
// Expects *rsa.PrivateKey for signing and *rsa.PublicKey for validation
type SigningMethodRSAPSS struct {
	name string
	hash crypto.Hash
	opts *rsa.PSSOptions
}

var (
	PS256 *SigningMethodRSAPSS
	PS384 *SigningMethodRSAPSS
	PS512 *SigningMethodRSAPSS
)

func init() {
	PS256 = &SigningMethodRSAPSS{
		name: "PS256",
		hash: crypto.SHA256,
		opts: &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash},
	}
	PS384 = &SigningMethodRSAPSS{
		name: "PS384",
		hash: crypto.SHA384,
		opts: &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash},
	}
	PS512 = &SigningMethodRSAPSS{
		name: "PS512",
		hash: crypto.SHA512,
		opts: &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash},
	}
}

func (s *SigningMethodRSAPSS) Name() string {
	return s.name
}

func (s *SigningMethodRSAPSS) GenerateKeyPair() *KeyPair {
	key, err := rsa.GenerateKey(rand.Reader, s.hash.Size()*8)
	if err != nil {
		panic(err)
	}
	return &KeyPair{
		PrivateKey: key,
		PublicKey:  &key.PublicKey,
	}
}

func (s *SigningMethodRSAPSS) Sign(partialToken []byte, key crypto.PrivateKey) ([]byte, error) {
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}
	if !s.hash.Available() {
		return nil, ErrHashUnavailable
	}
	hasher := s.hash.New()
	hasher.Write(partialToken)
	sig, err := rsa.SignPSS(rand.Reader, rsaKey, s.hash, hasher.Sum(nil), s.opts)
	if err != nil {
		return nil, err
	}
	return Base64Encode(sig), nil
}

func (s *SigningMethodRSAPSS) Verify(partialToken []byte, signature []byte, key crypto.PublicKey) error {
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
	return rsa.VerifyPSS(rsaKey, s.hash, hasher.Sum(nil), sig, s.opts)
}
