package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/taigrr/log-socket/v2/browser"
	logger "github.com/taigrr/log-socket/v2/log"
	"github.com/taigrr/log-socket/v2/ws"
)

const readHeaderTimeout = 5 * time.Second

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func generateLogs() {
	// Create loggers for different namespaces
	apiLogger := logger.NewLogger("api")
	dbLogger := logger.NewLogger("database")
	authLogger := logger.NewLogger("auth")

	for {
		logger.Info("This is a default namespace log!")
		apiLogger.Info("API request received")
		apiLogger.Debug("Processing API call")

		dbLogger.Info("Database query executed")
		dbLogger.Warn("Slow query detected")

		authLogger.Info("User authentication successful")
		authLogger.Error("Failed login attempt detected")

		logger.Trace("This is a trace log in default namespace!")
		logger.Warn("This is a warning in default namespace!")

		time.Sleep(2 * time.Second)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.LogSocketHandler)
	mux.HandleFunc("/api/namespaces", ws.NamespacesHandler)
	mux.HandleFunc("/", browser.LogSocketViewHandler)
	return mux
}

func newServer(address string) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           newMux(),
		ReadHeaderTimeout: readHeaderTimeout,
	}
}

func main() {
	defer logger.Flush()
	flag.Parse()
	go generateLogs()
	logger.Fatal(newServer(*addr).ListenAndServe())
}
