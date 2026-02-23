// Example: namespace-based logging for organizing logs by component.
//
// Namespaces let you tag log entries by subsystem (api, database, auth, etc.)
// and filter them in the web UI or via programmatic clients.
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/taigrr/log-socket/v2/browser"
	logger "github.com/taigrr/log-socket/v2/log"
	"github.com/taigrr/log-socket/v2/ws"
)

func main() {
	defer logger.Flush()

	// Create loggers for different subsystems
	apiLog := logger.NewLogger("api")
	dbLog := logger.NewLogger("database")
	authLog := logger.NewLogger("auth")
	cacheLog := logger.NewLogger("cache")

	// Simulate application activity
	go func() {
		for {
			apiLog.Infof("GET /api/users — 200 OK (%dms)", rand.Intn(200))
			apiLog.Debugf("Request headers: Accept=application/json")
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			dbLog.Infof("SELECT * FROM users — %d rows", rand.Intn(100))
			if rand.Float64() < 0.3 {
				dbLog.Warn("Slow query detected (>500ms)")
			}
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			if rand.Float64() < 0.7 {
				authLog.Info("User login successful")
			} else {
				authLog.Error("Failed login attempt from 192.168.1.42")
			}
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			cacheLog.Tracef("Cache hit ratio: %.1f%%", rand.Float64()*100)
			if rand.Float64() < 0.1 {
				cacheLog.Warn("Cache eviction triggered")
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// The /api/namespaces endpoint lists all active namespaces
	http.HandleFunc("/ws", ws.LogSocketHandler)
	http.HandleFunc("/api/namespaces", ws.NamespacesHandler)
	http.HandleFunc("/", browser.LogSocketViewHandler)
	fmt.Println("Log viewer with namespace filtering at http://localhost:8080")
	logger.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
