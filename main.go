package main

import (
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// APIResponse type
type APIResponse struct {
	Counter   int            `json:"counter"`
	Host      string         `json:"host"`
	Hostnames map[string]int `json:"hostnames"`
	Request   struct {
		Host       string   `json:"host"`
		RemoteAddr string   `json:"remote_addr"`
		RequestURI string   `json:"request_uri"`
		Method     string   `json:"method"`
		Proto      string   `json:"proto"`
		UserAgent  string   `json:"user_agent"`
		URL        *url.URL `json:"url"`
	} `json:"request"`
}

var (
	listenAddr string
	sessionKey string
	store      *sessions.CookieStore
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-go-hello-world")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	val := session.Values["apiResponse"]

	response, ok := val.(APIResponse)

	if !ok {
		log.Println("session not initialized (yet)")
	}

	response.Counter++

	hostname, _ := os.Hostname()
	response.Host = hostname

	if nil == response.Hostnames {
		response.Hostnames = make(map[string]int)
	}
	response.Hostnames[hostname]++

	response.Request.Host = r.Host
	response.Request.URL = r.URL
	response.Request.RemoteAddr = r.RemoteAddr
	response.Request.RequestURI = r.RequestURI
	response.Request.Method = r.Method
	response.Request.Proto = r.Proto
	response.Request.UserAgent = r.UserAgent()

	session.Values["apiResponse"] = response

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

func main() {
	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.StringVar(&sessionKey, "session-key", os.Getenv("SESSION_KEY"), "base64 encoded session key or SESSION_KEY env var")
	flag.Parse()

	if sessionKey == "" {
		fmt.Println("Please provide session key using: -session-key or SESSION_KEY env var")
		flag.Usage()
		os.Exit(1)
	}

	gob.Register(APIResponse{})

	store = sessions.NewCookieStore([]byte(sessionKey))

	router := mux.NewRouter()

	router.HandleFunc("/", Chain(rootHandler, Logging()))

	http.ListenAndServe(listenAddr, router)
}
