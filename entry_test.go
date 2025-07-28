package slogrus

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewEntry(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)

	if entry == nil {
		t.Error("NewEntry() returned nil")
	}
	if entry.logger != logger {
		t.Error("NewEntry() did not set logger correctly")
	}
	if entry.Data == nil {
		t.Error("NewEntry() did not initialize Data field")
	}
	if entry.Context == nil {
		t.Error("NewEntry() did not initialize Context field")
	}
}

func TestEntryWithField(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)
	entry2 := entry.WithField("key", "value")

	if entry2 == nil {
		t.Error("WithField() returned nil")
	}
	if entry2 == entry {
		t.Error("WithField() should return a new entry, not modify the original")
	}
	if len(entry2.Data) != 1 {
		t.Errorf("WithField() entry has %d fields, want 1", len(entry2.Data))
	}
	if entry2.Data["key"] != "value" {
		t.Errorf("WithField() entry.Data[\"key\"] = %v, want \"value\"", entry2.Data["key"])
	}
}

func TestEntryWithFields(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)
	fields := Fields{"key1": "value1", "key2": "value2"}
	entry2 := entry.WithFields(fields)

	if entry2 == nil {
		t.Error("WithFields() returned nil")
	}
	if entry2 == entry {
		t.Error("WithFields() should return a new entry, not modify the original")
	}
	if len(entry2.Data) != 2 {
		t.Errorf("WithFields() entry has %d fields, want 2", len(entry2.Data))
	}
	if entry2.Data["key1"] != "value1" {
		t.Errorf("WithFields() entry.Data[\"key1\"] = %v, want \"value1\"", entry2.Data["key1"])
	}
	if entry2.Data["key2"] != "value2" {
		t.Errorf("WithFields() entry.Data[\"key2\"] = %v, want \"value2\"", entry2.Data["key2"])
	}
}

func TestEntryWithContext(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)
	ctx := context.WithValue(context.Background(), "test", "value")
	entry2 := entry.WithContext(ctx)

	if entry2 == nil {
		t.Error("WithContext() returned nil")
	}
	if entry2 == entry {
		t.Error("WithContext() should return a new entry, not modify the original")
	}
	if entry2.Context != ctx {
		t.Error("WithContext() did not set context correctly")
	}
}

func TestEntryWithError(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)
	err := &ParseError{msg: "test error"}
	entry2 := entry.WithError(err)

	if entry2 == nil {
		t.Error("WithError() returned nil")
	}
	if entry2 == entry {
		t.Error("WithError() should return a new entry, not modify the original")
	}
	if len(entry2.Data) != 1 {
		t.Errorf("WithError() entry has %d fields, want 1", len(entry2.Data))
	}
	if entry2.Data["error"] != err {
		t.Errorf("WithError() entry.Data[\"error\"] = %v, want %v", entry2.Data["error"], err)
	}
}

func TestEntryWithTime(t *testing.T) {
	logger := New()
	entry := NewEntry(logger)
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	entry2 := entry.WithTime(testTime)

	if entry2 == nil {
		t.Error("WithTime() returned nil")
	}
	if entry2 == entry {
		t.Error("WithTime() should return a new entry, not modify the original")
	}
	if entry2.Time != testTime {
		t.Errorf("WithTime() entry.Time = %v, want %v", entry2.Time, testTime)
	}
}

func TestEntryLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug - 4})
	entry := logger.WithField("component", "test")

	// Test different logging levels
	entry.Trace("trace message")
	entry.Debug("debug message")
	entry.Info("info message")
	entry.Warn("warn message")
	entry.Error("error message")

	output := buf.String()

	// Check that messages were logged with field
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

	// Check that the field is present in the output
	if !strings.Contains(output, "component=test") {
		t.Error("Field 'component=test' not found in output")
	}
}

func TestEntryFormattedLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	entry := logger.WithField("component", "test")

	entry.Debugf("debug message %d", 42)
	entry.Infof("info message %s", "formatted")

	output := buf.String()

	if !strings.Contains(output, "debug message 42") {
		t.Error("Formatted debug message not found in output")
	}
	if !strings.Contains(output, "info message formatted") {
		t.Error("Formatted info message not found in output")
	}
	if !strings.Contains(output, "component=test") {
		t.Error("Field 'component=test' not found in formatted output")
	}
}

func TestEntryChaining(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	// Test method chaining
	logger.WithField("component", "test").
		WithField("operation", "chain").
		Info("chained message")

	output := buf.String()

	if !strings.Contains(output, "chained message") {
		t.Error("Chained message not found in output")
	}
	if !strings.Contains(output, "component=test") {
		t.Error("Field 'component=test' not found in chained output")
	}
	if !strings.Contains(output, "operation=chain") {
		t.Error("Field 'operation=chain' not found in chained output")
	}
}

func TestEntryLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	logger.SetLevel(WarnLevel)
	entry := logger.WithField("component", "test")

	entry.Debug("debug message")
	entry.Info("info message")
	entry.Warn("warn message")
	entry.Error("error message")

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
