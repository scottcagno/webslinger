package sessions

import (
	"context"
	"runtime"
	"sync"
	"time"
)

var (
	smOnce                sync.Once
	defaultSessionManager SessionManager
	DefaultSessionManager = &defaultSessionManager
)

func initSessionManager(timeout, lifetime time.Duration) *SessionManager {
	smOnce.Do(
		func() {
			defaultSessionManager = SessionManager{
				IdleTimeout: timeout,
				Lifetime:    lifetime,
				Cookie:      SessionCookie{},
				ErrorFunc:   nil,
				store:       openSessionStore(timeout),
				contextKey:  generateContextKey(),
			}
		},
	)
	DefaultSessionManager = &defaultSessionManager
	return DefaultSessionManager
}

var (
	ssOnce              sync.Once
	defaultSessionStore sessionStore
	DefaultSessionStore = &defaultSessionStore
)

func initSessionStore(timeout time.Duration) *sessionStore {
	ssOnce.Do(
		func() {
			if timeout < minimumTimeout {
				timeout = minimumTimeout
			}
			ctx, cancel := context.WithCancel(context.Background())
			defaultSessionStore = sessionStore{
				timeout:  timeout,
				sessions: new(sync.Map),
				ticker:   time.NewTicker(tickerDefault),
				ctx:      ctx,
				cancel:   cancel,
			}
			runtime.SetFinalizer(defaultSessionStore, (*sessionStore).close)
		},
	)
	DefaultSessionStore = &defaultSessionStore
	return DefaultSessionStore
}
