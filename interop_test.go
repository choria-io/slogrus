package logrus

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"testing"
)

// TestFromSlogLogger tests creating a slogrus logger from an existing slog logger
func TestFromSlogLogger(t *testing.T) {
	var buf bytes.Buffer

	// Create an slog logger
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slogLogger := slog.New(handler)

	// Create slogrus logger from slog logger
	logger := FromSlogLogger(slogLogger)

	// Test that it works
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected log message not found in output: %s", output)
	}

	// Should be JSON format since we used JSONHandler
	if !strings.Contains(output, `"msg":"test message"`) {
		t.Errorf("Expected JSON format not found in output: %s", output)
	}
}

// TestGetSlogLogger tests accessing the underlying slog logger
func TestGetSlogLogger(t *testing.T) {
	var buf bytes.Buffer

	// Create slogrus logger
	logger := NewTextLogger(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Get the underlying slog logger
	slogger := logger.GetSlogLogger()

	// Use slog logger directly
	slogger.Info("direct slog message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "direct slog message") {
		t.Errorf("Expected slog message not found in output: %s", output)
	}

	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected slog attribute not found in output: %s", output)
	}
}

// TestInteroperability tests mixed usage patterns
func TestInteroperability(t *testing.T) {
	var buf bytes.Buffer

	// Start with slog logger
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slogLogger := slog.New(handler)

	// Wrap in slogrus
	logger := FromSlogLogger(slogLogger)

	// Use logrus-style logging
	logger.WithField("component", "test").Info("logrus style message")

	// Get slog logger and use slog-style logging
	underlyingSlog := logger.GetSlogLogger()
	underlyingSlog.Info("slog style message", "service", "api")

	output := buf.String()

	// Check both messages are present
	if !strings.Contains(output, "logrus style message") {
		t.Errorf("Logrus-style message not found in output: %s", output)
	}

	if !strings.Contains(output, "slog style message") {
		t.Errorf("Slog-style message not found in output: %s", output)
	}

	// Check both have their respective attributes
	if !strings.Contains(output, "component=test") {
		t.Errorf("Logrus field not found in output: %s", output)
	}

	if !strings.Contains(output, "service=api") {
		t.Errorf("Slog attribute not found in output: %s", output)
	}
}

// TestFromSlogLoggerLevelHandling tests Level handling with FromSlogLogger
func TestFromSlogLoggerLevelHandling(t *testing.T) {
	var buf bytes.Buffer

	// Create slog logger with Warn Level
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	slogLogger := slog.New(handler)

	// Create slogrus logger
	logger := FromSlogLogger(slogLogger)

	// Debug should be filtered out by slog handler
	logger.Debug("debug message")

	// Warn should pass through
	logger.Warn("warn message")

	output := buf.String()

	// Debug should not appear
	if strings.Contains(output, "debug message") {
		t.Errorf("Debug message should have been filtered out: %s", output)
	}

	// Warn should appear
	if !strings.Contains(output, "warn message") {
		t.Errorf("Warn message should be present: %s", output)
	}
}

// BenchmarkFromSlogLogger benchmarks FromSlogLogger performance
func BenchmarkFromSlogLogger(b *testing.B) {
	handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	slogLogger := slog.New(handler)
	logger := FromSlogLogger(slogLogger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

// BenchmarkGetSlogLogger benchmarks GetSlogLogger performance
func BenchmarkGetSlogLogger(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	slogger := logger.GetSlogLogger()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		slogger.Info("benchmark message")
	}
}
