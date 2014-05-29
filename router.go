package main

import (
	"fmt"
	"github.com/gorilla/pat"
	"log"
	"net/http"
	"net/url"
)

func StatusHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Ok"))
}

func CopyQueryExcept(query url.Values, except []string) url.Values {
	q, _ := url.ParseQuery(query.Encode())
	for _, k := range except {
		q.Del(k)
	}
	return q
}

func NamedUserHandler(w http.ResponseWriter, req *http.Request) {
	userid := req.URL.Query().Get(":userid")
	eventname := req.URL.Query().Get(":eventname")
	parms := CopyQueryExcept(req.URL.Query(), []string{":userid", ":eventname"})
	NamedUser(userid, eventname, parms)
}

func InitRouter(port int) {
	r := pat.New()
	r.Get("/user/{userid}/{eventname}", NamedUserHandler)
	r.Get("/status", StatusHandler)

	http.Handle("/", r)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
