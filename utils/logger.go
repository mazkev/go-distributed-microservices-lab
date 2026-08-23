package utils

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// InitLogger menginisialisasi JSON Structured Logger standar Go (log/slog)
func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Output format JSON untuk kompatibilitas dengan Datadog / Grafana Loki / CloudWatch
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

// LogWithContext menyisipkan RequestID ke dalam log record jika tersedia di context
func LogWithContext(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if ctx != nil {
		if reqID, ok := ctx.Value(RequestIDKey).(string); ok && reqID != "" {
			return logger.With("request_id", reqID)
		}
	}
	return logger
}
