// Example: programmatic log client with namespace filtering.
//
// This shows how to create a Client that receives log entries
// programmatically, optionally filtered to specific namespaces.
// Useful for building custom log processors, alerting, or forwarding.
package main

import (
	"fmt"
	"time"

	logger "github.com/taigrr/log-socket/v2/log"
)

func main() {
	defer logger.Flush()

	// Create a client that receives ALL log entries
	allLogs := logger.CreateClient()
	allLogs.SetLogLevel(logger.LInfo)

	// Create a client that only receives "database" and "auth" logs
	securityLogs := logger.CreateClient("database", "auth")
	securityLogs.SetLogLevel(logger.LWarn) // Only warnings and above

	dbLog := logger.NewLogger("database")
	authLog := logger.NewLogger("auth")
	apiLog := logger.NewLogger("api")

	// Process all logs
	go func() {
		for {
			entry := allLogs.Get()
			fmt.Printf("[ALL] %s [%s] %s: %s\n",
				entry.Timestamp.Format(time.TimeOnly),
				entry.Namespace, entry.Level, entry.Output)
		}
	}()

	// Process only security-relevant warnings/errors
	go func() {
		for {
			entry := securityLogs.Get()
			if entry.Level == "ERROR" || entry.Level == "WARN" {
				fmt.Printf("🚨 SECURITY ALERT [%s] %s: %s\n",
					entry.Namespace, entry.Level, entry.Output)
			}
		}
	}()

	// Generate some logs
	for i := 0; i < 5; i++ {
		apiLog.Info("API request processed")
		dbLog.Info("Query executed successfully")
		dbLog.Warn("Connection pool running low")
		authLog.Error("Brute force attempt detected")
		time.Sleep(1 * time.Second)
	}

	// Clean up clients when done
	allLogs.Destroy()
	securityLogs.Destroy()
}
