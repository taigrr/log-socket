package browser

import (
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
