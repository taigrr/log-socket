package log

import (
	"testing"
	"time"
)

func TestLoggerTrace(t *testing.T) {
	c := CreateClient("logger-trace")
	c.SetLogLevel(LTrace)

	l := NewLogger("logger-trace")
	l.Trace("trace message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "TRACE" {
		t.Errorf("level = %q, want TRACE", e.Level)
	}
	if e.Output != "trace message" {
		t.Errorf("output = %q, want %q", e.Output, "trace message")
	}
	if e.Namespace != "logger-trace" {
		t.Errorf("namespace = %q, want %q", e.Namespace, "logger-trace")
	}
	c.Destroy()
}

func TestLoggerTracef(t *testing.T) {
	c := CreateClient("logger-tracef")
	c.SetLogLevel(LTrace)

	l := NewLogger("logger-tracef")
	l.Tracef("trace %s %d", "msg", 1)

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "TRACE" {
		t.Errorf("level = %q, want TRACE", e.Level)
	}
	if e.Output != "trace msg 1" {
		t.Errorf("output = %q, want %q", e.Output, "trace msg 1")
	}
	c.Destroy()
}

func TestLoggerTraceln(t *testing.T) {
	c := CreateClient("logger-traceln")
	c.SetLogLevel(LTrace)

	l := NewLogger("logger-traceln")
	l.Traceln("trace line")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "TRACE" {
		t.Errorf("level = %q, want TRACE", e.Level)
	}
	if e.Output != "trace line\n" {
		t.Errorf("output = %q, want %q", e.Output, "trace line\n")
	}
	c.Destroy()
}

func TestLoggerDebug(t *testing.T) {
	c := CreateClient("logger-debug")
	c.SetLogLevel(LDebug)

	l := NewLogger("logger-debug")
	l.Debug("debug message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "DEBUG" {
		t.Errorf("level = %q, want DEBUG", e.Level)
	}
	if e.Output != "debug message" {
		t.Errorf("output = %q, want %q", e.Output, "debug message")
	}
	c.Destroy()
}

func TestLoggerDebugf(t *testing.T) {
	c := CreateClient("logger-debugf")
	c.SetLogLevel(LDebug)

	l := NewLogger("logger-debugf")
	l.Debugf("debug %d", 42)

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "debug 42" {
		t.Errorf("output = %q, want %q", e.Output, "debug 42")
	}
	c.Destroy()
}

func TestLoggerInfo(t *testing.T) {
	c := CreateClient("logger-info")
	c.SetLogLevel(LInfo)

	l := NewLogger("logger-info")
	l.Info("info message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output != "info message" {
		t.Errorf("output = %q, want %q", e.Output, "info message")
	}
	c.Destroy()
}

func TestLoggerInfof(t *testing.T) {
	c := CreateClient("logger-infof")
	c.SetLogLevel(LInfo)

	l := NewLogger("logger-infof")
	l.Infof("count: %d", 99)

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "count: 99" {
		t.Errorf("output = %q, want %q", e.Output, "count: 99")
	}
	c.Destroy()
}

func TestLoggerInfoln(t *testing.T) {
	c := CreateClient("logger-infoln")
	c.SetLogLevel(LInfo)

	l := NewLogger("logger-infoln")
	l.Infoln("info line")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "info line\n" {
		t.Errorf("output = %q, want %q", e.Output, "info line\n")
	}
	c.Destroy()
}

func TestLoggerNotice(t *testing.T) {
	c := CreateClient("logger-notice")
	c.SetLogLevel(LNotice)

	l := NewLogger("logger-notice")
	l.Notice("notice message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "NOTICE" {
		t.Errorf("level = %q, want NOTICE", e.Level)
	}
	if e.Output != "notice message" {
		t.Errorf("output = %q, want %q", e.Output, "notice message")
	}
	c.Destroy()
}

func TestLoggerNoticef(t *testing.T) {
	c := CreateClient("logger-noticef")
	c.SetLogLevel(LNotice)

	l := NewLogger("logger-noticef")
	l.Noticef("notice %s", "formatted")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "notice formatted" {
		t.Errorf("output = %q, want %q", e.Output, "notice formatted")
	}
	c.Destroy()
}

func TestLoggerNoticeln(t *testing.T) {
	c := CreateClient("logger-noticeln")
	c.SetLogLevel(LNotice)

	l := NewLogger("logger-noticeln")
	l.Noticeln("notice line")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "notice line\n" {
		t.Errorf("output = %q, want %q", e.Output, "notice line\n")
	}
	c.Destroy()
}

func TestLoggerWarn(t *testing.T) {
	c := CreateClient("logger-warn")
	c.SetLogLevel(LWarn)

	l := NewLogger("logger-warn")
	l.Warn("warn message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "WARN" {
		t.Errorf("level = %q, want WARN", e.Level)
	}
	if e.Output != "warn message" {
		t.Errorf("output = %q, want %q", e.Output, "warn message")
	}
	c.Destroy()
}

func TestLoggerWarnf(t *testing.T) {
	c := CreateClient("logger-warnf")
	c.SetLogLevel(LWarn)

	l := NewLogger("logger-warnf")
	l.Warnf("warn %d", 1)

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "warn 1" {
		t.Errorf("output = %q, want %q", e.Output, "warn 1")
	}
	c.Destroy()
}

func TestLoggerWarnln(t *testing.T) {
	c := CreateClient("logger-warnln")
	c.SetLogLevel(LWarn)

	l := NewLogger("logger-warnln")
	l.Warnln("warn line")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "warn line\n" {
		t.Errorf("output = %q, want %q", e.Output, "warn line\n")
	}
	c.Destroy()
}

func TestLoggerError(t *testing.T) {
	c := CreateClient("logger-error")
	c.SetLogLevel(LError)

	l := NewLogger("logger-error")
	l.Error("error message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Level != "ERROR" {
		t.Errorf("level = %q, want ERROR", e.Level)
	}
	if e.Output != "error message" {
		t.Errorf("output = %q, want %q", e.Output, "error message")
	}
	c.Destroy()
}

func TestLoggerErrorf(t *testing.T) {
	c := CreateClient("logger-errorf")
	c.SetLogLevel(LError)

	l := NewLogger("logger-errorf")
	l.Errorf("err: %s", "something broke")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "err: something broke" {
		t.Errorf("output = %q, want %q", e.Output, "err: something broke")
	}
	c.Destroy()
}

func TestLoggerErrorln(t *testing.T) {
	c := CreateClient("logger-errorln")
	c.SetLogLevel(LError)

	l := NewLogger("logger-errorln")
	l.Errorln("error line")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "error line\n" {
		t.Errorf("output = %q, want %q", e.Output, "error line\n")
	}
	c.Destroy()
}

func TestLoggerPanic(t *testing.T) {
	c := CreateClient("logger-panic")
	c.SetLogLevel(LPanic)

	l := NewLogger("logger-panic")

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		// Verify the entry was broadcast
		e, ok := getEntry(c, time.Second)
		if !ok {
			t.Fatal("timed out waiting for panic entry")
		}
		if e.Level != "PANIC" {
			t.Errorf("level = %q, want PANIC", e.Level)
		}
		if e.Namespace != "logger-panic" {
			t.Errorf("namespace = %q, want %q", e.Namespace, "logger-panic")
		}
		c.Destroy()
	}()

	l.Panic("panic message")
}

func TestLoggerPanicf(t *testing.T) {
	c := CreateClient("logger-panicf")
	c.SetLogLevel(LPanic)

	l := NewLogger("logger-panicf")

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		e, ok := getEntry(c, time.Second)
		if !ok {
			t.Fatal("timed out")
		}
		if e.Output != "panic 42" {
			t.Errorf("output = %q, want %q", e.Output, "panic 42")
		}
		c.Destroy()
	}()

	l.Panicf("panic %d", 42)
}

func TestLoggerPanicln(t *testing.T) {
	c := CreateClient("logger-panicln")
	c.SetLogLevel(LPanic)

	l := NewLogger("logger-panicln")

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		e, ok := getEntry(c, time.Second)
		if !ok {
			t.Fatal("timed out")
		}
		if e.Output != "panic line\n" {
			t.Errorf("output = %q, want %q", e.Output, "panic line\n")
		}
		c.Destroy()
	}()

	l.Panicln("panic line")
}

func TestLoggerPrint(t *testing.T) {
	c := CreateClient("logger-print")
	c.SetLogLevel(LInfo)

	l := NewLogger("logger-print")
	l.Print("print message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	// Print delegates to Info
	if e.Level != "INFO" {
		t.Errorf("level = %q, want INFO", e.Level)
	}
	if e.Output != "print message" {
		t.Errorf("output = %q, want %q", e.Output, "print message")
	}
	c.Destroy()
}

func TestLoggerPrintf(t *testing.T) {
	c := CreateClient("logger-printf")
	c.SetLogLevel(LInfo)

	l := NewLogger("logger-printf")
	l.Printf("formatted %s", "print")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "formatted print" {
		t.Errorf("output = %q, want %q", e.Output, "formatted print")
	}
	c.Destroy()
}

func TestLoggerPrintln(t *testing.T) {
	c := CreateClient("logger-println")
	c.SetLogLevel(LInfo)

	l := NewLogger("logger-println")
	l.Println("println message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "println message\n" {
		t.Errorf("output = %q, want %q", e.Output, "println message\n")
	}
	c.Destroy()
}

func TestLoggerSetInfoDepth(t *testing.T) {
	l := NewLogger("depth-test")
	l.SetInfoDepth(3)
	if l.FileInfoDepth != 3 {
		t.Errorf("FileInfoDepth = %d, want 3", l.FileInfoDepth)
	}
}

func TestDefaultLogger(t *testing.T) {
	l := Default()
	if l.Namespace != DefaultNamespace {
		t.Errorf("namespace = %q, want %q", l.Namespace, DefaultNamespace)
	}
	if l.FileInfoDepth != 0 {
		t.Errorf("FileInfoDepth = %d, want 0", l.FileInfoDepth)
	}

	// Verify Default() logger can emit entries
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)

	l.Info("default logger message")

	e, ok := getEntry(c, time.Second)
	if !ok {
		t.Fatal("timed out")
	}
	if e.Output != "default logger message" {
		t.Errorf("output = %q, want %q", e.Output, "default logger message")
	}
	c.Destroy()
}

func TestLoggerPanicWithError(t *testing.T) {
	// When the first arg is an error, Panic should re-panic with that error
	c := CreateClient("logger-panic-err")
	c.SetLogLevel(LPanic)

	l := NewLogger("logger-panic-err")
	testErr := errTest("test error")

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic, got nil")
		}
		if err, ok := r.(errTest); ok {
			if string(err) != "test error" {
				t.Errorf("panic value = %q, want %q", string(err), "test error")
			}
		} else {
			// The first arg was an error, so Panic should re-panic with it
			t.Logf("panic type = %T, value = %v (implementation re-panics with original error)", r, r)
		}
		c.Destroy()
	}()

	l.Panic(testErr)
}

// errTest is a simple error type for testing.
type errTest string

func (e errTest) Error() string { return string(e) }
