package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// token = SigningMethod(
// 		Base64URLEncode(header) + "." +
// 		Base64URLEncode(payload),
// 		secret)

type Token struct {
	RawToken
	Header    TokenHeader
	Payload   MapClaims
	Method    SigningMethod
	Signature []byte
	Valid     bool
}

func ParseRawToken(raw RawToken) (*Token, error) {

	// Create error type
	var verr, err error

	// Split the raw token
	parts := bytes.Split(raw, []byte{'.'})
	if len(parts) != 3 {
		return nil, ErrTokenMalformed
	}

	// Initialize a token instance
	var token Token
	token.RawToken = raw

	// Parse the raw token header
	headerBytes := Base64Decode(parts[0])
	err = json.Unmarshal(headerBytes, &token.Header)
	if err != nil {
		verr = errors.Join(ErrTokenMalformed, err)
		return nil, verr
	}

	// Parse the raw token claims
	claimsBytes := Base64Decode(parts[1])
	err = json.Unmarshal(claimsBytes, &token.Payload)
	if err != nil {
		verr = errors.Join(ErrTokenMalformed, err)
		return nil, verr
	}

	// Get the signing method (even though it's not verifiable, yet)
	token.Method = GetSigningMethod(token.Header.Alg)
	if token.Method == nil {
		verr = errors.Join(ErrTokenMalformed, err)
		return nil, verr
	}

	// Add the signature
	token.Signature = parts[2]

	// Return out token (unverified raw)
	return &token, verr
}

func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(src)))
	base64.RawURLEncoding.Encode(buf, src)
	return buf
}

func Base64Decode(src []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(src)))
	n, err := base64.RawURLEncoding.Decode(buf, src)
	if err != nil {
		panic(err)
	}
	return buf[:n]
}
