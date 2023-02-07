package v3

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
)

// Algorithm represents a signing and verifying type
type Algorithm interface {
	// Name should return the string representation of the
	// JWT "alg" field
	Name() string

	// Sign should take the private key and a base64 encoded header
	// and payload data and return the signature.
	// Sign(header, payload []byte) ([]byte, error)
	Sign(key crypto.PrivateKey, input []byte) ([]byte, error)

	// Verify should take the public key, a base64 encoded header,
	// payload, and signature and verify the signature against
	// the header and payload.
	Verify(key crypto.PublicKey, header, payload, signature []byte) error
}

// NumericDate is a just a "UTC Unix Timestamp" (seconds since the epoch)
type NumericDate = int64

type ClaimsSet interface {
	GetEXP() (NumericDate, error)
	GetIAT() (NumericDate, error)
	GetNBF() (NumericDate, error)
	GetISS() (string, error)
	GetSUB() (string, error)
	GetAUD() ([]string, error)
}

type Token []byte

func NewToken(alg Algorithm, claims ClaimsSet, key crypto.PrivateKey) (Token, error) {

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
	input := bytes.Join([][]byte{header, payload}, []byte{'.'})

	// create and encode the signature
	dat, err = alg.Sign(key, input)
	if err != nil {
		return nil, err
	}
	signature := Base64Encode(dat)

	// create and return token
	return bytes.Join([][]byte{input, signature}, []byte{'.'}), nil
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
