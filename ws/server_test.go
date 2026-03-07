package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	logger "github.com/taigrr/log-socket/v2/log"
)

func TestSetUpgrader(t *testing.T) {
	custom := websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}
	SetUpgrader(custom)
	if upgrader.ReadBufferSize != 2048 {
		t.Errorf("ReadBufferSize = %d, want 2048", upgrader.ReadBufferSize)
	}
	if upgrader.WriteBufferSize != 2048 {
		t.Errorf("WriteBufferSize = %d, want 2048", upgrader.WriteBufferSize)
	}
	// Reset to default
	SetUpgrader(websocket.Upgrader{})
}

func TestLogSocketHandler_NonWebSocket(t *testing.T) {
	// A non-upgrade request should fail gracefully (upgrader returns error)
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	w := httptest.NewRecorder()
	LogSocketHandler(w, req)
	// The upgrader should return a 400-level error for non-websocket requests
	if w.Code == http.StatusOK || w.Code == http.StatusSwitchingProtocols {
		t.Errorf("expected error status for non-websocket request, got %d", w.Code)
	}
}

func TestLogSocketHandler_WebSocket(t *testing.T) {
	// Set upgrader with permissive origin check for testing
	SetUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	defer SetUpgrader(websocket.Upgrader{})

	server := httptest.NewServer(http.HandlerFunc(LogSocketHandler))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Send a log entry and verify it arrives over the websocket
	testLogger := logger.NewLogger("ws-test")
	testLogger.Info("test message for websocket")

	// Read messages until we find our test entry (the handler itself
	// logs "Websocket client attached." which may arrive first)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var found bool
	for i := 0; i < 10; i++ {
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read message: %v", err)
		}
		var entry logger.Entry
		if err := json.Unmarshal(message, &entry); err != nil {
			t.Fatalf("failed to unmarshal entry: %v", err)
		}
		if entry.Namespace == "ws-test" && entry.Level == "INFO" {
			found = true
			break
		}
	}
	if !found {
		t.Error("did not receive expected log entry with namespace ws-test")
	}
}

func TestLogSocketHandler_NamespaceFilter(t *testing.T) {
	SetUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	})
	defer SetUpgrader(websocket.Upgrader{})

	server := httptest.NewServer(http.HandlerFunc(LogSocketHandler))
	defer server.Close()

	// Connect with namespace filter
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?namespaces=filtered-ns"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Send a log to a different namespace — it should NOT be received
	otherLogger := logger.NewLogger("other-ns")
	otherLogger.Info("should not arrive")

	// Send a log to the filtered namespace — it SHOULD be received
	filteredLogger := logger.NewLogger("filtered-ns")
	filteredLogger.Info("should arrive")

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	var entry logger.Entry
	if err := json.Unmarshal(message, &entry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}
	if entry.Namespace != "filtered-ns" {
		t.Errorf("namespace = %q, want filtered-ns", entry.Namespace)
	}
	if !strings.Contains(entry.Output, "should arrive") {
		t.Errorf("output = %q, want to contain 'should arrive'", entry.Output)
	}
}
