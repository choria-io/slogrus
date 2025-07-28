# Code Style

 - The project uses the normal Go code style
 - Files are formatted using "go fmt"
 - Tests are written using standard go test Package
 - We commend code inside functions only when the code being commented is not obviously doing what is being commented
 - This project is used at times where allocation and performance really matter, thus ensure code is optimised for those properties
 - This project prefer to use `any` rather than `interface{}`

# Workflow

 - Be sure to use `go fmt` to format files
 - Ensure tests are passing regularly using `go test`
 - Ensure code quality using `staticcheck`
 - Check code coverage and ensure a good coverage is achieved

# Writing Tests

When writing tests focus on correctness of the main project code. If it appears there is a bug in the main code base
rather stop and report that there is a bug than trying to make the tests work despite the bug.

# Memory and Performance Architecture

## Overview

This document captures the significant architectural decisions and performance optimizations implemented in slogrus to achieve high-performance logging while maintaining full logrus API compatibility.

## Core Architecture

### Logger Structure
- `Logger` wraps `*slog.Logger` with additional fields for level management and output control
- Maintains internal level tracking alongside slog's level system for fast level checks
- Uses composition pattern to leverage slog's optimized implementation

### Entry System
- `Entry` represents individual log entries with fields, context, and metadata
- Implements lazy evaluation - entries are only processed when actually logged
- Uses copy-on-write semantics for field management

## Performance Optimizations

### 1. Direct Logger Call Optimization

**Problem**: Original implementation always created Entry objects even for simple logging calls without fields.

**Solution**: Implemented fast-path direct logging in `logger.go`:
- `log()`, `logf()`, and `logln()` methods bypass Entry allocation entirely
- Direct calls to `slog.Logger.Log()` using package-level context
- Only create Entry objects when fields are needed

**Impact**: 60-80% reduction in allocations for direct logging calls.

### 2. Package-Level Context Optimization

**Problem**: Repeated `context.Background()` calls created unnecessary allocations.

**Solution**:
- Added `var backgroundContext = context.Background()` in `entry.go:12`
- Reuse single context instance across all logging calls
- Used in both direct logger calls and Entry logging

**Impact**: Eliminates context allocation overhead in every log call.

### 3. WithField Double Allocation Fix

**Problem**: `WithField` was creating temporary maps that were immediately copied.

**Solution**: Direct map creation with pre-allocated capacity in `entry.go:44-58`:
```go
data := make(Fields, len(entry.Data)+1)
for k, v := range entry.Data {
    data[k] = v
}
data[key] = value
```

**Impact**: 40-60% reduction in allocations for field-based logging.

### 4. Fast/Slow Path Implementation

**Problem**: Attribute slice allocation occurred even when no fields were present.

**Solution**: Conditional processing in Entry logging methods:
```go
if len(entry.Data) == 0 {
    // Fast path - no attributes
    entry.logger.slogger.Log(entry.Context, level.toSlogLevel(), msg)
} else {
    // Slow path - with attributes
    attrs := make([]slog.Attr, 0, len(entry.Data))
    // ... attribute processing
}
```

**Impact**: Zero allocations for entries without fields.

### 5. Level Check Optimization

**Problem**: Level checking should be extremely fast for disabled levels.

**Solution**:
- Early return in all logging methods with `IsLevelEnabled()` check
- Leverages slog's optimized level checking internally
- Disabled level calls result in ~2ns/op with 0 allocations

## Memory Layout Considerations

### Entry Structure
- Fields map pre-allocated with capacity 6 (common case optimization)
- Time stamp captured once at entry creation
- Context reference (not copy) to avoid allocation

### Fields Management
- `Fields` type alias to `map[string]any` for logrus compatibility
- Copy-on-write semantics for field operations
- Efficient iteration using range loops over maps

## Benchmark Results

Key performance characteristics achieved:

- **Direct logging**: ~500ns/op, 16B/op, 1 alloc/op
- **Disabled level checks**: ~2ns/op, 0B/op, 0 allocs/op
- **Single field logging**: ~850ns/op, 528B/op, 6 allocs/op
- **Complex logging (7 fields)**: ~1600ns/op, 881B/op, 7 allocs/op

## Critical Performance Paths

### Hot Path: Direct Logger Calls
1. Level check (inline, ~2ns if disabled)
2. Message formatting (`fmt.Sprint` variants)
3. Direct slog call with pre-allocated context
4. Special handling for Fatal/Panic levels

### Warm Path: Single Field Logging
1. Entry creation with field
2. Level check
3. Fast/slow path determination
4. Attribute conversion and slog call

### Cold Path: Multi-Field Logging
1. Multiple field operations (copy-on-write)
2. Entry processing with full attribute conversion
3. slog structured logging call

## Design Principles

1. **Zero-cost abstractions**: When features aren't used, they don't impact performance
2. **Allocation minimization**: Pre-allocate, reuse, and avoid temporary objects
3. **Level-aware processing**: Expensive operations only when logging will occur
4. **slog delegation**: Leverage Go's optimized structured logging implementation
5. **API compatibility**: No breaking changes to logrus interface

## Interoperability Features

### slog Integration APIs

**FromSlogLogger()**: Creates slogrus Logger from existing slog.Logger
- Enables wrapping existing slog-based infrastructure
- Maintains original handler configuration and performance characteristics
- Useful for gradual migration scenarios

**GetSlogLogger()**: Exposes underlying slog.Logger for direct access
- Allows mixed usage patterns (logrus-style + slog-style in same codebase)
- Enables access to advanced slog features (groups, custom handlers, etc.)
- Maintains zero-cost when not used

### Migration Patterns Enabled

1. **Wrap existing slog infrastructure**: `FromSlogLogger(existingSlog)`
2. **Mixed API usage**: Use logrus patterns for compatibility, slog for advanced features
3. **Gradual transition**: Teams can migrate method-by-method while sharing same underlying logger

## Future Considerations

- Pool Entry objects for high-throughput scenarios
- Custom attribute types for common field patterns
- Specialized handlers for specific use cases
- Memory profiling integration for optimization validation
- Enhanced slog feature exposure (groups, custom LogValue types)