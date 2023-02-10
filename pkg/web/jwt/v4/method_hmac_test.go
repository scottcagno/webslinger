package v4

import (
	"bytes"
	"crypto"
	"fmt"
	"testing"
)

var hmacTestKey = HMACPrivateKey("0hmac1super2secret3code4right5here6buddy7so8what9you10gonna11do12about13it14")

var testClaims = RegisteredClaims{
	Issuer:         "joe",
	Subject:        "history",
	Audience:       "your mom",
	ExpirationTime: 1300819380,
	NotBeforeTime:  0,
	IssuedAtTime:   0,
	ID:             "",
}

var hmacTestData = []struct {
	name   string
	token  RawToken
	alg    string
	claims RegisteredClaims
	valid  bool
}{
	{
		"HS256",
		RawToken(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJleHAiOjEzMDA4MTkzODAsInN1YiI6Imhpc3RvcnkiLCJhdWQiOiJ5b3VyIG1vbSIsImlzcyI6ImpvZSJ9." +
				"7hAzBeo7T9U9VZXHOfmER0UvE5tG7Ar1rFPl9qjs7V0",
		),
		"HS256",
		testClaims,
		true,
	},
	{
		"HS384",
		RawToken(
			"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9." +
				"eyJleHAiOjEzMDA4MTkzODAsInN1YiI6Imhpc3RvcnkiLCJhdWQiOiJ5b3VyIG1vbSIsImlzcyI6ImpvZSJ9." +
				"PkH3Av3Gj5UT1WYYNBVyf0twwXS_M_IbrvfURENp6bzC9tj_jzg0Hcz7RO1YEn1u",
		),
		"HS384",
		testClaims,
		true,
	},
	{
		"HS512",
		RawToken(
			"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
				"eyJleHAiOjEzMDA4MTkzODAsInN1YiI6Imhpc3RvcnkiLCJhdWQiOiJ5b3VyIG1vbSIsImlzcyI6ImpvZSJ9." +
				"6AR9RUeUrxh_5WrwBND-lrAZCKOWpur17sbixhaos_rjXhQqysF7aI2e0qAHD1AqOFWCdjTzcHoiwUU9b8IQAA",
		),
		"HS512",
		testClaims,
		true,
	},
	{
		"web sample: invalid",
		RawToken(
			"eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9." +
				"eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ." +
				"dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXo",
		),
		"HS256",
		testClaims,
		false,
	},
}

func TestSigningMethodHMAC_Sign(t *testing.T) {
	for _, data := range hmacTestData {
		if !data.valid {
			continue
		}

		token := data.token
		header := token.Header()
		signingstr := token.SigningSection()
		signature := token.Signature()
		fmt.Printf(
			"token=\n\n%q\nheader=%q\nsigningstr=%q\nsignature=%q\n\n",
			token, header, signingstr, signature,
		)

		method := GetSigningMethod(header.Alg)

		fmt.Printf(">>> GetSigningMethod=%+v\n", method)

		sig, err := method.Sign(signingstr, crypto.PrivateKey(hmacTestKey))
		if err != nil {
			t.Errorf("[%v] Error signing token: %v\n", data.name, err)
		}
		if !bytes.Equal(sig, signature) {
			t.Errorf(
				"[%v] Incorrect signature:\nwanted...\n%v\ngot...\n%v\n",
				data.name, signature, sig,
			)
		}
	}
}

func TestSigningMethodHMAC_Verify(t *testing.T) {
	for _, data := range hmacTestData {

		token := data.token
		header := token.Header()
		signingstr := token.SigningSection()
		signature := token.Signature()

		method := GetSigningMethod(header.Alg)

		fmt.Printf(">>> GetSigningMethod=%+v\n", method)

		err := method.Verify(signingstr, signature, hmacTestKey)

		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %v\n", data.name, err)
		}
		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation\n", data.name)
		}
	}
}

// func BenchmarkHS256Signing(b *testing.B) {
// 	benchmarkSigning(b, jwt.SigningMethodHS256, hmacTestKey)
// }
//
// func BenchmarkHS384Signing(b *testing.B) {
// 	benchmarkSigning(b, jwt.SigningMethodHS384, hmacTestKey)
// }
//
// func BenchmarkHS512Signing(b *testing.B) {
// 	benchmarkSigning(b, jwt.SigningMethodHS512, hmacTestKey)
// }
