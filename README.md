# Log Socket v2

A real-time log viewer with WebSocket support and namespace filtering, written in Go.

## What's New in v2

**Breaking Changes:**
- Module path changed to `github.com/taigrr/log-socket/v2`
- `CreateClient()` now accepts variadic `namespaces ...string` parameter
- `Logger` type now includes `Namespace` field
- New `NewLogger(namespace string)` constructor for namespaced loggers

**New Features:**
- **Namespace support**: Organize logs by namespace (api, database, auth, etc.)
- **Namespace filtering**: Subscribe to specific namespaces via WebSocket
- **Frontend namespace selector**: Filter logs by namespace in the web UI
- **Namespace API**: GET `/api/namespaces` to list all active namespaces
- **Colored terminal output**: Log levels are color-coded in terminal (no external packages)

## Features

- Real-time log streaming via WebSocket
- Web-based log viewer with filtering capabilities
- **Namespace-based log organization**
- Support for multiple log levels (TRACE, DEBUG, INFO, WARN, ERROR, PANIC, FATAL)
- Color-coded log levels for better visibility
- Auto-scrolling with toggle option
- Log download functionality
- Log clearing capability
- File source tracking for each log entry

## Installation

```bash
go install github.com/taigrr/log-socket/v2@latest
```

## Quick Start

```go
package main

import (
	"net/http"
	logger "github.com/taigrr/log-socket/v2/log"
	"github.com/taigrr/log-socket/v2/ws"
	"github.com/taigrr/log-socket/v2/browser"
)

func main() {
	defer logger.Flush()
	
	// Set up HTTP handlers
	http.HandleFunc("/ws", ws.LogSocketHandler)
	http.HandleFunc("/api/namespaces", ws.NamespacesHandler)
	http.HandleFunc("/", browser.LogSocketViewHandler)
	
	// Use default namespace
	logger.Info("Application started")
	
	// Create namespaced loggers
	apiLogger := logger.NewLogger("api")
	dbLogger := logger.NewLogger("database")
	
	apiLogger.Info("API server ready")
	dbLogger.Debug("Database connected")
	
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Usage

### Starting the Server

```bash
log-socket
```

By default, the server runs on `0.0.0.0:8080`. Specify a different address:

```bash
log-socket -addr localhost:8080
```

### Web Interface

Open your browser and navigate to `http://localhost:8080`

**Namespace Filtering:**
- The namespace dropdown is automatically populated from `/api/namespaces`
- Select one or more namespaces to filter (hold Ctrl/Cmd to multi-select)
- Default is "All Namespaces" (shows everything)
- Click "Reconnect" to apply the filter

## API

### Logging Interface

The package provides two ways to log:

#### 1. Package-level functions (default namespace)

```go
logger.Trace("trace message")
logger.Debug("debug message")
logger.Info("info message")
logger.Notice("notice message")
logger.Warn("warning message")
logger.Error("error message")
logger.Panic("panic message")  // Logs and panics
logger.Fatal("fatal message")  // Logs and exits
```

Each has formatted (`f`) and line (`ln`) variants:
```go
logger.Infof("User %s logged in", username)
logger.Infoln("This adds a newline")
```

#### 2. Namespaced loggers

```go
apiLogger := logger.NewLogger("api")
apiLogger.Info("Request received")

dbLogger := logger.NewLogger("database")  
dbLogger.Warn("Slow query detected")
```

### Creating Clients with Namespace Filters

```go
// Listen to all namespaces
client := logger.CreateClient()

// Listen to specific namespace
client := logger.CreateClient("api")

// Listen to multiple namespaces
client := logger.CreateClient("api", "database", "auth")
```

### WebSocket API

#### Log Stream Endpoint

**URL:** `ws://localhost:8080/ws`

**Query Parameters:**
- `namespaces` (optional): Comma-separated list of namespaces to filter

**Examples:**
```
ws://localhost:8080/ws                     # All namespaces
ws://localhost:8080/ws?namespaces=api      # Only "api" namespace
ws://localhost:8080/ws?namespaces=api,database  # Multiple namespaces
```

**Message Format:**
```json
{
  "timestamp": "2024-11-10T15:42:49.777298-05:00",
  "output": "API request received",
  "file": "main.go:42",
  "level": "INFO",
  "namespace": "api"
}
```

#### Namespaces List Endpoint

**URL:** `GET http://localhost:8080/api/namespaces`

**Response:**
```json
{
  "namespaces": ["default", "api", "database", "auth"]
}
```

## Web Interface Features

- **Namespace Dropdown**: Dynamically populated from `/api/namespaces`, multi-select support
- **Text Search**: Filter logs by content, level, namespace, or source file
- **Auto-scroll**: Toggle auto-scrolling with checkbox
- **Download**: Save all logs as a JSON file
- **Clear**: Remove all logs from the viewer
- **Color Coding**: Different log levels are color-coded
- **Reconnect**: Reconnect WebSocket with new namespace filter

## Terminal Colors

Log output to stderr is automatically colorized when writing to a terminal. Colors are disabled when output is piped or redirected to a file.

### Color Scheme

| Level  | Color        |
|--------|--------------|
| TRACE  | Gray         |
| DEBUG  | Cyan         |
| INFO   | Green        |
| NOTICE | Blue         |
| WARN   | Yellow       |
| ERROR  | Red          |
| PANIC  | Bold Red     |
| FATAL  | Bold Red     |

### Controlling Colors

```go
// Disable colors (e.g., for CI/CD or file output)
logger.SetColorEnabled(false)

// Enable colors explicitly
logger.SetColorEnabled(true)

// Check current state
if logger.ColorEnabled() {
    // colors are on
}
```

Colors are implemented using standard ANSI escape codes with no external dependencies.

## Migration from v1

### Import Path

```go
// v1
import "github.com/taigrr/log-socket/log"

// v2
import "github.com/taigrr/log-socket/v2/log"
```

### CreateClient Changes

```go
// v1
client := log.CreateClient()

// v2 - specify namespace(s) or leave empty for all
client := log.CreateClient()              // All namespaces
client := log.CreateClient("api")         // Single namespace
client := log.CreateClient("api", "db")   // Multiple namespaces
```

### New Logger Constructor

```go
// v2 only - create namespaced logger
apiLogger := log.NewLogger("api")
apiLogger.Info("Message in api namespace")
```

## Dependencies

- [gorilla/websocket](https://github.com/gorilla/websocket) for WebSocket support

## Notes

The web interface is provided as an example implementation. Users are encouraged to customize it for their specific needs. The WebSocket endpoint (`/ws`) can be consumed by any WebSocket client.

## License

See LICENSE file for details.
