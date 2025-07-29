# slogrus

A drop-in replacement for [logrus](https://github.com/sirupsen/logrus) that uses Go's standard `log/slog` package under the hood. This library provides a migration path for projects that want to move from logrus to Go's structured logging without changing their existing code or to facilitate a gradual migration.

> [!NOTE]
> Claude Code almost entirely created this project. It is an experiment to see how this might solve the migration problem. Using this I migrated the entire `go-choria` code base in less than a hour with no code changes.

## Goals

- **Zero code changes**: Replace your logrus import and your code continues to work
- **No external dependencies**: Uses only Go's standard library `log/slog`
- **Full API compatibility**: Supports all logrus types, methods, and patterns
- **Modern logging**: Leverages Go 1.21+'s structured logging capabilities
- **Performance**: Benefits from slog's optimized implementation

## Installation

```bash
go get github.com/choria-io/slogrus
```

## Usage

### Basic Migration

Replace your logrus import:

```go
// Before
import "github.com/sirupsen/logrus"

// After
import "github.com/choria-io/slogrus"
```

All your existing code continues to work:

```go
slogrus.Info("Application started")
slogrus.WithField("user", "john").Info("User logged in")
slogrus.WithFields(slogrus.Fields{
    "method": "GET",
    "path":   "/api/users",
    "status": 200,
}).Info("Request completed")
```

### Logger Instances

Create and configure custom loggers:

```go
logger := slogrus.New()
logger.SetLevel(slogrus.DebugLevel)
logger.SetOutput(os.Stdout)

logger.Debug("Debug message")
logger.Info("Info message")
logger.Error("Error message")
```

### Structured Logging

Use fields for structured logging:

```go
log := slogrus.WithFields(slogrus.Fields{
    "component": "database",
    "operation": "query",
    "table":     "users",
})

log.Info("Query executed successfully")
log.WithField("duration", "45ms").Info("Query performance")
```

### Method Chaining

Chain methods for building complex log entries:

```go
slogrus.WithField("service", "api").
    WithField("version", "1.2.3").
    WithError(err).
    Error("Service startup failed")
```

### Context Support

Use context for request tracing:

```go
ctx := context.WithValue(context.Background(), "requestID", "req-123")
slogrus.WithContext(ctx).Info("Processing request")
```

### Logging Levels

All logrus levels are supported:

```go
slogrus.Trace("Very detailed information")
slogrus.Debug("Debug information") 
slogrus.Info("General information")
slogrus.Warn("Warning message")
slogrus.Error("Error occurred")
slogrus.Fatal("Fatal error - will exit")  // Calls os.Exit(1)
slogrus.Panic("Panic error - will panic") // Calls panic()
```

### Formatted Logging

Support for formatted strings:

```go
slogrus.Infof("User %s logged in", username)
slogrus.Errorf("Failed to connect to %s:%d", host, port)

// With fields
slogrus.WithField("component", "auth").
    Warnf("Failed login attempt for user %s", username)
```

## Advanced Usage

### Custom Handlers

Use slog handlers for advanced output formatting:

```go
// JSON output
logger := slogrus.NewJSONLogger(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})

// Text output with options
logger := slogrus.NewTextLogger(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelInfo,
    AddSource: true, // Include source file information
})

// Custom handler
handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
    Level: slog.LevelWarn,
})
logger := slogrus.NewWithHandler(handler)
```

### slog Interoperability

For mixed usage scenarios and gradual migration:

```go
// Create from existing slog logger
existingSlog := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger := slogrus.FromSlogLogger(existingSlog)

// Use logrus-style API
logger.WithField("component", "auth").Info("User logged in")

// Access underlying slog logger for advanced features
slogger := logger.GetSlogLogger()
slogger.Info("Direct slog usage", 
    "user", "john",
    "action", "login",
    slog.Group("metadata",
        "ip", "192.168.1.1",
        "agent", "curl/7.68.0",
    ),
)
```

### Formatter Compatibility

Basic formatter compatibility for easier migration:

```go
// Switch to JSON format
slogrus.SetFormatter(&slogrus.JSONFormatter{})

// Switch to text format  
slogrus.SetFormatter(&slogrus.TextFormatter{
    DisableColors: true,
    FullTimestamp: true,
})
```

### Level Management

```go
// Set global level
slogrus.SetLevel(slogrus.WarnLevel)

// Check if level is enabled
if slogrus.StandardLogger().IsLevelEnabled(slogrus.DebugLevel) {
    // Expensive debug operation
    slogrus.Debug("Detailed debug info")
}

// Parse level from string
level, err := slogrus.ParseLevel("info")
if err != nil {
    log.Fatal(err)
}
slogrus.SetLevel(level)
```

## Migration Guide

### Step 1: Update Import
```go
// Replace this
import "github.com/sirupsen/logrus"

// With this
import "github.com/choria-io/slogrus"
```

### Step 2: Test Your Application
Run your tests to ensure everything works as expected. The API is fully compatible.

### Step 3: Optimize (Optional)
Consider leveraging slog-specific features:

```go
// Take advantage of slog's structured logging
logger := slogrus.NewJSONLogger(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelInfo,
    AddSource: true,
})
```

## API Compatibility

This library implements the complete logrus API:

- **Types**: `Logger`, `Entry`, `Level`, `Fields`
- **Levels**: `PanicLevel`, `FatalLevel`, `ErrorLevel`, `WarnLevel`, `InfoLevel`, `DebugLevel`, `TraceLevel`
- **Methods**: All logging methods (`Info`, `Debug`, `Error`, etc.)
- **Formatted methods**: `Infof`, `Debugf`, `Errorf`, etc.
- **Line methods**: `Infoln`, `Debugln`, `Errorln`, etc.
- **Entry methods**: `WithField`, `WithFields`, `WithError`, `WithContext`, `WithTime`
- **Global functions**: All package-level logging functions
- **Configuration**: `SetLevel`, `SetOutput`, `SetFormatter`, `SetReportCaller`

## Performance

By using Go's standard `log/slog`, this library benefits from:

- Optimized structured logging implementation
- Efficient level checking
- Minimal allocations for disabled log levels
- Native JSON encoding

## Differences from Logrus

While maintaining API compatibility, there are some behavioral differences:

1. **Formatters**: Formatter interfaces are implemented but actual formatting is handled by slog handlers
2. **Hooks**: Hook system is not implemented (use slog handlers instead)
3. **Output types**: Only `io.Writer` outputs are supported (not syslog, etc.)

## Requirements

- Go 1.23 or later (for `log/slog` support)

## Examples

See the test files for comprehensive usage examples:
- `slogrus_test.go` - Basic functionality tests
- `logger_test.go` - Logger-specific tests  
- `entry_test.go` - Entry and structured logging tests
- `init_test.go` - Initialization and global function tests

## Benchmarks

This was written with a focus on performance and low allocations, extensive benchmarks are included:

| Benchmark                        | Iterations    | Time per Op  | Bytes per Op         | Allocs per Op |
|----------------------------------|---------------|--------------|----------------------|---------------|
| LoggerInfo                       | 1,907,768     | 617.4 ns/op  | 16 B/op              | 1 allocs/op   |
| LoggerInfof                      | 1,743,793     | 692.6 ns/op  | 32 B/op              | 1 allocs/op   |
| LoggerInfoln                     | 1,906,810     | 639.4 ns/op  | 16 B/op              | 1 allocs/op   |
| LoggerWithField                  | 1,000,000     | 1,014 ns/op  | 528 B/op             | 6 allocs/op   |
| LoggerWithFieldf                 | 1,000,000     | 1,090 ns/op  | 544 B/op             | 7 allocs/op   |
| LoggerWithFields                 | 930,514       | 1,269 ns/op  | 608 B/op             | 6 allocs/op   |
| LoggerWithFieldChaining          | 772,938       | 1,551 ns/op  | 1,289 B/op           | 11 allocs/op  |
| LoggerDebugDisabled              | 494,086,148   | 2.428 ns/op  | 0 B/op               | 0 allocs/op   |
| LoggerWithFieldDebugDisabled     | 6,154,128     | 202.3 ns/op  | 464 B/op             | 4 allocs/op   |
| GlobalInfo                       | 1,899,432     | 633.5 ns/op  | 16 B/op              | 1 allocs/op   |
| GlobalWithField                  | 1,000,000     | 1,073 ns/op  | 528 B/op             | 6 allocs/op   |
| LoggerInfoJSON                   | 2,346,199     | 508.7 ns/op  | 16 B/op              | 1 allocs/op   |
| LoggerWithFieldJSON              | 1,336,527     | 888.2 ns/op  | 528 B/op             | 6 allocs/op   |
| LoggerWithError                  | 931,095       | 1,225 ns/op  | 544 B/op             | 7 allocs/op   |
| ComplexLogging                   | 586,546       | 1,962 ns/op  | 881 B/op             | 7 allocs/op   |
| MemoryAllocation/DirectLog       | 1,874,785     | 637.9 ns/op  | 16 B/op              | 1 allocs/op   |
| MemoryAllocation/WithOneField    | 1,000,000     | 1,072 ns/op  | 536 B/op             | 6 allocs/op   |
| MemoryAllocation/WithThreeFields | 750,601       | 1,572 ns/op  | 1,297 B/op           | 10 allocs/op  |
| MemoryAllocation/WithFieldsMap   | 934,754       | 1,342 ns/op  | 616 B/op             | 6 allocs/op   |
| Throughput                       | 1,874,892     | 607.6 ns/op  | 750 B/op (1.65 MB/s) | 6 allocs/op   |
| LevelCheck                       | 1,000,000,000 | 0.3733 ns/op | 0 B/op               | 0 allocs/op   |
| FromSlogLogger                   | 1,820,745     | 657.6 ns/op  | 24 B/op              | 1 allocs/op   |
| GetSlogLogger                    | 2,019,850     | 597.1 ns/op  | 0 B/op               | 0 allocs/op   |