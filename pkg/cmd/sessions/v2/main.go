package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/scottcagno/webslinger/pkg/web/sessions"
)

var session MySession

type MySession struct {
	*sessions.SessionManager
}

func (s *MySession) LoadAndSaveHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headerKey := "X-Session"
			headerKeyExpiry := "X-Session-Expiry"

			ctx, err := s.Load(r.Context(), r.Header.Get(headerKey))
			if err != nil {
				log.Output(2, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			bw := &bufferedResponseWriter{ResponseWriter: w}
			sr := r.WithContext(ctx)
			next.ServeHTTP(bw, sr)

			if s.Status(ctx) == scs.Modified {
				token, expiry, err := s.Save(ctx)
				if err != nil {
					log.Output(2, err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				w.Header().Set(headerKey, token)
				w.Header().Set(headerKeyExpiry, expiry.Format(http.TimeFormat))
			}

			if bw.code != 0 {
				w.WriteHeader(bw.code)
			}
			w.Write(bw.buf.Bytes())
		},
	)
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf  bytes.Buffer
	code int
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	bw.code = code
}

func main() {
	session = MySession{scs.NewSession()}

	mux := http.NewServeMux()
	mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/get", getHandler)

	http.ListenAndServe(":4000", session.LoadAndSaveHeader(mux))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	session.Put(r.Context(), "message", "Hello from a session!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	msg := session.GetString(r.Context(), "message")
	io.WriteString(w, msg)
}
