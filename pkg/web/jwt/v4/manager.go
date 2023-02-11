package v4

import (
	"net/http"
	"time"
)

type TokenManager struct {
	method SigningMethod
	keys   *KeyPair
	Validator
}

func NewTokenManager(method SigningMethod, keys *KeyPair) *TokenManager {
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
	token, err := m.Validator.ValidateRawToken(raw, m.keys.PublicKey)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *TokenManager) ValidateTokenFromRequest(r *http.Request) (*Token, error) {
	var err error
	var raw RawToken
	var tok *Token

	raw, err = ExtractTokenFromRequest(r)
	if err != nil && err == ErrNoTokenInRequest {
		goto tryCookie
	}

tryCookie:
	raw, err = ExtractTokenFromCookie("token", r)
	if err != nil && err == ErrNoCookieFound {
		return nil, err
	}

	// we have a raw token, lets try to validate it
	tok, err = m.Validator.ValidateRawToken(raw, m.keys.PublicKey)
	if err != nil {
		return nil, err
	}

	return tok, nil
}
