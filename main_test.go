package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewMuxRegistersRoutes(t *testing.T) {
	mux := newMux()

	tests := []struct {
		name       string
		target     string
		assertions func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "browser viewer",
			target: "/",
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()
				if recorder.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200", recorder.Code)
				}
				if !strings.Contains(recorder.Body.String(), "<!DOCTYPE html>") {
					t.Fatal("expected viewer HTML in response")
				}
			},
		},
		{
			name:   "namespaces api",
			target: "/api/namespaces",
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()
				if recorder.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200", recorder.Code)
				}
				var payload map[string]any
				if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
					t.Fatalf("response is not valid JSON: %v", err)
				}
				if _, ok := payload["namespaces"]; !ok {
					t.Fatal("expected namespaces key in response")
				}
			},
		},
		{
			name:   "websocket endpoint",
			target: "/ws",
			assertions: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				t.Helper()
				if recorder.Code == http.StatusOK || recorder.Code == http.StatusSwitchingProtocols {
					t.Fatalf("expected websocket upgrade failure for plain HTTP request, got %d", recorder.Code)
				}
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, testCase.target, nil)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, request)
			testCase.assertions(t, recorder)
		})
	}
}

func TestNewServerConfiguresReadHeaderTimeout(t *testing.T) {
	server := newServer("127.0.0.1:9999")

	if server.Addr != "127.0.0.1:9999" {
		t.Fatalf("Addr = %q, want %q", server.Addr, "127.0.0.1:9999")
	}
	if server.ReadHeaderTimeout != readHeaderTimeout {
		t.Fatalf("ReadHeaderTimeout = %v, want %v", server.ReadHeaderTimeout, readHeaderTimeout)
	}
	if server.Handler == nil {
		t.Fatal("expected server handler to be configured")
	}
}
