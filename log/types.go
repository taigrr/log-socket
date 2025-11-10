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
