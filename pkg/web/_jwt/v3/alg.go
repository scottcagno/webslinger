package v3

import (
	"crypto"
)

type algRSA struct {
	name string
	hash crypto.Hash
}

var (
	AlgRSA256 = &algRSA{"RS256", crypto.SHA256}
	AlgRSA384 = &algRSA{"RS384", crypto.SHA384}
	AlgRSA512 = &algRSA{"RS512", crypto.SHA512}
)

func (a *algRSA) Name() string {
	return a.name
}

func (a *algRSA) Sign(key crypto.PrivateKey, input []byte) ([]byte, error) {
	// TODO: implement me..
	panic("implement me")
}

func (a *algRSA) Verify() {

}
