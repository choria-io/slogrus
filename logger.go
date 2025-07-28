package slogrus

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
	level   Level
	out     io.Writer
}

// New creates a new Logger instance with default text handler.
func New() *Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return &Logger{
		slogger: slog.New(handler),
		level:   InfoLevel,
		out:     os.Stderr,
	}
}

// NewWithHandler creates a new Logger with a custom slog.Handler.
func NewWithHandler(handler slog.Handler) *Logger {
	return &Logger{
		slogger: slog.New(handler),
		level:   InfoLevel,
		out:     os.Stderr,
	}
}

// FromSlogLogger creates a new Logger instance from an existing slog.Logger.
// This enables interoperability with existing slog-based code.
func FromSlogLogger(slogger *slog.Logger) *Logger {
	return &Logger{
		slogger: slogger,
		level:   InfoLevel, // Default level, can be changed with SetLevel
		out:     os.Stderr, // Default output, may not match slog handler's output
	}
}

// SetOutput sets the output destination for the logger.
func (logger *Logger) SetOutput(out io.Writer) {
	logger.out = out
	// Create a new handler with the new output
	opts := &slog.HandlerOptions{
		Level: logger.level.toSlogLevel(),
	}

	if _, ok := logger.slogger.Handler().(*slog.TextHandler); ok {
		logger.slogger = slog.New(slog.NewTextHandler(out, opts))
	} else if _, ok := logger.slogger.Handler().(*slog.JSONHandler); ok {
		logger.slogger = slog.New(slog.NewJSONHandler(out, opts))
	}
}

// SetLevel sets the logging level for the logger.
func (logger *Logger) SetLevel(level Level) {
	logger.level = level
	// Update the slog handler with new level
	opts := &slog.HandlerOptions{
		Level: level.toSlogLevel(),
	}

	// Recreate handler with new level
	if _, ok := logger.slogger.Handler().(*slog.TextHandler); ok {
		logger.slogger = slog.New(slog.NewTextHandler(logger.out, opts))
	} else if _, ok := logger.slogger.Handler().(*slog.JSONHandler); ok {
		logger.slogger = slog.New(slog.NewJSONHandler(logger.out, opts))
	}
}

// GetLevel returns the current logging level.
func (logger *Logger) GetLevel() Level {
	return logger.level
}

// IsLevelEnabled checks if the given level is enabled for logging.
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return level <= logger.level
}

// GetSlogLogger returns the underlying slog.Logger instance.
// This enables advanced slog operations and direct access to slog APIs.
func (logger *Logger) GetSlogLogger() *slog.Logger {
	return logger.slogger
}

// WithField creates an entry with a single field.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
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
func (logger *Logger) log(level Level, args ...interface{}) {
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
func (logger *Logger) logf(level Level, format string, args ...interface{}) {
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
func (logger *Logger) logln(level Level, args ...interface{}) {
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

// Trace logs a message at trace level.
func (logger *Logger) Trace(args ...interface{}) {
	logger.log(TraceLevel, args...)
}

// Debug logs a message at debug level.
func (logger *Logger) Debug(args ...interface{}) {
	logger.log(DebugLevel, args...)
}

// Info logs a message at info level.
func (logger *Logger) Info(args ...interface{}) {
	logger.log(InfoLevel, args...)
}

// Print logs a message at info level (alias for Info).
func (logger *Logger) Print(args ...interface{}) {
	logger.Info(args...)
}

// Warn logs a message at warning level.
func (logger *Logger) Warn(args ...interface{}) {
	logger.log(WarnLevel, args...)
}

// Warning logs a message at warning level (alias for Warn).
func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

// Error logs a message at error level.
func (logger *Logger) Error(args ...interface{}) {
	logger.log(ErrorLevel, args...)
}

// Fatal logs a message at fatal level and exits the program.
func (logger *Logger) Fatal(args ...interface{}) {
	logger.log(FatalLevel, args...)
}

// Panic logs a message at panic level and panics.
func (logger *Logger) Panic(args ...interface{}) {
	logger.log(PanicLevel, args...)
}

// Formatted logging methods

// Tracef logs a formatted message at trace level.
func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.logf(TraceLevel, format, args...)
}

// Debugf logs a formatted message at debug level.
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logf(DebugLevel, format, args...)
}

// Infof logs a formatted message at info level.
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logf(InfoLevel, format, args...)
}

// Printf logs a formatted message at info level (alias for Infof).
func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warnf logs a formatted message at warning level.
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logf(WarnLevel, format, args...)
}

// Warningf logs a formatted message at warning level (alias for Warnf).
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Errorf logs a formatted message at error level.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logf(ErrorLevel, format, args...)
}

// Fatalf logs a formatted message at fatal level and exits the program.
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logf(FatalLevel, format, args...)
}

// Panicf logs a formatted message at panic level and panics.
func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.logf(PanicLevel, format, args...)
}

// Line logging methods

// Traceln logs a message at trace level with newline handling.
func (logger *Logger) Traceln(args ...interface{}) {
	logger.logln(TraceLevel, args...)
}

// Debugln logs a message at debug level with newline handling.
func (logger *Logger) Debugln(args ...interface{}) {
	logger.logln(DebugLevel, args...)
}

// Infoln logs a message at info level with newline handling.
func (logger *Logger) Infoln(args ...interface{}) {
	logger.logln(InfoLevel, args...)
}

// Println logs a message at info level with newline handling (alias for Infoln).
func (logger *Logger) Println(args ...interface{}) {
	logger.Infoln(args...)
}

// Warnln logs a message at warning level with newline handling.
func (logger *Logger) Warnln(args ...interface{}) {
	logger.logln(WarnLevel, args...)
}

// Warningln logs a message at warning level with newline handling (alias for Warnln).
func (logger *Logger) Warningln(args ...interface{}) {
	logger.Warnln(args...)
}

// Errorln logs a message at error level with newline handling.
func (logger *Logger) Errorln(args ...interface{}) {
	logger.logln(ErrorLevel, args...)
}

// Fatalln logs a message at fatal level with newline handling and exits the program.
func (logger *Logger) Fatalln(args ...interface{}) {
	logger.logln(FatalLevel, args...)
}

// Panicln logs a message at panic level with newline handling and panics.
func (logger *Logger) Panicln(args ...interface{}) {
	logger.logln(PanicLevel, args...)
}
