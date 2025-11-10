package log

import (
	"strconv"
	"sync"
	"testing"
)

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

// Print prints out logs on info level
func TestPrint(t *testing.T) {
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
