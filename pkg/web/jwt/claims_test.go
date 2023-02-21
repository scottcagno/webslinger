package jwt

import (
	"errors"
	"testing"
	"time"
)

func TestRegisteredClaims(t *testing.T) {
	validator := &Validator{
		Margin:      5 * time.Minute,
		ValidateIAT: true,
		ExpectedAUD: "anyone",
		ExpectedISS: "me",
		ExpectedSUB: "this is a test",
	}
	claims1 := &RegisteredClaims{
		Issuer:         "me",
		Subject:        "this is a test",
		Audience:       "anyone",
		ExpirationTime: NumericDateNow().Add(12 * time.Hour),
		NotBeforeTime:  NumericDateNow(),
		IssuedAtTime:   NumericDateNow(),
		ID:             "21",
	}
	err := validator.ValidateClaims(claims1)
	if err != nil {
		t.Errorf("error validating calims: %s", err)
	}
	claims2 := &RegisteredClaims{
		Issuer:         "you",
		Subject:        "this is not a test",
		Audience:       "nobody",
		ExpirationTime: NumericDateNow().Add(12 * time.Hour),
		NotBeforeTime:  NumericDateNow(),
		IssuedAtTime:   NumericDateNow().Add(time.Hour - 1),
		ID:             "21",
	}
	err = validator.ValidateClaims(claims2)
	if err == nil {
		t.Errorf("validate claims should have failed and it didn't\n")
	}
}

type MyCustomClaims struct {
	SecretField string
	*RegisteredClaims
}

func (cc *MyCustomClaims) Validate() error {
	if cc.SecretField == "" {
		return errors.New("validation failed: secret field cannot be empty")
	}
	return nil
}

func TestRegisteredCustomClaims(t *testing.T) {
	validator := &Validator{
		Margin:      5 * time.Minute,
		ValidateIAT: true,
		ExpectedAUD: "anyone",
		ExpectedISS: "me",
		ExpectedSUB: "this is a test",
	}
	claims1 := &MyCustomClaims{
		SecretField: "abc123",
		RegisteredClaims: &RegisteredClaims{
			Issuer:         "me",
			Subject:        "this is a test",
			Audience:       "anyone",
			ExpirationTime: NumericDateNow().Add(12 * time.Hour),
			NotBeforeTime:  NumericDateNow(),
			IssuedAtTime:   NumericDateNow(),
			ID:             "21",
		},
	}
	err := validator.ValidateClaims(claims1)
	if err != nil {
		t.Errorf("error validating calims: %s", err)
	}
	claims2 := &MyCustomClaims{
		SecretField: "",
		RegisteredClaims: &RegisteredClaims{
			Issuer:         "me",
			Subject:        "this is a test",
			Audience:       "anyone",
			ExpirationTime: NumericDateNow().Add(12 * time.Hour),
			NotBeforeTime:  NumericDateNow(),
			IssuedAtTime:   NumericDateNow(),
			ID:             "21",
		},
	}
	err = validator.ValidateClaims(claims2)
	if err == nil {
		t.Errorf("validate claims should have failed and it didn't\n")
	}
}

func TestMapClaims(t *testing.T) {
	validator := &Validator{
		Margin:      5 * time.Minute,
		ValidateIAT: true,
		ExpectedAUD: "anyone",
		ExpectedISS: "me",
		ExpectedSUB: "this is a test",
	}
	claims1 := &MapClaims{
		"iss": "me",
		"sub": "this is a test",
		"aud": "anyone",
		"exp": NumericDateNow().Add(12 * time.Hour),
		"nbf": NumericDateNow(),
		"iat": NumericDateNow(),
		"jti": "21",
	}
	err := validator.ValidateClaims(claims1)
	if err != nil {
		t.Errorf("error validating calims: %s", err)
	}
	claims2 := &MapClaims{
		"iss": "you",
		"sub": "this is not a test",
		"aud": "nobody",
		"exp": NumericDateNow().Add(12 * time.Hour),
		"nbf": NumericDateNow(),
		"iat": NumericDateNow().Add(time.Hour - 1),
		"jti": "21",
	}
	err = validator.ValidateClaims(claims2)
	if err == nil {
		t.Errorf("validate claims should have failed and it didn't\n")
	}
}
