package v4

import (
	"errors"
	"time"
)

type NumericDate int64

func (n NumericDate) Time() time.Time {
	return time.Unix(int64(n), 0)
}

// SkipValidation can be  used as a return value from ValidateClaimFunc to
// indicate that the claim in the call is to be skipped. It is not returned
// as an error by any function.
var SkipValidation = errors.New("skip validation for this claim")

// ValidateClaimFunc is the type of function called by ValidateClaims in order
// to validate a claims set. It makes it possible to implement custom claims
// validation. If you do not want a claim to be validated the SkipValidation
// error is returned.
type ValidateClaimFunc func(name string, claim any) error

// ClaimsSet is an interface for a set of claims.
type ClaimsSet interface {

	// GetISS is the issuer of the JWT.
	GetISS() (string, error)

	// GetSUB is the subject of the JWT.
	GetSUB() (string, error)

	// GetAUD is the audience (Recipient for which the JWT is intended.)
	GetAUD() (string, error)

	// GetEXP is the time after which the JWT expires.
	GetEXP() (NumericDate, error)

	// GetNBF is the time before the JWT must not be accepted for processing.
	GetNBF() (NumericDate, error)

	// GetIAT is the time at which the JWT was issued.
	// It can be used to determine the age of the JWT.
	GetIAT() (NumericDate, error)

	// GetJTI is a unique identifier. It can be used to prevent the JWT
	// from being replayed; it allows a token to be used only once.
	GetJTI() (string, error)
}

// CustomClaimsSet is an interface for representing a custom claim, or a
// set of custom claims. It is intended to be used for building custom
// claim validation in addition to the integrated ClaimsSet.
type CustomClaimsSet interface {
	ClaimsSet
	Validate() error
}

// RegisteredClaims is the default set of registered claims. It can be used
// in addition to custom claims by embedding the registered claims in
type RegisteredClaims struct {
	Issuer         string      `json:"iss,omitempty"`
	Subject        string      `json:"sub,omitempty"`
	Audience       string      `json:"aud,omitempty"`
	ExpirationTime NumericDate `json:"exp,omitempty"`
	NotBeforeTime  NumericDate `json:"nbf,omitempty"`
	IssuedAtTime   NumericDate `json:"iat,omitempty"`
	ID             string      `json:"jti,omitempty"`
}

func (r *RegisteredClaims) GetISS() (string, error) {
	return r.Issuer, nil
}

func (r *RegisteredClaims) GetSUB() (string, error) {
	return r.Subject, nil
}

func (r *RegisteredClaims) GetAUD() (string, error) {
	return r.Audience, nil
}

func (r *RegisteredClaims) GetEXP() (NumericDate, error) {
	return r.ExpirationTime, nil
}

func (r *RegisteredClaims) GetNBF() (NumericDate, error) {
	return r.NotBeforeTime, nil
}

func (r *RegisteredClaims) GetIAT() (NumericDate, error) {
	return r.IssuedAtTime, nil
}

func (r *RegisteredClaims) GetJTI() (string, error) {
	return r.ID, nil
}
