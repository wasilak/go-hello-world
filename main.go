package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var (
	listenAddr  string
	sessionKey  string
	logFile     string
	otelEnabled bool
	store       *sessions.CookieStore
	file        *os.File
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -session-key XXXX\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.StringVar(&sessionKey, "session-key", os.Getenv("SESSION_KEY"), "base64 encoded session key or SESSION_KEY env var")
	flag.StringVar(&logFile, "log-file", os.Getenv("LOG_FILE"), "path to log")
	flag.BoolVar(&otelEnabled, "otel-enabled", false, "OpenTelemetry traces enabled")
	flag.Parse()

	if otelEnabled {
		initTracer()
	}

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

	gob.Register(APIStats{})

	store = sessions.NewCookieStore([]byte(sessionKey))

	router := mux.NewRouter()

	router.HandleFunc("/", Chain(rootHandler, Logging()))
	router.HandleFunc("/health", Chain(healthHandler, Logging()))

	router.Methods("GET").Path("/debug/statsviz/ws").Name("GET /debug/statsviz/ws").HandlerFunc(statsviz.Ws)
	router.Methods("GET").PathPrefix("/debug/statsviz/").Name("GET /debug/statsviz/").HandlerFunc(statsviz.Index)

	if otelEnabled {
		router.Use(otelmux.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	}

	http.ListenAndServe(listenAddr, router)
}
