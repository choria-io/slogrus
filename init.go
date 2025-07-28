package logrus

import (
	"io"
	"log/slog"
	"os"
)

// NewTextLogger creates a new Logger with a text handler.
func NewTextLogger(w io.Writer, opts *slog.HandlerOptions) *Logger {
	if w == nil {
		w = os.Stderr
	}
	if opts == nil {
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	handler := slog.NewTextHandler(w, opts)

	// Determine our internal Level based on slog handler Level
	var internalLevel Level = InfoLevel
	if opts.Level != nil {
		switch {
		case opts.Level.Level() <= slog.LevelDebug-4:
			internalLevel = TraceLevel
		case opts.Level.Level() <= slog.LevelDebug:
			internalLevel = DebugLevel
		case opts.Level.Level() <= slog.LevelInfo:
			internalLevel = InfoLevel
		case opts.Level.Level() <= slog.LevelWarn:
			internalLevel = WarnLevel
		case opts.Level.Level() <= slog.LevelError:
			internalLevel = ErrorLevel
		case opts.Level.Level() <= slog.LevelError+4:
			internalLevel = FatalLevel
		default:
			internalLevel = PanicLevel
		}
	}

	return &Logger{
		slogger:   slog.New(handler),
		Level:     internalLevel,
		Out:       w,
		Formatter: &TextFormatter{},
	}
}

// NewJSONLogger creates a new Logger with a JSON handler.
func NewJSONLogger(w io.Writer, opts *slog.HandlerOptions) *Logger {
	if w == nil {
		w = os.Stderr
	}
	if opts == nil {
		opts = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	handler := slog.NewJSONHandler(w, opts)

	// Determine our internal Level based on slog handler Level
	var internalLevel Level = InfoLevel
	if opts.Level != nil {
		switch {
		case opts.Level.Level() <= slog.LevelDebug-4:
			internalLevel = TraceLevel
		case opts.Level.Level() <= slog.LevelDebug:
			internalLevel = DebugLevel
		case opts.Level.Level() <= slog.LevelInfo:
			internalLevel = InfoLevel
		case opts.Level.Level() <= slog.LevelWarn:
			internalLevel = WarnLevel
		case opts.Level.Level() <= slog.LevelError:
			internalLevel = ErrorLevel
		case opts.Level.Level() <= slog.LevelError+4:
			internalLevel = FatalLevel
		default:
			internalLevel = PanicLevel
		}
	}

	return &Logger{
		slogger:   slog.New(handler),
		Level:     internalLevel,
		Out:       w,
		Formatter: &JSONFormatter{},
	}
}

// SetFormatter is a compatibility function for logrus that allows switching between text and JSON formatters.
// It recreates the standard logger with the appropriate handler.
func SetFormatter(formatter Formatter) {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: standardLogger.Level.toSlogLevel(),
	}

	switch formatter.(type) {
	case *TextFormatter:
		handler = slog.NewTextHandler(standardLogger.Out, opts)
		standardLogger.Formatter = formatter
	case *JSONFormatter:
		handler = slog.NewJSONHandler(standardLogger.Out, opts)
		standardLogger.Formatter = formatter
	default:
		// Default to text handler
		handler = slog.NewTextHandler(standardLogger.Out, opts)
		standardLogger.Formatter = &TextFormatter{}
	}

	standardLogger.slogger = slog.New(handler)
}

// Formatter interface for logrus compatibility.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// TextFormatter provides a logrus-compatible text formatter.
type TextFormatter struct {
	// DisableColors allows disabling colors in output.
	DisableColors bool
	// FullTimestamp enables full timestamp instead of just time.
	FullTimestamp bool
	// ForceColors forces colored output even when not in a TTY.
	ForceColors bool
}

// Format formats the entry as text (placeholder implementation).
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	// This is a placeholder - the actual formatting is handled by slog
	return []byte{}, nil
}

// JSONFormatter provides a logrus-compatible JSON formatter.
type JSONFormatter struct {
	// DisableTimestamp disables automatic timestamp field.
	DisableTimestamp bool
	// DisableHTMLEscape disables HTML escaping.
	DisableHTMLEscape bool
}

// Format formats the entry as JSON (placeholder implementation).
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	// This is a placeholder - the actual formatting is handled by slog
	return []byte{}, nil
}

// SetReportCaller enables or disables caller reporting for the standard logger.
func SetReportCaller(include bool) {
	// Create new handler options with caller reporting
	opts := &slog.HandlerOptions{
		Level:     standardLogger.Level.toSlogLevel(),
		AddSource: include,
	}

	// Recreate the handler based on current type
	var handler slog.Handler
	if _, ok := standardLogger.slogger.Handler().(*slog.JSONHandler); ok {
		handler = slog.NewJSONHandler(standardLogger.Out, opts)
		standardLogger.Formatter = &JSONFormatter{}
	} else {
		handler = slog.NewTextHandler(standardLogger.Out, opts)
		standardLogger.Formatter = &TextFormatter{}
	}

	standardLogger.slogger = slog.New(handler)
}
