package sessions

import (
	"context"
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

// session represents a server side session.
type session struct {
	token   string
	expires time.Time
	state   sessionState

	lock sync.Mutex
	data map[string]any
}

// newSessionData creates and returns a new *Session
func newSessionData(expires time.Duration) *session {
	return &session{
		expires: time.Now().Add(expires).UTC(),
		state:   unmodified,
		data:    make(map[string]any),
	}
}

// get takes a key of type string and returns a value of
// any type along with a boolean indicating true if the
// value was located and false if it was not found.
func (s *session) get(k string) (any, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok := s.data[k]
	return v, ok
}

// put takes a key of type string and value of any type
// and adds or updates it in the Session object.
func (s *session) put(k string, v any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[k] = v
	s.state = modified
}

// del takes a key of type string and attempts to locate
// and remove the key and associated value from the Session.
func (s *session) del(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.data[k]
	if !ok {
		return
	}
	delete(s.data, k)
	s.state = modified
}

// clear removes all data for the current session. The session token
// and lifetime are unaffected.
func (s *session) clear() {
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
func (s *session) String() string {
	var sb strings.Builder
	sb.Grow(64)
	sb.WriteString("token=")
	sb.WriteString(s.token)
	sb.WriteString(", expires=")
	sb.WriteString(s.expires.String())
	return sb.String()
}

func (sm *SessionManager) addSessionData(ctx context.Context, sess *session) context.Context {
	return context.WithValue(ctx, sm.ctxKey, sess)
}

func (sm *SessionManager) getSessionData(ctx context.Context) *session {
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	return sess
}

func (sm *SessionManager) getSessionState(ctx context.Context) sessionState {
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	sess.lock.Lock()
	defer sess.lock.Lock()
	var state sessionState
	state = sess.state
	return state
}
