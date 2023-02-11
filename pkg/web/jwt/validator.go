package jwt

import (
	"errors"
	"time"
)

// validator is the main validation structure for validating
// claims, etc.
type validator struct {

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

func (v *validator) ValidateClaims(claims ClaimsSet) error {

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

func (v *validator) checkIssuerClaim(iss string, err error) bool {
	// If expected is false or empty, skip (return true)
	if v.ExpectedISS == "" {
		return true
	}
	if err != nil && err != SkipValidation {
		return false
	}
	return iss == v.ExpectedISS
}

func (v *validator) checkSubjectClaim(sub string, err error) bool {
	// If expected is false or empty, skip (return true)
	if v.ExpectedSUB == "" {
		return true
	}
	if err != nil && err != SkipValidation {
		return false
	}
	return sub == v.ExpectedSUB
}

func (v *validator) checkAudienceClaim(aud string, err error) bool {
	// If expected is false or empty, skip (return true)
	if v.ExpectedAUD == "" {
		return true
	}
	if err != nil && err != SkipValidation {
		return false
	}
	return aud == v.ExpectedAUD
}

func (v *validator) checkExpiresAtClaim(claim func() (NumericDate, error), now time.Time) bool {
	exp, err := claim()
	if err != nil && err != SkipValidation {
		return false
	}
	return now.Before(exp.Time().Add(+v.Margin))
}

func (v *validator) checkIssuedAtClaim(claim func() (NumericDate, error), now time.Time) bool {
	if !v.ValidateIAT {
		return true
	}
	iat, err := claim()
	if err != nil && err != SkipValidation {
		return false
	}
	return !now.Before(iat.Time().Add(-v.Margin))
}

func (v *validator) checkNotBeforeClaim(claim func() (NumericDate, error), now time.Time) bool {
	nbf, err := claim()
	if err != nil && err != SkipValidation {
		return false
	}
	return !now.Before(nbf.Time().Add(-v.Margin))
}
