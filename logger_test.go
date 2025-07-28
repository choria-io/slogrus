package slogrus

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Error("New() returned nil")
		return
	}
	if logger.level != InfoLevel {
		t.Errorf("New() logger level = %v, want %v", logger.level, InfoLevel)
	}
}

func TestNewWithHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := NewWithHandler(handler)

	if logger == nil {
		t.Error("NewWithHandler() returned nil")
		return
	}
	if logger.level != InfoLevel {
		t.Errorf("NewWithHandler() logger level = %v, want %v", logger.level, InfoLevel)
	}
}

func TestLoggerSetLevel(t *testing.T) {
	logger := New()
	logger.SetLevel(DebugLevel)

	if logger.GetLevel() != DebugLevel {
		t.Errorf("GetLevel() = %v, want %v", logger.GetLevel(), DebugLevel)
	}
}

func TestIsLevelEnabled(t *testing.T) {
	logger := New()
	logger.SetLevel(WarnLevel)

	tests := []struct {
		level   Level
		enabled bool
	}{
		{PanicLevel, true},
		{FatalLevel, true},
		{ErrorLevel, true},
		{WarnLevel, true},
		{InfoLevel, false},
		{DebugLevel, false},
		{TraceLevel, false},
	}

	for _, test := range tests {
		if result := logger.IsLevelEnabled(test.level); result != test.enabled {
			t.Errorf("IsLevelEnabled(%v) = %v, want %v", test.level, result, test.enabled)
		}
	}
}

func TestWithField(t *testing.T) {
	logger := New()
	entry := logger.WithField("key", "value")

	if entry == nil {
		t.Error("WithField() returned nil")
		return
	}
	if len(entry.Data) != 1 {
		t.Errorf("WithField() entry has %d fields, want 1", len(entry.Data))
	}
	if entry.Data["key"] != "value" {
		t.Errorf("WithField() entry.Data[\"key\"] = %v, want \"value\"", entry.Data["key"])
	}
}

func TestWithFields(t *testing.T) {
	logger := New()
	fields := Fields{"key1": "value1", "key2": "value2"}
	entry := logger.WithFields(fields)

	if entry == nil {
		t.Error("WithFields() returned nil")
		return
	}
	if len(entry.Data) != 2 {
		t.Errorf("WithFields() entry has %d fields, want 2", len(entry.Data))
	}
	if entry.Data["key1"] != "value1" {
		t.Errorf("WithFields() entry.Data[\"key1\"] = %v, want \"value1\"", entry.Data["key1"])
	}
	if entry.Data["key2"] != "value2" {
		t.Errorf("WithFields() entry.Data[\"key2\"] = %v, want \"value2\"", entry.Data["key2"])
	}
}

func TestWithContext(t *testing.T) {
	logger := New()
	type testKeyType string
	const testKey testKeyType = "test"
	ctx := context.WithValue(context.Background(), testKey, "value")
	entry := logger.WithContext(ctx)

	if entry == nil {
		t.Error("WithContext() returned nil")
		return
	}
	if entry.Context != ctx {
		t.Error("WithContext() did not set context correctly")
	}
}

func TestWithError(t *testing.T) {
	logger := New()
	err := &ParseError{msg: "test error"}
	entry := logger.WithError(err)

	if entry == nil {
		t.Error("WithError() returned nil")
		return
	}
	if len(entry.Data) != 1 {
		t.Errorf("WithError() entry has %d fields, want 1", len(entry.Data))
	}
	if entry.Data["error"] != err {
		t.Errorf("WithError() entry.Data[\"error\"] = %v, want %v", entry.Data["error"], err)
	}
}

func TestLoggerBasicLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug - 4})

	// Test different logging levels
	logger.Trace("trace message")
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	// Check that messages were logged
	if !strings.Contains(output, "trace message") {
		t.Error("Trace message not found in output")
	}
	if !strings.Contains(output, "debug message") {
		t.Error("Debug message not found in output")
	}
	if !strings.Contains(output, "info message") {
		t.Error("Info message not found in output")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message not found in output")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message not found in output")
	}
}

func TestLoggerFormattedLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	logger.Debugf("debug message %d", 42)
	logger.Infof("info message %s", "test")

	output := buf.String()

	if !strings.Contains(output, "debug message 42") {
		t.Error("Formatted debug message not found in output")
	}
	if !strings.Contains(output, "info message test") {
		t.Error("Formatted info message not found in output")
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	logger.SetLevel(WarnLevel)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	// Debug and Info should be filtered out
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should be filtered out")
	}
	if strings.Contains(output, "info message") {
		t.Error("Info message should be filtered out")
	}

	// Warn and Error should be present
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message not found in output")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message not found in output")
	}
}
