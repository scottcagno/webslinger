package jwt

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractTokenFromRequest(r *http.Request) (RawToken, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return nil, ErrNoTokenInRequest
	}
	return RawToken(authHeader[7:]), nil
}

func ExtractTokenFromCookie(name string, r *http.Request) (RawToken, error) {
	c, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		return nil, ErrNoCookieFound
	}
	return RawToken(c.Value), nil
}
