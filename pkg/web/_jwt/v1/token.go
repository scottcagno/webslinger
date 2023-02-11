package v1

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Header represets the token header
type Header map[string]any

// Claims represents the typical registered claims
type Claims map[string]any

// Token is a parsed JSON Web Token.
type Token struct {
	Method *SigningMethod
	Header map[string]any
	Claims map[string]any
	Valid  bool
}

func NewToken(method *SigningMethod, claims map[string]any) *Token {
	return &Token{
		Method: method,
		Header: map[string]any{
			"typ": "JWT",
			"alg": method.Alg,
		},
		Claims: claims,
	}
}

// Sign is an interface for signing a JWT token, and return the JWT token
// as a signed encoded JWT token string.
func (t *Token) Sign(key PrivateKey) (string, error) {
	// Marshal the token header
	b, err := json.Marshal(t.Header)
	if err != nil {
		return "", err
	}
	// Marshal the token payload
	b, err = json.Marshal(t.Claims)
	if err != nil {
		return "", err
	}
	// Encode the header and the payload
	header := encodeSegment(b)
	claims := encodeSegment(b)
	// Pass the signing string into the signing method
	sig, err := t.Method.Sign(strings.Join([]string{header, claims}, "."), key)
	if err != nil {
		return "", err
	}
	_ = sig
	// nope..
	return "", nil
}

func (t *Token) Verify(token string) error {
	// nope...
	return nil
}

// encodeSegment encodes a JWT specific base64url encoding with padding stripped.
func encodeSegment(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}

// decodeSegment decodes a JWT specific base64url encoding with padding stripped.
func decodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}
