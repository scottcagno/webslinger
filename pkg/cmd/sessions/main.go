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

	http.HandleFunc("/secret", secret)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	log.Panic(http.ListenAndServe(":3333", nil))

}

func secret(w http.ResponseWriter, r *http.Request) {
	// check for a session
	s, err := sm.GetSession(w, r)
	if err != nil {
		fmt.Fprintf(w, "error getting session: %s\n", err)
		return
	}
	// check if the user is authenticated
	_, found := s.Get("auth")
	if !found {
		fmt.Fprintf(w, "error getting auth token: %s\n", err)
		return
	}
	// print secret message
	fmt.Fprintf(w, "The cake was a lie!\n")
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	// get or create a session
	s, err := sm.MustGetSession(w, r)
	if err != nil {
		fmt.Fprintf(w, "error getting or creating session: %s\n", err)
		return
	}
	// set an auth token
	s.Set("auth", true)
	// save the session
	err = sm.SaveSession(w, r, s)
	if err != nil {
		fmt.Fprintf(w, "error saving session: %s\n", err)
		return
	}
	fmt.Fprintf(w, "successful login\n")
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	// get the current session
	s, err := sm.GetSession(w, r)
	if err != nil {
		fmt.Fprintf(w, "error getting session: %s\n", err)
		return
	}
	// remove the session
	err = sm.KillSession(w, r, s)
	if err != nil {
		fmt.Fprintf(w, "error removing session: %s\n", err)
		return
	}
	fmt.Fprintf(w, "successfully logged out\n")
	return
}
