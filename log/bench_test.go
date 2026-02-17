package log

import (
	"fmt"
	"testing"
)

const benchNS = "bench-isolated-ns"

// benchClient creates a client with a continuous drain goroutine.
// Returns the client and a stop function to call after b.StopTimer().
func benchClient(ns string) (*Client, func()) {
	c := CreateClient(ns)
	c.SetLogLevel(LTrace)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-c.writer:
			}
		}
	}()
	return c, func() {
		close(done)
		c.Destroy()
	}
}

// BenchmarkCreateClient measures client creation overhead.
func BenchmarkCreateClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := CreateClient("bench")
		c.Destroy()
	}
}

// BenchmarkTrace benchmarks Logger.Trace.
func BenchmarkTrace(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Trace("benchmark trace message")
	}
}

// BenchmarkTracef benchmarks formatted Logger.Tracef.
func BenchmarkTracef(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Tracef("benchmark trace message %d", i)
	}
}

// BenchmarkDebug benchmarks Logger.Debug.
func BenchmarkDebug(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Debug("benchmark debug message")
	}
}

// BenchmarkInfo benchmarks Logger.Info.
func BenchmarkInfo(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("benchmark info message")
	}
}

// BenchmarkInfof benchmarks Logger.Infof with formatting.
func BenchmarkInfof(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Infof("user %s performed action %d", "testuser", i)
	}
}

// BenchmarkWarn benchmarks Logger.Warn.
func BenchmarkWarn(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Warn("benchmark warn message")
	}
}

// BenchmarkError benchmarks Logger.Error.
func BenchmarkError(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Error("benchmark error message")
	}
}

// BenchmarkMultipleClients measures fan-out to multiple consumers.
func BenchmarkMultipleClients(b *testing.B) {
	const numClients = 5
	l := NewLogger(benchNS)
	stops := make([]func(), numClients)
	for i := range stops {
		_, stops[i] = benchClient(benchNS)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("benchmark multi-client")
	}
	b.StopTimer()
	for _, stop := range stops {
		stop()
	}
}

// BenchmarkParallelInfo benchmarks concurrent Info calls from multiple goroutines.
func BenchmarkParallelInfo(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient(benchNS)
	defer stop()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("parallel benchmark info")
		}
	})
}

// BenchmarkNamespaceFiltering benchmarks logging when the consumer
// filters by a different namespace (messages are not delivered).
func BenchmarkNamespaceFiltering(b *testing.B) {
	l := NewLogger(benchNS)
	_, stop := benchClient("completely-different-ns")
	defer stop()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("filtered out message")
	}
}

// BenchmarkFileInfo measures the cost of runtime.Caller for file info.
func BenchmarkFileInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fileInfo(1)
	}
}

// BenchmarkEntryCreation measures raw fmt.Sprint overhead (baseline).
func BenchmarkEntryCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprint("benchmark message ", i)
	}
}
