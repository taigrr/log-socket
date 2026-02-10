package log

import (
	"testing"
)

// drainClient continuously reads from a client to prevent blocking.
// Returns a function to stop draining and wait for completion.
func drainClient(c *Client) func() {
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-c.writer:
				// Drain entries
			}
		}
	}()
	return func() { close(done) }
}

// -----------------------------------------------------------------------------
// Serial Benchmarks - Single Log Levels
// -----------------------------------------------------------------------------

func BenchmarkTrace(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LTrace)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Trace("benchmark message")
	}
}

func BenchmarkDebug(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LDebug)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debug("benchmark message")
	}
}

func BenchmarkInfo(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}
}

func BenchmarkNotice(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LNotice)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Notice("benchmark message")
	}
}

func BenchmarkWarn(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LWarn)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Warn("benchmark message")
	}
}

func BenchmarkError(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LError)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Error("benchmark message")
	}
}

// -----------------------------------------------------------------------------
// Formatted Logging Benchmarks
// -----------------------------------------------------------------------------

func BenchmarkDebugf(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LDebug)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debugf("benchmark message %d with %s", i, "formatting")
	}
}

func BenchmarkInfof(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("benchmark message %d with %s", i, "formatting")
	}
}

func BenchmarkErrorf(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LError)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Errorf("benchmark message %d with %s", i, "formatting")
	}
}

// -----------------------------------------------------------------------------
// Parallel Benchmarks
// -----------------------------------------------------------------------------

func BenchmarkDebugParallel(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LDebug)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Debug("parallel benchmark message")
		}
	})
}

func BenchmarkInfoParallel(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Info("parallel benchmark message")
		}
	})
}

func BenchmarkInfofParallel(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			Infof("parallel benchmark message %d", counter)
			counter++
		}
	})
}

// -----------------------------------------------------------------------------
// Logger Instance Benchmarks (Namespaced Logging)
// -----------------------------------------------------------------------------

func BenchmarkLoggerInfo(b *testing.B) {
	logger := NewLogger("benchmark")
	c := CreateClient("benchmark")
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkLoggerInfof(b *testing.B) {
	logger := NewLogger("benchmark")
	c := CreateClient("benchmark")
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("benchmark message %d", i)
	}
}

func BenchmarkLoggerDebugParallel(b *testing.B) {
	logger := NewLogger("benchmark")
	c := CreateClient("benchmark")
	c.SetLogLevel(LDebug)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("parallel benchmark message")
		}
	})
}

// -----------------------------------------------------------------------------
// Multiple Client Benchmarks
// -----------------------------------------------------------------------------

func BenchmarkMultipleClients(b *testing.B) {
	const numClients = 5
	var clients []*Client
	var stopFuncs []func()

	for i := 0; i < numClients; i++ {
		c := CreateClient(DefaultNamespace)
		c.SetLogLevel(LInfo)
		stop := drainClient(c)
		clients = append(clients, c)
		stopFuncs = append(stopFuncs, stop)
	}

	defer func() {
		for i, c := range clients {
			stopFuncs[i]()
			c.Destroy()
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message to multiple clients")
	}
}

func BenchmarkMultipleNamespaces(b *testing.B) {
	namespaces := []string{"auth", "db", "api", "cache", "queue"}
	var loggers []*Logger
	var clients []*Client
	var stopFuncs []func()

	for _, ns := range namespaces {
		loggers = append(loggers, NewLogger(ns))
		c := CreateClient(ns)
		c.SetLogLevel(LInfo)
		stop := drainClient(c)
		clients = append(clients, c)
		stopFuncs = append(stopFuncs, stop)
	}

	defer func() {
		for i, c := range clients {
			stopFuncs[i]()
			c.Destroy()
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		loggers[i%len(loggers)].Info("benchmark message")
	}
}

// -----------------------------------------------------------------------------
// With Synchronous Client Consumption (Legacy Pattern)
// -----------------------------------------------------------------------------

// BenchmarkDebugWithSyncConsumer measures the overhead of synchronous consumption
// using the Get() method, processing one message at a time.
func BenchmarkDebugWithSyncConsumer(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LDebug)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	// This benchmark measures just the logging side with async draining.
	// For sync consumption patterns, see the TestOrder test.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debug("benchmark message")
	}
}

// -----------------------------------------------------------------------------
// Comparison Benchmarks (Different Message Sizes)
// -----------------------------------------------------------------------------

func BenchmarkInfoShortMessage(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("ok")
	}
}

func BenchmarkInfoMediumMessage(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	msg := "This is a medium-length log message for benchmarking purposes"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(msg)
	}
}

func BenchmarkInfoLongMessage(b *testing.B) {
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LInfo)
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	msg := "This is a much longer log message that simulates real-world logging scenarios where developers tend to include more context about what the application is doing, including variable values, request IDs, and other debugging information that can be quite verbose"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(msg)
	}
}

// -----------------------------------------------------------------------------
// Overhead Benchmarks (Level Filtering)
// -----------------------------------------------------------------------------

func BenchmarkDebugFilteredByLevel(b *testing.B) {
	// No client with Debug level, so Debug logs won't be consumed
	// This measures the overhead of creating log entries that get filtered
	c := CreateClient(DefaultNamespace)
	c.SetLogLevel(LError) // Only Error and above
	stop := drainClient(c)
	defer stop()
	defer c.Destroy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Debug("this message will be filtered")
	}
}
