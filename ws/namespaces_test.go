package ws

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	logger "github.com/taigrr/log-socket/v2/log"
)

func TestNamespacesHandler(t *testing.T) {
	// Log to a known namespace to ensure it appears
	nsLogger := logger.NewLogger("ns-handler-test")
	nsLogger.Info("register namespace")

	req := httptest.NewRequest(http.MethodGet, "/api/namespaces", nil)
	w := httptest.NewRecorder()
	NamespacesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var result struct {
		Namespaces []string `json:"namespaces"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	found := false
	for _, ns := range result.Namespaces {
		if ns == "ns-handler-test" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("namespace 'ns-handler-test' not found in %v", result.Namespaces)
	}
}

func TestNamespacesHandler_ResponseFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/namespaces", nil)
	w := httptest.NewRecorder()
	NamespacesHandler(w, req)

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if _, ok := result["namespaces"]; !ok {
		t.Error("response missing 'namespaces' key")
	}
}
