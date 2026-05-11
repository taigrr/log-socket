package browser

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogSocketViewHandler_HTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	body := w.Body.String()
	// html/template escapes forward slashes in JS context
	if !strings.Contains(body, `ws:\/\/localhost:8080\/ws`) {
		t.Error("response should contain escaped ws://localhost:8080/ws URL")
	}
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("response should contain HTML doctype")
	}
}

func TestLogSocketViewHandler_CustomPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://myhost:9090/dashboard/", nil)
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `ws:\/\/myhost:9090\/dashboard\/ws`) {
		t.Error("expected escaped ws://myhost:9090/dashboard/ws in body")
	}
}

func TestLogSocketViewHandler_TrailingSlashTrimmed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	body := w.Body.String()
	// Should NOT have double slash before ws
	if strings.Contains(body, `\/\/ws`) {
		t.Error("should not have double slash before /ws")
	}
	if !strings.Contains(body, `ws:\/\/example.com\/ws`) {
		t.Error("expected escaped ws://example.com/ws in body")
	}
}

func TestLogSocketViewHandler_ForwardedHTTPSUsesWSS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/logs/", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `wss:\/\/example.com\/logs\/ws`) {
		t.Error("expected escaped wss://example.com/logs/ws in body")
	}
}

func TestLogSocketViewHandler_ForwardedHTTPOverridesTLS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/logs/", nil)
	req.TLS = &tls.ConnectionState{}
	req.Header.Set("X-Forwarded-Proto", "http")
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `ws:\/\/example.com\/logs\/ws`) {
		t.Error("expected escaped ws://example.com/logs/ws in body")
	}
}

func TestLogSocketViewHandler_ForwardedHostUsedForWebsocketURL(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://internal:8080/logs/", nil)
	req.TLS = &tls.ConnectionState{}
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "logs.example.com")
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `wss:\/\/logs.example.com\/logs\/ws`) {
		t.Error("expected escaped wss://logs.example.com/logs/ws in body")
	}
	if strings.Contains(body, `internal:8080`) {
		t.Error("response should not leak the internal host when X-Forwarded-Host is set")
	}
}

func TestLogSocketViewHandler_ForwardedHeadersUseFirstValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://internal:8080/logs/", nil)
	req.Header.Set("X-Forwarded-Proto", "HTTPS, http")
	req.Header.Set("X-Forwarded-Host", "logs.example.com, internal:8080")
	w := httptest.NewRecorder()
	LogSocketViewHandler(w, req)

	body := w.Body.String()
	if !strings.Contains(body, `wss:\/\/logs.example.com\/logs\/ws`) {
		t.Error("expected escaped wss://logs.example.com/logs/ws in body")
	}
}
