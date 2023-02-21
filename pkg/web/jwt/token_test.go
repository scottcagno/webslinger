package jwt

import (
	"bytes"
	"testing"
)

var rawToken = RawToken(
	"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
		"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
)

var validToken = &Token{
	RawToken: rawToken,
	Header: TokenHeader{
		Typ: "JWT",
		Alg: "HS256",
	},
	Payload: MapClaims{
		"iat":  NumericDate(1516239022),
		"name": "John Doe",
		"sub":  "1234567890",
	},
	Method:    HS256,
	Signature: []byte("SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"),
	Valid:     true,
}

func assert(t *testing.T, wanted, got any) {
	if wanted != got {
		t.Errorf("wanted=%s, got=%s\n", wanted, got)
	}
}

func TestParseRawToken(t *testing.T) {
	tok, err := ParseRawToken(rawToken)
	if err != nil {
		t.Errorf("error parsing raw token: %s", err)
	}
	if !bytes.Equal(tok.SigningSection(), validToken.SigningSection()) {
		t.Errorf("parsed token does not match valid token")
	}
}
