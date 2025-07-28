package logrus

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// Logger is the main logging struct that wraps slog.Logger for logrus compatibility.
type Logger struct {
	slogger *slog.Logger
	Level   Level

	// Out provides access to the configured output writer (logrus compatibility)
	Out io.Writer

	// Formatter stores the configured handler type (logrus compatibility)
	Formatter Formatter
}

// New creates a new Logger instance with default text handler.
func New() *Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &Logger{
		slogger:   slog.New(handler),
		Level:     InfoLevel,
		Out:       os.Stderr,
		Formatter: &TextFormatter{},
	}
}

// NewWithHandler creates a new Logger with a custom slog.Handler.
func NewWithHandler(handler slog.Handler) *Logger {
	// Determine formatter type based on handler
	var formatter Formatter = &TextFormatter{}
	if _, ok := handler.(*slog.JSONHandler); ok {
		formatter = &JSONFormatter{}
	}
	return &Logger{
		slogger:   slog.New(handler),
		Level:     InfoLevel,
		Out:       os.Stderr,
		Formatter: formatter,
	}
}

// FromSlogLogger creates a new Logger instance from an existing slog.Logger.
// This enables interoperability with existing slog-based code.
func FromSlogLogger(slogger *slog.Logger) *Logger {
	// Determine formatter type based on handler
	var formatter Formatter = &TextFormatter{}
	if _, ok := slogger.Handler().(*slog.JSONHandler); ok {
		formatter = &JSONFormatter{}
	}
	return &Logger{
		slogger:   slogger,
		Level:     InfoLevel, // Default Level, can be changed with SetLevel
		Out:       os.Stderr, // Default output, may not match slog handler's output
		Formatter: formatter,
	}
}

// SetOutput sets the output destination for the logger.
func (logger *Logger) SetOutput(out io.Writer) {
	logger.Out = out
	// Create a new handler with the new output
	opts := &slog.HandlerOptions{
		Level: logger.Level.toSlogLevel(),
	}

	if _, ok := logger.slogger.Handler().(*slog.TextHandler); ok {
		logger.slogger = slog.New(slog.NewTextHandler(logger.Out, opts))
		logger.Formatter = &TextFormatter{}
	} else if _, ok := logger.slogger.Handler().(*slog.JSONHandler); ok {
		logger.slogger = slog.New(slog.NewJSONHandler(logger.Out, opts))
		logger.Formatter = &JSONFormatter{}
	}
}

// SetLevel sets the logging Level for the logger.
func (logger *Logger) SetLevel(level Level) {
	logger.Level = level
	// Update the slog handler with new Level
	opts := &slog.HandlerOptions{
		Level: level.toSlogLevel(),
	}

	// Recreate handler with new Level
	if _, ok := logger.slogger.Handler().(*slog.TextHandler); ok {
		logger.slogger = slog.New(slog.NewTextHandler(logger.Out, opts))
		logger.Formatter = &TextFormatter{}
	} else if _, ok := logger.slogger.Handler().(*slog.JSONHandler); ok {
		logger.slogger = slog.New(slog.NewJSONHandler(logger.Out, opts))
		logger.Formatter = &JSONFormatter{}
	}
}

// IsLevelEnabled checks if the given Level is enabled for logging.
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return level <= logger.Level
}

// GetSlogLogger returns the underlying slog.Logger instance.
// This enables advanced slog operations and direct access to slog APIs.
func (logger *Logger) GetSlogLogger() *slog.Logger {
	return logger.slogger
}

// WithField creates an entry with a single field.
func (logger *Logger) WithField(key string, value any) *Entry {
	entry := NewEntry(logger)
	return entry.WithField(key, value)
}

// WithFields creates an entry with multiple fields.
func (logger *Logger) WithFields(fields Fields) *Entry {
	entry := NewEntry(logger)
	return entry.WithFields(fields)
}

// WithContext creates an entry with a context.
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	entry := NewEntry(logger)
	return entry.WithContext(ctx)
}

// WithError creates an entry with an error field.
func (logger *Logger) WithError(err error) *Entry {
	entry := NewEntry(logger)
	return entry.WithError(err)
}

// Direct logging methods

// log is the internal logging method
func (logger *Logger) log(level Level, args ...any) {
	if !logger.IsLevelEnabled(level) {
		return
	}

	// Fast path - direct slog call without Entry allocation
	msg := fmt.Sprint(args...)
	logger.slogger.Log(backgroundContext, level.toSlogLevel(), msg)

	// Handle Fatal and Panic levels
	if level == FatalLevel {
		os.Exit(1)
	} else if level == PanicLevel {
		panic(msg)
	}
}

// logf is the internal formatted logging method
func (logger *Logger) logf(level Level, format string, args ...any) {
	if !logger.IsLevelEnabled(level) {
		return
	}

	// Fast path - direct slog call without Entry allocation
	msg := fmt.Sprintf(format, args...)
	logger.slogger.Log(backgroundContext, level.toSlogLevel(), msg)

	// Handle Fatal and Panic levels
	if level == FatalLevel {
		os.Exit(1)
	} else if level == PanicLevel {
		panic(msg)
	}
}

// logln is the internal line logging method
func (logger *Logger) logln(level Level, args ...any) {
	if !logger.IsLevelEnabled(level) {
		return
	}

	// Fast path - direct slog call without Entry allocation
	msg := fmt.Sprintln(args...)
	// Remove trailing newline since slog will add its own formatting
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	logger.slogger.Log(backgroundContext, level.toSlogLevel(), msg)

	// Handle Fatal and Panic levels
	if level == FatalLevel {
		os.Exit(1)
	} else if level == PanicLevel {
		panic(msg)
	}
}

// Trace logs a message at trace Level.
func (logger *Logger) Trace(args ...any) {
	logger.log(TraceLevel, args...)
}

// Debug logs a message at debug Level.
func (logger *Logger) Debug(args ...any) {
	logger.log(DebugLevel, args...)
}

// Info logs a message at info Level.
func (logger *Logger) Info(args ...any) {
	logger.log(InfoLevel, args...)
}

// Print logs a message at info Level (alias for Info).
func (logger *Logger) Print(args ...any) {
	logger.Info(args...)
}

// Warn logs a message at warning Level.
func (logger *Logger) Warn(args ...any) {
	logger.log(WarnLevel, args...)
}

// Warning logs a message at warning Level (alias for Warn).
func (logger *Logger) Warning(args ...any) {
	logger.Warn(args...)
}

// Error logs a message at error Level.
func (logger *Logger) Error(args ...any) {
	logger.log(ErrorLevel, args...)
}

// Fatal logs a message at fatal Level and exits the program.
func (logger *Logger) Fatal(args ...any) {
	logger.log(FatalLevel, args...)
}

// Panic logs a message at panic Level and panics.
func (logger *Logger) Panic(args ...any) {
	logger.log(PanicLevel, args...)
}

// Formatted logging methods

// Tracef logs a formatted message at trace Level.
func (logger *Logger) Tracef(format string, args ...any) {
	logger.logf(TraceLevel, format, args...)
}

// Debugf logs a formatted message at debug Level.
func (logger *Logger) Debugf(format string, args ...any) {
	logger.logf(DebugLevel, format, args...)
}

// Infof logs a formatted message at info Level.
func (logger *Logger) Infof(format string, args ...any) {
	logger.logf(InfoLevel, format, args...)
}

// Printf logs a formatted message at info Level (alias for Infof).
func (logger *Logger) Printf(format string, args ...any) {
	logger.Infof(format, args...)
}

// Warnf logs a formatted message at warning Level.
func (logger *Logger) Warnf(format string, args ...any) {
	logger.logf(WarnLevel, format, args...)
}

// Warningf logs a formatted message at warning Level (alias for Warnf).
func (logger *Logger) Warningf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// Errorf logs a formatted message at error Level.
func (logger *Logger) Errorf(format string, args ...any) {
	logger.logf(ErrorLevel, format, args...)
}

// Fatalf logs a formatted message at fatal Level and exits the program.
func (logger *Logger) Fatalf(format string, args ...any) {
	logger.logf(FatalLevel, format, args...)
}

// Panicf logs a formatted message at panic Level and panics.
func (logger *Logger) Panicf(format string, args ...any) {
	logger.logf(PanicLevel, format, args...)
}

// Line logging methods

// Traceln logs a message at trace Level with newline handling.
func (logger *Logger) Traceln(args ...any) {
	logger.logln(TraceLevel, args...)
}

// Debugln logs a message at debug Level with newline handling.
func (logger *Logger) Debugln(args ...any) {
	logger.logln(DebugLevel, args...)
}

// Infoln logs a message at info Level with newline handling.
func (logger *Logger) Infoln(args ...any) {
	logger.logln(InfoLevel, args...)
}

// Println logs a message at info Level with newline handling (alias for Infoln).
func (logger *Logger) Println(args ...any) {
	logger.Infoln(args...)
}

// Warnln logs a message at warning Level with newline handling.
func (logger *Logger) Warnln(args ...any) {
	logger.logln(WarnLevel, args...)
}

// Warningln logs a message at warning Level with newline handling (alias for Warnln).
func (logger *Logger) Warningln(args ...any) {
	logger.Warnln(args...)
}

// Errorln logs a message at error Level with newline handling.
func (logger *Logger) Errorln(args ...any) {
	logger.logln(ErrorLevel, args...)
}

// Fatalln logs a message at fatal Level with newline handling and exits the program.
func (logger *Logger) Fatalln(args ...any) {
	logger.logln(FatalLevel, args...)
}

// Panicln logs a message at panic Level with newline handling and panics.
func (logger *Logger) Panicln(args ...any) {
	logger.logln(PanicLevel, args...)
}

// Writer returns an io.Writer that writes to the logger at the info log Level.
func (logger *Logger) Writer() *io.PipeWriter {
	return logger.WriterLevel(InfoLevel)
}

// WriterLevel returns an io.Writer that writes to the logger at the given log Level.
func (logger *Logger) WriterLevel(level Level) *io.PipeWriter {
	return NewEntry(logger).WriterLevel(level)
}
