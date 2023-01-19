package sessions

import (
	"net/http"
	"sync"
)

// SessionManager has a default implementation but is made available
// here in case you wish to create your own implementation.
type SessionManager interface {

	// NewSession creates and returns a new *Session.
	NewSession() *Session

	// GetSession takes a http.ResponseWriter and a
	// *http.Request and should attempt to return an
	// existing (active) *Session.
	GetSession(w http.ResponseWriter, r *http.Request) (*Session, error)

	// SaveSession takes a http.ResponseWriter, as well
	// as a *http.Request and a *Session and persists it
	// to the underlying sessionStore as well as writing anything
	// back to the ResponseWriter that the caller may need
	// to have access to.
	SaveSession(w http.ResponseWriter, r *http.Request, sess *Session)

	// KillSession takes a http.ResponseWriter, as well
	// as a *http.Request and a SessionID and removes it
	// from the underlying sessionStore and terminates it client
	// side as well.
	KillSession(w http.ResponseWriter, r *http.Request, sid SessionID)
}

var smOnce sync.Once

var DefaultSessionManager = &defaultSessionManager

var defaultSessionManager sessionManager

type sessionManager struct {
}

func openSessionManager() *sessionManager {
	smOnce.Do(
		func() {
			defaultSessionManager = *initSessionManager()
		},
	)
	return &defaultSessionManager
}

func initSessionManager() *sessionManager {
	sm := &sessionManager{}
	DefaultSessionManager = sm
	return sm
}
