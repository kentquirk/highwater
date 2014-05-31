package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/pat"
)

// Just returns Ok if the server is running
func StatusHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Ok"))
}

// CopyQueryExcept creates a clone of a url.Values object with some of the keys missing
// It's sort of the equivalent of _.omit() in js.
func CopyQueryExcept(query url.Values, except []string) url.Values {
	q, _ := url.ParseQuery(query.Encode())
	for _, k := range except {
		q.Del(k)
	}
	return q
}

// NamedUserHandler provides a handler for the route where the user is explicitly specified.
func NamedUserHandler(w http.ResponseWriter, req *http.Request) {
	userid := req.URL.Query().Get(":userid")
	eventname := req.URL.Query().Get(":eventname")
	parms := CopyQueryExcept(req.URL.Query(), []string{":userid", ":eventname"})
	NamedUser(userid, eventname, parms)
}

// TokenUserHandler provides a handler for the route where the user is implied by the token passed in
func TokenUserHandler(w http.ResponseWriter, req *http.Request) {
	eventname := req.URL.Query().Get(":eventname")
	parms := CopyQueryExcept(req.URL.Query(), []string{":eventname"})
	token := req.Header.Get("x-tidepool-session-token")
	TokenUser(token, eventname, parms)
}

// InitRouter sets up the server; it never returns unless there's an unrecoverable error.
func InitRouter(port int) {
	r := pat.New()
	r.Get("/user/{userid}/{eventname}", NamedUserHandler)
	r.Get("/status", StatusHandler)

	http.Handle("/", r)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
