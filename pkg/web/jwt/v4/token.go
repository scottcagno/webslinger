package v4

import (
	"bytes"
	"crypto"
	"encoding/json"
)

type Section uint8

const (
	HeaderSection    Section = 0
	PayloadSection   Section = 1
	ClaimsSection    Section = 1
	SignatureSection Section = 2
	PartialToken     Section = 3
	SigningString    Section = 3
	dot              byte    = '.'
)

type TokenHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type RawToken []byte

// GetSection returns the start and end index for section
// containing the header, the payload or the signature.
func (t RawToken) GetSection(section Section) (int, int) {
	switch section {
	case HeaderSection:
		return 0, bytes.IndexByte(t, dot)
	case PayloadSection | ClaimsSection:
		return bytes.IndexByte(t, dot), bytes.LastIndexByte(t, dot)
	case SignatureSection:
		return bytes.LastIndexByte(t, dot) + 1, len(t)
	case PartialToken | SigningString:
		return 0, bytes.LastIndexByte(t, dot)
	}
	return -1, -1
}

func (t RawToken) Header() *TokenHeader {
	var header TokenHeader
	beg, end := t.GetSection(HeaderSection)
	b := Base64Decode(t[beg:end])
	err := json.Unmarshal(b, &header)
	if err != nil {
		return nil
	}
	return &header
}

func (t RawToken) Claims() ClaimsSet {
	var claims ClaimsSet
	i, j := t.GetSection(ClaimsSection)
	b := Base64Decode(t[i:j])
	err := json.Unmarshal(b, &claims)
	if err != nil {
		return nil
	}
	return claims
}

func (t RawToken) SigningSection() []byte {
	i, j := t.GetSection(SigningString)
	b := make([]byte, j)
	copy(b, t[i:j])
	return b
}

func (t RawToken) Signature() []byte {
	i, j := t.GetSection(SignatureSection)
	b := make([]byte, len(t[i:j]))
	copy(b, t[i:j])
	return b
}

func NewToken(alg SigningMethod, claims ClaimsSet, key crypto.PrivateKey) (RawToken, error) {
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
