package gormx

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*Logger)(nil)

type Logger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

// LogMode implements logger.Interface.
func (log *Logger) LogMode(level logger.LogLevel) logger.Interface {
	log.LogLevel = level

	return log
}

// Info implements logger.Interface.
func (log *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if log.LogLevel >= logger.Info {
		slog.InfoContext(ctx, fmt.Sprintf(msg, data...))
	}
}

// Warn implements logger.Interface.
func (log *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if log.LogLevel >= logger.Warn {
		slog.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

// Error implements logger.Interface.
func (log *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if log.LogLevel >= logger.Error {
		slog.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

// Trace implements logger.Interface.
func (log *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if log.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)

	sql, rows := fc()

	switch {
	case err != nil && log.LogLevel >= logger.Error:
		if rows == -1 {
			slog.ErrorContext(ctx, fmt.Sprintf("[%.3fms] [rows:%s] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql), slog.String("error", err.Error()))
		} else {
			slog.ErrorContext(ctx, fmt.Sprintf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql), slog.String("error", err.Error()))
		}
	case elapsed > log.SlowThreshold && log.SlowThreshold != 0 && log.LogLevel >= logger.Warn:
		if rows == -1 {
			slog.WarnContext(ctx, fmt.Sprintf("[%.3fms] [rows:%s] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			slog.WarnContext(ctx, fmt.Sprintf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	case log.LogLevel == logger.Info:
		if rows == -1 {
			slog.InfoContext(ctx, fmt.Sprintf("[%.3fms] [rows:%s] %s", float64(elapsed.Nanoseconds())/1e6, "-", sql))
		} else {
			slog.InfoContext(ctx, fmt.Sprintf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql))
		}
	}
}
