package v1

import (
	"errors"
)

// Error constants
var (
	ErrInvalidKey      = errors.New("key is invalid")
	ErrInvalidKeyType  = errors.New("key is of invalid type")
	ErrHashUnavailable = errors.New("the requested hash function is unavailable")

	ErrTokenMalformed        = errors.New("token is malformed")
	ErrTokenUnverifiable     = errors.New("token is unverifiable")
	ErrTokenSignatureInvalid = errors.New("token signature is invalid")

	ErrTokenInvalid          = errors.New("token is invalid (could be anything)")
	ErrTokenInvalidAudience  = errors.New("token has invalid audience")
	ErrTokenExpired          = errors.New("token is expired")
	ErrTokenUsedBeforeIssued = errors.New("token used before issued")
	ErrTokenInvalidIssuer    = errors.New("token has invalid issuer")
	ErrTokenNotValidYet      = errors.New("token is not valid yet")
	ErrTokenInvalidId        = errors.New("token has invalid id")
	ErrTokenInvalidClaims    = errors.New("token has invalid claims")
)

// The errors that might occur when parsing and validating a token
const (
	ErrValidationMalformed uint32 = 1 << iota
	ErrValidationUnverifiable
	ErrValidationSignatureInvalid

	ErrValidationAudience      // AUD validation failed
	ErrValidationExpired       // EXP validation failed
	ErrValidationIssuedAt      // IAT validation failed
	ErrValidationIssuer        // ISS validation failed
	ErrValidationNotValidYet   // NBF validation failed
	ErrValidationId            // JTI validation failed
	ErrValidationClaimsInvalid // Generic claims validation error
)
