package v2

import (
	"crypto"
)

type VerifiedToken struct {
	Token     []byte
	Header    []byte
	Payload   []byte
	Signature []byte
	Claims    []byte
}

type Token struct {
	Header    []byte
	Payload   []byte
	Signature []byte
}

func Verify(alg Alg, k crypto.PublicKey, token []byte) (*VerifiedToken, error) {
	if len(token) == 0 {
		return nil, ErrMissing
	}
	tok, err := decode(alg, k, token)
	if err != nil {
		return nil, err
	}

	// TODO: finish this, maybe
	_ = tok
	return nil, nil
}
