package sessions

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// ErrNoSession is returned by the Session Manager's GetSession method
// when a Session is not found.
var ErrNoSession = errors.New("sessions: named session not present")

// SessionManager holds the configuration settings for your session and is
// the single structure providing management over your sessions and state.
type SessionManager struct {

	// IdleTimeout controls the maximum length of time a session can
	// be inactive before it expires. For example, some applications
	// may wish to set this so there is a timeout after 20 minutes of
	// inactivity. By default, IdleTimeout is not set and there is no
	// inactivity timeout.
	IdleTimeout time.Duration

	// Lifetime controls the maximum length of time that a session is valid
	// for before it expires. It is an absolute expiry which is set when the
	// session manager is first created and does not change.
	Lifetime time.Duration

	// Cookie contains the configuration settings for session cookies.
	Cookie SessionCookie

	// ErrorFunc allows you to control behavior when an error is encountered
	// by the LoadAndSave middleware. The default behavior is to respond with
	// a 500 http.StatusInternalServerError code. If a custom ErrorFunc is set,
	// then control will be passed to this instead. A typical use would be to
	// provide a function which logs the error and returns a customized HTML
	// error page, or redirects to a certain path.
	ErrorFunc func(http.ResponseWriter, *http.Request, error)

	Codec Codec

	// Store controls the session store, where the session data is persisted.
	store SessionStore

	// contextKey is the key used to set and retrieve the session data from a
	// context.Context. It's automatically generated to ensure uniqueness.
	contextKey contextKey
}

// OpenSessionManager instantiates and returns a new SessionManager
func OpenSessionManager(timeout, lifetime time.Duration) *SessionManager {
	return initSessionManager(timeout, lifetime)
}

// LoadAndSave provides middleware which automatically loads and saves session
// data for the current request, and communicates the session token to and from
// the client in a cookie.
func (sm *SessionManager) LoadAndSave(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			var token string

			// Look for a cookie that we can use to get the
			// current token string from
			cook, err := r.Cookie(sm.Cookie.Name)
			if err == nil {
				token = cook.Value
			}

			// Get an up-to-date version of the context.Context
			// that is associated with this session from the
			// session store.
			ctx, err := sm.Load(r.Context(), token)
			if err != nil {
				sm.ErrorFunc(w, r, err)
				return
			}

			// Update the current request.Context with our up-to-date
			// version of the context (containing this session data),
			// create a buffered response writer, and serve up the
			// handler.
			sr := r.WithContext(ctx)
			bw := &bufferedResponseWriter{ResponseWriter: w}
			next.ServeHTTP(bw, sr)

			// Clear out any form data, because we already served up
			// the handler.
			if sr.MultipartForm != nil {
				err = sr.MultipartForm.RemoveAll()
				if err != nil {
					sm.ErrorFunc(w, r, err)
					return
				}
			}

			// Get the current raw session data from our context, so
			// we can update the cookie for the client.
			sd := sm.getSessionDataFromContext(ctx)
			switch sd.state {
			case modified:
				token, expiry, err := sm.Commit(ctx)
				if err != nil {
					sm.ErrorFunc(w, r, err)
					return
				}
				sm.WriteSessionCookie(ctx, w, token, expiry)
			case destroyed:
				sm.WriteSessionCookie(ctx, w, "", time.Time{})
			}

			// Set our header, and write the contents of our buffered
			// writer to the actual ResponseWriter.
			w.Header().Add("Vary", "Cookie")
			if bw.code != 0 {
				w.WriteHeader(bw.code)
			}
			w.Write(bw.buf.Bytes())
		},
	)
}

// Load retrieves the session data for the given token from the session store,
// and returns a new context.Context containing the session data. If no matching
// token is found then this will create a new session.
func (sm *SessionManager) Load(ctx context.Context, token string) (context.Context, error) {
	if _, ok := ctx.Value(sm.contextKey).(*sessionData); ok {
		return ctx, nil
	}
	if token == "" {
		return sm.addSessionDataToContext(ctx, newSessionData(sm.Lifetime)), nil
	}
	b, found, err := sm.store.Find(token)
	if err != nil {
		return nil, err
	} else if !found {
		return sm.addSessionDataToContext(ctx, newSessionData(sm.Lifetime)), nil
	}

	sd := &sessionData{
		token: token,
		state: unmodified,
	}
	sd.deadline, sd.data, err = sm.Codec.Decode(b)
	if err != nil {
		return nil, err
	}

	// Mark the session data as modified if an idle timeout is being used. This
	// will force the session data to be re-committed to the session store with
	// a new expiry time.
	if sm.IdleTimeout > 0 {
		sd.state = modified
	}

	return sm.addSessionDataToContext(ctx, sd), nil
}

// Commit saves the session data to the session store and returns the session
// token, and expiry time.
func (sm *SessionManager) Commit(ctx context.Context) (string, time.Time, error) {
	sd := sm.getSessionDataFromContext(ctx)

	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.token == "" {
		var err error
		if sd.token, err = generateToken(); err != nil {
			return "", time.Time{}, err
		}
	}

	b, err := sm.Codec.Encode(sd.deadline, sd.values)
	if err != nil {
		return "", time.Time{}, err
	}

	expiry := sd.deadline
	if sm.IdleTimeout > 0 {
		ie := time.Now().Add(sm.IdleTimeout).UTC()
		if ie.Before(expiry) {
			expiry = ie
		}
	}

	if err = sm.doStoreCommit(ctx, sd.token, b, expiry); err != nil {
		return "", time.Time{}, err
	}

	return sd.token, expiry, nil
}

// WriteSessionCookie writes a cookie to the HTTP response with the provided
// token as the cookie value and expiry as the cookie expiry time. The expiry
// time will be included in the cookie only if the session is set to persist
// or has had RememberMe(true) called on it. If expiry is an empty time.Time
// struct (so that it's IsZero() method returns true) the cookie will be
// marked with a historical expiry time and negative max-age (so the browser
// deletes it).
func (sm *SessionManager) WriteSessionCookie(ctx context.Context, w http.ResponseWriter, tok string, exp time.Time) {
	cookie := &http.Cookie{
		Name:     sm.Cookie.Name,
		Value:    tok,
		Path:     sm.Cookie.Path,
		Domain:   sm.Cookie.Domain,
		Secure:   sm.Cookie.Secure,
		HttpOnly: sm.Cookie.HttpOnly,
		SameSite: sm.Cookie.SameSite,
	}
	switch {
	case exp.IsZero():
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	case sm.Cookie.Persist:
		cookie.Expires = time.Unix(exp.Unix()+1, 0)
		cookie.MaxAge = int(time.Until(exp).Seconds() + 1)
	}
	w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)
	http.SetCookie(w, cookie)
}

func (sm *SessionManager) addSessionDataToContext(ctx context.Context, sd *sessionData) context.Context {
	return context.WithValue(ctx, sm.contextKey, sd)
}

func (sm *SessionManager) getSessionDataFromContext(ctx context.Context) *sessionData {
	c, ok := ctx.Value(sm.contextKey).(*sessionData)
	if !ok {
		panic("scs: no session data in context")
	}
	return c
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf         bytes.Buffer
	code        int
	wroteHeader bool
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	if !bw.wroteHeader {
		bw.code = code
		bw.wroteHeader = true
	}
}

func (bw *bufferedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := bw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (bw *bufferedResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := bw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

/*

// NewSession creates and returns a new *Session
func (sm *SessionManager) NewSession() Session {
	return sm.store.newSession()
}

// MustGetSession checks for an existing session in the store using a cookie
// with the same name that the session manager was provided with. If one is
// not found, then it creates a new one and returns it.
func (sm *SessionManager) MustGetSession(w http.ResponseWriter, r *http.Request) (Session, error) {
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
func (sm *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) (Session, error) {
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
func (sm *SessionManager) SaveSession(w http.ResponseWriter, r *http.Request, sess sessionData) error {
	if sess == nil {
		return ErrNoSession
	}
	// persist the session to the store
	sm.store.saveSession(sess)
	// update the session cookie
	http.SetCookie(w, NewCookie(sm.name, sess.token, sm.domain, time.Unix(int64(sess.ExpiresIn()), 0)))
	return nil
}

// KillSession removes an existing session using the SessionID.
func (sm *SessionManager) KillSession(w http.ResponseWriter, r *http.Request, sess Session) error {
	if sess == nil {
		return ErrNoSession
	}
	// Check for an existing session by looking in the request for a cookie.
	// If we find a cookie we must expire it.
	c, err := r.Cookie(sm.name)
	if c == nil || err == http.ErrNoCookie {
		return nil
	}
	// Remove the session from the store, and put the updated cookie
	sm.store.killSession(sess)
	c.Expires = time.Now()
	c.MaxAge = -1
	http.SetCookie(w, c)
	return nil
}

*/
