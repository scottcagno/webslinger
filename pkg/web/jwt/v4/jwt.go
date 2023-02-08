package v4

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
)

// token = SigningMethod(
// 		Base64URLEncode(header) + "." +
// 		Base64URLEncode(payload),
// 		secret)

type Token []byte

func NewToken(alg SigningMethod, claims ClaimsSet, key crypto.PrivateKey) (Token, error) {
	// create and encode the header
	dat, err := json.Marshal(
		struct {
			Type string `json:"typ"`
			Alg  string `json:"alg"`
		}{
			Type: "JWT",
			Alg:  alg.Name(),
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
	Claims    ClaimsSet
	Method    SigningMethod
	Signature []byte
}

func ValidateClaims(token Token) (ClaimsSet, error) {
	// TODO: implement
	panic("implement me")

	// 1. split token
	// 2. parse header
	// 3. validate basic claims
	// 4. assemble and return claims
}

func Validate(token Token) (*ValidToken, error) {
	// TODO: implement
	panic("implement me")

	// 1. split token
	// 2. verify correct signing method
	// 3. check for key
	// 4. validate basic claims
	// 5. validate using the signing method
	// 6. assemble and return valid token
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
