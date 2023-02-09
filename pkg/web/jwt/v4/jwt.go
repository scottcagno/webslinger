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

type TokenHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type Token []byte

func NewToken(alg SigningMethod, claims ClaimsSet, key crypto.PrivateKey) (Token, error) {
	// create and encode the header
	dat, err := json.Marshal(
		TokenHeader{
			Typ: "JWT",
			Alg: alg.Name(),
		},
	)
	if err != nil {
		return nil, err
	}
	header := Base64Encode(dat)
	// create and encode the payload
	dat, err = json.Marshal(claims)
	if err != nil {
		return nil, err
	}
	payload := Base64Encode(dat)
	// create signing string input
	partialToken := bytes.Join([][]byte{header, payload}, []byte{'.'})
	// create and encode the signature
	dat, err = alg.Sign(partialToken, key)
	if err != nil {
		return nil, err
	}
	// create and return token
	return bytes.Join([][]byte{partialToken, dat}, []byte{'.'}), nil
}

type ValidToken struct {
	Raw       Token
	Header    TokenHeader
	Payload   ClaimsSet
	Method    SigningMethod
	Signature []byte
	Valid     bool
}

func (v *Validator) Validate(token Token, claims ClaimsSet, key crypto.PublicKey) (*ValidToken, error) {

	// Create error type
	var verr error

	// Parse the initial raw token
	raw, err := ParseRawToken(token)
	if err != nil {
		return nil, err
	}

	// Verify the signature method matches the provided
	// SigningMethod
	if raw.Header.Alg != v.Method.Name() {
		verr = errors.Join(ErrTokenUnverifiable, err)
		return nil, verr
	}

	// Validate the claims
	err = v.ValidateClaims(claims)
	if err != nil {
		verr = errors.Join(err, ErrTokenClaimsInvalid)
		// We should continue on to validating the signature
	}

	// Validate the final "validation" on the signature
	partialToken := token[:bytes.LastIndexByte(token, '.')]
	err = raw.Method.Verify(partialToken, raw.Signature, key)
	if err != nil {
		verr = errors.Join(err, ErrTokenSignatureInvalid)
		// continue
	}

	if verr != nil {
		return nil, verr
	}

	// We have a valid token, return it!
	raw.Valid = true
	return raw, nil
}

func ParseRawToken(token Token) (*ValidToken, error) {

	// Create error type
	var verr, err error

	// Split the raw token
	parts := bytes.Split(token, []byte{'.'})
	if len(parts) != 3 {
		return nil, ErrTokenMalformed
	}

	// Initialize a raw token instance
	var raw ValidToken
	raw.Raw = token

	// Parse the token header
	headerBytes := Base64Decode(parts[0])
	err = json.Unmarshal(headerBytes, &raw.Header)
	if err != nil {
		verr = errors.Join(ErrTokenMalformed, err)
		return nil, verr
	}

	// Parse the token claims
	claimsBytes := Base64Decode(parts[1])
	err = json.Unmarshal(claimsBytes, &raw.Payload)
	if err != nil {
		verr = errors.Join(ErrTokenMalformed, err)
		return nil, verr
	}

	// Get the signing method (even though it's not verifiable, yet)
	raw.Method = GetSigningMethod(raw.Header.Alg)
	if raw.Method == nil {
		verr = errors.Join(ErrTokenMalformed, err)
		return nil, verr
	}

	// Add this signature
	raw.Signature = parts[2]

	// Return out raw (unverified token)
	return &raw, verr
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
