package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/otakakot/sample-gorm-slog/internal/contextx"
	"github.com/otakakot/sample-gorm-slog/internal/gormx"
	"github.com/otakakot/sample-gorm-slog/internal/slogx"
)

var tracer = otel.Tracer("gorm")

func init() {
	provider := trace.NewTracerProvider()

	otel.SetTracerProvider(provider)

	log := slog.New(slogx.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))

	slog.SetDefault(log)

	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func main() {
	ctx, span := tracer.Start(context.Background(), "main")
	defer span.End()

	dsn := "postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &gormx.Logger{
			LogLevel:      logger.Info,
			SlowThreshold: 5 * time.Millisecond,
		},
	})
	if err != nil {
		panic("failed to connect database. error: " + err.Error())
	}

	user := &User{}

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		panic("failed to create user. error: " + err.Error())
	}

	ctx = context.WithValue(ctx, contextx.UserIDKey, user.ID)

	found := &User{
		Model: gorm.Model{
			ID: user.ID,
		},
	}

	ctx, span = tracer.Start(ctx, "find")
	defer span.End()

	if err := db.WithContext(ctx).First(found).Error; err != nil {
		panic("failed to find user. error: " + err.Error())
	}

	ctx, span = tracer.Start(ctx, "delete")
	defer span.End()

	if err := db.WithContext(ctx).Delete(found).Error; err != nil {
		panic("failed to delete user. error: " + err.Error())
	}

	ctx, span = tracer.Start(ctx, "unscoped_delete")
	defer span.End()

	if err := db.WithContext(ctx).Unscoped().Delete(found).Error; err != nil {
		panic("failed to unscoped delete user. error: " + err.Error())
	}

	fu := &User{
		Model: gorm.Model{
			ID: user.ID,
		},
	}

	ctx, span = tracer.Start(ctx, "not_found")
	defer span.End()

	db.WithContext(ctx).First(fu)
}

type User struct {
	gorm.Model
}
