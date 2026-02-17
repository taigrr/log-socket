package log

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func Default() *Logger {
	return &Logger{FileInfoDepth: 0, Namespace: DefaultNamespace}
}

func NewLogger(namespace string) *Logger {
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return &Logger{FileInfoDepth: 0, Namespace: namespace}
}

func (l *Logger) SetInfoDepth(depth int) {
	l.FileInfoDepth = depth
}

// Trace prints out logs on trace level
func (l Logger) Trace(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Formatted print for Trace
func (l Logger) Tracef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Trace prints out logs on trace level with newline
func (l Logger) Traceln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Debug prints out logs on debug level
func (l Logger) Debug(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Formatted print for Debug
func (l Logger) Debugf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Info prints out logs on info level
func (l Logger) Info(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Formatted print for Info
func (l Logger) Infof(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Info prints out logs on info level with newline
func (l Logger) Infoln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "INFO",
		level:     LInfo,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Notice prints out logs on notice level
func (l Logger) Notice(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Formatted print for Notice
func (l Logger) Noticef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Notice prints out logs on notice level with newline
func (l Logger) Noticeln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Warn prints out logs on warn level
func (l Logger) Warn(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Formatted print for Warn
func (l Logger) Warnf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Warn prints out logs on warn level with a newline
func (l Logger) Warnln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "WARN",
		level:     LWarn,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Error prints out logs on error level
func (l Logger) Error(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Formatted print for error
func (l Logger) Errorf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Error prints out logs on error level with a new line
func (l Logger) Errorln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "ERROR",
		level:     LError,
		Namespace: l.Namespace,
	}
	createLog(e)
}

// Panic prints out logs on panic level
func (l Logger) Panic(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: l.Namespace,
	}
	createLog(e)
	if len(args) > 0 {
		switch args[0].(type) {
		case error:
			panic(args[0])
		default:
			// falls through to default below
		}
	}
	Flush()
	panic(errors.New(output))
}

// Formatted print for panic
func (l Logger) Panicf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: l.Namespace,
	}
	createLog(e)
	if len(args) > 0 {
		switch args[0].(type) {
		case error:
			panic(args[0])
		default:
			// falls through to default below
		}
	}
	Flush()
	panic(errors.New(output))
}

// Panic prints out logs on panic level with a newline
func (l Logger) Panicln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: l.Namespace,
	}
	createLog(e)
	if len(args) > 0 {
		switch args[0].(type) {
		case error:
			panic(args[0])
		default:
			// falls through to default below
		}
	}
	Flush()
	panic(errors.New(output))
}

// Fatal prints out logs on fatal level
func (l Logger) Fatal(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: l.Namespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Formatted print for fatal
func (l Logger) Fatalf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: l.Namespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Fatal prints fatal level with a new line
func (l Logger) Fatalln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2 + l.FileInfoDepth),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: l.Namespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Handles print to info
func (l Logger) Print(args ...any) {
	l.Info(args...)
}

// Handles formatted print to info
func (l Logger) Printf(format string, args ...any) {
	l.Infof(format, args...)
}

// Handles print to info with new line
func (l Logger) Println(args ...any) {
	l.Infoln(args...)
}
