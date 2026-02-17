package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	clients        []*Client
	sliceTex       sync.Mutex
	stderrClient   *Client
	cleanup        sync.Once
	stderrFinished chan bool
	namespaces     map[string]bool
	namespacesMux  sync.RWMutex
)

func init() {
	namespaces = make(map[string]bool)
	initColorEnabled()
	stderrClient = CreateClient(DefaultNamespace)
	stderrClient.SetLogLevel(LTrace)
	stderrFinished = make(chan bool, 1)
	go stderrClient.logStdErr()
}

func (c *Client) logStdErr() {
	for e := range c.writer {
		if e.level >= c.LogLevel && c.matchesNamespace(e.Namespace) {
			levelStr := colorizeLevelText(e.Level, e.level)
			nsStr := colorize("["+e.Namespace+"]", colorPurple)
			fileStr := colorize(e.File, colorGray)
			fmt.Fprintf(os.Stderr, "%s\t%s\t%s\t%s\t%s\n", e.Timestamp.String(), levelStr, nsStr, e.Output, fileStr)
		}
	}
	stderrFinished <- true
}

func (c *Client) matchesNamespace(namespace string) bool {
	// Empty Namespaces slice means match all
	if len(c.Namespaces) == 0 {
		return true
	}
	for _, ns := range c.Namespaces {
		if ns == namespace {
			return true
		}
	}
	return false
}

func CreateClient(namespaces ...string) *Client {
	var client Client
	client.initialized = true
	client.Namespaces = namespaces
	client.writer = make(LogWriter, 1000)
	sliceTex.Lock()
	clients = append(clients, &client)
	sliceTex.Unlock()
	return &client
}

func Flush() {
	cleanup.Do(func() {
		close(stderrClient.writer)
		<-stderrFinished
		stderrClient.Destroy()
	})
}

func (c *Client) Destroy() error {
	var otherClients []*Client
	if !c.initialized {
		panic(errors.New("cannot delete uninitialized client, did you use CreateClient?"))
	}
	sliceTex.Lock()
	c.writer = nil
	c.initialized = false
	for _, x := range clients {
		if x.initialized {
			otherClients = append(otherClients, x)
		}
	}
	clients = otherClients
	sliceTex.Unlock()
	return nil
}

func (c *Client) GetLogLevel() Level {
	if !c.initialized {
		panic(errors.New("cannot get level for uninitialized client, use CreateClient instead"))
	}
	return c.LogLevel
}

func createLog(e Entry) {
	// Track namespace
	namespacesMux.Lock()
	namespaces[e.Namespace] = true
	namespacesMux.Unlock()
	
	sliceTex.Lock()
	for _, c := range clients {
		func(c *Client, e Entry) {
			// Filter by namespace if client has filters specified
			if !c.matchesNamespace(e.Namespace) {
				return
			}
			select {
			case c.writer <- e:
				// try to clear out one of the older entries
			default:
				select {
				case <-c.writer:
					c.writer <- e
				default:
				}
			}
		}(c, e)
	}
	sliceTex.Unlock()
}

// GetNamespaces returns a list of all namespaces that have been used
func GetNamespaces() []string {
	namespacesMux.RLock()
	defer namespacesMux.RUnlock()
	
	result := make([]string, 0, len(namespaces))
	for ns := range namespaces {
		result = append(result, ns)
	}
	return result
}

func SetLogLevel(level Level) {
	stderrClient.LogLevel = level
}

// SetLogLevel set log level of logger
func (c *Client) SetLogLevel(level Level) {
	if !c.initialized {
		panic(errors.New("cannot set level for uninitialized client, use CreateClient instead"))
	}
	c.LogLevel = level
}

func (c *Client) Get() Entry {
	if !c.initialized {
		panic(errors.New("cannot get logs for uninitialized client, did you use CreateClient?"))
	}
	return <-c.writer
}

// Trace prints out logs on trace level
func Trace(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "TRACE",
		Namespace: DefaultNamespace,
		level:     LTrace,
	}
	createLog(e)
}

// Formatted print for Trace
func Tracef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Trace prints out logs on trace level with newline
func Traceln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "TRACE",
		level:     LTrace,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Debug prints out logs on debug level
func Debug(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Formatted print for Debug
func Debugf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Debug prints out logs on debug level with a newline
func Debugln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "DEBUG",
		level:     LDebug,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Info prints out logs on info level
func Info(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Formatted print for Info
func Infof(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Info prints out logs on info level with a newline
func Infoln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "INFO",
		level:     LInfo,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Info prints out logs on info level
func Notice(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Formatted print for Info
func Noticef(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Info prints out logs on info level with a newline
func Noticeln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "NOTICE",
		level:     LNotice,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Warn prints out logs on warn level
func Warn(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Formatted print for Warn
func Warnf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Newline print for Warn
func Warnln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "WARN",
		level:     LWarn,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Error prints out logs on error level
func Error(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Formatted print for error
func Errorf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Error prints out logs on error level with a newline
func Errorln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "ERROR",
		level:     LError,
		Namespace: DefaultNamespace,
	}
	createLog(e)
}

// Panic prints out logs on panic level
func Panic(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	if len(args) >= 0 {
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
func Panicf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	if len(args) >= 0 {
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

func Panicln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "PANIC",
		level:     LPanic,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	if len(args) >= 0 {
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
func Fatal(args ...any) {
	output := fmt.Sprint(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

// Formatted print for fatal
func Fatalf(format string, args ...any) {
	output := fmt.Sprintf(format, args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

func Fatalln(args ...any) {
	output := fmt.Sprintln(args...)
	e := Entry{
		Timestamp: time.Now(),
		Output:    output,
		File:      fileInfo(2),
		Level:     "FATAL",
		level:     LFatal,
		Namespace: DefaultNamespace,
	}
	createLog(e)
	Flush()
	os.Exit(1)
}

func Print(args ...any) {
	Info(args...)
}

func Printf(format string, args ...any) {
	Infof(format, args...)
}

func Println(args ...any) {
	Infoln(args...)
}

// fileInfo for getting which line in which file
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
