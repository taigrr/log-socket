package slog

import (
	"context"
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
	defer c.Destroy()
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
	defer c.Destroy()
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
	defer c.Destroy()
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
