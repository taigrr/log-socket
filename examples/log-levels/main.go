// Example: log level filtering and all available levels.
//
// log-socket supports 8 log levels from TRACE (most verbose)
// to FATAL (least verbose). Setting a log level filters out
// everything below it.
package main

import (
	"fmt"

	logger "github.com/taigrr/log-socket/v2/log"
)

func main() {
	defer logger.Flush()

	fmt.Println("=== All log levels (TRACE and above) ===")
	logger.SetLogLevel(logger.LTrace)

	logger.Trace("Detailed execution trace — variable x = 42")
	logger.Debug("Processing request for user_id=123")
	logger.Info("Server started on :8080")
	logger.Notice("Configuration reloaded")
	logger.Warn("Disk usage at 85%")
	logger.Error("Failed to send email: SMTP timeout")
	// logger.Panic("...") — would panic
	// logger.Fatal("...") — would os.Exit(1)

	fmt.Println("\n=== Formatted variants ===")
	logger.Infof("Request took %dms", 42)
	logger.Warnf("Retrying in %d seconds (attempt %d/%d)", 5, 2, 3)
	logger.Errorf("HTTP %d: %s", 503, "Service Unavailable")

	fmt.Println("\n=== Only WARN and above ===")
	logger.SetLogLevel(logger.LWarn)

	logger.Debug("This will NOT appear")
	logger.Info("This will NOT appear either")
	logger.Warn("This WILL appear")
	logger.Error("This WILL appear too")

	fmt.Println("\n=== Per-logger levels via namespaced loggers ===")
	logger.SetLogLevel(logger.LTrace) // Reset global

	appLog := logger.NewLogger("app")
	appLog.Info("Namespaced loggers inherit the global output but tag entries")
	appLog.Warnf("Something needs attention in the %s subsystem", "app")
}
