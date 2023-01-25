package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"strings"
)

// SigningMethod strings
const (
	ALG_RS256 = "RS256"
	ALG_RS384 = "RS384"
	ALG_RS512 = "RS512"
	ALG_HS256 = "HS256"
	ALG_HS384 = "HS384"
	ALG_HS512 = "HS512"
	ALG_ES256 = "ES256"
	ALG_ES384 = "ES384"
	ALG_ES512 = "ES512"
)

// SigningMethod represents a signing method
type SigningMethod struct {
	Alg  string
	Hash crypto.Hash
}

// Sign should take a Token and a private key and generate
// the signature and return a fully encoded JWT token string.
func (s *SigningMethod) Sign(signingString string, key PrivateKey) (string, error) {
	switch s.Alg {
	case ALG_RS256, ALG_RS384, ALG_RS512:
		return s.signRSA(signingString, key)
	case ALG_HS256, ALG_HS384, ALG_HS512:
		return s.signHMAC(signingString, key)
	case ALG_ES256, ALG_ES384, ALG_ES512:
		return s.signECDSA(signingString, key)
	}
	return "", ErrInvalidKeyType
}

// Verify should take a signed JWT token string and verify
// that the signature is valid.
func (s *SigningMethod) Verify(token string, key PublicKey) error {
	switch s.Alg {
	case ALG_RS256, ALG_RS384, ALG_RS512:
		return s.verifyRSA(token, key)
	case ALG_HS256, ALG_HS384, ALG_HS512:
		return s.verifyHMAC(token, key)
	case ALG_ES256, ALG_ES384, ALG_ES512:
		return s.verifyECDSA(token, key)
	}
	return ErrInvalidKeyType
}

// RSA Signing Methods
var (
	RS256 = &SigningMethod{Alg: ALG_RS256, Hash: crypto.SHA256}
	RS384 = &SigningMethod{Alg: ALG_RS384, Hash: crypto.SHA384}
	RS512 = &SigningMethod{Alg: ALG_RS512, Hash: crypto.SHA512}
)

func (s *SigningMethod) signRSA(signingString string, key PrivateKey) (string, error) {
	// Validate type of key
	rsaKey, valid := key.(*rsa.PrivateKey)
	if !valid {
		return "", ErrInvalidKey
	}
	// Create the hasher
	if !s.Hash.Available() {
		return "", ErrHashUnavailable
	}

	hasher := s.Hash.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, s.Hash, hasher.Sum(nil))
	if err != nil {
		return "", err
	}
	return encodeSegment(sigBytes), nil
}

func (s *SigningMethod) verifyRSA(token string, key PublicKey) error {
	// Split token string
	parts, err := splitTokenString(token)
	if err != nil {
		return err
	}
	// Decode the signature
	var sig []byte
	if sig, err = decodeSegment(parts[SignaturePart]); err != nil {
		return err
	}

	var rsaKey *rsa.PublicKey
	var ok bool

	if rsaKey, ok = key.(*rsa.PublicKey); !ok {
		return ErrInvalidKeyType
	}

	// Create hasher
	if !s.Hash.Available() {
		return ErrHashUnavailable
	}
	hasher := s.Hash.New()
	hasher.Write([]byte(parts[SigningString]))

	// Verify the signature
	return rsa.VerifyPKCS1v15(rsaKey, s.Hash, hasher.Sum(nil), sig)
}

// HMAC Signing Methods
var (
	HS256 = &SigningMethod{Alg: ALG_HS256, Hash: crypto.SHA256}
	HS384 = &SigningMethod{Alg: ALG_HS384, Hash: crypto.SHA384}
	HS512 = &SigningMethod{Alg: ALG_HS512, Hash: crypto.SHA512}
)

func (s *SigningMethod) signHMAC(token *Token, key PrivateKey) (string, error) {

}

func (s *SigningMethod) verifyHMAC(token string, key PublicKey) error {

}

// ECDSA Signing Methods
var (
	ES256 = &SigningMethod{Alg: ALG_ES256, Hash: crypto.SHA256}
	ES384 = &SigningMethod{Alg: ALG_ES384, Hash: crypto.SHA384}
	ES512 = &SigningMethod{Alg: ALG_ES512, Hash: crypto.SHA512}
)

func (s *SigningMethod) signECDSA(token *Token, key PrivateKey) (string, error) {
	return "", ErrTokenSignatureInvalid
}

func (s *SigningMethod) verifyECDSA(token string, key PublicKey) error {
	return ErrTokenInvalid
}

const (
	HeaderPart    = 0
	PayloadPath   = 1
	SignaturePart = 2
	SigningString = 3
)

func splitTokenString(token string) ([]string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrTokenMalformed
	}
	parts = append(parts, strings.Join(parts[0:2], "."))
	return parts, nil
}
