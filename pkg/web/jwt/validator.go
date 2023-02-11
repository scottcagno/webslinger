package jwt

import (
	"bytes"
	"crypto"
	"errors"
	"time"
)

// Validator is the main validation structure for validating
// claims, etc.
type Validator struct {

	// Margin is an optional time margin that can be applied
	// to account for clock skew
	Margin time.Duration

	// ValidateIAT specifies whether the issued at time claim
	// will be validated.
	ValidateIAT bool

	// ExpectedAUD holds the audience this token expects; if
	// it is left as an empty string, audience validation will
	// be disabled.
	ExpectedAUD string

	// ExpectedISS holds the issuer this token expects; if
	// it is left as an empty string, issuer validation will
	// be disabled.
	ExpectedISS string

	// ExpectedSUB holds the subject this token expects; if
	// it is left as an empty string, subject validation will
	// be disabled.
	ExpectedSUB string

	Method SigningMethod
}

func (v *Validator) ValidateToken(raw RawToken, key crypto.PublicKey) (*Token, error) {

	// Create error type
	var verr error

	// Parse the initial raw rawToken
	token, err := ParseRawToken(raw)
	if err != nil {
		return nil, err
	}

	// Verify the signature method matches the provided
	// SigningMethod
	if token.Header.Alg != v.Method.Name() {
		verr = errors.Join(ErrTokenUnverifiable, err)
		return nil, verr
	}

	// ValidateRawToken the claims
	err = v.ValidateClaims(token.Payload)
	if err != nil {
		verr = errors.Join(err, ErrTokenClaimsInvalid)
		// We should continue on to validating the signature
	}

	// ValidateRawToken the final "validation" on the signature
	partialToken := raw[:bytes.LastIndexByte(raw, '.')]
	err = token.Method.Verify(partialToken, token.Signature, key)
	if err != nil {
		verr = errors.Join(err, ErrTokenSignatureInvalid)
		// continue
	}

	if verr != nil {
		return nil, verr
	}

	// We have a valid rawToken, return it!
	token.Valid = true
	return token, nil
}

func (v *Validator) ValidateClaims(claims ClaimsSet) error {

	// Create a new error
	var verr error

	// Get the current time
	now := time.Now()

	// ValidateRawToken expiration time claim (exp)
	if !v.checkExpiresAtClaim(claims.GetEXP, now) {
		verr = errors.Join(verr, ErrTokenExpired)
	}

	// ValidateRawToken not before time claim (nbf)
	if !v.checkNotBeforeClaim(claims.GetNBF, now) {
		verr = errors.Join(verr, ErrTokenNotValidYet)
	}

	// ValidateRawToken issued at time claim (iat)
	if !v.checkIssuedAtClaim(claims.GetIAT, now) {
		verr = errors.Join(verr, ErrTokenUsedBeforeIssued)
	}

	// ValidateRawToken audience claim (aud)
	if !v.checkAudienceClaim(claims.GetAUD()) {
		verr = errors.Join(verr, ErrTokenInvalidAudience)
	}

	// ValidateRawToken issuer claim (iss)
	if !v.checkIssuerClaim(claims.GetISS()) {
		verr = errors.Join(verr, ErrTokenInvalidIssuer)
	}

	// ValidateRawToken subject claim (sub)
	if !v.checkSubjectClaim(claims.GetSUB()) {
		verr = errors.Join(verr, ErrTokenInvalidSubject)
	}

	// ValidateRawToken any custom claims set that the
	// user may have implemented.
	if custom, ok := claims.(CustomClaimsSet); ok {
		if err := custom.Validate(); err != nil {
			if err != SkipValidation {
				verr = errors.Join(verr, ErrTokenInvalidCustomClaims)
				verr = errors.Join(verr, err)
			}
		}
	}

	return verr
}

func (v *Validator) checkIssuerClaim(iss string, err error) bool {
	// If expected is false or empty, skip (return true)
	if v.ExpectedISS == "" {
		return true
	}
	if err != nil && err != SkipValidation {
		return false
	}
	return iss == v.ExpectedISS
}

func (v *Validator) checkSubjectClaim(sub string, err error) bool {
	// If expected is false or empty, skip (return true)
	if v.ExpectedSUB == "" {
		return true
	}
	if err != nil && err != SkipValidation {
		return false
	}
	return sub == v.ExpectedSUB
}

func (v *Validator) checkAudienceClaim(aud string, err error) bool {
	// If expected is false or empty, skip (return true)
	if v.ExpectedAUD == "" {
		return true
	}
	if err != nil && err != SkipValidation {
		return false
	}
	return aud == v.ExpectedAUD
}

func (v *Validator) checkExpiresAtClaim(claim func() (NumericDate, error), now time.Time) bool {
	exp, err := claim()
	if err != nil && err != SkipValidation {
		return false
	}
	return now.Before(exp.Time().Add(+v.Margin))
}

func (v *Validator) checkIssuedAtClaim(claim func() (NumericDate, error), now time.Time) bool {
	if !v.ValidateIAT {
		return true
	}
	iat, err := claim()
	if err != nil && err != SkipValidation {
		return false
	}
	return !now.Before(iat.Time().Add(-v.Margin))
}

func (v *Validator) checkNotBeforeClaim(claim func() (NumericDate, error), now time.Time) bool {
	nbf, err := claim()
	if err != nil && err != SkipValidation {
		return false
	}
	return !now.Before(nbf.Time().Add(-v.Margin))
}
