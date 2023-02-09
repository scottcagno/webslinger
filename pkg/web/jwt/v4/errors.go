package v4

import (
	"errors"
)

// SigningMethod errors
var (
	ErrInvalidKeyType   = errors.New("alg: invalid key type")
	ErrHashUnavailable  = errors.New("alg: hash unavailable")
	ErrSignatureInvalid = errors.New("alg: signature is invalid")
)

// Validator errors
var (
	ErrTokenExpired             = errors.New("token expired")
	ErrTokenNotValidYet         = errors.New("token not valid yet")
	ErrTokenUsedBeforeIssued    = errors.New("token used before issued")
	ErrTokenInvalidAudience     = errors.New("token contains invalid audience")
	ErrTokenInvalidIssuer       = errors.New("token contains invalid issuer")
	ErrTokenInvalidSubject      = errors.New("token contains invalid subject")
	ErrTokenInvalidCustomClaims = errors.New("token contains invalid custom claims")
)

// Parser errors
var (
	ErrTokenMalformed        = errors.New("token is malformed")
	ErrTokenUnverifiable     = errors.New("token is unverifiable")
	ErrTokenClaimsInvalid    = errors.New("token claims validation error")
	ErrTokenSignatureInvalid = errors.New("token validation signature invalid")
)
