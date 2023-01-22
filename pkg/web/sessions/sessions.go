package sessions

import (
	"time"
)

type Session interface {

	// Decode should take a raw byte slice and return an expiry time,
	// session data, and any potential errors.
	Decode(b []byte) (expiry time.Time, data map[string]any, err error)

	// Encode should take an expiry time, session data, and return an
	// encoded byte slice, along with any potential errors.
	Encode(expiry time.Time, data map[string]any) ([]byte, error)
}

// SessionStore is an interface for custom session stores.
type SessionStore interface {

	// New should create and return a new session.
	New(token string) Session

	// Find should return a cached session, or an error if the session has expired or
	// cannot be found for some reason.
	Find(token string) (Session, error)

	// Save should persist the session on the underlying store. It should only return
	// an error if something goes wrong.
	Save(token string, session Session, expiry time.Time) error

	// Delete should remove the session token and corresponding session data from the
	// store. If the token does not exist, Delete should simply return nil.
	Delete(token string) error
}
