package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/scottcagno/webslinger/pkg/web/sessions"
)

const (
	// Configuration items
	sessionID = "MY_SESSIONS"
	domain
)

var sm = sessions.OpenSessionManager("MY_SESSIONS", "*", 15*time.Minute)
var id int

func main() {

	http.HandleFunc("/session/create", createSession)
	http.HandleFunc("/session/view", viewSession)
	http.HandleFunc("/session/kill", killSession)
	log.Panic(http.ListenAndServe(":3333", nil))

}

func createSession(w http.ResponseWriter, r *http.Request) {
	s := sm.NewSession()
	s.Set("id", id)
	sm.SaveSession(w, r, s)
	fmt.Fprintf(w, "created new session")
	return
}

func viewSession(w http.ResponseWriter, r *http.Request) {
	s, err := sm.GetSession(w, r)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "my session: %s\n", s.String())
	sm.SaveSession(w, r, s)
	return
}

func killSession(w http.ResponseWriter, r *http.Request) {
	s, err := sm.GetSession(w, r)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "killed session")
	sm.KillSession(w, r, s)
	return
}
