package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	logger "github.com/taigrr/log-socket/v2/log"
)

var upgrader = websocket.Upgrader{} // use default options

// SetUpgrader replaces the default [websocket.Upgrader] used by
// [LogSocketHandler].
func SetUpgrader(u websocket.Upgrader) {
	upgrader = u
}

func parseNamespaces(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	namespaces := make([]string, 0, len(parts))
	for _, part := range parts {
		namespace := strings.TrimSpace(part)
		if namespace == "" {
			continue
		}
		namespaces = append(namespaces, namespace)
	}
	if len(namespaces) == 0 {
		return nil
	}
	return namespaces
}

// LogSocketHandler upgrades the HTTP connection to a WebSocket and streams
// log entries to the client. An optional "namespaces" query parameter
// (comma-separated) filters which namespaces the client receives.
func LogSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Get namespaces from query parameter, comma-separated.
	// Empty or missing means all namespaces.
	namespaces := parseNamespaces(r.URL.Query().Get("namespaces"))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer conn.Close()

	lc := logger.CreateClient(namespaces...)
	defer lc.Destroy()
	lc.SetLogLevel(logger.LTrace)
	logger.Info("Websocket client attached.")

	// Start a read pump so the server detects client disconnects promptly.
	// Without this, a disconnected client is only noticed when WriteMessage
	// fails, which can be delayed indefinitely when no logs are produced.
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go func() {
		defer cancel()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		entry, ok := lc.GetContext(ctx)
		if !ok {
			// Context cancelled — client disconnected.
			return
		}
		logJSON, _ := json.Marshal(entry)
		if err := conn.WriteMessage(websocket.TextMessage, logJSON); err != nil {
			logger.Warn("write:", err)
			return
		}
	}
}
