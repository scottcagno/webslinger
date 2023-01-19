package sessions

import (
	"context"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	minimumTimeout = 1 * time.Minute
	tickerDefault  = 10 * time.Second
)

var ssOnce sync.Once

var DefaultSessionStore = &defaultSessionStore

var defaultSessionStore sessionStore

type sessionStore struct {
	storeKey string
	timeout  time.Duration
	sessions *sync.Map
	ticker   *time.Ticker
	ctx      context.Context
	cancel   context.CancelFunc
}

func openSessionStore(key string, timeout time.Duration) *sessionStore {
	ssOnce.Do(
		func() {
			defaultSessionStore = *initSessionStore(key, timeout)
		},
	)
	return &defaultSessionStore
}

func initSessionStore(key string, timeout time.Duration) *sessionStore {
	if timeout < minimumTimeout {
		timeout = minimumTimeout
	}
	ctx, cancel := context.WithCancel(context.Background())
	ss := &sessionStore{
		storeKey: key,
		timeout:  timeout,
		sessions: new(sync.Map),
		ticker:   time.NewTicker(tickerDefault),
		ctx:      ctx,
		cancel:   cancel,
	}
	DefaultSessionStore = ss
	go ss.cleanUpRoutine()
	runtime.SetFinalizer(ss, (*sessionStore).close)
	return ss
}

func (ss *sessionStore) newSession() *Session {
	return newSession(ss.timeout)
}

// getSession takes a SessionID and attempts to locate the
// matching *Session. If a matching *Session can be found
// it is returned along with a boolean indicating weather or
// not the session was found.
func (ss *sessionStore) getSession(sid SessionID) (*Session, bool) {
	session, found := ss.sessions.Load(sid)
	if !found {
		return nil, false
	}
	return session.(*Session), true
}

// saveSession takes a *Session and persists it to the underlying sessionStore.
func (ss *sessionStore) saveSession(session *Session) {
	// If the session is nil, do nothing the checkForExpiredSessions
	// will handle any of the extra cleanup necessary.
	if session == nil {
		return
	}
	// If the session is expired, remove it and return.
	if session.IsExpired() {
		ss.sessions.Delete(session.ID)
		return
	}
	// Otherwise, bump the expiry time and save it.
	session.ttl = time.Now().Add(ss.timeout)
	ss.sessions.Store(session.ID, session)
}

func (ss *sessionStore) cleanUpRoutine() {
	// When we receive a "tick", we should loop through the
	// sessions, checking to see if any of them are expired.
	// If we find any that are expired, we should remove them.
	for {
		select {
		case t := <-ss.ticker.C:
			// Clean up expired sessions
			log.Printf("Checking for expired sessions: %v\n", t)
			ss.sessions.Range(
				func(sid, session any) bool {
					if session.(*Session).IsExpired() {
						ss.sessions.Delete(sid)
					}
					return true
				},
			)
		case <-ss.ctx.Done():
			ss.ticker.Stop()
			return
		}
	}
}

func (ss *sessionStore) close() {
	log.Printf("*sessionStore.Close has been called.\n")
	// stop the ticker and free any other
	// resources.
	ss.ticker.Stop()
	ss.cancel()
}

func (ss *sessionStore) count() int {
	var sessionCount int
	ss.sessions.Range(
		func(k, v any) bool {
			if k != nil && v != nil {
				sessionCount++
			}
			return true
		},
	)
	return sessionCount
}

// String implements the Stringer interface for a session sessionStore.
func (ss *sessionStore) String() string {
	var sb strings.Builder
	ss.sessions.Range(
		func(k, v any) bool {
			if k != nil && v != nil {
				if session, ok := v.(*Session); ok {
					sb.WriteString(session.String())
					sb.WriteByte('\n')
				}
			}
			return true
		},
	)
	return sb.String()
}
