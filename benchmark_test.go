package slogrus

import (
	"bytes"
	"io"
	"log/slog"
	"testing"
)

// Benchmarks for direct logger calls (no fields)
func BenchmarkLoggerInfo(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("test message")
	}
}

func BenchmarkLoggerInfof(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Infof("test message %d", i)
	}
}

func BenchmarkLoggerInfoln(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Infoln("test", "message")
	}
}

// Benchmarks for Entry with single field
func BenchmarkLoggerWithField(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithField("key", "value").Info("test message")
	}
}

func BenchmarkLoggerWithFieldf(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithField("key", "value").Infof("test message %d", i)
	}
}

// Benchmarks for Entry with multiple fields
func BenchmarkLoggerWithFields(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	fields := Fields{
		"component": "benchmark",
		"operation": "test",
		"count":     42,
	}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithFields(fields).Info("test message")
	}
}

// Benchmarks for field chaining
func BenchmarkLoggerWithFieldChaining(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithField("component", "benchmark").
			WithField("operation", "test").
			WithField("count", i).
			Info("test message")
	}
}

// Benchmarks for disabled levels (should be very fast)
func BenchmarkLoggerDebugDisabled(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Debug("debug message that should be filtered")
	}
}

func BenchmarkLoggerWithFieldDebugDisabled(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithField("key", "value").Debug("debug message that should be filtered")
	}
}

// Benchmarks for global functions
func BenchmarkGlobalInfo(b *testing.B) {
	// Replace standard logger for benchmarking
	originalLogger := standardLogger
	standardLogger = NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	defer func() {
		standardLogger = originalLogger
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		Info("test message")
	}
}

func BenchmarkGlobalWithField(b *testing.B) {
	// Replace standard logger for benchmarking
	originalLogger := standardLogger
	standardLogger = NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	defer func() {
		standardLogger = originalLogger
	}()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		WithField("key", "value").Info("test message")
	}
}

// Benchmarks comparing JSON vs Text handlers
func BenchmarkLoggerInfoJSON(b *testing.B) {
	logger := NewJSONLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("test message")
	}
}

func BenchmarkLoggerWithFieldJSON(b *testing.B) {
	logger := NewJSONLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithField("key", "value").Info("test message")
	}
}

// Benchmark for WithError
func BenchmarkLoggerWithError(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	err := &ParseError{msg: "benchmark error"}
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithError(err).Info("test message")
	}
}

// Benchmark for complex logging scenario
func BenchmarkComplexLogging(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.WithFields(Fields{
			"user_id":    12345,
			"session_id": "sess_abc123",
			"request_id": "req_xyz789",
			"method":     "GET",
			"path":       "/api/users",
			"status":     200,
			"duration":   "45ms",
		}).Infof("Request completed: %s %s", "GET", "/api/users")
	}
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocation(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})

	b.Run("DirectLog", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.Info("simple message")
		}
	})

	b.Run("WithOneField", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.WithField("key", "value").Info("message with field")
		}
	})

	b.Run("WithThreeFields", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.WithField("key1", "value1").
				WithField("key2", "value2").
				WithField("key3", "value3").
				Info("message with three fields")
		}
	})

	b.Run("WithFieldsMap", func(b *testing.B) {
		fields := Fields{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			logger.WithFields(fields).Info("message with fields map")
		}
	})
}

// Throughput benchmark
func BenchmarkThroughput(b *testing.B) {
	var buf bytes.Buffer
	logger := NewTextLogger(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})

	b.SetBytes(1) // Roughly 1 log entry per iteration
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithField("worker", "benchmark").Info("throughput test message")
		}
	})
}

// Level checking benchmark
func BenchmarkLevelCheck(b *testing.B) {
	logger := NewTextLogger(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if logger.IsLevelEnabled(DebugLevel) {
			logger.Debug("this should not be called")
		}
	}
}
