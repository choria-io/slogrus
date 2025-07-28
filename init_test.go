package logrus

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewTextLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, nil)

	if logger == nil {
		t.Error("NewTextLogger() returned nil")
		return
	}
	if logger.Out != &buf {
		t.Error("NewTextLogger() did not set output correctly")
	}
}

func TestNewJSONLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewJSONLogger(&buf, nil)

	if logger == nil {
		t.Error("NewJSONLogger() returned nil")
		return
	}
	if logger.Out != &buf {
		t.Error("NewJSONLogger() did not set output correctly")
	}

	// Test that it produces JSON output
	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, `"msg":"test message"`) {
		t.Error("JSON logger did not produce JSON output")
	}
}

func TestSetFormatter(t *testing.T) {
	// Test setting text formatter
	SetFormatter(&TextFormatter{})

	// Test setting JSON formatter
	SetFormatter(&JSONFormatter{})

	// Test setting unknown formatter (should default to text)
	SetFormatter(nil)
}

func TestTextFormatter(t *testing.T) {
	formatter := &TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		ForceColors:   false,
	}

	logger := New()
	entry := NewEntry(logger)

	// Test that Format method exists and doesn't error
	_, err := formatter.Format(entry)
	if err != nil {
		t.Errorf("TextFormatter.Format() returned error: %v", err)
	}
}

func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{
		DisableTimestamp:  false,
		DisableHTMLEscape: true,
	}

	logger := New()
	entry := NewEntry(logger)

	// Test that Format method exists and doesn't error
	_, err := formatter.Format(entry)
	if err != nil {
		t.Errorf("JSONFormatter.Format() returned error: %v", err)
	}
}

func TestSetReportCaller(t *testing.T) {
	// Test enabling caller reporting
	SetReportCaller(true)

	// Test disabling caller reporting
	SetReportCaller(false)
}

func TestGlobalFunctions(t *testing.T) {
	var buf bytes.Buffer

	// Replace standard logger for testing
	originalLogger := standardLogger
	standardLogger = NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug - 4})
	defer func() {
		standardLogger = originalLogger
	}()

	// Test global logging functions
	Trace("global trace")
	Debug("global debug")
	Info("global info")
	Warn("global warn")
	Error("global error")

	output := buf.String()

	if !strings.Contains(output, "global trace") {
		t.Error("Global Trace() message not found in output")
	}
	if !strings.Contains(output, "global debug") {
		t.Error("Global Debug() message not found in output")
	}
	if !strings.Contains(output, "global info") {
		t.Error("Global Info() message not found in output")
	}
	if !strings.Contains(output, "global warn") {
		t.Error("Global Warn() message not found in output")
	}
	if !strings.Contains(output, "global error") {
		t.Error("Global Error() message not found in output")
	}
}

func TestGlobalFormattedFunctions(t *testing.T) {
	var buf bytes.Buffer

	// Replace standard logger for testing
	originalLogger := standardLogger
	standardLogger = NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	defer func() {
		standardLogger = originalLogger
	}()

	// Test global formatted logging functions
	Debugf("global debug %d", 42)
	Infof("global info %s", "formatted")

	output := buf.String()

	if !strings.Contains(output, "global debug 42") {
		t.Error("Global Debugf() message not found in output")
	}
	if !strings.Contains(output, "global info formatted") {
		t.Error("Global Infof() message not found in output")
	}
}

func TestGlobalWithFunctions(t *testing.T) {
	var buf bytes.Buffer

	// Replace standard logger for testing
	originalLogger := standardLogger
	standardLogger = NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	defer func() {
		standardLogger = originalLogger
	}()

	// Test global With* functions
	WithField("key", "value").Info("with field")
	WithFields(Fields{"key1": "value1", "key2": "value2"}).Info("with fields")
	WithError(&ParseError{msg: "test error"}).Info("with error")

	output := buf.String()

	if !strings.Contains(output, "with field") {
		t.Error("WithField() message not found in output")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("WithField() field not found in output")
	}
	if !strings.Contains(output, "with fields") {
		t.Error("WithFields() message not found in output")
	}
	if !strings.Contains(output, "key1=value1") {
		t.Error("WithFields() field key1 not found in output")
	}
	if !strings.Contains(output, "key2=value2") {
		t.Error("WithFields() field key2 not found in output")
	}
	if !strings.Contains(output, "with error") {
		t.Error("WithError() message not found in output")
	}
}

func TestSetOutput(t *testing.T) {
	var buf bytes.Buffer

	// Test setting output
	SetOutput(&buf)

	// Test that output was set correctly
	Info("test output")

	output := buf.String()
	if !strings.Contains(output, "test output") {
		t.Error("SetOutput() did not redirect output correctly")
	}
}

func TestSetLevel(t *testing.T) {
	// Test setting Level
	SetLevel(DebugLevel)

	if StandardLogger().Level != DebugLevel {
		t.Errorf("SetLevel() did not set Level correctly, got %v, want %v",
			StandardLogger().Level, DebugLevel)
	}
}
