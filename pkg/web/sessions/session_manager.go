package sessions

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

// ErrNoSession is returned by the Session Manager's GetSession method
// when a Session is not found.
var ErrNoSession = errors.New("sessions: named session not present")

// SessionManager has a default implementation but is made available
// here in case you wish to create your own implementation.
type SessionManager interface {

	// NewSession creates and returns a new *Session.
	NewSession() *Session

	// GetSession takes a http.ResponseWriter and a
	// *http.Request and should attempt to return an
	// existing (active) *Session or an error if the
	// named session can not be found.
	GetSession(w http.ResponseWriter, r *http.Request) (*Session, error)

	MustGetSession(w http.ResponseWriter, r *http.Request) (*Session, error)

	// SaveSession takes a http.ResponseWriter, as well
	// as a *http.Request and a *Session and persists it
	// to the underlying sessionStore as well as writing anything
	// back to the ResponseWriter that the caller may need
	// to have access to.
	SaveSession(w http.ResponseWriter, r *http.Request, sess *Session) error

	// KillSession takes a http.ResponseWriter, as well
	// as a *http.Request and a SessionID and removes it
	// from the underlying sessionStore and terminates it client
	// side as well.
	KillSession(w http.ResponseWriter, r *http.Request, sess *Session) error
}

var smOnce sync.Once

var DefaultSessionManager = &defaultSessionManager

var defaultSessionManager sessionManager

type sessionManager struct {
	name     string
	duration time.Duration
	domain   string
	store    *sessionStore
}

// OpenSessionManager instantiates and returns a new SessionManager
func OpenSessionManager(name, domain string, duration time.Duration) SessionManager {
	smOnce.Do(
		func() {
			defaultSessionManager = *initSessionManager(name, domain, duration)
		},
	)
	return &defaultSessionManager
}

func initSessionManager(name, domain string, duration time.Duration) *sessionManager {
	sm := &sessionManager{
		name:     name,
		domain:   domain,
		duration: duration,
		store:    openSessionStore(duration),
	}
	DefaultSessionManager = sm
	return sm
}

// NewSession creates and returns a new *Session
func (sm *sessionManager) NewSession() *Session {
	return sm.store.newSession()
}

// MustGetSession checks for an existing session in the store using a cookie
// with the same name that the session manager was provided with. If one is
// not found, then it creates a new one and returns it.
func (sm *sessionManager) MustGetSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	// Check for an existing session by looking in the request for a cookie.
	c, err := r.Cookie(sm.name)
	if err == http.ErrNoCookie {
		// No cookie was found, we will return a new session
		return sm.store.newSession(), nil
	}
	// Otherwise, we have found a session cookie, but we must check to ensure
	// that it is not expired.
	sess, found := sm.store.getSession(SessionID(c.Value))
	if !found {
		// No session has been found, we will return an error
		return nil, ErrNoSession
	}
	// Otherwise, we have successfully located an existing session that we can
	// return along with a nil error
	return sess, nil
}

// GetSession checks for an existing session in the store using a cookie
// with the same name that the session manager was provided with.
func (sm *sessionManager) GetSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	// Check for an existing session by looking in the request for a cookie.
	c, err := r.Cookie(sm.name)
	if err == http.ErrNoCookie {
		// No cookie was found, we will return an error
		return nil, err
	}
	// Otherwise, we have found a session cookie, but we must check to ensure
	// that it is not expired.
	sess, found := sm.store.getSession(SessionID(c.Value))
	if !found {
		// No session has been found, we will return an error
		return nil, ErrNoSession
	}
	// Otherwise, we have successfully located an existing session that we can
	// return along with a nil error
	return sess, nil
}

// SaveSession persists the provided session. If it receives a nil *Session, it will
// return an error.
func (sm *sessionManager) SaveSession(w http.ResponseWriter, r *http.Request, sess *Session) error {
	if sess == nil {
		return ErrNoSession
	}
	// persist the session to the store
	sm.store.saveSession(sess)
	// update the session cookie
	http.SetCookie(w, NewCookie(sm.name, string(sess.ID), sm.domain, sess.Expires))
	return nil
}

// KillSession removes an existing session using the SessionID.
func (sm *sessionManager) KillSession(w http.ResponseWriter, r *http.Request, sess *Session) error {
	if sess == nil {
		return ErrNoSession
	}
	// Check for an existing session by looking in the request for a cookie.
	// If we find a cookie we must expire it.
	c, err := r.Cookie(sm.name)
	if c == nil || err == http.ErrNoCookie {
		return nil
	}
	// Remove the session from the store, and set the updated cookie
	sm.store.killSession(sess)
	c.Expires = time.Now()
	c.MaxAge = -1
	http.SetCookie(w, c)
	return nil
}
