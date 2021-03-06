// ipify-api
//
// This software implements a basic REST API that provides users with a simple
// way to query their public IP address (IPv4 or IPv6).  This code assumes that
// you are running it on Heroku's platform (https://www.heroku.com/).

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
)

// IPAddress is a simple struct that we use to marshal our JSON responses.
type IPAddress struct {
	IP string `json:"ip"`
}

// getIP returns a user's public facing IP address (IPv4 OR IPv6).
//
// By default, it will return the IP address in plain text, but can also return
// data in both JSON and JSONP if requested to.
func getIP(w http.ResponseWriter, r *http.Request) {

	// Enable CORS support.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	ip := net.ParseIP(r.Header["X-Forwarded-For"][len(r.Header["X-Forwarded-For"])-1]).String()

	// If the user specifies a 'format' querystring, we'll try to return the
	// user's IP address in the specified format.
	if format, ok := r.Form["format"]; ok && len(format) > 0 {
		jsonStr, _ := json.Marshal(IPAddress{ip})

		switch format[0] {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, string(jsonStr))
			return
		case "jsonp":
			// If the user specifies a 'callback' parameter, we'll use that as
			// the name of our JSONP callback.
			callback := "callback"
			if val, ok := r.Form["callback"]; ok && len(val) > 0 {
				callback = val[0]
			}

			w.Header().Set("Content-Type", "application/javascript")
			fmt.Fprintf(w, callback+"("+string(jsonStr)+");")
			return
		}
	}

	// If no 'format' querystring was specified, we'll default to returning the
	// IP in plain text.
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, ip)
}

// main launches our web server which runs indefinitely.
func main() {
	http.HandleFunc("/", getIP)

	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
