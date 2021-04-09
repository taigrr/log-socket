package lambo_log_socket

import (
	"testing"
)

func TestCreateDestroy(t *testing.T) {
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, but found %d", len(clients))
	}
	c := CreateClient()
	if len(clients) != 2 {
		t.Errorf("Expected 2 clients, but found %d", len(clients))
	}
	c.Destroy()
	if len(clients) != 1 {
		t.Errorf("Expected 1 client, but found %d", len(clients))
	}
}

// SetLogLevel set log level of logger
func TestSetLogLevel(t *testing.T) {
	logLevels := [...]Level{LTrace, LDebug, LInfo, LWarn, LError, LPanic, LFatal}
	c := CreateClient()
	for _, x := range logLevels {
		c.SetLogLevel(x)
		if c.GetLogLevel() != x {
			t.Errorf("Got logLevel %d, but expected %d", int(c.GetLogLevel()), int(x))
		}
	}
	c.Destroy()
}

// Trace prints out logs on trace level
func TestTrace(t *testing.T) {
	testString := "Testing trace!"
	var c *Client
	c = CreateClient()
	c.SetLogLevel(LTrace)

	for i := 0; i < 5; i++ {
		Trace(testString)
	}
	for i := 0; i < 5; i++ {
		if testString != c.Get().Output {
			t.Error("Trace input doesn't match output")
		}
	}
}

// Debug prints out logs on debug level
func TestDebug(t *testing.T) {
	Debug("Test of Debug")
	//	if logLevel >= LDebug {
	//		entry := logger.WithFields(logrus.Fields{})
	//		entry.Data["file"] = fileInfo(2)
	//		entry.Debug(args...)
	//	}
}

// Info prints out logs on info level
func TestInfo(t *testing.T) {
	//	if logLevel >= LInfo {
	//		entry := logger.WithFields(logrus.Fields{})
	//		entry.Data["file"] = fileInfo(2)
	//		entry.Info(args...)
	//	}
}

// Warn prints out logs on warn level
func TestWarn(t *testing.T) {
	//	if logLevel >= LWarn {
	//		entry := logger.WithFields(logrus.Fields{})
	//		entry.Data["file"] = fileInfo(2)
	//		entry.Warn(args...)
	//	}
}

// Error prints out logs on error level
func TestError(t *testing.T) {
	//	if logLevel >= LError {
	//		entry := logger.WithFields(logrus.Fields{})
	//		entry.Data["file"] = fileInfo(2)
	//		entry.Error(args...)
	//	}
}

// Fatal prints out logs on fatal level
func TestFatal(t *testing.T) {
	//	if logLevel >= LFatal {
	//		entry := logger.WithFields(logrus.Fields{})
	//		entry.Data["file"] = fileInfo(2)
	//		entry.Fatal(args...)
	//	}
}

// Panic prints out logs on panic level
func TestPanic(t *testing.T) {
	//	if logLevel >= LPanic {
	//		entry := logger.WithFields(logrus.Fields{})
	//		entry.Data["file"] = fileInfo(2)
	//		entry.Panic(args...)
	//	}
}
func TestFlush(t *testing.T) {
	defer Flush()
}
