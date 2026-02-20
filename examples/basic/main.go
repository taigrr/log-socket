// Example: basic usage of log-socket as a drop-in logger.
//
// This demonstrates using the package-level logging functions,
// which work similarly to the standard library's log package.
package main

import (
	"fmt"
	"net/http"

	"github.com/taigrr/log-socket/v2/browser"
	logger "github.com/taigrr/log-socket/v2/log"
	"github.com/taigrr/log-socket/v2/ws"
)

func main() {
	defer logger.Flush()

	// Set the minimum log level (default is LTrace, showing everything)
	logger.SetLogLevel(logger.LDebug)

	// Package-level functions log to the "default" namespace
	logger.Info("Application starting up")
	logger.Debug("Debug mode enabled")
	logger.Warnf("Config file not found at %s, using defaults", "/etc/app/config.yaml")
	logger.Errorf("Failed to connect to database: %s", "connection refused")

	// Print/Printf/Println are aliases for Info
	logger.Println("This is equivalent to Infoln")

	// Start the web UI so you can view logs at http://localhost:8080
	http.HandleFunc("/ws", ws.LogSocketHandler)
	http.HandleFunc("/", browser.LogSocketViewHandler)
	fmt.Println("Log viewer available at http://localhost:8080")
	logger.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
