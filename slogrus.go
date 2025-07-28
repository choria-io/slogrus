package logrus

import (
	"io"
	"log/slog"
)

// Level represents the Level of severity for log events.
// It is compatible with logrus.Level.
type Level uint32

// These are the different logging levels in order of decreasing severity.
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// AllLevels represents all available log levels.
var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
	TraceLevel,
}

// String returns the string representation of the log Level.
func (level Level) String() string {
	switch level {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}
	return "unknown"
}

// toSlogLevel converts our Level to slog.Level
func (level Level) toSlogLevel() slog.Level {
	switch level {
	case TraceLevel:
		return slog.LevelDebug - 4
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case FatalLevel:
		return slog.LevelError + 4
	case PanicLevel:
		return slog.LevelError + 8
	}
	return slog.LevelInfo
}

// ParseLevel parses a Level string into a Level value.
func ParseLevel(lvl string) (Level, error) {
	switch lvl {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "trace":
		return TraceLevel, nil
	default:
		return InfoLevel, &ParseError{msg: "not a valid logrus Level: \"" + lvl + "\""}
	}
}

// ParseError represents an error encountered during Level parsing.
type ParseError struct {
	msg string
}

func (e *ParseError) Error() string {
	return e.msg
}

// Fields represents a map of key-value pairs for structured logging.
type Fields map[string]any

// standardLogger is the default logger instance, similar to logrus's global logger.
var standardLogger = New()

// StandardLogger returns the standard logger.
func StandardLogger() *Logger {
	return standardLogger
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(out io.Writer) {
	standardLogger.SetOutput(out)
}

// SetLevel sets the logging Level for the standard logger.
func SetLevel(level Level) {
	standardLogger.SetLevel(level)
}

// WithField creates an entry with a single field using the standard logger.
func WithField(key string, value any) *Entry {
	return standardLogger.WithField(key, value)
}

// WithFields creates an entry with multiple fields using the standard logger.
func WithFields(fields Fields) *Entry {
	return standardLogger.WithFields(fields)
}

// WithError creates an entry with an error field using the standard logger.
func WithError(err error) *Entry {
	return standardLogger.WithError(err)
}

// Global logging functions

// Trace logs a message at trace Level using the standard logger.
func Trace(args ...any) {
	standardLogger.Trace(args...)
}

// Debug logs a message at debug Level using the standard logger.
func Debug(args ...any) {
	standardLogger.Debug(args...)
}

// Info logs a message at info Level using the standard logger.
func Info(args ...any) {
	standardLogger.Info(args...)
}

// Print logs a message at info Level using the standard logger (alias for Info).
func Print(args ...any) {
	standardLogger.Print(args...)
}

// Warn logs a message at warning Level using the standard logger.
func Warn(args ...any) {
	standardLogger.Warn(args...)
}

// Warning logs a message at warning Level using the standard logger (alias for Warn).
func Warning(args ...any) {
	standardLogger.Warning(args...)
}

// Error logs a message at error Level using the standard logger.
func Error(args ...any) {
	standardLogger.Error(args...)
}

// Fatal logs a message at fatal Level using the standard logger and exits the program.
func Fatal(args ...any) {
	standardLogger.Fatal(args...)
}

// Panic logs a message at panic Level using the standard logger and panics.
func Panic(args ...any) {
	standardLogger.Panic(args...)
}

// Formatted global logging functions

// Tracef logs a formatted message at trace Level using the standard logger.
func Tracef(format string, args ...any) {
	standardLogger.Tracef(format, args...)
}

// Debugf logs a formatted message at debug Level using the standard logger.
func Debugf(format string, args ...any) {
	standardLogger.Debugf(format, args...)
}

// Infof logs a formatted message at info Level using the standard logger.
func Infof(format string, args ...any) {
	standardLogger.Infof(format, args...)
}

// Printf logs a formatted message at info Level using the standard logger (alias for Infof).
func Printf(format string, args ...any) {
	standardLogger.Printf(format, args...)
}

// Warnf logs a formatted message at warning Level using the standard logger.
func Warnf(format string, args ...any) {
	standardLogger.Warnf(format, args...)
}

// Warningf logs a formatted message at warning Level using the standard logger (alias for Warnf).
func Warningf(format string, args ...any) {
	standardLogger.Warningf(format, args...)
}

// Errorf logs a formatted message at error Level using the standard logger.
func Errorf(format string, args ...any) {
	standardLogger.Errorf(format, args...)
}

// Fatalf logs a formatted message at fatal Level using the standard logger and exits the program.
func Fatalf(format string, args ...any) {
	standardLogger.Fatalf(format, args...)
}

// Panicf logs a formatted message at panic Level using the standard logger and panics.
func Panicf(format string, args ...any) {
	standardLogger.Panicf(format, args...)
}

// Line global logging functions

// Traceln logs a message at trace Level using the standard logger with newline handling.
func Traceln(args ...any) {
	standardLogger.Traceln(args...)
}

// Debugln logs a message at debug Level using the standard logger with newline handling.
func Debugln(args ...any) {
	standardLogger.Debugln(args...)
}

// Infoln logs a message at info Level using the standard logger with newline handling.
func Infoln(args ...any) {
	standardLogger.Infoln(args...)
}

// Println logs a message at info Level using the standard logger with newline handling (alias for Infoln).
func Println(args ...any) {
	standardLogger.Println(args...)
}

// Warnln logs a message at warning Level using the standard logger with newline handling.
func Warnln(args ...any) {
	standardLogger.Warnln(args...)
}

// Warningln logs a message at warning Level using the standard logger with newline handling (alias for Warnln).
func Warningln(args ...any) {
	standardLogger.Warningln(args...)
}

// Errorln logs a message at error Level using the standard logger with newline handling.
func Errorln(args ...any) {
	standardLogger.Errorln(args...)
}

// Fatalln logs a message at fatal Level using the standard logger with newline handling and exits the program.
func Fatalln(args ...any) {
	standardLogger.Fatalln(args...)
}

// Panicln logs a message at panic Level using the standard logger with newline handling and panics.
func Panicln(args ...any) {
	standardLogger.Panicln(args...)
}
