package ws

import (
	"encoding/json"
	"net"
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
	current := getUpgrader()
	if current.ReadBufferSize != 2048 {
		t.Errorf("ReadBufferSize = %d, want 2048", current.ReadBufferSize)
	}
	if current.WriteBufferSize != 2048 {
		t.Errorf("WriteBufferSize = %d, want 2048", current.WriteBufferSize)
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

	waitForWebSocketEntry(t, conn, func(entry logger.Entry) bool {
		return entry.Namespace == logger.DefaultNamespace && entry.Output == "Websocket client attached."
	})

	// Send a log entry and verify it arrives over the websocket after the
	// client is fully attached.
	testLogger := logger.NewLogger("ws-test")
	testLogger.Info("test message for websocket")

	entry := waitForWebSocketEntry(t, conn, func(entry logger.Entry) bool {
		return entry.Namespace == "ws-test" && entry.Level == "INFO"
	})
	if !strings.Contains(entry.Output, "test message for websocket") {
		t.Errorf("output = %q, want to contain test message", entry.Output)
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

	// Send a log to a different namespace — it should NOT be received.
	otherLogger := logger.NewLogger("other-ns")
	otherLogger.Info("should not arrive")

	// Retry the matching namespace log a few times to avoid racing the server
	// goroutine that registers the filtered client after the websocket dial completes.
	filteredLogger := logger.NewLogger("filtered-ns")
	var entry logger.Entry
	for i := 0; i < 5; i++ {
		filteredLogger.Infof("should arrive %d", i)
		entry = waitForWebSocketEntry(t, conn, func(entry logger.Entry) bool {
			return entry.Namespace == "filtered-ns"
		})
		if entry.Namespace == "filtered-ns" {
			break
		}
	}
	if entry.Namespace != "filtered-ns" {
		t.Fatal("did not receive expected filtered namespace entry")
	}
	if !strings.Contains(entry.Output, "should arrive") {
		t.Errorf("output = %q, want to contain 'should arrive'", entry.Output)
	}
}

func waitForWebSocketEntry(t *testing.T, conn *websocket.Conn, match func(logger.Entry) bool) logger.Entry {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	for i := 0; i < 10; i++ {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				return logger.Entry{}
			}
			t.Fatalf("failed to read message: %v", err)
		}
		var entry logger.Entry
		if err := json.Unmarshal(message, &entry); err != nil {
			t.Fatalf("failed to unmarshal entry: %v", err)
		}
		if match(entry) {
			return entry
		}
	}
	return logger.Entry{}
}
