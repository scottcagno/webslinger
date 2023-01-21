package sessions

import (
	"time"
)

// Session is an interface for custom session types.
type SessionT interface {

	// ID should return the ID for the session.
	ID() string

	// Expires takes a time, and sets the expiry time for the session. It can be
	// used to update and keep the session alive.
	Expires(at time.Time)

	// ExpiresIn should return the remaining time left until the session is due
	// to expire.
	ExpiresIn() time.Duration
}

// SessionStore is an interface for custom session stores.
type SessionStoreT interface {

	// New should create and return a new session.
	New(token string) Session

	// Find should return a cached session, or an error if the session has expired or
	// cannot be found for some reason.
	Find(token string) (Session, error)

	// Save should persist the session on the underlying store. It should only return
	// an error if something goes wrong.
	Save(token string, session Session) error

	// Delete should remove the session token and corresponding session data from the
	// store. If the token does not exist, Delete should simply return nil.
	Delete(token string) error
}
