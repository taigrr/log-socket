package log

import "time"

const (
	LTrace Level = iota
	LDebug
	LInfo
	LNotice
	LWarn
	LError
	LPanic
	LFatal
)

const DefaultNamespace = "default"

// String returns the human-readable name of the log level (e.g. "INFO").
// It implements [fmt.Stringer].
func (l Level) String() string {
	switch l {
	case LTrace:
		return "TRACE"
	case LDebug:
		return "DEBUG"
	case LInfo:
		return "INFO"
	case LNotice:
		return "NOTICE"
	case LWarn:
		return "WARN"
	case LError:
		return "ERROR"
	case LPanic:
		return "PANIC"
	case LFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type (
	LogWriter chan Entry
	Level     int

	Client struct {
		LogLevel    Level    `json:"level"`
		Namespaces  []string `json:"namespaces"` // Empty slice means all namespaces
		writer      LogWriter
		initialized bool
	}
	Entry struct {
		Timestamp time.Time `json:"timestamp"`
		Output    string    `json:"output"`
		File      string    `json:"file"`
		Level     string    `json:"level"`
		Namespace string    `json:"namespace"`
		level     Level
	}
	Logger struct {
		FileInfoDepth int
		Namespace     string
	}
)
