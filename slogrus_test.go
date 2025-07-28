package slogrus

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{TraceLevel, "trace"},
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warning"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{PanicLevel, "panic"},
		{Level(99), "unknown"},
	}

	for _, test := range tests {
		if result := test.level.String(); result != test.expected {
			t.Errorf("Level(%d).String() = %q, want %q", test.level, result, test.expected)
		}
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		hasError bool
	}{
		{"panic", PanicLevel, false},
		{"fatal", FatalLevel, false},
		{"error", ErrorLevel, false},
		{"warn", WarnLevel, false},
		{"warning", WarnLevel, false},
		{"info", InfoLevel, false},
		{"debug", DebugLevel, false},
		{"trace", TraceLevel, false},
		{"invalid", InfoLevel, true},
		{"", InfoLevel, true},
	}

	for _, test := range tests {
		result, err := ParseLevel(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("ParseLevel(%q) expected error, got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseLevel(%q) unexpected error: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", test.input, result, test.expected)
			}
		}
	}
}

func TestParseError(t *testing.T) {
	err := &ParseError{msg: "test error"}
	if err.Error() != "test error" {
		t.Errorf("ParseError.Error() = %q, want %q", err.Error(), "test error")
	}
}

func TestAllLevels(t *testing.T) {
	expected := []Level{
		PanicLevel,
		FatalLevel,
		ErrorLevel,
		WarnLevel,
		InfoLevel,
		DebugLevel,
		TraceLevel,
	}

	if len(AllLevels) != len(expected) {
		t.Errorf("AllLevels has %d levels, want %d", len(AllLevels), len(expected))
	}

	for i, level := range AllLevels {
		if level != expected[i] {
			t.Errorf("AllLevels[%d] = %v, want %v", i, level, expected[i])
		}
	}
}

func TestStandardLogger(t *testing.T) {
	logger := StandardLogger()
	if logger == nil {
		t.Error("StandardLogger() returned nil")
	}

	if logger != standardLogger {
		t.Error("StandardLogger() did not return the standard logger instance")
	}
}
