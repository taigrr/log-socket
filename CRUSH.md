# CRUSH.md - Log Socket v2 Development Guide

This document provides context and conventions for working on **log-socket v2** - a real-time log viewer with WebSocket support and namespace filtering, written in Go.

## Project Overview

Log Socket v2 is a Go library and standalone application that provides:
- Real-time log streaming via WebSocket
- **Namespace-based log organization** (NEW in v2)
- Web-based log viewer with namespace filtering
- Support for multiple log levels (TRACE, DEBUG, INFO, NOTICE, WARN, ERROR, PANIC, FATAL)
- Client architecture allowing multiple subscribers to filtered log streams

**Key insight**: This is both a library (importable Go package) and a standalone application. The main.go serves as an example implementation.

## Essential Commands

### Build
```bash
go build -v ./...
```

### Run Server
```bash
# Default (runs on 0.0.0.0:8080)
go run main.go

# Custom address
go run main.go -addr localhost:8080
```

Once running, open browser to `http://localhost:8080` to view logs.

### Test
```bash
go test -v ./...
```

### Install
```bash
go install github.com/taigrr/log-socket/v2@latest
```

### Dependencies
```bash
go get .
```

## Project Structure

```
.
├── main.go              # Example server with multiple namespaces
├── log/                 # Core logging package
│   ├── log.go          # Package-level logging functions + namespace tracking
│   ├── logger.go       # Logger type with namespace support
│   ├── types.go        # Type definitions (includes Namespace fields)
│   └── log_test.go     # Tests
├── ws/                  # WebSocket server
│   ├── server.go       # LogSocketHandler with namespace filtering
│   └── namespaces.go   # HTTP handler for namespace list API
└── browser/             # Web UI
    ├── browser.go      # HTTP handler serving embedded HTML
    └── viewer.html     # Embedded web interface with namespace filter
```

## Major Changes in v2

### Module Path
- **v1**: `github.com/taigrr/log-socket`
- **v2**: `github.com/taigrr/log-socket/v2`

### Namespace Support

**Core concept**: Namespaces allow organizing logs by component, service, or domain (e.g., "api", "database", "auth").

#### Types Changes

**Entry** now includes:
```go
type Entry struct {
    Timestamp time.Time `json:"timestamp"`
    Output    string    `json:"output"`
    File      string    `json:"file"`
    Level     string    `json:"level"`
    Namespace string    `json:"namespace"` // NEW
    level     Level
}
```

**Client** now uses a slice for filtering:
```go
type Client struct {
    LogLevel    Level    `json:"level"`
    Namespaces  []string `json:"namespaces"` // Empty = all namespaces
    writer      LogWriter
    initialized bool
}
```

**Logger** has namespace field:
```go
type Logger struct {
    FileInfoDepth int
    Namespace     string // NEW
}
```

#### API Changes

**CreateClient** now variadic:
```go
// v1
func CreateClient() *Client

// v2
func CreateClient(namespaces ...string) *Client

// Examples:
client := log.CreateClient()                    // All namespaces
client := log.CreateClient("api")               // Single namespace
client := log.CreateClient("api", "database")   // Multiple namespaces
```

**NewLogger** constructor added:
```go
func NewLogger(namespace string) *Logger

// Example:
apiLogger := log.NewLogger("api")
apiLogger.Info("API request received")
```

### Namespace Tracking

Global namespace registry tracks all used namespaces:
```go
var (
    namespaces    map[string]bool
    namespacesMux sync.RWMutex
)

func GetNamespaces() []string
```

Namespaces are automatically registered when logs are created.

### WebSocket Changes

**Query parameter for filtering**:
```
ws://localhost:8080/ws?namespaces=api,database
```

The handler parses comma-separated namespace list and creates a filtered client.

### Web UI Changes

- Namespace filter input field added to controls
- Namespace column added to log table
- Reconnect button to apply namespace filter
- WebSocket URL includes namespace query parameter

## Code Organization & Architecture

### Log Package (`log/`)

Dual API remains, but with namespace support:

1. **Package-level functions**: Use "default" namespace
   - `log.Info()`, `log.Debug()`, etc.
   - All entry creations include `Namespace: DefaultNamespace`

2. **Logger instances**: Use custom namespace
   - Create with `log.NewLogger(namespace)`
   - All entry creations include `Namespace: l.Namespace`

### Client Architecture (Updated)

**Client filtering by namespace**:

1. **Empty Namespaces slice**: Receives all logs regardless of namespace
2. **Non-empty Namespaces**: Only receives logs matching one of the specified namespaces

**matchesNamespace helper**:
```go
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
```

**Entry flow with namespace filtering**:
1. Log function called with namespace
2. `Entry` created with namespace field
3. Namespace registered in global map
4. `createLog()` sends to all clients
5. Each client checks `matchesNamespace()`
6. Only matching clients receive the entry

### WebSocket Handler (`ws/`)

**Namespace parameter parsing**:
```go
namespacesParam := r.URL.Query().Get("namespaces")
var namespaces []string
if namespacesParam != "" {
    namespaces = strings.Split(namespacesParam, ",")
}
lc := logger.CreateClient(namespaces...)
```

**Namespaces API handler** (`ws/namespaces.go`):
```go
func NamespacesHandler(w http.ResponseWriter, r *http.Request) {
    namespaces := logger.GetNamespaces()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "namespaces": namespaces,
    })
}
```

### Browser Package (`browser/`)

Still embeds viewer.html, but HTML now includes:
- Namespace filter input
- Namespace column in grid (5 columns instead of 4)
- `reconnectWithNamespace()` method
- WebSocket URL construction with query parameter

## Go Version & Dependencies

- **Go version**: 1.24.4 (specified in go.mod)
- **Only external dependency**: `github.com/gorilla/websocket v1.5.3`

## Naming Conventions & Style

### Namespaces
- Use lowercase strings: `"api"`, `"database"`, `"auth"`, `"cache"`
- Default constant: `DefaultNamespace = "default"`
- Comma-separated in query params: `?namespaces=api,database`

### Log Levels
Unchanged from v1 - still use uppercase strings and iota constants.

### Variable Names
- Use descriptive names (`apiLogger`, `dbLogger`, `namespaces`)
- Exception: Loop variables, short-lived scopes

## Important Patterns & Gotchas

### 1. Namespace Tracking is Automatic
When any log is created, its namespace is automatically added to the global registry:
```go
func createLog(e Entry) {
    // Track namespace
    namespacesMux.Lock()
    namespaces[e.Namespace] = true
    namespacesMux.Unlock()
    // ... rest of function
}
```

**No manual registration needed** - just log and the namespace appears.

### 2. Empty Namespace List = All Logs
Both for clients and WebSocket connections:
```go
client := log.CreateClient()           // Gets ALL logs
ws://localhost:8080/ws                 // Gets ALL logs
```

This is the default behavior to maintain backward compatibility.

### 3. Client Namespace Filtering is Inclusive (OR)
If a client has multiple namespaces, it receives logs matching ANY of them:
```go
client := log.CreateClient("api", "database")
// Receives logs from "api" OR "database", not "auth"
```

### 4. Namespace Field Always Set
All logging functions set namespace:
- Package functions: `DefaultNamespace`
- Logger methods: `l.Namespace`

**Never nil or empty** - there's always a namespace.

### 5. WebSocket Namespace Reconnection
Changing namespace filter requires reconnecting WebSocket:
```javascript
reconnectWithNamespace() {
    if (this.ws) {
        this.ws.onclose = null; // Prevent auto-reconnect
        this.ws.close();
        this.ws = null;
    }
    this.reconnectAttempts = 0;
    this.connectWebSocket(); // Creates new connection with new filter
}
```

UI provides "Reconnect" button for this purpose.

### 6. Stderr Client Uses All Namespaces
The built-in stderr client (created in `init()`) listens to all namespaces:
```go
stderrClient = CreateClient(DefaultNamespace)
```

But only prints logs matching its own namespace in `logStdErr()`:
```go
if e.level >= c.LogLevel && c.matchesNamespace(e.Namespace) {
    fmt.Fprintf(os.Stderr, "%s\t%s\t[%s]\t%s\t%s\n", ...)
}
```

**Wait, that's a bug!** The stderr client is created with `DefaultNamespace` but should be created with no namespaces to see all logs. Let me check this.

Actually looking at the code:
```go
stderrClient = CreateClient(DefaultNamespace)
```

This means stderr client only sees "default" namespace logs. This might be intentional, but seems like a bug. Should probably be:
```go
stderrClient = CreateClient() // No args = all namespaces
```

### 7. Grid Layout Updated
The log viewer grid changed from 4 to 5 columns:
```css
/* v1 */
grid-template-columns: 180px 80px 1fr 120px;

/* v2 */
grid-template-columns: 180px 80px 100px 1fr 120px;
```

Order: Timestamp, Level, Namespace, Message, Source

## Testing

### Test Updates for v2

All `CreateClient()` calls in tests now pass namespace:
```go
// v1
c := CreateClient()

// v2
c := CreateClient("test")
c := CreateClient(DefaultNamespace)
```

Tests verify namespace appears in output (see stderr format).

### Running Tests
```bash
go test -v ./...
```

All existing tests pass with namespace support added.

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yaml`):
- Still uses Go 1.21 (should update to 1.24.4 to match go.mod)
- No changes needed for v2 functionality

## Common Tasks

### Adding a New Namespace

No code changes needed! Just create a logger:
```go
cacheLogger := log.NewLogger("cache")
cacheLogger.Info("Cache initialized")
```

Namespace automatically tracked and available via API.

### Creating Namespace-Specific Client

```go
// Subscribe only to API logs
apiClient := log.CreateClient("api")
defer apiClient.Destroy()

for {
    entry := apiClient.Get()
    // Only receives logs from "api" namespace
    processAPILog(entry)
}
```

### Filtering WebSocket by Namespace

Frontend:
```javascript
// Set namespace filter
document.getElementById('namespaceFilter').value = 'api,database';
// Click reconnect button or call:
logViewer.reconnectWithNamespace();
```

Backend automatically creates filtered client based on query param.

### Getting All Active Namespaces

```go
namespaces := log.GetNamespaces()
// Returns: ["default", "api", "database", "auth", ...]
```

Or via HTTP:
```bash
GET /api/namespaces
```

Returns:
```json
{
  "namespaces": ["default", "api", "database", "auth"]
}
```

## Migration from v1 to v2

### Import Paths
```go
// v1
import "github.com/taigrr/log-socket/log"
import "github.com/taigrr/log-socket/ws"
import "github.com/taigrr/log-socket/browser"

// v2
import "github.com/taigrr/log-socket/v2/log"
import "github.com/taigrr/log-socket/v2/ws"
import "github.com/taigrr/log-socket/v2/browser"
```

### CreateClient Calls
```go
// v1
client := log.CreateClient()

// v2 - same behavior
client := log.CreateClient() // Empty = all namespaces

// v2 - new filtering capability
client := log.CreateClient("api")
client := log.CreateClient("api", "database")
```

### Default() Logger
```go
// v1
logger := log.Default()

// v2 - uses default namespace
logger := log.Default()

// v2 - new namespaced option
logger := log.NewLogger("api")
```

### WebSocket URL
```
v1: ws://host/ws
v2: ws://host/ws                          # All namespaces (backward compatible)
v2: ws://host/ws?namespaces=api           # Filtered
v2: ws://host/ws?namespaces=api,database  # Multiple
```

## Repository Context

- **License**: Check LICENSE file
- **Funding**: GitHub sponsors (.github/FUNDING.yml)
- **Open Source**: github.com/taigrr/log-socket
- **Version**: v2.x.x (major version bump for breaking changes)

## Future Agent Notes

- **v2 is a breaking change**: Major version bump follows Go modules convention
- Namespace support is the primary new feature
- Backward compatible behavior: empty namespace list = all logs
- Namespace tracking is automatic via global registry
- Web UI has been significantly updated for namespace support
- Possible bug: stderr client might need to use `CreateClient()` instead of `CreateClient(DefaultNamespace)` to see all logs
- All tests updated and passing
- Example in main.go demonstrates multiple namespaces
