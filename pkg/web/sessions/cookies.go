package sessions

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// NewCookie is a helper that wraps the creation of a new cookie
// and returns a filled out *http.Cookie instance that can be
// modified if need be. This is meant to just set up some basic
// defaults.
func NewCookie(name, value, domain string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:       name,
		Value:      value,
		Path:       "/",
		Domain:     domain,
		Expires:    expires,
		RawExpires: "",
		MaxAge:     CookieMaxAge(expires),
		Secure:     false,                   // set to true, if using TLS (false otherwise)
		HttpOnly:   true,                    // protects against XSS attacks
		SameSite:   http.SameSiteStrictMode, // protects against CSRF attacks
		Raw:        "",
		Unparsed:   nil,
	}
}

// HasCookie checks if there is an existing *http.Cookie in the
// *http.Request that is associated with the provided name.
func HasCookie(r *http.Request, name string) bool {
	_, err := r.Cookie(name)
	return err == nil
}

// GetCookie retrieves a cookie from an HTTP request by its name,
// and returns the cookie. It is just a simple wrapper around r.Cookie
// and is really only meant to provide this API with a future forward
// way of introducing logic into returning a cookie if the need ever
// comes up.
//
// If the cookie is not found, it returns a nil cookie along with an
// error of http.ErrNoCookie otherwise, it will return a valid cookie
// and a nil error.
func GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	c, err := r.Cookie(name)
	if err != nil && err == http.ErrNoCookie {
		return nil, err
	}
	return c, nil
}

// TimeUntilExpires takes an expiration time and returns the
// remaining duration until the expiration time, minus 1 second.
// If the expiration time has already passed, it returns 0.
func TimeUntilExpires(expires time.Time) time.Duration {
	max := time.Until(expires)
	if max < 0 {
		return 0
	}
	return max - time.Second
}

// CookieMaxAge takes an expiration time and returns the
// remaining duration until the expiration time in seconds.
// If the expiration time has already passed, it returns -1
// indicating that the cookie should be removed.
func CookieMaxAge(expires time.Time) int {
	max := TimeUntilExpires(expires)
	if max == 0 {
		return -1
	}
	return int(max.Seconds())
}

// Base64Encode takes a plaintext string and returns a base64 encoded string
func Base64Encode(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

// Base64Decode takes a base64 encoded string and returns a plaintext string
func Base64Decode(s string) string {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		panic(fmt.Sprintf("base64 deocde: failed %q", err))
	}
	return string(b)
}

// URLEncode takes a plaintext string and returns a URL encoded string
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// URLDecode takes a URL encoded string and returns a plaintext string
func URLDecode(s string) string {
	us, err := url.QueryUnescape(s)
	if err != nil {
		panic(fmt.Sprintf("url decode: failed %q", err))
	}
	return us
}
