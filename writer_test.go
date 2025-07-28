package logrus

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestLoggerWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	writer := logger.Writer()
	if writer == nil {
		t.Error("Writer() returned nil")
		return
	}
	defer writer.Close()

	// Write a test message
	message := "test message from writer"
	_, err := fmt.Fprint(writer, message)
	if err != nil {
		t.Errorf("Error writing to writer: %v", err)
	}

	writer.Close()

	// Give some time for the goroutine to process
	time.Sleep(100 * time.Millisecond)

	output := buf.String()
	if !strings.Contains(output, message) {
		t.Errorf("Expected message '%s' not found in output: %s", message, output)
	}
}

func TestLoggerWriterLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	tests := []struct {
		level    Level
		message  string
		expected bool
	}{
		{DebugLevel, "debug message", true},
		{InfoLevel, "info message", true},
		{WarnLevel, "warn message", true},
		{ErrorLevel, "error message", true},
	}

	for _, test := range tests {
		t.Run(test.level.String(), func(t *testing.T) {
			buf.Reset()

			writer := logger.WriterLevel(test.level)
			if writer == nil {
				t.Error("WriterLevel() returned nil")
				return
			}
			defer writer.Close()

			_, err := fmt.Fprint(writer, test.message)
			if err != nil {
				t.Errorf("Error writing to writer: %v", err)
			}

			writer.Close()

			// Give some time for the goroutine to process
			time.Sleep(100 * time.Millisecond)

			output := buf.String()
			if test.expected && !strings.Contains(output, test.message) {
				t.Errorf("Expected message '%s' not found in output: %s", test.message, output)
			}
		})
	}
}

func TestEntryWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	entry := logger.WithField("component", "test")

	writer := entry.Writer()
	if writer == nil {
		t.Error("Entry Writer() returned nil")
		return
	}
	defer writer.Close()

	message := "entry writer test"
	_, err := fmt.Fprint(writer, message)
	if err != nil {
		t.Errorf("Error writing to entry writer: %v", err)
	}

	writer.Close()

	// Give some time for the goroutine to process
	time.Sleep(100 * time.Millisecond)

	output := buf.String()
	if !strings.Contains(output, message) {
		t.Errorf("Expected message '%s' not found in output: %s", message, output)
	}
	if !strings.Contains(output, "component=test") {
		t.Errorf("Expected field 'component=test' not found in output: %s", output)
	}
}

func TestEntryWriterLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	entry := logger.WithField("service", "writer-test")

	tests := []struct {
		level   Level
		message string
	}{
		{DebugLevel, "debug from entry writer"},
		{InfoLevel, "info from entry writer"},
		{WarnLevel, "warn from entry writer"},
		{ErrorLevel, "error from entry writer"},
	}

	for _, test := range tests {
		t.Run(test.level.String(), func(t *testing.T) {
			buf.Reset()

			writer := entry.WriterLevel(test.level)
			if writer == nil {
				t.Error("Entry WriterLevel() returned nil")
				return
			}
			defer writer.Close()

			_, err := fmt.Fprint(writer, test.message)
			if err != nil {
				t.Errorf("Error writing to entry writer: %v", err)
			}

			writer.Close()

			// Give some time for the goroutine to process
			time.Sleep(100 * time.Millisecond)

			output := buf.String()
			if !strings.Contains(output, test.message) {
				t.Errorf("Expected message '%s' not found in output: %s", test.message, output)
			}
			if !strings.Contains(output, "service=writer-test") {
				t.Errorf("Expected field 'service=writer-test' not found in output: %s", output)
			}
		})
	}
}

func TestWriterMultipleLines(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	writer := logger.Writer()
	if writer == nil {
		t.Error("Writer() returned nil")
		return
	}
	defer writer.Close()

	lines := []string{
		"First line",
		"Second line",
		"Third line",
	}

	for _, line := range lines {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			t.Errorf("Error writing line '%s': %v", line, err)
		}
	}

	writer.Close()

	// Give some time for the goroutine to process
	time.Sleep(100 * time.Millisecond)

	output := buf.String()
	for _, line := range lines {
		if !strings.Contains(output, line) {
			t.Errorf("Expected line '%s' not found in output: %s", line, output)
		}
	}
}

func TestWriterLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	logger.SetLevel(WarnLevel)

	// Test that debug Level messages are filtered out
	debugWriter := logger.WriterLevel(DebugLevel)
	if debugWriter == nil {
		t.Error("WriterLevel(DebugLevel) returned nil")
		return
	}
	defer debugWriter.Close()

	debugMessage := "this should be filtered out"
	_, err := fmt.Fprint(debugWriter, debugMessage)
	if err != nil {
		t.Errorf("Error writing debug message: %v", err)
	}
	debugWriter.Close()

	// Test that warn Level messages are not filtered
	warnWriter := logger.WriterLevel(WarnLevel)
	if warnWriter == nil {
		t.Error("WriterLevel(WarnLevel) returned nil")
		return
	}
	defer warnWriter.Close()

	warnMessage := "this should not be filtered"
	_, err = fmt.Fprint(warnWriter, warnMessage)
	if err != nil {
		t.Errorf("Error writing warn message: %v", err)
	}
	warnWriter.Close()

	// Give some time for the goroutines to process
	time.Sleep(100 * time.Millisecond)

	output := buf.String()

	// Debug message should be filtered out
	if strings.Contains(output, debugMessage) {
		t.Errorf("Debug message should be filtered out, but found in output: %s", output)
	}

	// Warn message should be present
	if !strings.Contains(output, warnMessage) {
		t.Errorf("Warn message should be present but not found in output: %s", output)
	}
}
