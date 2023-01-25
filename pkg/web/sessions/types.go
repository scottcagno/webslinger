package sessions

import (
	"errors"
	"time"
)

// ErrSessionNotFound is the error used with the Find method on the SessionStore interface
var ErrSessionNotFound = errors.New("session store: session data was not found or has expired")

// Codec is the interface for encoding and decoding session
// data to and from a byte slice for use by the session store.
type Codec interface {

	// Encode should take an expiry time, a session data map and
	// return an encoded byte slice or any errors encountered.
	Encode(time.Time, map[string]any) ([]byte, error)

	// Decode should take a raw byte slice and return an expiry time,
	// a session data map, or any errors encountered with decoding.
	Decode(b []byte) (time.Time, map[string]any, error)
}

// SessionStore is an interface for custom session stores.
type SessionStore interface {

	// Find should return the data for a session token from the underlying store. A
	// nil error should be returned unless the data cannot be found, is expired or
	// malformed in which case ErrSessionNotFound should be returned.
	Find(token string) ([]byte, error)

	// Save should persist the session to the underlying store with the provided expiry
	// time. It should only return an error if something goes wrong.
	Save(token string, b []byte, expiry time.Time) error

	// Delete should remove the session token and corresponding session data from the
	// store. If the token does not exist, Delete should simply return nil.
	Delete(token string) error
}
