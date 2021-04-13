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

// HealthResponse type
type HealthResponse struct {
	Status string `json:"status"`
}

// APIResponseRequest type
type APIResponseRequest struct {
	Host       string   `json:"host"`
	RemoteAddr string   `json:"remote_addr"`
	RequestURI string   `json:"request_uri"`
	Method     string   `json:"method"`
	Proto      string   `json:"proto"`
	UserAgent  string   `json:"user_agent"`
	URL        *url.URL `json:"url"`
	Headers http.Header `json:"headers"`
}

// APIResponse type
type APIResponse struct {
	Counter   int                `json:"counter"`
	Host      string             `json:"host"`
	Hostnames map[string]int     `json:"hostnames"`
	Request   APIResponseRequest `json:"request"`
}

var (
	listenAddr string
	sessionKey string
	logFile    string
	store      *sessions.FilesystemStore
	file       *os.File
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

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

	response.Request = APIResponseRequest{
		Host:       r.Host,
		URL:        r.URL,
		RemoteAddr: r.RemoteAddr,
		RequestURI: r.RequestURI,
		Method:     r.Method,
		Proto:      r.Proto,
		UserAgent:  r.UserAgent(),
		Headers: r.Header,
	}

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
			defer func() {
				logLine := fmt.Sprintf("%s %s\n", r.URL.Path, time.Since(start))
				log.Print(logLine)
				if logFile != "" {
					if _, err := file.WriteString(logLine); err != nil {
						log.Println(err)
					}
				}
			}()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -session-key XXXX\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.StringVar(&sessionKey, "session-key", os.Getenv("SESSION_KEY"), "base64 encoded session key or SESSION_KEY env var")
	flag.StringVar(&logFile, "log-file", os.Getenv("LOG_FILE"), "path to log")
	flag.Parse()

	var err error
	if logFile != "" {
		file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
	}

	if sessionKey == "" {
		fmt.Println("Please provide session key using: -session-key or SESSION_KEY env var")
		flag.Usage()
		os.Exit(1)
	}

	gob.Register(APIResponse{})

	store = sessions.NewFilesystemStore("", []byte(sessionKey))

	router := mux.NewRouter()

	router.HandleFunc("/", Chain(rootHandler, Logging()))
	router.HandleFunc("/health", Chain(healthHandler, Logging()))

	http.ListenAndServe(listenAddr, router)
}
