package jwt

import (
	"fmt"
	"testing"
)

var rawToken = RawToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
	"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
	"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

func TestParseRawToken(t *testing.T) {

	validToken, err := ParseRawToken(rawToken)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v\n", validToken)
}
