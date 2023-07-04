package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/arl/statsviz"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"golang.org/x/exp/slog"
)

var (
	listenAddr  string
	sessionKey  string
	logLevel    string
	logFormat   string
	otelEnabled bool
	store       *sessions.CookieStore
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -session-key XXXX\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&listenAddr, "listen-addr", ":5000", "server listen address")
	flag.StringVar(&sessionKey, "session-key", os.Getenv("SESSION_KEY"), "base64 encoded session key or SESSION_KEY env var")
	flag.StringVar(&logLevel, "log-level", os.Getenv("LOG_LEVEL"), "info")
	flag.StringVar(&logFormat, "log-format", os.Getenv("LOG_FORMAT"), "text")
	flag.BoolVar(&otelEnabled, "otel-enabled", false, "OpenTelemetry traces enabled")
	flag.Parse()

	LoggerInit(logLevel, logFormat)

	if sessionKey == "" {
		randomizedSessionKey, err := GenerateKey()
		if err != nil {
			panic(err)
		}
		slog.Info("Session key not provided, generating random one", "sessionKey", randomizedSessionKey)
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
