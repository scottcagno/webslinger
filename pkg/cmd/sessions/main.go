package main

import (
	"io"
	"net/http"

	"github.com/scottcagno/webslinger/pkg/web/sessions"
)

var sm = sessions.NewSessionManager()

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/put", putHandler)
	mux.HandleFunc("/get", getHandler)

	http.ListenAndServe(":4000", sm.LoadAndSave(mux))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	sm.Put(r.Context(), "message", "Hello from a session!")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	msg := sm.Get(r.Context(), "message").(string)
	io.WriteString(w, msg)
}
