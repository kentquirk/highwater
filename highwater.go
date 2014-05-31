// Highwater is a server that directs content to the kissmetrics metrics site
package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/kentquirk/userapiclient"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type highwater struct {
	host   string
	apikey string
	salt   string
	client *userapiclient.UserApiClient
}

var instance *highwater = nil

func InitMetrics(metricshost, apikey, salt,
	userapihost, serversecret string) {
	var once sync.Once
	once.Do(func() {
		instance = &highwater{
			metricshost,
			apikey,
			salt,
			userapiclient.NewApiClient("go_highwater", userapihost, serversecret),
		}
	})
}

func saveMetrics(user, event string, parms url.Values) {
	metricsUrl, err := url.Parse(instance.host)
	if err != nil {
		log.Fatal(err)
	}
	vals := make(url.Values)
	for k, v := range parms {
		vals[k] = v
	}
	vals.Add("_p", user)
	vals.Add("_k", instance.apikey)
	vals.Add("_n", event)
	metricsUrl.RawQuery = vals.Encode()

	mu := metricsUrl.String()
	_, err = http.Get(mu)
	log.Println(mu)
	if err != nil {
		log.Println(err) // we'll log the error but won't get upset about it for metrics
	}
}

// hash_id generates the same values as the hash function in the javascript highwater.
func hash_id(id string, length int) string {
	hash := sha1.New()
	io.WriteString(hash, instance.salt)
	io.WriteString(hash, id)
	result := hash.Sum(nil)
	out := fmt.Sprintf("%x", result)
	return out[:length]
}

// NamedUser creates a log entry with a hashed key, given a userid
func NamedUser(userid, event string, parms url.Values) {
	saveMetrics(hash_id(userid, 10), event, parms)
}

// TokenUser creates a log entry with a hashed key by extracting the userid
// from the token passed in
func TokenUser(token, event string, parms url.Values) {
	if td, err := instance.client.CheckToken(token); err == nil {
		saveMetrics(hash_id(td.UserID, 10), event, parms)
	}
}
