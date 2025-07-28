package slogrus

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// Package-level context to avoid repeated allocations
var backgroundContext = context.Background()

// Entry represents a single log entry, compatible with logrus.Entry.
type Entry struct {
	logger *Logger
	Data   Fields
	Time   time.Time
	Level  Level
	Caller *Caller

	// Context holds the context associated with this entry
	Context context.Context
}

// Caller represents caller information for a log entry.
type Caller struct {
	File     string
	Line     int
	Function string
}

// NewEntry creates a new Entry instance.
func NewEntry(logger *Logger) *Entry {
	return &Entry{
		logger:  logger,
		Data:    make(Fields, 6),
		Time:    time.Now(),
		Context: backgroundContext,
	}
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	data := make(Fields, len(entry.Data)+1)
	for k, v := range entry.Data {
		data[k] = v
	}
	data[key] = value
	return &Entry{
		logger:  entry.logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Caller:  entry.Caller,
		Context: entry.Context,
	}
}

// WithFields adds multiple fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data)+len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{
		logger:  entry.logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Caller:  entry.Caller,
		Context: entry.Context,
	}
}

// WithContext adds a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{
		logger:  entry.logger,
		Data:    dataCopy,
		Time:    entry.Time,
		Level:   entry.Level,
		Caller:  entry.Caller,
		Context: ctx,
	}
}

// WithError adds an error field to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return entry.WithField("error", err)
}

// WithTime adds a time field to the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	dataCopy := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		dataCopy[k] = v
	}
	return &Entry{
		logger:  entry.logger,
		Data:    dataCopy,
		Time:    t,
		Level:   entry.Level,
		Caller:  entry.Caller,
		Context: entry.Context,
	}
}

// log is the internal logging method that writes to slog
func (entry *Entry) log(level Level, args ...interface{}) {
	if !entry.logger.IsLevelEnabled(level) {
		return
	}

	// Get message
	msg := fmt.Sprint(args...)

	if len(entry.Data) == 0 {
		// Fast path - no attributes
		entry.logger.slogger.Log(entry.Context, level.toSlogLevel(), msg)
	} else {
		// Slow path - with attributes
		attrs := make([]slog.Attr, 0, len(entry.Data))
		for k, v := range entry.Data {
			attrs = append(attrs, slog.Any(k, v))
		}
		entry.logger.slogger.LogAttrs(entry.Context, level.toSlogLevel(), msg, attrs...)
	}

	// Handle Fatal and Panic levels
	if level == FatalLevel {
		os.Exit(1)
	} else if level == PanicLevel {
		panic(msg)
	}
}

// logf is the internal formatted logging method
func (entry *Entry) logf(level Level, format string, args ...interface{}) {
	if !entry.logger.IsLevelEnabled(level) {
		return
	}

	// Format message
	msg := fmt.Sprintf(format, args...)

	if len(entry.Data) == 0 {
		// Fast path - no attributes
		entry.logger.slogger.Log(entry.Context, level.toSlogLevel(), msg)
	} else {
		// Slow path - with attributes
		attrs := make([]slog.Attr, 0, len(entry.Data))
		for k, v := range entry.Data {
			attrs = append(attrs, slog.Any(k, v))
		}
		entry.logger.slogger.LogAttrs(entry.Context, level.toSlogLevel(), msg, attrs...)
	}

	// Handle Fatal and Panic levels
	if level == FatalLevel {
		os.Exit(1)
	} else if level == PanicLevel {
		panic(msg)
	}
}

// logln is the internal line logging method
func (entry *Entry) logln(level Level, args ...interface{}) {
	if !entry.logger.IsLevelEnabled(level) {
		return
	}

	// Get message - sprintln adds spaces between args and a newline
	msg := fmt.Sprintln(args...)
	// Remove trailing newline since slog will add its own formatting
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}

	if len(entry.Data) == 0 {
		// Fast path - no attributes
		entry.logger.slogger.Log(entry.Context, level.toSlogLevel(), msg)
	} else {
		// Slow path - with attributes
		attrs := make([]slog.Attr, 0, len(entry.Data))
		for k, v := range entry.Data {
			attrs = append(attrs, slog.Any(k, v))
		}
		entry.logger.slogger.LogAttrs(entry.Context, level.toSlogLevel(), msg, attrs...)
	}

	// Handle Fatal and Panic levels
	if level == FatalLevel {
		os.Exit(1)
	} else if level == PanicLevel {
		panic(msg)
	}
}

// Trace logs a message at trace level.
func (entry *Entry) Trace(args ...interface{}) {
	entry.log(TraceLevel, args...)
}

// Debug logs a message at debug level.
func (entry *Entry) Debug(args ...interface{}) {
	entry.log(DebugLevel, args...)
}

// Info logs a message at info level.
func (entry *Entry) Info(args ...interface{}) {
	entry.log(InfoLevel, args...)
}

// Print logs a message at info level (alias for Info).
func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

// Warn logs a message at warning level.
func (entry *Entry) Warn(args ...interface{}) {
	entry.log(WarnLevel, args...)
}

// Warning logs a message at warning level (alias for Warn).
func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

// Error logs a message at error level.
func (entry *Entry) Error(args ...interface{}) {
	entry.log(ErrorLevel, args...)
}

// Fatal logs a message at fatal level and exits the program.
func (entry *Entry) Fatal(args ...interface{}) {
	entry.log(FatalLevel, args...)
}

// Panic logs a message at panic level and panics.
func (entry *Entry) Panic(args ...interface{}) {
	entry.log(PanicLevel, args...)
}

// Formatted logging methods

// Tracef logs a formatted message at trace level.
func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.logf(TraceLevel, format, args...)
}

// Debugf logs a formatted message at debug level.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.logf(DebugLevel, format, args...)
}

// Infof logs a formatted message at info level.
func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.logf(InfoLevel, format, args...)
}

// Printf logs a formatted message at info level (alias for Infof).
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

// Warnf logs a formatted message at warning level.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.logf(WarnLevel, format, args...)
}

// Warningf logs a formatted message at warning level (alias for Warnf).
func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

// Errorf logs a formatted message at error level.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.logf(ErrorLevel, format, args...)
}

// Fatalf logs a formatted message at fatal level and exits the program.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.logf(FatalLevel, format, args...)
}

// Panicf logs a formatted message at panic level and panics.
func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.logf(PanicLevel, format, args...)
}

// Line logging methods

// Traceln logs a message at trace level with newline handling.
func (entry *Entry) Traceln(args ...interface{}) {
	entry.logln(TraceLevel, args...)
}

// Debugln logs a message at debug level with newline handling.
func (entry *Entry) Debugln(args ...interface{}) {
	entry.logln(DebugLevel, args...)
}

// Infoln logs a message at info level with newline handling.
func (entry *Entry) Infoln(args ...interface{}) {
	entry.logln(InfoLevel, args...)
}

// Println logs a message at info level with newline handling (alias for Infoln).
func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

// Warnln logs a message at warning level with newline handling.
func (entry *Entry) Warnln(args ...interface{}) {
	entry.logln(WarnLevel, args...)
}

// Warningln logs a message at warning level with newline handling (alias for Warnln).
func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

// Errorln logs a message at error level with newline handling.
func (entry *Entry) Errorln(args ...interface{}) {
	entry.logln(ErrorLevel, args...)
}

// Fatalln logs a message at fatal level with newline handling and exits the program.
func (entry *Entry) Fatalln(args ...interface{}) {
	entry.logln(FatalLevel, args...)
}

// Panicln logs a message at panic level with newline handling and panics.
func (entry *Entry) Panicln(args ...interface{}) {
	entry.logln(PanicLevel, args...)
}
