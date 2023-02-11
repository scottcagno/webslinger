package jwt

import (
	"net/http"
	"testing"
	"time"
)

func TestNewTokenManager(t *testing.T) {

	// RSA signing methods
	_ = RS256
	_ = RS384
	_ = RS512

	// HMAC signing methods
	_ = HS256
	_ = HS384
	_ = HS512

	// ECDSA signing methods
	_ = ES256
	_ = ES384
	_ = ES512

	// RSAPSS signing methods
	_ = PS256
	_ = PS384
	_ = PS512

	// generate key pair using RSA 512
	keys := RS512.GenerateKeyPair()

	// create a new token manager instance
	tm := NewTokenManager(RS512, keys)

	// generate a new (signed) token using default claims
	t1, err := tm.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	// see if we can validate the token
	_, err = tm.ValidateToken(t1)
	if err != nil {
		t.Fatal(err)
	}

	// generate a new signed token with custom claims
	t2, err := tm.GenerateToken(
		&RegisteredClaims{
			Issuer:         "jon doe",
			Subject:        "your mom goes to college",
			Audience:       "anyone",
			ExpirationTime: NumericDate(time.Now().Add(4 * time.Hour).Unix()),
			NotBeforeTime:  NumericDate(time.Now().Unix()),
			IssuedAtTime:   NumericDate(time.Now().Unix()),
			ID:             "25",
		},
	)

	// see if we can validate the token
	_, err = tm.ValidateToken(t2)
	if err != nil {
		t.Fatal(err)
	}

}

func TestTokenManager_ValidateTokenFromRequest(t *testing.T) {
	// generate key pair using RSA 512
	keys := RS512.GenerateKeyPair()

	// create a new token manager instance
	tm := NewTokenManager(RS512, keys)

	handler := func(w http.ResponseWriter, r *http.Request) {

		tok, err := tm.ValidateTokenFromRequest(r)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		_ = tok
	}

	_ = handler
}
