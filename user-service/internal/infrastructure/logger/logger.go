package logger

import (
	"context"
	"log/slog"
	"os"
	"user-service/internal/entity"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() entity.Logger {
	return &SlogLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (l *SlogLogger) logWithContext(ctx context.Context, level slog.Level, msg string, args ...interface{}) {
	if userID, ok := ctx.Value("user_id").(int64); ok {
		args = append(args, slog.Int64("userID", userID))
	}

	switch level {
	case slog.LevelDebug:
		l.logger.DebugContext(ctx, msg, args...)
	case slog.LevelInfo:
		l.logger.InfoContext(ctx, msg, args...)
	case slog.LevelWarn:
		l.logger.WarnContext(ctx, msg, args...)
	case slog.LevelError:
		l.logger.ErrorContext(ctx, msg, args...)
	default:
		l.logger.InfoContext(ctx, msg, args...)
	}
}

func (l *SlogLogger) Debug(ctx context.Context, msg string, args ...interface{}) {
	l.logWithContext(ctx, slog.LevelDebug, msg, args...)
}

func (l *SlogLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logWithContext(ctx, slog.LevelInfo, msg, args...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logWithContext(ctx, slog.LevelError, msg, args...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logWithContext(ctx, slog.LevelWarn, msg, args...)
}
