package main

import (
	"net/url"
	"testing"
)

func TestHighwater(t *testing.T) {
	InitMetrics("http://trk.kissmetrics.com/e",
		"28814957ee8160309f522a0bd0f2824de585befb",
		"gf78fSEI7tOQQP9xfXMO9HfRyMnW4Sx88Q",
		"http://localhost:9107",
		"This needs to be the same secret everywhere. YaHut75NsK1f9UKUXuWqxNN0RUwHFBCy",
	)
	v := url.Values{}
	v.Add("file", "highwater_test.go")
	ch := make(chan int)
	go func() {
		NamedUser("b@b.com", "Go Highwater Test", v)
		ch <- 1
	}()
	<-ch
}
