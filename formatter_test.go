package logrus

import (
	"bytes"
	"log/slog"
	"testing"
)

// TestFormatterField tests that the Formatter field is properly set and maintained
func TestFormatterField(t *testing.T) {
	// Test New() creates TextFormatter
	logger := New()
	if _, ok := logger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected New() to create TextFormatter")
	}

	// Test NewTextLogger creates TextFormatter
	var buf bytes.Buffer
	textLogger := NewTextLogger(&buf, nil)
	if _, ok := textLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected NewTextLogger() to create TextFormatter")
	}

	// Test NewJSONLogger creates JSONFormatter
	jsonLogger := NewJSONLogger(&buf, nil)
	if _, ok := jsonLogger.Formatter.(*JSONFormatter); !ok {
		t.Error("Expected NewJSONLogger() to create JSONFormatter")
	}

	// Test NewWithHandler with TextHandler
	textHandler := slog.NewTextHandler(&buf, nil)
	handlerLogger := NewWithHandler(textHandler)
	if _, ok := handlerLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected NewWithHandler(TextHandler) to create TextFormatter")
	}

	// Test NewWithHandler with JSONHandler
	jsonHandler := slog.NewJSONHandler(&buf, nil)
	jsonHandlerLogger := NewWithHandler(jsonHandler)
	if _, ok := jsonHandlerLogger.Formatter.(*JSONFormatter); !ok {
		t.Error("Expected NewWithHandler(JSONHandler) to create JSONFormatter")
	}

	// Test FromSlogLogger with text handler
	slogTextLogger := slog.New(slog.NewTextHandler(&buf, nil))
	fromTextLogger := FromSlogLogger(slogTextLogger)
	if _, ok := fromTextLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected FromSlogLogger(text) to create TextFormatter")
	}

	// Test FromSlogLogger with JSON handler
	slogJSONLogger := slog.New(slog.NewJSONHandler(&buf, nil))
	fromJSONLogger := FromSlogLogger(slogJSONLogger)
	if _, ok := fromJSONLogger.Formatter.(*JSONFormatter); !ok {
		t.Error("Expected FromSlogLogger(JSON) to create JSONFormatter")
	}
}

// TestSetFormatterUpdatesField tests that SetFormatter updates the Formatter field
func TestSetFormatterUpdatesField(t *testing.T) {
	// Replace standard logger for testing
	originalLogger := standardLogger
	standardLogger = New()
	defer func() {
		standardLogger = originalLogger
	}()

	// Test setting TextFormatter
	textFormatter := &TextFormatter{DisableColors: true}
	SetFormatter(textFormatter)
	if standardLogger.Formatter != textFormatter {
		t.Error("Expected SetFormatter to set TextFormatter field")
	}

	// Test setting JSONFormatter
	jsonFormatter := &JSONFormatter{DisableTimestamp: true}
	SetFormatter(jsonFormatter)
	if standardLogger.Formatter != jsonFormatter {
		t.Error("Expected SetFormatter to set JSONFormatter field")
	}

	// Test setting nil formatter defaults to TextFormatter
	SetFormatter(nil)
	if _, ok := standardLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected SetFormatter(nil) to default to TextFormatter")
	}
}

// TestFormatterFieldMaintainedDuringOperations tests that Formatter field is maintained during SetOutput/SetLevel
func TestFormatterFieldMaintainedDuringOperations(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	// Test with TextLogger
	textLogger := NewTextLogger(&buf1, nil)

	// Verify initial formatter
	if _, ok := textLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected initial TextFormatter")
	}

	// Test SetOutput maintains formatter
	textLogger.SetOutput(&buf2)
	if _, ok := textLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected SetOutput to maintain TextFormatter")
	}

	// Test SetLevel maintains formatter
	textLogger.SetLevel(ErrorLevel)
	if _, ok := textLogger.Formatter.(*TextFormatter); !ok {
		t.Error("Expected SetLevel to maintain TextFormatter")
	}

	// Test with JSONLogger
	jsonLogger := NewJSONLogger(&buf1, nil)

	// Verify initial formatter
	if _, ok := jsonLogger.Formatter.(*JSONFormatter); !ok {
		t.Error("Expected initial JSONFormatter")
	}

	// Test SetOutput maintains formatter
	jsonLogger.SetOutput(&buf2)
	if _, ok := jsonLogger.Formatter.(*JSONFormatter); !ok {
		t.Error("Expected SetOutput to maintain JSONFormatter")
	}

	// Test SetLevel maintains formatter
	jsonLogger.SetLevel(ErrorLevel)
	if _, ok := jsonLogger.Formatter.(*JSONFormatter); !ok {
		t.Error("Expected SetLevel to maintain JSONFormatter")
	}
}
