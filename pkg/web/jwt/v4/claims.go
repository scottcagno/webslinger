package v4

type NumericDate = int64

type ClaimsSet interface {

	// GetISS is the issuer of the JWT.
	GetISS() (string, error)

	// GetSUB is the subject of the JWT.
	GetSUB() (string, error)

	// GetAUD is the audience (Recipient for which the JWT is intended.)
	GetAUD() ([]string, error)

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

type RegisteredClaims struct {
	Issuer         string   `json:"iss,omitempty"`
	Subject        string   `json:"sub,omitempty"`
	Audience       []string `json:"aud,omitempty"`
	ExpirationTime int64    `json:"exp,omitempty"`
	NotBeforeTime  int64    `json:"nbf,omitempty"`
	IssuedAtTime   int64    `json:"iat,omitempty"`
	JWTID          string   `json:"jti,omitempty"`
}

func (r *RegisteredClaims) GetISS() (string, error) {
	return r.Issuer, nil
}

func (r *RegisteredClaims) GetSUB() (string, error) {
	return r.Subject, nil
}

func (r *RegisteredClaims) GetAUD() ([]string, error) {
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
	return r.JWTID, nil
}
