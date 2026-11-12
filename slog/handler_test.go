package slog

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/taigrr/log-socket/v2/log"
)

// getWithTimeout reads from a log client with a timeout to avoid hanging tests.
func getWithTimeout(c *log.Client, timeout time.Duration) (log.Entry, bool) {
	ch := make(chan log.Entry, 1)
	go func() { ch <- c.Get() }()
	select {
	case e := <-ch:
		return e, true
	case <-time.After(timeout):
		return log.Entry{}, false
	}
}

func TestHandler_Enabled(t *testing.T) {
	h := NewHandler(WithLevel(slog.LevelWarn))
	if h.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected Info to be disabled when min level is Warn")
	}
	if !h.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("expected Warn to be enabled")
	}
	if !h.Enabled(context.Background(), slog.LevelError) {
		t.Error("expected Error to be enabled")
	}
}

func TestHandler_Handle(t *testing.T) {
	c := log.CreateClient()
	defer func() { _ = c.Destroy() }() // Destroy currently always returns nil.
	c.SetLogLevel(log.LTrace)

	h := NewHandler(WithNamespace("test-ns"))
	logger := slog.New(h)

	logger.Info("hello world", "key", "value")

	e, ok := getWithTimeout(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for log entry")
	}
	if e.Namespace != "test-ns" {
		t.Errorf("namespace = %q, want %q", e.Namespace, "test-ns")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output == "" {
		t.Error("output should not be empty")
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	c := log.CreateClient()
	defer func() { _ = c.Destroy() }() // Destroy currently always returns nil.
	c.SetLogLevel(log.LTrace)

	h := NewHandler()
	h2 := h.WithAttrs([]slog.Attr{slog.String("service", "api")})
	logger := slog.New(h2)

	logger.Info("request")

	e, ok := getWithTimeout(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "request service=api" {
		t.Errorf("output = %q, want %q", e.Output, "request service=api")
	}
}

func TestHandler_WithGroup(t *testing.T) {
	c := log.CreateClient()
	defer func() { _ = c.Destroy() }() // Destroy currently always returns nil.
	c.SetLogLevel(log.LTrace)

	h := NewHandler()
	h2 := h.WithGroup("http").WithAttrs([]slog.Attr{slog.Int("status", 200)})
	logger := slog.New(h2)

	logger.Info("done")

	e, ok := getWithTimeout(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "done http.status=200" {
		t.Errorf("output = %q, want %q", e.Output, "done http.status=200")
	}
}

func TestHandler_WithAttrsBeforeWithGroup(t *testing.T) {
	c := log.CreateClient()
	defer func() { _ = c.Destroy() }() // Destroy currently always returns nil.
	c.SetLogLevel(log.LTrace)

	h := NewHandler()
	h2 := h.WithAttrs([]slog.Attr{slog.String("service", "api")}).
		WithGroup("http").
		WithAttrs([]slog.Attr{slog.Int("status", 200)})
	logger := slog.New(h2)

	logger.Info("done", "method", "GET")

	e, ok := getWithTimeout(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "done service=api http.status=200 http.method=GET" {
		t.Errorf("output = %q, want %q", e.Output, "done service=api http.status=200 http.method=GET")
	}
}

func TestHandler_WithNamespaceEmptyUsesDefault(t *testing.T) {
	c := log.CreateClient(log.DefaultNamespace)
	defer func() { _ = c.Destroy() }() // Destroy currently always returns nil.
	c.SetLogLevel(log.LTrace)

	h := NewHandler(WithNamespace(""))
	logger := slog.New(h)

	logger.Info("default namespace")

	e, ok := getWithTimeout(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Namespace != log.DefaultNamespace {
		t.Errorf("namespace = %q, want %q", e.Namespace, log.DefaultNamespace)
	}
}

func TestSlogLevelMapping(t *testing.T) {
	tests := []struct {
		level slog.Level
		want  string
	}{
		{slog.LevelDebug, "DEBUG"},
		{slog.LevelInfo, "INFO"},
		{slog.LevelWarn, "WARN"},
		{slog.LevelError, "ERROR"},
	}
	for _, tt := range tests {
		got := slogLevelToString(tt.level)
		if got != tt.want {
			t.Errorf("slogLevelToString(%v) = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestHandler_GroupAttrs(t *testing.T) {
	tests := []struct {
		name string
		attr slog.Attr
		want string
	}{
		{"nested", slog.Group("request", slog.Int("status", 200), slog.Group("peer", slog.String("host", "localhost"))), "message outer.request.status=200 outer.request.peer.host=localhost"},
		{"inline", slog.Group("", slog.Int("status", 200)), "message outer.status=200"},
		{"empty", slog.Group("request"), "message"},
		{"empty children", slog.Group("request", slog.Attr{}, slog.Group("peer")), "message"},
	}
	for _, tt := range tests {
		for _, bound := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/bound=%t", tt.name, bound), func(t *testing.T) {
				client := log.CreateClient(t.Name())
				defer func() { _ = client.Destroy() }() // Destroy currently always returns nil.
				handler := NewHandler(WithNamespace(t.Name())).WithGroup("outer")
				record := slog.NewRecord(time.Now(), slog.LevelInfo, "message", 0)
				if bound {
					handler = handler.WithAttrs([]slog.Attr{tt.attr})
				} else {
					record.AddAttrs(tt.attr)
				}
				if err := handler.Handle(context.Background(), record); err != nil {
					t.Fatal(err)
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				entry, ok := client.GetContext(ctx)
				if !ok {
					t.Fatal("no log entry received")
				}
				if entry.Output != tt.want {
					t.Fatalf("output = %q, want %q", entry.Output, tt.want)
				}
			})
		}
	}
}
