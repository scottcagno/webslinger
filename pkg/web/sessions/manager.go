package sessions

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
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
	Cookie CookieConfig

	// ErrorFunc allows you to control behavior when an error is encountered
	// by the LoadAndSave middleware. The default behavior is to respond with
	// a 500 http.StatusInternalServerError code. If a custom ErrorFunc is set,
	// then control will be passed to this instead. A typical use would be to
	// provide a function which logs the error and returns a customized HTML
	// error page, or redirects to a certain path.
	ErrorFunc func(http.ResponseWriter, *http.Request, error)

	// Codec is the Codec that is used to encode and decode data to and from
	// the underlying Store and the local session manager session type.
	Codec Codec

	// Store controls the session Store, where the session data is persisted.
	Store SessionStore

	// ctxKey is the key used to set and retrieve the session data from a
	// context.Context. It's automatically generated to ensure uniqueness.
	ctxKey ctxKey
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		IdleTimeout: 0,
		Lifetime:    24 * time.Hour,
		Cookie: CookieConfig{
			Name:     "session",
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Persist:  true,
		},
		ErrorFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Output(2, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		},
		Codec:  GobCodec{},
		Store:  NewMemoryStoreWithInterval(15 * time.Minute),
		ctxKey: generateContextKey(),
	}
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
			// session Store.
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

			// Get the current session data state, so we know how what
			// to do with our cookies.
			switch sm.getSessionState(ctx) {
			case modified:
				token, expiry, err := sm.Save(ctx)
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
			_, err = w.Write(bw.buf.Bytes())
			if err != nil {
				sm.ErrorFunc(w, r, err)
				return
			}
		},
	)
}

// Load retrieves the session data for the given token from the session Store,
// and returns a new context.Context containing the session data. If no matching
// token is found then this will create a new session.
func (sm *SessionManager) Load(ctx context.Context, token string) (context.Context, error) {
	// Check the context for a cached session and return it.
	_, ok := ctx.Value(sm.ctxKey).(*session)
	if ok {
		return ctx, nil
	}
	// If we don't have a cached session, and we don't have a token then we need
	// to create a brand-new session.
	if token == "" {
		// Return a new session instance wrapped inside a context
		return context.WithValue(ctx, sm.ctxKey, newSessionData(sm.Lifetime)), nil
	}
	// Otherwise, we need to check the Store using the provided token.
	b, err := sm.Store.Find(token)
	if err != nil {
		// We go an error from the Store
		if err == ErrSessionNotFound {
			// Session was not found, we have to create a new instance and return it
			// inside a new context
			return context.WithValue(ctx, sm.ctxKey, newSessionData(sm.Lifetime)), nil
		}
		// Otherwise, it's a bad error, and we should exit and return
		return nil, err
	}
	// Decode our raw session data
	expires, data, err := sm.Codec.Decode(b)
	if err != nil {
		return nil, err
	}
	// Initialize a new session data type
	sess := &session{
		token:   token,
		expires: expires,
		state:   unmodified,
		data:    data,
	}
	// Mark the session data as modified if an idle timeout is being used. This
	// will force the session data to be re-committed to the session Store with
	// a new expiry time.
	if sm.IdleTimeout > 0 {
		sess.state = modified
	}
	// Add it to our context, and return
	return context.WithValue(ctx, sm.ctxKey, sess), nil
}

var errNoSessionDataFoundInContext = errors.New("session manager: no session data found in context")

// Save saves the session data to the session Store and returns the session
// token, and expiry time.
func (sm *SessionManager) Save(ctx context.Context) (string, time.Time, error) {
	// Get the session data from the context
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	// Lock it up!
	sess.lock.Lock()
	defer sess.lock.Unlock()
	// Generate a fresh token
	if sess.token == "" {
		token, err := generateToken()
		if err != nil {
			return "", time.Time{}, err
		}
		sess.token = token
	}
	// Encode the session data, so we can save it back to the Store
	b, err := sm.Codec.Encode(sess.expires, sess.data)
	if err != nil {
		return "", time.Time{}, err
	}
	// For security purposes, we should ensure that the session expiry
	// time is not set too far in the future.
	expiry := sess.expires
	if sm.IdleTimeout > 0 {
		ie := time.Now().Add(sm.IdleTimeout).UTC()
		if ie.Before(expiry) {
			// Session duration is too long, lets bring it back within
			// the idle timeout range before we save.
			expiry = ie
		}
	}
	// Save the session data to the underlying Store
	err = sm.Store.Save(sess.token, b, expiry)
	if err != nil {
		return "", time.Time{}, err
	}
	return sess.token, expiry, nil
}

// Destroy deletes the session data from the underlying Store and sets
// the session status to destroyed. Any further action in the same
// request chain will result in the creation of a new session.
func (sm *SessionManager) Destroy(ctx context.Context) error {
	// Get the session data from the context
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	// Lock it up!
	sess.lock.Lock()
	defer sess.lock.Unlock()
	// Call the stores delete method
	err := sm.Store.Delete(sess.token)
	if err != nil {
		return err
	}
	// Update the session details
	sess.token = ""
	sess.expires = time.Now().Add(sm.Lifetime).UTC()
	sess.state = destroyed
	// for k := range sess.data {
	// 	delete(sess.data, k)
	// }
	sess.data = nil
	return nil
}

// Put adds a key and corresponding value to the session data. Any existing
// key and value that matches will be replaced, and the state will be changed
// to the `modified` state.
func (sm *SessionManager) Put(ctx context.Context, key string, val any) {
	// Get the session data from the context
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	// Call the `put` method on session data. All session methods lock and
	// automatically update the state accordingly.
	sess.put(key, val)
}

func (sm *SessionManager) Get(ctx context.Context, key string) any {
	// Get the session data from the context
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	// Call the `get` method on session data. All session methods lock and
	// automatically update the state accordingly.
	val, found := sess.get(key)
	if !found {
		return nil
	}
	return val
}

// Del deletes the given key and associated value from the session data and
// updates the state to `modified` accordingly.
func (sm *SessionManager) Del(ctx context.Context, key string) {
	// Get the session data from the context
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	// Call the `del` method on session data. All session methods lock and
	// automatically update the state accordingly.
	sess.del(key)
}

// Clear deletes all the keys and associated values from the session data
// and updates the state to `modified` accordingly.
func (sm *SessionManager) Clear(ctx context.Context) {
	// Get the session data from the context
	sess, ok := ctx.Value(sm.ctxKey).(*session)
	if !ok {
		panic(errNoSessionDataFoundInContext)
	}
	// Call the `clear` method on session data. All session methods lock and
	// automatically update the state accordingly.
	sess.clear()
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

// generateToken generates a new unique token
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ctxKey is the unique context key that the session manager
// uses to keep track of the sessions.
type ctxKey string

var (
	ctxKeyID     uint64
	ctxKeyIDLock = new(sync.Mutex)
)

func generateContextKey() ctxKey {
	ctxKeyIDLock.Lock()
	defer ctxKeyIDLock.Unlock()
	atomic.AddUint64(&ctxKeyID, 1)
	return ctxKey("session." + strconv.FormatUint(ctxKeyID, 10))
}
