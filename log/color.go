package log

import (
	"os"
	"sync"
)

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"

	// Bold variants
	colorBoldRed    = "\033[1;31m"
	colorBoldYellow = "\033[1;33m"
	colorBoldWhite  = "\033[1;37m"
)

var (
	colorEnabled     = true
	colorEnabledOnce sync.Once
	colorMux         sync.RWMutex
)

// SetColorEnabled enables or disables colored output for stderr logging.
// By default, color is enabled when stderr is a terminal.
func SetColorEnabled(enabled bool) {
	colorMux.Lock()
	colorEnabled = enabled
	colorMux.Unlock()
}

// ColorEnabled returns whether colored output is currently enabled.
func ColorEnabled() bool {
	colorMux.RLock()
	defer colorMux.RUnlock()
	return colorEnabled
}

// isTerminal checks if the given file descriptor is a terminal.
// This is a simple heuristic that works on Unix-like systems.
func isTerminal(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// initColorEnabled sets the default color state based on whether stderr is a terminal.
func initColorEnabled() {
	colorEnabledOnce.Do(func() {
		colorEnabled = isTerminal(os.Stderr)
	})
}

// levelColor returns the ANSI color code for a given log level.
func levelColor(level Level) string {
	switch level {
	case LTrace:
		return colorGray
	case LDebug:
		return colorCyan
	case LInfo:
		return colorGreen
	case LNotice:
		return colorBlue
	case LWarn:
		return colorYellow
	case LError:
		return colorRed
	case LPanic:
		return colorBoldRed
	case LFatal:
		return colorBoldRed
	default:
		return colorReset
	}
}

// colorize wraps text with ANSI color codes if color is enabled.
func colorize(text string, color string) string {
	if !ColorEnabled() {
		return text
	}
	return color + text + colorReset
}

// colorizeLevelText returns the level string with appropriate color.
func colorizeLevelText(level string, lvl Level) string {
	return colorize(level, levelColor(lvl))
}
