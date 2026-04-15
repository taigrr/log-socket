package log

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"
)

// getEntry reads from a client with a timeout to avoid hanging tests.
func getEntry(c *Client, timeout time.Duration) (Entry, bool) {
	ch := make(chan Entry, 1)
	go func() { ch <- c.Get() }()
	select {
	case e := <-ch:
		return e, true
	case <-time.After(timeout):
		return Entry{}, false
	}
}

// Test CreateClient() and Client.Destroy()
func TestCreateDestroy(t *testing.T) {
	// Ensure only stderr exists at the beginning
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, but found %d", len(clients))
	}
	// Create a new client, ensure it's added
	c := CreateClient("test")
	if len(clients) != 2 {
		t.Errorf("Expected 2 clients, but found %d", len(clients))
	}
	// Destroy it and ensure it's actually removed from the array
	c.Destroy()
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, but found %d", len(clients))
	}
}

// SetLogLevel set log level of logger
func TestSetLogLevel(t *testing.T) {
	logLevels := [...]Level{LTrace, LDebug, LInfo, LWarn, LError, LPanic, LFatal}
	c := CreateClient("test")
	for _, x := range logLevels {
		c.SetLogLevel(x)
		if c.GetLogLevel() != x {
			t.Errorf("Got logLevel %d, but expected %d", int(c.GetLogLevel()), int(x))
		}
	}
	c.Destroy()
}

func BenchmarkDebugSerial(b *testing.B) {
	c := CreateClient("test")
	var x sync.WaitGroup
	x.Add(b.N)
	for i := 0; i < b.N; i++ {
		Debug(i)
		go func() {
			c.Get()
			x.Done()
		}()
	}
	x.Wait()
	c.Destroy()
}

// Trace ensure logs come out in the right order
func TestOrder(t *testing.T) {
	testString := "Testing trace: "
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LTrace)

	for i := 0; i < 5000; i++ {
		Trace(testString + strconv.Itoa(i))
		if testString+strconv.Itoa(i) != c.Get().Output {
			t.Error("Trace input doesn't match output")
		}
	}
	c.Destroy()
}

func TestDebug(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LDebug)

	Debug("debug message")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for debug entry")
	}
	if e.Level != "DEBUG" {
		t.Errorf("level = %q, want DEBUG", e.Level)
	}
	if e.Output != "debug message" {
		t.Errorf("output = %q, want %q", e.Output, "debug message")
	}
	if e.Namespace != DefaultNamespace {
		t.Errorf("namespace = %q, want %q", e.Namespace, DefaultNamespace)
	}
	c.Destroy()
}

func TestDebugf(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LDebug)

	Debugf("hello %s %d", "world", 42)
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "hello world 42" {
		t.Errorf("output = %q, want %q", e.Output, "hello world 42")
	}
	c.Destroy()
}

func TestInfo(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)

	Info("info message")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for info entry")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output != "info message" {
		t.Errorf("output = %q, want %q", e.Output, "info message")
	}
	c.Destroy()
}

func TestInfof(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)

	Infof("count: %d", 99)
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "count: 99" {
		t.Errorf("output = %q, want %q", e.Output, "count: 99")
	}
	c.Destroy()
}

func TestPrint(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)

	Print("print message")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	// Print is an alias for Info
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output != "print message" {
		t.Errorf("output = %q, want %q", e.Output, "print message")
	}
	c.Destroy()
}

func TestPrintf(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)

	Printf("formatted %s", "print")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "formatted print" {
		t.Errorf("output = %q, want %q", e.Output, "formatted print")
	}
	c.Destroy()
}

func TestNotice(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LNotice)

	Notice("notice message")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "NOTICE" {
		t.Errorf("level = %q, want NOTICE", e.Level)
	}
	c.Destroy()
}

func TestWarn(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LWarn)

	Warn("warning message")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for warn entry")
	}
	if e.Level != "WARN" {
		t.Errorf("level = %q, want WARN", e.Level)
	}
	if e.Output != "warning message" {
		t.Errorf("output = %q, want %q", e.Output, "warning message")
	}
	c.Destroy()
}

func TestWarnf(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LWarn)

	Warnf("warn %d", 1)
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "warn 1" {
		t.Errorf("output = %q, want %q", e.Output, "warn 1")
	}
	c.Destroy()
}

func TestError(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LError)

	Error("error message")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for error entry")
	}
	if e.Level != "ERROR" {
		t.Errorf("level = %q, want ERROR", e.Level)
	}
	if e.Output != "error message" {
		t.Errorf("output = %q, want %q", e.Output, "error message")
	}
	c.Destroy()
}

func TestErrorf(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LError)

	Errorf("err: %s", "something broke")
	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "err: something broke" {
		t.Errorf("output = %q, want %q", e.Output, "err: something broke")
	}
	c.Destroy()
}

func TestPanic(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LPanic)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		c.Destroy()
	}()

	Panic("panic message")
}

func TestPanicf(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LPanic)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		c.Destroy()
	}()

	Panicf("panic %d", 42)
}

func TestPanicln(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LPanic)

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		c.Destroy()
	}()

	Panicln("panic line")
}

// TestLogLevelFiltering verifies that the client's log level is stored correctly.
// Note: level filtering only applies to stderr output, not to client channels.
// All entries matching the namespace are delivered to the client channel regardless of level.
func TestLogLevelFiltering(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LWarn)

	if c.GetLogLevel() != LWarn {
		t.Errorf("expected log level LWarn, got %d", c.GetLogLevel())
	}

	// Both entries arrive at the client channel (level filtering is stderr-only)
	Info("info message")
	Warn("warn message")

	e1, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for first entry")
	}
	if e1.Output != "info message" {
		t.Errorf("expected 'info message', got %q", e1.Output)
	}

	e2, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for second entry")
	}
	if e2.Output != "warn message" {
		t.Errorf("expected 'warn message', got %q", e2.Output)
	}
	c.Destroy()
}

// TestNamespaceFiltering verifies clients only receive matching namespaces.
func TestNamespaceFiltering(t *testing.T) {
	c := CreateClient("api")
	c.SetLogLevel(LTrace)

	apiLogger := NewLogger("api")
	dbLogger := NewLogger("database")

	// Log to database namespace — should not arrive at "api" client
	dbLogger.Info("db message")

	// Log to api namespace — should arrive
	apiLogger.Info("api message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for api entry")
	}
	if e.Output != "api message" {
		t.Errorf("expected 'api message', got %q", e.Output)
	}
	if e.Namespace != "api" {
		t.Errorf("namespace = %q, want 'api'", e.Namespace)
	}
	c.Destroy()
}

// TestMultiNamespaceClient verifies a client subscribed to multiple namespaces.
func TestMultiNamespaceClient(t *testing.T) {
	c := CreateClient("api", "auth")
	c.SetLogLevel(LTrace)

	apiLogger := NewLogger("api")
	authLogger := NewLogger("auth")
	dbLogger := NewLogger("database")

	dbLogger.Info("db message")     // filtered out
	apiLogger.Info("api message")   // should arrive
	authLogger.Info("auth message") // should arrive

	e1, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for first entry")
	}
	if e1.Output != "api message" {
		t.Errorf("first entry = %q, want 'api message'", e1.Output)
	}

	e2, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out waiting for second entry")
	}
	if e2.Output != "auth message" {
		t.Errorf("second entry = %q, want 'auth message'", e2.Output)
	}
	c.Destroy()
}

// TestGetNamespaces verifies the namespace registry.
func TestGetNamespaces(t *testing.T) {
	l := NewLogger("test-ns-registry")
	l.Info("register this namespace")

	nss := GetNamespaces()
	found := false
	for _, ns := range nss {
		if ns == "test-ns-registry" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'test-ns-registry' in GetNamespaces(), got %v", nss)
	}
}

// TestLoggerDebugln verifies the Debugln method on Logger.
func TestLoggerDebugln(t *testing.T) {
	c := CreateClient("debugln-test")
	c.SetLogLevel(LDebug)

	l := NewLogger("debugln-test")
	l.Debugln("debugln message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "DEBUG" {
		t.Errorf("level = %q, want DEBUG", e.Level)
	}
	// Sprintln appends a newline
	if e.Output != "debugln message\n" {
		t.Errorf("output = %q, want %q", e.Output, "debugln message\n")
	}
	c.Destroy()
}

// TestNewLoggerEmptyNamespace verifies empty namespace defaults to DefaultNamespace.
func TestNewLoggerEmptyNamespace(t *testing.T) {
	l := NewLogger("")
	if l.Namespace != DefaultNamespace {
		t.Errorf("namespace = %q, want %q", l.Namespace, DefaultNamespace)
	}
}

// TestFileInfo verifies fileInfo returns a non-empty file:line string.
func TestFileInfo(t *testing.T) {
	fi := fileInfo(1)
	if fi == "" || fi == "<???>:1" {
		t.Errorf("fileInfo returned unexpected value: %q", fi)
	}
}

// TestColorize verifies color wrapping.
func TestColorize(t *testing.T) {
	SetColorEnabled(true)
	result := colorize("hello", colorRed)
	expected := colorRed + "hello" + colorReset
	if result != expected {
		t.Errorf("colorize with color enabled: got %q, want %q", result, expected)
	}

	SetColorEnabled(false)
	result = colorize("hello", colorRed)
	if result != "hello" {
		t.Errorf("colorize with color disabled: got %q, want %q", result, "hello")
	}

	// Restore default
	SetColorEnabled(true)
}

// TestParseLevelString verifies level string parsing.
func TestParseLevelString(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"TRACE", LTrace},
		{"DEBUG", LDebug},
		{"INFO", LInfo},
		{"NOTICE", LNotice},
		{"WARN", LWarn},
		{"ERROR", LError},
		{"PANIC", LPanic},
		{"FATAL", LFatal},
		{"UNKNOWN", LInfo}, // default
		{"", LInfo},        // default
	}
	for _, tt := range tests {
		got := parseLevelString(tt.input)
		if got != tt.want {
			t.Errorf("parseLevelString(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

// TestBroadcast verifies the public Broadcast function.
func TestBroadcast(t *testing.T) {
	c := CreateClient("broadcast-ns")
	c.SetLogLevel(LTrace)

	e := Entry{
		Timestamp: time.Now(),
		Output:    "broadcast test",
		File:      "test.go:1",
		Level:     "WARN",
		Namespace: "broadcast-ns",
	}
	Broadcast(e)

	got, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if got.Output != "broadcast test" {
		t.Errorf("output = %q, want %q", got.Output, "broadcast test")
	}
	if got.Level != "WARN" {
		t.Errorf("level = %q, want WARN", got.Level)
	}
	c.Destroy()
}

// TestMatchesNamespace verifies the namespace matching helper.
func TestMatchesNamespace(t *testing.T) {
	// Client with no namespace filter matches everything
	c := CreateClient()
	if !c.matchesNamespace("anything") {
		t.Error("empty Namespaces should match all")
	}
	c.Destroy()

	// Client with specific namespaces
	c2 := CreateClient("api", "auth")
	if !c2.matchesNamespace("api") {
		t.Error("should match 'api'")
	}
	if !c2.matchesNamespace("auth") {
		t.Error("should match 'auth'")
	}
	if c2.matchesNamespace("database") {
		t.Error("should not match 'database'")
	}
	c2.Destroy()
}

// TestGetContext verifies context cancellation stops blocking Get.
func TestGetContext(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LTrace)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, ok := c.GetContext(ctx)
	if ok {
		t.Error("expected GetContext to return false on cancelled context")
	}
	c.Destroy()
}

// TestGetContextReceivesEntry verifies GetContext delivers entries normally.
func TestGetContextReceivesEntry(t *testing.T) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LTrace)

	go func() {
		time.Sleep(10 * time.Millisecond)
		Info("context entry")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	e, ok := c.GetContext(ctx)
	if !ok {
		t.Fatal("expected GetContext to return entry")
	}
	if e.Output != "context entry" {
		t.Errorf("output = %q, want %q", e.Output, "context entry")
	}
	c.Destroy()
}

// TestLevelString verifies the Level.String() method.
func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LTrace, "TRACE"},
		{LDebug, "DEBUG"},
		{LInfo, "INFO"},
		{LNotice, "NOTICE"},
		{LWarn, "WARN"},
		{LError, "ERROR"},
		{LPanic, "PANIC"},
		{LFatal, "FATAL"},
		{Level(99), "UNKNOWN"},
	}
	for _, tt := range tests {
		got := tt.level.String()
		if got != tt.want {
			t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestFlush(t *testing.T) {
	defer Flush()
}
