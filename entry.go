package logrus

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Package-Level context to avoid repeated allocations
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

	// Logger provides access to the logger instance (logrus compatibility)
	Logger *Logger
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
		Logger:  logger,
	}
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value any) *Entry {
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
		Logger:  entry.logger,
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
		Logger:  entry.logger,
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
func (entry *Entry) log(level Level, args ...any) {
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
func (entry *Entry) logf(level Level, format string, args ...any) {
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
func (entry *Entry) logln(level Level, args ...any) {
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

// Trace logs a message at trace Level.
func (entry *Entry) Trace(args ...any) {
	entry.log(TraceLevel, args...)
}

// Debug logs a message at debug Level.
func (entry *Entry) Debug(args ...any) {
	entry.log(DebugLevel, args...)
}

// Info logs a message at info Level.
func (entry *Entry) Info(args ...any) {
	entry.log(InfoLevel, args...)
}

// Print logs a message at info Level (alias for Info).
func (entry *Entry) Print(args ...any) {
	entry.Info(args...)
}

// Warn logs a message at warning Level.
func (entry *Entry) Warn(args ...any) {
	entry.log(WarnLevel, args...)
}

// Warning logs a message at warning Level (alias for Warn).
func (entry *Entry) Warning(args ...any) {
	entry.Warn(args...)
}

// Error logs a message at error Level.
func (entry *Entry) Error(args ...any) {
	entry.log(ErrorLevel, args...)
}

// Fatal logs a message at fatal Level and exits the program.
func (entry *Entry) Fatal(args ...any) {
	entry.log(FatalLevel, args...)
}

// Panic logs a message at panic Level and panics.
func (entry *Entry) Panic(args ...any) {
	entry.log(PanicLevel, args...)
}

// Formatted logging methods

// Tracef logs a formatted message at trace Level.
func (entry *Entry) Tracef(format string, args ...any) {
	entry.logf(TraceLevel, format, args...)
}

// Debugf logs a formatted message at debug Level.
func (entry *Entry) Debugf(format string, args ...any) {
	entry.logf(DebugLevel, format, args...)
}

// Infof logs a formatted message at info Level.
func (entry *Entry) Infof(format string, args ...any) {
	entry.logf(InfoLevel, format, args...)
}

// Printf logs a formatted message at info Level (alias for Infof).
func (entry *Entry) Printf(format string, args ...any) {
	entry.Infof(format, args...)
}

// Warnf logs a formatted message at warning Level.
func (entry *Entry) Warnf(format string, args ...any) {
	entry.logf(WarnLevel, format, args...)
}

// Warningf logs a formatted message at warning Level (alias for Warnf).
func (entry *Entry) Warningf(format string, args ...any) {
	entry.Warnf(format, args...)
}

// Errorf logs a formatted message at error Level.
func (entry *Entry) Errorf(format string, args ...any) {
	entry.logf(ErrorLevel, format, args...)
}

// Fatalf logs a formatted message at fatal Level and exits the program.
func (entry *Entry) Fatalf(format string, args ...any) {
	entry.logf(FatalLevel, format, args...)
}

// Panicf logs a formatted message at panic Level and panics.
func (entry *Entry) Panicf(format string, args ...any) {
	entry.logf(PanicLevel, format, args...)
}

// Line logging methods

// Traceln logs a message at trace Level with newline handling.
func (entry *Entry) Traceln(args ...any) {
	entry.logln(TraceLevel, args...)
}

// Debugln logs a message at debug Level with newline handling.
func (entry *Entry) Debugln(args ...any) {
	entry.logln(DebugLevel, args...)
}

// Infoln logs a message at info Level with newline handling.
func (entry *Entry) Infoln(args ...any) {
	entry.logln(InfoLevel, args...)
}

// Println logs a message at info Level with newline handling (alias for Infoln).
func (entry *Entry) Println(args ...any) {
	entry.Infoln(args...)
}

// Warnln logs a message at warning Level with newline handling.
func (entry *Entry) Warnln(args ...any) {
	entry.logln(WarnLevel, args...)
}

// Warningln logs a message at warning Level with newline handling (alias for Warnln).
func (entry *Entry) Warningln(args ...any) {
	entry.Warnln(args...)
}

// Errorln logs a message at error Level with newline handling.
func (entry *Entry) Errorln(args ...any) {
	entry.logln(ErrorLevel, args...)
}

// Fatalln logs a message at fatal Level with newline handling and exits the program.
func (entry *Entry) Fatalln(args ...any) {
	entry.logln(FatalLevel, args...)
}

// Panicln logs a message at panic Level with newline handling and panics.
func (entry *Entry) Panicln(args ...any) {
	entry.logln(PanicLevel, args...)
}

// Writer returns an io.Writer that writes to the logger at the info log Level.
func (entry *Entry) Writer() *io.PipeWriter {
	return entry.WriterLevel(InfoLevel)
}

// WriterLevel returns an io.Writer that writes to the logger at the given log Level.
func (entry *Entry) WriterLevel(level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...any)

	switch level {
	case TraceLevel:
		printFunc = entry.Trace
	case DebugLevel:
		printFunc = entry.Debug
	case InfoLevel:
		printFunc = entry.Info
	case WarnLevel:
		printFunc = entry.Warn
	case ErrorLevel:
		printFunc = entry.Error
	case FatalLevel:
		printFunc = entry.Fatal
	case PanicLevel:
		printFunc = entry.Panic
	default:
		printFunc = entry.Print
	}

	go entry.writerScanner(reader, printFunc)

	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

// writerScanner scans the input from the reader and writes it to the logger.
func (entry *Entry) writerScanner(reader *io.PipeReader, printFunc func(args ...any)) {
	scanner := bufio.NewScanner(reader)

	// Use a reasonable buffer size for scanning
	scanner.Buffer(make([]byte, bufio.MaxScanTokenSize), bufio.MaxScanTokenSize)

	for scanner.Scan() {
		printFunc(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		entry.Error("Error while reading from Writer: ", err)
	}

	reader.Close()
}

// writerFinalizer is called when the writer is garbage collected.
func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}
