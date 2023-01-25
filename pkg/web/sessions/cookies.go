package sessions

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// CookieConfig contains the configuration settings for session cookies.
type CookieConfig struct {

	// Name sets the name of the session cookie. It should not contain
	// whitespace, commas, colons, semicolons, backslashes, the equals
	// sign or control characters as per RFC6265. The default cookie name
	// is "SESSION". If your application uses two different sessions, you
	// must make sure that the cookie name for each is unique.
	Name string

	// Path sets the 'Path' attribute on the session cookie. The default
	// value is "/". Passing the empty string "" will result in it being
	// set to the path that the cookie was issued from.
	Path string

	// Domain sets the 'Domain' attribute on the session cookie. By default,
	// it will be set to the domain name that the cookie was issued from.
	Domain string

	// Secure sets the 'Secure' attribute on the session cookie. The default
	// value is false. It's recommended that you set this to true and serve
	// all requests over HTTPS in production environments.
	Secure bool

	// HttpOnly sets the 'HttpOnly' attribute on the session cookie. The
	// default value is true. Having HttpOnly set to true can help prevent
	// against XSS attacks.
	HttpOnly bool

	// SameSite controls the value of the 'SameSite' attribute on the session
	// cookie. By default, this is set to http.SameSiteLaxMode. If you want
	// no SameSite attribute or value in the session cookie then you should
	// set this to 0. Available options for the 'SameSite' attribute are
	// SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode and SameSiteNoneMode.
	// Having SameSite set to http.SameSiteStrictMode can help protect against
	// CSRF attacks.
	SameSite http.SameSite

	// Persist sets whether the session cookie should be retained after a
	// user closes their browser (default value is true.) The appropriate
	// 'Expires' and 'MaxAge' values will be added to the session cookie.
	Persist bool
}

// NewCookie is a helper that wraps the creation of a new cookie
// and returns a filled out *http.Cookie instance that can be
// modified if need be. This is meant to just put up some basic
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
		Secure:     false,                // put to true, if using TLS (false otherwise)
		HttpOnly:   true,                 // protects against XSS attacks
		SameSite:   http.SameSiteLaxMode, // protects against CSRF attacks
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
