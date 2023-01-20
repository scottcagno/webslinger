package sessions

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Session represents a server side session.
type Session struct {
	ID      SessionID
	data    *sync.Map
	Expires time.Time
}

// newSession creates and returns a new *Session
func newSession(ttl time.Duration) *Session {
	return &Session{
		ID:      NewSessionID(),
		data:    new(sync.Map),
		Expires: time.Now().Add(ttl),
	}
}

// ExpiresIn returns a duration of time until this
// Session is to be marked as expired.
func (s *Session) ExpiresIn() time.Duration {
	return time.Until(s.Expires)
}

// IsExpired returns a boolean resulting in true if
// the Session time to live is expired.
func (s *Session) IsExpired() bool {
	return time.Until(s.Expires) < 1
}

// Set takes a key of type string and value of any type
// and adds or updates it in the Session object.
func (s *Session) Set(k string, v any) {
	s.data.Store(k, v)
}

// Get takes a key of type string and returns a value of
// any type along with a boolean indicating true if the
// value was located and false if it was not found.
func (s *Session) Get(k string) (any, bool) {
	v, ok := s.data.Load(k)
	return v, ok
}

// Del takes a key of type string and attempts to locate
// and remove the key and associated value from the Session.
func (s *Session) Del(k string) {
	s.data.Delete(k)
}

// String implements the Stringer interface for a Session
func (s *Session) String() string {
	var sb strings.Builder
	sb.Grow(64)
	sb.WriteString("ID=")
	sb.WriteString(string(s.ID))
	sb.WriteString(", TTL=")
	sb.WriteString(s.ExpiresIn().String())
	return sb.String()
}

// ascii is a constant value of all the basic alphanumeric
// characters and is used in the NewSessionID function to
// create and return a new SessionID.
const ascii = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// SessionID represents a primary session ID key.
// type SessionID [32]byte
type SessionID string

// NewSessionID creates and returns a (hopefully) unique ID
// that can be used to ID a Session object.
func NewSessionID() SessionID {
	sid := make([]byte, 32, 32)
	for i := range sid {
		sid[i] = ascii[rand.Intn(len(ascii))]
	}
	return SessionID(sid)
}
