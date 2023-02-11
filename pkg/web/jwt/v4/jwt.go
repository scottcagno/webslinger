package v4

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// token = SigningMethod(
// 		Base64URLEncode(header) + "." +
// 		Base64URLEncode(payload),
// 		secret)

type Token struct {
	Raw       RawToken
	Header    TokenHeader
	Payload   *RegisteredClaims
	Method    SigningMethod
	Signature []byte
	Valid     bool
}

func (v *Validator) ValidateRawToken(raw RawToken, key crypto.PublicKey) (*Token, error) {

	// Create error type
	var verr error

	// Parse the initial raw rawToken
	token, err := ParseRawToken(raw)
	if err != nil {
		return nil, err
	}

	// Verify the signature method matches the provided
	// SigningMethod
	if token.Header.Alg != v.Method.Name() {
		verr = errors.Join(ErrTokenUnverifiable, err)
		return nil, verr
	}

	// ValidateRawToken the claims
	err = v.ValidateClaims(token.Payload)
	if err != nil {
		verr = errors.Join(err, ErrTokenClaimsInvalid)
		// We should continue on to validating the signature
	}

	// ValidateRawToken the final "validation" on the signature
	partialToken := raw[:bytes.LastIndexByte(raw, '.')]
	err = token.Method.Verify(partialToken, token.Signature, key)
	if err != nil {
		verr = errors.Join(err, ErrTokenSignatureInvalid)
		// continue
	}

	if verr != nil {
		return nil, verr
	}

	// We have a valid rawToken, return it!
	token.Valid = true
	return token, nil
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
	token.Raw = raw

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
