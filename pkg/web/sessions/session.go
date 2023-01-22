package sessions

import (
	"encoding/base64"
	"encoding/gob"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// state represents the state of the session data during a request
type sessionState uint8

const (
	unmodified sessionState = iota
	modified
	destroyed
)

// sessionData represents a server side session.
type sessionData struct {
	token   string
	expires time.Time
	state   sessionState

	lock sync.Mutex
	data map[string]any
}

// newSessionData creates and returns a new *Session
func newSessionData(expires time.Duration) *sessionData {
	return &sessionData{
		expires: time.Now().Add(expires).UTC(),
		state:   unmodified,
		data:    make(map[string]any),
	}
}

// get takes a key of type string and returns a value of
// any type along with a boolean indicating true if the
// value was located and false if it was not found.
func (s *sessionData) get(k string) (any, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok := s.data[k]
	return v, ok
}

// put takes a key of type string and value of any type
// and adds or updates it in the Session object.
func (s *sessionData) put(k string, v any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[k] = v
	s.state = modified
}

// del takes a key of type string and attempts to locate
// and remove the key and associated value from the Session.
func (s *sessionData) del(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[k]
	if !ok {
		return
	}
	delete(s.data, k)
	s.state = modified
}

func (s *sessionData) clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.data) == 0 {
		return
	}
	for k := range s.data {
		delete(s.data, k)
	}
	s.state = modified
}

// String implements the Stringer interface for a Session
func (s *sessionData) String() string {
	var sb strings.Builder
	sb.Grow(64)
	sb.WriteString("token=")
	sb.WriteString(s.token)
	sb.WriteString(", expires=")
	sb.WriteString(s.expires.String())
	return sb.String()
}

// Decode should take a raw byte slice and return an expiry time,
// session data, and any potential errors.
func (s *sessionData) Decode(b []byte) (expiry time.Time, data map[string]any, err error) {
	gob.NewEncoder()
}

// Encode should take an expiry time, session data, and return an
// encoded byte slice, along with any potential errors.
func (s *sessionData) Encode(expiry time.Time, data map[string]any) ([]byte, error)

// ExpiresIn returns a duration of time until this
// Session is to be marked as expired.
func (s *sessionData) ExpiresIn() time.Duration {
	return time.Until(s.expires)
}

// IsExpired returns a boolean resulting in true if
// the Session time to live is expired.
func (s *sessionData) IsExpired() bool {
	return time.Until(s.expires) < 1
}

func (s *sessionData) Expires(at time.Time) {
	s.expires = at
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

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
