package v4

import (
	"time"
)

type TokenManager struct {
	method SigningMethod
	keys   KeyPair
	Validator
}

func NewTokenManager(method SigningMethod, keys KeyPair) *TokenManager {
	m := &TokenManager{
		method: method,
		keys:   keys,
	}
	m.Validator = Validator{
		Margin:      time.Minute,
		ValidateIAT: false,
		ExpectedAUD: "",
		ExpectedISS: "",
		ExpectedSUB: "",
		Method:      method,
	}
	return m
}

func (m *TokenManager) GenerateToken(claims ClaimsSet) (RawToken, error) {
	token, err := NewToken(m.method, claims, m.keys.PrivateKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *TokenManager) ValidateToken(raw RawToken) (*Token, error) {
	token, err := m.Validator.Validate(raw, m.keys.PublicKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}
